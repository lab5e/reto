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
	Init initCommand `kong:"cmd,help='Initialise release tool'"`
}

type RunContext interface {
}
