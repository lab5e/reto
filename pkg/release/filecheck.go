package release

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lab5e/reto/pkg/toolbox"
)

// NewFileVersions checks that the binaries/artifacts are changed since last
// release. This ignores the common artifacts (with ID set to '-') and will
// only check the actual binaries.
func NewFileVersions(config Config, printErrors bool) bool {
	// This holds the checksums for the new files
	newChecksums := make(map[string]string)
	for _, v := range config.Files {
		if v.Target == anyTarget {
			// ignore common files. These will be the same from release to release
			continue
		}

		buf, err := ioutil.ReadFile(v.Name)
		if err != nil {
			fmt.Printf("%sUnable to read %s%s: %v\n", toolbox.Red, v.Name, toolbox.Reset, err)
			return false
		}
		sum := sha256.Sum256(buf)

		newChecksums[filepath.Base(v.Name)] = fmt.Sprintf("%x", sum)
	}

	// Look in releases for the previous version
	var releasedVersions []string

	fileinfos, err := ioutil.ReadDir(archiveDir)
	if err != nil {
		fmt.Printf("%sCould not read release directory%s: %v\n", toolbox.Red, toolbox.Reset, err)
		return false
	}

	for _, fi := range fileinfos {
		checksumFile := checksumFileName(config.Name, fi.Name())
		if fi.IsDir() && toolbox.IsFile(checksumFile) {
			releasedVersions = append(releasedVersions, fi.Name())
		}
	}
	if len(releasedVersions) == 0 {
		fmt.Println("Note: Found no old versions")
		return true
	}
	sort.Strings(releasedVersions)

	oldVersion := releasedVersions[len(releasedVersions)-1]

	// open the checksum file and extract the checksums for each file
	oldChecksumFile := checksumFileName(config.Name, oldVersion)

	buf, err := ioutil.ReadFile(oldChecksumFile)
	if err != nil {
		fmt.Printf("%sCould not read previous checksum file %s%s: %v\n", toolbox.Red, oldChecksumFile, toolbox.Reset, err)
		return false
	}
	lines := strings.Split(string(buf), "\n")

	sameChecksum := 0
	for _, v := range lines {
		tuples := strings.Split(strings.TrimSpace(v), "  ")
		if len(tuples) != 2 {
			continue
		}
		newChecksum := newChecksums[tuples[1]]
		if newChecksum == tuples[0] {
			if printErrors {
				fmt.Printf("%sFile %s%s has the same checksum as the previous version (%s) -- %s", toolbox.Red, tuples[1], toolbox.Reset, oldVersion, tuples[0])
			}
			sameChecksum++
		}
	}

	return sameChecksum == 0
}
