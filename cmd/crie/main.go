package main

import (
	"github.com/tyhal/crie/cmd/crie/config/language"
	"os"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	crie_cmds "github.com/tyhal/crie/cmd/crie/cmd"
	"github.com/tyhal/crie/cmd/crie/config/project"
)

//`
//	|> crie: the act of crying and dying at the same time
//
//	|> "this unformated code makes me want to crie"
//
//	|> It's more important about picking a standard than it is to pick any certain one.
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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		viper.AddConfigPath(projectConfigPath)
		viper.SetConfigType("yml")
		viper.SetEnvPrefix("CRIE")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		_ = viper.ReadInConfig()
		err := viper.Unmarshal(&projectConfig)
		if err != nil {
			return err
		}

		setLogging()

		return nil
	},
}

var languageConfigPath string

var projectConfigPath string
var projectConfig project.Config

func setCrieConfig(cmd *cobra.Command, args []string) error {
	langConfig, err := language.LoadFile(languageConfigPath)
	if err != nil {
		return err
	}
	crie_cmds.SetCrie(crie_cmds.SaveConfiguration(&projectConfig, langConfig))
	return nil
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
	if projectConfig.Log.Trace {
		log.SetLevel(log.TraceLevel)
	}
	if projectConfig.Log.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if projectConfig.Log.Quiet {
		log.SetLevel(log.FatalLevel)
	}
	if projectConfig.Log.JSON {
		log.SetFormatter(&log.JSONFormatter{})
		projectConfig.Lint.StrictLogging = true
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
	cmd.PersistentFlags().BoolVarP(&projectConfig.Lint.Continue, "continue", "e", false, "show all errors rather than stopping at the first")
	errFatal(viper.BindPFlag("Lint.Continue", cmd.PersistentFlags().Lookup("continue")))
	cmd.PersistentFlags().BoolVarP(&projectConfig.Lint.Passes, "passes", "p", false, "show files that passed")
	errFatal(viper.BindPFlag("Lint.Casses", cmd.PersistentFlags().Lookup("passes")))

	cmd.PersistentFlags().BoolVarP(&projectConfig.Lint.GitDiff, "git-diff", "g", false, "only use files from the current commit to (git-target)")
	errFatal(viper.BindPFlag("Lint.GitDiff", cmd.PersistentFlags().Lookup("git-diff")))
	cmd.PersistentFlags().StringVarP(&projectConfig.Lint.GitTarget, "git-target", "t", "origin/main", "the branch to compare against to find changed files")
	errFatal(viper.BindPFlag("Lint.GitTarget", cmd.PersistentFlags().Lookup("git-target")))

	cmd.PersistentFlags().StringVar(&projectConfig.Lint.Lang, "lang", "", "run with only one language (see `crie ls` for available options)")
	errFatal(viper.BindPFlag("Lint.Lang", cmd.PersistentFlags().Lookup("lang")))

	cmd.PreRunE = setCrieConfig

	rootCmd.AddCommand(cmd)
}

// addCrieCommand is the same as addLintCommand but only to ensure the languages are loaded
func addCrieCommand(cmd *cobra.Command) {
	cmd.PreRunE = setCrieConfig
	rootCmd.AddCommand(cmd)
}

func errFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&projectConfigPath, "project-config", "crie.yml", "optional project config location")
	errFatal(viper.BindPFlag("projectConfigPath", rootCmd.PersistentFlags().Lookup("project-config")))
	rootCmd.PersistentFlags().StringVar(&languageConfigPath, "language-config", "crie_lang.yml", "optional language override config location")
	errFatal(viper.BindPFlag("languageConfigPath", rootCmd.PersistentFlags().Lookup("language-config")))

	rootCmd.PersistentFlags().BoolVarP(&projectConfig.Log.JSON, "json", "j", projectConfig.Log.JSON, "turn on json output")
	errFatal(viper.BindPFlag("Log.JSON", rootCmd.PersistentFlags().Lookup("json")))

	rootCmd.PersistentFlags().BoolVarP(&projectConfig.Log.Verbose, "verbose", "v", projectConfig.Log.Verbose, "turn on verbose printing for reports")
	errFatal(viper.BindPFlag("Log.Verbose", rootCmd.PersistentFlags().Lookup("verbose")))

	rootCmd.PersistentFlags().BoolVarP(&projectConfig.Log.Quiet, "quiet", "q", projectConfig.Log.Quiet, "turn off extra prints from failures (suppresses verbose)")
	errFatal(viper.BindPFlag("Log.Trace", rootCmd.PersistentFlags().Lookup("quiet")))

	rootCmd.PersistentFlags().BoolVar(&projectConfig.Log.Trace, "trace", projectConfig.Log.Trace, "turn on trace printing for reports")
	errFatal(viper.BindPFlag("Log.Trace", rootCmd.PersistentFlags().Lookup("trace")))
	errFatal(rootCmd.PersistentFlags().MarkHidden("trace"))

	addLintCommand(crie_cmds.ChkCmd)
	addLintCommand(crie_cmds.FmtCmd)
	addLintCommand(crie_cmds.LntCmd)

	addLintCommand(crie_cmds.ConfCmd)
	addLintCommand(crie_cmds.InitCmd)

	addCrieCommand(crie_cmds.NonCmd)
	addCrieCommand(crie_cmds.LsCmd)

	rootCmd.AddCommand(crie_cmds.SchemaCmd)
	crie_cmds.SchemaCmd.AddCommand(crie_cmds.SchemaLangCmd)
	crie_cmds.SchemaCmd.AddCommand(crie_cmds.SchemaProjectCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
