package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/config/language"
	"github.com/tyhal/crie/cmd/crie/config/project"
	"github.com/tyhal/crie/cmd/crie/config/run_instance"
	"gopkg.in/yaml.v3"
)

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run formatters",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := run_instance.Crie.Run("fmt")

		if err != nil {
			log.Fatal(err)
		}
	},
}

// LsCmd List support languages command
var LsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List languages",
	Long:    `List all languages available and the commands run when used`,
	Run: func(cmd *cobra.Command, args []string) {
		run_instance.Crie.List()
	},
}

// ChkCmd Run all code checking commands
var ChkCmd = &cobra.Command{
	Use:     "chk",
	Aliases: []string{"check"},
	Short:   "Run checkers",
	Long:    `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {
		err := run_instance.Crie.Run("chk")

		if err != nil {
			log.Fatal(err)
		}
	},
}

// NonCmd List every type of file that just passes through
var NonCmd = &cobra.Command{
	Use:     "non",
	Aliases: []string{"not-linted"},
	Short:   "List what isn't supported for this project",
	Long: `List what isn't supported for this project

Find the file extensions that dont have an associated regex match within crie`,
	Run: func(cmd *cobra.Command, args []string) {
		run_instance.Crie.NoStandards()
	},
}

// InitCmd command will create a project project file for Crie
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an optional project project file",
	Long:  `Create an optional project project file`,

	Run: func(cmd *cobra.Command, args []string) {
		err := project.Config.NewProjectConfigFile()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var ConfCmd = &cobra.Command{
	Use:     "conf",
	Aliases: []string{"config", "cnf"},
	Short:   "Print all the currently configured project",
	Long:    "Takes all the information from env, flags, and the project project file to show the complete configuration",
	Run: func(cmd *cobra.Command, args []string) {
		displaySettings, err := yaml.Marshal(project.Config)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(displaySettings))
	},
}

var SchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"sch"},
	Short:   "Print jsonschema of the crie project config",
	Long:    `Print jsonschema of the crie project config`,
}

var SchemaLangCmd = &cobra.Command{
	Use:     "lang",
	Aliases: []string{"language", "lng"},
	Short:   "Print the schema for language's configurations",
	Long:    `Print json schema for cries configuration format used override language project`,
	Run: func(cmd *cobra.Command, args []string) {
		schema := language.LanguagesSchema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	},
}

var SchemaProjectCmd = &cobra.Command{
	Use:     "proj",
	Aliases: []string{"project", "prj"},
	Short:   "Print the schema for language's configurations",
	Long:    `Print json schema for cries configuration format used override language project`,
	Run: func(cmd *cobra.Command, args []string) {
		schema := project.ProjectSchema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	},
}

func stage(stageName string) {
	log.Info("❨ " + stageName + " ❩")
	err := run_instance.Crie.Run(stageName)
	if err != nil {
		if run_instance.Crie.ContinueOnError {
			log.Error(err)
		} else {
			log.Fatal(err)
		}
	}
}

// LntCmd Runs all commands
var LntCmd = &cobra.Command{
	Use:     "lnt",
	Aliases: []string{"lint", "all"},
	Short:   "Run everything",
	Long:    `Runs both format and then check`,
	Run: func(cmd *cobra.Command, args []string) {
		stage("fmt")
		stage("chk")
	},
}
