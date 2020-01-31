package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/releasetool/pkg/changelog"
	"github.com/ExploratoryEngineering/releasetool/pkg/release"
)

type statusCommand struct {
}

func (c *statusCommand) Run(rc RunContext) error {
	ctx, err := release.Verify()
	if err != nil {
		return err
	}

	released := "NO"
	if ctx.Released {
		released = "YES"
	}

	fmt.Printf("Active version: %s\n", ctx.Version)
	fmt.Printf("Released:       %s\n", released)
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

	if changelog.WorkingCopyIsComplete() == nil {
		fmt.Println("Changelog is OK")
	}

	if release.VerifyConfig(ctx.Config) == nil {
		fmt.Println("Configuration is OK")
	}
	return nil
}
