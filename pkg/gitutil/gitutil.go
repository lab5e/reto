package gitutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/lab5e/reto/pkg/toolbox"
)

// HasChanges returns true if there's uncomitted or unstaged changes on the
// current branch.
// Using the regular git command here since the Worktree() and Status() methods
// are *really* slow on even medium-sized repositories.
func HasChanges(rootDir string, verbose bool) bool {
	cmd := exec.Command("git", "-C", rootDir, "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sCould not read Git repo at %s%s: %v\n", toolbox.Red, rootDir, toolbox.Reset, err)
		return true
	}
	lines := strings.Split(out.String(), "\n")
	ret := false
	for _, v := range lines {
		if strings.TrimSpace(v) == "" {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(v), "??") {
			continue
		}
		if verbose {
			fmt.Printf("%sUncommitted changes%s: %s\n", toolbox.Red, toolbox.Reset, strings.TrimSpace(v))
			ret = true
		}
	}
	return ret
}

// GetCurrentHash returns the current hash for HEAD by digging through
// the .git directory. The hash is stored somewhere in .git/refs/heads and
// the file .git/HEAD points to the current branch
func GetCurrentHash(rootDir string) (string, error) {
	cmd := exec.Command("git", "-C", rootDir, "rev-parse", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sCould not read Git repo at %s%s: %v\n", toolbox.Red, rootDir, toolbox.Reset, err)
		return "", err
	}
	return out.String(), nil
}

// GetHash returns the hash for the given version. The hash is found by
// reading .git/refs/tags/<version>. If the file isn't found it will return
// an error
func GetHash(rootDir, version string) (string, error) {
	tagFile := fmt.Sprintf(".git/refs/tags/v%s", version)
	if rootDir != "" {
		tagFile = fmt.Sprintf("%s/%s", rootDir, tagFile)
	}
	buf, err := ioutil.ReadFile(tagFile)
	if err != nil {
		fmt.Printf("%sCould not find a version named%s %s%s%s in %s%s%s\n",
			toolbox.Red, toolbox.Reset, toolbox.Cyan, version, toolbox.Reset, toolbox.Cyan, rootDir, toolbox.Reset)
		return "", err
	}
	return strings.TrimSpace(string(buf)), nil
}

// TagVersion creates a version tag in Git
func TagVersion(rootDir, tagName, message string) error {
	cmd := exec.Command("git", "-C", rootDir, "tag", tagName, "-m", message)
	return cmd.Run()
}

// CreateCommit creates a new commit.
func CreateCommit(rootDir, message string, files ...string) (string, error) {
	for _, v := range files {
		cmd := exec.Command("git", "-C", rootDir, "add", v)
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	err := exec.Command("git", "-C", rootDir, "commit", "-m", message).Run()
	if err != nil {
		return "", err
	}
	return GetCurrentHash(rootDir)
}
