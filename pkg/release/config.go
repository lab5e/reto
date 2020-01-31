package release

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

// File is the configuration setting for a single file
type File struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Target string `json:"target"`
}

// Config is the tool configuration
type Config struct {
	SourceRoot string   `json:"sourceRoot"`
	Name       string   `json:"name"`
	Targets    []string `json:"targets"`
	Files      []File   `json:"files"`
}

// ConfigPath is the path to the configuration file
const ConfigPath = "release/config.json"

// WriteSampleConfig writes a sample configuration to the release directory
func WriteSampleConfig() error {
	_, err := os.Stat(ConfigPath)
	if !os.IsNotExist(err) {
		toolbox.PrintError("Configuration file already exists")
		return err
	}

	c := sampleConfig()
	buf, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(ConfigPath, buf, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not write sample config: %v", err)
		return err
	}
	return nil
}

// sampleConfig is the sample configuration file.
func sampleConfig() Config {
	return Config{
		SourceRoot: ".",
		Name:       "TODO set your product name",
		Targets:    []string{"TODO: set target (amd64-darwin, arm-linux, mips-plan9...)"},
		Files: []File{
			File{
				ID:     "TODO: set ID for file",
				Name:   "TODO: Add your built files here",
				Target: "TODO: Set target for file here, '-' if it doesn't apply",
			},
		},
	}
}

func readConfig() (Config, error) {
	buf, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		toolbox.PrintError("Could not read configuration: %v", err)
		return Config{}, err
	}
	ret := Config{}
	if err := json.Unmarshal(buf, &ret); err != nil {
		toolbox.PrintError("Configuration file format error: %v", err)
		return Config{}, err
	}
	return ret, nil
}

const anyTarget = "-"

// VerifyConfig verifies that the artifact config is correct
func VerifyConfig(config Config) error {
	if err := toolbox.CheckForTODO(ConfigPath); err != nil {
		return err
	}
	if len(config.Targets) == 0 {
		toolbox.PrintError("There are no targets in the configuration file")
		return errors.New("no targets")
	}
	if len(config.Files) == 0 {
		toolbox.PrintError("There are no output files in the configuration file")
		return errors.New("no targets")
	}

	fileTargets := make(map[string]map[string]bool)
	for _, v := range config.Files {
		if v.Target == anyTarget {
			continue
		}
		existing, ok := fileTargets[v.ID]
		if !ok {
			existing = make(map[string]bool)
		}
		existing[v.Target] = true
		fileTargets[v.ID] = existing
	}

	errs := 0
	for id, v := range fileTargets {
		targets := make(map[string]bool)
		for _, t := range config.Targets {
			if t == anyTarget {
				continue
			}
			targets[t] = true
		}
		for target := range v {
			if target == anyTarget {
				continue
			}
			_, ok := targets[target]
			if !ok {
				toolbox.PrintError("File with ID '%s' have unknown target %s", id, target)
				errs++
			}
			delete(targets, target)
		}
		if len(targets) > 0 {
			for target := range targets {
				toolbox.PrintError("File with ID '%s' is missing target %s", id, target)
				errs++
			}
		}
	}
	if errs > 0 {
		return errors.New("target missing")
	}
	return nil
}
