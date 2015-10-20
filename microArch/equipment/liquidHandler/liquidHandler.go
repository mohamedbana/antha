// /equipment/liquidHandler/liquidHandler.go: Part of the Antha language
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
// 2 Royal College St, London NW1 0NH UK

//package liquidHandler defines a liquid handler implementation as an Antha compatible
// equipment.
package liquidHandler

import (
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
)

//AnthaLiquidHandler represents a liquidHandler that can be identified as an antha compatible
// device, represented by and ID and that responds to a certain set of Behaviours
type AnthaLiquidHandler struct {
	ID         string
	Behaviours []equipment.Behaviour
}

//NewAnthaLiquidHandler instantiates a LiquidHandler identified by id and supporting the following behaviours:
// action.LH_MOVE
// action.LH_MOVE_EXPLICIT
// action.LH_MOVE_RAW
// action.LH_ASPIRATE
// action.LH_DISPENSE
// action.LH_LOAD_TIPS
// action.LH_UNLOAD_TIPS
// action.LH_SET_PIPPETE_SPEED
// action.LH_SET_DRIVE_SPEED
// action.LH_STOP
// action.LH_SET_POSITION_STATE
// action.LH_RESET_PISTONS
// action.LH_WAIT
// action.LH_MIX
func NewAnthaLiquidHandler(id string) *AnthaLiquidHandler {
	//Our liquid handler is going to be able to mix and move liquid
	be := make([]equipment.Behaviour, 0)
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_EXPLICIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_RAW, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_ASPIRATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_DISPENSE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_LOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_UNLOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_PIPPETE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_DRIVE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_STOP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_POSITION_STATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_RESET_PISTONS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_WAIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MIX, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_CONFIG, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_READ, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_END, ""))

	eq := new(AnthaLiquidHandler)
	eq.Behaviours = be
	eq.ID = id

	return eq
}

//GetID returns the unique id for this liquid handler
func (e AnthaLiquidHandler) GetID() string {
	return e.ID
}

//GetEquipmentDefinition returns the equipment definition for the liquid handler.
// This funcionality is still T.B.D. in terms of parameters and returns
func (e AnthaLiquidHandler) GetEquipmentDefinition() {
	//TODO
}

// This funcionality is still T.B.D. in terms of parameters and returns
func (e AnthaLiquidHandler) Do(actionDescription equipment.ActionDescription) error {
	//TODO
	return nil
}

//Status should give a description of the current execution status and any future actions queued to the device
func (e *AnthaLiquidHandler) Status() string {
	//TODO implement properly
	return "OK"
}

//Can queries a piece of equipment about an action execution. The description of the action must meet the constraints
// of the piece of equipment.
func (e *AnthaLiquidHandler) Can(b equipment.ActionDescription) bool {
	for _, eb := range e.Behaviours {
		if eb.Matches(b) {
			return true
		}
	}
	return false
}

//Init driver will be initialized when registered
func (e *AnthaLiquidHandler) Init() error {
	return nil
}

//Shutdown disconnect, turn off, signal whatever is necessary for a graceful shutdown
func (e *AnthaLiquidHandler) Shutdown() error {
	return nil
}
