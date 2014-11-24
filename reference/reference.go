// antha/reference/reference.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

// Example of threading model

package reference

import (
	"github.com/Synthace/goflow"
	"sync"
	"time"
)

// support function to fire when a full bag of values has arrived
type AsyncCompleter interface {
	Complete(interface{})
}

// support function to map into a concrete struct
type AsyncMapper interface {
	Map(map[string]interface{}) interface{}
}

// Simple structure to coordinate the asynchronous aggregation of multiple
// values that have to be fired together
type AsyncBag struct {
	bag       map[string]interface{}
	keys      int
	completer AsyncCompleter
	mapper    AsyncMapper
	lock      sync.Mutex
}

// makes a new AsyncBag which requires keys to fire f
func (a *AsyncBag) init(keys int, completer AsyncCompleter, mapper AsyncMapper) {
	a.keys = keys
	a.completer = completer
	a.mapper = mapper
	a.bag = make(map[string]interface{})
}

// adds value and returns true if the bag was fired
// TODO: Should the competion be wrapped in a sync.Once in case
// there are duplicate params flowing through the network with the
// same threadID?
func (a *AsyncBag) AddValue(key string, value interface{}) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.bag[key] = value

	// if we completed our bag, fire competion
	if len(a.bag) == a.keys {
		// should we unlock mutex during this?
		go func() {
			values := a.mapper.Map(a.bag)
			a.completer.Complete(values)
		}()
		return true
	}
	return false
}

var params = [...]string{
	"Color",
	"SleepTime"}

type ThreadID string

type ThreadParam struct {
	Value interface{}
	ID    ThreadID
}

// channel interfaces
// with threadID grouped types
type Example struct {
	flow.Component                    // component "superclass" embedded
	Color          <-chan ThreadParam // color to make this well
	SleepTime      <-chan ThreadParam // amount of time to randomly wait
	WellColor      chan<- ThreadParam // output color
	lock           sync.Mutex
	params         map[ThreadID]*AsyncBag
}

// single execution thread variables
// with concrete types
type paramBlock struct {
	Color     string
	SleepTime time.Duration
	WellColor string
	ID        ThreadID
}

// could handle mapping in the threadID better...
func (e *Example) Map(m map[string]interface{}) interface{} {
	var res paramBlock
	res.Color = m["Color"].(ThreadParam).Value.(string)
	res.SleepTime = m["SleepTime"].(ThreadParam).Value.(time.Duration)
	res.ID = m["Color"].(ThreadParam).ID
	return res
}

func (e *Example) OnColor(param ThreadParam) {
	e.lock.Lock()
	var bag *AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(AsyncBag)
		bag.init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Color", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

func (e *Example) OnSleepTime(param ThreadParam) {
	e.lock.Lock()
	var bag *AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(AsyncBag)
		bag.init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SleepTime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

// AsyncBag functions
func (e *Example) Complete(params interface{}) {
	p := params.(paramBlock)
	go e.steps(p)
}

// generic typing for interface support
func (e *Example) anthaElement() {}

// init function, read characterization info from seperate file to validate ranges?
func (e *Example) init() {
	e.params = make(map[ThreadID]*AsyncBag)
}

func NewExample() *Example {
	e := new(Example)
	e.init()
	return e
}

// main function for use in goroutines
func (e *Example) steps(p paramBlock) {
	time.Sleep(p.SleepTime)
	e.WellColor <- ThreadParam{p.Color, p.ID}
}