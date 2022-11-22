package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func getConfigFiles(c *cli.Context) []string {
	defaultFile := "dotbot.yaml"
	ret := c.StringSlice("config")

	if len(ret) == 0 {
		st, err := os.Stat(defaultFile)
		if err == nil && st.Mode().IsRegular() {
			ret = []string{defaultFile}
			log.Tracef("loading default %s", defaultFile)
		}
	}
	return ret
}

type Context struct {
	app *cli.App
}

func (me *Context) addCommand(cmd *cli.Command) {
	me.app.Commands = append(me.app.Commands, cmd)
}
