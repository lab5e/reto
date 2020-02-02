package commands

import (
	"errors"
	"fmt"

	"github.com/ExploratoryEngineering/reto/pkg/gitutil"
	"github.com/ExploratoryEngineering/reto/pkg/release"
)

type statusCommand struct {
	Verbose    bool `kong:"short='V',help='Verbose output'"`
	NoBinCheck bool `kong:"short='c',help='Skip change check on release files'"`
}

func okNotOK(v bool) string {
	if v {
		return "OK"
	}
	return "NOT OK"
}

func (c *statusCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}

	fmt.Println("Checking templates")
	templateErr := release.TemplatesComplete(ctx, rc.ReleaseCommands().Status.Verbose)
	fmt.Println("Checking config")
	configErr := release.VerifyConfig(ctx.Config, rc.ReleaseCommands().Status.Verbose)
	fmt.Println("Checking old releases")
	changedFiles := release.NewFileVersions(ctx.Config, rc.ReleaseCommands().Status.Verbose)
	fmt.Println("Checking SCM")
	gitChanges := !gitutil.HasChanges(ctx.Config.SourceRoot)

	fmt.Println()
	fmt.Printf("Configuration:       %s\n", okNotOK(configErr == nil))
	fmt.Printf("Working templates:   %s\n", okNotOK(templateErr == nil))
	fmt.Printf("Version number:      %s\n", okNotOK(!ctx.Released))
	fmt.Printf("Uncommitted changes: %s\n", okNotOK(gitChanges))
	fmt.Printf("Changed binaries:    %s\n", okNotOK(changedFiles))
	fmt.Println()
	fmt.Printf("Active version:      %s\n", ctx.Version)
	fmt.Printf("Commit Hash:         %s\n", ctx.CommitHash)
	fmt.Printf("Name:                %s\n", ctx.Name)

	readyToRelease := configErr == nil && templateErr == nil && gitChanges && changedFiles

	if rc.ReleaseCommands().Status.Verbose {
		fmt.Println()
		fmt.Println("Configuration:")
		fmt.Println("  Targets: ")
		for _, v := range ctx.Config.Targets {
			fmt.Printf("  - %s\n", v)
		}
		fmt.Println("  Files:")
		for _, v := range ctx.Config.Files {
			fmt.Printf("  - %s/%s\n", v.Name, v.Target)
		}
		fmt.Println()
	}

	if readyToRelease {
		fmt.Println()
		fmt.Println("Ready to release.")
		return nil
	}

	return errors.New("notready")
}
