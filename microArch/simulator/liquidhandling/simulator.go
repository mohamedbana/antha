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

// Simulate a liquid handler
type VirtualLiquidHandler struct {
    properties LHProperties 
}

//Create a new VirtualLiquidHandler which mimics the given properties
func NewVirtualLiquidHandler(props LHProperties) (*VirtualLiquidHandler, error) {
    var vlh VirtualLiquidHandler

    vlh.properties = props

    return &vlh
}

//Move command
func (self *VirtualLiquidHandler) Move(deckposition []string, wellcoords []string, reference []int, 
                                       offsetX, offsetY, offsetZ []float64, plate_type []string, 
                                       head int) driver.CommandStatus {
    panic("unimplemented")
}

//Move raw
func (self *VirtualLiquidHandler) MoveRaw(head int, x, y, z float64) driver.CommandStatus {
    panic("unimplemented")
}

//Aspirate
func (self *VirtualLiquidHandler) Aspirate(volume []float64, overstroke []bool, head int, multi int, 
                                           platetype []string, what []string, llf []bool) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//Dispense
func (self *VirtualLiquidHandler) Dispense(volume []float64, blowout []bool, head int, multi int, 
                                           platetype []string, what []string, llf []bool) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//Load Tips
func (self *VirtualLiquidHandler) LoadTips(channels []int, head, multi int, 
                                           platetype, position, well []string) driver.CommandStatus {
    panic("unimplemented")
}

//Unload Tips
func (self *VirtualLiquidHandler) UnloadTips(channels []int, head, multi int, 
                                             platetype, position, well []string) driver.CommandStatus {
    panic("unimplemented")
}

//SetPipetteSpeed
func (self *VirtualLiquidHandler) SetPipetteSpeed(head, channel int, rate float64) driver.CommandStatus {
    panic("unimplemented")
}

//SetDriveSpeed
func (self *VirtualLiquidHandler) SetDriveSpeed(drive string, rate float64) driver.CommandStatus {
    panic("unimplemented")
}

//Stop
func (self *VirtualLiquidHandler) Stop() driver.CommandStatus {
    panic("unimplemented")
}

//Go
func (self *VirtualLiquidHandler) Go() driver.CommandStatus {
    panic("unimplemented")
}

//Init
func (self *VirtualLiquidHandler) Initialize() driver.CommandStatus {
    panic("unimplemented")
}

//Finalize
func (self *VirtualLiquidHandler) Finalize() driver.CommandStatus {
    panic("unimplemented")
}

//Wait
func (self *VirtualLiquidHandler) Wait(time float64) driver.CommandStatus {
    panic("unimplemented")
}

//Mix
func (self *VirtualLiquidHandler) Mix(head int, volume []float64, platetype []string, cycles []int, 
                                      multi int, what []string, blowout []bool) driver.CommandStatus {
    panic("unimplemented")
}

//ResetPistons
func (self *VirtualLiquidHandler) ResetPistons(head, channel int) driver.CommandStatus {
    panic("unimplemented")
}

//AddPlateTo
func (self *VirtualLiquidHandler) AddPlateTo(position string, plate interface{}, name string) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//RemoveAllPlates
func (self *VirtualLiquidHandler) RemoveAllPlates() driver.CommandStatus {
    panic("unimplemented")
}

//RemovePlateAt
func (self *VirtualLiquidHandler) RemovePlateAt(position string) driver.CommandStatus {
    panic("unimplemented")
}

//SetPositionState
func (self *VirtualLiquidHandler) SetPositionState(position string, state driver.PositionState) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//GetCapabilites
func (self *VirtualLiquidHandler) GetCapabilities() (LHProperties, driver.CommandStatus) {
    panic("unimplemented")
}

//GetCurrentPosition
func (self *VirtualLiquidHandler) GetCurrentPosition(head int) (string, driver.CommandStatus) {
    panic("unimplemented")
}

//GetPositionState
func (self *VirtualLiquidHandler) GetPositionState(position string) (string, driver.CommandStatus) {
    panic("unimplemented")
}

//GetHeadState
func (self *VirtualLiquidHandler) GetHeadState(head int) (string, driver.CommandStatus) {
    panic("unimplemented")
}

//GetStatus
func (self *VirtualLiquidHandler) GetStatus() (driver.Status, driver.CommandStatus) {
    panic("unimplemented")
}

//UpdateMetaData
func (self *VirtualLiquidHandler) UpdateMetaData(props *LHProperties) driver.CommandStatus {
    panic("unimplemented")
}

//UnloadHead
func (self *VirtualLiquidHandler) UnloadHead(param int) driver.CommandStatus {
    panic("unimplemented")
}

//LoadHead
func (self *VirtualLiquidHandler) LoadHead(param int) driver.CommandStatus {
    panic("unimplemented")
}

//Lights On
func (self *VirtualLiquidHandler) LightsOn() driver.CommandStatus {
    panic("unimplemented")
}

//Lights Off
func (self *VirtualLiquidHandler) LightsOff() driver.CommandStatus {
    panic("unimplemented")
}

//Load Adaptor
func (self *VirtualLiquidHandler) LoadAdaptor(param int) driver.CommandStatus {
    panic("unimplemented")
}

//Unload Adaptor
func (self *VirtualLiquidHandler) UnloadAdaptor(param int) driver.CommandStatus {
    panic("unimplemented")
}

//Open
func (self *VirtualLiquidHandler) Open() driver.CommandStatus {
    panic("unimplemented")
}

//Close
func (self *VirtualLiquidHandler) Close() driver.CommandStatus {
    panic("unimplemented")
}

//Message
func (self *VirtualLiquidHandler) Message(level int, title, text string, showcancel bool) 
                                  driver.CommandStatus {
    panic("unimplemented")
}

//GetOutputFile
func (self *VirtualLiquidHandler) GetOutputFile() (string, driver.CommandStatus) {
    panic("unimplemented")
}



