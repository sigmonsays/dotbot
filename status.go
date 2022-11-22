package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

type Status struct {
	ctx *Context
}

func (me *Status) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (me *Status) Run(c *cli.Context) error {
	configfiles := getConfigFiles(c)

	for _, filename := range configfiles {
		err := me.RunFile(filename)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Status) RunFile(path string) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}

	run, err := CompileRun(cfg.Symlinks)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(run)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", buf)

	return nil

}
