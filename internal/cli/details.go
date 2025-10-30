package cli

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/internal/errchain"
)

// addCrieCommand is the same as addLintCommand but only to ensure the languages are loaded
func addCrieCommand(cmd *cobra.Command) {
	cmd.PreRunE = setCrie
	RootCmd.AddCommand(cmd)
}

// LsCmd Show support languages command
var LsCmd = &cobra.Command{
	Use:               "ls",
	Aliases:           []string{"list"},
	Short:             "Show languages",
	Long:              `Show all languages available and the commands run when used`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	Run: func(cmd *cobra.Command, args []string) {
		err := crieRun.Languages.Show(os.Stdout)
		if err != nil {
			log.Fatal(errchain.From(err).Link("crie list failed"))
		}
	},
}

// NonCmd Show every type of file that just passes through
var NonCmd = &cobra.Command{
	Use:     "non",
	Aliases: []string{"not-linted"},
	Short:   "Show what isn't supported for this project",
	Long: `Show what isn't supported for this project

Find the file extensions that dont have an associated regex match within crie`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	Run: func(cmd *cobra.Command, args []string) {
		err := crieRun.NoStandards()
		if err != nil {
			log.Fatal(errchain.From(err).Link("finding unassociated files"))
		}
	},
}
