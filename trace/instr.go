package trace

import (
	"sync"
)

type Value interface {
	Name() Name
	Get() interface{}
}

type Promise struct {
	lock      sync.Mutex
	construct func(interface{}) Value
	value     Value
	err       error
	out       chan interface{}
}

func (a *Promise) Err() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.err
}

func (a *Promise) Value() Value {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.value
}

func (a *Promise) set(v interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.value != nil {
		a.value = a.construct(v)
	}
}

type basicValue struct {
	name Name
	v    interface{}
}

func (a *basicValue) Get() interface{} {
	return a.v
}

func (a *basicValue) Name() Name {
	return a.name
}

type fromValue struct {
	name Name
	v    interface{}
	from []Value
}

func (a *fromValue) Get() interface{} {
	return a.v
}

func (a *fromValue) Name() Name {
	return a.name
}

type promiseValue struct {
	name Name
	v    interface{}
	from []Value
	op   string
}

func (a *promiseValue) Get() interface{} {
	return a.v
}

func (a *promiseValue) Name() Name {
	return a.name
}

type Scope struct {
	lock   sync.Mutex
	parent *Scope
	pidx   int
	count  int
}

func (a *Scope) MakeScope() *Scope {
	a.lock.Lock()
	defer a.lock.Unlock()
	s := &Scope{parent: a, pidx: a.count}
	a.count += 1
	return s
}

func (a *Scope) MakeName(desc string) Name {
	a.lock.Lock()
	defer a.lock.Unlock()
	n := Name{scope: a, idx: a.count, desc: desc}
	a.count += 1
	return n
}

type Name struct {
	scope *Scope
	idx   int
	desc  string
}

type instr struct {
	op      string
	args    []Value
	promise *Promise
}
