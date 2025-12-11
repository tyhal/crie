package language

import (
	"errors"

	"github.com/tyhal/crie/internal/runner"
)

// Language is used to map customer yaml decoders for implementations of Crie linters
type Language struct {
	FileMatch *Regex `json:"match,flow" yaml:"match,flow" jsonschema:"oneof_type=string;array" jsonschema_required:"true" jsonschema_description:"a regex to match files against to know they are the target of this language"`
	Fmt       Linter `json:"fmt,omitempty" yaml:"fmt,omitempty" jsonschema_description:"used for the given language when formatting"`
	Chk       Linter `json:"chk,omitempty" yaml:"chk,omitempty" jsonschema_description:"used for the given language when linting/checking"`
}

// ToRunFormat will convert the yaml friendly version to an internal representation used by crie
func (l Language) ToRunFormat() (*runner.Language, error) {
	if l.FileMatch == nil {
		return nil, errors.New("field match is required")
	}
	return &runner.Language{
		FileMatch: l.FileMatch.Regexp,
		Fmt:       l.Fmt.Linter,
		Chk:       l.Chk.Linter,
	}, nil
}
