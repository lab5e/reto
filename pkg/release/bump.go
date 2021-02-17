package release

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lab5e/reto/pkg/toolbox"
)

// BumpVersion bumps the version. Errors are printed on stderr
// The new release context is returned
func BumpVersion(major, minor, patch bool) (*Context, error) {

	ctx, err := GetContext()
	if err != nil {
		return nil, err
	}

	tuples := strings.Split(ctx.Version, ".")
	if len(tuples) != 3 {
		fmt.Printf("%sInvalid version string in version file%s: %s\n", toolbox.Red, toolbox.Reset, ctx.Version)
		return nil, errors.New("invalid version")
	}

	bumps := 0
	if major {
		ctx.Major++
		ctx.Minor = 0
		ctx.Patch = 0
		bumps++
	}

	if minor {
		ctx.Minor++
		ctx.Patch = 0
		bumps++
	}

	if patch {
		ctx.Patch++
		bumps++
	}

	if bumps == 0 {
		fmt.Printf("Must specify which version to bump\n")
		return nil, errors.New("no bump")
	}

	if bumps != 1 {
		fmt.Printf("Only one of bump major, minor or patch can be bumped\n")
		return nil, errors.New("arg error")
	}

	ctx.Version = fmt.Sprintf("%d.%d.%d", ctx.Major, ctx.Minor, ctx.Patch)

	if err := ioutil.WriteFile(VersionFile, []byte(ctx.Version), toolbox.DefaultFilePerm); err != nil {
		fmt.Printf("Error writing version file: %v\n", err)
		return nil, err
	}
	return ctx, nil
}
