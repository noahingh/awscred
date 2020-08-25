package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/hanjunlee/awscred/api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	// SetCommand set the configuration of profile..
	SetCommand = &cli.Command{
		Name:  "set",
		Usage: "set the configuration which is related with the session token generation. e.g) awscred set --serial SERIAL PROFILE",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "on",
				Value: false,
				Usage: "set enalbed after setting.",
			},
			&cli.StringFlag{
				Name:     "serial",
				Aliases:  []string{"s"},
				Value:    "",
				Usage:    "the identification number of the MFA device that is associated with the IAM user.",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "duration",
				Aliases: []string{"c"},
				Value:   43200,
				Usage:   "the  duration, in seconds, that the credentials should remain valid.",
			},
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
				address  string
				profile  string
				serial   string
				duration int64
			)
			if c.Bool("debug") {
				setDebugMode()
			}

			address = "localhost:" + strconv.Itoa(c.Int("port"))
			profile = c.Args().Get(0)
			serial = c.String("serial")
			duration = c.Int64("duration")

			if err := set(address, profile, serial, duration); err != nil {
				return err
			}

			if c.Bool("on") {
				return on(address, profile)
			}
			return nil
		},
	}
)

func set(address, profile, serial string, duration int64) error {
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
	_, err = c.SetConfig(ctx, &pb.SetConfigRequest{
		Profile:  profile,
		Serial:   serial,
		Duration: duration,
	})
	if err != nil {
		return fmt.Errorf("couldn't configure: %s", err)
	}

	fmt.Printf("set the config [serial: \"%s\", duration: \"%d\"]: %s \n", serial, duration, profile)
	return nil
}
