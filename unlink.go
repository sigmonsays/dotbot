package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Unlink struct {
	ctx *Context
}

func (me *Unlink) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "config",
			Usage:   "config file",
			Aliases: []string{"c"},
		},
		&cli.BoolFlag{
			Name:    "pretend",
			Usage:   "pretend mode",
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "auto",
			Usage:   "auto mode",
			Aliases: []string{"a"},
		},
	}
	return nil
}

type UnlinkOptions struct {
	Pretend  bool
	AutoMode bool
}

func (me *Unlink) Run(c *cli.Context) error {
	opts := &UnlinkOptions{}
	configfiles := me.ctx.getConfigFiles(c)
	opts.Pretend = c.Bool("pretend")
	opts.AutoMode = c.Bool("auto")
	log.Tracef("%d files to execute", len(configfiles))

	if len(configfiles) == 0 && opts.AutoMode == false {
		log.Warnf("Nothing to do, try passing -c dotbot.yaml ")
		return nil
	}

	for _, filename := range configfiles {
		err := me.RunFile(opts, filename)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	if opts.AutoMode {
		err := me.RunAutoMode(opts)
		if err != nil {
			log.Warnf("RunAutoMode: %s", err)
		}
	}

	return nil
}

func (me *Unlink) RunFile(opts *UnlinkOptions, path string) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}

	return me.RunConfig(opts, cfg, path)
}

func (me *Unlink) RunConfig(opts *UnlinkOptions, cfg *AppConfig, path string) error {
	run, err := CompileRun(path, cfg.Symlinks, cfg.WalkDir, cfg.Script, cfg.Include)
	if err != nil {
		return err
	}

	err = Mkdirs(run.HomeDir, cfg.Mkdirs)
	if err != nil {
		return err
	}

	err = DoUnlinks(opts, run)
	if err != nil {
		return err
	}
	return nil
}
func DoUnlinks(opts *UnlinkOptions, run *Run) error {
	var (
		err     error
		removed int
	)

	for _, li := range run.Links {
		if opts.Pretend {
			if li.NeedsCreate == false {
				log.Infof("unlink %q %q", li.AbsLink, li.Target)
			}
			continue
		}

		if li.NeedsCreate == false {
			log.Tracef("Remove %s", li.Target)
			err = os.Remove(li.Target)
			if err != nil {
				log.Warnf("Remove %s: %s", li.Target, err)
				continue
			}

			removed++
		}
	}
	if removed > 0 {
		log.Infof("removed %d links", removed)
	}
	return nil
}

func (me *Unlink) RunAutoMode(opts *UnlinkOptions) error {
	cfg := GetDefaultConfig()
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filenames, err := ListDir(cwd)
	if err != nil {
		return err
	}

	// build config using current directory listing
	for _, filename := range filenames {
		if filename == ".git" {
			continue
		}
		cfg.Symlinks["~/"+filename] = filename
	}

	run, err := CompileRun("", cfg.Symlinks, cfg.WalkDir, cfg.Script, cfg.Include)
	if err != nil {
		return err
	}

	err = DoUnlinks(opts, run)
	if err != nil {
		return err
	}

	return nil
}
