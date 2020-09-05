package cmd

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	pb "github.com/hanjunlee/awscred/api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	// ExportCommand set the profile enabled
	ExportCommand = &cli.Command{
		Name:  "export",
		Usage: "return the shell command to export AWS environment variables. e.g) awscred export [PROFILE]",
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
			profile = c.Args().Get(0)

			return export(address, profile)
		},
	}
)

func export(address, profile string) error {
	if profile == "" {
		exportCredentialsFile(address)
		return nil
	}

	exportCredentialsProfile(address, profile)
	return nil
}

func exportCredentialsFile(address string) error {
	log.Debugf("set the conn: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	c := pb.NewAWSCredClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	r, err := c.GetCredentialsFile(ctx, &pb.GetCredentialsFileRequest{})
	if err != nil {
		return err
	}

	if runtime.GOOS == "window" {
		fmt.Printf("set AWS_SHARED_CREDENTIALS_FILE=%s", r.Path)
	} else {
		fmt.Printf("export AWS_SHARED_CREDENTIALS_FILE=%s", r.Path)
	}
	return nil
}

func exportCredentialsProfile(address, profile string) error {
	log.Debugf("set the conn: %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err)
	}

	c := pb.NewAWSCredClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	r, err := c.GetCredentialsProfile(ctx, &pb.GetCredentialsProfileRequest{Profile: profile})
	if err != nil {
		return err
	}

	if runtime.GOOS == "window" {
		fmt.Printf("set AWS_ACCESS_KEY_ID=%s ; set AWS_SECRET_ACCESS_KEY=%s ; set AWS_SESSION_TOKEN=%s", r.AccessKeyID, r.SecretAccessKey, r.SessionToken)
	} else {
		fmt.Printf("export AWS_ACCESS_KEY_ID=%s ; export AWS_SECRET_ACCESS_KEY=%s ; export AWS_SESSION_TOKEN=%s", r.AccessKeyID, r.SecretAccessKey, r.SessionToken)
	}
	return nil
}
