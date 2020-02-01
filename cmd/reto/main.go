package main

import (
	"os"

	"github.com/ExploratoryEngineering/reto/pkg/commands"
	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
	"github.com/alecthomas/kong"
)

func main() {
	/*defer func() {
		// The Kong parser panics when there's a sole dash in the argument list
		// I'm not sure if this is a bug or a feature :)
		if r := recover(); r != nil {
			fmt.Println("Error parsing command line: ", r)
		}
	}()*/

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
