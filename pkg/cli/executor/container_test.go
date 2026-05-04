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

// testContainerExecutor runs common integration tests for container executors
func testContainerExecutor(t *testing.T, newExecutor func(image string) Executor) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Run("basic execution", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("hello"), 0o644))

		// Change to tmpDir so the container can mount it
		origDir, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(tmpDir))
		defer os.Chdir(origDir)

		e := newExecutor("alpine:latest")

		ctx := context.Background()
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
	})

	t.Run("ChDir", func(t *testing.T) {
		tmpDir := t.TempDir()
		subDir := filepath.Join(tmpDir, "subdir")
		require.NoError(t, os.MkdirAll(subDir, 0o755))
		filePath := filepath.Join(subDir, "test.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("content"), 0o644))

		origDir, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(tmpDir))
		defer os.Chdir(origDir)

		e := newExecutor("alpine:latest")

		ctx := context.Background()
		err = e.Setup(ctx, Instance{
			Bin:   "sh",
			Start: []string{"-c", "basename $(pwd)", "_"},
			ChDir: true,
		})
		require.NoError(t, err)
		defer e.Cleanup(ctx)

		var stdout bytes.Buffer
		err = e.Exec(filepath.Join("subdir", "test.txt"), &stdout, &bytes.Buffer{})
		require.NoError(t, err)

		// Should be in subdir when ChDir=true
		assert.Contains(t, strings.TrimSpace(stdout.String()), "subdir")
	})
}
