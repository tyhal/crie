package shfmt

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyhal/crie/pkg/linter"
)

func TestLint_Name(t *testing.T) {
	l := &LintShfmt{}
	assert.Equal(t, "shfmt", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &LintShfmt{}
	assert.NoError(t, l.WillRun())
}

func TestLint_Cleanup(_ *testing.T) {
	l := &LintShfmt{}
	var wg sync.WaitGroup
	wg.Add(1)
	l.Cleanup(&wg)
	wg.Wait()
}

func TestLint_MaxConcurrency(t *testing.T) {
	l := &LintShfmt{}
	assert.Equal(t, math.MaxInt32, l.MaxConcurrency())
}

func TestLint_Run(t *testing.T) {

	// Don't test shfmt but make sure switching lang works
	tests := []struct {
		name     string
		lang     string
		input    string
		expected string
		error    bool
	}{
		{
			name: "basic sh",
			lang: "sh",
			input: `#!/bin/sh
set -x
echo           "hello world"
`,
			expected: `#!/bin/sh
set -x
echo "hello world"
`,
		},
		{
			name: "basic bash",
			lang: "bash",
			input: `#!/bin/bash
set -x
echo           "hello world"
`,
			expected: `#!/bin/bash
set -x
echo "hello world"
`,
		},
		{
			name:  "unknown language",
			lang:  "unknown",
			error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpDir := t.TempDir()
			testFilePath := filepath.Join(tmpDir, "test.sh")
			err := os.WriteFile(testFilePath, []byte(tt.input), 0644)
			assert.NoError(t, err)

			l := &LintShfmt{
				Language: tt.lang,
			}
			rep := make(chan linter.Report, 1)

			l.Run(testFilePath, rep)

			report := <-rep

			assert.Equal(t, testFilePath, report.File)
			if tt.error {
				assert.Error(t, report.Err)
			} else {
				assert.NoError(t, report.Err)
			}

			actual, err := os.ReadFile(testFilePath)
			require.Equal(t, tt.expected, string(actual))

		})
	}
}
