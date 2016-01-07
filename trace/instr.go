package trace

import (
	"sync"
)

// Wrapper around instruction inputs and outputs
type Value interface {
	Name() Name
	Get() interface{}
}

// Placeholder for value to be returned by an instruction
type Promise struct {
	lock      sync.Mutex
	construct func(interface{}) Value
	v         Value
	e         error
	out       chan interface{}
}

func (a *Promise) err() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.e
}

func (a *Promise) value() Value {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.v
}

func (a *Promise) set(v interface{}) {
	if a.construct == nil {
		return
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.v != nil {
		a.v = a.construct(v)
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

// Scope for a name
type Scope struct {
	lock   sync.Mutex
	parent *Scope
	pidx   int
	count  int
}

// Create a child name scope
func (a *Scope) MakeScope() *Scope {
	a.lock.Lock()
	defer a.lock.Unlock()
	s := &Scope{parent: a, pidx: a.count}
	a.count += 1
	return s
}

// Make a name in this scope
func (a *Scope) MakeName(desc string) Name {
	a.lock.Lock()
	defer a.lock.Unlock()
	n := Name{scope: a, idx: a.count, desc: desc}
	a.count += 1
	return n
}

// A name
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
