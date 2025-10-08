package cli

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli/exec"

	"sync"
	"testing"
)

func TestLintCli_isContainer(t *testing.T) {
	{
		l := &LintCli{Img: "docker.io/tyhal/crie-dep-apk:latest"}
		assert.True(t, l.isContainer())
	}
	{
		l := &LintCli{}
		assert.False(t, l.isContainer())
	}
}

func TestLint_Name(t *testing.T) {
	l := &LintCli{Exec: exec.ExecInstance{Bin: "test"}}
	assert.Equal(t, "test", l.Name())
}

func TestLint_Cleanup(t *testing.T) {
	l := &LintCli{executor: &exec.NoopExecutor{}}
	var wg sync.WaitGroup
	wg.Add(1)
	l.Cleanup(&wg)
	wg.Wait()
}

func TestLint_Run(t *testing.T) {
	l := &LintCli{executor: &exec.NoopExecutor{}} // TODO test with no executor setup
	rep := make(chan linter.Report, 1)

	l.Run("test.txt", rep)

	report := <-rep
	assert.Equal(t, "test.txt", report.File)
	assert.NoError(t, report.Err)
	assert.Equal(t, "stdout", report.StdOut.(*bytes.Buffer).String())
	assert.Equal(t, "stderr", report.StdErr.(*bytes.Buffer).String())
}
