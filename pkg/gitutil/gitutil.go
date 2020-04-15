package gitutil

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
)

// HasChanges returns true if there's uncomitted or unstaged changes on the
// current branch.
// Using the regular git command here since the Worktree() and Status() methods
// are *really* slow on even medium-sized repositories.
func HasChanges(rootDir string) bool {
	cmd := exec.Command("git", "-C", rootDir, "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		toolbox.PrintError("Could not read Git repo at %s: %v", rootDir, err)
		return true
	}
	lines := strings.Split(out.String(), "\n")
	for _, v := range lines {
		if strings.HasPrefix("??", v) {
			continue
		}
		return true
	}
	return false
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
		toolbox.PrintError("Could not read Git repo at %s: %v", rootDir, err)
		return "", err
	}
	return out.String(), nil
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
