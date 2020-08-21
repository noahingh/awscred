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
	// OffCommand set the profile disabled
	OffCommand = &cli.Command{
		Name:  "off",
		Usage: "set disabled the session token of profile to be reflected on the awscred credentials.",
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

			return off(address, profile)
		},
	}
)

func off(address, profile string) error {
	log.Debugf("set the conn: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	c := pb.NewAWSCredClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Debug("grpc call to the server.")
	_, err = c.SetOff(ctx, &pb.SetOffRequest{Profile: profile})
	if err != nil {
		return fmt.Errorf("couldn't set enabled: %s", err)
	}

	log.Printf("set the profile disabled: %s\n", profile)

	return nil
}


