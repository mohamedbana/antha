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
    "io/ioutil"
    "strings"
    "fmt"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
)

// Simulate a liquid handler Driver
type VirtualLiquidHandler struct {
    properties *liquidhandling.LHProperties 
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

    log []string
}

//Create a new VirtualLiquidHandler which mimics an LHDriver
func NewVirtualLiquidHandler(props *liquidhandling.LHProperties) (*VirtualLiquidHandler, error) {
    var vlh VirtualLiquidHandler

    vlh.properties = props.Dup()
    vlh.log = make([]string, 0)

    return &vlh, nil
}

//Write to the log
func (self *VirtualLiquidHandler) LogLine(line string) {
    self.log = append(self.log, line)
}

//save the log
func (self *VirtualLiquidHandler) SaveLog(filename string) {
    ioutil.WriteFile(filename, []byte(strings.Join(self.log, "\n")), 0644)
}

//Move command - used
func (self *VirtualLiquidHandler) Move(deckposition []string, wellcoords []string, reference []int, 
                                       offsetX, offsetY, offsetZ []float64, plate_type []string, 
                                       head int) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`Move(
    deckposition = %v,
    wellcoords = %v,
    reference = %v,
    offsetX,Y,Z = (%v, %v, %v),
    plate_type = %v,
    head = %v)`, deckposition, wellcoords, reference, offsetX, offsetY, offsetZ, plate_type, head))
    //Asserts:
    //deckposition exists - why is it a list?
    //wellcoords exists within the plate at deckposition
    //reference is in allowable range
    //offsetX, offsetY, offsetZ are within the well (but I guess they needn't be...)
    //plate_type matches the type of plate at deckposition
    //head is valid
    return driver.CommandStatus{true, driver.OK, "MOVE ACK"}
}

//Move raw - not yet implemented in compositerobotinstruction
func (self *VirtualLiquidHandler) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`MoveRaw(
    head = %v,
    offsetX,Y,Z = (%v, %v, %v))`, head, x,y,z))
    //Asserts:
    //head exists
    //x,y,x are within the machine
    return driver.CommandStatus{true, driver.OK, "MOVERAW ACK"}
}

//Aspirate - used
func (self *VirtualLiquidHandler) Aspirate(volume []float64, overstroke []bool, head int, multi int, 
                                           platetype []string, what []string, llf []bool) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`Aspirate(
    volume = %v,
    overstroke = %v,
    head = %v,
    multi = %v,
    platetype = %v,
    what = %v,
    llf = %v)`, volume, overstroke, head, multi, platetype, what, llf))
    //volumes are equal if adapter isn't independent
    //tips are loaded in each adapter location the aspirates
    //volume is smaller that the tips' maximum capacity
    //the tip hasn't been used for a different liquid
    //head exists
    //multi matches number of tips loaded
    //platetype matches the plate at the location we moved to
    //what matches the expected liquid class
    //llf is the right size, cannot vary unless independent
    return driver.CommandStatus{true, driver.OK, "ASPIRATE ACK"}
}

//Dispense - used
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int, 
                                           platetype []string, what []string, llf []bool) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`Dispense(
    volume = %v,
    blowout = %v,
    head = %v,
    multi = %v,
    platetype = %v,
    what = %v,
    llf = %v)`, volume, blowout, head, multi, platetype, what, llf))
    //Volumes are equal if adapter isn't indepentent
    //Volumes are at most equal to the volume in the tip
    //blowout is the right length
    //head exists
    //multi is valid
    //platetype matches the type of plate that we're next to
    //what matches the liquid class that was aspirated
    //llf is the right size and follows independence constraint
    return driver.CommandStatus{true, driver.OK, "DISPENSE ACK"}
}

//LoadTips - used
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int, 
                                           platetype, position, well []string) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`LoadTips(
    channels = %v,
    head = %v,
    multi = %v,
    platetype = %v,
    position = %v,
    well = %v)`, channels, head, multi, platetype, position, well))
    //channels is correct length and value (how does this work)
    //head exists
    //multi is in correct range for adaptor
    //platetype matches the plate we're over
    //position is correct, tips still exists there
    //well exists (difference between platetype and well?)
    return driver.CommandStatus{true, driver.OK, "LOADTIPS ACK"}
}

//UnloadTips - used
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int, 
                                             platetype, position, well []string) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`UnloadTips(
    channels = %v,
    head = %v,
    multi = %v,
    platetype = %v,
    position = %v,
    well = %v)`, channels, head, multi, platetype, position, well))
    //Tips are loaded in channels
    //independence constraints are met
    //head exists
    //multi is correct
    //platetype matches the plate we're over
    //platetype is tip-waste
    //position and well are correct
    return driver.CommandStatus{true, driver.OK, "UNLOADTIPS ACK"}
}

//SetPipetteSpeed - used
func (self *VirtualLiquidHandler) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`SetPipetteSpeed(
    head = %v,
    channel = %v,
    rate = %v)`, head, channel, rate))
    //head exists
    //channel exists
    //speed is within allowable range
    return driver.CommandStatus{true, driver.OK, "SETPIPETTESPEED ACK"}
}

//SetDriveSpeed - used
func (self *VirtualLiquidHandler) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`SetDriveSpeed(
    drive = %v,
    rate = %v)`, drive, rate))
    //drive string?
    //rate is within allowable range (what is this?)
    return driver.CommandStatus{true, driver.OK, "SETDRIVESPEED ACK"}
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
    self.LogLine("Initialize()")
    //check that this is called before anything else?
    return driver.CommandStatus{true, driver.OK, "INITIALIZE ACK"}
}

//Finalize - used
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
    self.LogLine("Finalize()")
    //check that this is called last, no more calls
    return driver.CommandStatus{true, driver.OK, "FINALIZE ACK"}
}

//Wait - used
func (self *VirtualLiquidHandler) Wait(time float64) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`Wait(time = %v)`, time))
    //time is positive
    //maybe a warning if it's super-long
    return driver.CommandStatus{true, driver.OK, "WAIT ACK"}
}

//Mix - used
func (self *VirtualLiquidHandler) Mix(head int, volume []float64, platetype []string, cycles []int, 
                                      multi int, what []string, blowout []bool) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`Mix(
    head = %v,
    volume = %v,
    platetype = %v,
    cycles = %v,
    multi = %v,
    what = %v,
    blowout = %v)`, head, volume, platetype, cycles, multi, what, blowout))
    //head exists
    //volume is lte volume in wells
    //platetype matches the plate we're over
    //muli is correct
    //what matches expected liquidclass
    //volume, platetype, what, blowout match independence constraint
    return driver.CommandStatus{true, driver.OK, "MIX ACK"}
}

//ResetPistons - used
func (self *VirtualLiquidHandler) ResetPistons(head, channel int) driver.CommandStatus {
    self.LogLine("ResetPistons()")
    //head exists
    //channel exists
    //what does this do again? probably need to make sure it gets called appropriately
    return driver.CommandStatus{true, driver.OK, "RESETPISTONS ACK"}
}

//AddPlateTo - used
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) driver.CommandStatus {
    self.LogLine(fmt.Sprintf(`AddPlateTo(
    position = %v,
    plate = %v,
    name = %v)`, position, plate, name))
    //position exists
    //position can accept a plate of this type
    //plate can be cast to LHPlate. plate type mathes position type
    return driver.CommandStatus{true, driver.OK, "ADDPLATETO ACK"}
}

//RemoveAllPlates - used
func (self *VirtualLiquidHandler) RemoveAllPlates() driver.CommandStatus {
    self.LogLine("RemoveAllPlates()")
    //remove plates, no checks required.
    return driver.CommandStatus{true, driver.OK, "REMOVEALLPLATES ACK"}
}

//RemovePlateAt - unused
func (self *VirtualLiquidHandler) RemovePlateAt(position string) driver.CommandStatus {
    self.LogLine(fmt.Sprintf("RemovePlateAt(position = %v)", position))
    //plate exists at position
    return driver.CommandStatus{true, driver.OK, "REMOVEPLATEAT ACK"}
}

//SetPositionState - unused
func (self *VirtualLiquidHandler) SetPositionState(position string, state driver.PositionState) driver.CommandStatus {
    panic("unimplemented")
}

//GetCapabilites - used
func (self *VirtualLiquidHandler) GetCapabilities() (liquidhandling.LHProperties, driver.CommandStatus) {
    self.LogLine("GetCapabilities()")
    //no checks required
    return *self.properties, driver.CommandStatus{true, driver.OK, ""} 
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
func (self *VirtualLiquidHandler) UpdateMetaData(props *liquidhandling.LHProperties) driver.CommandStatus {
    self.LogLine("ResetPistons(props *LHProperties)")
    //check that the props and self.props are the same...
    return driver.CommandStatus{true, driver.OK, "UPDATEMETADATA ACK"}
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
func (self *VirtualLiquidHandler) Message(level int, title, text string, showcancel bool) driver.CommandStatus {
    panic("unimplemented")
}

//GetOutputFile - used, but not in instruction stream
func (self *VirtualLiquidHandler) GetOutputFile() (string, driver.CommandStatus) {
    //Probably won't get called on the simulator just yet...
    panic("unimplemented")
}



