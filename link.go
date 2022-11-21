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
	configfiles := c.StringSlice("config")
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

type Run struct {
	HomeDir string
	Links   []*LinkInfo
}
type LinkInfo struct {
	Target      string
	Link        string
	AbsLink     string
	DestExists  bool
	PointsTo    string
	IsValid     bool
	NeedsCreate bool
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

	run, err := me.CompileRun(cfg.Symlinks)
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

func (me *Link) CompileRun(symlinks map[string]string) (*Run, error) {
	ret := &Run{}

	usr, _ := user.Current()
	homedir := usr.HomeDir
	log.Tracef("homedir is %s", homedir)
	ret.HomeDir = homedir

	for target, link := range symlinks {

		li := &LinkInfo{}
		ret.Links = append(ret.Links, li)

		// resolve target tilde prefix
		if strings.HasPrefix(target, "~/") {
			target = filepath.Join(homedir, target[2:])
		}

		abslink, err := filepath.Abs(link)
		if err != nil {
			log.Warnf("Abs %s: %s", link, err)
			continue
		}
		log.Tracef("---")
		log.Tracef("symlink %s to %s", abslink, target)
		log.Tracef("abslink %s", abslink)
		log.Tracef("target %s", target)

		li.Target = target
		li.Link = link
		li.AbsLink = abslink

		// stat dest and check if destination exists and is a symlink
		dest, err := os.Lstat(target)
		dest_exists := (err == nil)
		li.DestExists = dest_exists
		log.Tracef("target exists %v", dest_exists)
		is_symlink := false
		pointsto := ""
		if dest_exists && dest.Mode()&os.ModeSymlink != 0 {
			is_symlink = true
			link, err := os.Readlink(target)
			if err == nil {
				pointsto = link
			} else {
				log.Warnf("Readlink %s", err)
			}
		}
		li.PointsTo = pointsto

		is_valid := true
		if is_symlink {
			// evaluate what it points to
			if pointsto != abslink {
				is_valid = false
				log.Tracef("invalid link points to %s but should be %s", abslink, pointsto)
			}
		}
		li.IsValid = is_valid

		if dest_exists && is_symlink && is_valid {
			log.Tracef("%s already created", abslink)
			log.Tracef("points to %s", pointsto)
		} else {
			li.NeedsCreate = true
		}
		log.Tracef("needs_create:%v", li.NeedsCreate)

	}

	return ret, nil
}
