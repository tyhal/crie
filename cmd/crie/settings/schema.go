package settings

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/tyhal/crie/pkg/linter/cli"
	"github.com/tyhal/crie/pkg/linter/noop"
	"github.com/tyhal/crie/pkg/linter/shfmt"
	"maps"
)

func PrintProjectSchema() {
	schema := jsonschema.Reflect(&ConfigProject{})
	cliSchema := jsonschema.Reflect(&cli.LintCli{})
	shfmtSchema := jsonschema.Reflect(&shfmt.LintShfmt{})
	noopSchema := jsonschema.Reflect(&noop.LintNoop{})
	
	maps.Copy(schema.Definitions, cliSchema.Definitions)
	maps.Copy(schema.Definitions, shfmtSchema.Definitions)
	maps.Copy(schema.Definitions, noopSchema.Definitions)

	jsonBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(jsonBytes))
}
