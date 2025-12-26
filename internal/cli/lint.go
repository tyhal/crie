package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyhal/crie/internal/config/language"
	"github.com/tyhal/crie/internal/runner"
	"github.com/tyhal/crie/pkg/errchain"
)

var crieRun runner.RunConfiguration

// setCrie pushes the Languages to the crie.RunConfiguration
func setCrie(_ *cobra.Command, _ []string) error {

	langs, err := language.LoadFile(languageConfigPath)
	if err != nil {
		return err
	}

	crieLanguages, err := langs.ToRunFormat()
	if err != nil {
		return err
	}

	var ignore *regexp.Regexp
	if projectConfig.Ignore != nil && len(projectConfig.Ignore) > 0 {
		ignore = regexp.MustCompile(strings.Join(projectConfig.Ignore, "|"))
	}

	crieRun = runner.RunConfiguration{Options: projectConfig.Lint, Ignore: ignore, Languages: crieLanguages}

	return nil
}

func addLintCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&projectConfig.Lint.Continue, "continue", "a", false, "show all errors rather than stopping at the first")
	errFatal(viper.BindPFlag("Lint.Continue", cmd.PersistentFlags().Lookup("continue")))

	cmd.PersistentFlags().BoolVarP(&projectConfig.Lint.Passes, "passes", "p", false, "show files that passed")
	errFatal(viper.BindPFlag("Lint.Passes", cmd.PersistentFlags().Lookup("passes")))

	cmd.PersistentFlags().BoolVarP(&projectConfig.Lint.GitDiff, "git-diff", "g", false, "only check files changed in git")
	errFatal(viper.BindPFlag("Lint.GitDiff", cmd.PersistentFlags().Lookup("git-diff")))

	cmd.PersistentFlags().StringVarP(&projectConfig.Lint.GitTarget, "git-target", "t", "", "a target branch to compare against e.g 'remote/branch' or 'branch'")
	errFatal(viper.BindPFlag("Lint.GitTarget", cmd.PersistentFlags().Lookup("git-target")))
	errFatal(cmd.RegisterFlagCompletionFunc("git-target", completeGitTarget))

	cmd.PersistentFlags().StringVar(&projectConfig.Lint.Only, "only", "", "run with only one language (see `crie ls` for available options)")
	errFatal(viper.BindPFlag("Lint.Only", cmd.PersistentFlags().Lookup("only")))
	errFatal(cmd.RegisterFlagCompletionFunc("only", completeLanguage))

	cmd.PreRunE = setCrie

	RootCmd.AddCommand(cmd)
}

// ChkCmd Run all code checking commands
var ChkCmd = &cobra.Command{
	Use:               "chk",
	Aliases:           []string{"check"},
	Short:             "Run linters that only check code",
	Long:              `Check all code standards for coding conventions`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		err := crieRun.Run(cmd.Context(), runner.LintTypeChk)
		if err != nil {
			return errchain.From(err).Link("crie check")
		}
		return nil
	},
}

// FmtCmd Format code command
var FmtCmd = &cobra.Command{
	Use:               "fmt",
	Short:             "Run formatters",
	Long:              `Run all formatters in the list`,
	Args:              cobra.NoArgs,
	SilenceUsage:      true,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	RunE: func(cmd *cobra.Command, _ []string) error {
		err := crieRun.Run(cmd.Context(), runner.LintTypeFmt)
		if err != nil {
			return errchain.From(err).Link("crie format")
		}
		return nil
	},
}

func stage(ctx context.Context, lintType runner.LintType) error {
	log.Infof("❨ %s ❩", lintType.String())
	err := crieRun.Run(ctx, lintType)
	if err != nil {
		return errchain.From(err).LinkF("crie %s", lintType)
	}
	return nil
}

// LntCmd Runs all commands
var LntCmd = &cobra.Command{
	Use:               "lnt",
	Aliases:           []string{"lint", "all"},
	Short:             "Runs both fmt and then chk",
	Long:              `Runs both format and then check`,
	Args:              cobra.NoArgs,
	ValidArgsFunction: cobra.FixedCompletions(nil, cobra.ShellCompDirectiveNoFileComp),
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		stages := []runner.LintType{runner.LintTypeFmt, runner.LintTypeChk}
		var failedStages []string

		for _, lintType := range stages {
			if err := stage(cmd.Context(), lintType); err != nil {
				if crieRun.Options.Continue {
					failedStages = append(failedStages, lintType.String())
				} else {
					return err
				}
			}
		}

		var err error
		if len(failedStages) > 0 {
			err = fmt.Errorf("crie stages failed: %s", strings.Join(failedStages, ", "))
			log.Error(err)
		}
		return err
	},
}
