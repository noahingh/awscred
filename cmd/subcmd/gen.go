package subcmd

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
	// GenCommand generate a new session token.
	GenCommand = &cli.Command{
		Name:  "gen",
		Usage: "generate a new session token and cache the token in the config file. e.g) awscred gen --code CODE PROFILE",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "code",
				Aliases:  []string{"c"},
				Value:    "",
				Usage:    "The  value  provided  by  the MFA device.",
				Required: true,
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
				address string
				profile string
				code    string
			)
			if c.Bool("debug") {
				setDebugMode()
			}

			address = "localhost:" + strconv.Itoa(c.Int("port"))
			profile = c.Args().Get(0)
			code = c.String("code")

			return gen(address, profile, code)
		},
	}
)

func gen(address, profile, code string) error {
	log.Debugf("set the conn: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	c := pb.NewAWSCredClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Debug("grpc call to the server.")
	_, err = c.SetGenerate(ctx, &pb.SetGenerateRequest{
		Profile: profile,
		Token:   code,
	})
	if err != nil {
		return fmt.Errorf("couldn't generate a session token: %s", err)
	}

	fmt.Printf("generate a session token.\n")
	return nil
}
