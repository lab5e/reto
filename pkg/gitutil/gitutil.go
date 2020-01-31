package gitutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// GetCurrentHash returns the current hash for HEAD by digging through
// the .git directory. The hash is stored somewhere in .git/refs/heads and
// the file .git/HEAD points to the current branch
func GetCurrentHash(rootDir string) (string, error) {
	buf, err := ioutil.ReadFile(filepath.Join(rootDir, ".git/HEAD"))
	if err != nil {
		return "", err
	}
	refs := strings.Split(strings.TrimSpace(string(buf)), " ")
	if len(refs) != 2 {
		return "", errors.New("can't grok .git/HEAD")
	}

	hashBuf, err := ioutil.ReadFile(fmt.Sprintf(".git/%s", refs[1]))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(hashBuf)), nil
}
