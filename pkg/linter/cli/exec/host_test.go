package exec

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillHost(t *testing.T) {
	// Only verify error path to avoid environment-dependent binary checks
	assert.Error(t, WillHost("definitely-not-a-real-binary-xyz"))
}

func TestHostExecutor_Setup(t *testing.T) {
	e := &hostExecutor{}
	assert.NoError(t, e.Setup())
}

func TestHostExecutor_Run(t *testing.T) {
	// Decide which shell to use
	bin := "sh"
	if WillHost(bin) != nil {
		if WillHost("bash") == nil {
			bin = "bash"
		} else {
			t.Skip("neither 'sh' nor 'bash' is available on host; skipping")
		}
	}

	e := &hostExecutor{}

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "hello.txt")
	assert.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))

	// Script echoes PWD and the first arg, and ensures the file exists
	script := `echo "PWD=$(pwd)"; echo "ARG=$1"; test -f "$1"`
	// Note: with 'sh -c', the first argument after the script becomes $0; add a dummy so filePath maps to $1
	front := []string{"-c", script, "_"}
	var out bytes.Buffer

	ei := Instance{
		Bin:   bin,
		Start: front,
		End:   nil,
	}

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
			out.Reset()
			ei.ChDir = tc.chdir
			err = e.Exec(ei, filePath, &out, &out)
			assert.NoError(t, err, "exec with %s should succeed", tc.name)
			stdout := out.String()
			assert.Contains(t, stdout, "PWD="+tc.expPWD)
			assert.Contains(t, stdout, "ARG="+tc.expARG)
			// Also ensure outputs have the expected lines at least once
			assert.GreaterOrEqual(t, strings.Count(stdout, "PWD="), 1)
			assert.GreaterOrEqual(t, strings.Count(stdout, "ARG="), 1)
		})
	}
}

func TestHostExecutor_Cleanup(t *testing.T) {
	e := &hostExecutor{}
	assert.NoError(t, e.Cleanup())
}
