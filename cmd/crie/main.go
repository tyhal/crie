package main

import "os"
import "github.com/tyhal/crie/internal/cli"

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
