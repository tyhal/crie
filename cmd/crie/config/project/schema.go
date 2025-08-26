package project

import (
	"github.com/invopop/jsonschema"
)

// ProjectSchema
func ProjectSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Config{})
}
