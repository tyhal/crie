package linter

import (
	"fmt"
	"io"
	"path"
	"strings"
)

// GitHubReporter writes each report's findings as GitHub Actions workflow commands, which annotate the lines
// they're about on a pull request, then passes the report on to the Reporter it wraps
type GitHubReporter struct {
	Reporter
	Out io.Writer
	// Dir is prefixed to finding paths to make them relative to the repository root, when crie runs in a subdirectory
	Dir string
}

// NewGitHubReporter annotates findings on out before reporting them with next
func NewGitHubReporter(next Reporter, out io.Writer, dir string) Reporter {
	return &GitHubReporter{Reporter: next, Out: out, Dir: dir}
}

// Log writes a workflow command for each finding with a line, as an error if the report failed or else a warning.
// A failed file without any is still annotated, on the file as a whole.
func (r *GitHubReporter) Log(rep *Report) error {
	level := "warning"
	if rep.Err != nil {
		level = "error"
	}
	annotated := false
	for _, f := range rep.Findings {
		// without a line it's more likely status output (e.g. "path: written") than a problem
		if f.Line == 0 {
			continue
		}
		props := []string{"file=" + escapeProperty(path.Join(r.Dir, f.File)), fmt.Sprintf("line=%d", f.Line)}
		if f.Col > 0 {
			props = append(props, fmt.Sprintf("col=%d", f.Col))
		}
		props = append(props, "title="+escapeProperty(f.Tool))
		_, _ = fmt.Fprintf(r.Out, "::%s %s::%s\n", level, strings.Join(props, ","), escapeData(f.Message))
		annotated = true
	}
	if rep.Err != nil && !annotated && rep.File != "" {
		_, _ = fmt.Fprintf(r.Out, "::error file=%s::%s\n", escapeProperty(path.Join(r.Dir, rep.File)), escapeData(rep.Target+" failed, see the logs"))
	}
	return r.Reporter.Log(rep)
}

// escapeData escapes a workflow command's message, as the actions toolkit does
func escapeData(s string) string {
	return strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A").Replace(s)
}

// escapeProperty escapes a workflow command's property value, as the actions toolkit does
func escapeProperty(s string) string {
	return strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A", ":", "%3A", ",", "%2C").Replace(s)
}
