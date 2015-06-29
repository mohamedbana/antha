// /equipment/equipment.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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

//package equipment defines the data representation of a piece of equipment and every necessary bits to communicate
// with them.
package equipment

import (
	"github.com/antha-lang/antha/microArch/equipment/action"
)

//ActionDescription is the representation of an action that is going to be executed and extends the particular Action
// with data on how to carry it out
type ActionDescription struct {
	Action     action.Action
	ActionData string //TODO probably make a struct with proper fields
	Params     map[string]string
}

//NewActionDescription instantiates a new action with the given data
func NewActionDescription(action action.Action, data string, params map[string]string) *ActionDescription {
	ac := new(ActionDescription)
	ac.Action = action
	ac.ActionData = data
	ac.Params = params
	return ac
}

//Behaviour represents the capabilities of a piece of Equipment to perform an action. This is related to how you can
// ask an equipment to carry out an action. The actionData in ActionDescription should meet the constraints that
// a piece of equipment describes in its behaviour
type Behaviour struct {
	Action      action.Action
	Constraints string //TODO probably make a struct with proper fields
}

//NewBehaviour will instantiate a new Behaviour matching the given action and constraints
func NewBehaviour(action action.Action, constraints string) *Behaviour {
	b := new(Behaviour)
	b.Action = action
	b.Constraints = constraints
	return b
}

//Matches checks whether the action description can be carried out by this behaviour
func (b *Behaviour) Matches(ac ActionDescription) bool {
	if b.Action != ac.Action {
		return false
	}
	//TODO do something with the constrains!!!!!!
	return true
}

//Equipment is something capable of performing different actions under different restrictions and explaining what its
// capabilities and the restrictions on them are
type Equipment interface {
	//GetID returns the string that identifies a piece of equipment. Ideally uuids v4 should be used.
	GetID() string
	//GetEquipmentDefinition returns a description of the equipment device in terms of
	// operations it can handle, restrictions, configuration options ...
	GetEquipmentDefinition()
	//Perform an action in the equipment. Actions might be transmitted in blocks to the equipment
	// The grouping of the actions (as a set, plate or whatever) is not performed at the equipment driver level
	// or is it?
	Do(actionDescription ActionDescription) error
	//Can queries a piece of equipment about an action execution. The description of the action must meet the constraints
	// of the piece of equipment.
	Can(ac ActionDescription) bool
	//Status should give a description of the current execution status and any future actions queued to the device
	Status() string
	//Init driver will be initialized when registered
	Init() error
	//Shutdown disconnect, turn off, signal whatever is necessary for a graceful shutdown
	Shutdown() error
}
