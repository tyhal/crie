package settings

import "github.com/tyhal/crie/pkg/crie/linter"

// ConfigLanguage is used to map customer yaml decoders for implementations of Crie linters
type ConfigLanguage struct {
	Regex *ConfigRegex `json:"match,flow,omitempty" yaml:"match,flow,omitempty" jsonschema:"oneof_type=string;array" jsonschema_description:"a regex to match files against to know they are the target of this language"`
	Fmt   ConfigLinter `json:"fmt,omitempty" yaml:"fmt,omitempty" jsonschema_description:"used for the given language when formatting"`
	Chk   ConfigLinter `json:"chk,omitempty" yaml:"chk,omitempty" jsonschema_description:"used for the given language when linting/checking"`
}

func (l ConfigLanguage) toLanguage() *linter.Language {
	return &linter.Language{
		Regex: l.Regex.Regexp,
		Fmt:   l.Fmt.Linter,
		Chk:   l.Chk.Linter,
	}
}
