package runner

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func commitHelper(t *testing.T, files []string, dir string, repo *git.Repository, contents string) {

	tree, err := repo.Worktree()
	assert.NoError(t, err)

	for _, file := range files {
		testFilePath := filepath.Join(dir, file)
		err = os.WriteFile(testFilePath, []byte(contents), 0644)
		assert.NoError(t, err)

		_, err = tree.Add(file)
		assert.NoError(t, err)
	}

	_, err = tree.Commit("commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
}

func TestGit_IsRepo(t *testing.T) {
	tmpDir := t.TempDir()
	crieConfig := RunConfiguration{}
	assert.False(t, crieConfig.IsRepo(tmpDir))

	_, err := git.PlainInit(tmpDir, false)
	assert.NoError(t, err)

	assert.True(t, crieConfig.IsRepo(tmpDir))
}

func TestGit_fileListRepoAll(t *testing.T) {
	tmpDir := t.TempDir()
	crieConfig := RunConfiguration{}

	files := []string{"README.md", "another.md"}

	repo, err := git.PlainInit(tmpDir, false)
	assert.NoError(t, err)

	commitHelper(t, files, tmpDir, repo, "test")

	changed, err := crieConfig.fileListRepoAll(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, files, changed)
}

// TODO table test
func TestGit_fileListRepoChanged(t *testing.T) {
	tmpDir := t.TempDir()

	repo, err := git.PlainInit(tmpDir, false)
	assert.NoError(t, err)

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"/dev/null"},
	})
	assert.NoError(t, err)

	initialFiles := []string{"README.md", "another.md"}
	commitHelper(t, initialFiles, tmpDir, repo, "test")

	head, err := repo.Head()
	assert.NoError(t, err)
	err = repo.Storer.SetReference(plumbing.NewHashReference(
		plumbing.NewRemoteReferenceName("origin", "main"),
		head.Hash(),
	))
	assert.NoError(t, err)

	changeFiles := []string{"README.md", "changed.md"}
	commitHelper(t, changeFiles, tmpDir, repo, "changed")

	invalidConfig := RunConfiguration{
		Options: Options{
			GitTarget: "abc",
		},
	}
	_, err = invalidConfig.fileListRepoChanged(tmpDir)
	assert.Error(t, err)

	validConfig := RunConfiguration{
		Options: Options{
			GitTarget: "origin/main",
		},
	}
	changed, err := validConfig.fileListRepoChanged(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, changeFiles, changed)
}
