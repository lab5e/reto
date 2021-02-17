package commands

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/lab5e/reto/pkg/version"
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
