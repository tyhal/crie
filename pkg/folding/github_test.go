package folding

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubFolder(t *testing.T) {
	f := NewGithub(os.Stdout)
	id, err := f.Start("a/b", "hello", true)
	assert.NoError(t, err)
	assert.Equal(t, "", id)
	err = f.Stop(id)
	assert.NoError(t, err)
}
