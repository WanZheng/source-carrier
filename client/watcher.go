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

func (c *SyncClient) watchRouting() {
	for {
		ch, err := c.w.Read()
		if err != nil {
			log.Fatal("failed to read change:", ch, err)
		}

		if c.ignore.IsIgnored(ch.Path) {
			log.Printf("Ignore change: %v %s", ch.Op, ch.Path)
			continue
		}

		c.queue <- ch

		if filepath.Base(ch.Path) == ".gitignore" {
			log.Print("update ignore patterns")
			if err := c.ignore.Update(); err != nil {
				log.Fatal("failed to update ignore:", err)
			}
		}
	}
}

func (c *SyncClient) dbRouting() {
	for {
		ch := <-c.queue
		if err := c.db.Write(ch); err != nil {
			log.Fatal("failed to write db:", err)
		}
	}
}
