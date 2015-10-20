// /anthalib/driver/incubating/manual/manualdriver.go: Part of the Antha language
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

package manual

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipmentManager"
)

type ManualDriver struct {
	eq equipment.Equipment
}

func (m *ManualDriver) sendActionToEquipment(ac equipment.ActionDescription) error {
	go m.eq.Do(ac)
	return nil
}

//NewManualDriver returns a new instance of a manual driver pointing to the right piece of equipment
func NewManualDriver() *ManualDriver {
	ret := new(ManualDriver)
	eqm := equipmentManager.GetEquipmentManager()
	params := make(map[string]string, 0)
	ret.eq = eqm.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", params)) //TODO make this into something more meaningful
	return ret
}
func (m *ManualDriver) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus {
	params := make(map[string]string)
	params["deckposition"] = fmt.Sprintf("%v", deckposition[0])
	params["wellcoords"] = fmt.Sprintf("%v", wellcoords[0])
	params["reference"] = fmt.Sprintf("%v", reference)
	params["offsetX"] = fmt.Sprintf("%v", offsetX)
	params["offsetY"] = fmt.Sprintf("%v", offsetY)
	params["offsetZ"] = fmt.Sprintf("%v", offsetZ)
	params["plate_type"] = fmt.Sprintf("%v", plate_type)
	params["head"] = fmt.Sprintf("%v", head)

	desc := fmt.Sprintf("Deck Postition %v @well %v with reference %v", deckposition, wellcoords, reference)
	ad := *equipment.NewActionDescription(action.LH_MOVE, desc, params)
	m.sendActionToEquipment(ad)
	//	go m.eq.Do(ad)
	//	err := m.eq.Do(ad)
	//	if err != nil {
	//		return driver.CommandStatus{
	//			OK:        false,
	//			Errorcode: 1,
	//			Msg:       err.Error(),
	//		}
	//	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Initialize() driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Finalize() driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}

}
func (m *ManualDriver) Incubate(what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) driver.CommandStatus {
	params := make(map[string]string)
	params["what"] = fmt.Sprintf("%v", what)
	params["temp"] = fmt.Sprintf("%v", temp)
	params["time"] = fmt.Sprintf("%v", time)
	params["shaking"] = fmt.Sprintf("%v", shaking)
	var desc string
	var act action.Action
	if shaking {
		desc = fmt.Sprintf("Incubate shaking %s for %s at a temperature of %s.", params["what"], params["time"], params["temp"])
		act = action.IN_INCUBATE_SHAKE
	} else {
		desc = fmt.Sprintf("Incubate %s for %s at a temperature of %s.", params["what"], params["time"], params["temp"])
		act = action.IN_INCUBATE
	}
	ad := *equipment.NewActionDescription(act, desc, params)
	m.sendActionToEquipment(ad)
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
