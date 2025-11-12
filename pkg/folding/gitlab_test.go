package folding

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitlabFolder(t *testing.T) {
	f := NewGitlab(os.Stdout)
	id, err := f.Start("a/b", "hello", true)
	assert.NoError(t, err)
	assert.NotEqual(t, "", id)
	err = f.Stop(id)
	assert.NoError(t, err)
}
