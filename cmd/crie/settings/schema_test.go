package settings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProjectSchema(t *testing.T) {
	schema := ProjectSchema()

	for _, ref := range linterRefs {
		assert.Contains(t, schema.Definitions, ref)
	}
}
