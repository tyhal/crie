package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/cmd"
	"github.com/tyhal/crie/cmd/crie/conf"
	"github.com/tyhal/crie/pkg/crie"
	"os"
	"sort"
)

//`
//	|> crie: the act of crying and dying at the same time
//
//	|> "this unformated code makes me want to crie"
//
//	|> Its more important about picking a standard than it is to pick any certain one.
//
//	>>-
//		Does a good farmer neglect a crop he has planted?
//		Does a good teacher overlook even the most humble student?
//		Does a good father allow a single child to starve?
//		Does a good programmer refuse to maintain his code?
//	>>-
//`

var rootCmd = &cobra.Command{
	Use:     "crie",
	Short:   "crie is a formatter for many languages.",
	Example: "check all python files in the last commit 'crie chk --git-diff 1 --lang python'",
	Long: `
	crie brings together a vast collection of formatters and linters
	to create a handy tool that can prettify any codebase.`,
}

var quiet = false
var verbose = false
var trace = false
var json = false
var state crie.RunConfiguration

func msgLast(fields []string) {
	sort.Slice(fields, func(i, j int) bool {
		if fields[i] == "msg" {
			return false
		}
		if fields[j] == "msg" {
			return true
		}
		return fields[i] < fields[j]
	})
}

func setLogging() {
	if trace {
		log.SetLevel(log.TraceLevel)
	}
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	if quiet {
		log.SetLevel(log.FatalLevel)
	}
	if json {
		log.SetFormatter(&log.JSONFormatter{})
		state.StrictLogging = true
	} else {
		log.SetFormatter(&log.TextFormatter{
			SortingFunc:      msgLast,
			DisableQuote:     true,
			DisableTimestamp: true,
			DisableSorting:   false,
		})
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

	cmd.Config = &state

	rootCmd.PersistentFlags().BoolVarP(&json, "json", "j", json, "turn on json output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "turn on verbose printing for reports")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", quiet, "turn off extra prints from failures (suppresses verbose)")
	rootCmd.PersistentFlags().BoolVarP(&state.StrictLogging, "strict-logging", "s", false, "ensure all messages use the structured logger (set true if using json output)")
	rootCmd.PersistentFlags().StringVar(&state.ConfPath, "config", "crie.yml", "project config file location")

	rootCmd.PersistentFlags().BoolVarP(&trace, "trace", "t", trace, "turn on trace printing for reports")
	err := rootCmd.PersistentFlags().MarkHidden("trace")
	if err != nil {
		log.Fatal(err)
	}

	addLintCommand(cmd.ChkCmd)
	addLintCommand(cmd.FmtCmd)
	addLintCommand(cmd.LntCmd)

	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.NonCmd)
	rootCmd.AddCommand(cmd.LsCmd)

	cobra.OnInitialize(setLogging)
}

func main() {

	state.Languages = conf.LanguageList

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
