package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/version"
	"github.com/alecthomas/kong"
)

type versionFlag string

func (v versionFlag) Decode(ctx *kong.DecodeContext) error {
	return nil
}
func (v versionFlag) IsBool() bool {
	return true
}
func (v versionFlag) BeforeApply(app *kong.Kong) error {
	fmt.Printf("%s: %s\n", version.Number, version.Name)
	app.Exit(1)
	return nil
}
