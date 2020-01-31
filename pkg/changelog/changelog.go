package changelog

import (
	"errors"
	"io/ioutil"

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

const TemplatePath = "release/templates/changelog-template.md"
const WorkingChangelog = "release/changelog.md"

// MakeTemplate creates the changelog template file
func MakeTemplate() error {
	if err := ioutil.WriteFile(
		TemplatePath,
		[]byte(DefaultChangeLogTemplate), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Unable to create the release log template: %v", err)
		return err
	}
	return nil
}

// CopyTemplate copies the template into the release folder
func CopyTemplate() error {
	if err := toolbox.CopyFile(TemplatePath, WorkingChangelog); err != nil {
		toolbox.PrintError("Could not copy changelog template to release directory: %v", err)
		return err
	}
	return nil
}

// WorkingCopyIsComplete verifies that there's no TODO statements in the change log
// It will print an error message on stderr witht the line number and return an error
// if one or more TODO strings are found. It's simple but for a reason :)
func WorkingCopyIsComplete() error {
	return toolbox.CheckForTODO(WorkingChangelog)
}

// ReleaseChangelog makes a copy of the working copy, puts it in the release directory
// and copies the template in to the working copy
func ReleaseChangelog() error {
	return errors.New("not implemented")
}
