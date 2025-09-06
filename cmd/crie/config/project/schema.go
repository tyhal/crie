package project

import (
	"github.com/invopop/jsonschema"
)

// Schema generates a jsonscema for a project Config
func Schema() *jsonschema.Schema {
	return jsonschema.Reflect(&Config{})
}
