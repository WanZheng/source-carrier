package data

type Op int

const (
	UPDATE Op = iota
	DEL
)

func (o Op) String() string {
	if o == UPDATE {
		return "UPD"
	} else {
		return "DEL"
	}
}

type Change struct {
	Seq   int64
	Op    Op
	Path  string
	Size  int64
	Mtime int64
	Data  []byte
}
