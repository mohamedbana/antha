package trace

import (
	"golang.org/x/net/context"
	"sync"
)

type Trace struct {
	lock    sync.Mutex
	issue   []instr
	retired [][]instr
}

func (a *Trace) execute(context.Context) error {
	// XXX: Sort instructions
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, v := range a.issue {
		v.promise.set(nil)
		close(v.promise.out)
	}

	a.retired = append(a.retired, a.issue)
	a.issue = nil

	return nil
}
