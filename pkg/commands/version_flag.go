package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/reto/pkg/version"
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
	fmt.Printf("%s: %s (%s)\n", version.Number, version.Name, version.BuildTime)
	app.Exit(1)
	return nil
}
