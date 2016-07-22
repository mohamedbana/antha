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
    tip             *wtype.LHTip        //Nil if no tip loaded, otherwise the tip that's loaded
    contents        *wtype.LHComponent  //What's in the tip?
    position        wtype.Coordinates  //position relative to the adaptor
    adaptor         *AdaptorState       //the channel's adaptor
}

func NewChannelState(adaptor *AdaptorState, position wtype.Coordinates) *ChannelState {
    r := ChannelState{}
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
func (self *ChannelState) LoadTip(tip *wtype.LHTip) error {
    
}

//UnloadTip
func (self *ChannelState) UnloadTip() *wtype.LHTip {
    
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
        as.channels = append(as.channels, NewChannelState(&as, channel_offset.Multiply(i)))
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
//                            RobotState
// -------------------------------------------------------------------------------

//RobotState Represent the physical state of a liquidhandling robot
type RobotState struct {
    positions       map[string]*wtype.Coordinates
    plates          map[string]interface{}
    adaptors        []*AdaptorState 
}

type AdaptorParams struct {
    initial_position    wtype.Coordinates,
    channels            int,
    channel_offset      wtype.Coordinates,
}

func NewRobotState(positions map[string]*wtype.Coordinates,
                   adaptors []AdaptorParams) *RobotState {
    rs := RobotState{
        positions,
        make(map[string]interface{},
        make([]*AdaptorState, 0, len(adaptors)),
    }
    for _,ap := range adaptors {
        rs.adaptors = append(rs.adaptors, NewAdaptorState(&rs, ap.initial_position, ap.channels, ap.channel_offset))
    }
    return &rs
}

//                            Accessors
//                            ---------

//GetPositionByString
func (self *RobotState) GetPositionByString(position string) wtype.Coordinates {
    return self.positions[position]
}

//GetPlateByString
func (self *RobotState) GetPlateByString(position string) interface{} {
    return self.plates[position]
}

//GetPlateBelow
func (self *RobotState) GetPlateBelow(position wtype.Coordinates) interface{} {
}

//GetWellBelow
func (self *RobotState) GetWellBelow(position wtype.Coordinates) interface{}, *wtype.LHWell {
}
