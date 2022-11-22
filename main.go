package main

import (
	"os"

	"github.com/urfave/cli/v2"

	gologging "github.com/sigmonsays/go-logging"
)

func main() {

	ctx := &Context{}
	link := &Link{ctx}
	status := &Status{ctx}
	cleanup := &Cleanup{ctx}

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
		Name:    "status",
		Aliases: []string{"s"},
		Usage:   "print status as json",
		Action:  status.Run,
		Flags:   status.Flags(),
	})
	ctx.addCommand(&cli.Command{
		Name:   "cleanup",
		Usage:  "show unreferenced files",
		Action: cleanup.Run,
		Flags:  cleanup.Flags(),
	})
	app.Run(os.Args)
}
