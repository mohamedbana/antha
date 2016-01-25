// Package trace provides deferred execution of instruction streams generated
// at runtime. IssueCommand adds a command to be executed while the actual
// execution is deferred until all goroutines in a trace context are blocked
// on Reads of instruction results (i.e., Promises).
package trace

import (
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"sync"
)

type trace struct {
	lock    sync.Mutex
	issue   []instp
	retired [][]instp
}

func (a *trace) signal(lockedPool *poolCtx) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.signalWithLock(lockedPool)
}

func (a *trace) signalWithLock(lockedPool *poolCtx) error {
	// Don't execute until all active all blocked
	if lockedPool.alive == len(lockedPool.blocked) {
		return a.execute(lockedPool.Context)
	}
	return nil
}

func (a *trace) execute(ctx context.Context) error {
	// TODO: Deterministic sort of instructions

	if err := run(ctx, a.issue); err != nil {
		return err
	}

	for _, v := range a.issue {
		// TODO: Update this when clients actually use results of promises
		v.promise.set(nil)
		close(v.promise.out)
	}

	a.retired = append(a.retired, a.issue)
	a.issue = nil

	return nil
}
