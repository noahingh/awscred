package subcmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/hanjunlee/awscred/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	// OnCommand set the profile enabled
	OnCommand = &cli.Command{
		Name:  "on",
		Usage: "set the session token of profile to be reflected on the awscred credentials.",
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
				profile string
			)
			if c.Bool("debug") {
				setDebugMode()
			}

			address = "localhost:" + strconv.Itoa(c.Int("port"))
			profile = c.Args().Get(0); 

			return on(address, profile)
		},
	}
)

func on(address, profile string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	c := pb.NewAWSCredClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SetOn(ctx, &pb.SetOnRequest{Profile: profile})
	if err != nil {
		return fmt.Errorf("couldn't set enabled: %s", err)
	}

	log.Infof("set \"%s\" enabled", profile)

	return nil
}
