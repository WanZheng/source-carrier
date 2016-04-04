package server

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syncfile/data"
)

type SyncServer struct {
	root string
	db   data.DB
}

func NewSyncServer(root string, db data.DB) *SyncServer {
	return &SyncServer{root, db}
}

func (s *SyncServer) Seq(seq int64, reply *int64) error {
	latest, err := s.db.Seq()
	*reply = latest
	return err
}

func (s *SyncServer) Put(ch data.Change, reply *int64) error {
	if err := s.writeFile(ch); err != nil {
		return err
	}
	if err := s.db.Write(ch); err != nil {
		return err
	}
	*reply = ch.Seq
	return nil
}

func (s *SyncServer) writeFile(ch data.Change) error {
	path := filepath.Join(s.root, ch.Path)
	if ch.Op == data.DEL {
		err := removeAll(path)
		if err == nil || os.IsNotExist(err) {
			return nil
		}
		return err
	}

	os.MkdirAll(filepath.Dir(path), 0755)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(ch.Data)
	return err
}

func removeAll(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return removeDir(path)
	}
	return os.Remove(path)
}

func removeDir(path string) error {
	l, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, info := range l {
		child := filepath.Join(path, info.Name())
		if info.IsDir() {
			if err := removeDir(child); err != nil {
				return err
			}
		} else {
			if err := os.Remove(child); err != nil {
				return err
			}
		}
	}
	return nil
}
