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
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
	"math"
	"sort"
	"strings"
)

func pTips(N int) string {
	if N == 1 {
		return "tip"
	}
	return "tips"
}

func summariseWell2Channel(well []string, channels []int) string {
	ret := make([]string, 0, len(channels))
	for ch := range channels {
		ret = append(ret, fmt.Sprintf("%s->channel%v", well[ch], ch))
	}
	return strings.Join(ret, ", ")
}

func summariseChannels(channels []int) string {
	if len(channels) == 1 {
		return fmt.Sprintf("channel %d", channels[0])
	}
	sch := make([]string, 0, len(channels))
	for _, ch := range channels {
		sch = append(sch, fmt.Sprintf("%d", ch))
	}
	return fmt.Sprintf("channels %s", strings.Join(sch, ","))
}

func summariseVolumes(vols []float64) string {
	equal := true
	for _, v := range vols {
		if v != vols[0] {
			equal = false
			break
		}
	}

	if equal {
		return wunit.NewVolume(vols[0], "ul").String()
	}

	s_vols := make([]string, len(vols))
	for i, v := range vols {
		s_vols[i] = wunit.NewVolume(v, "ul").String()
		s_vols[i] = s_vols[i][:len(s_vols[i])-2]
	}
	return fmt.Sprintf("{%s}ul", strings.Join(s_vols, ","))
}

func summariseStrings(s []string) string {
	m := map[string]bool{}
	for _, v := range s {
		if v != "" {
			m[v] = true
		}
	}
	s2 := make([]string, 0, len(m))
	for k, _ := range m {
		s2 = append(s2, k)
	}

	if len(s2) == 0 {
		return "nil"
	}
	if len(s2) == 1 {
		return s2[0]
	}
	return fmt.Sprintf("{%s}", strings.Join(s2, ","))
}

type adaptorCollision struct {
	channel int
	objects []wtype.LHObject
}

type adaptorCollisions []adaptorCollision

func (self adaptorCollisions) String() string {
	channels := make([]string, 0, len(self))
	objects := []wtype.LHObject{}

	seen := func(o wtype.LHObject) bool {
		for _, O := range objects {
			if o == O {
				return true
			}
		}
		return false
	}

	for _, ac := range self {
		channels = append(channels, fmt.Sprintf("%d", ac.channel))
		for _, obj := range ac.objects {
			if !seen(obj) {
				objects = append(objects, obj)
			}
		}
	}

	s_obj := make([]string, 0, len(objects))
	for _, o := range objects {
		s_obj = append(s_obj, fmt.Sprintf("%s \"%s\"", wtype.ClassOf(o), wtype.NameOf(o)))
	}

	if len(self) == 1 {
		return fmt.Sprintf("channel %s collides with %s", channels[0], strings.Join(s_obj, " and "))
	}
	return fmt.Sprintf("channels %s collide with %s", strings.Join(channels, ","), strings.Join(s_obj, " and "))
}

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
	simulator.ErrorReporter
	state             *RobotState
	tipbox_collisions bool
}

//Create a new VirtualLiquidHandler which mimics an LHDriver
func NewVirtualLiquidHandler(props *liquidhandling.LHProperties) *VirtualLiquidHandler {
	var vlh VirtualLiquidHandler

	vlh.tipbox_collisions = true

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

	//Make the deck
	deck := wtype.NewLHDeck("simulated deck", props.Mnfr, props.Model)
	for name, pos := range props.Layout {
		//size not given un LHProperties, assuming standard 96well size
		deck.AddSlot(name, pos, wtype.Coordinates{127.76, 85.48, 0})
		//deck.SetSlotAccepts(name, "riser")
	}

	for _, name := range props.Tip_preferences {
		deck.SetSlotAccepts(name, "tipbox")
	}
	for _, name := range props.Input_preferences {
		deck.SetSlotAccepts(name, "plate")
	}
	for _, name := range props.Output_preferences {
		deck.SetSlotAccepts(name, "plate")
	}
	for _, name := range props.Tipwaste_preferences {
		deck.SetSlotAccepts(name, "tipwaste")
	}

	vlh.state.SetDeck(deck)

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

func contains(v int, s []int) bool {
	for _, val := range s {
		if v == val {
			return true
		}
	}
	return false
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

func (self *VirtualLiquidHandler) GetObjectAt(slot string) wtype.LHObject {
	ret, _ := self.state.GetDeck().GetChild(slot)
	return ret
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

	if ret {
		for i := range platetype {
			if contains(i, channels) {
				if platetype[i] == "" && well[i] == "" && position[i] == "" {
					self.AddErrorf(f_name, "Command given for channel %d, but no platetype, well or position given", i)
					return false
				}
			} else {
				if !(platetype[i] == "" && well[i] == "" && position[i] == "") {
					self.AddWarningf(f_name, "No command for channel %d, but platetype, well or position given", i)
				}
			}
		}
	}
	return ret
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

	target, ok := self.state.GetDeck().GetChild(deckposition)
	if !ok {
		self.AddErrorf(fname, "Unknown location \"%s\"", deckposition)
		return ret, false
	}
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

//return the unique elements of a string slice
func getUnique(ss []string) []string {
	r := []string{}
	is_in := func(s string) bool {
		for _, v := range r {
			if v == s {
				return true
			}
		}
		return false
	}

	for _, s := range ss {
		if !is_in(s) && s != "" {
			r = append(r, s)
		}
	}
	return r
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
			//if there's a tip, take account of it
			if tip := adaptor.GetChannel(i).GetTip(); tip != nil {
				coords[i] = coords[i].Add(wtype.Coordinates{0., 0., tip.GetSize().Z})
			}
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
			//get slice of well coords
			wc := make([]wtype.WellCoords, len(wellcoords))
			for i := range wellcoords {
				wc[i] = wtype.MakeWellCoords(wellcoords[i])
			}
			self.AddErrorf("Move", "Non-independent head '%d' can't move adaptors to \"%s\" positions %s, layout mismatch",
				head, strings.Join(getUnique(platetype), "\",\""), wtype.HumanizeWellCoords(wc))
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

	ret := driver.CommandStatus{true, driver.OK, "ASPIRATE ACK"}

	if !self.testSliceLength("Aspirate", map[string]int{
		"volume":     len(volume),
		"overstroke": len(overstroke),
		"platetype":  len(platetype),
		"what":       len(what),
		"llf":        len(llf),
	}, multi) {
		return ret
	}

	//which channels are specified
	channels := make([]bool, multi)
	s_chan := make([]int, 0)
	s_vols := make([]float64, 0)
	negative := false
	for i := range channels {
		channels[i] = !(platetype[i] == "" && what[i] == "")
		if channels[i] {
			s_chan = append(s_chan, i)
			s_vols = append(s_vols, volume[i])
			negative = negative || volume[i] < 0.
		}
	}

	describe := func() string {
		return fmt.Sprintf("aspirating %s of %s to head %d %s",
			summariseVolumes(s_vols), summariseStrings(what), head, summariseChannels(s_chan))
	}

	if negative {
		self.AddErrorf("Aspirate", "Error %s - cannot aspirate negative volume")
		return ret
	}

	//Verify arguments
	if head < 0 || head >= self.state.GetNumberOfAdaptors() {
		self.AddErrorf("Aspirate", "While %s - unknown head %d", describe(), head)
		return ret
	}
	adaptor := self.state.GetAdaptor(head)
	if multi != adaptor.GetChannelCount() {
		self.AddErrorf("Aspirate", "Invalid arguments - multi[=%d] must equal number of channels on head %d [=%d]", multi, head, adaptor.GetChannelCount())
		return ret
	}

	deck := self.state.GetDeck()

	//get the position of tips
	tip_pos := make([]wtype.Coordinates, multi)
	wells := make([]*wtype.LHWell, multi)
	tip_missing := []int{}
	for i, explicit := range channels {
		ch := adaptor.GetChannel(i)
		if ch.HasTip() {
			tip_pos[i] = ch.GetAbsolutePosition().Subtract(wtype.Coordinates{0., 0., ch.GetTip().GetSize().Z})

			objs := deck.GetPointIntersections(tip_pos[i])
			for _, o := range objs {
				if well, ok := o.(*wtype.LHWell); ok {
					wells[i] = well
				}
			}
		} else if explicit {
			tip_missing = append(tip_missing, i)
		}
	}
	if len(tip_missing) > 0 {
		self.AddErrorf("Aspirate", "While %s - missing %s on %s", describe(), pTips(len(tip_missing)), summariseChannels(tip_missing))
		return ret
	}

	//independence constraints
	if !adaptor.IsIndependent() {
		different := false
		v := -1.
		for ch, b := range channels {
			if b {
				if v >= 0 {
					different = different || v != volume[ch]
				} else {
					v = volume[ch]
				}
			} else if wells[ch] != nil {
				//a non-explicitly requested tip is in a well. If the well has stuff in it, it'll get aspirated
				if c := wells[ch].Contents(); !c.IsZero() {
					self.AddErrorf("Aspirate",
						"While %s - channel %d will inadvertantly aspirate %s from well %s as head is not independent",
						describe(), ch, c.GetType(), wells[ch].GetName())
				}
			}
		}
		if different {
			self.AddErrorf("Aspirate", "While %s - channels cannot aspirate different volumes in non-independent head", describe())
			return ret
		}
	}

	//check liquid type
	for i := range channels {
		if !channels[i] || wells[i] == nil {
			continue
		}
		if wells[i].Contents().GetType() != what[i] {
			self.AddWarningf("Aspirate", "While %s - well %s contains %s, not %s",
				describe(), wells[i].GetName(), wells[i].Contents().GetType(), what[i])
		}
	}

	//move liquid
	no_well := []int{}
	for i := range channels {
		if !channels[i] {
			continue
		}
		v := wunit.NewVolume(volume[i], "ul")
		tip := adaptor.GetChannel(i).GetTip()
		fv := tip.CurrentVolume()
		fv.Add(v)

		if wells[i] == nil {
			no_well = append(no_well, i)
		} else if wells[i].WorkingVolume().LessThan(v) {
			self.AddErrorf("Aspirate", "While %s - well %s only contains %s working volume",
				describe(), wells[i].GetName(), wells[i].WorkingVolume())
		} else if fv.GreaterThan(tip.MaxVol) {
			self.AddErrorf("Aspirate", "While %s - channel %d contains %s, command exceeds maximum volume %s",
				describe(), i, tip.CurrentVolume(), tip.MaxVol)
		} else if c, err := wells[i].Remove(v); err != nil {
			self.AddErrorf("Aspirate", "While %s - unexpected well error \"%s\"", describe(), err.Error())
		} else if fv.LessThan(tip.MinVol) {
			self.AddWarningf("Aspirate", "While %s - minimum tip volume is %s",
				describe(), tip.MinVol)
			//will get an error here, but ignore it since we're already raising a warning
			tip.Add(c)
		} else if err := tip.Add(c); err != nil {
			self.AddErrorf("Aspirate", "While %s - unexpected tip error \"%s\"", describe(), err.Error())
		}
	}

	if len(no_well) > 0 {
		self.AddErrorf("Aspirate", "While %s - %s on %s not in a well", describe(), pTips(len(no_well)), summariseChannels(no_well))
	}

	return ret
}

//Dispense - used
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int,
	platetype []string, what []string, llf []bool) driver.CommandStatus {

	ret := driver.CommandStatus{true, driver.OK, "DISPENSE ACK"}

	//check that arguments are valid
	if head < 0 || head >= self.state.GetNumberOfAdaptors() {
		self.AddErrorf("Dispense", "Cannot dispense from unknown head %d", head)
		return ret
	}
	adaptor := self.state.GetAdaptor(head)
	deck := self.state.GetDeck()

	if multi != adaptor.GetChannelCount() {
		self.AddErrorf("Dispense", "Multi[=%d] doesn't match number of channels on head %d[=%d]",
			multi, head, adaptor.GetChannelCount())
		return ret
	}

	if !self.testSliceLength("Dispense", map[string]int{
		"volume":    len(volume),
		"blowout":   len(blowout),
		"platetype": len(platetype),
		"what":      len(what),
		"llf":       len(llf),
	}, multi) {
		return ret
	}

	//find which channels are explicitly commanded
	channels := make([]bool, multi)
	for i := range channels {
		channels[i] = !(platetype[i] == "" && what[i] == "")
	}

	//find the position of each tip
	tip_pos := make([]wtype.Coordinates, multi)
	wells := make([]*wtype.LHWell, multi)
	for i := range channels {
		if ch := adaptor.GetChannel(i); ch.HasTip() {
			tip_pos[i] = ch.GetAbsolutePosition().Subtract(wtype.Coordinates{0., 0., ch.GetTip().GetSize().Z})

			for _, o := range deck.GetPointIntersections(tip_pos[i]) {
				if w, ok := o.(*wtype.LHWell); ok {
					wells[i] = w
					break
				}
			}
		}
	}

	//independence contraints

	//dispense
	for i, explicit := range channels {
		if !explicit {
			continue
		}
		v := wunit.NewVolume(volume[i], "ul")
		tip := adaptor.GetChannel(i).GetTip()

		if wells[i] == nil {
			self.AddErrorf("Dispense", "tip not in well")
		} else if c, err := tip.Remove(v); err != nil {
			self.AddErrorf("Dispense", "Unexpected tip error \"%s\"", err.Error())
		} else if err := wells[i].Add(c); err != nil {
			self.AddErrorf("Dispense", "Unexpected well error \"%s\"", err.Error())
		}
	}

	return ret
}

//LoadTips - used
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int,
	platetype, position, well []string) driver.CommandStatus {
	ret := driver.CommandStatus{true, driver.OK, "LOADTIPS ACK"}

	//check that the command is valid
	if !self.testTipArgs("LoadTips", channels, head, multi, platetype, position, well) {
		return ret
	}

	//make well coords
	wc := make([]wtype.WellCoords, multi)
	for i := range well {
		wc[i] = wtype.MakeWellCoords(well[i])
	}

	//check that all channels are empty
	adaptor := self.state.GetAdaptor(head)
	if tc := adaptor.GetTipCount(); tc > 0 {
		//get a slice of channels with tips on
		loaded := make([]string, 0, tc)
		for i := 0; i < adaptor.GetChannelCount(); i++ {
			if adaptor.GetChannel(i).HasTip() {
				loaded = append(loaded, fmt.Sprintf("%d", i))
			}
		}
		if tc > 1 {
			self.AddErrorf("LoadTips", "Cannot load tips to Head%d when channels %s already have tips loaded",
				head, strings.Join(loaded, ","))
		} else {
			self.AddErrorf("LoadTips", "Cannot load tips to Head%d when channel %s already has a tip loaded",
				head, loaded[0])
		}
		return ret
	}

	//get the deck
	deck := self.state.GetDeck()

	//Get the tip at each requested location
	tips := make([]*wtype.LHTip, multi)
	for _, i := range channels {
		if o, ok := deck.GetChild(position[i]); !ok {
			self.AddErrorf("LoadTips", "No known location \"%s\"", position[i])
			return ret
		} else if o == nil {
			self.AddErrorf("LoadTips", "No tipbox found at position %s, empty deck position", position[i])
			return ret
		} else if tipbox, ok := o.(*wtype.LHTipbox); !ok {
			self.AddErrorf("LoadTips", "No tipbox found at position %s, instead found %s \"%s\"",
				position[i], wtype.ClassOf(o), wtype.NameOf(o))
			return ret
		} else if !tipbox.AddressExists(wc[i]) {
			self.AddErrorf("LoadTips", "Request for tip at %s in tipbox of size [%dx%d]",
				wc[i].FormatA1(), tipbox.NCols(), tipbox.NRows())
			return ret
		} else {
			tips[i] = tipbox.GetChildByAddress(wc[i]).(*wtype.LHTip)
			if tips[i] == nil {
				self.AddErrorf("LoadTips", "Cannot load to channel %d as no tip at %s in tipbox \"%s\"",
					i, wc[i].FormatA1(), tipbox.GetName())
				return ret
			}
		}
	}

	//check alignment
	z_off := make([]float64, multi)
	misaligned := []string{}
	target := []string{}
	amount := []string{}
	for _, ch := range channels {
		tip_s := tips[ch].GetSize()
		tip_p := tips[ch].GetPosition().Add(wtype.Coordinates{0.5 * tip_s.X, 0.5 * tip_s.Y, tip_s.Z})
		ch_p := adaptor.GetChannel(ch).GetAbsolutePosition()
		delta := ch_p.Subtract(tip_p)
		if xy := delta.AbsXY(); xy > 0.5 {
			misaligned = append(misaligned, fmt.Sprintf("%d", ch))
			target = append(target, wc[ch].FormatA1())
			amount = append(amount, fmt.Sprintf("%v", xy))
		}
		z_off[ch] = delta.Z
		if delta.Z < 0. {
			self.AddError("LoadTips", "Request to load tip which is above adaptor - CRASH")
			return ret
		}
	}
	if len(misaligned) == 1 {
		self.AddErrorf("LoadTips", "Channel %s is misaligned with tip at %s by %smm",
			misaligned[0], target[0], amount[0])
		return ret
	} else if len(misaligned) > 1 {
		self.AddErrorf("LoadTips", "Channels %s are misaligned with tips at %s by %s mm respectively",
			strings.Join(misaligned, ","), strings.Join(target, ","), strings.Join(amount, ","))
		return ret
	}

	//if not independent, check there are no other tips in the way
	if !adaptor.IsIndependent() {
		collisions := adaptorCollisions{}
		zo_max := 0.
		zo_min := math.MaxFloat64
		for _, ch := range channels {
			if z_off[ch] > zo_max {
				zo_max = z_off[ch]
			}
			if z_off[ch] < zo_min {
				zo_min = z_off[ch]
			}
		}
		if zo_max != zo_min {
			self.AddErrorf("LoadTips", "Distance between channels and tips varies from %v to %v mm in non-independent head",
				zo_min, zo_max)
			return ret
		}
		for i := 0; i < adaptor.GetChannelCount(); i++ {
			if contains(i, channels) {
				continue
			}
			ch_pos := adaptor.GetChannel(i).GetAbsolutePosition()
			size := wtype.Coordinates{0., 0., zo_max + 0.5}
			box := wtype.NewBBox(ch_pos.Subtract(size), size)
			objects := deck.GetBoxIntersections(*box)
			//filter out tipboxes if we're meant to be ignoring them
			//(hack to prevent dubious tipbox geometry messing this up)
			if !self.tipbox_collisions {
				no_tipboxes := objects[:0]
				for _, o := range objects {
					if _, ok := o.(*wtype.LHTipbox); !ok {
						no_tipboxes = append(no_tipboxes, o)
					}
				}
				objects = no_tipboxes
			}
			if len(objects) > 0 {
				collisions = append(collisions, adaptorCollision{i, objects})
			}
		}

		if len(collisions) > 0 {
			self.AddErrorf("LoadTips", "Cannot load %s, %v (Head%d not independent)",
				summariseWell2Channel(well, channels), collisions, head)
		}
	}

	//move the tips to the adaptors
	for _, ch := range channels {
		tips[ch].GetParent().(*wtype.LHTipbox).RemoveTip(wc[ch])
		adaptor.GetChannel(ch).LoadTip(tips[ch])
		tips[ch].SetParent(nil)
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

	adaptor := self.state.GetAdaptor(head)
	deck := self.state.GetDeck()

	//Raise a warning if we're trying to eject tips that aren't there
	missing := []string{}
	for _, ch := range channels {
		if !adaptor.GetChannel(ch).HasTip() {
			missing = append(missing, fmt.Sprintf("%d", ch))
		}
	}
	if len(missing) == 1 {
		self.AddWarningf("UnloadTips", "No tip present at Head%d channel%s to eject", head, missing[0])
	} else if len(missing) > 0 {
		self.AddWarningf("UnloadTips", "No tips present on Head%d channels %s to eject", head, strings.Join(missing, ","))
	}

	//Check that this is possible
	if !adaptor.IsIndependent() {
		extra := []int{}
		for ch := 0; ch < multi; ch++ {
			if contains(ch, channels) {
				continue
			}
			if adaptor.GetChannel(ch).HasTip() {
				extra = append(extra, ch)
			}
		}
		if len(extra) == 1 {
			self.AddErrorf("UnloadTips", "Cannot unload tips from head%d %s without unloading tip from channel %s (head isn't independent)",
				head, summariseChannels(channels), extra[0])
			return ret
		} else if len(extra) > 1 {
			self.AddErrorf("UnloadTips", "Cannot unload tips from head%d %s without unloading tips from %s (head isn't independent)",
				head, summariseChannels(channels), summariseChannels(extra))
			return ret
		}
	}

	for _, ch := range channels {
		//get the target
		if target, ok := deck.GetChild(position[ch]); !ok {
			self.AddErrorf("UnloadTips", "Unknown deck position \"%s\"", position[ch])
			break
		} else if target == nil {
			self.AddErrorf("UnloadTips", "Cannot unload to empty deck location \"%s\"", position[ch])
			break
		} else if addr, ok := target.(wtype.Addressable); !ok {
			self.AddErrorf("UnloadTips", "Cannot unload tips to %s \"%s\" at location %s",
				wtype.ClassOf(target), wtype.NameOf(target), position[ch])
		} else {
			//get the location of the channel
			ch_pos := adaptor.GetChannel(ch).GetAbsolutePosition()
			//parse the wellcoords
			wc := wtype.MakeWellCoords(well[ch])
			if wc.IsZero() {
				self.AddErrorf("UnloadTips", "Cannot parse well coordinates \"%s\"", well[ch])
				break
			}
			if !addr.AddressExists(wc) {
				self.AddErrorf("UnloadTips", "Cannot unload to address %s in %s \"%s\" size [%dx%d]",
					wc.FormatA1(), wtype.ClassOf(target), wtype.NameOf(target), addr.NRows(), addr.NCols())
				break
			}
			//get the child - *LHTip or *LHWell
			child := addr.GetChildByAddress(wc)
			well_p, _ := addr.WellCoordsToCoords(wc, wtype.TopReference)
			delta := ch_pos.Subtract(well_p)

			switch target := target.(type) {
			case *wtype.LHTipbox:
				//put the tip in the tipbox
				if child.(*wtype.LHTip) != nil {
					self.AddErrorf("UnloadTips", "Cannot unload to tipbox \"%s\" %s, tip already present there",
						target.GetName(), wc.FormatA1())
				} else if delta.AbsXY() > 0.25 {
					self.AddErrorf("UnloadTips", "Head%d channel%d misaligned from tipbox \"%s\" %s by %.2fmm",
						head, ch, target.GetName(), wc.FormatA1(), delta.AbsXY())
				} else if delta.Z > target.GetSize().Z/2. {
					self.AddWarningf("UnloadTips", "Ejecting tip from Head%d channel%d to tipbox \"%s\" %s from height of %.2fmm",
						head, ch, target.GetName(), wc.FormatA1(), delta.Z)
				} else {
					target.PutTip(wc, adaptor.GetChannel(ch).UnloadTip())
				}

			case *wtype.LHTipwaste:
				//put the tip in the tipwaste
				if child == nil {
					//I don't think this should happen, but it would be embarressing to segfault...
					self.AddWarningf("UnloadTips", "Tipwaste \"%s\" well %s was nil, cannot check head alignment",
						target.GetName(), wc.FormatA1())
					adaptor.GetChannel(ch).UnloadTip()
				} else if max_delta := child.GetSize(); delta.X > max_delta.X || delta.Y > max_delta.Y {
					self.AddErrorf("UnloadTips", "Cannot unload, head%d channel%d is not above tipwaste \"%s\"",
						head, ch, target.GetName())
				} else if target.SpaceLeft() <= 0 {
					self.AddErrorf("UnloadTips", "Cannot unload tip to overfull tipwaste \"%s\", contains %d tips",
						target.GetName(), target.Contents)
				} else {
					target.Dispose(1)
					adaptor.GetChannel(ch).UnloadTip()
				}
			default:
				self.AddErrorf("UnloadTips", "Cannot unload tips to %s \"%s\" at location %s",
					wtype.ClassOf(target), wtype.NameOf(target), position[ch])
				break
			}
		}
		if self.HasError() {
			break
		}
	}

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
	if self.state.IsInitialized() {
		self.AddWarning("Initialize", "Call to initialize when robot is already initialized")
	}
	self.state.Initialize()
	return driver.CommandStatus{true, driver.OK, "INITIALIZE ACK"}
}

//Finalize - used
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
	if !self.state.IsInitialized() {
		self.AddWarning("Finalize", "Call to finalize when robot is not inisialized")
	}
	if self.state.IsFinalized() {
		self.AddWarning("Finalize", "Call to finalize when robot is already finalized")
	}
	self.state.Finalize()
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

		if err := self.state.GetDeck().SetChild(position, obj); err != nil {
			self.AddError("AddPlateTo", err.Error())
			return ret
		}

		if tb, ok := obj.(*wtype.LHTipbox); ok {
			//check that the height of the tips is greater than the height of the tipbox
			if tb.GetSize().Z >= (tb.TipZStart+tb.Tiptype.GetSize().Z) && self.tipbox_collisions {
				self.AddWarningf("AddPlateTo",
					"Tipbox \"%s\" is taller than the tips it holds (%.2fmm > %.2fmm), disabling tipbox collision detection",
					tb.GetName(), tb.GetSize().Z, tb.TipZStart+tb.Tiptype.GetSize().Z)
				self.tipbox_collisions = false
			}
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
