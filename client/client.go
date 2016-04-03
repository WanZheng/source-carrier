package client

import (
	"net/rpc"
	"path/filepath"
	"syncfile/data"
	"syncfile/file-watcher"
)

type SyncClient struct {
	root     string
	servAddr string
	port     int
	db       data.DB
	c        *rpc.Client
	w        file_watcher.FileWatcher
}

func NewSyncClient(root, servAddr string, port int) *SyncClient {
	return &SyncClient{
		root:     root,
		port:     port,
		servAddr: servAddr,
	}
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
		return err
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

func (c *SyncClient) Run() error {
	go c.runFSWatcher()
	return c.runHttpServer()
}
