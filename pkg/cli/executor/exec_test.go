package executor

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testHelperExecutor is a helper function for testing executors
func testHelperExecutor(t *testing.T, newExec func() Executor) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	subdir := "sub"
	err := os.Mkdir(filepath.Join(tmpDir, "sub"), 0o755)
	require.NoError(t, err)

	filePath := filepath.Join(subdir, "hello.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))
	otherPath := filepath.Join(subdir, "world.txt")
	require.NoError(t, os.WriteFile(otherPath, []byte("world"), 0o644))

	// Note: with 'sh -c', the first argument after the script becomes $0; add a dummy so the targets map to $1...
	front := []string{"-c", `echo "PWD=$(pwd)"; echo "ARG=$*"; for f; do test -f "$f" || exit 1; done`, "_"}
	var out bytes.Buffer

	tests := []struct {
		name    string
		chdir   bool
		batch   bool
		targets []string
		expPWD  string
		expARG  string
	}{
		{"chdir=false", false, false, []string{filePath}, tmpDir, filePath},
		{"chdir=true", true, false, []string{filePath}, filepath.Join(tmpDir, "sub"), filepath.Base(filePath)},
		{"batch=true", false, true, []string{filePath, otherPath}, tmpDir, filePath + " " + otherPath},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			out.Reset()
			e := newExec()
			err := e.Setup(t.Context(), Instance{
				Bin:   "sh",
				Start: front,
				End:   nil,
				ChDir: tc.chdir,
				Batch: tc.batch,
			})
			defer func(e Executor, ctx context.Context) {
				err := e.Cleanup(ctx)
				require.NoError(t, err, "cleanup failed")
			}(e, t.Context())
			require.NoError(t, err)
			err = e.Exec(t.Context(), &out, &out, tc.targets...)
			if assert.NoError(t, err) {
				stdout := out.String()
				assert.Contains(t, stdout, "PWD="+tc.expPWD)
				assert.Contains(t, stdout, "ARG="+tc.expARG)
				// Also, ensure outputs have the expected lines at least once
				assert.GreaterOrEqual(t, strings.Count(stdout, "PWD="), 1)
				assert.GreaterOrEqual(t, strings.Count(stdout, "ARG="), 1)
			} else {
				t.Logf("stdout: %s", out.String())
			}
		})
	}
}

func TestInstance_buildParams(t *testing.T) {
	tests := []struct {
		name    string
		inst    Instance
		targets []string
		exp     []string
	}{
		{"single", Instance{Start: Par{"-a"}, End: Par{"-z"}}, []string{"d/f"}, []string{"-a", "d/f", "-z"}},
		{"many", Instance{Start: Par{"-a"}, End: Par{"-z"}, Batch: true}, []string{"d/f", "g"}, []string{"-a", "d/f", "g", "-z"}},
		{"batch start for many", Instance{Start: Par{"-a"}, BatchStart: Par{"-f", "gcc"}, Batch: true}, []string{"d/f", "g"}, []string{"-a", "-f", "gcc", "d/f", "g"}},
		{"no batch start for one", Instance{Start: Par{"-a"}, BatchStart: Par{"-f", "gcc"}, Batch: true}, []string{"d/f"}, []string{"-a", "d/f"}},
		{"chdir", Instance{ChDir: true}, []string{"d/f"}, []string{"f"}},
		{"no file arg", Instance{Start: Par{"-a"}, NoFileArg: true}, []string{"d"}, []string{"-a"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp, tc.inst.buildParams(tc.targets...))
		})
	}
}

func TestInstance_Validate(t *testing.T) {
	assert.NoError(t, Instance{Batch: true}.Validate())
	assert.NoError(t, Instance{ChDir: true, NoFileArg: true}.Validate())
	assert.ErrorIs(t, Instance{Batch: true, ChDir: true}.Validate(), ErrBadInstance)
	assert.ErrorIs(t, Instance{Batch: true, NoFileArg: true}.Validate(), ErrBadInstance)
}

func TestExec_NoTargets(t *testing.T) {
	e := NewHost()
	require.NoError(t, e.Setup(t.Context(), Instance{Bin: "true"}))
	assert.ErrorIs(t, e.Exec(t.Context(), io.Discard, io.Discard), ErrNoTargets)
}
