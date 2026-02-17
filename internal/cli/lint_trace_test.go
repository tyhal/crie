//go:build trace
// +build trace

package cli

import (
	"context"
	"io"
	"os"
	"runtime/trace"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/tyhal/crie/internal/config/language"
	"github.com/tyhal/crie/internal/runner"
)

// disableLogging changes the logger output and returns a restore function.
// The caller is expected to defer the returned function.
func disableLogging() func() {
	originalOutput := logrus.StandardLogger().Out
	logrus.SetOutput(io.Discard)
	return func() {
		logrus.SetOutput(originalOutput)
	}
}

func TestRunConfiguration_trace_Run(t *testing.T) {
	defer disableLogging()()

	// Change directory up two levels
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	err = os.Chdir("../..")
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	f, err := os.Create("real_trace.out")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	langs, err := language.DefaultLanguageConfig().ToRunFormat()
	if err != nil {
		panic(err)
	}
	config := &runner.RunConfiguration{
		Options: runner.Options{
			StrictLogging: true,
			Continue:      true,
			Passes:        true,
		},
		Languages: langs,
	}

	ctx := t.Context()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	_ = config.Run(ctx, runner.LintTypeChk)
}
