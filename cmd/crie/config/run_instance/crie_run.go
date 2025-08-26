package run_instance

import (
	"github.com/tyhal/crie/cmd/crie/config/language"
	"github.com/tyhal/crie/cmd/crie/config/project"
	"github.com/tyhal/crie/pkg/crie"
	"github.com/tyhal/crie/pkg/crie/linter"
	"regexp"
	"strings"
)

// TODO stop using globals for dumping config somewhere
var Crie crie.RunConfiguration // 3. final config

// SaveConfiguration pushes the ConfigLanguages to the crie.RunConfiguration
func SaveConfiguration(proj *project.ConfigProject, langs *language.ConfigLanguages) {
	Crie.Ignore = regexp.MustCompile(strings.Join(proj.Ignore, "|"))

	Crie.Languages = make(map[string]*linter.Language)
	for langName, lang := range langs.Languages {
		Crie.Languages[langName] = lang.ToCrieLanguage()
	}
}
