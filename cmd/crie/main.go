package main

import (
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/cmd"
	"github.com/tyhal/crie/cmd/crie/config/language"
	"github.com/tyhal/crie/cmd/crie/config/project"
	"github.com/tyhal/crie/cmd/crie/config/run_instance"
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

		if err := projectConfig.LoadFile(); err != nil {
			log.Fatalf("Failed to load project config: %v", err)
		}

		langConfig, err := language.LoadFile(languageConfigPath)
		if err != nil {
			log.Fatalf("Failed to load language config: %v", err)
		}

		run_instance.SaveConfiguration(&projectConfig, langConfig)
	},
}

var languageConfigPath string

// Stuttering AF
var projectConfigPath string
var projectConfig project.Config

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
	if projectConfig.Trace {
		log.SetLevel(log.TraceLevel)
	}
	if projectConfig.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if projectConfig.Quiet {
		log.SetLevel(log.FatalLevel)
	}
	if projectConfig.JSON {
		log.SetFormatter(&log.JSONFormatter{})
		run_instance.Crie.StrictLogging = true
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
	cmd.PersistentFlags().BoolVarP(&run_instance.Crie.ContinueOnError, "continue", "e", false, "show all errors rather than stopping at the first")
	cmd.PersistentFlags().BoolVarP(&run_instance.Crie.ShowPasses, "passes", "p", false, "show files that passed")

	cmd.PersistentFlags().BoolVarP(&run_instance.Crie.GitDiff, "git-diff", "g", false, "only use files from the current commit to (git-target)")
	cmd.PersistentFlags().StringVarP(&run_instance.Crie.GitTarget, "git-target", "t", "origin/main", "the branch to compare against to find changed files")

	cmd.PersistentFlags().StringVar(&run_instance.Crie.SingleLang, "lang", "", "run with only one language (see `crie ls` for available options)")

	rootCmd.AddCommand(cmd)
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(setLogging)

	rootCmd.PersistentFlags().BoolVarP(&projectConfig.JSON, "json", "j", projectConfig.JSON, "turn on json output")
	rootCmd.PersistentFlags().BoolVarP(&projectConfig.Verbose, "verbose", "v", projectConfig.Verbose, "turn on verbose printing for reports")
	rootCmd.PersistentFlags().BoolVarP(&projectConfig.Quiet, "quiet", "q", projectConfig.Quiet, "turn off extra prints from failures (suppresses verbose)")
	rootCmd.PersistentFlags().BoolVarP(&run_instance.Crie.StrictLogging, "strict-logging", "s", false, "ensure all messages use the structured logger (set true if using json output)")
	rootCmd.PersistentFlags().StringVar(&projectConfigPath, "project-config", "crie.yml", "project config location")
	rootCmd.PersistentFlags().StringVar(&languageConfigPath, "language-config", "crie_lang.yml", "language override config location")

	rootCmd.PersistentFlags().BoolVar(&projectConfig.Trace, "trace", projectConfig.Trace, "turn on trace printing for reports")
	err := rootCmd.PersistentFlags().MarkHidden("trace")
	if err != nil {
		log.Fatal(err)
	}

	addLintCommand(cmd.ChkCmd)
	addLintCommand(cmd.FmtCmd)
	addLintCommand(cmd.LntCmd)

	//rootCmd.AddCommand(cmd.InitCmd)
	//rootCmd.AddCommand(cmd.ConfCmd)
	rootCmd.AddCommand(cmd.NonCmd)
	rootCmd.AddCommand(cmd.LsCmd)

	rootCmd.AddCommand(cmd.SchemaCmd)
	cmd.SchemaCmd.AddCommand(cmd.SchemaLangCmd)
	cmd.SchemaCmd.AddCommand(cmd.SchemaProjectCmd)
}

func initConfig() {
	// 1. crie cli project
	// 2. project project
	// 	1. crie language override project
	// 	2. ignore file project
	// 3. crie's internal project

	// crie cli project do map to crie's internal project too
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
