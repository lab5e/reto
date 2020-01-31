package commands

import (
	"fmt"
)

type statusCommand struct {
}

func (c *statusCommand) Run(rc RunContext) error {
	config, err := verifySetup()
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
