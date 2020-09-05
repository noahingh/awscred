package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	hc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	// PingCommand send the ping to server
	PingCommand = &cli.Command{
		Name:  "ping",
		Usage: "send the ping to daemon",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   defaultPort,
				Usage:   "the port of server.",
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
				address string
			)
			if c.Bool("debug") {
				setDebugMode()
			}

			address = "localhost:" + strconv.Itoa(c.Int("port"))

			return ping(address)
		},
	}
)

func ping(address string) error {
	const (
		okMessage = "pong"
		failedMessage = "failed"
	)

	log.Debugf("set the conn: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	c := hc.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Check(ctx, &hc.HealthCheckRequest{})
	if err != nil || r.Status != hc.HealthCheckResponse_SERVING {
		fmt.Print(failedMessage)
		return nil
	}

	fmt.Print(okMessage)
	return nil
}
