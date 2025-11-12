package noop

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter"
)

func TestLint_Name(t *testing.T) {
	l := &LintNoop{}
	assert.Equal(t, "noop", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &LintNoop{}
	assert.NoError(t, l.WillRun())
}

func TestLint_Cleanup(_ *testing.T) {
	l := &LintNoop{}
	var wg sync.WaitGroup
	wg.Add(1)
	l.Cleanup(&wg)
	wg.Wait()
}

func TestLint_MaxConcurrency(t *testing.T) {
	l := &LintNoop{}
	assert.Equal(t, math.MaxInt32, l.MaxConcurrency())
}

func TestLint_Run(t *testing.T) {
	l := &LintNoop{}
	rep := make(chan linter.Report, 1)

	l.Run("test.txt", rep)

	report := <-rep
	assert.Equal(t, "test.txt", report.File)
	assert.NoError(t, report.Err)
	assert.Nil(t, report.StdOut)
	assert.Nil(t, report.StdErr)
}
