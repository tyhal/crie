package filelist

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func commit(t *testing.T, tree *git.Worktree) {
	t.Helper()
	_, err := tree.Commit("commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	assert.NoError(t, err)
}

func removeHelper(t *testing.T, files []string, dir string, repo *git.Repository) {
	t.Helper()

	tree, err := repo.Worktree()
	assert.NoError(t, err)

	for _, file := range files {
		testFilePath := filepath.Join(dir, file)
		err := os.Remove(testFilePath)
		assert.NoError(t, err)

		_, err = tree.Add(file)
		assert.NoError(t, err)
	}

	commit(t, tree)
}

func commitHelper(t *testing.T, files []string, dir string, repo *git.Repository, contents string) {
	t.Helper()

	tree, err := repo.Worktree()
	assert.NoError(t, err)

	for _, file := range files {
		testFilePath := filepath.Join(dir, file)
		if filepath.Dir(file) != "" {
			err := os.MkdirAll(path.Join(dir, filepath.Dir(file)), 0755)
			assert.NoError(t, err)
		}
		err = os.WriteFile(testFilePath, []byte(contents), 0644)
		assert.NoError(t, err)

		_, err = tree.Add(file)
		assert.NoError(t, err)
	}

	commit(t, tree)
}

func baseFiles(files ...string) []string {
	var filesList []string
	for _, file := range files {
		filesList = append(filesList, path.Base(file))
	}
	return filesList
}

type testRepo struct {
	initialCommit     []string
	afterTargetCommit []string
	removeCommit      []string
	changeUncommited  []string
	newUncommited     []string
	ignored           []string
}

const gitIgnoreFile = ".gitignore"

func (tr testRepo) Setup(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	repo, err := git.PlainInit(tmpDir, false)
	assert.NoError(t, err)

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"/dev/null"},
	})
	assert.NoError(t, err)

	if len(tr.ignored) > 0 {
		commitHelper(t, []string{gitIgnoreFile}, tmpDir, repo, strings.Join(tr.ignored, "\n")+"\n")
		plainHelper(t, tr.ignored, tmpDir, "ignored")
	}

	commitHelper(t, tr.initialCommit, tmpDir, repo, "initial")
	head, err := repo.Head()
	assert.NoError(t, err)
	err = repo.Storer.SetReference(plumbing.NewHashReference(
		plumbing.NewRemoteReferenceName("origin", "main"),
		head.Hash(),
	))
	assert.NoError(t, err)
	commitHelper(t, tr.afterTargetCommit, tmpDir, repo, "secondary")

	removeHelper(t, tr.removeCommit, tmpDir, repo)
	plainHelper(t, tr.changeUncommited, tmpDir, "uncommitted_file_in_repo")
	plainHelper(t, tr.newUncommited, tmpDir, "uncommitted_file_not_in_repo")

	return tmpDir
}

func TestFromGitRepo(t *testing.T) {
	aInitFile := "a_init"             // an example of an unchanged file from the initial commit
	bChangedFile := "b_changed"       // an example of a changed and commited file
	cRemovedFile := "c_removed"       // an example of a removed file
	dAddedFile := "d_added"           // an example of a newly commited file
	eChangedFile := "e_changed"       // an example of a previously committed file with uncommitted changes
	fUncommitedFile := "f_uncommited" // an example of an uncommited file
	gIgnoredFile := "g_ignored"       // an example of a file ignored by git

	hInitFile := path.Join("dir", "h_init")             // an example of an unchanged file from the initial commit in a subdirectory
	iChangedFile := path.Join("dir", "i_changed")       // an example of a changed and commited file in a subdirectory
	jUncommitedFile := path.Join("dir", "j_uncommited") // an example of an uncommited file in a subdirectory

	tr := testRepo{
		initialCommit:     []string{aInitFile, bChangedFile, cRemovedFile, hInitFile, iChangedFile},
		afterTargetCommit: []string{bChangedFile, dAddedFile, eChangedFile, iChangedFile},
		removeCommit:      []string{cRemovedFile},

		changeUncommited: []string{eChangedFile},
		newUncommited:    []string{fUncommitedFile, jUncommitedFile},
		ignored:          []string{gIgnoredFile},
	}
	tmpDir := tr.Setup(t)

	baseCommitted := []string{gitIgnoreFile, aInitFile, bChangedFile, dAddedFile, eChangedFile}
	dirCommited := []string{hInitFile, iChangedFile}

	tests := []struct {
		name      string
		diff      bool
		target    string
		dir       string
		expectErr bool
		errMsg    string
		expect    []string
	}{
		// basic error cases
		{
			name:      "target missing",
			target:    "invalid",
			expectErr: true,
			diff:      true,
			errMsg:    "reference not found",
		},
		{
			name:      "invalid target",
			target:    "a/b/c",
			expectErr: true,
			diff:      true,
			errMsg:    "invalid git target",
		},
		{
			name:      "target set without diff",
			target:    "target",
			diff:      false,
			expectErr: true,
			errMsg:    "invalid parameter combination",
		},
		// basic happy path
		{
			name:      "simple diff of uncommitted changes",
			diff:      true,
			expectErr: false,
			expect:    append(tr.changeUncommited, tr.newUncommited...),
		},
		{
			name:      "diff to target",
			diff:      true,
			target:    "origin/main",
			expectErr: false,
			expect:    tr.afterTargetCommit,
		},
		{
			name:      "all files",
			expectErr: false,
			expect:    append(baseCommitted, dirCommited...),
		},
		// subdirectory happy path
		{
			name:      "simple diff of uncommitted changes in subdirectory",
			dir:       "dir",
			diff:      true,
			expectErr: false,
			expect:    baseFiles(jUncommitedFile),
		},
		{
			name:      "diff to target in subdirectory",
			dir:       "dir",
			diff:      true,
			target:    "origin/main",
			expectErr: false,
			expect:    baseFiles(iChangedFile),
		},
		{
			name:   "all files in subdirectory",
			dir:    "dir",
			expect: baseFiles(dirCommited...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currDir := path.Join(tmpDir, tt.dir)

			repo, err := GetGitRepo(currDir)
			assert.NoError(t, err)
			changed, err := FromGitRepo(currDir, repo, tt.diff, tt.target)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.expect, changed)
		})
	}
}

func TestGetGitRepo(t *testing.T) {
	tests := []struct {
		name      string
		init      bool
		dir       string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "no git repo",
			expectErr: true,
			errMsg:    "does not exist",
		},
		{
			name:      "git repo",
			init:      true,
			expectErr: false,
		},
		{
			name:      "git repo in subdir",
			init:      true,
			dir:       "subdir",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if tt.init {
				_, err := git.PlainInit(tmpDir, false)
				assert.NoError(t, err)
			}
			_, err := GetGitRepo(path.Join(tmpDir, tt.dir))
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}
