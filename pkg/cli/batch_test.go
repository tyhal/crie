package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tyhal/crie/pkg/cli/executor"
	"github.com/tyhal/crie/pkg/linter"
)

// batchExecutor is a test double that writes fixed output and returns a fixed error, recording each run's targets
type batchExecutor struct {
	out  string
	err  error
	runs [][]string
}

func (b *batchExecutor) Setup(_ context.Context, _ executor.Instance) error { return nil }
func (b *batchExecutor) Cleanup(_ context.Context) error                    { return nil }
func (b *batchExecutor) Exec(_ context.Context, stdout io.Writer, _ io.Writer, targets ...string) error {
	b.runs = append(b.runs, targets)
	_, _ = io.WriteString(stdout, b.out)
	return b.err
}

func TestLintCli_RunBatch(t *testing.T) {
	files := []string{"a.go", "dir/b.go", "c.go"}
	lintErr := linter.Result(errors.New("exit code 1"))
	execErr := errors.New("container gone")
	finding := func(file string, line, col int, msg string) linter.Finding {
		return linter.Finding{Tool: "tool", File: file, Line: line, Col: col, Message: msg}
	}

	tests := []struct {
		name        string
		out         string
		err         error
		expRerun    map[string][]linter.Finding
		expFindings map[string][]linter.Finding
		expFail     bool
	}{
		{
			name: "a passing batch passes every file",
		},
		{
			name:        "a passing batch keeps the findings it printed",
			out:         "dir/b.go:4:2: just a warning\n",
			expFindings: map[string][]linter.Finding{"dir/b.go": {finding("dir/b.go", 4, 2, "just a warning")}},
		},
		{
			// dir/b.go may only have a warning, so it's rerun for its own result like c.go, each keeping its findings
			name: "a failing batch reruns the files it names",
			out:  "  /src/dir/b.go:1:1: warning\n  snippet\n[c.go:2:1] error\n",
			err:  lintErr,
			expRerun: map[string][]linter.Finding{
				"dir/b.go": {finding("dir/b.go", 1, 1, "warning")},
				"c.go":     {finding("c.go", 2, 1, "error")},
			},
		},
		{
			name:     "a failing batch naming no file reruns them all",
			out:      "panic: something broke\n",
			err:      lintErr,
			expRerun: map[string][]linter.Finding{"a.go": nil, "dir/b.go": nil, "c.go": nil},
		},
		{
			name:        "a batch that couldn't run fails every file",
			out:         "a.go:3:1: bad\n",
			err:         execErr,
			expFindings: map[string][]linter.Finding{"a.go": {finding("a.go", 3, 1, "bad")}},
			expFail:     true,
		},
	}
	t.Run("one file is a normal run", func(t *testing.T) {
		xctr := &batchExecutor{out: "a.go:1:1: bad\n", err: lintErr}
		l := &LintCli{Exec: executor.Instance{Bin: "tool", Batch: true}, executor: xctr}
		reports, rerun := l.RunBatch(t.Context(), files[:1])
		assert.Empty(t, rerun)
		require.Len(t, reports, 1)
		assert.ErrorIs(t, reports[0].Err, lintErr)
		assert.Equal(t, "a.go:1:1: bad\n", fmt.Sprint(reports[0].StdOut), "its own output is its report")
		assert.Equal(t, []linter.Finding{finding("a.go", 1, 1, "bad")}, reports[0].Findings, "a normal run's path led output is read as findings too")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xctr := &batchExecutor{out: tt.out, err: tt.err}
			l := &LintCli{Exec: executor.Instance{Bin: "tool", Batch: true}, executor: xctr}
			require.True(t, l.Batched())

			reports, rerun := l.RunBatch(t.Context(), files)
			assert.Equal(t, [][]string{files}, xctr.runs, "every file should be passed to a single exec")
			assert.Equal(t, tt.expRerun, rerun)
			require.Len(t, reports, len(files)-len(rerun), "every file is either reported or rerun")
			for _, rep := range reports {
				assert.NotNil(t, rep.StdOut)
				assert.NotNil(t, rep.StdErr)
				assert.Equal(t, tt.expFindings[rep.Target[len("tool "):]], rep.Findings, rep.Target)
				if tt.expFail {
					assert.ErrorIs(t, rep.Err, execErr, rep.Target)
				} else {
					assert.NoError(t, rep.Err, rep.Target)
				}
			}
		})
	}
}

func TestParseFinding(t *testing.T) {
	tests := []struct {
		rest string
		exp  linter.Finding
	}{
		{"3:1: bad", linter.Finding{Line: 3, Col: 1, Message: "bad"}},
		{"12: no column", linter.Finding{Line: 12, Message: "no column"}},
		{"24 DL3003 warning: Use WORKDIR", linter.Finding{Line: 24, Message: "DL3003 warning: Use WORKDIR"}},
		{"2:1] Message name must be UpperCamelCase", linter.Finding{Line: 2, Col: 1, Message: "Message name must be UpperCamelCase"}},
		{"31:9:proselint.Very:Remove 'very'.", linter.Finding{Line: 31, Col: 9, Message: "proselint.Very:Remove 'very'."}},
		{"13:1: [error] duplication of key", linter.Finding{Line: 13, Col: 1, Message: "[error] duplication of key"}},
		{" written", linter.Finding{Message: "written"}},
		{"7:x: odd", linter.Finding{Line: 7, Message: "x: odd"}},
	}
	for _, tt := range tests {
		tt.exp.Tool, tt.exp.File = "tool", "f.go"
		assert.Equal(t, tt.exp, parseFinding("tool", "f.go", tt.rest), tt.rest)
	}
}

func TestLineOwner(t *testing.T) {
	files := []string{"a.go", "./b.go", "dir/a.go"}
	tests := []struct {
		line    string
		exp     int
		expRest string
	}{
		{"a.go:1:1: x", 0, "1:1: x"},
		{"./a.go:1:1: x", 0, "1:1: x"},
		{"b.go:1:1: x", 1, "1:1: x"},
		{"./b.go:1: x", 1, "1: x"},
		{"dir/a.go:2: x", 2, "2: x"},
		{"./dir/a.go:2: x", 2, "2: x"},
		{"  /src/dir/a.go:1:1: x\n", 2, "1:1: x"},
		{"\t/src/a.go:1:1: x", 0, "1:1: x"},
		{"[dir/a.go:2:1] x", 2, "2:1] x"},
		{"[a.go:2:1] x", 0, "2:1] x"},
		{"/src/other/b.go:1: x", 1, "1: x"},
		{"a.go", -1, ""},
		{"a.go 1 x", -1, ""},
		{"a.gox:1: x", -1, ""},
		{"xa.go:1: x", -1, ""},
		{"note: a.go:1: x", -1, ""},
		{":a.go", -1, ""},
		{"", -1, ""},
	}
	for _, tt := range tests {
		owner, rest := lineOwner(files, []byte(tt.line))
		assert.Equal(t, tt.exp, owner, tt.line)
		if tt.exp >= 0 {
			assert.Equal(t, tt.expRest, rest, tt.line)
		}
	}
}
