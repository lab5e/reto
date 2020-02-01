package commands

import "github.com/ExploratoryEngineering/releasetool/pkg/release"

type checksumCommand struct {
}

func (c *checksumCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}
	return release.GenerateSHA256File(ctx)
}

type verifyCommand struct {
	Archive    string `kong:"required,help='Zip archive to verify',type='existingfile'"`
	SHA256File string `kong:"required,help='SHA256 Checksum file',type='existingfile'"`
}

func (c *verifyCommand) Run(rc RunContext) error {
	return release.VerifyChecksums(rc.ReleaseCommands().Verify.SHA256File, rc.ReleaseCommands().Verify.Archive)
}
