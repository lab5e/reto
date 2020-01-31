package release

import "fmt"

func archiveName(ctx *Context, target string) string {
	return fmt.Sprintf("release/archives/%s/%s-%s_%s.zip", ctx.Version, ctx.Version, ctx.Config.Name, target)
}

// BuildRelease builds a release archive for a particular target
func BuildRelease(ctx *Context, target string) {
	fmt.Printf("Building release archive %s \n", archiveName(ctx, target))
	for _, v := range ctx.Config.Files {
		if v.Target == anyTarget || v.Target == target {
			fmt.Printf(" - [%s] %s\n", v.ID, v.Name)
		}
	}
	// build the changelog for the release
}
