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
// 2 Royal College St, London NW1 0NH UK

// support package with wrapper classes for marshalling parameters into
// elements
package execute

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// type to allow generic access to Antha Elements
type AnthaElement interface {
	anthaElement()
}

type ThreadID string
type BlockID struct {
	ThreadID    ThreadID
	OutputCount int
}

func (b BlockID) String() string {
	return fmt.Sprintf("%s,%d", string(b.ThreadID), b.OutputCount)
}
func StringToBlockID(in string) (*BlockID, error) {
	//TODO add format checking etc etc
	s := strings.Split(in, ",")
	count, _ := strconv.Atoi(s[1])
	return &BlockID{ThreadID: ThreadID(s[0]), OutputCount: count}, nil
}

type ThreadParam struct {
	Value   interface{}
	ID      ThreadID
	BlockID BlockID
	Error   bool
}

type BlockConfig struct {
	BlockID BlockID
	Threads map[string]string
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
	if _, ok := a.bag[key]; !ok {
		a.bag[key] = value
	}

	// if we completed our bag, fire competion
	if len(a.bag) == a.keys {
		// should we unlock mutex during this?
		values := a.mapper.Map(a.bag)
		a.completer.Complete(values)
		return true
	}
	return false
}

//PortInfo describes a port from a ComponentInfo
type PortInfo struct {
	Id          string        `json:"id"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Addressable bool          `json:"addressable"` // ignored
	Required    bool          `json:"required"`
	Values      []interface{} `json:"values"`  // ignored
	Default     interface{}   `json:"default"` // ignored
}

//ComponentInfo describes a protocol as a fbp component
type ComponentInfo struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Subgraph    bool       `json:"subgraph"`
	InPorts     []PortInfo `json:"inPorts"`
	OutPorts    []PortInfo `json:"outPorts"`
}

//NewComponentInfo returns a new ComponentInfo initalized with the given information
func NewComponentInfo(pkgName string, description string, icon string, subgraph bool, inPorts []PortInfo, outPorts []PortInfo) *ComponentInfo {
	ret := new(ComponentInfo)
	ret.Name = pkgName
	ret.Description = description
	ret.Subgraph = subgraph
	ret.Icon = icon
	ret.InPorts = inPorts
	ret.OutPorts = outPorts

	return ret
}

//NewPortInfo returns a new PortInfo struct initialized with the given values
func NewPortInfo(id string, portInfoType string, description string, addressable bool, required bool, values []interface{}, defaultValue interface{}) *PortInfo {
	ret := new(PortInfo)
	ret.Id = id
	ret.Type = portInfoType
	ret.Description = description
	ret.Addressable = addressable
	ret.Required = required
	ret.Values = values
	ret.Default = defaultValue

	return ret
}

//GraphConnection describes a connection between two processes inside a GraphDescription
type GraphConnection struct {
	Data interface{} `json:",omitempty"`
	Src  struct {
		Process string
		Port    string
	} `json:",omitempty"`
	Tgt struct {
		Process string
		Port    string
	}
	Metadata struct {
		Buffer int `json:",omitempty"`
	} `json:",omitempty"`
}

//GraphDescription describes a fbp graph.
type GraphDescription struct {
	Properties struct {
		Name string
	}
	Processes map[string]struct {
		Component string
		Metadata  struct {
			Sync                 bool   `json:",omitempty"`
			PoolSize             int64  `json:",omitempty"`
			NameSpaceClass       string `json:",omitempty"`
			ComponentLibraryName string `json:",omitempty"`
		} `json:",omitempty"`
	}
	Connections []GraphConnection
	Exports     []struct {
		Private string
		Public  string
	}
	InPorts map[string]struct {
		Process string
		Port    string
	}
	OutPorts map[string]struct {
		Process string
		Port    string
	}
}

//JSONValue holds information for a pair key value inside a JSONBlock
type JSONValue struct {
	Name       string
	JSONString string
}

//JSONBlock holds information from a JSON string in a key value fashion except for ID and Error keys
// can be used to unmarshal unknown structs and decide types on the fly
type JSONBlock struct {
	ID     *ThreadID
	Error  *bool
	Values map[string]interface{}
}

//UnmarshalJson JSONBlock from a json string saving ID and Error and the rest of information as a
// key value pair
func (j *JSONBlock) UnmarshalJSON(in []byte) error {
	tmp := make(map[string]interface{})
	json.Unmarshal(in, &tmp)

	var id ThreadID
	var ids string
	ids = tmp["ID"].(string)
	id = ThreadID(ids)
	j.ID = &id
	var err bool
	err = tmp["Error"].(bool)
	j.Error = &err
	delete(tmp, "ID")
	delete(tmp, "Error")
	j.Values = tmp
	return nil
}

//MarshalJSON builds a json string containing a JSONBlock structure in which ID and Error are added explicitly
// and the rest of fields are added from a key/value pair map
func (j *JSONBlock) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for k, v := range j.Values {
		tmp[k] = v
	}
	tmp["ID"] = j.ID
	tmp["Error"] = j.Error

	return json.Marshal(&tmp)
}
