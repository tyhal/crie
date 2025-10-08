package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodman_SocketGet(t *testing.T) {
	if willPodman() != nil {
		t.Skip()
	}

	socket, err := getPodmanMachineSocket()
	assert.NoError(t, err)
	assert.NotEqual(t, "", socket)
}
