package executor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWillDocker(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = WillDocker()
	}, "WillDocker should not panic")
}

func TestDockerExecutor_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if WillDocker() != nil {
		t.Skip("skipping test in non-docker mode")
	}

	if err := WillDocker(); err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))

	// Change to tmpDir so the container can mount it
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer os.Chdir(origDir)

	e := NewDocker("alpine:latest")

	ctx := t.Context()
	err = e.Setup(ctx, Instance{
		Bin:   "sh",
		Start: []string{"-c", "echo PWD=$(pwd); echo ARG=$1; cat $1", "_"},
	})
	require.NoError(t, err)
	defer e.Cleanup(ctx)

	var stdout, stderr bytes.Buffer
	err = e.Exec("test.txt", &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	assert.Contains(t, output, "PWD=")
	assert.Contains(t, output, "ARG=test.txt")
	assert.Contains(t, output, "hello")
	assert.Empty(t, stderr.String())
}
