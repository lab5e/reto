package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lab5e/reto/pkg/commands"
	"github.com/lab5e/reto/pkg/toolbox"
)

func main() {
	var param commands.Root
	ctx := kong.Parse(&param,
		kong.Name("reto"),
		kong.Description("Release tool"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			NoAppSummary: true,
			Summary:      true,
		}), kong.BindTo(&param, (*commands.RunContext)(nil)))

	if param.Root != "" {
		if err := os.Chdir(param.Root); err != nil {
			toolbox.PrintError("Couldn't change directory to %s", param.Root)
			os.Exit(1)
		}
	}
	if err := ctx.Run(); err != nil {
		// Won't print the error since the commands will do it
		os.Exit(1)
	}
}
