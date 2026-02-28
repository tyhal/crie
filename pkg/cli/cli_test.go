package cli

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"bytes"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/cli/executor"
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
	l := &LintCli{Exec: executor.Instance{Bin: "test"}}
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
		l := &LintCli{executor: executor.NewNoop()}
		assert.NotPanics(t, func() {
			err := l.Cleanup(t.Context())
			assert.NoError(t, err)
		})
	}
	{
		l := &LintCli{executor: nil}
		assert.NotPanics(t, func() {
			err := l.Cleanup(t.Context())
			assert.NoError(t, err)
		})
	}
}

func TestLint_Run(t *testing.T) {
	l := &LintCli{executor: executor.NewNoop()} // TODO test with no executor setup

	report := l.Run("test.txt")
	assert.Equal(t, "test.txt", report.Target)
	assert.NoError(t, report.Err)
	assert.Equal(t, "stdout", report.StdOut.(*bytes.Buffer).String())
	assert.Equal(t, "stderr", report.StdErr.(*bytes.Buffer).String())
}
