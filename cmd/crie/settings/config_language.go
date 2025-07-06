package settings

import "github.com/tyhal/crie/pkg/crie/linter"

// ConfigLanguage is used to map customer yaml decoders for implementations of Crie linters
type ConfigLanguage struct {
	Regex *ConfigRegex `yaml:"match,flow,omitempty"`
	Fmt   ConfigLinter `yaml:"fmt,omitempty"`
	Chk   ConfigLinter `yaml:"chk,omitempty"`
}

func (l ConfigLanguage) toLanguage() *linter.Language {
	return &linter.Language{
		Regex: l.Regex.Regexp,
		Fmt:   l.Fmt.Linter,
		Chk:   l.Chk.Linter,
	}
}
