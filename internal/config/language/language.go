package language

import (
	"errors"

	"github.com/tyhal/crie/internal/runner"
)

// Language is used to map customer yaml decoders for implementations of Crie linters
type Language struct {
	FileMatch *Regex `json:"match" yaml:"match,flow" jsonschema:"a regex to match files against to know they are the target of this language"`
	GroupBy   string `json:"group_by,omitzero" yaml:"group_by" jsonschema:"module marker filename; when set, matched files are grouped by nearest ancestor directory containing this file (e.g. go.mod) and the tool runs once per group"`
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
		GroupBy:   l.GroupBy,
		Fmt:       l.Fmt.Linter,
		Chk:       l.Chk.Linter,
	}, nil
}
