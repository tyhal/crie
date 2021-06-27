package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/pkg/crie/project"
)

// Config is a reference to an all ready setup configuration that these commands will utilise
var Config *project.LintConfiguration

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run formatters",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Config.Run("fmt")

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
		Config.List()
	},
}

// ChkCmd Run all code checking commands
var ChkCmd = &cobra.Command{
	Use:     "chk",
	Aliases: []string{"check"},
	Short:   "Run checkers",
	Long:    `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Config.Run("chk")

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
		Config.NoStandards()
	},
}

func stage(stageName string) {
	log.Info("❨ " + stageName + " ❩")
	err := Config.Run(stageName)
	if err != nil {
		if Config.ContinueOnError {
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

		log.Info("❨ proj ❩")
		Config.CheckProjects()
	},
}
