package release

import (
	"io/ioutil"
	"os"

	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
)

const (
	initialVersion = "0.0.0"
)

// InitTool initializes the directory structure for the tool. Errors are printed
// to stderr.
func InitTool() error {
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

	templateDir := "release/templates"

	if err := os.Mkdir(templateDir, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create the template directory: %v", err)
		return err
	}

	if err := initChangelog(); err != nil {
		return err
	}

	if err := WriteSampleConfig(); err != nil {
		return err
	}

	archiveDir := "release/archives"
	if err := os.Mkdir(archiveDir, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create the archives directory: %v", err)
		return err
	}

	if err := ioutil.WriteFile("release/.gitignore", []byte("archivea\n"), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create .gitignore file in release directory: %v", err)
		return err
	}
	return nil
}
