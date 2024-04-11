package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Context struct {
	app     *cli.App
	configs []string
}

func (me *Context) addCommand(cmd *cli.Command) {
	me.app.Commands = append(me.app.Commands, cmd)
}

func (me *Context) getConfigFiles(c *cli.Context) []string {
	defaultFile := "dotbot.yaml"
	ret := me.configs

	if len(ret) == 0 {
		st, err := os.Stat(defaultFile)
		if err == nil && st.Mode().IsRegular() {
			ret = []string{defaultFile}
			log.Tracef("loading default %s", defaultFile)
		}
	}

	// add any additional files
	for _, f := range c.Args().Slice() {
		st, err := os.Stat(f)
		if err != nil {
			log.Errorf("Stat %s: %s", f, err)
			continue
		}
		if st.IsDir() {
			log.Errorf("Stat %s: is a directory", f)
			continue
		}
		ret = append(ret, f)
	}

	return ret
}

func ListDir(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(10000)
	if err != nil {
		return nil, err
	}
	return names, nil
}
