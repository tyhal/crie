package cli

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	Run: func(_ *cobra.Command, _ []string) {
		err := crieRun.NamedMatches.Show(os.Stdout)
		if err != nil {
			log.Fatal(fmt.Errorf("crie list failed: %w", err))
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
	RunE: func(cmd *cobra.Command, _ []string) error {
		err := crieRun.NoStandards(cmd.OutOrStdout())
		if err != nil {
			err = fmt.Errorf("finding unassociated files: %w", err)
			log.Error(err)
		}
		return err
	},
}
