// /anthalib/simulator/liquidhandling/robotstate.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/simulator"
	"math"
)

// -------------------------------------------------------------------------------
//                            ChannelState
// -------------------------------------------------------------------------------

//ChannelState Represent the physical state of a single channel
type ChannelState struct {
	number   int
	tip      *wtype.LHTip       //Nil if no tip loaded, otherwise the tip that's loaded
	contents *wtype.LHComponent //What's in the tip?
	position wtype.Coordinates  //position relative to the adaptor
	adaptor  *AdaptorState      //the channel's adaptor
}

func NewChannelState(number int, adaptor *AdaptorState, position wtype.Coordinates) *ChannelState {
	r := ChannelState{}
	r.number = number
	r.position = position
	r.adaptor = adaptor

	return &r
}

//                            Accessors
//                            ---------

//HasTip is a tip loaded
func (self *ChannelState) HasTip() bool {
	return self.tip != nil
}

//GetTip get the loaded tip, returns nil if none loaded
func (self *ChannelState) GetTip() *wtype.LHTip {
	return self.tip
}

//IsEmpty returns true only if a tip is loaded and contains liquid
func (self *ChannelState) IsEmpty() bool {
	return self.HasTip() && self.contents != nil && self.contents.IsZero()
}

//GetContents get the contents of the loaded tip, retuns nil if no contents or no tip
func (self *ChannelState) GetContents() *wtype.LHComponent {
	return self.contents
}

//GetRelativePosition get the channel's position relative to the head
func (self *ChannelState) GetRelativePosition() wtype.Coordinates {
	return self.position
}

//SetRelativePosition get the channel's position relative to the head
func (self *ChannelState) SetRelativePosition(v wtype.Coordinates) {
	self.position = v
}

//GetAbsolutePosition get the channel's absolute position
func (self *ChannelState) GetAbsolutePosition() wtype.Coordinates {
	return self.position.Add(self.adaptor.GetPosition())
}

//GetTarget get the LHObject below the adaptor
func (self *ChannelState) GetTarget() wtype.LHObject {
	pos := self.GetAbsolutePosition()
	zbox := wtype.NewZBox4f(pos.X, pos.Y, 0, 0)

	objs := self.adaptor.GetRobot().GetObjectsIn(zbox)
	if len(objs) == 0 {
		return nil
	}

	z_min := math.MaxFloat64
	i_min := -1
	for i := range objs {
		bb := objs[i].GetBounds()
		if gap := pos.Z - (bb.GetPosition().Z + bb.GetSize().Z); gap > 0 && gap < z_min {
			z_min = gap
			i_min = i
		}
	}
	if i_min < 0 {
		return nil
	}
	return objs[i_min]
}

//                            Actions
//                            -------

//Aspirate
func (self *ChannelState) Aspirate(volume wunit.Volume) error {

	return nil
}

//Dispense
func (self *ChannelState) Dispense(volume *wunit.Volume) error {

	return nil
}

//LoadTip
func (self *ChannelState) LoadTip(tip *wtype.LHTip) {
	self.tip = tip
}

//UnloadTip
func (self *ChannelState) UnloadTip() *wtype.LHTip {
	tip := self.tip
	self.tip = nil
	return tip
}

// -------------------------------------------------------------------------------
//                            AdaptorState
// -------------------------------------------------------------------------------

//AdaptorState Represent the physical state and layout of the adaptor
type AdaptorState struct {
	channels    []*ChannelState
	position    wtype.Coordinates
	independent bool
	robot       *RobotState
}

func NewAdaptorState(independent bool,
	channels int,
	channel_offset wtype.Coordinates) *AdaptorState {
	as := AdaptorState{
		make([]*ChannelState, 0, channels),
		wtype.Coordinates{},
		independent,
		nil,
	}

	for i := 0; i < channels; i++ {
		as.channels = append(as.channels, NewChannelState(i, &as, channel_offset.Multiply(float64(i))))
	}

	return &as
}

//                            Accessors
//                            ---------

//GetPosition
func (self *AdaptorState) GetPosition() wtype.Coordinates {
	return self.position
}

//GetChannelCount
func (self *AdaptorState) GetChannelCount() int {
	return len(self.channels)
}

//GetChannel
func (self *AdaptorState) GetChannel(ch int) *ChannelState {
	return self.channels[ch]
}

//GetTipCount
func (self *AdaptorState) GetTipCount() int {
	r := 0
	for _, ch := range self.channels {
		if ch.HasTip() {
			r++
		}
	}
	return r
}

//IsIndependent
func (self *AdaptorState) IsIndependent() bool {
	return self.independent
}

//GetRobot
func (self *AdaptorState) GetRobot() *RobotState {
	return self.robot
}

//SetRobot
func (self *AdaptorState) SetRobot(r *RobotState) {
	self.robot = r
}

func (self *AdaptorState) SetPosition(p wtype.Coordinates) {
	self.position = p
}

//                            Actions
//                            -------

func (self *AdaptorState) Move(target wtype.LHObject, wc []wtype.WellCoords, ref []wtype.WellReference, off []wtype.Coordinates) *simulator.SimulationError {
	addr, ok := target.(wtype.Addressable)
	if !ok {
		if n, nok := target.(wtype.Named); nok {
			return simulator.NewErrorf("", "Target object \"%s\" is not addressable", n.GetName())
		} else {
			return simulator.NewErrorf("", "Target object is not addressable")
		}
	}

	//find the origin
	origin := wtype.Coordinates{}
	for i := range wc {
		if addr.HasLocation(wc[i]) {
			origin, _ = addr.WellCoordsToCoords(wc[i], ref[i])
			origin = origin.Add(off[i]).Subtract(self.channels[i].GetRelativePosition())
			break
		}
	}

	//find the relative positions
	positions := make([]wtype.Coordinates, len(self.channels))
	for i := range self.channels {
		if wc[i].IsZero() {
			positions[i] = self.GetChannel(i).GetRelativePosition()
		} else {
			if pos, ok := addr.WellCoordsToCoords(wc[i], ref[i]); ok {
				positions[i] = pos.Add(off[i]).Subtract(origin)
			} else {
				return simulator.NewErrorf("", "No well \"%s\" in target", wc[i].FormatA1())
			}
		}
	}

	//if not independent, relative positions should not change
	if !self.independent {
		for i := range positions {
			if positions[i] != self.channels[i].GetRelativePosition() {
				return simulator.NewError("", "Failed to adjust channel offset in non-independent adaptor")
			}
		}
	}

	//do it
	self.position = origin
	for i := range self.channels {
		self.channels[i].SetRelativePosition(positions[i])
	}
	return nil
}

// -------------------------------------------------------------------------------
//                            RobotState
// -------------------------------------------------------------------------------

//RobotState Represent the physical state of a liquidhandling robot
type RobotState struct {
	slots       map[string]wtype.LHSlot
	adaptors    []*AdaptorState
	initialized bool
	finalized   bool
}

func NewRobotState() *RobotState {
	rs := RobotState{}
	rs.slots = make(map[string]wtype.LHSlot)
	rs.adaptors = make([]*AdaptorState, 0)
	rs.initialized = false
	rs.finalized = false
	return &rs
}

//                            Accessors
//                            ---------

//GetAdaptor
func (self *RobotState) GetAdaptor(num int) *AdaptorState {
	return self.adaptors[num]
}

//GetNumberOfAdaptors
func (self *RobotState) GetNumberOfAdaptors() int {
	return len(self.adaptors)
}

//AddAdaptor
func (self *RobotState) AddAdaptor(a *AdaptorState) {
	a.SetRobot(self)
	self.adaptors = append(self.adaptors, a)
}

//AddSlot
func (self *RobotState) AddSlot(s wtype.LHSlot) {
	self.slots[s.(wtype.Named).GetName()] = s
}

//GetSlot
func (self *RobotState) GetSlot(name string) wtype.LHSlot {
	return self.slots[name]
}

//IsInitialized
func (self *RobotState) IsInitialized() bool {
	return self.initialized
}

//IsFinalized
func (self *RobotState) IsFinalized() bool {
	return self.finalized
}

//GetObjectsIn Return all slots that intersect with the bounding box
func (self *RobotState) GetObjectsIn(bb *wtype.BBox) []wtype.LHObject {
	ret := make([]wtype.LHObject, 0)
	for _, slot := range self.slots {
		if slot.GetChild() != nil && bb.Intersects(slot.GetChild().GetBounds()) {
			ret = append(ret, slot.GetChild())
		}
	}
	return ret
}

//                            Actions
//                            -------

//Initialize
func (self *RobotState) Initialize() *simulator.SimulationError {
	if self.initialized {
		return simulator.NewError("", "Called Initialize on already initialised liquid handler")
	}
	self.initialized = true
	return nil
}

//Finalize
func (self *RobotState) Finalize() *simulator.SimulationError {
	if self.finalized {
		return simulator.NewError("", "Called Finalize on already finalized liquid handler")
	}
	if !self.initialized {
		return simulator.NewError("", "Called Finalize on uninitialized liquidhandler")
	}
	self.finalized = true
	return nil
}

//AddObject
func (self *RobotState) AddObject(slot_name string, o wtype.LHObject) *simulator.SimulationError {
	if sl, ok := self.slots[slot_name]; ok {
		//check that the slot is empty and can hold a child of this type
		if child := sl.GetChild(); child != nil {
			//In the future, we'll check if the child can accept another LHObject, for now barf
			cname := "unknown"
			oname := "unknown"
			if n, ok := child.(wtype.Named); ok {
				cname = n.GetName()
			}
			if n, ok := o.(wtype.Named); ok {
				oname = n.GetName()
			}
			return simulator.NewErrorf("",
				"Couldn't add \"%s\" to location \"%s\" which already contains \"%s\"",
				oname, slot_name, cname)
		} else if !sl.Accepts(o) {
			return simulator.NewError("", sl.SetChild(o).Error())
		}

		//check for intersections with other objects
		bb := o.GetBounds()
		bb.SetPosition(sl.GetChildPosition())
		for name, slot := range self.slots {
			if c := slot.GetChild(); c != nil && bb.Intersects(c.GetBounds()) {
				return simulator.NewErrorf("", "Object intersects with object at position \"%s\"", name)
			}
		}
		sl.SetChild(o)
	} else {
		return simulator.NewErrorf("", "Robot contains no locations named \"%s\"", slot_name)
	}
	return nil
}
