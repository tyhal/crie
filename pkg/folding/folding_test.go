package folding

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getFoldWithEnv(t *testing.T, env string) Folder {
	t.Helper()
	folder := newFrom(os.Stdout, func(s string) bool {
		return s == env
	})
	return folder
}

func TestSetEnv(t *testing.T) {
	t.Setenv("CRIE_RANDOM_ENV", "RANDOM_VALUE")
	is := isEnvSet("CRIE_RANDOM_ENV")
	assert.True(t, is)
	isnt := isEnvSet("CRIE_UNSET_ENV")
	assert.False(t, isnt)
}

func TestNewFolder(t *testing.T) {
	var b bytes.Buffer
	folder := NewW(&b)
	folder.Start("a/b", "hello", true)
	assert.Greater(t, b.Len(), 0)

	folder = New()
	assert.NotNil(t, folder)

	folder = getFoldWithEnv(t, "FALLBACK")
	assert.IsType(t, &plainFolder{}, folder)

	folder = getFoldWithEnv(t, "GITLAB_CI")
	assert.IsType(t, &gitlabFolder{}, folder)

	folder = getFoldWithEnv(t, "GITHUB_ACTIONS")
	assert.IsType(t, &githubFolder{}, folder)
}
