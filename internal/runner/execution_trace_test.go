//go:build trace
// +build trace

package runner

import (
	"context"
	"os"
	"runtime/trace"
	"testing"
)

func TestRunConfiguration_trace_runLinters(t *testing.T) {
	defer disableLogging()()

	f, err := os.Create("artificial_trace.out")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	opts := Options{
		StrictLogging: true,
	}

	test := struct {
		config *RunConfiguration
		files  []string
	}{
		config: &RunConfiguration{
			Options:      opts,
			NamedMatches: genLangs(10),
		},
		files: genFilenames(100),
	}

	ctx := t.Context()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	_ = test.config.runLinters(ctx, LintTypeChk, test.files)
}
