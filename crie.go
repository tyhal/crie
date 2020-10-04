package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
	"github.com/tyhal/crie/cli"
	"github.com/tyhal/crie/imp"
	"os"
	"strconv"
)

var majorNum = "0"
var minorOffset = 0
var patchNum = "44"

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
	Use:     "crie",
	Short:   "crie is a formatter for many languages.",
	Version: majorNum + "." + strconv.Itoa(len(imp.LanguageList)-minorOffset) + "." + patchNum,
	Example: "check all python files in the last commit 'crie chk --git-diff 1 --lang python'",
	Long: `
	crie brings together a vast collection of formatters and linters
	to create a handy tool that can prettify any codebase.`,
}

var quiet = false
var verbose = false
var state api.ProjectLintConfiguration

func setLogLevel() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	if quiet {
		log.SetLevel(log.FatalLevel)
	}
}

func addLintCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&state.ContinueOnError, "continue", "e", false, "show all errors rather than stopping at the first")
	cmd.PersistentFlags().BoolVarP(&state.ShowPasses, "passes", "p", false, "show files that passed")
	cmd.PersistentFlags().IntVarP(&state.GitDiff, "git-diff", "g", 0, "check files that changed in the last X commits")
	cmd.PersistentFlags().StringVar(&state.SingleLang, "lang", "", "run with only one language (see `crie ls` for available options)")

	rootCmd.AddCommand(cmd)
}

func init() {

	cli.Config = &state

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "turn on verbose printing for reports")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", quiet, "turn off extra prints from failures (suppresses verbose)")
	rootCmd.PersistentFlags().StringVar(&state.ConfPath, "config", "crie.yml", "config file location")

	addLintCommand(cli.ChkCmd)
	addLintCommand(cli.FmtCmd)
	addLintCommand(cli.LntCmd)

	rootCmd.AddCommand(cli.NonCmd)
	rootCmd.AddCommand(cli.LsCmd)

	cobra.OnInitialize(setLogLevel)
}

func main() {

	// You could change this to your own implementation of standards
	state.Languages = imp.LanguageList

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
