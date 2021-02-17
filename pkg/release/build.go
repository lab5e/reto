package release

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/lab5e/reto/pkg/gitutil"
	"github.com/lab5e/reto/pkg/toolbox"
)

// Build builds a new release from the current setup
func Build(tagVersion, commitNewRelease bool) error {
	fi, err := os.Stat(archiveDir)
	if err != nil {
		if err := os.Mkdir(archiveDir, toolbox.DefaultDirPerm); err != nil {
			toolbox.PrintError("Could not create archive directory: %v", err)
			return err
		}
	}
	if err != nil {
		toolbox.PrintError("Can't check status of archive dir: %v", err)
		return err
	}
	if !fi.IsDir() {
		toolbox.PrintError("%s is not a directory", archiveDir)
		return errors.New("no archive")
	}
	ctx, err := GetContext()
	if err != nil {
		return err
	}
	if ctx.Released {
		toolbox.PrintError("This version is already released. Bump the version and try again.")
		return errors.New("already released")
	}

	if err := TemplatesComplete(ctx, true); err != nil {
		return err
	}

	if err := VerifyConfig(ctx.Config, true); err != nil {
		return err
	}
	if gitutil.HasChanges(ctx.Config.SourceRoot, true) {
		toolbox.PrintError("There are uncommitted or unstaged changes in the current Git branch")
		return errors.New("uncommitted changes")
	}

	if !NewFileVersions(ctx.Config, true) {
		return errors.New("no file changes")
	}

	if err := os.Mkdir(fmt.Sprintf("%s/%s", releaseDir, ctx.Version), toolbox.DefaultDirPerm); err != nil {
		toolbox.PrintError("Could not create release directory: %v", err)
		return err
	}

	archivePath := fmt.Sprintf("%s/%s", archiveDir, ctx.Version)
	if err := os.Mkdir(archivePath, toolbox.DefaultDirPerm); err != nil {
		toolbox.PrintError("Unable to create archive directory %s: %v", archivePath, err)
		return err
	}

	var checksumFiles []string // Files that should be checksummed
	var tempFiles []string     // temp files that should be deleted when done
	for _, template := range ctx.Config.Templates {
		workingCopy := fmt.Sprintf("release/%s", template.Name)
		releasedCopy := fmt.Sprintf("%s/%s/%s", releaseDir, ctx.Version, template.Name)
		archiveCopy := fmt.Sprintf("%s/%s/%s", archiveDir, ctx.Version, template.Name)
		if err := mergeTemplate(workingCopy, releasedCopy, ctx); err != nil {
			return err
		}
		if err := os.Remove(workingCopy); err != nil {
			toolbox.PrintError("Could not remove template %s: %v", template, err)
			return err
		}
		if err := toolbox.CopyFile(fmt.Sprintf("%s/%s", templateDir, template.Name), workingCopy); err != nil {
			toolbox.PrintError("Could not copy %s to release directory: %v", workingCopy, err)
			return err
		}
		if template.TemplateAction == ConcatenateAction {
			if err := concatenateTemplate(template.Name, archiveCopy); err != nil {
				return err
			}
			tempFiles = append(tempFiles, archiveCopy)
			checksumFiles = append(checksumFiles, archiveCopy)
		} else {
			checksumFiles = append(checksumFiles, releasedCopy)
		}
	}
	for _, target := range ctx.Config.Targets {
		if err := buildRelease(ctx, target, archivePath, checksumFiles); err != nil {
			return err
		}
	}

	for _, file := range ctx.Config.Files {
		checksumFiles = append(checksumFiles, file.Name)
	}
	if err := generateChecksumFile(ctx, checksumFiles); err != nil {
		return err
	}

	// Remove the generate files in the archive folder
	for _, v := range tempFiles {
		if err := os.Remove(v); err != nil {
			toolbox.PrintError("Could not remove temporary file at %s: %v", v, err)
			return err
		}
	}

	// Tag the commit with the new version
	if tagVersion {
		if err := gitutil.TagVersion(
			ctx.Config.SourceRoot,
			fmt.Sprintf("v%s", ctx.Version),
			fmt.Sprintf("Release v%s (%s)", ctx.Version, ctx.Name)); err != nil {
			return err
		}
	}

	if commitNewRelease {
		commitMessage := fmt.Sprintf(
			`Release %s

			Released version %s (%s)
			`, ctx.Version, ctx.Version, ctx.Name)
		var filesToCommit []string
		for _, v := range ctx.Config.Templates {
			filesToCommit = append(filesToCommit, fmt.Sprintf("%s/%s/%s", releaseDir, ctx.Version, filepath.Base(v.Name)))
			filesToCommit = append(filesToCommit, fmt.Sprintf("release/%s", filepath.Base(v.Name)))
		}
		hash, err := gitutil.CreateCommit(
			ctx.Config.SourceRoot,
			commitMessage,
			filesToCommit...)
		if err != nil {
			toolbox.PrintError("Could not commit the new release files: %v", err)
			return err
		}
		fmt.Printf("New change log is committed as %s\n", hash[:6])
		newCtx, err := BumpVersion(false, false, true)
		if err != nil {
			toolbox.PrintError("Could not autobump version: %v", err)
			return nil
		}
		fmt.Printf("auto-bumped new version to %s\n", newCtx.Version)
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
func buildRelease(ctx *Context, target, archivePath string, tempFiles []string) error {
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

	for _, tempFile := range tempFiles {
		buf, err := ioutil.ReadFile(tempFile)
		if err != nil {
			toolbox.PrintError("Could not read temp file %s: %v", tempFile, err)
			return err
		}
		fmt.Printf(" - [template] %s\n", filepath.Base(tempFile))
		tempHeader := &zip.FileHeader{
			Name:               filepath.Base(tempFile),
			Method:             zip.Deflate,
			UncompressedSize64: uint64(len(buf)),
			Modified:           time.Now(),
		}
		if err := writeZipped(zipWriter, tempHeader, buf); err != nil {
			return err
		}

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
