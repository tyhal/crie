package noop

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLint_Name(t *testing.T) {
	l := &LintNoop{}
	assert.Equal(t, "noop", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &LintNoop{}
	assert.NoError(t, l.Setup(t.Context()))
}

func TestLint_Cleanup(t *testing.T) {
	l := &LintNoop{}
	assert.NoError(t, l.Setup(t.Context()))
	assert.NoError(t, l.Cleanup(t.Context()))
}

func TestLint_Run(t *testing.T) {
	l := &LintNoop{}
	assert.NoError(t, l.Setup(t.Context()))
	report := l.Run("test.txt")
	assert.Equal(t, "test.txt", report.Target)
	assert.NoError(t, report.Err)
	assert.Nil(t, report.StdOut)
	assert.Nil(t, report.StdErr)
}
