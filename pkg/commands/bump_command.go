package commands

import (
	"fmt"

	"github.com/lab5e/reto/pkg/release"
)

type bumpCommand struct {
	Major bool `kong:"short='M',help='Bump major version'"`
	Minor bool `kong:"short='m',help='Bump minor version'"`
	Patch bool `kong:"short='p',help='Bump patch version'"`
}

func (c *bumpCommand) Run(rc RunContext) error {
	ctx, err := release.BumpVersion(rc.ReleaseCommands().Bump.Major, rc.ReleaseCommands().Bump.Minor, rc.ReleaseCommands().Bump.Patch)
	if err != nil {
		return err
	}
	fmt.Printf("New version is now %s\n", ctx.Version)
	return nil
}
