package cli

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/crie/linter"
	"sync"
	"testing"
)

func TestLint_Name(t *testing.T) {
	l := &LintCli{Bin: "test"}
	assert.Equal(t, "test", l.Name())
}

func TestLint_Cleanup(t *testing.T) {
	l := &LintCli{executor: &noopExecutor{}}
	var wg sync.WaitGroup
	wg.Add(1)
	l.Cleanup(&wg)
	wg.Wait()
}

func TestLint_Run(t *testing.T) {
	l := &LintCli{executor: &noopExecutor{}} // TODO test with no executor setup
	rep := make(chan linter.Report, 1)

	l.Run("test.txt", rep)

	report := <-rep
	assert.Equal(t, "test.txt", report.File)
	assert.NoError(t, report.Err)
	assert.Equal(t, "stdout", report.StdOut.(*bytes.Buffer).String())
	assert.Equal(t, "stderr", report.StdErr.(*bytes.Buffer).String())

}
