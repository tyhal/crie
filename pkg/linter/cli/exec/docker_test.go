package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillDocker(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = WillDocker()
	}, "WillPodman should not panic")
}
