package data

import (
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	path := "/tmp/test.sqlite3"
	os.Remove(path)
	db, err := OpenSqlDB(path)
	if err != nil {
		t.Error(err)
	}

	// empty
	seq, err := db.Seq()
	if err != nil {
		t.Error("failed to query seq: ", err)
	}
	if seq != 0 {
		t.Error("seq != 0 for an empty db")
	}

	l, err := db.List(0)
	if err != nil {
		t.Error(err)
	}
	if len(l) != 0 {
		t.Error("should be empty")
	}

	// 1 change
	ch := Change{
		Op:    UPDATE,
		Path:  "a.file",
		Size:  1,
		Mtime: 2}
	if err := db.Write(ch); err != nil {
		t.Error(err)
	}

	seq, err = db.Seq()
	if seq != 1 {
		t.Error("seq != 1")
	}

	l, err = db.List(0)
	if err != nil {
		t.Error(err)
	}
	if len(l) != 1 {
		t.Error("should be empty")
	}

	// remove 1
	ch.Op = DEL
	db.Write(ch)
	seq, err = db.Seq()
	if err != nil {
		t.Error(err)
	}
	if seq != 2 {
		t.Error("seq != 2, seq=", seq)
	}
	l, err = db.List(0)
	if err != nil {
		t.Error(err)
	}
	if len(l) != 1 {
		t.Error("should be empty")
	}
	ch0 := l[0]
	if ch0.Op != ch.Op || ch0.Path != ch.Path {
		t.Error("incorrect change: ", ch0)
	}

	// synced
	if err := db.Synced(2); err != nil {
		t.Error(err)
	}
	l, _ = db.List(0)
	if len(l) != 0 {
		t.Error("should be empty: ", l)
	}

	defer db.Close()
}
