package main

import (
	"os"

	"github.com/hanjunlee/awscred/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	version = "v0.2.0"
)

func main() {
	app := &cli.App{
		Name:    "awscred",
		Version: version,
		Usage:   "awscred is a tool to generate a AWS session token and manage it",
		Commands: []*cli.Command{
			cmd.RunCommand,
			cmd.TerminateCommand,
			cmd.OnCommand,
			cmd.OffCommand,
			cmd.SetCommand,
			cmd.GenCommand,
			cmd.InfoCommand,
			cmd.ExportCommand,
			cmd.PingCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
