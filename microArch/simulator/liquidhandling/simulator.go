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
	"github.com/antha-lang/antha/microArch/driver"
)

func Simulate(properties *LHProperties, instructions []TerminalRobotInstruction)
{
    vlh := NewVirtualLiquidHandler(properties)

    for _, ins = range instructions {
        ins.OutputTo(vlh)
    }
}

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
    properties LHProperties 
    //Need to store:
    // LHProperties
    // plate(s) at each layout position
    //   contents of each well of each plate / tip box
    //     liquid type / tip type
    // contents of each loaded tip
    //   liquid type
    // tips on each adaptor
    //   LHAdapter know how many tips it has, and what type they are (assumed equal)
    //   but it can't tell which position they're on
    //   tip properties
    // adaptors on each head
    //   Adapeter properties
    // head location
    //   Head properties
}

//Create a new VirtualLiquidHandler which mimics an LHDriver
func NewVirtualLiquidHandler(props LHProperties) (*VirtualLiquidHandler, error) {
    var vlh VirtualLiquidHandler

    vlh.properties = props.Dup()

    return &vlh
}

//Move command - used
func (self *VirtualLiquidHandler) Move(deckposition []string, wellcoords []string, reference []int, 
                                       offsetX, offsetY, offsetZ []float64, plate_type []string, 
                                       head int) driver.CommandStatus {
    //Asserts:
    //deckposition exists - why is it a list?
    //wellcoords exists within the plate at deckposition
    //reference is in allowable range
    //offsetX, offsetY, offsetZ are within the well (but I guess they needn't be...)
    //plate_type matches the type of plate at deckposition
    //head is valid
    driver.CommandStatus{true, driver.OK, "MOVE ACK"}
}

//Move raw - not yet implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
    //Asserts:
    //head exists
    //x,y,x are within the machine
    panic("unimplemented")
}

//Aspirate - used
func (self *VirtualLiquidHandler) Aspirate(volume []float64, overstroke []bool, head int, multi int, 
                                           platetype []string, what []string, llf []bool) 
                                  driver.CommandStatus {
    //volumes are equal if adapter isn't independent
    //tips are loaded in each adapter location the aspirates
    //volume is smaller that the tips' maximum capacity
    //the tip hasn't been used for a different liquid
    //head exists
    //multi matches number of tips loaded
    //platetype matches the plate at the location we moved to
    //what matches the expected liquid class
    //llf is the right size, cannot vary unless independent
    panic("unimplemented")
}

//Dispense - used
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int, 
                                           platetype []string, what []string, llf []bool) 
                                  driver.CommandStatus {
    //Volumes are equal if adapter isn't indepentent
    //Volumes are at most equal to the volume in the tip
    //blowout is the right length
    //head exists
    //multi is valid
    //platetype matches the type of plate that we're next to
    //what matches the liquid class that was aspirated
    //llf is the right size and follows independence constraint
    panic("unimplemented")
}

//LoadTips - used
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int, 
                                           platetype, position, well []string) driver.CommandStatus {
    //channels is correct length and value (how does this work)
    //head exists
    //multi is in correct range for adaptor
    //platetype matches the plate we're over
    //position is correct, tips still exists there
    //well exists (difference between platetype and well?)
    panic("unimplemented")
}

//UnloadTips - used
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int, 
                                             platetype, position, well []string) driver.CommandStatus {
    //Tips are loaded in channels
    //independence constraints are met
    //head exists
    //multi is correct
    //platetype matches the plate we're over
    //platetype is tip-waste
    //position and well are correct
    panic("unimplemented")
}

//SetPipetteSpeed - used
func (self *VirtualLiquidHandler) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
    //head exists
    //channel exists
    //speed is within allowable range
    panic("unimplemented")
}

//SetDriveSpeed - used
func (self *VirtualLiquidHandler) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
    //drive string?
    //rate is within allowable range (what is this?)
    panic("unimplemented")
}

//Stop - unused
func (self *VirtualLiquidHandler) Stop() driver.CommandStatus {
    panic("unimplemented")
}

//Go - unused
func (self *VirtualLiquidHandler) Go() driver.CommandStatus {
    panic("unimplemented")
}

//Initialize - used
func (self *VirtualLiquidHandler) Initialize() driver.CommandStatus {
    //check that this is called before anything else?
    panic("unimplemented")
}

//Finalize - used
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
    //check that this is called last, no more calls
    panic("unimplemented")
}

//Wait - used
func (self *VirtualLiquidHandler) Wait(time float64) driver.CommandStatus {
    //time is positive
    //maybe a warning if it's super-long
    panic("unimplemented")
}

//Mix - used
func (self *VirtualLiquidHandler) Mix(head int, volume []float64, platetype []string, cycles []int, 
                                      multi int, what []string, blowout []bool) driver.CommandStatus {
    //head exists
    //volume is lte volume in wells
    //platetype matches the plate we're over
    //muli is correct
    //what matches expected liquidclass
    //volume, platetype, what, blowout match independence constraint
    panic("unimplemented")
}

//ResetPistons - used
func (self *VirtualLiquidHandler) ResetPistons(head, channel int) driver.CommandStatus {
    //head exists
    //channel exists
    //what does this do again? probably need to make sure it gets called appropriately
    panic("unimplemented")
}

//AddPlateTo - used
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) 
                                  driver.CommandStatus {
    //position exists
    //position can accept a plate of this type
    //plate can be cast to LHPlate. plate type mathes position type
    panic("unimplemented")
}

//RemoveAllPlates - used
func (self *VirtualLiquidHandler) RemoveAllPlates() driver.CommandStatus {
    //remove plates, no checks required.
    panic("unimplemented")
}

//RemovePlateAt - unused
func (self *VirtualLiquidHandler) RemovePlateAt(position string) driver.CommandStatus {
    //plate exists at position
    panic("unimplemented")
}

//SetPositionState - unused
func (self *VirtualLiquidHandler) SetPositionState(position string, state driver.PositionState) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//GetCapabilites - used
func (self *VirtualLiquidHandler) GetCapabilities() (LHProperties, driver.CommandStatus) {
    //no checks requireds
    return (self.properties, driver.CommandStatus{true, driver.OK, ""}) 
}

//GetCurrentPosition - unused
func (self *VirtualLiquidHandler) GetCurrentPosition(head int) (string, driver.CommandStatus) {
    panic("unimplemented")
}

//GetPositionState - unused
func (self *VirtualLiquidHandler) GetPositionState(position string) (string, driver.CommandStatus) {
    panic("unimplemented")
}

//GetHeadState - unused
func (self *VirtualLiquidHandler) GetHeadState(head int) (string, driver.CommandStatus) {
    panic("unimplemented")
}

//GetStatus - unused
func (self *VirtualLiquidHandler) GetStatus() (driver.Status, driver.CommandStatus) {
    panic("unimplemented")
}

//UpdateMetaData - used
func (self *VirtualLiquidHandler) UpdateMetaData(props *LHProperties) driver.CommandStatus {
    panic("unimplemented")
}

//UnloadHead - unused
func (self *VirtualLiquidHandler) UnloadHead(param int) driver.CommandStatus {
    panic("unimplemented")
}

//LoadHead - unused
func (self *VirtualLiquidHandler) LoadHead(param int) driver.CommandStatus {
    panic("unimplemented")
}

//Lights On - not implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOn() driver.CommandStatus {
    panic("unimplemented")
}

//Lights Off - notimplemented in compositerobotinstruction
func (self *VirtualLiquidHandler) LightsOff() driver.CommandStatus {
    panic("unimplemented")
}

//LoadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) LoadAdaptor(param int) driver.CommandStatus {
    panic("unimplemented")
}

//UnloadAdaptor - notimplemented in CRI
func (self *VirtualLiquidHandler) UnloadAdaptor(param int) driver.CommandStatus {
    panic("unimplemented")
}

//Open - notimplemented in CRI
func (self *VirtualLiquidHandler) Open() driver.CommandStatus {
    panic("unimplemented")
}

//Close - notimplement in CRI
func (self *VirtualLiquidHandler) Close() driver.CommandStatus {
    panic("unimplemented")
}

//Message - unused
func (self *VirtualLiquidHandler) Message(level int, title, text string, showcancel bool) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//GetOutputFile - used, but not in instruction stream
func (self *VirtualLiquidHandler) GetOutputFile() (string, driver.CommandStatus) {
    //Probably won't get called on the simulator just yet...
    panic("unimplemented")
}



