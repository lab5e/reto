package release

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"sort"

	"github.com/lab5e/reto/pkg/toolbox"
)

// sampleChangelog is the default template for the changelog
var sampleChangelog = `## v{{ .Version }}: {{ .Name }}

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

// TemplatesComplete verifies that there's no TODO statements in the change log
// It will print an error message on stderr witht the line number and return an error
// if one or more TODO strings are found. It's simple but for a reason :)
func TemplatesComplete(ctx *Context, printErrors bool) error {
	errs := 0
	for _, template := range ctx.Config.Templates {
		workingCopy := fmt.Sprintf("release/%s", template.Name)
		if err := toolbox.CheckForTODO(workingCopy, printErrors); err != nil {
			errs++
		}
	}
	if errs > 0 {
		return errors.New("incomplete")
	}
	return nil
}

func initTemplates() error {
	templateFile := fmt.Sprintf("%s/changelog.md", templateDir)
	if err := ioutil.WriteFile(templateFile, []byte(sampleChangelog), toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not create sample file %s", templateFile)
	}

	for _, template := range defaultConfig().Templates {
		templateFile := fmt.Sprintf("%s/%s", templateDir, template.Name)
		workingCopy := fmt.Sprintf("release/%s", template.Name)
		if err := toolbox.CopyFile(templateFile, workingCopy); err != nil {
			toolbox.PrintError("Could not copy %s to release directory: %v", workingCopy, err)
			return err
		}
	}
	return nil
}

// Expand the template vars in the working copy of the changelog.
func mergeTemplate(source string, dest string, ctx *Context) error {
	buf, err := ioutil.ReadFile(source)
	if err != nil {
		toolbox.PrintError("Could not open working copy of template: %v", err)
		return err
	}

	t, err := template.New(source).Parse(string(buf))
	if err != nil {
		toolbox.PrintError("Could not parse working copy of template: %v", err)
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		toolbox.PrintError("Could not create template at %s: %v", dest, err)
		return err
	}
	defer f.Close()

	if err := t.Execute(f, ctx); err != nil {
		toolbox.PrintError("Could not merge template: %v", err)
		return err
	}
	return nil
}

// MergeChangelogs merges all of the changelogs into one big buffer
func concatenateTemplate(basename string, destination string) error {
	fmt.Printf("Merging %s...\n", basename)

	fileinfos, err := ioutil.ReadDir(releaseDir)
	if err != nil {
		toolbox.PrintError("Could not read release directory: %v", err)
		return err
	}

	inputs := make([]string, 0)
	for _, fi := range fileinfos {
		templatePath := fmt.Sprintf("%s/%s/%s", releaseDir, fi.Name(), basename)
		if fi.IsDir() && toolbox.IsFile(templatePath) {
			inputs = append(inputs, templatePath)
		}
	}
	sort.Strings(inputs)
	var buffer []byte
	for i := len(inputs) - 1; i >= 0; i-- {
		buf, err := ioutil.ReadFile(inputs[i])
		if err != nil {
			toolbox.PrintError("Could not read changelog at %s: %v", inputs[i], err)
			return err
		}
		buffer = append(buffer, buf...)
	}
	fmt.Printf("%d files merged\n", len(inputs))
	if err := ioutil.WriteFile(destination, buffer, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Error writing %s: %v", destination, err)
		return err
	}
	return nil
}
