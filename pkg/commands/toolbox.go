package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
	ret := releaseConfig{
		Version: lines[0],
	}

	var versionErr = errors.New("invalid version content")
	tuples := strings.Split(ret.Version, ".")
	if len(tuples) != 3 {
		printError("Version string is malformed: %s", ret.Version)
		return nil, versionErr
	}
	v, err := strconv.ParseInt(tuples[0], 10, 63)
	if err != nil {
		printError("Major version is not an integer: %s", ret.Version)
		return nil, versionErr
	}
	ret.Major = int(v)

	v, err = strconv.ParseInt(tuples[1], 10, 63)
	if err != nil {
		printError("Minor version is not an integer: %s", ret.Version)
		return nil, versionErr
	}
	ret.Minor = int(v)

	v, err = strconv.ParseInt(tuples[2], 10, 63)
	if err != nil {
		printError("Patch version is not an integer: %s", ret.Version)
		return nil, versionErr
	}
	ret.Patch = int(v)

	return &ret, nil
}
