// anthalib//liquidhandling/basicrobotinstruction.go: Part of the Antha language
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

package liquidhandling

import "fmt"

const (
	ASP int = iota
	DSP
	MOV
	LOD
	ULD
	TFR
)

type RobotInstruction interface{
	InstructionType()int
	GetParameter(name string)interface{}
}

type AspirateInstruction struct{
	instructionType int
	Vol float64
	Volunit string
	Speed float64
	Speedunit string
	ComponentType string
}

func (ai AspirateInstruction)InstructionType()int{
	return ai.instructionType
}
func (ins AspirateInstruction)GetParameter(s string)interface{}{
	switch(s){
		case "VOLUME" : return ins.Vol
		case "SPEED"  : return ins.Speed
		case "VOLUNIT": return ins.Volunit
		case "SPDUNIT": return ins.Speedunit
		case "CTYPE"  : return ins.ComponentType
		default: raiseError(fmt.Sprintf("Aspirate: illegal parameter: %s", s))
	}
	return nil
}

type DispenseInstruction struct{
	instructionType int
	Vol float64
	Volunit string
	Speed float64
	Speedunit string
	ComponentType string
}

func (di DispenseInstruction)InstructionType()int{
	return di.instructionType
}
func (ins DispenseInstruction)GetParameter(s string)interface{}{
	switch(s){
		case "VOLUME" : return ins.Vol
		case "SPEED"  : return ins.Speed
		case "VOLUNIT": return ins.Volunit
		case "SPDUNIT": return ins.Speedunit
		case "CTYPE"  : return ins.ComponentType
		default: raiseError(fmt.Sprintf("Dispense: illegal parameter: %s", s))
	}
	return nil
}

type MoveInstruction struct{
	instructionType int
	Pos int
	Well string
	Height int
	OffsetX float64
	OffsetY float64
	OffsetZ float64
	ComponentType string
}

func (mi MoveInstruction)InstructionType()int{
	return mi.instructionType
}

func (ins MoveInstruction)GetParameter(s string)interface{}{
	switch(s){
		case "POSITION"    : return ins.Pos
		case "WELL"   : return ins.Well
		case "HEIGHT" : return ins.Height
		case "OFFSETX": return ins.OffsetX
		case "OFFSETY": return ins.OffsetY
		case "OFFSETZ": return ins.OffsetZ
		default: raiseError(fmt.Sprintf("Move: illegal parameter: %s", s))
	}
	return nil
}

type LoadInstruction struct{
	instructionType int
}

func (li LoadInstruction)InstructionType()int{
	return li.instructionType
}

func (ins LoadInstruction)GetParameter(s string) interface{}{
	raiseError(fmt.Sprintf("Load: illegal parameter: %s", s))
	// props to the Go compiler for forcing me to put in unreachable statements!
	// </snark>
	return nil
}

type UnloadInstruction struct{
	instructionType int
}

func (ui UnloadInstruction)InstructionType()int{
	return ui.instructionType
}

func (ins UnloadInstruction)GetParameter(s string)interface{}{
	raiseError(fmt.Sprintf("Unload: illegal parameter: %s", s))
	return nil
}

func Aspirate(vol float64, volunit string, speed float64, speedunit string, what string)AspirateInstruction{
	return AspirateInstruction{ASP, vol, volunit, speed, speedunit, what}
}

func Dispense(vol float64, volunit string, speed float64, speedunit string, what string)DispenseInstruction{
	return DispenseInstruction{DSP, vol, volunit, speed, speedunit, what}
}

func Move(pos int, well string, height int, offsetX, offsetY, offsetZ float64, what string)MoveInstruction{
	return MoveInstruction{MOV, pos, well, height, offsetX, offsetY, offsetZ, what}
}

func Load()LoadInstruction{
	return LoadInstruction{LOD}
}

func Unload()UnloadInstruction{
	return UnloadInstruction{ULD}
}