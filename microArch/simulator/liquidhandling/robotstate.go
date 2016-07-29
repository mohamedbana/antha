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
    "github.com/antha-lang/antha/microArch/simulator"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
    "math"
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

//SetTip set the tip, mainly for debug and testing
func (self *ChannelState) SetTip(tip *wtype.LHTip) {
    self.tip = tip
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

type GetPlateRet struct {
    Plate       interface{}
    Position    wtype.WellCoords
    Offset      wtype.Coordinates 
}

//GetPlate Get the next plate below the adaptor and check that the position, type and well match
func (self *ChannelState) GetPlate(platetype, position, well string) (*GetPlateRet, *simulator.SimulationError) {
    //get the plate interface from the robot deck
    abs_p := self.GetAbsolutePosition()
    plate_loc := self.adaptor.GetRobot().GetDeck().GetPlateBelow(abs_p)

    //if there's no plate
    if plate_loc == nil {
        return nil, simulator.NewErrorf("", "No Plate below channel %v, expected type \"%s\" at \"%s\"", self.number, platetype, position)
    }

    //check the type of the plate
    if pt := plate_loc.plate.(wtype.Typed).GetType(); pt != platetype {
        return nil, simulator.NewErrorf("", "Plate below channel %v is of type \"%s\" not \"%s\"", self.number, pt, platetype)
    }

    //check the position of the plate
    if plate_loc.location_name != position {
        return nil, simulator.NewErrorf("", "Plate below channel %v is in location \"%s\" not \"%s\"",
                               self.number, plate_loc.location_name, position)
    }

    //check the well
    wc, offset := plate_loc.plate.(wtype.LHDeckObject).CoordsToWellCoords(abs_p.Subtract(plate_loc.GetOffset()))
    if wc.FormatA1() != well {
        return nil, simulator.NewErrorf("", "Channel %v is above well %s not well %s", self.number, wc.FormatA1(), well)
    }

    //return everything
    r := GetPlateRet{
        plate_loc.plate,
        wc,
        offset}
    return &r, nil

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
func (self *ChannelState) LoadTip(platetype, position, well string) *simulator.SimulationError {
    if self.tip != nil {
        return simulator.NewErrorf("", "Tip already loaded on channel %v", self.number)
    }

    //is this channel expected to remain empty
    if platetype == "" && position == "" && well == "" {
        //Todo: check that I'm not going to collide with anything while tip loading motion is made

        return nil
    }

    p, err := self.GetPlate(platetype, position, well)
    if err != nil {
        return err
    }

    //check that we've got a tipbox
    tipbox, ok := p.Plate.(*wtype.LHTipbox)
    if !ok {
        return simulator.NewErrorf("", "No Tipbox beneath channel %v", self.number)
    }

    //Check that we're lined up
    if misalignment := p.Offset.AbsXY(); misalignment > 0.5 {
        return simulator.NewErrorf("", "Channel %v Misaligned from tip by %smm", self.number, misalignment)
    }

    //if there's a tip, move it to the adaptor
    //at first, I thought there being no tip (tip==nil) was an error, but
    //actually it could be deliberate in the case of a non-independent adaptor
    //and the callee should verify afterwards that the expected tips were loaded
    wc := wtype.MakeWellCoords(well)
    tip, _ := tipbox.GetCoords(wc)
    tipbox.RemoveTip(wc)
    self.tip = tip.(*wtype.LHTip)
    return nil
}

//UnloadTip
func (self *ChannelState) UnloadTip(platetype, position, well string) *simulator.SimulationError {
    p, err := self.GetPlate(platetype, position, well)
    if err != nil {
        return err
    }

    //check that we've got a tipwaste
    tipwaste, ok := p.Plate.(*wtype.LHTipwaste)
    if !ok {
        return simulator.NewError("", "No tipwaste below adaptor")
    }

    //check that we're actually over the well
    if self.HasTip() && 
        (math.Abs(p.Offset.X) > 0.5*tipwaste.AsWell.Xdim ||
         math.Abs(p.Offset.Y) > 0.5*tipwaste.AsWell.Ydim) {
           return simulator.NewError("", "Ejecting a tip while not over TipWaste well")
    }

    //do it!
    self.tip = nil
    self.contents = nil

    if !tipwaste.Dispose(1) {
        return simulator.NewError("", "Tipbox full")
    }
    return nil
}

// -------------------------------------------------------------------------------
//                            AdaptorState
// -------------------------------------------------------------------------------

//AdaptorState Represent the physical state and layout of the adaptor
type AdaptorState struct {
    channels        []*ChannelState
    position        wtype.Coordinates
    independent     bool
    robot           *RobotState 
}

func NewAdaptorState(robot *RobotState, 
                     initial_position wtype.Coordinates, 
                     independent bool,
                     channels int, 
                     channel_offset wtype.Coordinates) *AdaptorState {
    as := AdaptorState{
        make([]*ChannelState,0,channels),
        initial_position,
        independent,
        robot}

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
    for _,ch := range self.channels {
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

//                            Actions
//                            -------

func (self *AdaptorState) LoadTips(platetype, position, well []string) *simulator.SimulationError {
    for i := range self.channels {
        if !(self.independent && platetype[i] == "" && position[i] == "" && well[i] == "") {
            if err := self.channels[i].LoadTip(platetype[i], position[i], well[i]); err != nil {
                return err
            }            
        }
    }
    return nil
}

func (self *AdaptorState) UnLoadTips(platetype, position, well []string) error {
    for i := range self.channels {
        if !(self.independent && platetype[i] == "" && position[i] == "" && well[i] == "") {
            if err := self.channels[i].UnloadTip(platetype[i], position[i], well[i]); err != nil {
                return err
            }            
        }
    }
    return nil
}

func (self *AdaptorState) Move(platetype, position, well []string, 
        reference []wtype.WellReference, offset []wtype.Coordinates) *simulator.SimulationError {
    positions := make([]*wtype.Coordinates, len(self.channels))

    for i := range self.channels {
        //leave position nil if we don't care
        if platetype[i] == "" && position[i] == "" && well[i] == "" {
            continue
        }

        pl := self.robot.GetDeck().GetPlateByPosition(position[i], platetype[i])
        if pl == nil {
            return simulator.NewErrorf("No plate of type \"%s\" at location \"%s\"", platetype[i], position[i])
        }
        do := pl.GetPlate().(wtype.LHDeckObject)

        wc := wtype.MakeWellCoords(well[i])
        if !do.HasCoords(wc) {
            return simulator.NewErrorf("Plate has no coordinates %s", well[i])
        }

        rel_pos, ok := do.WellCoordsToCoords(wc, reference[i])
        if !ok {
            return simulator.NewErrorf("", "Could not get location of well %s", wc.Format1A())
        }

        //TODO Check for collision with plate/well

        pos_i := pl.GetOffset().Add(rel_pos).Add(offset[i])
        positions[i] = &pos_i
    }

    //convert positions to relative
    //find the origin of the adaptor
    var origin *wtype.Coordinates 
    for i := range self.channels {
        if positions[i] != nil {
            o := positions[i].Subtract(self.channels[i].GetRelativePosition())
            origin = &o
        }
    }
    if origin == nil {
        return simulator.NewWarning("", "Ignoring empty command")
    }
    for i := range positions {
        if positions[i] != nil { 
            *positions[i] = positions[i].Subtract(*origin)
        } else { //if position hasn't been specified, leave relative position alone
            p := self.channels[i].GetRelativePosition()
            positions[i] = &p
        }
    }

    //if the head isn't independent, check that relative positions are the same
    if !self.independent {
        for i := range positions {
            if *positions[i] != self.channels[i].GetRelativePosition() {
                return simulator.NewError("", "Channels can't move independently")
            }
        }
    }

    self.position = *origin
    for i := range self.channels {
        self.channels[i].SetRelativePosition(*positions[i])
    }
    return nil
}

// -------------------------------------------------------------------------------
//                            PlateLocation
// -------------------------------------------------------------------------------

//PlateLocation a wrapper for an LHObject
type PlateLocation struct {
    plate_name      string
    offset          wtype.Coordinates
    size            wtype.Coordinates
    location_name   string
    plate           interface{}
}

func NewPlateLocation(plate_name string, plate wtype.LHDeckObject, position wtype.Coordinates, position_name string) *PlateLocation {
    r := PlateLocation{
        plate_name,
        position,
        plate.GetSize(),
        position_name,
        plate.(interface{})}
    return &r
}

func (self *PlateLocation) GetPlate() interface{} {
    return self.plate
}

func (self *PlateLocation) GetSize() wtype.Coordinates {
    return self.size
}

func (self *PlateLocation) GetOffset() wtype.Coordinates {
    return self.offset
}

func (self *PlateLocation) GetName() string {
    return self.plate_name
}

func (self *PlateLocation) GetLocationName() string {
    return self.location_name
}

func (self *PlateLocation) ContainsXY(p wtype.Coordinates) bool {
    d := p.Subtract(self.offset)
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
    named_positions         map[string]wtype.Coordinates
    plate_locations         []*PlateLocation
}

func NewDeckState(positions map[string]wtype.Coordinates) *DeckState {
    r := DeckState{
        positions,
        make([]*PlateLocation, 0),
    }
    return &r
}

func (self *DeckState) HasPosition(position_name string) bool {
    if _, ok := self.named_positions[position_name]; ok {
        return true
    }
    return false
}

func (self *DeckState) AddPlate(plate_name string, plate interface{}, position wtype.Coordinates, position_name string) *simulator.SimulationError {
    to_add := NewPlateLocation(plate_name, plate.(wtype.LHDeckObject), position, position_name)
    for _,pl := range self.plate_locations {
        if pl.Intersects(to_add) {
            return simulator.NewErrorf("", "Cannot add \"%s\" to location \"%s\", intersects with \"%s\" at \"%s\"", 
                             plate_name, position_name, pl.GetName(), pl.GetLocationName())
        }
    }
    self.plate_locations = append(self.plate_locations, to_add)
    return nil
}

func (self *DeckState) AddPlateToNamed(plate_name string, plate interface{}, position string) *simulator.SimulationError {
    return self.AddPlate(plate_name, plate, self.named_positions[position], position)
}

//GetPlateBelow get the next plate below the position and the relative position within the plate
func (self *DeckState) GetPlateBelow(pos wtype.Coordinates) *PlateLocation {
    var p *PlateLocation = nil
    z := math.MaxFloat64
    for _,pl := range self.plate_locations {
        if pl.ContainsXY(pos) {
            if dz := pos.Z - pl.GetOffset().Z; dz > 0 && dz < z {
                p = pl
                z = dz
            }
        }
    }
    return p
}

func (self *DeckState) GetPlateByPosition(position, platetype string) *PlateLocation {
    ret := make([]*PlateLocation,0)
    for _,pl := range self.plate_locations {
        if pl.GetLocationName() == position && pl.GetPlate().(wtype.Typed).GetType() == platetype {
            ret = append(ret, pl)
        }
    }
    if len(ret) != 1 {
        return nil
    }
    return ret[0]
}

func (self *DeckState) GetPlateAt(position string) interface{} {
    ret := make([]*PlateLocation,0)
    for _,pl := range self.plate_locations {
        if pl.GetLocationName() == position {
            ret = append(ret, pl)
        }
    }
    if len(ret) == 0 {
        return nil
    }
    return ret[0].GetPlate()
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
    Initial_position    wtype.Coordinates
    Independent         bool
    Channels            int
    Channel_offset      wtype.Coordinates
}

func NewRobotState(positions map[string]wtype.Coordinates,
                   adaptors []AdaptorParams) *RobotState {
    rs := RobotState{}
    rs.deck = NewDeckState(positions)
    rs.adaptors = make([]*AdaptorState, 0, len(adaptors))
    rs.initialized = false
    rs.finalized = false
    for _,ap := range adaptors {
        rs.adaptors = append(rs.adaptors, NewAdaptorState(&rs, ap.Initial_position, ap.Independent, ap.Channels, ap.Channel_offset))
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

//GetNumberOfAdaptors
func (self *RobotState) GetNumberOfAdaptors() int {
    return len(self.adaptors)
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
func (self *RobotState) Initialize() *simulator.SimulationError {
    if self.initialized {
        return simulator.NewError("", "Called Initialize on already initialised liquid handler")
    }
    self.initialized = true
    return nil
}

//Finalize
func (self *RobotState) Finalize() error {
    if self.finalized {
        return simulator.NewError("", "Called Finalize on already finalized liquid handler")
    }
    if !self.initialized {
        return simulator.NewError("", "Called Finalize on uninitialized liquidhandler")
    }
    self.finalized = true
    return nil
}

