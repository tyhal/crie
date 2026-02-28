package noop

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"errors"
	"testing"
	"testing/synctest"
	"time"

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

func TestLint_WithSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := WithSleep(200, 100)

		go l.Setup(t.Context())
		synctest.Wait()
		time.Sleep(100)

		go l.Run("test.txt")
		synctest.Wait()
		time.Sleep(200)

		go l.Cleanup(t.Context())
		synctest.Wait()
		time.Sleep(100)
	})
}

func TestLint_WithErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := WithErr(errors.New("run error"), errors.New("run error"))
		err := l.Setup(t.Context())
		assert.Error(t, err)
		rep := l.Run("test.txt")
		assert.Error(t, rep.Err)
	})
}

func TestLint_Ordering(t *testing.T) {
	l := &LintNoop{}
	err := l.Cleanup(t.Context())
	assert.ErrorIs(t, err, ErrMissedSetup)

	l = &LintNoop{}
	rep := l.Run("test.txt")
	assert.ErrorIs(t, rep.Err, ErrMissedSetup)
}
