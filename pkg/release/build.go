package release

import (
	"fmt"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

func archiveName(ctx *Context, target string) string {
	return fmt.Sprintf("release/archives/%s/%s-%s_%s.zip", ctx.Version, ctx.Version, ctx.Config.Name, target)
}

// BuildRelease builds a release archive for a particular target
func BuildRelease(ctx *Context, target string) error {
	if err := os.Mkdir(fmt.Sprintf("release/%s", ctx.Version), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create release directory: %v", err)
		return err
	}

	ctx.Target = target
	// build the changelog for the release
	if err := releaseChangelog(ctx); err != nil {
		return err
	}

	fmt.Printf("Building release archive %s \n", archiveName(ctx, target))
	for _, v := range ctx.Config.Files {
		if v.Target == anyTarget || v.Target == target {
			fmt.Printf(" - [%s] %s\n", v.ID, v.Name)
		}
	}

	return nil
}
