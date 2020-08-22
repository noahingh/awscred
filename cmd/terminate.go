package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
)

var (
	// TerminateCommand terminate the daemon
	TerminateCommand = &cli.Command{
		Name:  "terminate",
		Usage: "terminate the daemon if the daemon is running.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "homedir",
				Aliases: []string{"dir"},
				Value:   filepath.Join(homeDir, ".awscred"),
				Usage:   "the path of awscred home directory.",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "debug mode.",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				setDebugMode()
			}

			homedir := c.String("homedir")
			if _, err := os.Stat(homedir); os.IsNotExist(err) {
				return fmt.Errorf("there's no dir: %s", homedir)
			}

			return terminate(homedir)
		},
	}
)

func terminate(homedir string) error {
	var (
		b = true
	)
	daemondir := filepath.Join(homedir, "daemon")

	daemonCtx := &daemon.Context{
		PidFileName: filepath.Join(daemondir, "daemon.pid"),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(daemondir, "daemon.log"),
		LogFilePerm: 0640,
		WorkDir:     daemondir,
		Umask:       027,
	}

	daemon.AddFlag(daemon.BoolFlag(&b), syscall.SIGTERM)
	d, err := daemonCtx.Search()
	if err != nil {
		return fmt.Errorf("Unable to search the daemon: %s", err)
	}

	if err := daemon.SendCommands(d); err != nil {
		return fmt.Errorf("Unable to terminate the daemon: %s", err)
	}

	fmt.Println("terminate the daemon.")
	return nil
}
