package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	// build is the git version of this program. It is set using build flags in the makefile.
	build = "develop"
)

func main() {
	app := &cli.App{
		Name:  "api",
		Usage: "run api service",
		Version: build,
		Commands: []*cli.Command{
			APICommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
