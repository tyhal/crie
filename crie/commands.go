package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
	"log"
)

func all() {
	err := api.Fmt()

	if err != nil {
		log.Fatal(err)
	}

	err = api.Chk()

	if err != nil {
		log.Fatal(err)
	}
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "fmt -> chk -> chk",
	Long:  `Runs format then lint and finally check`,
	Run: func(cmd *cobra.Command, args []string) {
		all()
	},
}

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Run crie formatters in current dir",
	Long:  `Run all formatters in the list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Fmt()

		if err != nil {
			log.Fatal(err)
		}
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all languages available",
	Long:  `List all languages available and the commands run when used`,
	Run: func(cmd *cobra.Command, args []string) {
		api.List()
	},
}

var chkCmd = &cobra.Command{
	Use:   "chk",
	Short: "Run crie checkers in current dir",
	Long:  `Check all code standards for coding conventions`,
	Run: func(cmd *cobra.Command, args []string) {

		err := api.Chk()

		if err != nil {
			log.Fatal(err)
		}
	},
}

var non = &cobra.Command{
	Use:   "non",
	Short: "List filetypes in this dir crie doesn't support",
	Long:  `List filetypes in this dir crie doesn't support`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
		api.NoStandards()
	},
}
