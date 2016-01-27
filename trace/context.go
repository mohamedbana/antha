package trace

import (
	"errors"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"runtime/debug"
	"sync"
)

type scopeKey int
type instrKey int
type poolKey int

const theScopeKey scopeKey = 0
const theInstrKey instrKey = 0
const thePoolKey poolKey = 0

var poolDoneError error = errors.New("trace: done error") // Dummy error to signal sucessful execution

type poolCtx struct {
	context.Context
	lock    sync.Mutex
	alive   int
	blocked map[*Promise]bool
	err     error
	done    chan struct{}
}

func (a *poolCtx) Value(key interface{}) interface{} {
	if key == thePoolKey {
		return a
	}
	return a.Context.Value(key)
}

func (a *poolCtx) Err() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.err == nil {
		return a.Context.Err()
	}
	if a.err == poolDoneError {
		return nil
	}
	return a.err
}

func (a *poolCtx) Done() <-chan struct{} {
	return a.done
}

func (a *poolCtx) cancelWithLock(err error) {
	if a.err == nil {
		a.err = err
		close(a.done)
	}
}

func (a *poolCtx) cancel(err error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.cancelWithLock(err)
}

func (a *poolCtx) remove(p *Promise) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.blocked, p)
}

type DoneFunc func() <-chan struct{}

// Create a pool context. It is done when all Go()-created routines return
// normally or when any return an error.
func WithPool(parent context.Context) (context.Context, context.CancelFunc, DoneFunc) {
	tr := getTrace(parent)
	done := make(chan struct{})
	pctx := &poolCtx{
		Context: parent,
		done:    done,
		alive:   1,
		blocked: make(map[*Promise]bool),
	}
	rootPromise := &Promise{}
	pctx.blocked[rootPromise] = true

	if parent.Done() != nil {
		go func() {
			select {
			case <-parent.Done():
				pctx.cancel(parent.Err())
			case <-done:
			}
		}()
	}

	return pctx,
		func() { pctx.cancel(context.Canceled) },
		func() <-chan struct{} {
			pctx.remove(rootPromise)
			decrement(pctx, tr, 1, nil)
			return done
		}
}

func withScope(parent context.Context) context.Context {
	pscope, _ := parent.Value(theScopeKey).(*Scope)
	var s *Scope
	if pscope == nil {
		s = &Scope{}
	} else {
		s = pscope.MakeScope()
	}

	return context.WithValue(parent, theScopeKey, s)
}

func withTrace(parent context.Context) context.Context {
	return context.WithValue(parent, theInstrKey, &trace{})
}

// Create a new trace and pool context. The general template for clients is:
//   ctx, cancel, allDone := NewContext(parent)
//   defer cancel()
//   ...
//   Go(ctx, ...)
//   ...
//   select {
//     case <-allDone():
//       return ctx.Err()
//     ...
//   }
func NewContext(parent context.Context) (context.Context, context.CancelFunc, DoneFunc) {
	c, cfn, dfn := WithPool(withTrace(withScope(parent)))
	return c, cfn, dfn
}

func getPool(ctx context.Context) *poolCtx {
	c := ctx.Value(thePoolKey).(*poolCtx)
	if c == nil {
		panic("trace: pool not found")
	}
	return c
}

func getScope(ctx context.Context) *Scope {
	s := ctx.Value(theScopeKey).(*Scope)
	if s == nil {
		panic("trace: scope not defined")
	}
	return s
}

func getTrace(ctx context.Context) *trace {
	t := ctx.Value(theInstrKey).(*trace)
	if t == nil {
		panic("trace: trace not found")
	}
	return t
}

func tryUnblock(tr *trace, pctx *poolCtx) {
	if err := tr.signal(pctx); err != nil {
		pctx.cancelWithLock(err)
	}
}

func decrement(pctx *poolCtx, tr *trace, delta int, err error) {
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	pctx.alive -= delta
	tryUnblock(tr, pctx)

	var cancel error
	if err != nil {
		cancel = err
	}
	if pctx.alive == 0 && cancel == nil {
		cancel = poolDoneError
	}
	if cancel != nil {
		pctx.cancelWithLock(cancel)
	}
}

// Read a promise value
func Read(ctx context.Context, p *Promise) (Value, error) {
	if v := p.value(); v != nil {
		return v, nil
	}

	pctx := getPool(ctx)
	tr := getTrace(ctx)

	pctx.lock.Lock()
	pctx.blocked[p] = true
	tryUnblock(tr, pctx)
	pctx.lock.Unlock()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case v, ok := <-p.out:
		if ok {
			p.set(v)
		}
	}

	if err == nil {
		err = p.err()
	}

	pctx.remove(p)

	return p.value(), err
}

// Read all promises concurrently
func ReadAll(parent context.Context, ps ...*Promise) ([]Value, error) {
	var mapLock sync.Mutex
	vs := make(map[int]Value)

	pctx, cancel, allDone := WithPool(parent)
	defer cancel()

	for idx, promise := range ps {
		i := idx
		p := promise
		Go(pctx, func(ctx context.Context) error {
			v, err := Read(ctx, p)
			if err != nil {
				return err
			}

			mapLock.Lock()
			defer mapLock.Unlock()
			vs[i] = v

			return nil
		})
	}

	<-allDone()
	if err := pctx.Err(); err != nil {
		return nil, err
	}

	var ret []Value
	for i := 0; i < len(ps); i += 1 {
		ret = append(ret, vs[i])
	}
	return ret, nil
}

// Create a new goroutine in the pool context
func Go(parent context.Context, fn func(ctx context.Context) error) {
	pctx := getPool(parent)
	ctx := withScope(parent)
	tr := getTrace(parent)

	pctx.lock.Lock()
	defer pctx.lock.Unlock()

	pctx.alive += 1

	go func() {
		var err error
		defer func() {
			decrement(pctx, tr, 1, err)
		}()
		defer func() {
			if res := recover(); res != nil {
				err = &goError{BaseError: res, Stack: debug.Stack()}
			}
		}()
		err = fn(ctx)
	}()
}
