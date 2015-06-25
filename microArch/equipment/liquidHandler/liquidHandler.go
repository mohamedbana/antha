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
// 1 Royal College St, London NW1 0NH UK

package liquidHandler

import (
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
)

type AnthaLiquidHandler struct {
	ID         string
	Behaviours []equipment.Behaviour
}

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

	//	var eq AnthaLiquidHandler
	//	eq = AnthaLiquidHandler{*equipment.NewAnthaEquipment(id, be)}

	eq := new(AnthaLiquidHandler)
	eq.Behaviours = be
	eq.ID = id

	return eq
}

//func NewAnthaEquipment(id string, bs []Behaviour) *AnthaEquipment {
//	ret := new(AnthaEquipment)
//	ret.ID = id
//	ret.Behaviours = bs
//	return ret
//}

func (e AnthaLiquidHandler) GetID() string {
	return e.ID
}

func (e AnthaLiquidHandler) GetEquipmentDefinition() {
	//TODO
}
func (e AnthaLiquidHandler) Do(actionDescription equipment.ActionDescription) error {
	//TODO
	return nil
}

func (e AnthaLiquidHandler) Can(b equipment.ActionDescription) bool {
	for _, eb := range e.Behaviours {
		if eb.Matches(b) {
			return true
		}
	}
	return false
}

func (e AnthaLiquidHandler) Status() string {
	//TODO implement properly
	return "OK"
}
