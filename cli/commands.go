package cli

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
)

// Config is a reference to an all ready setup configuration that these commands will utilise
var Config *api.ProjectLintConfiguration

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run crie formatters in current dir",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		Config.LintType = "fmt"
		err := Config.Run()

		if err != nil {
			log.Fatal(err)
		}
	},
}

// LsCmd List support languages command
var LsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all languages available",
	Long:    `List all languages available and the commands run when used`,
	Run: func(cmd *cobra.Command, args []string) {
		Config.List()
	},
}

// ChkCmd Run all code checking commands
var ChkCmd = &cobra.Command{
	Use:     "chk",
	Aliases: []string{"check"},
	Short:   "Run crie checkers in current dir",
	Long:    `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {

		err := Config.Chk()

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
	Long: `Find the file extensions that dont
			have an associated regex match within crie`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		Config.NoStandards()
	},
}

// LntCmd Runs all commands
var LntCmd = &cobra.Command{
	Use:     "lnt",
	Aliases: []string{"lint", "all"},
	Short:   "Fully Lint the code",
	Long:    `Runs both format and then check`,
	Run: func(cmd *cobra.Command, args []string) {
		Config.LintType = "fmt"
		err := Config.Run()

		if err != nil {
			log.Fatal(err)
		}

		err = Config.Chk()

		if err != nil {
			log.Fatal(err)
		}
	},
}
