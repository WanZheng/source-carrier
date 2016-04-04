package client

import (
	"log"
	"path/filepath"
	"syncfile/watcher"
)

func (c *SyncClient) openFSWatcher() error {
	c.w = watcher.NewFSWatcher()
	return c.w.Open(c.root)
}

func (c *SyncClient) runFSWatcher() {
	for {
		ch, err := c.w.Read()
		if err != nil {
			log.Fatal(err)
		}

		if c.ignore.IsIgnored(ch.Path) {
			log.Printf("Ignore change: %v %s", ch.Op, ch.Path)
			continue
		}

		log.Printf("Change: %v %s", ch.Op, ch.Path)
		if err := c.db.Write(ch); err != nil {
			log.Fatal(err)
		}

		if filepath.Base(ch.Path) == ".gitignore" {
			log.Print("update ignore patterns")
			if err := c.updateIgnoreTable(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
