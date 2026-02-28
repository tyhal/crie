package project

import (
	"bytes"
	"os"

	"github.com/tyhal/crie/internal/runner"
	"gopkg.in/yaml.v3"
)

// LoggingConfig is a grouping of log settings
type LoggingConfig struct {
	Quiet   bool `json:"quiet,omitempty" yaml:"quiet,omitempty" jsonschema_description:"disable all output except for errors"`
	Verbose bool `json:"verbose,omitempty" yaml:"verbose,omitempty" jsonschema_description:"enable debug logging"`
	Trace   bool `json:"trace,omitempty" yaml:"trace,omitempty" jsonschema_description:"(hidden opt) enable all logging (very very verbose)"`
	JSON    bool `json:"json,omitempty" yaml:"json,omitempty" jsonschema_description:"change format to json structured logging"`
}

// Config are all the things for crie cli
type Config struct {
	Log    LoggingConfig  `json:"log" yaml:"log,omitempty" jsonschema_description:"logging options"`
	Dir    string         `json:"dir" yaml:"dir,omitempty" jsonschema_description:"the directory to run crie in"`
	Lint   runner.Options `json:"lint" yaml:"lint,omitempty" jsonschema_description:"options for commands that lint"`
	Ignore []string       `json:"ignore,omitempty" yaml:"ignore,omitempty" jsonschema_description:"list of regexes matched against the file list to ignore them (exact paths also work)"`
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
