package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
	"strconv"
)

var majorNum = "0"
var minorOffset = 0
var patchNum = "0"

var major bool
var minor bool
var patch bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version for the crie command line",
	Long:  `Print the current version for the crie command line`,
	Run: func(cmd *cobra.Command, args []string) {
		if major {
			fmt.Println(majorNum)
		} else if minor {
			fmt.Println(strconv.Itoa(api.Version() - minorOffset))
		} else if patch {
			fmt.Println(patchNum)
		} else {
			fmt.Println(majorNum + "." + strconv.Itoa(api.Version()-minorOffset) + "." + patchNum)
		}
	},
}

func init() {
	versionCmd.Flags().BoolVar(&major, "major", false, "Get major number")
	versionCmd.Flags().BoolVar(&minor, "minor", false, "Get minor number")
	versionCmd.Flags().BoolVar(&patch, "patch", false, "Get patch number")
}
