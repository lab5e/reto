package release

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

// DefaultChangeLogTemplate is the default template for the changelog
var DefaultChangeLogTemplate = `
## v{{ .Version }}: {{ .Name }}

### Features

[TODO: Write new features]

### API

[TODO: Changes to the API]

### Command line

[TODO: Command line changes]

### Other

[TODO: Write other changes here]

Commit hash: {{ .CommitHash }}
`

// TemplatePath is the path to the changelog template
const TemplatePath = "release/templates/changelog-template.md"

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
		return err
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

// MergeChangelogs merges all of the changelogs into one big buffer
func MergeChangelogs() ([]byte, error) {
	fmt.Println("Merging changelogs...")

	fileinfos, err := ioutil.ReadDir("release")
	if err != nil {
		toolbox.PrintError("Could not read release directory: %v", err)
		return nil, err
	}

	inputs := make([]string, 0)
	for _, fi := range fileinfos {
		changelogPath := fmt.Sprintf("release/%s/changelog.md", fi.Name())
		if fi.IsDir() && toolbox.IsFile(changelogPath) {
			inputs = append(inputs, changelogPath)
		}
	}
	sort.Strings(inputs)
	var ret []byte
	for i := len(inputs) - 1; i >= 0; i-- {
		fmt.Println(inputs[i])
		ret = append(ret, []byte("\n\n")...)
		buf, err := ioutil.ReadFile(inputs[i])
		if err != nil {
			toolbox.PrintError("Could not read changelog at %s: %v", inputs[i], err)
			return nil, err
		}
		ret = append(ret, buf...)
	}
	ret = append(ret, byte('\n'))

	// Remove multi-line separators just to make markdown linters happy
	consecutiveSpaces := 0
	mergedFile := "# Changelog"

	for _, line := range strings.Split(string(ret), "\n") {
		if strings.TrimSpace(line) == "" {
			consecutiveSpaces++
		} else {
			consecutiveSpaces = 0
		}
		if consecutiveSpaces > 1 {
			continue
		}
		mergedFile = mergedFile + "\n" + line
	}
	fmt.Printf("%d changelogs merged\n", len(inputs))
	return []byte(mergedFile), nil
}
