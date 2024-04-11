package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Clean struct {
	ctx *Context
}

func (me *Clean) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (me *Clean) Run(c *cli.Context) error {
	configfiles := me.ctx.getConfigFiles(c)

	for _, filename := range configfiles {
		err := me.RunFile(filename)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Clean) RunFile(path string) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}
	return me.RunConfig(cfg, path)
}

func (me *Clean) RunConfig(cfg *AppConfig, path string) error {

	err := me.CleanUnreferenced(cfg, path)
	if err != nil {
		log.Warnf("Clean unreferenced %s", err)
	}

	return nil
}
func (me *Clean) CleanUnreferenced(cfg *AppConfig, path string) error {
	p := NewRunParamsConfig(cfg)
	run, err := CompileRun(path, p)
	if err != nil {
		return err
	}

	configRef := make(map[string]bool, 0)
	for _, l := range run.Links {
		configRef[l.Link] = true
	}

	log.Infof("Finding files that are not referenced in config")

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	names, err := ListDir(pwd)
	for _, name := range names {
		_, found := configRef[name]
		if found == false {
			log.Infof("%s not mentioned in config", name)
		}
	}

	return nil

}
