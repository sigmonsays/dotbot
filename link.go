package main

import (
	"os"
	"path/filepath"
	"strings"

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
		&cli.BoolFlag{
			Name:    "auto",
			Usage:   "auto mode",
			Aliases: []string{"a"},
		},
	}
	return nil
}

type LinkOptions struct {
	Pretend  bool
	AutoMode bool
}

func (me *Link) Run(c *cli.Context) error {
	opts := &LinkOptions{}
	configfiles := getConfigFiles(c)
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

	return me.RunConfig(opts, cfg)
}

func (me *Link) RunConfig(opts *LinkOptions, cfg *AppConfig) error {
	run, err := CompileRun(cfg.Symlinks)
	if err != nil {
		return err
	}

	err = Mkdirs(run.HomeDir, cfg.Mkdirs)
	if err != nil {
		return err
	}

	err = CreateLinks(opts, run)
	if err != nil {
		return err
	}
	return nil
}
func CreateLinks(opts *LinkOptions, run *Run) error {
	var (
		err     error
		created int
	)

	for _, li := range run.Links {
		if opts.Pretend {
			if li.NeedsCreate {
				log.Infof("pretend create link %q %q", li.AbsLink, li.Target)
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

func Mkdirs(homedir string, paths []string) error {
	dirMode := os.FileMode(0755)
	for _, path := range paths {
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(homedir, path[2:])
		}

		st, err := os.Stat(path)
		if err == nil && st.Mode().IsRegular() {
			log.Warnf("file exists for mkdir %s", path)
			continue
		}
		if err != nil {
			log.Infof("mkdir %s", path)
			os.MkdirAll(path, dirMode)
		}
	}
	return nil
}

func (me *Link) RunAutoMode(opts *LinkOptions) error {
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

	run, err := CompileRun(cfg.Symlinks)
	if err != nil {
		return err
	}

	err = CreateLinks(opts, run)
	if err != nil {
		return err
	}

	return nil
}
