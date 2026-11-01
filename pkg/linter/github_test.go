package linter

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// logRecorder is a Reporter that records the reports passed on to it
type logRecorder struct{ reports []*Report }

func (r *logRecorder) Log(rep *Report) error {
	r.reports = append(r.reports, rep)
	return rep.Err
}

func TestGitHubReporter(t *testing.T) {
	failed := errors.New("exit code 1")
	tests := []struct {
		name string
		dir  string
		rep  Report
		exp  string
	}{
		{
			name: "a failed report's findings are errors",
			rep: Report{Target: "revive a.go", Err: failed, Findings: []Finding{
				{Tool: "revive", File: "a.go", Line: 90, Col: 6, Message: "don't use underscores in Go names"},
				{Tool: "revive", File: "a.go", Line: 91, Message: "no column"},
			}},
			exp: "::error file=a.go,line=90,col=6,title=revive::don't use underscores in Go names\n" +
				"::error file=a.go,line=91,title=revive::no column\n",
		},
		{
			name: "a passing report's findings are warnings",
			rep:  Report{Target: "yamllint a.yml", Findings: []Finding{{Tool: "yamllint", File: "a.yml", Line: 2, Col: 1, Message: "truthy value"}}},
			exp:  "::warning file=a.yml,line=2,col=1,title=yamllint::truthy value\n",
		},
		{
			name: "a failed file without findings is annotated as a whole",
			dir:  "sub",
			rep:  Report{Target: "shellcheck script/test", File: "script/test", Err: failed},
			exp:  "::error file=sub/script/test::shellcheck script/test failed, see the logs\n",
		},
		{
			name: "a failed file with findings isn't annotated again as a whole",
			rep:  Report{Target: "revive a.go", File: "a.go", Err: failed, Findings: []Finding{{Tool: "revive", File: "a.go", Line: 3, Message: "bad"}}},
			exp:  "::error file=a.go,line=3,title=revive::bad\n",
		},
		{
			name: "a passing file without findings isn't annotated",
			rep:  Report{Target: "revive a.go", File: "a.go"},
			exp:  "",
		},
		{
			name: "findings without a line aren't annotated",
			rep:  Report{Target: "remark a.md", Findings: []Finding{{Tool: "remark", File: "a.md", Message: "written"}}},
			exp:  "",
		},
		{
			name: "messages and properties are escaped",
			rep:  Report{Target: "x a,b.go", Err: failed, Findings: []Finding{{Tool: "x:y", File: "a,b.go", Line: 1, Message: "50% done\nnext line"}}},
			exp:  "::error file=a%2Cb.go,line=1,title=x%3Ay::50%25 done%0Anext line\n",
		},
		{
			name: "paths are made relative to the repository root",
			dir:  "sub/dir",
			rep:  Report{Target: "revive a.go", Err: failed, Findings: []Finding{{Tool: "revive", File: "a.go", Line: 3, Message: "bad"}}},
			exp:  "::error file=sub/dir/a.go,line=3,title=revive::bad\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			next := &logRecorder{}
			r := NewGitHubReporter(next, &out, tt.dir)

			err := r.Log(&tt.rep)
			assert.Equal(t, tt.exp, out.String())
			assert.Equal(t, []*Report{&tt.rep}, next.reports, "the report should still be passed on to the wrapped reporter")
			assert.Equal(t, tt.rep.Err, err)
		})
	}
}
