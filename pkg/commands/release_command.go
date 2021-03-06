package commands

import (
	"fmt"

	"github.com/lab5e/reto/pkg/release"
	"github.com/lab5e/reto/pkg/toolbox"
)

type releaseCommand struct {
	CreateTag bool `kong:"short='t',help='Tag version in Git',default='true'"`
	Commit    bool `kong:"short='c',help='Commit new change log after release',default='true'"`
}

func (c *releaseCommand) Run(rc RunContext) error {
	if err := release.Build(
		rc.ReleaseCommands().Release.CreateTag,
		rc.ReleaseCommands().Release.Commit); err != nil {
		return err
	}
	fmt.Print("Release built. ")
	if rc.ReleaseCommands().Release.CreateTag {
		fmt.Printf("%sPush tags with 'git push --tags'%s.", toolbox.Cyan, toolbox.Reset)
	}
	fmt.Println()
	return nil
}
