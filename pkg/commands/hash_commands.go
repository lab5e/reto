package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lab5e/reto/pkg/gitutil"
	"github.com/lab5e/reto/pkg/hashname"
	"github.com/lab5e/reto/pkg/release"
	"github.com/lab5e/reto/pkg/toolbox"
)

type hashCommand struct {
	Version string `kong:"short='v',help='Version number to use. Will use last commit otherwise'"`
}

type hashNameCommand struct {
	Version string `kong:"short='v',help='Version number to use. Will use last commit otherwise'"`
}

func (c *hashCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}
	if c.Version == "" {
		fmt.Println(ctx.CommitHash)
		return nil
	}

	hash, err := gitutil.GetHash(ctx.Config.SourceRoot, c.Version)
	if err != nil {
		toolbox.PrintError(err.Error())
		return err
	}
	fmt.Println(hash)
	return nil
}

func (c *hashNameCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}
	if c.Version == "" {
		fmt.Println(ctx.Name)
		return nil
	}

	hash, err := gitutil.GetHash(ctx.Config.SourceRoot, c.Version)
	if err != nil {
		toolbox.PrintError(err.Error())
		return err
	}
	fmt.Println(hashname.HashToName(hash))
	return nil
}

type nameHashCommand struct {
	Name string `kong:"arg,help='Name to convert into hash'"`
}

func (c *nameHashCommand) Run(rc RunContext) error {
	hash, err := hashname.NameToHash(rc.ReleaseCommands().Namehash.Name)
	if err != nil {
		toolbox.PrintError("Unable to convert name into hash: %v", err)
		return err
	}
	fmt.Println(hash)
	return nil
}

type nameVersionCommand struct {
	Name string `kong:"arg, help='Look up version from hash name'"`
}

func (c *nameVersionCommand) Run(rc RunContext) error {
	hash, err := hashname.NameToHash(rc.ReleaseCommands().Nameversion.Name)
	if err != nil {
		toolbox.PrintError("Unable to convert name into hash: %v", err)
		return err
	}

	ctx, err := release.GetContext()
	if err != nil {
		return err
	}

	tagDir := ".git/refs/tags"
	if ctx.Config.SourceRoot != "" {
		tagDir = fmt.Sprintf("%s/%s", ctx.Config.SourceRoot, tagDir)
	}
	infos, err := ioutil.ReadDir(tagDir)
	if err != nil {
		toolbox.PrintError("Could not read tag directory in %s", ctx.Config.SourceRoot)
		return err
	}

	for _, fi := range infos {
		if fi.IsDir() {
			continue
		}
		buf, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", tagDir, fi.Name()))
		if err != nil {
			toolbox.PrintError("Could not read hash for tag %s: %v", fi.Name(), err)
			return err
		}
		if strings.HasPrefix(string(buf), hash) {
			fmt.Println(fi.Name()[1:])
			return nil
		}
	}
	return errors.New("could not find a matching version")
	// Find the mapped tag version
}
