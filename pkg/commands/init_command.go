package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/templates"
)

type initCommand struct {
}

func (c *initCommand) Run(rc RunContext) error {
	// Make sure the release directory doesn't exist
	err := os.Mkdir(releaseDir, defaultFilePerm)
	if os.IsExist(err) {
		printError("The 'release' directory already exists.")
		return err
	}
	if err != nil {
		printError("Error creating the release directory: %v", err)
		return err
	}

	f, err := os.Create(versionFile)
	if os.IsExist(err) {
		printError("The VERSION file already exists in the release directory")
		return err
	}
	if err != nil {
		printError("Error creating the %s file: %v", versionFile, err)
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(initialVersion))
	if os.IsPermission(err) {
		printError("Permission denied on the %s file. Can't write initial version", versionFile)
		return err
	}
	if err != nil {
		printError("Error writing initial version to the %s file: %v", versionFile, err)
		return err
	}

	templateDir := fmt.Sprintf("%s%ctemplates", releaseDir, os.PathSeparator)

	if err := os.Mkdir(templateDir, defaultFilePerm); err != nil {
		printError("Could not create the template directory: %v", err)
		return err
	}

	// Make template files
	templateFile := fmt.Sprintf("%s%cchangelog-template.md", templateDir, os.PathSeparator)
	if err := ioutil.WriteFile(
		templateFile,
		[]byte(templates.DefaultChangeLogTemplate), defaultFilePerm); err != nil {
		printError("Unable to create the release log template: %v", err)
		return err
	}

	// Copy the template into the release folder
	workingChangelog := fmt.Sprintf("%s%cchangelog.md", releaseDir, os.PathSeparator)
	if err := copyFile(templateFile, workingChangelog); err != nil {
		printError("Could not copy changelog template to release directory: %v", err)
		return err
	}

	fmt.Printf("Initialized version to %s. Working change log is in %s%cchangelog.md\n", initialVersion, releaseDir, os.PathSeparator)
	return nil
}
