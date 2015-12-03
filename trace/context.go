package trace

import (
	"errors"
	"golang.org/x/net/context"
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

type DoneFunc func() <-chan struct{}

func WithPool(parent context.Context) (context.Context, context.CancelFunc, DoneFunc) {
	done := make(chan struct{})
	pctx := &poolCtx{
		Context: parent,
		done:    done,
		blocked: make(map[*Promise]bool),
	}
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
			return done
		}
}

func WithScope(parent context.Context) context.Context {
	pscope, _ := parent.Value(theScopeKey).(*Scope)
	var s *Scope
	if pscope == nil {
		s = &Scope{}
	} else {
		s = pscope.MakeScope()
	}

	return context.WithValue(parent, theScopeKey, s)
}

func WithTrace(parent context.Context) context.Context {
	return context.WithValue(parent, theInstrKey, &Trace{})
}

func NewContext(parent context.Context) (context.Context, context.CancelFunc, DoneFunc) {
	c, cfn, dfn := WithPool(WithTrace(WithScope(parent)))
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

func getTrace(ctx context.Context) *Trace {
	t := ctx.Value(theInstrKey).(*Trace)
	if t == nil {
		panic("trace: trace not found")
	}
	return t
}

func MakeValue(ctx context.Context, desc string, v interface{}) Value {
	return &basicValue{
		name: getScope(ctx).MakeName(desc),
		v:    v,
	}
}

func MakeValueFrom(ctx context.Context, desc string, v interface{}, from ...Value) Value {
	return &fromValue{
		name: getScope(ctx).MakeName(desc),
		v:    v,
		from: from,
	}
}

func IssueCommand(ctx context.Context, op string, args ...Value) *Promise {
	result := getScope(ctx).MakeName("")
	out := make(chan interface{})

	p := &Promise{
		construct: func(v interface{}) Value {
			return &promiseValue{
				name: result,
				op:   op,
				from: args,
				v:    v,
			}
		},
		out: out,
	}

	t := getTrace(ctx)
	t.lock.Lock()
	defer t.lock.Unlock()

	t.issue = append(t.issue, instr{op: op, args: args, promise: p})

	return p
}

func tryUnblock(ctx context.Context, pctx *poolCtx, tr *Trace) {
	if pctx.alive == len(pctx.blocked) {
		if err := tr.execute(ctx); err != nil {
			pctx.cancelWithLock(err)
		}
	}
}

func decrement(ctx context.Context, pctx *poolCtx, tr *Trace, err error) {
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	pctx.alive -= 1
	tryUnblock(ctx, pctx, tr)

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

func Read(ctx context.Context, p *Promise) (Value, error) {
	if v := p.Value(); v != nil {
		return v, nil
	}

	pctx := getPool(ctx)
	tr := getTrace(ctx)

	pctx.lock.Lock()
	pctx.blocked[p] = true
	tryUnblock(ctx, pctx, tr)
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
		err = p.Err()
	}

	pctx.lock.Lock()
	delete(pctx.blocked, p)
	pctx.lock.Unlock()

	return p.Value(), err
}

func ReadAll(parent context.Context, ps ...*Promise) ([]Value, error) {
	vs := make(map[int]Value)

	pctx, cancel, allDone := WithPool(parent)
	defer cancel()

	for idx, p := range ps {
		i := idx
		Go(pctx, func(ctx context.Context) error {
			v, err := Read(ctx, p)
			if err != nil {
				return err
			}
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

func Go(parent context.Context, fn func(ctx context.Context) error) {
	pctx := getPool(parent)
	ctx := WithScope(parent)
	tr := getTrace(parent)

	pctx.lock.Lock()
	defer pctx.lock.Unlock()

	pctx.alive += 1

	go func() {
		err := fn(ctx)
		decrement(ctx, pctx, tr, err)
	}()
}
