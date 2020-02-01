package commands

// Root is the root command
type Root struct {
	Ver      versionFlag     `kong:"name='ver',short='v',help='Show release tool version'"`
	Root     string          `kong:"short='r',help='Root directory for the tool'"`
	Init     initCommand     `kong:"cmd,help='Initialise release tool'"`
	Version  versionCommand  `kong:"cmd,help='Show current version'"`
	Bump     bumpCommand     `kong:"cmd,help='Version bumping'"`
	Hash     hashCommand     `kong:"cmd,help='Display current git hash'"`
	Hashname hashNameCommand `kong:"cmd,help='Display current git hash as human readable name'"`
	Namehash nameHashCommand `kong:"cmd,help='Display human readable name as git hash'"`
	Status   statusCommand   `kong:"cmd,help='Display current status'"`
	Release  releaseCommand  `kong:"cmd,help='Generate a release from existing binaries'"`
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
