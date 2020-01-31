package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/release"
)

type releaseCommand struct {
}

func (c *releaseCommand) Run(rc RunContext) error {
	if err := release.Build(); err != nil {
		return err
	}
	fmt.Println("Release built.")
	return nil
}
