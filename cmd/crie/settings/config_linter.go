package settings

import (
	"fmt"
	"github.com/tyhal/crie/pkg/crie/linter"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"gopkg.in/yaml.v3"
)

// ConfigLinter attaches a type discriminator field to make a Crie Linter implementation yaml parsable
type ConfigLinter struct {
	Type string `yaml:"type"`
	linter.Linter
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

		var cliLinter cli.Lint
		if err := value.Decode(&cliLinter); err != nil {
			return err
		}
		cl.Linter = &cliLinter

	case "shfmt":
		var shfmtLinter shfmt.Lint
		if err := value.Decode(&shfmtLinter); err != nil {
			return err
		}
		cl.Linter = &shfmtLinter

	default:
		return fmt.Errorf("unknown linter type: %s", typeOnly.Type)
	}

	return nil
}
