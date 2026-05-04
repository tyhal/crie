package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillDocker(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = WillDocker()
	}, "WillDocker should not panic")
}

func TestDockerExecutor_Integration(t *testing.T) {
	if err := WillDocker(); err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	testContainerExecutor(t, NewDocker)
}
