package cli

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/tyhal/crie/pkg/linter"
)

// Batched reports if the tool is configured to lint many files per run, see executor.Instance Batch
func (e *LintCli) Batched() bool {
	return e.Exec.Batch
}

// RunBatch lints all filePaths in one run of the tool.
// A tool's severity rules only reach us through its exit code, which a batch shares between all its files
// (e.g. a file with only warnings would fail alongside one with an error), so when a batch fails
// each file its output names is returned to rerun on its own for its own result, or every file if it names none.
// The batch's machine readable findings are kept for each file, as a rerun's output is for people to read.
func (e *LintCli) RunBatch(ctx context.Context, filePaths []string) (reports []linter.Report, rerun map[string][]linter.Finding) {
	// one file is just a normal run, its exit code is its own
	if len(filePaths) == 1 {
		return []linter.Report{e.Run(ctx, filePaths[0])}, nil
	}

	var out bytes.Buffer
	// stdout and stderr share a buffer so lines keep the order the tool wrote them in
	err := e.executor.Exec(ctx, &out, &out, filePaths...)

	_, lintFailed := errors.AsType[*linter.FailedResultError](err)
	found := e.findings(filePaths, out.Bytes())

	for i, filePath := range filePaths {
		rep := linter.Report{Target: e.Name() + " " + filePath, File: filePath, StdOut: &bytes.Buffer{}, StdErr: &bytes.Buffer{}, Findings: found[i]}
		switch {
		case lintFailed && (len(found) == 0 || len(found[i]) > 0):
			if rerun == nil {
				rerun = map[string][]linter.Finding{}
			}
			rerun[filePath] = found[i]
			continue
		case err != nil && !lintFailed:
			// the tool couldn't run at all (e.g. the exec failed), which fails every file
			rep.Err, rep.StdOut = err, bytes.NewReader(out.Bytes())
		}
		reports = append(reports, rep)
	}
	return reports, rerun
}

// findings returns each line of out led by one of filePaths (see lineOwner) as a finding for that file
func (e *LintCli) findings(filePaths []string, out ...[]byte) map[int][]linter.Finding {
	found := map[int][]linter.Finding{}
	for _, o := range out {
		for line := range bytes.Lines(o) {
			if i, rest := lineOwner(filePaths, line); i >= 0 {
				found[i] = append(found[i], parseFinding(e.Name(), filePaths[i], rest))
			}
		}
	}
	return found
}

// parseFinding reads the "line[:col]" and message following a finding's leading path,
// e.g. "3:1: msg", "24 DL3003 msg" (hadolint) or "2:1] msg" (protolint). A missing line is left as 0.
func parseFinding(tool, filePath, rest string) linter.Finding {
	f := linter.Finding{Tool: tool, File: filePath}
	num := func() int {
		end := 0
		for end < len(rest) && rest[end] >= '0' && rest[end] <= '9' {
			end++
		}
		n, _ := strconv.Atoi(rest[:end])
		rest = rest[end:]
		return n
	}
	f.Line = num()
	if f.Line > 0 && len(rest) > 1 && rest[0] == ':' && rest[1] >= '0' && rest[1] <= '9' {
		rest = rest[1:]
		f.Col = num()
	}
	f.Message = strings.TrimSpace(strings.TrimLeft(rest, ":] "))
	return f
}

// lineOwner returns the index of the file path that leads line, or -1 if there isn't one, and the rest of the line after it.
// The leading path is the token before the first ':', after any indentation and an opening '[' (e.g. "[path:1:2] msg").
// It belongs to a file when it's the file's path, with or without a leading "./", or ends with "/" and the path,
// as tools in a container may print absolute paths. An exact match wins over the longest suffix match.
func lineOwner(filePaths []string, line []byte) (int, string) {
	token := strings.TrimPrefix(strings.TrimLeft(string(line), " \t"), "[")
	end := strings.IndexByte(token, ':')
	if end <= 0 {
		return -1, ""
	}
	rest := strings.TrimRight(token[end+1:], "\r\n")
	token = strings.TrimPrefix(token[:end], "./")

	owner, longest := -1, 0
	for i, filePath := range filePaths {
		filePath = strings.TrimPrefix(filePath, "./")
		if token == filePath {
			return i, rest
		}
		if len(filePath) > longest && strings.HasSuffix(token, "/"+filePath) {
			owner, longest = i, len(filePath)
		}
	}
	return owner, rest
}
