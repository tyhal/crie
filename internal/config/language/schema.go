package language

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/dockfmt"
	"github.com/tyhal/crie/pkg/linter/noop"
	"github.com/tyhal/crie/pkg/linter/shfmt"
)

// these references are used by Linter to give hints to configuring linter implementations when its only an interface
var linterRefs = []string{"LintCli", "LintShfmt", "LintDockFmt", "LintNoop"}

// Schema is used to generate a complete jsonschema for the Languages struct
func Schema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeFor[Linter](): Linter{}.JSONSchema(),
		},
	}
	schema, err := jsonschema.For[Languages](opts)
	if err != nil {
		panic(err)
	}

	if schema == nil {
		panic("jsonschema.For[Languages] returned nil schema")
	}

	if schema.Defs == nil {
		schema.Defs = make(map[string]*jsonschema.Schema)
	}

	// Add the definitions for each implementation of a crie Linter

	// LintCli
	cliSchema, _ := jsonschema.For[cli.LintCli](nil)
	schema.Defs["LintCli"] = cliSchema
	// LintShfmt
	shfmtSchema, _ := jsonschema.For[shfmt.LintShfmt](nil)
	schema.Defs["LintShfmt"] = shfmtSchema
	// LintDockFmt
	dockfmtSchema, _ := jsonschema.For[dockfmt.LintDockFmt](nil)
	schema.Defs["LintDockFmt"] = dockfmtSchema
	// LintNoop
	noopSchema, _ := jsonschema.For[noop.LintNoop](nil)
	schema.Defs["LintNoop"] = noopSchema

	return schema
}
