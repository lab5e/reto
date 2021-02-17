package commands

import (
	"github.com/lab5e/reto/pkg/release"
)

type verifyCommand struct {
	Archive    string `kong:"required,help='Zip archive to verify',type='existingfile'"`
	SHA256File string `kong:"required,help='SHA256 Checksum file',type='existingfile'"`
}

func (c *verifyCommand) Run(rc RunContext) error {
	return release.VerifyChecksums(rc.ReleaseCommands().Verify.SHA256File, rc.ReleaseCommands().Verify.Archive)
}
