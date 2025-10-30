package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyhal/crie/internal/config/language"
	"github.com/tyhal/crie/internal/config/project"
	"gopkg.in/yaml.v3"
)

func completeYml(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
	return []cobra.Completion{"yml", "yaml"}, cobra.ShellCompDirectiveFilterFileExt
}

// ConfCmd is used to show configuration settings after flags, env, configs are parsed
var ConfCmd = &cobra.Command{
	Use:               "conf",
	Aliases:           []string{"config", "cnf"},
	Short:             "Print configuration settings",
	Long:              "Print what crie has parsed from flags, env, the project file, and then defaults",
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	RunE: func(_ *cobra.Command, _ []string) error {
		var projectConfig project.Config
		err := viper.Unmarshal(&projectConfig)
		if err != nil {
			return err
		}
		marshal, err := yaml.Marshal(projectConfig)
		if err != nil {
			return err
		}

		fmt.Println(string(marshal))
		return nil
	},
}

// InitCmd command will create a project project file for crieRun
var InitCmd = &cobra.Command{
	Use:               "init",
	Short:             "Create optional config files",
	Long:              `Create an optional project file and an extra optional language override file`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	RunE: func(_ *cobra.Command, _ []string) error {

		err := language.NewLanguageConfigFile(viper.GetString("Language.Conf"))
		if err != nil {
			return err
		}
		fmt.Printf("new language file created: %s\nused to overide crie internal language settings (optional / can be deleted)\n", viper.GetString("Language.Conf"))

		fmt.Println()

		var projectConfig project.Config
		err = viper.Unmarshal(&projectConfig)
		if err != nil {
			return err
		}
		err = projectConfig.NewProjectConfigFile(viper.GetString("Project.Conf"))
		if err != nil {
			return err
		}
		fmt.Printf("new project file created: %s\nthis will be treated as your project defaults (overiden by flags and env)\n", viper.GetString("Project.Conf"))
		return nil
	},
}
