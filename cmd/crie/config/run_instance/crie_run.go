package run_instance

import (
	"regexp"
	"strings"

	"github.com/tyhal/crie/cmd/crie/config/language"
	"github.com/tyhal/crie/cmd/crie/config/project"
	"github.com/tyhal/crie/pkg/crie"
	"github.com/tyhal/crie/pkg/crie/linter"
)

// TODO stop using globals for dumping config somewhere
var Crie crie.RunConfiguration // 3. final config

// SaveConfiguration pushes the Languages to the crie.RunConfiguration
func SaveConfiguration(proj *project.Config, langs *language.Languages) {
	Crie.Ignore = regexp.MustCompile(strings.Join(proj.Ignore, "|"))

	Crie.Languages = make(map[string]*linter.Language, len(langs.Languages))
	for langName, lang := range langs.Languages {
		Crie.Languages[langName] = lang.ToCrieLanguage()
	}
}
