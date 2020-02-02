package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/reto/pkg/release"
)

type releaseCommand struct {
	CreateTag bool `kong:"short='t',help='Tag version in Git',default='true'"`
}

func (c *releaseCommand) Run(rc RunContext) error {
	if err := release.Build(rc.ReleaseCommands().Release.CreateTag); err != nil {
		return err
	}
	fmt.Print("Release built. ")
	if rc.ReleaseCommands().Release.CreateTag {
		fmt.Print("Push tags with 'git push --tags'.")
	}
	fmt.Println()
	return nil
}
