package language

import (
	"regexp"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"gopkg.in/yaml.v3"
)

// Regex wraps regexp.Regexp with custom YAML unmarshaling
type Regex struct {
	*regexp.Regexp
}

// JSONSchema returns the JSON schema for the Regex type
func (cr Regex) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "array",
		Items: &jsonschema.Schema{
			Type: "string",
		},
	}
}

// UnmarshalYAML implements custom YAML unmarshalling
func (cr *Regex) UnmarshalYAML(value *yaml.Node) error {
	var patterns []string
	if err := value.Decode(&patterns); err != nil {
		return err
	}

	if len(patterns) > 0 {
		compiled := regexp.MustCompile(strings.Join(patterns, "|"))
		cr.Regexp = compiled
	}

	return nil
}
