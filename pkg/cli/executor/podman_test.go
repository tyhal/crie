package executor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWillPodman(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = WillPodman(t.Context())
	}, "WillPodman should not panic")
}

func TestPodman_SocketGet(t *testing.T) {
	if WillPodman(t.Context()) != nil {
		t.Skip()
	}

	socket, err := getPodmanMachineSocket()
	assert.NoError(t, err)
	assert.NotEmpty(t, socket)
}

func TestPodmanExecutor_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	ctx := t.Context()
	if WillPodman(ctx) != nil {
		t.Skip("skipping test in non-podman mode")
	}

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))

	// Change to tmpDir so the container can mount it
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer os.Chdir(origDir)

	e := NewPodman("alpine:latest")

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
