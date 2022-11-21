package main

import (
	"os"
	"os/user"
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
	}
	return nil
}
func (me *Link) Run(c *cli.Context) error {

	configfiles := c.StringSlice("config")
	log.Tracef("%d files to execute", len(configfiles))

	for _, filename := range configfiles {
		err := me.RunFile(filename)
		if err != nil {
			log.Warnf("RunFile %s: %s", filename, err)
		}
	}

	return nil
}

func (me *Link) RunFile(path string) error {
	log.Tracef("runfile %s", path)
	cfg := GetDefaultConfig()
	err := cfg.LoadYaml(path)
	if err != nil {
		return err
	}
	if log.IsTrace() {
		cfg.PrintConfig()
	}

	usr, _ := user.Current()
	homedir := usr.HomeDir
	log.Tracef("homedir is %s", homedir)

	for target, link := range cfg.Symlinks {

		// resolve target tilde prefix
		if strings.HasPrefix(target, "~/") {
			target = filepath.Join(homedir, target[2:])
		}

		abslink, err := filepath.Abs(link)
		if err != nil {
			log.Warnf("Abs %s: %s", link, err)
			continue
		}
		log.Tracef("symlink %s to %s", abslink, target)

		// stat dest and check if destination exists and is a symlink
		dest, err := os.Lstat(target)
		dest_exists := (err == nil)
		log.Tracef("target exists %v", dest_exists)
		if err == nil {
			ftype := dest.Mode().Type()
			log.Tracef("target filetype %d", ftype)
		}
		is_symlink := false
		pointsto := ""
		if dest.Mode()&os.ModeSymlink != 0 {
			is_symlink = true
			link, err := os.Readlink(target)
			if err == nil {
				pointsto = link
			} else {
				log.Warnf("Readlink %s", err)
			}
		}

		if dest_exists && is_symlink {
			log.Tracef("%s already created", abslink)
			log.Tracef("points to %s", pointsto)
			continue
		}

		err = os.Symlink(abslink, target)
		if err != nil {
			log.Warnf("Symlink %s", err)
		}

	}
	return nil
}
