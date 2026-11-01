package executor

import (
	"os"
	"path/filepath"
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

	sock, err := getPodmanMachineSocket()
	assert.NoError(t, err)
	assert.True(t, sock.configured)
	if _, found := os.LookupEnv("CONTAINER_HOST"); !found {
		assert.NotEmpty(t, sock.uri)
	}
}

func TestPodman_ConfiguredButUnreachable(t *testing.T) {
	t.Setenv("CONTAINER_HOST", "unix://"+filepath.Join(t.TempDir(), "missing.sock"))
	_, err := probePodman(t.Context())
	assert.ErrorIs(t, err, ErrRuntimeUnavailable)
}

func TestPodmanExecutor_Integration(t *testing.T) {
	if err := WillPodman(t.Context()); err != nil {
		t.Skipf("Podman not available: %v", err)
	}

	testHelperExecutor(t, func() Executor {
		return NewPodman("alpine:latest")
	})
}
