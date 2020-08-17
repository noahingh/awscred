package main

import (
	"os"

	"github.com/hanjunlee/awscred/cmd/subcmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "awscred",
		Usage: "awscred is a tool to generate a AWS session token and manage it",
		Commands: []*cli.Command{
			subcmd.RunCommand,
			subcmd.OnCommand,
			subcmd.OffCommand,
			subcmd.SetCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
