package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/settings"
)

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run formatters",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := settings.Cli.Crie.Run("fmt")

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
		settings.Cli.Crie.List()
	},
}

// ChkCmd Run all code checking commands
var ChkCmd = &cobra.Command{
	Use:     "chk",
	Aliases: []string{"check"},
	Short:   "Run checkers",
	Long:    `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {
		err := settings.Cli.Crie.Run("chk")

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
		settings.Cli.Crie.NoStandards()
	},
}

// InitCmd command will create a project settings file for Crie
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an optional project settings file",
	Long:  `Create an optional project settings file`,

	Run: func(cmd *cobra.Command, args []string) {
		err := settings.Cli.CreateNewProjectSettings()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var SchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Print jsonschema of the crie project config",
	Long:  `Print jsonschema of the crie project config`,
	Run: func(cmd *cobra.Command, args []string) {
		schema := settings.ProjectSchema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))

	},
}

func stage(stageName string) {
	log.Info("❨ " + stageName + " ❩")
	err := settings.Cli.Crie.Run(stageName)
	if err != nil {
		if settings.Cli.Crie.ContinueOnError {
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
