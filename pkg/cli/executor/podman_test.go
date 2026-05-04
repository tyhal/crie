package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillPodman(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = WillPodman(t.Context())
	}, "WillPodman should not panic")
}

func TestPodman_SocketGet(t *testing.T) {
	if WillPodman(t.Context()) != nil {
		t.Skip()
	}

	socket, err := getPodmanMachineSocket()
	assert.NoError(t, err)
	assert.NotEmpty(t, socket)
}

func TestPodmanExecutor_Integration(t *testing.T) {
	ctx := t.Context()
	if err := WillPodman(ctx); err != nil {
		t.Skipf("Podman not available: %v", err)
	}

	testContainerExecutor(t, NewPodman)
}
