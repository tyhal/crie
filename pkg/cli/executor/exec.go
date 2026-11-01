// Package executor allows for a CLI tool to be executed in different contexts
package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"

	"go.opentelemetry.io/contrib/propagators/envcar"
	"go.opentelemetry.io/otel"
)

// ErrRuntimeUnavailable is returned when a container runtime is configured but can't be used,
// callers should stop rather than silently falling back to another executor
var ErrRuntimeUnavailable = errors.New("container runtime unavailable")

// Par represents command-line parameters passed to a linter executable.
type Par []string

// Instance describes a single CLI execution configuration for a linter tool.
// It specifies the binary to run, positional arguments to prepend/append,
// and whether execution should change to the file's directory first.
type Instance struct {
	Bin        string `json:"bin" yaml:"bin" jsonschema:"the binary or command to use"`
	Start      Par    `json:"start,omitzero" yaml:"start,flow" jsonschema:"parameters that will be put in front of the file path"`
	End        Par    `json:"end,omitzero" yaml:"end,flow" jsonschema:"parameters that will be put behind the file path"`
	ChDir      bool   `json:"chdir,omitzero" yaml:"chdir" jsonschema:"if true the tool will change directory to where the target file is located"`
	WillWrite  bool   `json:"will_write,omitzero" yaml:"will_write" jsonschema:"whether the tool will write to the file system, if unset for container execution will mount with read-only privileges"`
	NoFileArg  bool   `json:"no_file_arg,omitzero" yaml:"no_file_arg" jsonschema:"if true the target path is not inserted into command arguments; combine with chdir to use the target as the working directory"`
	Batch      bool   `json:"batch,omitzero" yaml:"batch" jsonschema:"if true the tool is given many file paths in one run, it must print each finding on a line led by the file path and a colon (e.g. path:line:col: msg) so results can be reported per file; can't be combined with chdir or no_file_arg"`
	BatchStart Par    `json:"batch_start,omitzero" yaml:"batch_start,flow" jsonschema:"parameters put after start only when a run covers many files, e.g. an output format leading each finding with its file path, so single file runs keep the tool's usual output"`
}

// ErrBadInstance is returned when an Instance has options that can't be used together
var ErrBadInstance = errors.New("invalid exec options")

// Validate checks the Instance's options can be used together
func (i Instance) Validate() error {
	if i.Batch && (i.ChDir || i.NoFileArg) {
		return fmt.Errorf("%w: batch can't be combined with chdir or no_file_arg", ErrBadInstance)
	}
	return nil
}

// buildParams constructs the command arguments after the binary name.
// targets should already be in the format expected by the executor (OS path for host, Linux path for containers).
func (i Instance) buildParams(targets ...string) []string {
	params := append([]string{}, i.Start...)
	if len(targets) > 1 {
		params = append(params, i.BatchStart...)
	}
	if !i.NoFileArg {
		for _, target := range targets {
			if i.ChDir {
				target = filepath.Base(target)
			}
			params = append(params, target)
		}
	}
	return append(params, i.End...)
}

// apiObjectID matches container and exec ids in container API paths
var apiObjectID = regexp.MustCompile(`/[0-9a-f]{64}`)

// apiSpanName names container API client spans "METHOD path" with ids replaced,
// so it's the same for every call of the same operation
func apiSpanName(_ string, r *http.Request) string {
	return r.Method + " " + apiObjectID.ReplaceAllString(r.URL.Path, "/{id}")
}

// traceEnv returns the trace context held in ctx as KEY=value environment variables (TRACEPARENT etc.)
// so that tools which understand OpenTelemetry can parent their spans to crie's.
func traceEnv(ctx context.Context) []string {
	var env []string
	otel.GetTextMapPropagator().Inject(ctx, &envcar.Carrier{
		SetEnvFunc: func(k, v string) { env = append(env, k+"="+v) },
	})
	return env
}

// ErrNoTargets is returned when Exec is called without any target paths
var ErrNoTargets = errors.New("no targets to execute against")

// Executor abstracts how a CLI tool is executed (host, Docker, Podman, etc.).
// Implementations should prepare resources in Setup, execute the tool in Exec,
// and free resources in Cleanup.
// Exec runs the tool once against all targets, more than one target is only expected for a batch Instance.
type Executor interface {
	Setup(ctx context.Context, i Instance) error
	Exec(ctx context.Context, stdout io.Writer, stderr io.Writer, targets ...string) error
	Cleanup(ctx context.Context) error
}
