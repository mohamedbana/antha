package trace

import (
	//"log"
	"sync"
)

type Trace struct {
	lock    sync.Mutex
	issue   []instr
	retired [][]instr
}

func (a *Trace) signal(lockedPool *poolCtx) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.signalWithLock(lockedPool)
}

func (a *Trace) signalWithLock(lockedPool *poolCtx) error {
	// XXX don't execute until all active all blocked
	if lockedPool.alive == len(lockedPool.blocked) {
		return a.execute()
	}
	return nil
}

func (a *Trace) execute() error {
	// XXX: Sort instructions
	for _, v := range a.issue {
		v.promise.set(nil)
		close(v.promise.out)
	}

	a.retired = append(a.retired, a.issue)
	a.issue = nil

	return nil
}
