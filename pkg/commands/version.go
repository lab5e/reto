package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/reto/pkg/release"
)

// versionCommand displays the current version
type versionCommand struct {
}

func (c *versionCommand) Run(rc RunContext) error {
	config, err := release.GetContext()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", config.Version)
	return nil
}
