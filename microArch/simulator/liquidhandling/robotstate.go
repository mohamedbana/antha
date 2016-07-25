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
    "errors"
    "fmt"
)

// -------------------------------------------------------------------------------
//                            ChannelState
// -------------------------------------------------------------------------------

//ChannelState Represent the physical state of a single channel
type ChannelState struct {
    number          int
    tip             *wtype.LHTip        //Nil if no tip loaded, otherwise the tip that's loaded
    contents        *wtype.LHComponent  //What's in the tip?
    position        wtype.Coordinates  //position relative to the adaptor
    adaptor         *AdaptorState       //the channel's adaptor
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
func (self *ChannelState) GetContents() wtype.LHComponent {
    return self.contents
}

//GetRelativePosition get the channel's position relative to the head
func (self *ChannelState) GetRelativePosition() wtype.Coordiantes {
    return self.position
}

//GetAbsolutePosition get the channel's absolute position
func (self *ChannelState) GetAbsolutePosition() wtype.Coordinates {
    return self.position.Add(self.adaptor.GetPosition())
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
func (self *ChannelState) LoadTip(platetype, position, well string) error {
    p := self.GetAbsolutePosition()
    plate, rel_p := self.adaptor.GetRobot().GetDeck().GetPlateBelow(p)
}

//UnloadTip
func (self *ChannelState) UnloadTip(platetype, position, well string) *wtype.LHTip, error {
    
}

// -------------------------------------------------------------------------------
//                            AdaptorState
// -------------------------------------------------------------------------------

//AdaptorState Represent the physical state and layout of the adaptor
type AdaptorState struct {
    channels        []*ChannelState
    position        wtype.Coordinates
    robot           *RobotState 
}

func NewAdaptorState(robot *RobotState, 
                     initial_position wtype.Coordinates, 
                     channels int, 
                     channel_offset wtype.Coordinates) *AdaptorState {
    as := AdaptorState{
        make([]*ChannelState,0,channels),
        initial_position,
        robot}

    for i := 0; i < channels; i++ {
        as.channels = append(as.channels, NewChannelState(i, &as, channel_offset.Multiply(i)))
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
    return len(channels)
}

//GetTipCount
func (self *AdaptorState) GetTipCount() int {
    r := 0
    for _,ch := range self.channels {
        if ch.HasTip() {
            r++
        }
    }
    return r
}

//GetRobot
func (self *AdaptorState) GetRobot() *RobotState {
    return self.robot
}

// -------------------------------------------------------------------------------
//                            PlateLocation
// -------------------------------------------------------------------------------

//PlateLocation a wrapper for an LHObject
type PlateLocation struct {
    offset  wtype.Coordinates
    size    wtype.Coordinates
    plate   interface{}
}

func NewPlateLocation(plate wtype.LHObject, position wtype.Coordinates) *PlateLocation {
    r := PlateLocation{
        position,
        plate.GetSize(),
        plate.(interface{})}
    return &r
}

func (self *PlateLocation) GetPlate() interface{} {
    return self.plate
}

func (self *PlateLocation) GetOffset() wtype.Coordinates {
    return self.offset
}

func (self *PlateLocation) ContainsXY(p wtype.Coordinates) bool {
    d := self.offset.Subtract(p)
    return d.X >=0 && d.Y >= 0 && d.X < self.size.X && d.Y < self.size.X
}

func (self *PlateLocation) Intersects(rhs *PlateLocation) bool {
    //test a single dimension. 
    //(a,b) are the start and end of the first position
    //(c,d) are the start and end of the second pos
    // assert(a > b  and  d > c)
    f := func(a,b,c,d float64) bool {
        return !(c >= b || d <= a)
    }

    return (f(self.offset.X, self.offset.X + self.size.X,
               rhs.offset.X,  rhs.offset.X +  rhs.size.X) &&
            f(self.offset.Y, self.offset.Y + self.size.Y,
               rhs.offset.Y,  rhs.offset.Y +  rhs.size.Y) &&
            f(self.offset.Z, self.offset.Z + self.size.Z,
               rhs.offset.Z,  rhs.offset.Z +  rhs.size.Z))
}

// -------------------------------------------------------------------------------
//                            DeckState
// -------------------------------------------------------------------------------

//DeckState Represent the physical state of an LH robot's deck
type DeckState struct {
    named_positions         map[string]*wtype.Coordinates
    positions               map[string][]interface{}
    plate_locations         []*PlateLocation
}

func NewDeckState(positions map[string]*wtype.Coordinates) *DeckState {
    r := DeckState{
        positions,
        make(map[string][]interface{}),
        make([]*PlateLocation, 0)
    }
    for k := range positions {
        r.positions[k] = make([]interface{}, 0)
    }
    return &r
}

func (self *DeckState) AddPlate(plate interface{}, position wtype.Coordinates) bool {
    to_add := NewPlateLocation(plate.(wtype.LHObject), position)
    for _,pl := range self.plate_locations {
        if pl.Intersects(to_add) {
            return false
        }
    }
    self.plate_locations = append(self.plate_locations, to_add)
    return true
}

func (self *DeckState) AddPlateToNamed(plate interface{}, position string) bool {
    return self.AddPlate(plate, self.named_positions[position])
}

//GetPlateBelow get the next plate below the position and the relative position within the plate
func (self *DeckState) GetPlateBelow(pos wtype.Coordinates) interface{}, wtype.Coordinates {
    var p *PlateLocation = nil
    z := math.MAXFloat64
    for _,pl := range self.plate_locations {
        if pl.ContainsXY(pos) {
            if dz := pos.Z - pl.Offset().Z; dz > 0 && dz < z {
                p = pl
                z = dz
            }
        }
    }
    if p == nil {
        return nil, wtype.Coordinates{}
    }
    return p.GetPlate(), pos.Subtract(p.GetOffset())
}

func (self *DeckState) GetPlatesName(position string) []interface{} {
    return self.positions[position]
}

// -------------------------------------------------------------------------------
//                            RobotState
// -------------------------------------------------------------------------------

//RobotState Represent the physical state of a liquidhandling robot
type RobotState struct {
    deck            *DeckState
    adaptors        []*AdaptorState 
    initialized     bool
    finalized       bool
}

type AdaptorParams struct {
    initial_position    wtype.Coordinates,
    channels            int,
    channel_offset      wtype.Coordinates,
}

func NewRobotState(positions map[string]*wtype.Coordinates,
                   adaptors []AdaptorParams) *RobotState {
    rs := RobotState{}
    rs.deck = NewDeckState(positions)
    rs.adaptors = make([]*AdaptorState, 0, len(adaptors))
    rs.initialised = false
    rs.finalised = false
    for _,ap := range adaptors {
        rs.adaptors = append(rs.adaptors, NewAdaptorState(&rs, ap.initial_position, ap.channels, ap.channel_offset))
    }
    return &rs
}

//                            Accessors
//                            ---------

//GetDeck
func (self *RobotState) GetDeck() *DeckState {
    return self.deck
}

//GetAdaptor
func (self *RobotState) GetAdaptor(num int) *AdaptorState {
    return self.adaptors[num]
}

//IsInitialized
func (self *RobotState) IsInitialized() bool {
    return self.initialized
}

//IsFinalized
func (self *RobotState) IsFinalized() bool {
    return self.finalized
}

//                            Actions
//                            -------

//Initialize
func (self *RobotState) Initialize() error {
    if self.initialized {
        return errors.New("Called Initialize on already initialised liquid handler")
    }
    self.initialized = true
    return nil
}

//Finalize
func (self *RobotState) Finalize() error {
    if self.finalized {
        return errors.New("Called Finalize on already finalized liquid handler")
    }
    if !self.initialized {
        return errors.New("Called Finalize on uninitialized liquidhandler")
    }
    self.finalized = true
    return nil
}

