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
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Architectures: ")
	for _, v := range config.Config.Architectures {
		fmt.Printf("  - %s\n", v)
	}
	fmt.Println("  OSes:")
	for _, v := range config.Config.OSes {
		fmt.Printf("  - %s\n", v)
	}
	fmt.Println("  Files:")
	for _, v := range config.Config.Files {
		fmt.Printf("  - %s (%s, %s)\n", v.Name, v.Arch, v.OS)
	}
	fmt.Println()

	return nil
}
