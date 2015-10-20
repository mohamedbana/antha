// microArch/equipment/action/types.go: Part of the Antha language
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

package action

//Action describes a particular function that an equipment can perform. It is the same concept as an interface, but
// since it is a representation of the real world, we cannot model it with an actual one
type Action int

const (
	NONE = iota
	MESSAGE
	LH_SETUP
	LH_MOVE
	LH_MOVE_EXPLICIT
	LH_MOVE_RAW
	LH_ASPIRATE
	LH_DISPENSE
	LH_LOAD_TIPS
	LH_UNLOAD_TIPS
	LH_SET_PIPPETE_SPEED
	LH_SET_DRIVE_SPEED
	LH_STOP
	LH_SET_POSITION_STATE
	LH_RESET_PISTONS
	LH_WAIT
	LH_MIX
	LH_ADD_PLATE
	LH_REMOVE_PLATE
	LH_REMOVE_ALL_PLATES //?? maybe not necessary
	LH_CONFIG
	LH_READ
	LH_END
	MLH_CHANGE_TIPS
	IN_INCUBATE
	IN_INCUBATE_SHAKE
)

func (a Action) String() string {
	switch a {
	case NONE:
		return "NONE"
	case LH_SETUP:
		return "SETUP"
	case LH_MOVE:
		return "MOVE"
	case LH_MOVE_EXPLICIT:
		return "MOVE EXPLICIT"
	case LH_MOVE_RAW:
		return "MOVE RAW"
	case LH_ASPIRATE:
		return "ASPIRATE"
	case LH_DISPENSE:
		return "DISPENSE"
	case LH_LOAD_TIPS:
		return "LOAD TIPS"
	case LH_UNLOAD_TIPS:
		return "UNLOAD TIPS"
	case LH_SET_PIPPETE_SPEED:
		return "SET PIPPETE SPEED"
	case LH_SET_DRIVE_SPEED:
		return "SET DRIVE SPEED"
	case LH_STOP:
		return "STOP"
	case LH_SET_POSITION_STATE:
		return "SET POSTION STATE"
	case LH_RESET_PISTONS:
		return "RESET PISTONS"
	case LH_WAIT:
		return "WAIT"
	case LH_MIX:
		return "MIX"
	case LH_ADD_PLATE:
		return "ADD PLATE"
	case LH_REMOVE_PLATE:
		return "REMOVE PLATE"
	case LH_REMOVE_ALL_PLATES:
		return "REMOVE ALL PLATES"
	case LH_CONFIG:
		return "CONFIG"
	case IN_INCUBATE:
		return "INCUBATE"
	case IN_INCUBATE_SHAKE:
		return "INCUBATE SHAKING"
	case MLH_CHANGE_TIPS:
		return "CHANGE TIPS"
	case LH_READ:
		return "READ"
	case LH_END:
		return "END"

	case MESSAGE:
		return "MESSAGE"
	default:
		return ""
	}
}
