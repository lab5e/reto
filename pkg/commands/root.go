package commands

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
}

func (r *Root) ReleaseCommands() Root {
	return *r
}

type RunContext interface {
	ReleaseCommands() Root
}
