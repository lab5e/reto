package toolbox

import "os"

// IsFile returns true if the file exists, false on error
func IsFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
