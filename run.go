package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func NewRun() *Run {
	ret := &Run{}
	usr, _ := user.Current()
	homedir := usr.HomeDir
	log.Tracef("homedir is %s", homedir)
	ret.HomeDir = homedir
	return ret
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

func CompileRun(symlinks map[string]string) (*Run, error) {
	ret := NewRun()

	for target, link := range symlinks {

		li := &LinkInfo{}
		ret.Links = append(ret.Links, li)

		// resolve target tilde prefix
		if strings.HasPrefix(target, "~/") {
			target = filepath.Join(ret.HomeDir, target[2:])
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
