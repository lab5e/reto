package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/reto/pkg/release"
)

type releaseCommand struct {
	CreateTag      bool   `kong:"short='t',help='Tag version in Git',default='true'"`
	Commit         bool   `kong:"short='c',help='Commit new change log after release',default='true'"`
	CommitterName  string `kong:"short='e',help='Override committer name from config',default=''"`
	CommitterEmail string `kong:"short='e',help='Override committer email from config',default=''"`
}

func (c *releaseCommand) Run(rc RunContext) error {
	if err := release.Build(
		rc.ReleaseCommands().Release.CreateTag,
		rc.ReleaseCommands().Release.Commit,
		rc.ReleaseCommands().Release.CommitterName,
		rc.ReleaseCommands().Release.CommitterEmail); err != nil {
		return err
	}
	fmt.Print("Release built. ")
	if rc.ReleaseCommands().Release.CreateTag {
		fmt.Print("Push tags with 'git push --tags'.")
	}
	fmt.Println()
	return nil
}
