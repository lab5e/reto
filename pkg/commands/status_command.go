package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/release"
)

type statusCommand struct {
}

func (c *statusCommand) Run(rc RunContext) error {
	config, err := release.Verify()
	if err != nil {
		return err
	}

	released := "NO"
	if config.Released {
		released = "YES"
	}

	fmt.Printf("Active version: %s\n", config.Version)
	fmt.Printf("Released:       %s\n", released)
	return nil
}
