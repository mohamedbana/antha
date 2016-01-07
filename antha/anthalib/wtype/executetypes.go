package wtype

type ThreadID string

type BlockID struct {
	ThreadID ThreadID
	_        int
}
