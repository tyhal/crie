package printer

// github.com/errata-ai/vale@v1.7.1/ui/color.go
// COPY of above (MIT) licensed project
// CHANGES:
// Stopped printing to stdout
// Reduced formatting
// Returns err not bool

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/errata-ai/vale/core"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io"
	"regexp"
	"strings"
)

const (
	errorColor      = color.FgRed
	warningColor    = color.FgYellow
	suggestionColor = color.FgBlue
)

func pluralize(s string, n int) string {
	if n != 1 {
		return s + "s"
	}
	return s
}

var spaces = regexp.MustCompile(" +")

func fixOutputSpacing(msg string) string {
	msg = strings.Replace(msg, "\n", " ", -1)
	msg = spaces.ReplaceAllString(msg, " ")
	return msg
}

// GetVerboseAlerts prints Alerts in verbose format.
func GetVerboseAlerts(linted []*core.File, wrap bool) (io.Reader, error) {
	var lintErrors, lintWarnings, lintSuggestions int
	var e, w, s int
	var symbol string
	var err error
	buf := new(bytes.Buffer)

	for _, f := range linted {
		e, w, s = printVerboseAlert(f, wrap, buf)
		lintErrors += e
		lintWarnings += w
		lintSuggestions += s
	}

	etotal := fmt.Sprintf("%d %s", lintErrors, pluralize("error", lintErrors))
	wtotal := fmt.Sprintf("%d %s", lintWarnings, pluralize("warning", lintWarnings))
	stotal := fmt.Sprintf("%d %s", lintSuggestions, pluralize("suggestion", lintSuggestions))

	if lintErrors > 0 {
		err = errors.New("linting errors found")
	}

	if lintErrors > 0 || lintWarnings > 0 {
		symbol = "\u2716"
	} else {
		symbol = "\u2714"
	}

	n := len(linted)
	_, formatErr := fmt.Fprintf(buf, "%s %s, %s and %s in %d %s.\n", symbol,
		colorize(etotal, errorColor), colorize(wtotal, warningColor),
		colorize(stotal, suggestionColor), n, pluralize("file", n))

	if formatErr != nil {
		err = formatErr
	}

	return buf, err
}

// printVerboseAlert includes an alert's line, column, level, and message.
func printVerboseAlert(f *core.File, wrap bool, writer io.Writer) (int, int, int) {
	var loc, level string
	var errorCount, warningCount, notifyCount int

	alerts := f.SortedAlerts()
	if len(alerts) == 0 {
		return 0, 0, 0
	}

	table := tablewriter.NewWriter(writer)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetAutoWrapText(!wrap)

	for _, a := range alerts {
		a.Message = fixOutputSpacing(a.Message)
		if a.Severity == "suggestion" {
			level = colorize(a.Severity, suggestionColor)
			notifyCount++
		} else if a.Severity == "warning" {
			level = colorize(a.Severity, warningColor)
			warningCount++
		} else {
			level = colorize(a.Severity, errorColor)
			errorCount++
		}
		loc = fmt.Sprintf("%d:%d", a.Line, a.Span[0])
		table.Append([]string{loc, level, a.Message, a.Check})
	}
	table.Render()
	return errorCount, warningCount, notifyCount
}

func colorize(message string, textColor color.Attribute) string {
	colorPrinter := color.New(textColor)
	f := colorPrinter.SprintFunc()
	return f(message)
}
