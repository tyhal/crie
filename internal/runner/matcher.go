// Package runner contains the logic to run the linters
package runner

import (
	"io"
	"regexp"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/tyhal/crie/pkg/linter"
)

// LinterMatch is used to associate a file pattern with the relevant tools to check and format
type LinterMatch struct {
	FileMatch *regexp.Regexp
	Fmt       linter.Linter
	Chk       linter.Linter
}

// NamedMatches store the name for a LinterMatch to make referencing them easier
type NamedMatches map[string]LinterMatch

func getName(l linter.Linter) string {
	if l == nil {
		return ""
	}
	return l.Name()
}

// Show all linters and their associated file types
func (s NamedMatches) Show(w io.Writer) error {
	table := tablewriter.NewWriter(w)
	table.Header([]string{"language", "checker", "formatter", "associated files"})

	// GetFiles sorted language names
	langNames := make([]string, 0, len(s))
	for langName := range s {
		langNames = append(langNames, langName)
	}
	sort.Strings(langNames)

	for _, langName := range langNames {
		l := s[langName]
		err := table.Append([]string{langName, getName(l.Chk), getName(l.Fmt), l.FileMatch.String()})
		if err != nil {
			return err
		}
	}
	err := table.Render()
	return err
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
func (l *LinterMatch) GetLinter(which LintType) linter.Linter {
	switch which {
	case LintTypeFmt:
		return l.Fmt
	case LintTypeChk:
		return l.Chk
	default:
		return nil
	}
}
