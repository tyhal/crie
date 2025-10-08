package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillDocker(t *testing.T) {
	assert.NoError(t, WillDocker())
}
