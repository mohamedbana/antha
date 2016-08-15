// /anthalib/simulator/liquidhandling/simulator.go: Part of the Antha language
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

package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
	"sort"
	"strings"
)

func summariseWell2Channel(well []string, channels []int) string {
	ret := make([]string, len(well))
	for i := range well {
		ret[i] = fmt.Sprintf("%s->channel%v", well[i], channels[i])
	}
	return strings.Join(ret, ", ")
}

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
	simulator.ErrorReporter
	state *RobotState
}

//Create a new VirtualLiquidHandler which mimics an LHDriver
func NewVirtualLiquidHandler(props *liquidhandling.LHProperties) *VirtualLiquidHandler {
	var vlh VirtualLiquidHandler

	vlh.validateProperties(props)
	//if the properties are that bad, don't bother building RobotState
	if vlh.HasError() {
		return &vlh
	}
	vlh.state = NewRobotState()

	//add the adaptors
	for _, head := range props.Heads {
		p := head.Adaptor.Params
		//9mm spacing currently hardcoded.
		//At some point we'll either need to fetch this from the driver or
		//infer it from the type of tipboxes/plates accepted
		spacing := wtype.Coordinates{0, 0, 0}
		if p.Orientation == wtype.LHVChannel {
			spacing.Y = 9.
		} else if p.Orientation == wtype.LHVChannel {
			spacing.X = 9.
		}
		vlh.state.AddAdaptor(NewAdaptorState(p.Independent, p.Multi, spacing))
	}

	is_in := func(s string, l []string) bool {
		for _, k := range l {
			if s == k {
				return true
			}
		}
		return false
	}
	//add the slots
	for name, coords := range props.Layout {
		//assuming default size for now
		sl := wtype.NewDeckSlot(name, coords, wtype.Coordinates{127.76, 85.48, 0.})
		if is_in(name, props.Tipwaste_preferences) {
			sl.SetAcceptsTipwaste(true)
		}
		if is_in(name, props.Tip_preferences) {
			sl.SetAcceptsTip(true)
		}
		if is_in(name, props.Input_preferences) || is_in(name, props.Output_preferences) {
			sl.SetAcceptsPlate(true)
		}
		vlh.state.AddSlot(sl)
	}

	return &vlh
}

// ------------------------------------------------------------------------------- Useful Utilities

func (self *VirtualLiquidHandler) validateProperties(props *liquidhandling.LHProperties) {

	//check a property
	check_prop := func(l []string, name string) {
		//is empty
		if len(l) == 0 {
			self.AddWarningf("NewVirtualLiquidHandler", "No %s specified", name)
		}
		//all locations defined
		for _, loc := range l {
			if _, ok := props.Layout[loc]; !ok {
				self.AddWarningf("NewVirtualLiquidHandler", "Undefined location \"%s\" referenced in %s", loc, name)
			}
		}
	}

	check_prop(props.Tip_preferences, "tip preferences")
	check_prop(props.Input_preferences, "input preferences")
	check_prop(props.Output_preferences, "output preferences")
	check_prop(props.Tipwaste_preferences, "tipwaste preferences")
	check_prop(props.Wash_preferences, "wash preferences")
	check_prop(props.Waste_preferences, "waste preferences")
}

func elements_equal(slice []string) bool {
	for i := range slice {
		if slice[i] != slice[0] {
			return false
		}
	}
	return true
}

//testSliceLength test that a bunch of slices are the correct length
func (self *VirtualLiquidHandler) testSliceLength(f_name string, slice_lengths map[string]int, exp_length int) bool {

	wrong := []string{}
	for name, actual_length := range slice_lengths {
		if actual_length != exp_length {
			wrong = append(wrong, name)
		}
	}

	if len(wrong) == 1 {
		self.AddErrorf(f_name, "Slice %s is not of expected length %v", wrong[0], exp_length)
		return false
	} else if len(wrong) > 1 {
		//for unit testing, names need to always be in the same order
		sort.Strings(wrong)
		self.AddErrorf(f_name, "Slices %s are not of expected length %v", strings.Join(wrong, ","), exp_length)
		return false
	}
	return true
}

//getAdaptorState
func (self *VirtualLiquidHandler) getAdaptorState(f_name string, h int) *AdaptorState {
	if h < 0 || h >= self.state.GetNumberOfAdaptors() {
		self.AddErrorf(f_name, "Unknown head %v", h)
		return nil
	}
	return self.state.GetAdaptor(h)
}

func (self *VirtualLiquidHandler) GetAdaptorState(head int) *AdaptorState {
	return self.state.GetAdaptor(head)
}

//testTipArgs check that load/unload tip arguments are valid insofar as they won't crash in RobotState
func (self *VirtualLiquidHandler) testTipArgs(f_name string, channels []int, head, multi int,
	platetype, position, well []string) bool {
	//head should exist
	adaptor := self.getAdaptorState(f_name, head)
	if adaptor == nil {
		return false
	}

	n_channels := adaptor.GetChannelCount()
	ret := true

	if multi != n_channels {
		self.AddErrorf(f_name, "Multi(=%v) doesn't match number of channels on Head%v(=%v)", multi, head, n_channels)
		ret = false
	}

	bad_channels := []string{}
	mchannels := map[int]bool{}
	dup_channels := []string{}
	for _, ch := range channels {
		if ch < 0 || ch >= n_channels {
			bad_channels = append(bad_channels, fmt.Sprintf("%v", ch))
		} else {
			if mchannels[ch] {
				dup_channels = append(dup_channels, fmt.Sprintf("%v", ch))
			} else {
				mchannels[ch] = true
			}
		}
	}
	if len(bad_channels) == 1 {
		self.AddErrorf(f_name, "Unknown channel \"%v\"", bad_channels[0])
		ret = false
	} else if len(bad_channels) > 1 {
		self.AddErrorf(f_name, "Unknown channels \"%v\"", strings.Join(bad_channels, "\",\""))
		ret = false
	}
	if len(dup_channels) == 1 {
		self.AddErrorf(f_name, "Channel%v appears more than once", dup_channels[0])
		ret = false
	} else if len(dup_channels) == 1 {
		self.AddErrorf(f_name, "Channels {%s} appear more than once", strings.Join(dup_channels, "\",\""))
		ret = false
	}

	ret = ret && self.testSliceLength(f_name, map[string]int{
		"platetype": len(platetype),
		"position":  len(position),
		"well":      len(well)},
		n_channels)
	return ret
}

func (self *VirtualLiquidHandler) GetObjectAt(slot_name string) wtype.LHObject {
	return self.state.GetSlot(slot_name).GetChild()
}

//getAbsolutePosition get a position within the liquidhandler, adding any errors as neccessary
//bool is false if the instruction shouldn't continue (e.g. missing deckposition e.t.c)
func (self *VirtualLiquidHandler) getAbsolutePosition(fname, deckposition, well string, reference int, platetype string) (wtype.Coordinates, bool) {
	ret := wtype.Coordinates{}

	var ref wtype.WellReference
	switch reference {
	case 0:
		ref = wtype.BottomReference
	case 1:
		ref = wtype.TopReference
	case 2:
		ref = wtype.LiquidReference
	default:
		self.AddErrorf(fname, "Invalid reference %d", reference)
		return ret, false
	}

	slot := self.state.GetSlot(deckposition)
	if slot == nil {
		self.AddErrorf(fname, "Unknown location \"%s\"", deckposition)
		return ret, false
	}

	target := slot.GetChild()
	if target == nil {
		self.AddErrorf(fname, "No object found at position %s", deckposition)
		return ret, false
	}
	if platetype != wtype.TypeOf(target) {
		self.AddWarningf(fname, "Object found at %s was type \"%s\" not type \"%s\" as expected",
			deckposition, wtype.TypeOf(target), platetype)
	}

	addr, ok := target.(wtype.Addressable)
	if !ok {
		self.AddErrorf(fname, "Object \"%s\" at \"%s\" is not addressable", wtype.NameOf(target), deckposition)
		return ret, false
	}

	wc := wtype.MakeWellCoords(well)
	if wc.IsZero() {
		self.AddErrorf(fname, "Couldn't parse well \"%s\"", well)
		return ret, false
	}

	if !addr.AddressExists(wc) {
		self.AddErrorf(fname, "Request for well %s in object \"%s\" at \"%s\" which is of size [%dx%d]",
			wc.FormatA1(), wtype.NameOf(target), deckposition, addr.NRows(), addr.NCols())
		return ret, false
	}

	ret, ok = addr.WellCoordsToCoords(wc, ref)
	if !ok {
		//since we already checked that the address exists, this must be a bad reference
		self.AddErrorf(fname, "Object type %s at %s doesn't support reference \"%s\"",
			wtype.TypeOf(target), deckposition, wtype.WellReferenceNames[ref])
		return ret, false
	}
	return ret, true
}

func makeOffsets(Xs, Ys, Zs []float64) []wtype.Coordinates {
	ret := make([]wtype.Coordinates, len(Xs))
	for i := range Xs {
		ret[i].X = Xs[i]
		ret[i].Y = Ys[i]
		ret[i].Z = Zs[i]
	}
	return ret
}

// ------------------------------------------------------------------------ ExtendedLHDriver

//Move command - used
func (self *VirtualLiquidHandler) Move(deckposition []string, wellcoords []string, reference []int,
	offsetX, offsetY, offsetZ []float64, platetype []string,
	head int) driver.CommandStatus {
	ret := driver.CommandStatus{true, driver.OK, "MOVE ACK"}

	//get the adaptor
	adaptor := self.getAdaptorState("Move", head)
	if adaptor == nil {
		return ret
	}

	//check slice length
	ok := self.testSliceLength("Move", map[string]int{
		"deckposition": len(deckposition),
		"wellcoords":   len(wellcoords),
		"reference":    len(reference),
		"offsetX":      len(offsetX),
		"offsetY":      len(offsetY),
		"offsetZ":      len(offsetZ),
		"plate_type":   len(platetype)},
		adaptor.GetChannelCount())
	if !ok {
		return ret
	}

	//find the coordinates of each explicitly requested position
	coords := make([]wtype.Coordinates, adaptor.GetChannelCount())
	offsets := makeOffsets(offsetX, offsetY, offsetZ)
	explicit := make([]bool, adaptor.GetChannelCount())
	exp_count := 0
	for i := range deckposition {
		if deckposition[i] == "" {
			if wellcoords[i] != "" {
				self.AddWarningf("Move", "deckposition was blank, but well was \"%s\"", wellcoords[i])
			}
			if platetype[i] != "" {
				self.AddWarningf("Move", "deckposition was blank, but platetype was \"%s\"", platetype[i])
			}
			explicit[i] = false
		} else {
			coords[i], ok = self.getAbsolutePosition("Move", deckposition[i], wellcoords[i], reference[i], platetype[i])
			if !ok {
				return ret
			}
			coords[i] = coords[i].Add(offsets[i])
			explicit[i] = true
			exp_count++
		}
	}
	if exp_count == 0 {
		self.AddWarning("Move", "Ignoring blank move command")
	}

	//find the head location, origin
	origin := wtype.Coordinates{}
	//for now, assuming that the relative position of the first explicitly provided channel and the head stay
	//the same. This seems sensible for the Glison, but might turn out not to be how other robots with independent channels work
	for i, c := range coords {
		if explicit[i] {
			origin = c.Subtract(adaptor.GetChannel(i).GetRelativePosition())
			break
		}
	}

	//fill in implicit locations
	for i := range coords {
		if !explicit[i] {
			coords[i] = origin.Add(adaptor.GetChannel(i).GetRelativePosition())
		}
	}

	//Get relative locations
	rel_coords := make([]wtype.Coordinates, adaptor.GetChannelCount())
	for i := range coords {
		rel_coords[i] = coords[i].Subtract(origin)
	}

	//check that the requested position is possible given the head/adaptor capabilities
	if !adaptor.IsIndependent() {
		//i.e. the channels can't move relative to each other or the head, so relative locations must remain the same
		moved := []string{}
		for i, rc := range rel_coords {
			if rc != adaptor.GetChannel(i).GetRelativePosition() {
				moved = append(moved, fmt.Sprintf("%d", i))
			}
		}
		if len(moved) > 0 {
			self.AddErrorf("Move", "Head %d Channels %s cannot move independently", head, strings.Join(moved, ","))
			return ret
		}
	}

	//check for collisions in the new location

	//update the head position accordingly
	adaptor.SetPosition(origin)
	for i, rc := range rel_coords {
		adaptor.GetChannel(i).SetRelativePosition(rc)
	}

	return ret
}

//Move raw - not yet implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
	self.AddWarning("MoveRaw", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "MOVERAW ACK"}
}

//Aspirate - used
func (self *VirtualLiquidHandler) Aspirate(volume []float64, overstroke []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {
	self.AddWarning("Aspirate", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "ASPIRATE ACK"}
}

//Dispense - used
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {
	self.AddWarning("Dispense", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "DISPENSE ACK"}
}

//LoadTips - used
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int,
	platetype, position, well []string) driver.CommandStatus {
	ret := driver.CommandStatus{true, driver.OK, "LOADTIPS ACK"}

	//check that RobotState won't crash
	if !self.testTipArgs("LoadTips", channels, head, multi, platetype, position, well) {
		return ret
	}

	adaptor := self.state.GetAdaptor(head)
	//if err := adaptor.LoadTips(platetype, position, well); err != nil {
	//    err.SetFunctionName("LoadTips")
	//    self.AddSimulationError(err)
	//}

	//Check that tips were loaded onto the specified channels
	mchannels := map[int]bool{}
	for _, ch := range channels {
		mchannels[ch] = true
	}
	extra := []string{}
	missing := []string{}
	for ch := 0; ch < adaptor.GetChannelCount(); ch++ {
		if ht, et := adaptor.GetChannel(ch).HasTip(), mchannels[ch]; ht && !et {
			extra = append(extra, fmt.Sprintf("%v", ch))
		} else if !ht && et {
			missing = append(missing, fmt.Sprintf("%v", ch))
		}
	}
	if len(extra) == 1 {
		self.AddErrorf("LoadTips", "Unexpected tip loaded to channel %s", extra[0])
	} else if len(extra) > 1 {
		self.AddErrorf("LoadTips", "Unexpected tips loaded to channels {%s}", strings.Join(extra, ","))
	}
	if len(missing) == 1 {
		self.AddErrorf("LoadTips", "Failed to load tip to channel %s", missing[0])
	} else if len(missing) > 1 {
		self.AddErrorf("LoadTips", "Failed to load tips to channels {%s}", strings.Join(missing, ","))
	}

	return ret
}

//UnloadTips - used
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int,
	platetype, position, well []string) driver.CommandStatus {
	ret := driver.CommandStatus{true, driver.OK, "UNLOADTIPS ACK"}

	//check that RobotState won't crash
	if !self.testTipArgs("UnloadTips", channels, head, multi, platetype, position, well) {
		return ret
	}

	//if err := self.state.GetAdaptor(head).UnloadTips(platetype, position, well); err != nil {
	//    self.AddError("UnloadTips", err.Error())
	//}

	//TODO: check that specified channels have no tips, any other tips remain

	return ret
}

//SetPipetteSpeed - used
func (self *VirtualLiquidHandler) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
	self.AddWarning("SetPipetteSpeed", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "SETPIPETTESPEED ACK"}
}

//SetDriveSpeed - used
func (self *VirtualLiquidHandler) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
	self.AddWarning("SetDriveSpeed", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "SETDRIVESPEED ACK"}
}

//Stop - unused
func (self *VirtualLiquidHandler) Stop() driver.CommandStatus {
	self.AddWarning("Stop", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "STOP ACK"}
}

//Go - unused
func (self *VirtualLiquidHandler) Go() driver.CommandStatus {
	self.AddWarning("Go", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "GO ACK"}
}

//Initialize - used
func (self *VirtualLiquidHandler) Initialize() driver.CommandStatus {
	err := self.state.Initialize()
	if err != nil {
		self.AddError("Initialize", err.Error())
	}
	return driver.CommandStatus{true, driver.OK, "INITIALIZE ACK"}
}

//Finalize - used
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
	err := self.state.Finalize()
	if err != nil {
		self.AddError("Finalize", err.Error())
	}
	return driver.CommandStatus{true, driver.OK, "FINALIZE ACK"}
}

//Wait - used
func (self *VirtualLiquidHandler) Wait(time float64) driver.CommandStatus {
	self.AddWarning("Wait", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "WAIT ACK"}
}

//Mix - used
func (self *VirtualLiquidHandler) Mix(head int, volume []float64, platetype []string, cycles []int,
	multi int, what []string, blowout []bool) driver.CommandStatus {
	self.AddWarning("Mix", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "MIX ACK"}
}

//ResetPistons - used
func (self *VirtualLiquidHandler) ResetPistons(head, channel int) driver.CommandStatus {
	self.AddWarning("ResetPistons", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "RESETPISTONS ACK"}
}

//AddPlateTo - used
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {

	ret := driver.CommandStatus{true, driver.OK, "ADDPLATETO ACK"}

	if obj, ok := plate.(wtype.LHObject); ok {
		if n, nok := obj.(wtype.Named); nok && n.GetName() != name {
			self.AddWarningf("AddPlateTo", "Object name(=%s) doesn't match argument name(=%s)", n.GetName(), name)
		}

		if err := self.state.AddObject(position, obj); err != nil {
			err.SetFunctionName("AddPlateTo")
			self.AddSimulationError(err)
		}

	} else {
		self.AddErrorf("AddPlateTo", "Couldn't add object of type %T to %s", plate, position)
	}

	return ret
}

//RemoveAllPlates - used
func (self *VirtualLiquidHandler) RemoveAllPlates() driver.CommandStatus {
	self.AddWarning("RemoveAllPlates", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "REMOVEALLPLATES ACK"}
}

//RemovePlateAt - unused
func (self *VirtualLiquidHandler) RemovePlateAt(position string) driver.CommandStatus {
	self.AddWarning("RemovePlateAt", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "REMOVEPLATEAT ACK"}
}

//SetPositionState - unused
func (self *VirtualLiquidHandler) SetPositionState(position string, state driver.PositionState) driver.CommandStatus {
	self.AddWarning("SetPositionState", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "SETPOSITIONSTATE ACK"}
}

//GetCapabilites - used
func (self *VirtualLiquidHandler) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
	self.AddWarning("SetPositionState", "Not yet implemented")
	return liquidhandling.LHProperties{}, driver.CommandStatus{true, driver.OK, "GETCAPABILITIES ACK"}
}

//GetCurrentPosition - unused
func (self *VirtualLiquidHandler) GetCurrentPosition(head int) (string, driver.CommandStatus) {
	self.AddWarning("GetCurrentPosition", "Not yet implemented")
	return "", driver.CommandStatus{true, driver.OK, "GETCURRNETPOSITION ACK"}
}

//GetPositionState - unused
func (self *VirtualLiquidHandler) GetPositionState(position string) (string, driver.CommandStatus) {
	self.AddWarning("GetPositionState", "Not yet implemented")
	return "", driver.CommandStatus{true, driver.OK, "GETPOSITIONSTATE ACK"}
}

//GetHeadState - unused
func (self *VirtualLiquidHandler) GetHeadState(head int) (string, driver.CommandStatus) {
	self.AddWarning("GetHeadState", "Not yet implemented")
	return "I'm fine thanks, how are you?", driver.CommandStatus{true, driver.OK, "GETHEADSTATE ACK"}
}

//GetStatus - unused
func (self *VirtualLiquidHandler) GetStatus() (driver.Status, driver.CommandStatus) {
	self.AddWarning("GetStatus", "Not yet implemented")
	return driver.Status{}, driver.CommandStatus{true, driver.OK, "GETSTATUS ACK"}
}

//UpdateMetaData - used
func (self *VirtualLiquidHandler) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus {
	self.AddWarning("UpdateMetaData", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "UPDATEMETADATA ACK"}
}

//UnloadHead - unused
func (self *VirtualLiquidHandler) UnloadHead(param int) driver.CommandStatus {
	self.AddWarning("UnloadHead", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "UNLOADHEAD ACK"}
}

//LoadHead - unused
func (self *VirtualLiquidHandler) LoadHead(param int) driver.CommandStatus {
	self.AddWarning("LoadHead", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "LOADHEAD ACK"}
}

//Lights On - not implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOn() driver.CommandStatus {
	return driver.CommandStatus{true, driver.OK, "LIGHTSON ACK"}
}

//Lights Off - notimplemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOff() driver.CommandStatus {
	return driver.CommandStatus{true, driver.OK, "LIGHTSOFF ACK"}
}

//LoadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) LoadAdaptor(param int) driver.CommandStatus {
	self.AddWarning("LoadAdaptor", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "LOADADAPTOR ACK"}
}

//UnloadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) UnloadAdaptor(param int) driver.CommandStatus {
	self.AddWarning("UnloadAdaptor", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "UNLOADADAPTOR ACK"}
}

//Open - notimplemented in CRI
func (self *VirtualLiquidHandler) Open() driver.CommandStatus {
	self.AddWarning("Open", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "OPEN ACK"}
}

//Close - notimplement in CRI
func (self *VirtualLiquidHandler) Close() driver.CommandStatus {
	self.AddWarning("Close", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "CLOSE ACK"}
}

//Message - unused
func (self *VirtualLiquidHandler) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
	self.AddWarning("Message", "Not yet implemented")
	return driver.CommandStatus{true, driver.OK, "MESSAGE ACK"}
}

//GetOutputFile - used, but not in instruction stream
func (self *VirtualLiquidHandler) GetOutputFile() (string, driver.CommandStatus) {
	self.AddWarning("GetOutputFile", "Not yet implemented")
	return "You forgot to say 'please'", driver.CommandStatus{true, driver.OK, "GETOUTPUTFILE ACK"}
}
