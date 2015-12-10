// Package trace provides deferred execution of instruction streams generated
// at runtime. IssueCommand adds a command to be executed while the actual
// execution is deferred until all goroutines in a trace context are blocked
// on Reads of instruction results (i.e., Promises).
package trace

import (
	"sync"
)

type trace struct {
	lock    sync.Mutex
	issue   []instr
	retired [][]instr
}

func (a *trace) signal(lockedPool *poolCtx) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.signalWithLock(lockedPool)
}

func (a *trace) signalWithLock(lockedPool *poolCtx) error {
	// XXX don't execute until all active all blocked
	if lockedPool.alive == len(lockedPool.blocked) {
		return a.execute()
	}
	return nil
}

func (a *trace) execute() error {
	// XXX: Sort instructions

	for _, v := range a.issue {
		v.promise.set(nil)
		close(v.promise.out)
	}

	a.retired = append(a.retired, a.issue)
	a.issue = nil

	return nil
}
