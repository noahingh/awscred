package subcmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	pb "github.com/hanjunlee/awscred/api"
	"github.com/hanjunlee/awscred/internal/daemon/server"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	defaultPort = 5126
)

var (
	homeDir, _ = homedir.Dir()

	// RunCommand run a Daemon.
	RunCommand = &cli.Command{
		Name:  "run",
		Usage: "start a daemon to reflect session tokens on a new credentials.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "credentials",
				Aliases: []string{"c"},
				Value:   filepath.Join(homeDir, ".aws", "credentials"),
				Usage:   "the path of aws credentials file.",
			},
			&cli.StringFlag{
				Name:    "homedir",
				Aliases: []string{"dir"},
				Value:   filepath.Join(homeDir, ".awscred"),
				Usage:   "the path of awscred home directory.",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   defaultPort,
				Usage:   "the port of server.",
			},
			&cli.BoolFlag{
				Name:    "server-mode",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "run as the gRPC server, not daemon.",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "debug mode.",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				err      error
				origCred string
				cred     string
				conf     string
				port     string
			)
			if c.Bool("debug") {
				setDebugMode()
			}

			if origCred, err = filepath.Abs(c.String("credentials")); err != nil {
				return fmt.Errorf("failed to the abs path of aws credentials: %s", err)
			}

			homedir := c.String("homedir")
			if _, err := os.Stat(homedir); os.IsNotExist(err) {
				os.Mkdir(homedir, 0755)
			}

			if cred, err = filepath.Abs(filepath.Join(homedir, "credentials")); err != nil {
				return fmt.Errorf("failed to the abs path of awscred credentials: %s", err)
			}

			if conf, err = filepath.Abs(filepath.Join(homedir, "config")); err != nil {
				return fmt.Errorf("failed to the abs path of awscred config: %s", err)
			}

			port = ":" + strconv.Itoa(c.Int("port"))
			return run(origCred, cred, conf, port)
		},
	}
)

func run(orig, cred, conf, port string) error {
	log.Infof("start a daemon: [aws credentials: \"%s\", awscred credentials: \"%s\", awscred config: \"%s\"]", orig, cred, conf)

	i := server.NewInteractor(orig, cred, conf)
	i.StartWatch(context.Background())
	i.Reflect()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	// TODO: support a daemon.
	s := grpc.NewServer()
	pb.RegisterAWSCredServer(s, server.NewServer(i))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to server: %s", err)
	}
	return nil
}
