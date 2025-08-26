package language

import (
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/noop"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"gopkg.in/yaml.v3"
)

// ConfigLinter attaches a type discriminator field to make a Crie Linter implementation yaml parsable
type ConfigLinter struct {
	linter.Linter
}

func (c ConfigLinter) JSONSchema() *jsonschema.Schema {

	var schema jsonschema.Schema

	schema.OneOf = make([]*jsonschema.Schema, 3)
	// linterRefs are manually added from schema.go
	for i, ref := range linterRefs {
		schema.OneOf[i] = &jsonschema.Schema{
			Ref: fmt.Sprintf("#/$defs/%s", ref),
		}
	}

	return &schema
}

func decodeLinter[T linter.Linter](value *yaml.Node, dst *linter.Linter) error {
	var src T
	if err := value.Decode(&src); err != nil {
		return err
	}
	*dst = src
	return nil
}

// UnmarshalYAML implements custom YAML unmarshalling
func (cl *ConfigLinter) UnmarshalYAML(value *yaml.Node) error {
	var typeOnly struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&typeOnly); err != nil {
		return err
	}

	switch typeOnly.Type {
	case "cli":
		return decodeLinter[*cli.LintCli](value, &cl.Linter)
	case "shfmt":
		return decodeLinter[*shfmt.LintShfmt](value, &cl.Linter)
	case "noop":
		return decodeLinter[*noop.LintNoop](value, &cl.Linter)
	case "":
		return fmt.Errorf("field missing 'type'")
	default:
		return fmt.Errorf("unknown linter type: %s", typeOnly.Type)
	}
}
