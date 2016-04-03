package client

import (
	"log"
	"syncfile/file-watcher"
)

func (c *SyncClient) openFSWatcher() error {
	c.w = file_watcher.NewFSWatcher()
	return c.w.Open(c.root)
}

func (c *SyncClient) runFSWatcher() {
	for {
		ch, err := c.w.Read()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Change: %v %s", ch.Op, ch.Path)
		if err := c.db.Write(ch); err != nil {
			log.Fatal(err)
		}
	}
}
