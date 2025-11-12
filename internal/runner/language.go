// Package runner contains the logic to run the linters
package runner

import (
	"fmt"
	"io"
	"regexp"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/tyhal/crie/pkg/linter"
)

// Language is used to associate a file pattern with the relevant tools to check and format
type Language struct {
	FileMatch *regexp.Regexp
	Fmt       linter.Linter
	Chk       linter.Linter
}

// Languages store the name to a singular language configuration within crie
type Languages map[string]*Language

// Show to print all languages chkConf fmt and always commands
func (s Languages) Show(w io.Writer) error {
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
