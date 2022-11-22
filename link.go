package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

type Link struct {
	ctx *Context
}

func (me *Link) Flags() []cli.Flag {
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
	}
	return nil
}

type LinkOptions struct {
	Pretend bool
}

func (me *Link) Run(c *cli.Context) error {
	opts := &LinkOptions{}
	configfiles := getConfigFiles(c)
	opts.Pretend = c.Bool("pretend")
	log.Tracef("%d files to execute", len(configfiles))

	for _, filename := range configfiles {
		err := me.RunFile(opts, filename)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Link) RunFile(opts *LinkOptions, path string) error {
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

	created := 0

	for _, li := range run.Links {
		if opts.Pretend {
			if li.NeedsCreate {
				log.Infof("create link %q %q", li.AbsLink, li.Target)
			}
			continue
		}

		if li.NeedsCreate {

			if li.DestExists {
				os.Remove(li.Target)
			}

			log.Infof("symlink %s", li.Target)
			err = os.Symlink(li.AbsLink, li.Target)
			if err != nil {
				log.Warnf("Symlink %s", err)
				continue
			}
			created++
		}
	}
	if created > 0 {
		log.Infof("created %d links", created)
	}

	return nil

}
