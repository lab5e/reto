package release

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

var VersionFile = "release/VERSION"

// ReleaseContext is a general release configuration type
type ReleaseContext struct {
	Version  string
	Major    int
	Minor    int
	Patch    int
	Released bool
}

// Verify verifies that the release tool is initialized
func Verify() (*ReleaseContext, error) {
	if _, err := os.Stat(VersionFile); err != nil {
		toolbox.PrintError("Can't read the version file: %v", err)
		return nil, errors.New("no version file")
	}
	buf, err := ioutil.ReadFile(VersionFile)
	if err != nil {
		toolbox.PrintError("Unable to read version file: %v", err)
		return nil, err
	}
	lines := strings.Split(string(buf), "\n")
	if len(lines) == 0 {
		toolbox.PrintError("Version file does not contain a version")
		return nil, errors.New("no version found")
	}
	ret := ReleaseContext{
		Version: lines[0],
	}

	var versionErr = errors.New("invalid version content")
	tuples := strings.Split(ret.Version, ".")
	if len(tuples) != 3 {
		toolbox.PrintError("Version string is malformed: %s", ret.Version)
		return nil, versionErr
	}
	v, err := strconv.ParseInt(tuples[0], 10, 63)
	if err != nil {
		toolbox.PrintError("Major version is not an integer: %s", ret.Version)
		return nil, versionErr
	}
	ret.Major = int(v)

	v, err = strconv.ParseInt(tuples[1], 10, 63)
	if err != nil {
		toolbox.PrintError("Minor version is not an integer: %s", ret.Version)
		return nil, versionErr
	}
	ret.Minor = int(v)

	v, err = strconv.ParseInt(tuples[2], 10, 63)
	if err != nil {
		toolbox.PrintError("Patch version is not an integer: %s", ret.Version)
		return nil, versionErr
	}
	ret.Patch = int(v)

	_, err = os.Stat(fmt.Sprintf("release/%s", ret.Version))
	ret.Released = true
	if os.IsNotExist(err) {
		ret.Released = false
	}
	if err != nil && !os.IsNotExist(err) {
		toolbox.PrintError("Could not read release directory: %v", err)
		return nil, err
	}

	return &ret, nil
}
