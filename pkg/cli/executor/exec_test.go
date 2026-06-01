package executor

import (
	"bytes"
	"context"
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
	filePath := filepath.Join(tmpDir, "hello.txt")
	assert.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))

	// Note: with 'sh -c', the first argument after the script becomes $0; add a dummy so filePath maps to $1
	front := []string{"-c", `echo "PWD=$(pwd)"; echo "ARG=$1"; test -f "$1"`, "_"}
	var out bytes.Buffer

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	tests := []struct {
		name   string
		chdir  bool
		expPWD string
		expARG string
	}{
		{"chdir=false", false, cwd, filePath},
		{"chdir=true", true, tmpDir, filepath.Base(filePath)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			out.Reset()
			e := newExec()
			err = e.Setup(t.Context(), Instance{
				Bin:   "sh",
				Start: front,
				End:   nil,
				ChDir: tc.chdir,
			})
			defer func(e Executor, ctx context.Context) {
				err := e.Cleanup(ctx)
				require.NoError(t, err, "cleanup failed")
			}(e, t.Context())
			require.NoError(t, err)
			err = e.Exec(filePath, &out, &out)
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
