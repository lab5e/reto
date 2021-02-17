package release

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/lab5e/reto/pkg/gitutil"
	"github.com/lab5e/reto/pkg/hashname"
)

// VersionFile is the path to the fil containing the version
var VersionFile = "release/VERSION"

// Context is a general release configuration type
type Context struct {
	Config     Config
	Version    string
	Major      int
	Minor      int
	Patch      int
	Released   bool
	CommitHash string
	Name       string
}

// GetContext verifies that the release tool is initialized
func GetContext() (*Context, error) {
	if _, err := os.Stat(VersionFile); err != nil {
		fmt.Printf("Can't read the version file: %v\n", err)
		return nil, errors.New("no version file")
	}
	buf, err := ioutil.ReadFile(VersionFile)
	if err != nil {
		fmt.Printf("Unable to read version file: %v\n", err)
		return nil, err
	}
	lines := strings.Split(string(buf), "\n")
	if len(lines) == 0 {
		fmt.Printf("Version file does not contain a version\n")
		return nil, errors.New("no version found")
	}
	ret := Context{
		Version: lines[0],
	}

	var versionErr = errors.New("invalid version content")
	tuples := strings.Split(ret.Version, ".")
	if len(tuples) != 3 {
		fmt.Printf("Version string is malformed: %s\n", ret.Version)
		return nil, versionErr
	}
	v, err := strconv.ParseInt(tuples[0], 10, 63)
	if err != nil {
		fmt.Printf("Major version is not an integer: %s\n", ret.Version)
		return nil, versionErr
	}
	ret.Major = int(v)

	v, err = strconv.ParseInt(tuples[1], 10, 63)
	if err != nil {
		fmt.Printf("Minor version is not an integer: %s\n", ret.Version)
		return nil, versionErr
	}
	ret.Minor = int(v)

	v, err = strconv.ParseInt(tuples[2], 10, 63)
	if err != nil {
		fmt.Printf("Patch version is not an integer: %s\n", ret.Version)
		return nil, versionErr
	}
	ret.Patch = int(v)

	_, err = os.Stat(fmt.Sprintf("%s/%s", releaseDir, ret.Version))
	ret.Released = true
	if os.IsNotExist(err) {
		ret.Released = false
	}
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Could not read release directory: %v\n", err)
		return nil, err
	}

	ret.Config, err = readConfig()
	if err != nil {
		return nil, err
	}

	ret.CommitHash, err = gitutil.GetCurrentHash(ret.Config.SourceRoot)
	if err != nil {
		return nil, err
	}
	ret.Name = hashname.HashToName(ret.CommitHash)
	return &ret, nil
}
