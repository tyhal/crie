// Package main is the main entrypoint for the crie CLI.
package main

import (
	"context"
	"os"

	"charm.land/fang/v2"

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
