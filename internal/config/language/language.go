package language

import (
	"github.com/tyhal/crie/internal/runner"
)

// Language is used to map customer yaml decoders for implementations of Crie linters
type Language struct {
	Regex *Regex `json:"match,flow,omitempty" yaml:"match,flow,omitempty" jsonschema:"oneof_type=string;array" jsonschema_description:"a regex to match files against to know they are the target of this language"`
	Fmt   Linter `json:"fmt,omitempty" yaml:"fmt,omitempty" jsonschema_description:"used for the given language when formatting"`
	Chk   Linter `json:"chk,omitempty" yaml:"chk,omitempty" jsonschema_description:"used for the given language when linting/checking"`
}

// ToCrieLanguage will convert the yaml friendly version to an internal representation used by crie
func (l Language) ToCrieLanguage() *runner.Language {
	return &runner.Language{
		Regex: l.Regex.Regexp,
		Fmt:   l.Fmt.Linter,
		Chk:   l.Chk.Linter,
	}
}
