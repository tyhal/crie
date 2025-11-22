// Package main is the main entrypoint for the crie CLI.
package main

import (
	"os"

	"github.com/tyhal/crie/internal/cli"
	lintercli "github.com/tyhal/crie/pkg/linter/cli"
)

var version = "latest"

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cli.RootCmd.Version = version
	lintercli.Version = version
}
