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
	"errors"

	"github.com/antha-lang/antha/antha/anthalib/driver"
	"github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipmentManager"
)
//ManualDriver represents a piece of equipment to substitute a human performing the different actions possible in the
// antha ecosystem. I can substitute a piece of equipment or help with the usage of a non-connected device
type ManualDriver struct {
	//eq is the equipment instance to which they send the actions
	eq equipment.Equipment
	//ag aggregates the different calls into groups with meaning to humans and ignores robot specific actions
	ag Aggregator
	//plateLookup lookup dictionary to give fancier names to plates
	plateLookup	TranslateDictionary
	//tipboxlookup lookup dictionary to give fancier names to tipbox
	tipboxlookup TranslateDictionary
	//tipwastelookup lookup dictionary to give fancier names to tipwaste
	tipwastelookup TranslateDictionary
}
//TranslateDictionary keeps a dictionary of numerated items of the type prefix and gives names and holds a lookup table
// the names are `prefix num count`
type TranslateDictionary struct {
	dictionary map[string]string
	prefix		string
	count		int
}
//NewTranslateDictionary creates a new dictionary with the given prefix
func NewTranslateDictionary(prefix string) *TranslateDictionary {
	ret := new(TranslateDictionary)
	ret.count = 1
	ret.prefix = prefix
	ret.dictionary = make(map[string]string)
	return ret
}
//lookupName gives the given reference to a name in the dictionary. Gives a new one when it does not exist and stores it
func (t *TranslateDictionary) lookupName(ref string) string {
	name, ok := t.dictionary[ref]
	if !ok {
		t.dictionary[ref] = fmt.Sprintf("%s num %d", t.prefix, t.count)
		t.count = t.count + 1
		name = t.dictionary[ref]
	}
	return name
}
//lookupPlateName looks the name given to a plate in the underlying plate dictionary
func (m *ManualDriver)lookupPlateName(ref string) string {
	return m.plateLookup.lookupName(ref)
}
//lookupTipboxName looks the name given to a plate in the underlying tipbox dictionary
func (m *ManualDriver)lookupTipboxName(ref string) string {
	return m.tipboxlookup.lookupName(ref)
}
//lookupTipwasteName looks the name given to a plate in the underlying tipwaste dictionary
func (m *ManualDriver)lookupTipwasteName(ref string) string {
	return m.tipwastelookup.lookupName(ref)
}
//sendActionToEquipment sends a specific actionDescription to the underlying equipment driver. Can have a synchronous
// or asynchronous operation
func (m *ManualDriver) sendActionToEquipment(ac equipment.ActionDescription) error {
	r := m.ag.addAction(ac)
	if r != nil {
		for _, a := range r {
			go m.eq.Do(a)
		}
	}
	return nil
}

//NewManualDriver returns a new instance of a manual driver pointing to the right piece of equipment
func NewManualDriver() *ManualDriver {
	ret := new(ManualDriver)
	ret.ag = *NewAggregator()
	ret.plateLookup = *NewTranslateDictionary("Plate")
	ret.tipboxlookup = *NewTranslateDictionary("Tipbox")
	ret.tipwastelookup = *NewTranslateDictionary("Tipwaste")
	eqm := *equipmentManager.GetEquipmentManager()
	params := make(map[string]string,0)
	ret.eq = *eqm.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", params))
	return ret
}
func (m *ManualDriver) Move(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []string, head int) driver.CommandStatus {
	params := make(map[string]string)
	params["deckposition"] = fmt.Sprintf("%v", m.lookupPlateName(deckposition[0]))
	params["wellcoords"] = fmt.Sprintf("%v", wellcoords[0])
	params["reference"] = fmt.Sprintf("%v", reference)
	params["offsetX"] = fmt.Sprintf("%v", offsetX)
	params["offsetY"] = fmt.Sprintf("%v", offsetY)
	params["offsetZ"] = fmt.Sprintf("%v", offsetZ)
	params["plate_type"] = fmt.Sprintf("%v", plate_type)
	params["head"] = fmt.Sprintf("%v", head)

	desc := fmt.Sprintf("Deck Postition %v @well %v with reference %v", m.lookupPlateName(deckposition[0]), wellcoords, reference)
	ad := *equipment.NewActionDescription(action.LH_MOVE, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) MoveExplicit(deckposition []string, wellcoords []string, reference []int, offsetX, offsetY, offsetZ []float64, plate_type []*wtype.LHPlate, head int) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Deck Postition %v @well %v with reference %v", deckposition, wellcoords, reference)
	ad := *equipment.NewActionDescription(action.LH_MOVE_EXPLICIT, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
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
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Aspirate(volume []float64, overstroke []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	params := make(map[string]string)
	params["volume"] = fmt.Sprintf("%g", volume[0])
	params["overstroke"] = fmt.Sprintf("%v", overstroke[0])
	params["head"] = fmt.Sprintf("%v", head)
	params["multi"] = fmt.Sprintf("%v", multi)
	params["platetype"] = fmt.Sprintf("%v", platetype[0])
	params["what"] = fmt.Sprintf("%v", what[0])
//	params["llf"] = fmt.Sprintf("%v", llf[0])

	desc := fmt.Sprintf("Aspirate volumes %v", volume) //TOOD make a meaning of the values
	ad := *equipment.NewActionDescription(action.LH_ASPIRATE, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}

}
func (m *ManualDriver) Dispense(volume []float64, blowout []bool, head int, multi int, platetype []string, what []string, llf []bool) driver.CommandStatus {
	params := make(map[string]string)
	params["volume"] = fmt.Sprintf("%g", volume[0])
	params["blowout"] = fmt.Sprintf("%t", blowout[0])
	params["head"] = fmt.Sprintf("%v", head)
	params["multi"] = fmt.Sprintf("%v", multi)
	params["platetype"] = fmt.Sprintf("%v", platetype[0])
	params["what"] = fmt.Sprintf("%s", what[0])
	params["llf"] = fmt.Sprintf("%t", llf)

	desc := fmt.Sprintf("Dispense volumes %v", volume)
	ad := *equipment.NewActionDescription(action.LH_DISPENSE, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
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
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
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
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
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
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
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
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
//doesNotApply will return an OK driver.CommandStatus, it should be used in actions that are not implemented in this
// driver as a means to annotate this
func (m *ManualDriver) doesNotApply() driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) Stop() driver.CommandStatus {
	return m.doesNotApply()
}
func (m *ManualDriver) Go() driver.CommandStatus {
	return m.doesNotApply()
}
func (m *ManualDriver) Initialize() driver.CommandStatus {
	return m.doesNotApply()
}
func (m *ManualDriver) Finalize() driver.CommandStatus {
	return m.doesNotApply()
}
func (m *ManualDriver) SetPositionState(position string, state driver.PositionState) driver.CommandStatus {
	desc := fmt.Sprintf("Set position %s to state %v", position, state)
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_SET_POSITION_STATE, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}

//TODO implement changing capabilities for manual driver
func (m *ManualDriver) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	return liquidhandling.LHProperties{}, driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) GetCurrentPosition(head int) (string, driver.CommandStatus) {
	return "", m.doesNotApply()
}
func (m *ManualDriver) GetPositionState(position string) (string, driver.CommandStatus) {
	return "", m.doesNotApply()
}
func (m *ManualDriver) GetHeadState(head int) (string, driver.CommandStatus) {
	return "", m.doesNotApply()
}
func (m *ManualDriver) GetStatus() (driver.Status, driver.CommandStatus) {
	st := new(driver.Status)
	return *st, m.doesNotApply()
}
func (m *ManualDriver) ResetPistons(head, channel int) driver.CommandStatus {
	return m.doesNotApply()
}
func (m *ManualDriver) Wait(time float64) driver.CommandStatus {
	return m.doesNotApply()
}
func (m *ManualDriver) Mix(head int, volume []float64, fvolume []float64, platetype []string, cycles []int, multi int, prms map[string]interface{}) driver.CommandStatus {
	params := make(map[string]string)

	desc := fmt.Sprintf("Mix with head num %d", head)
	ad := *equipment.NewActionDescription(action.LH_MIX, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}

func (m *ManualDriver) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
	//plate can be lhplate, lhtipbox or lhtipwaste
	params := make(map[string]string)
	var desc string
	var ad equipment.ActionDescription
	switch w := plate.(type){
	case *wtype.LHPlate:
		desc = fmt.Sprintf("Get a %s of type %s and label it as %s. We will refer to it by that name", "Plate", w.Type, m.lookupPlateName(w.GetName()))
	case wtype.LHPlate:
		desc = fmt.Sprintf("Get a %s of type %s and label it as %s. We will refer to it by that name", "Plate", w.Type, m.lookupPlateName(w.GetName()))
	case *wtype.LHTipbox:
		desc = fmt.Sprintf("Get a %s of type %s and label it as %s. We will refer to it by that name", "Tipbox", w.Type, m.lookupTipboxName(w.GetName()))
	case wtype.LHTipbox:
		desc = fmt.Sprintf("Get a %s of type %s and label it as %s. We will refer to it by that name", "Tipbox", w.Type, m.lookupTipboxName(w.GetName()))
	case *wtype.LHTipwaste:
		desc = fmt.Sprintf("Get a %s of type %s and label it as %s. We will refer to it by that name", "Tipwaste", w.Type, m.lookupTipwasteName(w.GetName()))
	case wtype.LHTipwaste:
		desc = fmt.Sprintf("Get a %s of type %s and label it as %s. We will refer to it by that name", "Tipwaste", w.Type, m.lookupTipwasteName(w.GetName()))
	default:
		//panic(reflect.TypeOf(plate))
		panic(errors.New("Unknow plate type"))
	}

	ad = *equipment.NewActionDescription(action.LH_ADD_PLATE, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) RemoveAllPlates() driver.CommandStatus {
	desc := fmt.Sprintf("Remove all Plates from the liquid handler") //TODO
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_REMOVE_ALL_PLATES, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
func (m *ManualDriver) RemovePlateAt(position string) driver.CommandStatus {
	desc := fmt.Sprintf("Remove Plate at position %s", position) //TODO
	params := make(map[string]string)

	ad := *equipment.NewActionDescription(action.LH_REMOVE_PLATE, desc, params)
	err := m.sendActionToEquipment(ad)
	if err != nil {
		return driver.CommandStatus{
			OK:        false,
			Errorcode: 1,
			Msg:       err.Error(),
		}
	}
	return driver.CommandStatus{
		OK:        true,
		Errorcode: 0,
		Msg:       "OK",
	}
}
//Aggregator has a list of actions in which it tries to stablish relations to make them more human friendly
type Aggregator struct {
	//Calls is the list of actionDescriptions
	Calls	[]equipment.ActionDescription
}
//NewAggregator instantiates a new aggregator
func NewAggregator() *Aggregator {
	ret := new(Aggregator)
	ret.Calls = make([]equipment.ActionDescription, 0)
	return ret
}
//AddAction appends a new action to our Call action and tries to stablish some kind of aggregation between the already
// present actions and this new one. If the action could be aggregated with future actions it is left in our call list
// it is returned otherwise. The addition of an action could result in more than one action being returned
func (a *Aggregator) addAction(ac equipment.ActionDescription) []equipment.ActionDescription {
	a.Calls = append(a.Calls, ac)
	return a.aggregate()
}
//purge returns all the remaining actions and reinits the aggregator
func (a *Aggregator) purge() []equipment.ActionDescription {
	ret := a.Calls[:]
	a.Calls = make([]equipment.ActionDescription,0)
	return ret
}
//aggregate tries to merge and return actions.
func (a *Aggregator) aggregate() (ret []equipment.ActionDescription) {
	out := "initial values "
	for k, v := range a.Calls {
		out = out + fmt.Sprintf("%d - %s, ",k, v.Action)
	}
	if len(a.Calls) > 3 {
		panic(errors.New("Unexpected Behaviour"))
	}
	if len(a.Calls) == 0 {
		return nil
	} else if len(a.Calls) == 1 {
		if !actionIsMove(a.Calls[0].Action) && a.Calls[0].Action != action.LH_DISPENSE && !actionIsSetup(a.Calls[0].Action) && !actionIsHandleTips(a.Calls[0].Action){
			//return the action immediately
			ret = append(ret, a.Calls[0])
			a.Calls = make([]equipment.ActionDescription,0)
			return
		}
	} else if len(a.Calls) == 2 { //more than 1 action, we check for aggregates, we should not have more than 2 actually
		first := a.Calls[0]
		second := a.Calls[1]
		//cases are
		// move + move = move (preserve last)
		// move + asp = asp with pos
		// move + disp = disp with pos
		// move = disp(blowout) = nothing
		// disp + move = disp + move
		// disp + anything = disp out + check anything
		// move + other = return other, clean queue
		if actionIsMove(first.Action) && actionIsMove(second.Action) {
			a.Calls = a.Calls[1:]
			return nil// to be aggregated
		} else if actionIsMove(first.Action) && second.Action == action.LH_ASPIRATE {
			data := fmt.Sprintf("Aspirate from %s well %s, %sul. of %s.", first.Params["deckposition"], first.Params["wellcoords"], second.Params["volume"], second.Params["what"])
			newAc := equipment.ActionDescription{
				Action: action.LH_ASPIRATE,
				ActionData: data,
				Params: second.Params,
			}
			a.Calls = make([]equipment.ActionDescription, 0) //clean the queue
			ret = append(ret, newAc)
			return
		} else if actionIsMove(first.Action) && second.Action == action.LH_DISPENSE && second.Params["blowout"] == "true" {
			a.Calls = make([]equipment.ActionDescription, 0) //clean the queue
			return nil //command mean nothing for a human
		} else if actionIsMove(first.Action) && second.Action == action.LH_DISPENSE && second.Params["blowout"] != "true" {
			data := fmt.Sprintf("Dispense in %s well %s, %sul. of %s.", first.Params["deckposition"], first.Params["wellcoords"], second.Params["volume"], second.Params["what"])
			newAc := equipment.ActionDescription{
				Action: action.LH_DISPENSE,
				ActionData: data,
				Params: first.Params,
			}
			a.Calls = make([]equipment.ActionDescription, 0) //clean the queue
			a.Calls = append(a.Calls, newAc)
			return nil // we can have moves after the dispense
		} else if first.Action == action.LH_DISPENSE && actionIsMove(second.Action) {
			if second.Params["reference"] == "0" {
				data := fmt.Sprintf("%s. Touch off after dispense", first.ActionData)
				newAc := equipment.ActionDescription{
					Action: action.LH_DISPENSE,
					ActionData: data,
					Params: first.Params,
				}
				a.Calls = make([]equipment.ActionDescription, 0) //clean the queue
				ret = append(ret, newAc)
				return
			}
		} else if first.Action == action.LH_DISPENSE { //obviously second action is not a move
			ret = append(ret, first)
			a.Calls = a.Calls[1:]
			ret = append(ret, a.aggregate()...)
			return
		}else if actionIsMove(first.Action) { //Move + unknown, ignore move
			a.Calls = a.Calls[1:]
			return a.aggregate()
		} else if actionIsSetup(first.Action) && (actionIsSetup(second.Action) || actionIsHandleTips(second.Action)) {
			newDesc := fmt.Sprintf("%s\n%s", first.ActionData, second.ActionData)
			newAc := equipment.NewActionDescription(action.LH_SETUP, newDesc, first.Params)
			a.Calls = make([]equipment.ActionDescription, 0)
			a.Calls = append(a.Calls, *newAc)
			return
		} else if actionIsSetup(first.Action) && actionIsMove(second.Action) {
			return
		} else if actionIsSetup(first.Action) && !actionIsSetup(second.Action) { //second is not handle tips
			ret = append(ret, first)
			a.Calls = a.Calls[1:]
			return
		} else if actionIsHandleTips(first.Action) && actionIsHandleTips(second.Action) {
			//unload and then load
			newDesc := fmt.Sprintf("Change tips with new tip of type %s", second.Params["platetype"])
			newAc := equipment.NewActionDescription(action.MLH_CHANGE_TIPS, newDesc, first.Params)
			a.Calls = make([]equipment.ActionDescription, 0)
			ret = append(ret, *newAc)
			return
		} else if actionIsHandleTips(first.Action) && actionIsMove(second.Action) {
			return
		} else if actionIsHandleTips(first.Action) && !actionIsHandleTips(second.Action) {
			ret = append(ret, first)
			a.Calls = a.Calls[1:]
			ret = append(ret, a.aggregate()...)
			return
		}
	} else {
		//cases to get here are
		// disp + move + move = disp + move (merge the two moves and ignore the disp until a ref 0 is found
		// disp + move + otherthing = output disp, deal with the others
		// setup + move + tips = setup + tips
		// tips + move + tips = tips + tips
		// setup + move + otherthing = output tips + recursive
		first	:= a.Calls[0]
		second	:= a.Calls[1]
		third	:= a.Calls[2]
		if first.Action == action.LH_DISPENSE && !actionIsMove(third.Action) { //spit the dispense and merge the other two
			ret = append(ret, first)
			a.Calls = a.Calls[1:]
			ret = append(ret, a.aggregate()...)
			return
		} else if first.Action == action.LH_DISPENSE && actionIsMove(second.Action) && actionIsMove(third.Action) { //disp + two moves
			if third.Params["reference"] == "0" {
				//take away the man in the middle // we are basically ignoring the first move, and calling recursively
				a.Calls = append(a.Calls[0:1], third)
				ret = a.aggregate()
				return
			} else { //normal two move merge, ignore the second
				a.Calls = append(a.Calls[0:1], third)
				return
			}
		} else if actionIsSetup(first.Action) && actionIsMove(second.Action) && actionIsHandleTips(third.Action){
			a.Calls = append(a.Calls[0:1], a.Calls[2])
			return a.aggregate()
		} else if actionIsSetup(first.Action) && actionIsMove(second.Action) { //third is anything
			ret = append(ret, first)
			a.Calls = a.Calls[1:]
			ret = append(ret, a.aggregate()...)
			return
		} else if actionIsHandleTips(first.Action) && actionIsMove(second.Action) && actionIsHandleTips(third.Action) {
			a.Calls = append(a.Calls[0:1], a.Calls[2])
			return a.aggregate()
		}
		actionStack := fmt.Sprintf("first %s, second %s, third %s", first.Action, second.Action, third.Action)
		panic(errors.New("Unhandled situation" + actionStack))
	}
	return
}
//actionIsMove returns true when an action that represents a movement is fed as an argument
func actionIsMove(ac action.Action) bool {
	return actionInActionList(ac, []action.Action{action.LH_MOVE, action.LH_MOVE_RAW, action.LH_MOVE_EXPLICIT})
}
//actionIsSetup returns true when an action that represents a movement is fed as an argument
func actionIsSetup(ac action.Action) bool {
	return actionInActionList(ac, []action.Action{action.LH_REMOVE_ALL_PLATES, action.LH_REMOVE_PLATE, action.LH_ADD_PLATE, action.LH_SETUP})
}
//actionIsTips returns true when an action depicts some kind of tip handling
func actionIsHandleTips(ac action.Action) bool {
	return actionInActionList(ac, []action.Action{action.LH_LOAD_TIPS, action.LH_UNLOAD_TIPS})
}
//actionInActionList checks if action is contained in the given slice of actions, returns true if so
func actionInActionList(ac action.Action, list []action.Action) bool {
	for _, v := range list {
		if v.String() == ac.String() {
			return true
		}
	}
	return false
}
