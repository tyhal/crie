package project

import (
	"github.com/google/jsonschema-go/jsonschema"
)

// Schema generates a jsonscema for a project Config
func Schema() *jsonschema.Schema {
	s, _ := jsonschema.For[Config](nil)
	return s
}
