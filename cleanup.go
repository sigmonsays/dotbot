package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Cleanup struct {
	ctx *Context
}

func (me *Cleanup) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (me *Cleanup) Run(c *cli.Context) error {
	configfiles := me.ctx.getConfigFiles(c)

	for _, filename := range configfiles {
		err := me.RunFile(filename)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Cleanup) RunFile(path string) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}
	return me.RunConfig(cfg)
}

func (me *Cleanup) RunConfig(cfg *AppConfig) error {

	run, err := CompileRun(cfg.Symlinks)
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

func ListDir(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(1000)
	if err != nil {
		return nil, err
	}
	return names, nil
}
