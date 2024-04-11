package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	gologging "github.com/sigmonsays/go-logging"
)

// These variables are populated by goreleaser when the binary is built.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {

	ctx := &Context{}
	link := &Link{ctx}
	unlink := &Unlink{ctx}
	status := &Status{ctx}
	clean := &Clean{ctx}

	app := &cli.App{
		Name:  "dotbot",
		Usage: "manage dot files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"l"},
				Value:   "info",
				Usage:   "set log level",
			},
			&cli.StringSliceFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "config file",
			},
		},
		Before: func(c *cli.Context) error {
			loglevel := c.String("loglevel")
			if loglevel != "" {
				gologging.SetLogLevel(loglevel)
			}

			ctx.configs = c.StringSlice("config")

			return nil
		},
	}
	ctx.app = app
	ctx.addCommand(&cli.Command{
		Name:    "link",
		Aliases: []string{"l"},
		Usage:   "create symlinks",
		Action:  link.Run,
		Flags:   link.Flags(),
	})
	ctx.addCommand(&cli.Command{
		Name:    "unlink",
		Aliases: []string{"u"},
		Usage:   "remove symlinks",
		Action:  unlink.Run,
		Flags:   unlink.Flags(),
	})
	ctx.addCommand(&cli.Command{
		Name:    "status",
		Aliases: []string{"s"},
		Usage:   "print status table",
		Action:  status.Run,
		Flags:   status.Flags(),
	})
	ctx.addCommand(&cli.Command{
		Name:   "clean",
		Usage:  "show unreferenced files",
		Action: clean.Run,
		Flags:  clean.Flags(),
	})
	ctx.addCommand(&cli.Command{
		Name:  "version",
		Usage: "print version",
		Action: func(c *cli.Context) error {
			fmt.Printf("version %s\n", version)
			fmt.Printf("commit  %s\n", commit)
			fmt.Printf("date    %s\n", date)
			return nil
		},
	})
	app.Run(os.Args)
}
