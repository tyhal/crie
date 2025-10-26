package filelist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func plainHelper(t *testing.T, files []string, dir string, contents string) {
	t.Helper()
	for _, file := range files {
		testFilePath := filepath.Join(dir, file)
		err := os.WriteFile(testFilePath, []byte(contents), 0644)
		assert.NoError(t, err)
	}
}

func TestFromDir(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	assert.NoError(t, os.Mkdir(subDir, 0755))

	dirFile := filepath.Join("sub", "dir.txt")
	expFiles := []string{"a.txt", "b.txt", "c.txt", dirFile}
	plainHelper(t, expFiles, tmpDir, "test")

	files, err := FromDir(tmpDir)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expFiles, files)

	files, err = FromDir(subDir)
	assert.NoError(t, err)
	assert.ElementsMatch(t, baseFiles(dirFile), files)
}
