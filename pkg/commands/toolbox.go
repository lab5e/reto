package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func printError(fmtString string, opts ...interface{}) {
	if opts == nil {
		fmt.Fprintf(os.Stderr, "releasetool: %s\n", fmt.Sprintf(fmtString, opts...))
		return
	}
	fmt.Fprintf(os.Stderr, "releasetool: %s\n", fmt.Sprintf(fmtString, opts...))
}

// verifySetup verifies that the release tool is initialized
func verifySetup() (*releaseConfig, error) {
	if _, err := os.Stat(versionFile); err != nil {
		printError("Can't read the version file: %v", err)
		return nil, errors.New("no version file")
	}
	buf, err := ioutil.ReadFile(versionFile)
	if err != nil {
		printError("Unable to read version file: %v", err)
		return nil, err
	}
	lines := strings.Split(string(buf), "\n")
	if len(lines) == 0 {
		printError("Version file does not contain a version")
		return nil, errors.New("no version found")
	}
	return &releaseConfig{
		Version: lines[0],
	}, nil
}
