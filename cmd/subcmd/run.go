package subcmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	pb "github.com/hanjunlee/awscred/api"
	"github.com/hanjunlee/awscred/internal/daemon/server"
	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
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

			// run as server-mode
			if c.Bool("server-mode") {
				return runServer(origCred, cred, conf, port)
			}

			return runDaemon(homedir, origCred, cred, conf, port)
		},
	}
)

func runDaemon(homedir, orig, cred, conf, port string) error {
	daemondir := filepath.Join(homedir, "daemon")

	if _, err := os.Stat(daemondir); os.IsNotExist(err) {
		log.Warnf("the dir doesn't exist, create a new dir: %s", daemondir)
		os.Mkdir(daemondir, 0755)
	}

	daemonCtx := &daemon.Context{
		PidFileName: filepath.Join(daemondir, "daemon.pid"),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(daemondir, "daemon.log"),
		LogFilePerm: 0640,
		WorkDir:     daemondir,
		Umask:       027,
	}

	d, err := daemonCtx.Reborn()
	if err != nil {
		return err
	}
	if d != nil {
		fmt.Println("run the daemon successfully.")
		return nil
	}
	defer daemonCtx.Release()

	// ------ start daemon ------
	ctx, cancle := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	defer cancle()

	go g.Go(func() error {
		log.Infof("start a server ... \n - aws credentials: \"%s\"\n - awscred credentials: \"%s\"\n - awscred config: \"%s\"", orig, cred, conf)

		lis, err := net.Listen("tcp", port)
		if err != nil {
			return err
		}
		s := getServer(orig, cred, conf, port)

		return s.Serve(lis)
	})

	// set the signal handler
	daemon.SetSigHandler(func(sig os.Signal) error {
		cancle()

		log.Info("remove the awscred credential")
		os.Remove(cred)

		return daemon.ErrStop
	}, os.Interrupt, syscall.SIGTERM)
	err = daemon.ServeSignals()
	if err != nil {
		log.Errorf("failed to server signal: %s", err)
	}

	return nil
}

func runServer(orig, cred, conf, port string) error {
	var (
		s *grpc.Server
	)

	// signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	ctx, cancle := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	defer cancle()

	go g.Go(func() error {
		log.Infof("start a server ... \n - aws credentials: \"%s\"\n - awscred credentials: \"%s\"\n - awscred config: \"%s\"", orig, cred, conf)

		lis, err := net.Listen("tcp", port)
		if err != nil {
			return err
		}
		s = getServer(orig, cred, conf, port)

		return s.Serve(lis)
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	cancle()

	log.Info("stop the server gracefully.")
	s.GracefulStop()

	log.Info("remove the awscred credentials.")
	os.Remove(cred)

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func getServer(orig, cred, conf, port string) *grpc.Server {
	i := server.NewInteractor(orig, cred, conf)
	i.StartWatch(context.Background())
	i.Reflect()

	s := grpc.NewServer()
	pb.RegisterAWSCredServer(s, server.NewServer(i))
	return s
}
