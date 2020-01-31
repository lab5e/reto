package commands

import (
	"fmt"
	"os"
)

const (
	releaseDir      = "release"
	initialVersion  = "0.0.0"
	versionFilePerm = 0600
)

var versionFile = fmt.Sprintf("%s%cVERSION", releaseDir, os.PathSeparator)

type Root struct {
	Init    initCommand    `kong:"cmd,help='Initialise release tool'"`
	Version versionCommand `kong:"cmd,help='Show current version'"`
	Bump    bumpCommand    `kong:"cmd,help='Version bumping'"`
}

func (r *Root) ReleaseCommands() Root {
	return *r
}

type RunContext interface {
	ReleaseCommands() Root
}

type releaseConfig struct {
	Version string
	Major   int
	Minor   int
	Patch   int
}
