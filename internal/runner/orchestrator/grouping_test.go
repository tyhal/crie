package orchestrator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindModuleRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte{}, 0o644))
	sub := filepath.Join(root, "internal", "pkg")
	require.NoError(t, os.MkdirAll(sub, 0o755))

	got, err := findModuleRoot(filepath.Join(sub, "foo.go"), "go.mod")
	require.NoError(t, err)
	assert.Equal(t, root, got)

	_, err = findModuleRoot(filepath.Join(t.TempDir(), "foo.go"), "go.mod")
	assert.Error(t, err)
}

func TestGroupByModule(t *testing.T) {
	base := t.TempDir()
	modA := filepath.Join(base, "a")
	modB := filepath.Join(base, "b")
	require.NoError(t, os.MkdirAll(modA, 0o755))
	require.NoError(t, os.MkdirAll(modB, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(modA, "go.mod"), []byte{}, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(modB, "go.mod"), []byte{}, 0o644))

	fileA := filepath.Join(modA, "foo.go")
	fileB := filepath.Join(modB, "bar.go")
	orphan := filepath.Join(t.TempDir(), "orphan.go")

	groups := groupByModule([]string{fileA, fileB, orphan}, "go.mod")

	assert.Len(t, groups, 2)
	assert.Equal(t, []string{fileA}, groups[modA])
	assert.Equal(t, []string{fileB}, groups[modB])
}
