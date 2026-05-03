// Package main is the main entrypoint for the crie CLI.
package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"

	"github.com/tyhal/crie/internal/cli"
	lintercli "github.com/tyhal/crie/pkg/cli"
)

var version = "latest"

func main() {
	if err := fang.Execute(context.Background(), cli.RootCmd, fang.WithVersion(version)); err != nil {
		os.Exit(1)
	}
}

func init() {
	lintercli.Version = version
}
