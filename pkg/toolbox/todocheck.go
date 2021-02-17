package toolbox

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// CheckForTODO checks if the string TODO is somewhere in a file. Prints the
// line number and file name if so.
func CheckForTODO(file string, printErrors bool) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("%sCould not read %s%s: %v", Red, file, Reset, err)
		return err
	}

	lines := strings.Split(string(buf), "\n")
	todos := 0
	for i, v := range lines {
		if strings.Index(v, "TODO") > 0 {
			if printErrors {
				fmt.Printf("%s%s%s: TODO statement on line %s%d%s", Cyan, file, Reset, Cyan, i+1, Reset)
			}
			todos++
		}
	}
	if todos > 0 {
		return errors.New("not completed")
	}
	return nil
}
