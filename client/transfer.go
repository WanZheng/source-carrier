package client

import (
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"syncfile/data"
)

func (c *SyncClient) openRpc() error {
	cli, err := rpc.DialHTTP("tcp", c.servAddr)
	if err != nil {
		return err
	}
	c.c = cli
	return nil
}

func (c *SyncClient) sync() error {
	var seq int64
	if err := c.c.Call("Sync.Seq", 0, &seq); err != nil {
		return err
	}

	for {
		l, err := c.db.List(seq)
		if err != nil {
			return err
		}
		if len(l) <= 0 {
			break
		}

		for _, ch := range l {
			log.Printf("Syncing: #%d %v %s", ch.Seq, ch.Op, ch.Path)
			if err := c.fillData(&ch); err != nil {
				log.Printf("Failed to read #%d %v %s: %v", ch.Seq, ch.Op, ch.Path, err)
				return err
			}
			if err := c.c.Call("Sync.Put", ch, &seq); err != nil {
				log.Printf("Failed to sync #%d %v %s: %v", ch.Seq, ch.Op, ch.Path, err)
				return err
			}
			seq = ch.Seq
		}
	}

	return c.db.Synced(seq)
}

func (c *SyncClient) fillData(ch *data.Change) error {
	if ch.Op != data.UPDATE {
		return nil
	}
	f, err := os.Open(filepath.Join(c.root, ch.Path))
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	ch.Data = buf
	return nil
}
