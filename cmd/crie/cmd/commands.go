package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tyhal/crie/pkg/crie"
	"github.com/tyhal/crie/pkg/crie/linter"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/config/language"
	"github.com/tyhal/crie/cmd/crie/config/project"
)

var crieRun crie.RunConfiguration

// SetCrie is used to ensure the entire crie configuration is set together
func SetCrie(cr crie.RunConfiguration) {
	crieRun = cr
}

// SaveConfiguration pushes the Languages to the crie.RunConfiguration
func SaveConfiguration(proj *project.Config, langs *language.Languages) crie.RunConfiguration {
	crieConfig := crie.RunConfiguration{
		ContinueOnError: proj.Lint.Continue,
		ShowPasses:      proj.Lint.Passes,
		GitDiff:         proj.Lint.GitDiff,
		GitTarget:       proj.Lint.GitTarget,
		SingleLang:      proj.Lint.Lang,
	}

	crieConfig.Ignore = regexp.MustCompile(strings.Join(proj.Ignore, "|"))

	crieConfig.Languages = make(map[string]*linter.Language, len(langs.Languages))
	for langName, lang := range langs.Languages {
		crieConfig.Languages[langName] = lang.ToCrieLanguage()
	}

	return crieConfig
}

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run formatters",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := crieRun.Run("fmt")

		if err != nil {
			log.Fatal(err)
		}
	},
}

// LsCmd List support languages command
var LsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List languages",
	Long:    `List all languages available and the commands run when used`,
	Run: func(cmd *cobra.Command, args []string) {
		crieRun.List()
	},
}

// ChkCmd Run all code checking commands
var ChkCmd = &cobra.Command{
	Use:     "chk",
	Aliases: []string{"check"},
	Short:   "Run checkers",
	Long:    `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {
		err := crieRun.Run("chk")

		if err != nil {
			log.Fatal(err)
		}
	},
}

// NonCmd List every type of file that just passes through
var NonCmd = &cobra.Command{
	Use:     "non",
	Aliases: []string{"not-linted"},
	Short:   "List what isn't supported for this project",
	Long: `List what isn't supported for this project

Find the file extensions that dont have an associated regex match within crieRun`,
	Run: func(cmd *cobra.Command, args []string) {
		crieRun.NoStandards()
	},
}

// InitCmd command will create a project project file for crieRun
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an optional project project file",
	Long:  `Create an optional project project file`,

	RunE: func(_ *cobra.Command, _ []string) error {

		err := language.NewLanguageConfigFile(viper.GetString("languageConfigPath"))
		if err != nil {
			return err
		}
		fmt.Printf("new language file created: %s\nused to overide crie internal language settings (optional / can be deleted)\n", viper.GetString("languageConfigPath"))

		fmt.Println()

		var projectConfig project.Config
		err = viper.Unmarshal(&projectConfig)
		if err != nil {
			return err
		}
		err = projectConfig.NewProjectConfigFile(viper.GetString("projectConfigPath"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("new project file created: %s\nthis will be treated as your project defaults (overiden by flags and env)\n", viper.GetString("projectConfigPath"))
		return nil
	},
}

// ConfCmd is used to show configuration settings after flags, env, configs are parsed
var ConfCmd = &cobra.Command{
	Use:     "conf",
	Aliases: []string{"config", "cnf"},
	Short:   "Print what crie has parsed from flags, env, the project file, and then defaults",
	Long:    "Print what crie has parsed from flags, env, the project file, and then defaults",
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

// SchemaCmd is used to hold jsonschema generator commands
var SchemaCmd = &cobra.Command{
	Use:     "schema",
	Aliases: []string{"sch"},
	Short:   "Print jsonschema's of the crie cli configs",
	Long:    `Print jsonschema's' of the crie cli configs`,
}

// SchemaLangCmd is the jsonschema generator for the crie languages configuration
var SchemaLangCmd = &cobra.Command{
	Use:     "lang",
	Aliases: []string{"language", "lng"},
	Short:   "Print the schema for language's configurations",
	Long:    `Print json schema for cries configuration format used override language project`,
	Run: func(_ *cobra.Command, _ []string) {
		schema := language.Schema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	},
}

// SchemaProjectCmd is the jsonschema generator for the crie project configuration
var SchemaProjectCmd = &cobra.Command{
	Use:     "proj",
	Aliases: []string{"project", "prj"},
	Short:   "Print the schema for language's configurations",
	Long:    `Print json schema for cries configuration format used override language project`,
	Run: func(_ *cobra.Command, _ []string) {
		schema := project.Schema()
		jsonBytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(jsonBytes))
	},
}

func stage(stageName string) {
	log.Info("❨ " + stageName + " ❩")
	err := crieRun.Run(stageName)
	if err != nil {
		if crieRun.ContinueOnError {
			log.Error(err)
		} else {
			log.Fatal(err)
		}
	}
}

// LntCmd Runs all commands
var LntCmd = &cobra.Command{
	Use:     "lnt",
	Aliases: []string{"lint", "all"},
	Short:   "Run everything",
	Long:    `Runs both format and then check`,
	Run: func(_ *cobra.Command, _ []string) {
		stage("fmt")
		stage("chk")
	},
}
