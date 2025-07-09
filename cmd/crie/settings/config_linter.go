package settings

import (
	"fmt"
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/noop"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"gopkg.in/yaml.v3"
)

// ConfigLinter attaches a type discriminator field to make a Crie Linter implementation yaml parsable
type ConfigLinter struct {
	Type string `yaml:"type"`
	linter.Linter
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

	cl.Type = typeOnly.Type

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
