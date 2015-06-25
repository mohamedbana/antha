// /anthalib/driver/liquidhandling/manual/manualdriver.go: Part of the Antha language
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

	"os"

	"github.com/antha-lang/antha/antha/anthalib/driver"
	"github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipmentManager"
)

type ManualDriver struct {
	eq equipment.Equipment
	ag Aggregator
}

func (m *ManualDriver) sendActionToEquipment(ac equipment.ActionDescription) error {
	r := m.ag.addAction(ac)
	if r != nil {
		go m.eq.Do(*r)
	}
	return nil
}

//NewManualDriver returns a new instance of a manual driver pointing to the right piece of equipment
func NewManualDriver() *ManualDriver {
	ret := new(ManualDriver)
	eqm := *equipmentManager.GetEquipmentManager()
	params := make(map[string]string, 0)
	ret.eq = *eqm.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", params)) //TODO make this into something more meaningful
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

	desc := fmt.Sprintf("Deck Position %v @well %v with reference %v", deckposition, wellcoords, reference)
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
func (m *ManualDriver) MoveExplicit(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []*wtype.LHPlate, head int) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Deck Position %v @well %v with reference %v", deckposition, wellcoords, reference)
	ad := *equipment.NewActionDescription(action.LH_MOVE_EXPLICIT, desc, params)
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
func (m *ManualDriver) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
	desc := fmt.Sprintf("Move to coordinates %v, %v, %v", x, y, z)
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_MOVE_RAW, desc, params)
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
func (m *ManualDriver) Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	params := make(map[string]string)
	params["volume"] = fmt.Sprintf("%g ml.", volume[0])
	params["overstroke"] = fmt.Sprintf("%v", overstroke[0])
	params["head"] = fmt.Sprintf("%v", head)
	params["multi"] = fmt.Sprintf("%v", multi)
	params["platetype"] = fmt.Sprintf("%v", platetype[0])
	params["what"] = fmt.Sprintf("%v", what[0])
	//	params["llf"] = fmt.Sprintf("%v", llf[0])

	desc := fmt.Sprintf("Aspirate volumes %v", volume) //TOOD make a meaning of the values
	ad := *equipment.NewActionDescription(action.LH_ASPIRATE, desc, params)
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
func (m *ManualDriver) Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Dispense volumes %v", volume)
	ad := *equipment.NewActionDescription(action.LH_DISPENSE, desc, params)
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
func (m *ManualDriver) LoadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Load tips for channels %v", channels)
	ad := *equipment.NewActionDescription(action.LH_LOAD_TIPS, desc, params)
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
func (m *ManualDriver) UnloadTips(channels []int, head, multi int, platetype, position, well []string) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Unload tips for channels %v", channels)
	ad := *equipment.NewActionDescription(action.LH_UNLOAD_TIPS, desc, params)
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
func (m *ManualDriver) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	desc := fmt.Sprintf("Set pippete speed to %v for head num %d, channel num %d", rate, head, channel)
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_SET_PIPPETE_SPEED, desc, params)
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
func (m *ManualDriver) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	desc := fmt.Sprintf("Set drive %s to speed %v", drive, rate)
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_SET_DRIVE_SPEED, desc, params)
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
func (m *ManualDriver) Stop() driver.CommandStatus {
	//TODO command not sent in manual case
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Go() driver.CommandStatus {
	//TODO check the meaning of this particular call
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Initialize() driver.CommandStatus {
	//TODO chekc if additional messages need to be sent for preparation
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Finalize() driver.CommandStatus {
	//TODO check if additional things need to be done
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) SetPositionState(position string, state driver.PositionState) driver.CommandStatus {
	desc := fmt.Sprintf("Set position %s to state %v", position, state)
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_SET_POSITION_STATE, desc, params)
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

//TODO implement capabilities for manual driver
func (m *ManualDriver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	return liquidhandling.LHProperties{}, driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) GetCurrentPosition(head int) (string, driver.CommandStatus) {
	return "", driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) GetPositionState(position string) (string, driver.CommandStatus) {
	return "", driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) GetHeadState(head int) (string, driver.CommandStatus) {
	return "", driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) GetStatus() (driver.Status, driver.CommandStatus) {
	return driver.Status{}, driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) ResetPistons(head, channel int) driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Wait(time float64) driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Mix(head int, volume []float64, fvolume []float64, platetype []string, cycles []int, multi int, prms map[string]interface{}) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Mix with head num %d", head)
	ad := *equipment.NewActionDescription(action.LH_MIX, desc, params)
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

func (m *ManualDriver) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Add plate %s with name %s in position %s", plate, name, position)
	ad := *equipment.NewActionDescription(action.LH_ADD_PLATE, desc, params)
	go m.eq.Do(ad)

	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) RemoveAllPlates() driver.CommandStatus {
	desc := fmt.Sprintf("Remove all Plates from the liquid handler")
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_REMOVE_ALL_PLATES, desc, params)
	go m.eq.Do(ad)

	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) RemovePlateAt(position string) driver.CommandStatus {
	desc := fmt.Sprintf("Remove Plate at position %s", position)
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_REMOVE_PLATE, desc, params)
	go m.eq.Do(ad)

	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}

type Aggregator struct {
	Calls []equipment.ActionDescription
}

func NewAggregator() *Aggregator {
	return new(Aggregator)
}
func (a *Aggregator) addAction(ac equipment.ActionDescription) *equipment.ActionDescription {
	a.Calls = append(a.Calls, ac)
	return a.aggregate()
}
func (a *Aggregator) purge() []equipment.ActionDescription {
	ret := a.Calls[:]
	a.Calls = make([]equipment.ActionDescription, 0)
	return ret // TODO really purge
}

func (a *Aggregator) aggregate() *equipment.ActionDescription {
	if len(a.Calls) == 0 {
		return nil
	}
	var ret *equipment.ActionDescription
	ret = nil
	offPos := make(map[int]bool, 0)
	var last *equipment.ActionDescription
	last = nil
	for pos, ac := range a.Calls {
		if last == nil {
			last = &a.Calls[pos]
		} else {
			if last.Action == action.LH_MOVE || last.Action == action.LH_MOVE_EXPLICIT || last.Action == action.LH_MOVE_RAW {
				if ac.Action == action.LH_DISPENSE {
					//TODO check for blowout
					offPos[pos] = true
					offPos[pos-1] = true
					data := last.ActionData + " " + ac.ActionData
					newAc := equipment.ActionDescription{
						Action:     action.LH_DISPENSE,
						ActionData: data,
						Params:     nil,
					}
					ret = &newAc
					break
				} else if ac.Action == action.LH_ASPIRATE {
					offPos[pos] = true
					offPos[pos-1] = true

					//eval data as arrays when multi?
					data := fmt.Sprintf("Aspirate from %s/%s %s of %s.", last.Params["deckposition"], last.Params["wellcoords"], ac.Params["volume"], ac.Params["what"])
					newAc := equipment.ActionDescription{
						Action:     action.LH_ASPIRATE,
						ActionData: data,
						Params:     nil,
					}
					ret = &newAc
					break
				} else if ac.Action == action.LH_MOVE || ac.Action == action.LH_MOVE_EXPLICIT || ac.Action == action.LH_MOVE_RAW {
					offPos[pos-1] = true
					break
				}
			}
		}
	}
	if len(offPos) > 0 {
		newCalls := make([]equipment.ActionDescription, 0)
		for pos, a := range a.Calls {
			if ex, _ := offPos[pos]; !ex {
				newCalls = append(newCalls, a)
			}
		}
		a.Calls = newCalls
	}
	return ret
}

func stringToFile(str string) {
	f, err := os.OpenFile("./note.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}
