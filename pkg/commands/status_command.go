package commands

import (
	"fmt"
	"os"
)

type statusCommand struct {
}

func (c *statusCommand) Run(rc RunContext) error {
	config, err := verifySetup()
	if err != nil {
		return err
	}

	_, err = os.Stat(fmt.Sprintf("%s%c%s", releaseDir, os.PathSeparator, config.Version))
	released := "YES"
	if os.IsNotExist(err) {
		released = "NO"
	}
	if err != nil && !os.IsNotExist(err) {
		printError("Could not read release directory: %v", err)
		return err
	}
	fmt.Printf("Active version: %s\n", config.Version)
	fmt.Printf("Released:       %s\n", released)
	return nil
}
