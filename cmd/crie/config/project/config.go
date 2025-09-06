package project

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

// LoggingConfig is a grouping of log settings
type LoggingConfig struct {
	Quiet   bool `json:"quiet" yaml:"quiet" jsonschema_description:"disable all output except for errors"`
	Verbose bool `json:"verbose" yaml:"verbose" jsonschema_description:"enable debug logging"`
	Trace   bool `json:"trace" yaml:"trace" jsonschema_description:"(hidden opt) enable all logging (very very verbose)"`
	JSON    bool `json:"json" yaml:"json" jsonschema_description:"change format to json structured logging"`
}

// LintConfig is a grouping of settings that affect crie directly
type LintConfig struct {
	Continue      bool   `json:"continue" yaml:"continue"`
	Passes        bool   `json:"passes" yaml:"passes"`
	GitTarget     string `json:"gitTarget" yaml:"gitTarget"`
	GitDiff       bool   `json:"gitDiff" yaml:"gitDiff"`
	Lang          string `json:"lang" yaml:"lang"`
	StrictLogging bool   `json:"-" yaml:"-"`
}

// Config are all the things for crie cli
type Config struct {
	Log    LoggingConfig `json:"log" yaml:"log" jsonschema_description:"logging options"`
	Lint   LintConfig    `json:"lint" yaml:"lint" jsonschema_description:"options for commands that lint"`
	Ignore []string      `json:"ignore" yaml:"ignore" jsonschema_description:"list of regexes matched against the file list to ignore them (exact paths also work)"`
}

// NewProjectConfigFile Creates the project file locally
func (cli *Config) NewProjectConfigFile(path string) error {
	yamlOut, err := yaml.Marshal(cli)

	if err != nil {
		return err
	}

	var buf bytes.Buffer
	// TODO add versioning
	buf.WriteString("# yaml-language-server: $schema=https://raw.githubusercontent.com/tyhal/crie/main/doc/schema_proj.json\n")
	buf.Write(yamlOut)
	yamlContent := buf.Bytes()
	err = os.WriteFile(path, yamlContent, 0644)

	if err != nil {
		return err
	}

	return nil
}
