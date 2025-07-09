package shfmt

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
)

func TestLint_Name(t *testing.T) {
	l := &Lint{}
	assert.Equal(t, "shfmt", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &Lint{}
	assert.NoError(t, l.WillRun())
}

func TestLint_Cleanup(t *testing.T) {
	l := &Lint{}
	var wg sync.WaitGroup
	wg.Add(1)
	l.Cleanup(&wg)
	wg.Wait()
}

func TestLint_MaxConcurrency(t *testing.T) {
	l := &Lint{}
	assert.Equal(t, math.MaxInt32, l.MaxConcurrency())
}

// TODO either mock or make dummy files
//func TestLint_Run(t *testing.T) {
//	l := &Lint{}
//	rep := make(chan linter.Report, 1)
//
//	l.Run("test.txt", rep)
//
//	report := <-rep
//	assert.Equal(t, "test.txt", report.File)
//	assert.NoError(t, report.Err)
//	assert.Nil(t, report.StdOut)
//	assert.Nil(t, report.StdErr)
//}
