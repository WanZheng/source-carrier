package watcher

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syncfile/data"
)

type FileWatcher interface {
	Open(path string) error
	Read() (data.Change, error)
	Close()
}

type FSWatcher struct {
	root string
	ch   chan string
	cmd  *exec.Cmd
}

func NewFSWatcher() *FSWatcher {
	w := FSWatcher{}
	w.ch = make(chan string, 16)
	return &w
}

func (w *FSWatcher) Open(path string) error {
	root, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	root, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	w.root = root

	w.cmd = exec.Command("fswatch", "-e", "\\.git\\/", w.root)
	go func() {
		out, err := w.cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
			return
		}
		if err = w.cmd.Start(); err != nil {
			log.Fatal(err)
			return
		}
		r := bufio.NewReader(out)
		for {
			l, err := r.ReadString('\n')
			if err != nil {
				break
			}
			l = strings.TrimRight(l, "\n\r")
			w.ch <- l
		}
	}()
	return nil
}

func (w *FSWatcher) Close() {
	w.cmd.Process.Kill()
	w.cmd.Wait()
	close(w.ch)
}

func (w *FSWatcher) Read() (data.Change, error) {
	for {
		f, ok := <-w.ch
		if !ok {
			return data.Change{}, nil
		}
		rel, err := filepath.Rel(w.root, f)
		if err != nil {
			return data.Change{}, err
		}
		info, err := os.Stat(f)
		if err != nil {
			if os.IsNotExist(err) {
				return data.Change{Op: data.DEL, Path: rel}, nil
			}
			return data.Change{}, err
		}
		if info.IsDir() {
			continue
		}
		return data.Change{Op: data.UPDATE,
			Path:  rel,
			Size:  info.Size(),
			Mtime: info.ModTime().Unix(),
		}, nil
	}
}
