package release

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ExploratoryEngineering/reto/pkg/gitutil"
	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
)

// Build builds a new release from the current setup
func Build(tagVersion, commitNewRelease bool, overrideName, overrideEmail string) error {
	ctx, err := GetContext()
	if err != nil {
		return err
	}
	if ctx.Released {
		toolbox.PrintError("This version is already released. Bump the version and try again.")
		return errors.New("already released")
	}

	// Do a quick sanity check on the change log
	if err := ChangelogComplete(true); err != nil {
		return err
	}

	if overrideName != "" {
		ctx.Config.CommitterName = overrideName
	}
	if overrideEmail != "" {
		ctx.Config.CommitterEmail = overrideEmail
	}

	if err := VerifyConfig(ctx.Config, true); err != nil {
		return err
	}
	if gitutil.HasChanges(ctx.Config.SourceRoot) {
		toolbox.PrintError("There are uncommitted or unstaged changes in the current Git branch")
		return errors.New("uncommitted changes")
	}

	if !NewFileVersions(ctx.Config) {
		return errors.New("no file changes")
	}

	if err := os.Mkdir(fmt.Sprintf("release/%s", ctx.Version), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create release directory: %v", err)
		return err
	}

	// Tag the commit with the new version
	if tagVersion {
		if err := gitutil.TagVersion(
			ctx.Config.SourceRoot,
			ctx.Config.CommitterName,
			ctx.Config.CommitterEmail,
			fmt.Sprintf("v%s", ctx.Version),
			fmt.Sprintf("Release v%s (%s)", ctx.Version, ctx.Name)); err != nil {
			return err
		}
	}

	if err := releaseChangelog(ctx); err != nil {
		return err
	}

	if err := os.Remove(WorkingChangelog); err != nil {
		toolbox.PrintError("Could not remove released changelog in %s: %v", WorkingChangelog, err)
		return err
	}

	if err := copyChangelogTemplate(); err != nil {
		return err
	}

	// Merge all changelogs into one
	buf, err := MergeChangelogs()
	if err != nil {
		return err
	}

	archivePath := fmt.Sprintf("release/archives/%s", ctx.Version)
	if err := os.Mkdir(archivePath, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Unable to create archive directory %s: %v", archivePath, err)
		return err
	}

	for _, target := range ctx.Config.Targets {
		if err := buildRelease(ctx, target, archivePath, buf); err != nil {
			return err
		}
	}

	if err := generateChecksumFile(ctx, buf); err != nil {
		return err
	}

	if commitNewRelease {
		commitMessage := fmt.Sprintf(
			`Release %s

			Released version %s (%s)
			`, ctx.Version, ctx.Version, ctx.Name)
		hash, err := gitutil.CreateCommit(
			ctx.Config.SourceRoot,
			ctx.Config.CommitterName,
			ctx.Config.CommitterEmail,
			commitMessage,
			WorkingChangelog,
			fmt.Sprintf("release/%s/changelog.md", ctx.Version))
		if err != nil {
			toolbox.PrintError("Could not commit the new release files: %v", err)
			return err
		}
		fmt.Printf("New change log is committed as %s\n", hash[:6])
	}
	return nil
}

func writeZipped(z *zip.Writer, header *zip.FileHeader, buf []byte) error {
	header.Method = zip.Deflate
	zf, err := z.CreateHeader(header)
	if err != nil {
		toolbox.PrintError("Could not create zip entry %s: %v", header.Name, err)
		return err
	}
	if _, err := zf.Write(buf); err != nil {
		return err
	}
	return nil
}

// buildRelease builds a release archive for a particular target
func buildRelease(ctx *Context, target, archivePath string, changeLog []byte) error {
	archive := fmt.Sprintf("%s/%s-%s_%s.zip", archivePath, ctx.Version, ctx.Config.Name, target)

	fmt.Printf("Building release archive %s \n", archive)

	f, err := os.Create(archive)
	if err != nil {
		toolbox.PrintError("Unable to create archive file %s: %v", archive, err)
		return err
	}
	zipWriter := zip.NewWriter(f)
	comment := fmt.Sprintf("This archive contains an %s build for %s v%s (%s). Please see changelog for details.", ctx.Config.Name, target, ctx.Version, ctx.Name)
	zipWriter.SetComment(comment)
	defer zipWriter.Close()

	fmt.Printf(" - [changelog] changelog.md\n")
	logHeader := &zip.FileHeader{
		Name:               "changelog.md",
		Method:             zip.Deflate,
		UncompressedSize64: uint64(len(changeLog)),
		Modified:           time.Now(),
	}
	if err := writeZipped(zipWriter, logHeader, changeLog); err != nil {
		return err
	}
	for _, v := range ctx.Config.Files {
		if v.Target == anyTarget || v.Target == target {
			fmt.Printf(" - [%s] %s\n", v.ID, v.Name)
			buf, err := ioutil.ReadFile(v.Name)
			if err != nil {
				toolbox.PrintError("Unable to read file %s: %v", v.Name, err)
				return err
			}
			fi, err := os.Stat(v.Name)
			if err != nil {
				toolbox.PrintError("Could not stat %s: %v", v.Name, err)
				return err
			}
			header, err := zip.FileInfoHeader(fi)
			if err != nil {
				toolbox.PrintError("Could not create file info header for %s: %v", v.Name, err)
				return err
			}
			writeZipped(zipWriter, header, buf)
		}
	}
	return nil
}
