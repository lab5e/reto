package commands

import "github.com/ExploratoryEngineering/releasetool/pkg/release"

// Root is the root command
type Root struct {
	Ver      versionFlag     `kong:"name='ver',short='v',help='Show release tool version'"`
	Init     initCommand     `kong:"cmd,help='Initialise release tool'"`
	Version  versionCommand  `kong:"cmd,help='Show current version'"`
	Bump     bumpCommand     `kong:"cmd,help='Version bumping'"`
	Hash     hashCommand     `kong:"cmd,help='Display current git hash'"`
	Hashname hashNameCommand `kong:"cmd,help='Display current git hash as human readable name'"`
	Namehash nameHashCommand `kong:"cmd,help='Display human readable name as git hash'"`
	Status   statusCommand   `kong:"cmd,help='Display current status'"`
	Release  releaseCommand  `kong:"cmd,help='Generate a release from existing binaries'"`
	Test     testCommand     `kong:"cmd,help='Command for testing.'"`
	Checksum checksumCommand `kong:"cmd,help='Show checksum for built files'"`
	Verify   verifyCommand   `kong:"cmd,help='Verify checksum on archive and files inside archive'"`
}

// ReleaseCommands is the command parameters
func (r *Root) ReleaseCommands() Root {
	return *r
}

// RunContext is a context for the commands. This decouples the parameter struct
// from the commands... just slightly
type RunContext interface {
	ReleaseCommands() Root
}

type testCommand struct {
}

func (c *testCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}
	if err := release.GenerateSHA256File(ctx); err != nil {
		return err
	}
	return nil
}
