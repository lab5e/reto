package commands

import (
	"errors"

	"github.com/ExploratoryEngineering/reto/pkg/release"
	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
)

type checksumCommand struct {
}

func (c *checksumCommand) Run(rc RunContext) error {
	toolbox.PrintError("not implemented")
	return errors.New("not implemented")
}

type verifyCommand struct {
	Archive    string `kong:"required,help='Zip archive to verify',type='existingfile'"`
	SHA256File string `kong:"required,help='SHA256 Checksum file',type='existingfile'"`
}

func (c *verifyCommand) Run(rc RunContext) error {
	return release.VerifyChecksums(rc.ReleaseCommands().Verify.SHA256File, rc.ReleaseCommands().Verify.Archive)
}
