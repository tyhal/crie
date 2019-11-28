package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
	"github.com/tyhal/crie/cli"
)

// Execute the commands that are parsed
func Execute() error {
	return rootCmd.Execute()
}

var quote = `
	|> crie: the act of crying and dying at the same time

	|> "this unformated code makes me want to crie"

	|> Its more important about picking a standard than it is to pick any certain one.

	>>-
		Does a good farmer neglect a crop he has planted?
		Does a good teacher overlook even the most humble student?
		Does a good father allow a single child to starve?
		Does a good programmer refuse to maintain his code? 
	>>-
`

var rootCmd = &cobra.Command{
	Use:   "crie",
	Short: "crie is a formatter for many languages.",
	Long: `
	
	crie brings together a vast collection of formatters and linters
	to create a handy tool that can prettify any codebase.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&api.Verbose, "verbose", "v", false, "turn on verbose printing for reports")
	rootCmd.PersistentFlags().BoolVarP(&api.Quiet, "quiet", "q", false, "turn off extra prints from failures (suppresses verbose)")
	rootCmd.PersistentFlags().BoolVarP(&api.ContinueOnError, "continue", "e", false, "show all errors rather than stopping at the first")
	rootCmd.PersistentFlags().BoolVarP(&api.GitDiff, "git-diff", "g", false, "use the last 10 commits to check files")
	rootCmd.PersistentFlags().StringVar(&api.GlobalState.ConfName, "config", "crie.yml", "config file location")
	rootCmd.PersistentFlags().StringVar(&api.SingleLang, "lang", "", "run with only one language (see list for available options)")

	rootCmd.AddCommand(cli.ChkCmd)
	rootCmd.AddCommand(cli.FmtCmd)
	rootCmd.AddCommand(cli.AllCmd)
	rootCmd.AddCommand(cli.VersionCmd)
	rootCmd.AddCommand(cli.LsCmd)
	rootCmd.AddCommand(cli.NonCmd)

	cobra.OnInitialize(api.Initialise)
}

func main() {
	if err := Execute(); err != nil {
		log.Fatal(err)
	}
}
