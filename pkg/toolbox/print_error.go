package toolbox

import (
	"fmt"
	"os"
)

// PrintError prints an error message to stderr
func PrintError(fmtString string, opts ...interface{}) {
	if opts == nil {
		fmt.Fprintf(os.Stderr, "reto: %s\n", fmt.Sprintf(fmtString, opts...))
		return
	}
	fmt.Fprintf(os.Stderr, "reto: %s\n", fmt.Sprintf(fmtString, opts...))
}
