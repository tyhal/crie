package runner

import (
	"fmt"
	"regexp"

	"github.com/tyhal/crie/pkg/linter"
)

// Language is used to associate a file pattern with the relevant tools to check and format
type Language struct {
	Regex *regexp.Regexp
	Fmt   linter.Linter
	Chk   linter.Linter
}

// LintType indicates which linter to retrieve for a language (formatter or checker).
type LintType int

// Supported LintType values for selecting a linter.
const (
	// LintTypeFmt selects the formatter linter for a language.
	LintTypeFmt LintType = iota
	// LintTypeChk selects the checker linter for a language.
	LintTypeChk
)

func (lt LintType) String() string {
	switch lt {
	case LintTypeFmt:
		return "fmt"
	case LintTypeChk:
		return "chk"
	default:
		return "unknown"
	}
}

// GetLinter returns the linter associated with the provided LintType
// (either the formatter or the checker) for this language.
func (l *Language) GetLinter(which LintType) (linter.Linter, error) {
	switch which {
	case LintTypeFmt:
		return l.Fmt, nil
	case LintTypeChk:
		return l.Chk, nil
	default:
		return nil, fmt.Errorf("invalid linter type: %d", which)
	}
}
