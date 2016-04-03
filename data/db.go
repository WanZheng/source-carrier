package data

type DB interface {
	Write(change Change) error
	List(fromSeq int64) ([]Change, error)
	Seq() (int64, error)
	Synced(seq int64) error
	Close()
}
