package cli

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"bytes"

	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter"
	"github.com/tyhal/crie/pkg/linter/cli/exec"
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
	l := &LintCli{Exec: exec.Instance{Bin: "test"}}
	assert.Equal(t, "test", l.Name())
}

func TestLint_imgTagged(t *testing.T) {
	{
		l := &LintCli{Img: "a", TagCrieVersion: true}
		assert.Equal(t, "a:latest", l.imgTagged())
	}
	{
		l := &LintCli{Img: "a", TagCrieVersion: false}
		assert.Equal(t, "a", l.imgTagged())
	}
}

func TestLint_Cleanup(t *testing.T) {
	{
		l := &LintCli{executor: &exec.noopExecutor{}}
		var wg sync.WaitGroup
		wg.Add(1)
		assert.NotPanics(t, func() { l.Cleanup(&wg) })
		wg.Wait()
	}
	{
		l := &LintCli{executor: nil}
		var wg sync.WaitGroup
		wg.Add(1)
		assert.NotPanics(t, func() { l.Cleanup(&wg) })
		wg.Wait()
	}
}

func TestLint_Run(t *testing.T) {
	l := &LintCli{executor: &exec.noopExecutor{}} // TODO test with no executor setup
	rep := make(chan linter.Report, 1)

	l.Run("test.txt", rep)

	report := <-rep
	assert.Equal(t, "test.txt", report.Target)
	assert.NoError(t, report.Err)
	assert.Equal(t, "stdout", report.StdOut.(*bytes.Buffer).String())
	assert.Equal(t, "stderr", report.StdErr.(*bytes.Buffer).String())
}
