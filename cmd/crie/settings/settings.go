package settings

import (
	"github.com/tyhal/crie/pkg/crie"
	"github.com/tyhal/crie/pkg/crie/linter"
	"regexp"
	"strings"
)

// Cli is the entire current Crie Cli settings
var Cli CliSettings

// CliSettings are all the things for crie cli
type CliSettings struct {
	ConfigProject ConfigProject         // 1. defaults loaded
	ConfigPath    string                // 2. user config overrides defaults
	Crie          crie.RunConfiguration // 3. final config
	Quiet         bool
	Verbose       bool
	Trace         bool
	JSON          bool
}

// SaveConfiguration pushes the ConfigProject to the crie.RunConfiguration
func (cli *CliSettings) SaveConfiguration() {
	cli.Crie.Ignore = regexp.MustCompile(strings.Join(cli.ConfigProject.Ignore, "|"))

	cli.Crie.Languages = make(map[string]*linter.Language)
	for langName, lang := range cli.ConfigProject.Languages {
		cli.Crie.Languages[langName] = lang.toLinter()
	}
}
