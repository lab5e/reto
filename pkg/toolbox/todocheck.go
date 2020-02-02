package toolbox

import (
	"errors"
	"io/ioutil"
	"strings"
)

// Check if the string TODO is somewhere in a file. Prints the line number
// and file name if so.
func CheckForTODO(file string, printErrors bool) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		PrintError("Could not read %s: %v", file, err)
		return err
	}

	lines := strings.Split(string(buf), "\n")
	todos := 0
	for i, v := range lines {
		if strings.Index(v, "TODO") > 0 {
			if printErrors {
				PrintError("%s: TODO statement on line %d", file, i+1)
			}
			todos++
		}
	}
	if todos > 0 {
		return errors.New("not completed")
	}
	return nil
}
