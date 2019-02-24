package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
	"os"
)

// Execute the commands that are parsed
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "crie",
	Short: "crie is a formatter for many languages.",
	Long: `crie brings together a vast collection of formatters and linters
to create a handy tool that can prettify any codebase.

	crie: the act of crying and dying at the same time

	"this unformated code makes me want to crie"

	Its more important about picking a standard than it is to pick any certain one.

	Does a good farmer neglect a crop he has planted?
	Does a good teacher overlook even the most humble student?
	Does a good father allow a single child to starve?
	Does a good programmer refuse to maintain his code? `,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&api.Verbose, "verbose", "v", false, "turn on verbose printing for reports")
	rootCmd.PersistentFlags().BoolVarP(&api.Quiet, "quiet", "q", false, "turn off extra prints from failures (suppresses verbose)")
	rootCmd.PersistentFlags().BoolVarP(&api.ContinueOnError, "continue", "e", false, "show all errors rather than stopping at the first")
	rootCmd.PersistentFlags().BoolVarP(&api.GitDiff, "git-diff", "g", false, "Use the last 10 commits to check files")
	rootCmd.PersistentFlags().BoolVarP(&api.CheckIgnores, "ignores", "i", false, "find files in crie.yml which don't need to be there")
	rootCmd.PersistentFlags().StringVar(&api.GlobalState.ConfName, "config", "crie.yml", "config file")
	rootCmd.PersistentFlags().StringVar(&api.SingleLang, "lang", "", "run with only one language (see list for available options)")

	rootCmd.AddCommand(chkCmd)
	rootCmd.AddCommand(chkCmd)
	rootCmd.AddCommand(fmtCmd)
	rootCmd.AddCommand(allCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(non)

	cobra.OnInitialize(api.Initialise)
}

func main() {
	if err := Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
