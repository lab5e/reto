package commands

import (
	"fmt"
	"os"
)

type initCommand struct {
}

func (c *initCommand) Run(rc RunContext) error {
	fmt.Println("This is the init command")
	// Make sure the release directory exists
	err := os.Mkdir(releaseDir, versionFilePerm)
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

	fmt.Printf("Initialized version to %s\n", initialVersion)
	return nil
}
