package client

import (
	"log"
	"net/rpc"
	"path/filepath"
	"syncfile/data"
	"syncfile/gitignore"
	"syncfile/watcher"
)

type SyncClient struct {
	root     string
	servAddr string
	port     int

	queue  chan data.Change
	db     data.DB
	c      *rpc.Client
	w      watcher.FileWatcher
	ignore *gitignore.Gitignore
}

func NewSyncClient(root, servAddr string, port int) *SyncClient {
	c := SyncClient{
		root:     root,
		port:     port,
		servAddr: servAddr,
		ignore:   gitignore.NewGitignore(root),
	}
	c.queue = make(chan data.Change, 1000)
	return &c
}

func (c *SyncClient) Close() {
	c.db.Close()
	c.w.Close()
}

func (c *SyncClient) Open() error {
	if err := c.openDB(); err != nil {
		return err
	}

	if err := c.openRpc(); err != nil {
		log.Print("[Warning] failed to connect to server: ", err)
	}

	if err := c.openFSWatcher(); err != nil {
		return err
	}

	return nil
}

func (c *SyncClient) openDB() error {
	dbPath := filepath.Join(c.root, "../client.sqlite3")
	db, err := data.OpenSqlDB(dbPath)
	c.db = db
	return err
}

func (c *SyncClient) updateIgnoreTable() error {
	return c.ignore.Update()
}

func (c *SyncClient) Run() error {
	if err := c.updateIgnoreTable(); err != nil {
		return err
	}

	go c.watchRouting()

	if err := c.reScan(); err != nil {
		return err
	}

	go c.dbRouting()

	return c.runHttpServer()
}
