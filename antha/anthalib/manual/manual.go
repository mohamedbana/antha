// /anthalib/manual/manual.go: Part of the Antha language
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

package manual

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/driver"
	"github.com/antha-lang/antha/internal/github.com/nu7hatch/gouuid"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipment/manual"
	"github.com/antha-lang/antha/microArch/equipmentManager"
)

type Manual struct {
	Manager        equipmentManager.AnthaEquipmentManager
	ManualInstance manual.AnthaManual
}

func NewManual() *Manual {
	manual := new(Manual)
	return manual
}

func (m *Manual) Init() { //TODO add the proper return statements with a commandStatus
	id, err := uuid.NewV4()
	if err != nil {
		panic(err) //TODO
	}
	m.Manager = *equipmentManager.NewAnthaEquipmentManager(id.String())
	idm, err := uuid.NewV4()
	if err != nil {
		panic(err) //TODO
	}
	var md equipment.Equipment
	md = manual.NewAnthaManual(idm.String())
	m.Manager.RegisterEquipment(&md)
}

func (m *Manual) Message(message string) driver.CommandStatus {
	var act action.Action
	act = action.MESSAGE
	ret := new(driver.CommandStatus)
	var eq *equipment.Equipment
	params := make(map[string]string, 0)
	//	fmt.Println(m.Manager)
	acd := equipment.NewActionDescription(act, message, params)
	//	fmt.Println(*acd)
	eq = m.Manager.GetActionCandidate(*acd)
	if eq == nil {
		ret.Errorcode = driver.ERR
		ret.Msg = fmt.Sprintf("Could not find suitable equipment for action %v : %v.", act, message)
		ret.OK = false
		return *ret
	}
	var actionEq equipment.Equipment
	actionEq = *eq
	err := actionEq.Do(*acd)
	if err != nil {
		ret.Errorcode = driver.ERR
		ret.Msg = err.Error()
		ret.OK = false
		return *ret

	}
	ret.Errorcode = driver.OK
	ret.OK = true
	return *ret
}
