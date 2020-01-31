package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/release"
)

// versionCommand displays the current version
type versionCommand struct {
}

func (c *versionCommand) Run(rc RunContext) error {
	config, err := release.Verify()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", config.Version)
	return nil
}
