package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillPodman(t *testing.T) {
	assert.NoError(t, WillPodman())
}

func TestPodman_SocketGet(t *testing.T) {
	if WillPodman() != nil {
		t.Skip()
	}

	socket, err := getPodmanMachineSocket()
	assert.NoError(t, err)
	assert.NotEqual(t, "", socket)
}
