package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	gologging "github.com/sigmonsays/go-logging"
)

func main() {
	fmt.Printf("dotbot")

	ctx := &Context{}
	link := &Link{ctx}
	app := &cli.App{
		Name:  "dotbot",
		Usage: "manage dot files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"l"},
				Value:   "",
				Usage:   "set log level",
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
		Usage:   "create links",
		Action:  link.Run,
		Flags:   link.Flags(),
	})
	app.Run(os.Args)
}

type Context struct {
	app *cli.App
}

func (me *Context) addCommand(cmd *cli.Command) {
	me.app.Commands = append(me.app.Commands, cmd)
}
