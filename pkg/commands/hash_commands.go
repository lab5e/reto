package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/gitutil"
	"github.com/ExploratoryEngineering/releasetool/pkg/hashname"
)

type hashNameCommand struct {
}

func (c *hashNameCommand) Run(rc RunContext) error {
	hash, err := gitutil.GetCurrentHash()
	if err != nil {
		printError("Unable to read git hash: %v", err)
		return err
	}
	fmt.Printf("%s\n", hashname.HashToName(hash))
	return nil
}

type nameHashCommand struct {
	Name string `kong:"arg,help='Name to convert into hash'"`
}

func (c *nameHashCommand) Run(rc RunContext) error {
	hash, err := hashname.NameToHash(rc.ReleaseCommands().Namehash.Name)
	if err != nil {
		printError("Unable to convert name into hash: %v", err)
		return err
	}
	fmt.Printf("%s\n", hash)
	return nil
}
