package commands

import (
	"fmt"
)

// versionCommand displays the current version
type versionCommand struct {
}

func (c *versionCommand) Run(rc RunContext) error {
	config, err := verifySetup()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", config.Version)
	return nil
}
