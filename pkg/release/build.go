package release

import (
	"errors"
	"fmt"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

func archiveName(ctx *Context, target string) string {
	return fmt.Sprintf("release/archives/%s/%s-%s_%s.zip", ctx.Version, ctx.Version, ctx.Config.Name, target)
}

// Build builds a new release from the current setup
func Build() error {
	ctx, err := GetContext()
	if err != nil {
		return err
	}
	if ctx.Released {
		toolbox.PrintError("This version is already released. Bump the version and try again.")
		return errors.New("already released")
	}

	// Do a quick sanity check on the change log
	if err := ChangelogComplete(); err != nil {
		return err
	}
	if err := VerifyConfig(ctx.Config); err != nil {
		return err
	}
	if err := os.Mkdir(fmt.Sprintf("release/%s", ctx.Version), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create release directory: %v", err)
		return err
	}

	// build the changelog for the release
	if err := releaseChangelog(ctx); err != nil {
		return err
	}

	if err := os.Remove("release/changelog.md"); err != nil {
		toolbox.PrintError("Could not remove released changelog in release/changelog.md: %v", err)
		return err
	}

	if err := copyChangelogTemplate(); err != nil {
		return err
	}

	// Merge all changelogs into one
	for _, target := range ctx.Config.Targets {
		if err := buildRelease(ctx, target); err != nil {
			return err
		}
	}
	return nil
}

// buildRelease builds a release archive for a particular target
func buildRelease(ctx *Context, target string) error {

	fmt.Printf("Building release archive %s \n", archiveName(ctx, target))
	fmt.Printf(" - [changelog] changelog.md\n")
	for _, v := range ctx.Config.Files {
		if v.Target == anyTarget || v.Target == target {
			fmt.Printf(" - [%s] %s\n", v.ID, v.Name)
		}
	}
	return nil
}
