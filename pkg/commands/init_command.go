package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/release"
)

type initCommand struct {
}

func (c *initCommand) Run(rc RunContext) error {
	if err := release.InitTool(); err != nil {
		return err
	}
	fmt.Printf("Initialized releasetool.\n")
	return nil
}
