package commands

import (
	"errors"

	"github.com/ExploratoryEngineering/releasetool/pkg/release"
	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

type releaseCommand struct {
}

func (c *releaseCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}
	if ctx.Released {
		toolbox.PrintError("This version is already released. Bump the version and try again.")
		return errors.New("already released")
	}

	// Do a quick sanity check on the change log
	if err := release.ChangelogComplete(); err != nil {
		return err
	}
	if err := release.VerifyConfig(ctx.Config); err != nil {
		return err
	}
	for _, target := range ctx.Config.Targets {
		release.BuildRelease(ctx, target)
	}
	return nil
}
