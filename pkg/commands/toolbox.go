package commands

import (
	"errors"
	"fmt"
	"os"
)

// verifySetup verifies that the release tool is initialized
func verifySetup() error {
	if _, err := os.Stat(versionFile); err != nil {
		fmt.Printf("Can't read the version file: %v\n", err)
		return errors.New("no version file")
	}
	return nil
}
