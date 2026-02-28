package language

import (
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/tyhal/crie/pkg/cli"
	"github.com/tyhal/crie/pkg/dockfmt"
	"github.com/tyhal/crie/pkg/linter"
	"github.com/tyhal/crie/pkg/noop"
	"github.com/tyhal/crie/pkg/shfmt"
	"gopkg.in/yaml.v3"
)

// Linter attaches a type discriminator field to make a Crie Linter implementation yaml parsable
type Linter struct {
	linter.Linter
}

// JSONSchema is used to parse a valid jsonschema just for a Linter
func (l Linter) JSONSchema() *jsonschema.Schema {
	var schema jsonschema.Schema

	schema.OneOf = make([]*jsonschema.Schema, len(linterRefs))
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
func (l *Linter) UnmarshalYAML(value *yaml.Node) error {
	var typeOnly struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&typeOnly); err != nil {
		return err
	}

	switch typeOnly.Type {
	case "cli":
		return decodeLinter[*cli.LintCli](value, &l.Linter)
	case "shfmt":
		return decodeLinter[*shfmt.LintShfmt](value, &l.Linter)
	case "dockfmt":
		return decodeLinter[*dockfmt.LintDockFmt](value, &l.Linter)
	case "noop":
		return decodeLinter[*noop.LintNoop](value, &l.Linter)
	case "":
		return fmt.Errorf("field missing 'type'")
	default:
		return fmt.Errorf("unknown linter type: %s", typeOnly.Type)
	}
}
