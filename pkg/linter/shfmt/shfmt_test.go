package shfmt

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyhal/crie/pkg/crie/linter"
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestLint_Name(t *testing.T) {
	l := &LintShfmt{}
	assert.Equal(t, "shfmt", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &LintShfmt{}
	assert.NoError(t, l.WillRun())
}

func TestLint_Cleanup(t *testing.T) {
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

	smolSh := `#!/bin/sh
set -x
echo           "hello world"
`
	expected := `#!/bin/sh
set -x
echo "hello world"
`

	tmpDir := t.TempDir()
	testFilePath := filepath.Join(tmpDir, "test.sh")
	err := os.WriteFile(testFilePath, []byte(smolSh), 0644)
	assert.NoError(t, err)

	l := &LintShfmt{
		Language: "sh",
	}
	rep := make(chan linter.Report, 1)

	l.Run(testFilePath, rep)

	report := <-rep
	assert.Equal(t, testFilePath, report.File)
	assert.NoError(t, report.Err)

	actual, err := os.ReadFile(testFilePath)
	require.Equal(t, expected, string(actual))
}
