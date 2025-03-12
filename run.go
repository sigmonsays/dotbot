package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func NewRun() *Run {
	run := &Run{}
	usr, _ := user.Current()
	homedir := usr.HomeDir
	log.Tracef("homedir is %s", homedir)
	run.HomeDir = homedir
	run.Script = make([]*Script, 0)
	return run
}

type Run struct {
	HomeDir string
	Links   []*LinkInfo
	Script  []*Script
	Mkdir   []string
}
type LinkInfo struct {
	OrigTarget  string
	Target      string
	Link        string
	FileType    string
	AbsLink     string
	DestExists  bool
	PointsTo    string
	IsValid     bool
	NeedsCreate bool
}

func NewRunParamsConfig(c *AppConfig) *RunParams {
	ret := &RunParams{}
	ret.Symlinks = c.Symlinks
	ret.Walkdir = c.WalkDir
	ret.Script = c.Script
	ret.Mkdir = c.Mkdirs
	ret.Include = c.Include
	return ret
}

type RunParams struct {
	Symlinks map[string]string
	Walkdir  map[string]string
	Script   []*Script
	Mkdir    []string
	Include  []string
}

func CompileRun(path string, p *RunParams) (*Run, error) {
	run := NewRun()
	err := CompileRunWithRun(path, run, p)
	if err != nil {
		return nil, err
	}
	return run, nil
}

func CompileRunWithRun(path string, run *Run, p *RunParams) error {

	// change dir and then change back before returning
	dir, _ := os.Getwd()
	ChdirToFile(path)
	defer func() {
		log.Tracef("chdir to old dir %s", dir)
		os.Chdir(dir)
	}()

	//  mkdirs
	for _, m := range p.Mkdir {
		run.Mkdir = append(run.Mkdir, m)
	}

	//  scripts
	for _, s := range p.Script {
		if s.Disabled {
			continue
		}
		run.Script = append(run.Script, s)
	}

	err := CompileRunSymlinks(run, p.Symlinks)
	if err != nil {
		log.Warnf("CompileRunSymlinks %s", err)
		return err
	}

	err = CompileRunWalkDir(run, p.Walkdir)
	if err != nil {
		log.Warnf("CompileRunWalkDir %s", err)
		return err
	}

	includes, err := GetIncludes(path, p.Include)
	if err != nil {
		log.Warnf("GetIncludes %s", err)
		return err
	}
	for _, include := range includes {
		cfg2 := GetDefaultConfig()

		st, err := os.Stat(include)
		if err != nil {
			log.Warnf("Stat %s: %s", include, err)
			continue
		}
		if st.IsDir() {
			include = include + "/dotbot.yaml"
		}

		err = cfg2.LoadYaml(include)
		if err != nil {
			log.Errorf("Include %s: %s", include, err)
			continue
		}
		p2 := NewRunParamsConfig(cfg2)
		err = CompileRunWithRun(include, run, p2)
		if err != nil {
			log.Errorf("Include CompileRunWithRun %s: %s", include, err)
			continue
		}
	}

	return nil
}

func CompileRunSymlinks(run *Run, symlinks map[string]string) error {

	// symlinks
	for target, link := range symlinks {

		li := &LinkInfo{}
		run.Links = append(run.Links, li)
		li.OrigTarget = target

		// resolve target tilde prefix
		if strings.HasPrefix(target, "~/") {
			target = filepath.Join(run.HomeDir, target[2:])
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

		target_stat, err := os.Stat(link)
		li.FileType = "unknown"
		if err == nil {
			if target_stat.Mode().IsDir() {
				li.FileType = "dir"
			} else if target_stat.Mode().IsRegular() {
				li.FileType = "file"
			}
		}

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

	return nil
}

func CompileRunWalkDir(run *Run, walkdir map[string]string) error {

	// compile a walkdir into a CompileRunSymlinks
	symlinks := make(map[string]string, 0)

	for targetdir, srcdir := range walkdir {
		filenames, err := ListDir(srcdir)
		if err != nil {
			log.Warnf("ListDir %s: %s", srcdir, err)
			continue
		}
		for _, filename := range filenames {
			log.Tracef("ProcessWalkDir %s", filename)
			src := filepath.Join(srcdir, filename)
			target := filepath.Join(targetdir, filename)
			log.Tracef("symlink %s to %s", src, target)
			symlinks[target] = src
		}
	}
	return CompileRunSymlinks(run, symlinks)
}

func GetIncludes(path string, includes []string) ([]string, error) {
	ret := make([]string, 0)

	if len(includes) == 0 {
		log.Tracef("no includes to process in %s", path)
		return nil, nil
	}
	for _, include := range includes {
		matches, err := filepath.Glob(include)
		if err != nil {
			log.Warnf("Glob %s: %", include, err)
			continue
		}
		ret = append(ret, matches...)
	}
	return ret, nil
}

// change directory to a filename
func ChdirToFile(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Warnf("abs %s: %s", path, err)
		return err
	}
	dir := filepath.Dir(abs)
	log.Tracef("chdir %s", dir)
	err = os.Chdir(dir)
	if err != nil {
		log.Errorf("Chdir %s: %s", dir, err)
		return err
	}

	return nil
}
