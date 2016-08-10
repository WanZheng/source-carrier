package client

import (
	"log"
	"os"
	"path/filepath"
	"syncfile/data"
)

func (c *SyncClient) reScan() error {
	l, err := c.db.List(0)
	if err != nil {
		return err
	}

	// 1. read records from DB, and build a map
	m := make(map[string]data.Change, 0)
	for _, ch := range l {
		m[ch.Path] = ch
	}

	// 2. traval file tree, mark exist files and find out updated files
	updated := make([]data.Change, 0, 0)
	err = filepath.Walk(c.root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path, err = filepath.Rel(c.root, path)
			if err != nil {
				return err
			}
			if c.ignore.IsIgnored(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if info.IsDir() {
				return nil
			}

			ch, ok := m[path]
			if ok && path != ch.Path {
				log.Fatalf("path not match: %v, %v", path, ch)
			}

			if ok {
				ch.Seq = -1
				m[path] = ch
			}
			if !ok || ch.Op != data.UPDATE || info.Size() != ch.Size || info.ModTime().Unix() != ch.Mtime {
				if ok {
					log.Printf("file changed: %v %v", info, ch)
				}
				updated = append(updated, data.Change{Op: data.UPDATE, Path: path, Size: info.Size(), Mtime: info.ModTime().Unix()})
			}
			return nil
		})
	if err != nil {
		return err
	}

	// 3. traval map, and find deleted file
	for _, ch := range m {
		if ch.Seq == -1 || ch.Op == data.DEL {
			continue
		}
		ch.Op = data.DEL
		ch.Seq = -1
		updated = append(updated, ch)
	}

	// 4. update db
	for _, ch := range updated {
		if err := c.db.Write(ch); err != nil {
			return err
		}
	}

	return nil
}
