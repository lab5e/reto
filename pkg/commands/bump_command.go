package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ExploratoryEngineering/releasetool/pkg/release"
	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

type bumpCommand struct {
	Major bool `kong:"short='M',help='Bump major version'"`
	Minor bool `kong:"short='m',help='Bump minor version'"`
	Patch bool `kong:"short='p',help='Bump patch version'"`
}

func (c *bumpCommand) Run(rc RunContext) error {
	config, err := release.Verify()
	if err != nil {
		return err
	}

	tuples := strings.Split(config.Version, ".")
	if len(tuples) != 3 {
		toolbox.PrintError("Invalid version string in version file: %s", config.Version)
		return errors.New("invalid version")
	}

	bumps := 0
	if rc.ReleaseCommands().Bump.Major {
		config.Major++
		config.Minor = 0
		config.Patch = 0
		bumps++
	}

	if rc.ReleaseCommands().Bump.Minor {
		config.Minor++
		config.Patch = 0
		bumps++
	}

	if rc.ReleaseCommands().Bump.Patch {
		config.Patch++
		bumps++
	}

	if bumps == 0 {
		toolbox.PrintError("Must specify which version to bump")
		return errors.New("no bump")
	}

	if bumps != 1 {
		toolbox.PrintError("Only onf of bump major, minor or patch can be bumped")
		return errors.New("arg error")
	}

	newVersion := fmt.Sprintf("%d.%d.%d", config.Major, config.Minor, config.Patch)

	if err := ioutil.WriteFile(release.VersionFile, []byte(newVersion), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Error writing version file: %v", err)
		return err
	}
	fmt.Printf("New version is now %s\n", newVersion)
	return nil
}
