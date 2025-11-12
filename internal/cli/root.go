package cli

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyhal/crie/internal/config/project"
	"github.com/tyhal/crie/pkg/errchain"
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

// RootCmd is the root Cobra command for the crie CLI.
// It wires configuration, global flags, and subcommands together.
var RootCmd = &cobra.Command{
	Use:   "crie",
	Short: "crie is a formatter and linter for many languages.",
	Example: `
check all files changes since the target branch 
	$ crie chk --git-diff --git-target=origin/main

format all python files
	$ crie fmt --only python
`,
	Long: `
	crie brings together a vast collection of formatters and linters
	to create a handy tool that can sanity check any codebase.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {

		viper.SetConfigFile(projectConfigPath)
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

		// enable git diff if target is set
		if projectConfig.Lint.GitTarget != "" {
			projectConfig.Lint.GitDiff = true
		}

		return nil
	},
}

var languageConfigPath string

var projectConfigPath string
var projectConfig project.Config

func errFatal(err error) {
	if err != nil {
		log.Fatal(errchain.From(err).Link("incorrect viper configuration"))
	}
}

func init() {

	RootCmd.PersistentFlags().StringVarP(&projectConfigPath, "conf", "c", "crie.yml", "project configuration file")
	errFatal(viper.BindPFlag("Project.Conf", RootCmd.PersistentFlags().Lookup("conf")))
	errFatal(RootCmd.RegisterFlagCompletionFunc("conf", completeYml))

	RootCmd.PersistentFlags().StringVarP(&languageConfigPath, "lang-conf", "l", "crie.lang.yml", "language configuration file")
	errFatal(viper.BindPFlag("Language.Conf", RootCmd.PersistentFlags().Lookup("lang-conf")))
	errFatal(RootCmd.RegisterFlagCompletionFunc("lang-conf", completeYml))

	RootCmd.PersistentFlags().BoolVarP(&projectConfig.Log.JSON, "json", "j", projectConfig.Log.JSON, "turn on json output")
	errFatal(viper.BindPFlag("Log.JSON", RootCmd.PersistentFlags().Lookup("json")))

	RootCmd.PersistentFlags().BoolVarP(&projectConfig.Log.Verbose, "verbose", "v", projectConfig.Log.Verbose, "turn on verbose printing for reports")
	errFatal(viper.BindPFlag("Log.Verbose", RootCmd.PersistentFlags().Lookup("verbose")))

	RootCmd.PersistentFlags().BoolVarP(&projectConfig.Log.Quiet, "quiet", "q", projectConfig.Log.Quiet, "only prints critical errors (suppresses verbose)")
	errFatal(viper.BindPFlag("Log.Quiet", RootCmd.PersistentFlags().Lookup("quiet")))

	RootCmd.PersistentFlags().BoolVar(&projectConfig.Log.Trace, "trace", projectConfig.Log.Trace, "turn on trace printing for reports")
	errFatal(viper.BindPFlag("Log.Trace", RootCmd.PersistentFlags().Lookup("trace")))
	errFatal(RootCmd.PersistentFlags().MarkHidden("trace"))

	viper.SetDefault("Ignore", []string{})

	addLintCommand(ChkCmd)
	addLintCommand(FmtCmd)
	addLintCommand(LntCmd)

	addLintCommand(ConfCmd)
	addLintCommand(InitCmd)

	addCrieCommand(NonCmd)
	addCrieCommand(LsCmd)

	RootCmd.AddCommand(SchemaCmd)
	SchemaCmd.AddCommand(SchemaLangCmd)
	SchemaCmd.AddCommand(SchemaProjectCmd)
}
