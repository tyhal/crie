package settings

import (
	"github.com/tyhal/crie/pkg/crie/linter"
)

// ConfigLanguage is used to map customer yaml decoders for implementations of Crie linters
type ConfigLanguage struct {
	Regex *ConfigRegex `yaml:"match,flow,omitempty"`
	Fmt   ConfigLinter `yaml:"fmt,omitempty"`
	Chk   ConfigLinter `yaml:"chk,omitempty"`
}

func (l ConfigLanguage) toLinter() *linter.Language {
	return &linter.Language{
		Regex: l.Regex.Regexp,
		Fmt:   l.Fmt.Linter,
		Chk:   l.Chk.Linter,
	}
}

// ConfigProject is the schema for a projects' settings file
type ConfigProject struct {
	Languages map[string]ConfigLanguage `yaml:"languages"`
	Ignore    []string                  `yaml:"ignore"`
}

func (c *ConfigProject) merge(src ConfigProject) {

	for langName, lang := range src.Languages {
		if existing, exists := c.Languages[langName]; exists {
			c.Languages[langName] = existing
		} else {
			c.Languages[langName] = lang
		}
	}

	c.Ignore = append(c.Ignore, src.Ignore...)
}
