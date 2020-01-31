package commands

import (
	"fmt"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/changelog"
	"github.com/ExploratoryEngineering/releasetool/pkg/release"
	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

type initCommand struct {
}

func (c *initCommand) Run(rc RunContext) error {
	// Make sure the release directory doesn't exist
	err := os.Mkdir("release", toolbox.DefaultFilePerm)
	if os.IsExist(err) {
		toolbox.PrintError("The 'release' directory already exists.")
		return err
	}
	if err != nil {
		toolbox.PrintError("Error creating the release directory: %v", err)
		return err
	}

	f, err := os.Create(release.VersionFile)
	if os.IsExist(err) {
		toolbox.PrintError("The VERSION file already exists in the release directory")
		return err
	}
	if err != nil {
		toolbox.PrintError("Error creating the %s file: %v", release.VersionFile, err)
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(initialVersion))
	if os.IsPermission(err) {
		toolbox.PrintError("Permission denied on the %s file. Can't write initial version", release.VersionFile)
		return err
	}
	if err != nil {
		toolbox.PrintError("Error writing initial version to the %s file: %v", release.VersionFile, err)
		return err
	}

	templateDir := "release/templates"

	if err := os.Mkdir(templateDir, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create the template directory: %v", err)
		return err
	}

	// Make template files
	if err := changelog.MakeTemplate(); err != nil {
		return err
	}
	if err := changelog.CopyTemplate(); err != nil {
		return err
	}

	if err := release.WriteSampleConfig(); err != nil {
		return err
	}

	fmt.Printf("Initialized version to %s. Working change log is in release/changelog.md\n", initialVersion)
	return nil
}
