package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Sops struct {
	Directory string
	Decrypt   map[string]string
}

func CompileRunSops(run *Run, sops []*Sops) error {

	// symlinks
	for _, sop := range sops {

		for link, encfile := range sop.Decrypt {

		// decrypt file

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

	}

	return nil
}
