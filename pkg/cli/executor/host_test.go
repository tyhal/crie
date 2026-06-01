package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillHost(t *testing.T) {
	// Only verify error path to avoid environment-dependent binary checks
	assert.Error(t, WillHost("definitely-not-a-real-binary-xyz"))
}

func TestHostExecutor_Setup(t *testing.T) {
	e := &hostExecutor{}
	assert.NoError(t, e.Setup(t.Context(), Instance{}))
}

func TestHostExecutor_Run(t *testing.T) {
	if WillHost("sh") != nil {
		t.Skip("neither 'sh' not available on host; skipping")
	}
	testHelperExecutor(t, NewHost)
}

func TestHostExecutor_Cleanup(t *testing.T) {
	e := &hostExecutor{}
	assert.NoError(t, e.Cleanup(t.Context()))
}
