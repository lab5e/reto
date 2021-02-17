package release

import (
	"io/ioutil"
	"os"

	"github.com/lab5e/reto/pkg/toolbox"
)

const (
	initialVersion = "0.0.0"
	archiveDir     = "release/archives"
	releaseDir     = "release/releases"
	templateDir    = "release/templates"
)

// InitTool initializes the directory structure for the tool. Errors are printed
// to stderr.
func InitTool() error {
	// Make sure the release directory doesn't exist
	err := os.MkdirAll(releaseDir, toolbox.DefaultDirPerm)
	if err != nil {
		toolbox.PrintError("Error creating the release directory: %v", err)
		return err
	}
	if err := os.MkdirAll(templateDir, toolbox.DefaultDirPerm); err != nil {
		toolbox.PrintError("Could not create the template directory: %v", err)
		return err
	}
	if err := os.MkdirAll(archiveDir, toolbox.DefaultDirPerm); err != nil {
		toolbox.PrintError("Could not create the archive directory: %v", err)
		return err
	}

	f, err := os.Create(VersionFile)
	if os.IsExist(err) {
		toolbox.PrintError("The VERSION file already exists in the release directory")
		return err
	}
	if err != nil {
		toolbox.PrintError("Error creating the %s file: %v", VersionFile, err)
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(initialVersion))
	if os.IsPermission(err) {
		toolbox.PrintError("Permission denied on the %s file. Can't write initial version", VersionFile)
		return err
	}
	if err != nil {
		toolbox.PrintError("Error writing initial version to the %s file: %v", VersionFile, err)
		return err
	}

	if err := initTemplates(); err != nil {
		return err
	}

	if err := writeSampleConfig(); err != nil {
		return err
	}

	if err := ioutil.WriteFile("release/.gitignore", []byte("archives\n"), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create .gitignore file in release directory: %v", err)
		return err
	}
	return nil
}
