package file_watcher

import (
	"log"
	"testing"
)

func TestFSWatcher(t *testing.T) {
	var w FileWatcher
	w = NewFSWatcher()
	if err := w.Open("/tmp/a"); err != nil {
		t.Error("Failed to open:", err)
	}
	for i := 0; i < 10; i++ {
		ch, err := w.Read()
		log.Printf("ch=%v, err=%v", ch, err)
		if err != nil {
			t.Error("Failed to read change:", err)
		}
		t.Log("change: ", *ch)
	}
	w.Close()
}
