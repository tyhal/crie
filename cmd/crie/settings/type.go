package settings

import (
	"github.com/tyhal/crie/pkg/crie"
	"github.com/tyhal/crie/pkg/crie/linter"
)

// ProjectConfigFile is the schema for a projects' settings file
type ProjectConfigFile struct {
	Languages []linter.Language `yaml:"languages"`
	Ignore    []string          `yaml:"ignore"`
}

// CliSettings are all the things for crie cli
type CliSettings struct {
	Crie          crie.RunConfiguration
	ProjectConfig ProjectConfigFile
	ConfigPath    string
	Quiet         bool
	Verbose       bool
	Trace         bool
	JSON          bool
}
