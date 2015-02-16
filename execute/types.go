// antha/execute/types.go: Part of the Antha language
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

// support package with wrapper classes for marshalling parameters into
// elements
package execute

import (
	"sync"
)

// type to allow generic access to Antha Elements
type AnthaElement interface {
	anthaElement()
}

type ThreadID string

type ThreadParam struct {
	Value interface{}
	ID    ThreadID
}

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
func (a *AsyncBag) Init(keys int, completer AsyncCompleter, mapper AsyncMapper) {
	a.keys = keys
	a.completer = completer
	a.mapper = mapper
	a.bag = make(map[string]interface{})
}

// adds value and returns true if the bag was fired
// TODO: Should the completion be wrapped in a sync.Once in case
// there are duplicate params flowing through the network with the
// same threadID?
func (a *AsyncBag) AddValue(key string, value interface{}) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.bag[key] = value

	// if we completed our bag, fire competion
	if len(a.bag) == a.keys {
		// should we unlock mutex during this?
		values := a.mapper.Map(a.bag)
		a.completer.Complete(values)
		return true
	}
	return false
}
