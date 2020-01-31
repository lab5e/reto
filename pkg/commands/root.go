package commands

import (
	"fmt"
	"os"
)

const (
	releaseDir     = "release"
	initialVersion = "0.0.0"
)

var versionFile = fmt.Sprintf("%s%cVERSION", releaseDir, os.PathSeparator)

type Root struct {
	Init    initCommand    `kong:"cmd,help='Initialise release tool'"`
	Version versionCommand `kong:"cmd,help='Show current version'"`
}

type RunContext interface {
}

type releaseConfig struct {
	Version string
}
