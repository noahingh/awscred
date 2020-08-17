package subcmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	pb "github.com/hanjunlee/awscred/api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	// InfoCommand generate a new session token.
	InfoCommand = &cli.Command{
		Name:  "info",
		Usage: "show the information for each profile.",
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
			return info(address)
		},
	}
)

func info(address string) error {
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
	r, err := c.GetProfileList(ctx, &pb.GetProfileListRequest{})
	if err != nil {
		return fmt.Errorf("couldn't generate a session token: %s", err)
	}

	printProfileList(r.Profiles)

	return nil
}

func printProfileList(pl []*pb.Profile) {
	const (
		padding = 4
	)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tON\tSERIAL\tDURATION\tEXPIRED\t")
	for _, p := range pl {
		var (
			name     = p.Name
			on       = strconv.FormatBool(p.On)
			serial   = p.Serial
			duration = strconv.Itoa(int(p.Duration))
			expired  = p.Expired
		)

		t := time.Time{}
		if expired != t.Format(time.RFC3339) {
			et, err := time.Parse(time.RFC3339, expired)
			if err == nil {
				// append the time left.
				left := et.Sub(time.Now())
				expired = strings.Join([]string{
					expired,
					" (",
					strconv.FormatFloat(left.Hours(), 'f', 1, 64),
					"h)",
				}, "")
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", name, on, serial, duration, expired)
	}
	w.Flush()
}
