package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	assert.NotPanics(t, func() {
		Schema()
	})
}
