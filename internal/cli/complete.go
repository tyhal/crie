package cli

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/internal/config/language"
	"github.com/tyhal/crie/internal/filelist"
)

func completeYml(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
	return []cobra.Completion{"yml", "yaml"}, cobra.ShellCompDirectiveFilterFileExt
}

func completeGitTarget(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	repo, err := filelist.GetGitRepo(cwd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var comps []cobra.Completion
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	for _, r := range remotes {
		refName := plumbing.ReferenceName("refs/remotes/" + r.Config().Name + "/HEAD")
		if ref, err := repo.Reference(refName, true); err == nil {
			comps = append(comps, ref.Name().Short())
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func completeLanguage(cmd *cobra.Command, _ []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	filter := func(_ language.Language) bool { return true }
	switch cmd.Name() {
	case "fmt":
		filter = func(l language.Language) bool {
			return l.Fmt.Linter != nil
		}
	case "chk":
		filter = func(l language.Language) bool {
			return l.Chk.Linter != nil
		}
	}

	langs, err := language.LoadFile(languageConfigPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var comps []cobra.Completion
	for langName, lang := range langs.Languages {
		if strings.HasPrefix(langName, toComplete) && filter(lang) {
			comps = append(comps, langName)
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}
