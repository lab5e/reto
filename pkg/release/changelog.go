package release

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

// DefaultChangeLogTemplate is the default template for the changelog
var DefaultChangeLogTemplate = `
# Changelog {{ version }}: {{ name }}

## Features

[TODO: Write new features]

## API

[TODO: Changes to the API]

## Command line

[TODO: Command line changes]

## Other

[TODO: Write other changes here]
`

// TemplatePath is the path to the changelog template
const TemplatePath = "releases/templates/changelog-template.md"

// WorkingChangelog is the working version of the changelog
const WorkingChangelog = "release/changelog.md"

// ChangelogComplete verifies that there's no TODO statements in the change log
// It will print an error message on stderr witht the line number and return an error
// if one or more TODO strings are found. It's simple but for a reason :)
func ChangelogComplete() error {
	return toolbox.CheckForTODO(WorkingChangelog)
}

// MakeTemplate creates the changelog template file
func makeChangelogTemplate() error {
	if err := ioutil.WriteFile(
		TemplatePath,
		[]byte(DefaultChangeLogTemplate), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Unable to create the release log template: %v", err)
		return err
	}
	return nil
}

// CopyTemplate copies the template into the release folder
func copyChangelogTemplate() error {
	if err := toolbox.CopyFile(TemplatePath, WorkingChangelog); err != nil {
		toolbox.PrintError("Could not copy changelog template to release directory: %v", err)
		return err
	}
	return nil
}

// initChangelog initalises the change log and template
func initChangelog() error {

	// Make template files
	if err := makeChangelogTemplate(); err != nil {
		return err
	}
	if err := copyChangelogTemplate(); err != nil {
		return err
	}

	return nil
}

// Expand the template vars in the working copy of the changelog.
func releaseChangelog(ctx *Context) error {
	buf, err := ioutil.ReadFile("release/changelog.md")
	if err != nil {
		toolbox.PrintError("Could not open working copy of changelog: %v", err)
		return err
	}

	t, err := template.New("changelog").Parse(string(buf))
	if err != nil {
		toolbox.PrintError("Could not parse working copy of changelog: %v", err)
	}

	releaseChangelog := fmt.Sprintf("release/%s/changelog.md", ctx.Version)
	f, err := os.Create(releaseChangelog)
	if err != nil {
		toolbox.PrintError("Could not create changelog at %s: %v", releaseChangelog, err)
		return err
	}
	defer f.Close()

	if err := t.Execute(f, ctx); err != nil {
		toolbox.PrintError("Could not merge template for changelog: %v", err)
		return err
	}
	return nil
}
