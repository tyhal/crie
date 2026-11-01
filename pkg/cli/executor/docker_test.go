package executor

import (
	"net/http/httptest"
	"strings"
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
	testHelperExecutor(t, func() Executor {
		return NewDocker("alpine:latest")
	})
}

func TestAPISpanName(t *testing.T) {
	id := strings.Repeat("ab12", 16)
	tests := map[string]string{
		"/v1.41/containers/" + id + "/exec": "POST /v1.41/containers/{id}/exec",
		"/v1.41/exec/" + id + "/json":       "POST /v1.41/exec/{id}/json",
		"/v1.41/containers/create":          "POST /v1.41/containers/create",
	}
	for path, want := range tests {
		assert.Equal(t, want, apiSpanName("", httptest.NewRequest("POST", path, nil)))
	}
}
