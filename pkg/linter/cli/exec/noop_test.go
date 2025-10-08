package exec

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopExecutor_Setup(t *testing.T) {
	executor := NoopExecutor{}
	assert.NoError(t, executor.Setup())
}

func TestNoopExecutor_Exec(t *testing.T) {
	executor := NoopExecutor{}

	var outB, errB bytes.Buffer

	err := executor.Exec(ExecInstance{}, "test.txt", &outB, &errB)
	assert.NoError(t, err)

	assert.Equal(t, "stdout", outB.String())
	assert.Equal(t, "stderr", errB.String())
}

func TestNoopExecutor_Cleanup(t *testing.T) {
	executor := NoopExecutor{}
	assert.NoError(t, executor.Cleanup())
}
