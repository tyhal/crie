package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"github.com/tyhal/crie/pkg/crie"
	"github.com/tyhal/crie/pkg/crie/linter"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/cmd/crie/config/language"
	"github.com/tyhal/crie/cmd/crie/config/project"
)

var crieRun crie.RunConfiguration

// SetCrie pushes the Languages to the crie.RunConfiguration
func SetCrie(proj *project.Config, langs *language.Languages) {

	languages := make(map[string]*linter.Language, len(langs.Languages))
	for langName, lang := range langs.Languages {
		languages[langName] = lang.ToCrieLanguage()
	}

	var ignore *regexp.Regexp
	if proj.Ignore != nil && len(proj.Ignore) > 0 {
		ignore = regexp.MustCompile(strings.Join(proj.Ignore, "|"))
	}

	crieRun = crie.RunConfiguration{Options: proj.Lint, Ignore: ignore, Languages: languages}
}

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run formatters",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := crieRun.Run("fmt")

		if err != nil {
			log.Fatal(fmt.Errorf("crie format failed: %w", err))
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
	Short:   "Run linters that only check code",
	Long:    `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {
		err := crieRun.Run("chk")

		if err != nil {
			log.Fatal(fmt.Errorf("crie check failed: %w", err))
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
	Short: "Create optional config files",
	Long:  `Create an optional project file and an extra optional language override file`,

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

// ConfCmd is used to show configuration settings after flags, env, configs are parsed
var ConfCmd = &cobra.Command{
	Use:     "conf",
	Aliases: []string{"config", "cnf"},
	Short:   "Print configuration settings",
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
	Short:   "Print JSON schemas for config files",
	Long:    `Print JSON schemas for config files`,
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
		if crieRun.Options.Continue {
			log.Error(err)
		} else {
			log.Fatal(fmt.Errorf("crie %s failed: %w", stageName, err))
		}
	}
}

// LntCmd Runs all commands
var LntCmd = &cobra.Command{
	Use:     "lnt",
	Aliases: []string{"lint", "all"},
	Short:   "Runs both fmt and then chk",
	Long:    `Runs both format and then check`,
	Run: func(_ *cobra.Command, _ []string) {
		stage("fmt")
		stage("chk")
	},
}
