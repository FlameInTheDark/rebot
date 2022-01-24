package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	// build is the git version of this program. It is set using build flags in the makefile.
	build = "develop"
)

func main() {
	app := &cli.App{
		Name:    "weather",
		Usage:   "run weather command service",
		Version: build,
		Commands: []*cli.Command{
			RunCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
