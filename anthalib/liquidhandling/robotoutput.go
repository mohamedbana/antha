// anthalib//liquidhandling/robotoutput.go: Part of the Antha language
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

import (
	"fmt"
	"regexp"
	"strings"
)

type RobotOutputInterface struct {
	InstructionOutputs []string
}

func NewOutputInterface(filename string) RobotOutputInterface {
	roi := RobotOutputInterface{make([]string, 10, 10)}

	roi.InstructionOutputs[0] = "Aspirate volume %-6.2f_VOLUME%"
	roi.InstructionOutputs[1] = "Dispense volume %-6.2f_VOLUME%"
	roi.InstructionOutputs[2] = "Move to position %d_POSITION% well %s_WELL% height %d_HEIGHT% Offsets: %-6.2f_OFFSETX% X %-6.2f_OFFSETY% Y %-6.2f_OFFSETZ% Z"
	roi.InstructionOutputs[3] = "Load tips"
	roi.InstructionOutputs[4] = "Unload tips"
	roi.InstructionOutputs[5] = "Should not be seen"

	return roi
}

func (self RobotOutputInterface) Output(ins RobotInstruction) string {
	// we get the appropriate output type
	s := self.InstructionOutputs[ins.InstructionType()]
	// then change the meta string into the appropriate one
	return self.ReplacePlaceholders(s, ins)
}

func (self RobotOutputInterface) ReplacePlaceholders(s string, ins RobotInstruction) string {
	rx, _ := regexp.Compile("%-?(\\d+(\\.\\d)?)?[defgs]_[A-Za-z]+%")
	loc := rx.FindIndex([]byte(s))

	if loc == nil {
		return s
	}

	// we have a match

	match := s[loc[0] : loc[1]-1]
	pre := s[0:loc[0]]
	post := s[loc[1]:len(s)]

	// make the replacement
	var replacement string
	tx := strings.Split(match, "_")

	/*
		switch(tx[1]){
			case "VOLUME" : replacement=fmt.Sprintf(tx[0], ins.Vol)
			case "SPEED"  : replacement=fmt.Sprintf(tx[0], ins.Speed)
			case "POS"    : replacement=fmt.Sprintf(tx[0], ins.Pos)
			case "WELL"   : replacement=fmt.Sprintf(tx[0], ins.Well)
			case "HEIGHT" : replacement=fmt.Sprintf(tx[0], ins.Height)
			case "OFFSETX": replacement=fmt.Sprintf(tx[0], ins.OffsetX)
			case "OFFSETY": replacement=fmt.Sprintf(tx[0], ins.OffsetY)
			case "OFFSETZ": replacement=fmt.Sprintf(tx[0], ins.OffsetZ)
			default: raiseError(fmt.Sprintf("Illegal parameter: %s", tx[1]))
		}
	*/

	// simpler way

	replacement = fmt.Sprintf(tx[0], ins.GetParameter(tx[1]))

	// rebuild the string
	s = pre + replacement + post

	// now do the next one
	return self.ReplacePlaceholders(s, ins)
}
