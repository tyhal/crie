package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/cmd"
	"github.com/tyhal/crie/cmd/crie/settings"
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
//		Does a good farmer neglect a crop they have planted?
//		Does a good teacher overlook even the most humble student?
//		Does a good father allow a single child to starve?
//		Does a good programmer refuse to maintain his code?
//	>>-
//`

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "crie",
	Short:   "crie is a formatter and linter for many languages.",
	Example: "check all files changes since the target branch 'crie chk --git-diff --git-target=origin/main --lang python'",
	Long: `
	crie brings together a vast collection of formatters and linters
	to create a handy tool that can prettify any codebase.`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := settings.Cli.LoadConfigFile(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		settings.Cli.SaveConfiguration()
	},
}

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
	if settings.Cli.Trace {
		log.SetLevel(log.TraceLevel)
	}
	if settings.Cli.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if settings.Cli.Quiet {
		log.SetLevel(log.FatalLevel)
	}
	if settings.Cli.JSON {
		log.SetFormatter(&log.JSONFormatter{})
		settings.Cli.Crie.StrictLogging = true
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
	cmd.PersistentFlags().BoolVarP(&settings.Cli.Crie.ContinueOnError, "continue", "e", false, "show all errors rather than stopping at the first")
	cmd.PersistentFlags().BoolVarP(&settings.Cli.Crie.ShowPasses, "passes", "p", false, "show files that passed")

	cmd.PersistentFlags().BoolVarP(&settings.Cli.Crie.GitDiff, "git-diff", "g", false, "only use files from the current commit to (git-target)")
	cmd.PersistentFlags().StringVarP(&settings.Cli.Crie.GitTarget, "git-target", "t", "origin/main", "the branch to compare against to find changed files")

	cmd.PersistentFlags().StringVar(&settings.Cli.Crie.SingleLang, "lang", "", "run with only one language (see `crie ls` for available options)")

	rootCmd.AddCommand(cmd)
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(setLogging)

	rootCmd.PersistentFlags().BoolVarP(&settings.Cli.JSON, "json", "j", settings.Cli.JSON, "turn on json output")
	rootCmd.PersistentFlags().BoolVarP(&settings.Cli.Verbose, "verbose", "v", settings.Cli.Verbose, "turn on verbose printing for reports")
	rootCmd.PersistentFlags().BoolVarP(&settings.Cli.Quiet, "quiet", "q", settings.Cli.Quiet, "turn off extra prints from failures (suppresses verbose)")
	rootCmd.PersistentFlags().BoolVarP(&settings.Cli.Crie.StrictLogging, "strict-logging", "s", false, "ensure all messages use the structured logger (set true if using json output)")
	rootCmd.PersistentFlags().StringVar(&settings.Cli.ConfigPath, "settings", "crie.yml", "project settings file location")

	rootCmd.PersistentFlags().BoolVar(&settings.Cli.Trace, "trace", settings.Cli.Trace, "turn on trace printing for reports")
	err := rootCmd.PersistentFlags().MarkHidden("trace")
	if err != nil {
		log.Fatal(err)
	}

	addLintCommand(cmd.ChkCmd)
	addLintCommand(cmd.FmtCmd)
	addLintCommand(cmd.LntCmd)

	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.SchemaCmd)
	rootCmd.AddCommand(cmd.NonCmd)
	rootCmd.AddCommand(cmd.LsCmd)
}

func initConfig() {
	// TODO one config file, for two purposes means I need to parse it twice partially
	// 1. crie cli settings
	// 2. project settings
	// 	1. crie language override settings
	// 	2. ignore file settings
	// 3. crie's internal settings

	// crie cli settings do map to crie's internal settings too
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
