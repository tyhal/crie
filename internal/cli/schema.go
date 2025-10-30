package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tyhal/crie/internal/config/language"
	"github.com/tyhal/crie/internal/config/project"
)

// SchemaCmd is used to hold jsonschema generator commands
var SchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"sch"},
	Short:   "Print JSON schemas for config files",
	Long:    `Print JSON schemas for config files`,
}

// SchemaLangCmd is the jsonschema generator for the crie languages configuration
var SchemaLangCmd = &cobra.Command{
	Use:               "lang",
	Aliases:           []string{"language", "lng"},
	Short:             "Print the schema for language's configurations",
	Long:              `Print json schema for cries configuration format used override language project`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	Run: func(_ *cobra.Command, _ []string) {
		schema := language.Schema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	},
}

// SchemaProjectCmd is the jsonschema generator for the crie project configuration
var SchemaProjectCmd = &cobra.Command{
	Use:               "proj",
	Aliases:           []string{"project", "prj"},
	Short:             "Print the schema for language's configurations",
	Long:              `Print json schema for cries configuration format used override language project`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	Run: func(_ *cobra.Command, _ []string) {
		schema := project.Schema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	},
}
