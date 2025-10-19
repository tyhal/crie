package linter

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner_listen(t *testing.T) {
	tests := []struct {
		name     string
		runner   *Runner
		reports  []Report
		wantErrs []bool
	}{

		{"success", &Runner{ShowPass: true}, []Report{{File: "test.go", Err: nil}}, []bool{false}},
		{"error", &Runner{}, []Report{{File: "test.go", Err: fmt.Errorf("fail")}}, []bool{true}},
		{"mixed", &Runner{ShowPass: true, StrictLogging: true}, []Report{{File: "ok.go", Err: nil}, {File: "bad.go", Err: fmt.Errorf("fail")}}, []bool{false, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := make(chan error, len(tt.reports))
			linterReport := make(chan Report, len(tt.reports))

			go tt.runner.listen(results, linterReport)

			for _, report := range tt.reports {
				linterReport <- report
			}
			close(linterReport)

			for i, wantErr := range tt.wantErrs {
				err := <-results
				if (err != nil) != wantErr {
					t.Errorf("result %d: wantErr=%v, got=%v", i, wantErr, err)
				}
			}
		})
	}
}

func TestRunner_LintFileList(t *testing.T) {
	tests := []struct {
		name  string
		files []string
	}{
		{"no files", []string{}},
		{"single file", []string{"test.go"}},
		{"multiple files", []string{"test1.go", "test2.go", "test3.go"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &Runner{}
			linter := &mockLinter{}

			err := runner.LintFileList(linter, tt.files)

			assert.NoError(t, err, "linting %s", tt.name)
		})
	}
}

type mockLinter struct{}

func (m *mockLinter) Name() string                         { return "mock" }
func (m *mockLinter) WillRun() error                       { return nil }
func (m *mockLinter) Cleanup(wg *sync.WaitGroup)           { wg.Done() }
func (m *mockLinter) MaxConcurrency() int                  { return 2 }
func (m *mockLinter) Run(filePath string, rep chan Report) { rep <- Report{File: filePath} }
