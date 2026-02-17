package dockfmt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLint_Name(t *testing.T) {
	l := &LintDockFmt{}
	assert.Equal(t, "dockfmt", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &LintDockFmt{}
	assert.NoError(t, l.Setup(t.Context()))
}

func TestLint_Cleanup(t *testing.T) {
	l := &LintDockFmt{}
	err := l.Cleanup(t.Context())
	assert.NoError(t, err)
}

func TestLint_Run(t *testing.T) {

	// Don't test shfmt but make sure switching lang works
	tests := []struct {
		name     string
		linter   LintDockFmt
		input    string
		expected string
		error    bool
	}{
		{
			name:   "dockerfile default",
			linter: LintDockFmt{},
			input: `from alpine as base
RUN echo    "hello world"
RUN ls > dirlist   2>&1
RUN apk add --no-cache \
bash \
git`,
			expected: `FROM alpine as base
RUN echo "hello world"
RUN ls >dirlist 2>&1
RUN apk add --no-cache \
    bash \
    git`,
		},
		{
			name:   "dockerfile with settings",
			linter: LintDockFmt{IndentSize: 2, TrailingNewline: true, SpaceRedirects: true},
			input: `from alpine as base
RUN echo    "hello world"
RUN ls > dirlist   2>&1
RUN apk add --no-cache \
bash \
git`,
			expected: `FROM alpine as base
RUN echo "hello world"
RUN ls > dirlist 2>&1
RUN apk add --no-cache \
  bash \
  git
`,
		},
		{
			name:   "not a dockerfile",
			linter: LintDockFmt{},
			input: `this isn't a dockerfile
its a test to see what happens when you run dockfmt on a file that isn't a dockerfile'`,
			expected: `this isn't a dockerfile
its a test to see what happens when you run dockfmt on a file that isn't a dockerfile'`,
			error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpDir := t.TempDir()
			testFilePath := filepath.Join(tmpDir, "Dockerfile")
			err := os.WriteFile(testFilePath, []byte(tt.input), 0644)
			assert.NoError(t, err)

			report := tt.linter.Run(testFilePath)

			assert.Equal(t, testFilePath, report.Target)
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
