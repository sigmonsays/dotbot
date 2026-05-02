package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Clean struct {
	ctx *Context
}

func (me *Clean) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "target-dir",
			Usage:   "override target directory",
			Aliases: []string{"t"},
		},
	}
}

func (me *Clean) Run(c *cli.Context) error {
	configfiles := me.ctx.getConfigFiles(c)
	opts := &LinkOptions{}
	opts.TargetDir = c.String("target-dir")

	for _, filename := range configfiles {
		err := me.RunFile(filename, opts)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Clean) RunFile(path string, opts *LinkOptions) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}
	return me.RunConfig(cfg, opts, path)
}

func (me *Clean) RunConfig(cfg *AppConfig, opts *LinkOptions, path string) error {
	err := me.CleanUnreferenced(cfg, opts, path)
	if err != nil {
		log.Warnf("Clean unreferenced %s", err)
	}

	return nil
}
func (me *Clean) CleanUnreferenced(cfg *AppConfig, opts *LinkOptions, path string) error {
	p := NewRunParamsConfig(cfg)
	run, err := CompileRun(path, p, opts)
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
