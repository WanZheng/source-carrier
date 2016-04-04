package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	TABLE_NAME = "changes"
	COL_SEQ    = "seq"
	COL_OP     = "op"
	COL_PATH   = "path"
	COL_SIZE   = "size"
	COL_MTIME  = "mtime"
)

type SqlDB struct {
	db *sql.DB
}

func OpenSqlDB(path string) (*SqlDB, error) {
	db := SqlDB{}
	err := db.OpenOrCreate(path)
	return &db, err
}

func (db *SqlDB) OpenOrCreate(path string) error {
	newFile := false
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		newFile = true
	} else if err != nil {
		return err
	}

	h, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	db.db = h
	if newFile {
		return db.CreateTable()
	}
	return nil
}

func (db *SqlDB) CreateTable() error {
	_, err := db.db.Exec(
		fmt.Sprintf(`
		CREATE TABLE %s (
			%s INTEGER PRIMARY KEY AUTOINCREMENT,
			%s INTEGER,
			%s TEXT UNIQUE,
			%s INTEGER,
			%s INTEGER)
			`,
			TABLE_NAME,
			COL_SEQ, COL_OP, COL_PATH, COL_SIZE, COL_MTIME))
	return err
}

func (db *SqlDB) Write(ch Change) error {
	log.Printf("Record change: %v %s", ch.Op, ch.Path)
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE %s=?", TABLE_NAME, COL_PATH),
		ch.Path)
	if err != nil {
		tx.Rollback()
		return err
	}

	if ch.Seq == 0 {
		_, err = tx.Exec(
			fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s) VALUES (?,?,?,?)", TABLE_NAME, COL_PATH, COL_OP, COL_SIZE, COL_MTIME),
			ch.Path, ch.Op, ch.Size, ch.Mtime)
	} else {
		_, err = tx.Exec(
			fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s) VALUES (?,?,?,?,?)", TABLE_NAME, COL_SEQ, COL_PATH, COL_OP, COL_SIZE, COL_MTIME),
			ch.Seq, ch.Path, ch.Op, ch.Size, ch.Mtime)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *SqlDB) List(fromSeq int64) ([]Change, error) {
	rows, err := db.db.Query(
		fmt.Sprintf("SELECT %s,%s,%s,%s,%s FROM %s WHERE %s>? ORDER BY %s",
			COL_SEQ, COL_OP, COL_PATH, COL_SIZE, COL_MTIME, TABLE_NAME, COL_SEQ, COL_SEQ),
		fromSeq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seq int64
	var op Op
	var path string
	var size, mtime int64
	l := make([]Change, 0, 0)
	for rows.Next() {
		if err := rows.Scan(&seq, &op, &path, &size, &mtime); err != nil {
			return nil, err
		}
		l = append(l, Change{Seq: seq, Op: op, Path: path, Size: size, Mtime: mtime})
	}
	return l, rows.Err()
}

func (db *SqlDB) Seq() (int64, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var count int
	if err := tx.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", TABLE_NAME)).Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, nil
	}

	var seq int64
	err = tx.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", COL_SEQ, TABLE_NAME)).Scan(&seq)
	return seq, err
}

func (db *SqlDB) Synced(seq int64) error {
	_, err := db.db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE %s<=? AND %s=?", TABLE_NAME, COL_SEQ, COL_OP),
		seq, DEL)
	return err
}

func (db *SqlDB) Close() {
	db.db.Close()
}
