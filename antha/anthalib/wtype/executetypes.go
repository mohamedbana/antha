package wtype

type ThreadID string

func (a ThreadID) String() string {
	return string(a)
}

type BlockID struct {
	ThreadID ThreadID
	_        int
}
