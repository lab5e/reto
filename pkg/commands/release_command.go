package commands

import (
	"errors"
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/changelog"
	"github.com/ExploratoryEngineering/releasetool/pkg/release"
	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

type releaseCommand struct {
}

func (c *releaseCommand) Run(rc RunContext) error {
	ctx, err := release.Verify()
	if err != nil {
		return err
	}
	if ctx.Released {
		toolbox.PrintError("This version is already released. Bump the version and try again.")
		return errors.New("already released")
	}

	// Do a quick sanity check on the change log
	if err := changelog.WorkingCopyIsComplete(); err != nil {
		return err
	}
	if err := release.VerifyConfig(ctx.Config); err != nil {
		return err
	}
	fmt.Println("Doing release")

	return nil
}
