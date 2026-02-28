package language

import (
	"errors"

	"github.com/tyhal/crie/internal/runner"
)

// Language is used to map customer yaml decoders for implementations of Crie linters
type Language struct {
	FileMatch *Regex `json:"match,flow" yaml:"match,flow" jsonschema:"a regex to match files against to know they are the target of this language"`
	Fmt       Linter `json:"fmt,omitzero" yaml:"fmt" jsonschema:"used for the given language when formatting"`
	Chk       Linter `json:"chk,omitzero" yaml:"chk" jsonschema:"used for the given language when linting/checking"`
}

// ToRunFormat will convert the yaml friendly version to an internal representation used by crie
func (l Language) ToRunFormat() (runner.LinterMatch, error) {
	if l.FileMatch == nil {
		return runner.LinterMatch{}, errors.New("field match is required")
	}
	return runner.LinterMatch{
		FileMatch: l.FileMatch.Regexp,
		Fmt:       l.Fmt.Linter,
		Chk:       l.Chk.Linter,
	}, nil
}
