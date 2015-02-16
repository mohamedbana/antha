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
	"bytes"
	"encoding/json"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/goflow"
	"io"
	"log"
	"sync"
	"time"
)

var params = [...]string{
	"Color",
	"SleepTime"}

type WellColorParam struct {
	WellColor string
}

// channel interfaces
// with threadID grouped types
type Example struct {
	flow.Component                            // component "superclass" embedded
	Color          <-chan execute.ThreadParam // color to make this well
	SleepTime      <-chan execute.ThreadParam // amount of time to randomly wait
	WellColor      chan<- execute.ThreadParam // output color
	lock           sync.Mutex
	params         map[execute.ThreadID]*execute.AsyncBag
}

// single execution thread variables
// with concrete types
type ParamBlock struct {
	Color     string
	SleepTime time.Duration
	WellColor string
	ID        execute.ThreadID
}

type JSONBlock struct {
	Color     *string
	SleepTime *time.Duration
	WellColor *string
	ID        *execute.ThreadID
}

// support function for wire format
func (p *ParamBlock) ToJSON() (b bytes.Buffer) {
	enc := json.NewEncoder(&b)
	if err := enc.Encode(p); err != nil {
		log.Fatalln(err) // currently fatal error
	}
	return
}

// helper generator function
func ParamsFromJSON(r io.Reader) (p *ParamBlock) {
	p = new(ParamBlock)
	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		log.Fatalln(err)
	}
	return
}

// could handle mapping in the threadID better...
func (e *Example) Map(m map[string]interface{}) interface{} {
	var res ParamBlock
	res.Color = m["Color"].(execute.ThreadParam).Value.(string)
	res.SleepTime = m["SleepTime"].(execute.ThreadParam).Value.(time.Duration)
	res.ID = m["Color"].(execute.ThreadParam).ID
	return res
}

func (e *Example) OnColor(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
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

func (e *Example) OnSleepTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
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

// execute.AsyncBag functions
func (e *Example) Complete(params interface{}) {
	p := params.(ParamBlock)
	e.steps(p)
}

// generic typing for interface support
func (e *Example) anthaElement() {}

// init function, read characterization info from seperate file to validate ranges?
func (e *Example) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func NewExample() *Example {
	e := new(Example)
	e.init()
	return e
}

// main function for use in goroutines
func (e *Example) steps(p ParamBlock) {
	time.Sleep(p.SleepTime)
	e.WellColor <- execute.ThreadParam{p.Color, p.ID}
}
