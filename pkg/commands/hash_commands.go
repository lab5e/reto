package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/hashname"
	"github.com/ExploratoryEngineering/releasetool/pkg/release"
	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

type hashCommand struct {
}

type hashNameCommand struct {
}

func (c *hashCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", ctx.CommitHash)
	return nil
}

func (c *hashNameCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", hashname.HashToName(ctx.CommitHash))
	return nil
}

type nameHashCommand struct {
	Name string `kong:"arg,help='Name to convert into hash'"`
}

func (c *nameHashCommand) Run(rc RunContext) error {
	hash, err := hashname.NameToHash(rc.ReleaseCommands().Namehash.Name)
	if err != nil {
		toolbox.PrintError("Unable to convert name into hash: %v", err)
		return err
	}
	fmt.Printf("%s\n", hash)
	return nil
}
