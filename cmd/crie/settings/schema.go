package settings

import (
	"github.com/invopop/jsonschema"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/noop"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"maps"
)

// these references are used by ConfigLinter to give hints to configuring linter implementations when its only an interface
var linterRefs = []string{"LintCli", "LintShfmt", "LintNoop"}

func ProjectSchema() *jsonschema.Schema {
	schema := jsonschema.Reflect(&ConfigProject{})

	// Add the definitions for each implementation of a crie Linter

	// LintCli
	maps.Copy(schema.Definitions, jsonschema.Reflect(&cli.LintCli{}).Definitions)
	// LintShfmt
	maps.Copy(schema.Definitions, jsonschema.Reflect(&shfmt.LintShfmt{}).Definitions)
	// LintNoop
	maps.Copy(schema.Definitions, jsonschema.Reflect(&noop.LintNoop{}).Definitions)

	return schema
}
