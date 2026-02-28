package project

import (
	"bytes"
	"os"

	"github.com/tyhal/crie/internal/runner"
	"gopkg.in/yaml.v3"
)

// LoggingConfig is a grouping of log settings
type LoggingConfig struct {
	Quiet   bool `json:"quiet,omitzero" yaml:"quiet" jsonschema:"disable all output except for errors"`
	Verbose bool `json:"verbose,omitzero" yaml:"verbose" jsonschema:"enable debug logging"`
	Trace   bool `json:"trace,omitzero" yaml:"trace" jsonschema:"(hidden opt) enable all logging (very very verbose)"`
	JSON    bool `json:"json,omitzero" yaml:"json" jsonschema:"change format to json structured logging"`
}

// Config are all the things for crie cli
type Config struct {
	Log    LoggingConfig  `json:"log,omitzero" yaml:"log" jsonschema:"logging options"`
	Dir    string         `json:"dir,omitzero" yaml:"dir" jsonschema:"the directory to run crie in"`
	Lint   runner.Options `json:"lint,omitzero" yaml:"lint" jsonschema:"options for commands that lint"`
	Ignore []string       `json:"ignore,omitzero" yaml:"ignore" jsonschema:"list of regexes matched against the file list to ignore them (exact paths also work)"`
}

// NewProjectConfigFile Creates the project file locally
func (cli *Config) NewProjectConfigFile(path string) error {
	yamlOut, err := yaml.Marshal(cli)

	if err != nil {
		return err
	}

	var buf bytes.Buffer
	// TODO add versioning
	buf.WriteString("# yaml-language-server: $schema=https://raw.githubusercontent.com/tyhal/crie/main/res/schema/proj.json\n")
	buf.Write(yamlOut)
	yamlContent := buf.Bytes()
	err = os.WriteFile(path, yamlContent, 0644)

	if err != nil {
		return err
	}

	return nil
}
