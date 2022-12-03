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
	run, err := CompileRun(cfg.Symlinks, cfg.Script)
	if err != nil {
		return err
	}

	err = Mkdirs(run.HomeDir, cfg.Mkdirs)
	if err != nil {
		return err
	}

	err = DoCreateLinks(opts, run)
	if err != nil {
		return err
	}
	err = CleanLinks(opts, cfg.Clean, run.HomeDir)
	if err != nil {
		return err
	}

	return nil
}

func DoCreateLinks(opts *LinkOptions, run *Run) error {
	err := RunScripts(opts, run, "pre")
	if err != nil {
		return err
	}
	err = CreateLinks(opts, run)
	if err != nil {
		return err
	}

	err = RunScripts(opts, run, "post")
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

		// validate we are not linking to ourselves
		if li.AbsLink == li.Target {
			log.Debugf("Ignoring symlink to self %s", li.AbsLink)
			continue
		}

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

	run, err := CompileRun(cfg.Symlinks, cfg.Script)
	if err != nil {
		return err
	}

	err = DoCreateLinks(opts, run)
	if err != nil {
		return err
	}

	return nil
}

func RunScripts(opts *LinkOptions, run *Run, stype string) error {
	log.Tracef("running scripts of type %s", stype)
	for _, script := range run.Script {
		if script.Type != stype {
			log.Tracef("skip script %s type %s", script.Id, script.Type)
			continue
		}
		err := script.Validate()
		if err != nil {
			log.Warnf("%s-script %s validate: %s", script.Type, script.Id, err)
			continue
		}

		sres, err := script.Run()
		if err != nil {
			log.Warnf("%s-script %s: run: %s", script.Type, script.Id, err)
			continue
		}

		log.Tracef("%s-script %s returned %s", script.Type, script.Id, sres)
	}

	return nil
}

// go through each glob and ensure its a directory
func CleanLinks(opts *LinkOptions, dirs []string, homedir string) error {
	for _, glob_pattern := range dirs {
		if strings.HasPrefix(glob_pattern, "~") {
			glob_pattern = filepath.Join(homedir, glob_pattern[1:])
		}
		matches, err := filepath.Glob(glob_pattern)
		if err != nil {
			return err
		}
		err = CleanLinksGlob(opts, glob_pattern, matches)
		if err != nil {
			log.Warnf("Glob %s: %s", glob_pattern, err)
			continue
		}
	}

	return nil
}

func CleanLinksGlob(opts *LinkOptions, dir_pattern string, matches []string) error {
	log.Tracef("CleanLinksGlob %s", dir_pattern)
	for _, filename := range matches {
		log.Tracef("filename %s", filename)

		st, err := os.Stat(filename)
		if err != nil {
			log.Warnf("Stat %s: %s", filename, err)
			continue
		}
		if st.IsDir() == false {
			continue // we only want directories
		}

		dir := filename

		ls, err := ListDir(dir)
		if err != nil {
			log.Warnf("ListDir %s: %s", dir, err)
			continue
		}

		// go through each file and ensure it's a symlink and valid
		for _, filename := range ls {
			fullpath := filepath.Join(dir, filename)
			log.Tracef("clean symlink %s", fullpath)

			// stat the link
			st, err := os.Lstat(fullpath)
			if err != nil {
				log.Warnf("Lstat %s: %s", fullpath, err)
			}
			is_symlink := false
			if err == nil && st.Mode()&os.ModeSymlink == os.ModeSymlink {
				is_symlink = true
			}

			// stat the file (and what the link points to)
			_, err = os.Stat(fullpath)
			stat_ok := (err == nil)

			// if it's a symlink, get what it points to
			var link string
			if is_symlink {
				link, err = os.Readlink(fullpath)
				if err != nil {
					log.Warnf("Readlink %s", err)
				}
			}
			log.Tracef("points to %s", link)

			// check if the symlink points to something invalid
			dangling := false
			if is_symlink == true && link != "" {
				_, err := os.Stat(link)
				if err != nil {
					dangling = true
				}
			}
			if dangling && stat_ok == false {
				log.Tracef("dangling symlink %s is invalid, points to %s",
					fullpath, link)

				log.Infof("Remove dangling symlink %s", fullpath)
				os.Remove(fullpath)
			}

		}
	}
	return nil
}
