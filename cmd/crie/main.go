package main

import "os"
import "github.com/tyhal/crie/internal/cli"

var version = "dev"

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cli.RootCmd.Version = version
}
