package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tyhal/crie/api"
	"strconv"
)

var majorNum = "0"
var minorOffset = 0
var patchNum = "6"

var major bool
var minor bool
var patch bool

// VersionCmd Print the current version for the crie command line
var VersionCmd = &cobra.Command{
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
	VersionCmd.Flags().BoolVar(&major, "major", false, "Get major number")
	VersionCmd.Flags().BoolVar(&minor, "minor", false, "Get minor number")
	VersionCmd.Flags().BoolVar(&patch, "patch", false, "Get patch number")
}
