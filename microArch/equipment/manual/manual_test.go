// /equipment/manual/manual_test.go: Part of the Antha language
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

package manual

import (
	"testing"

	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
)

func TestCapabilities(t *testing.T) {
	anthaManual := NewAnthaManualCUI("")
	expectedBehaviours := [...]equipment.ActionDescription{
		*equipment.NewActionDescription(action.MESSAGE, "", nil),
		*equipment.NewActionDescription(action.LH_SETUP, "", nil),
		*equipment.NewActionDescription(action.LH_MOVE, "", nil),
		*equipment.NewActionDescription(action.LH_MOVE_EXPLICIT, "", nil),
		*equipment.NewActionDescription(action.LH_MOVE_RAW, "", nil),
		*equipment.NewActionDescription(action.LH_ASPIRATE, "", nil),
		*equipment.NewActionDescription(action.LH_DISPENSE, "", nil),
		*equipment.NewActionDescription(action.LH_LOAD_TIPS, "", nil),
		*equipment.NewActionDescription(action.LH_UNLOAD_TIPS, "", nil),
		*equipment.NewActionDescription(action.LH_SET_PIPPETE_SPEED, "", nil),
		*equipment.NewActionDescription(action.LH_SET_DRIVE_SPEED, "", nil),
		*equipment.NewActionDescription(action.LH_STOP, "", nil),
		*equipment.NewActionDescription(action.LH_SET_POSITION_STATE, "", nil),
		*equipment.NewActionDescription(action.LH_RESET_PISTONS, "", nil),
		*equipment.NewActionDescription(action.LH_WAIT, "", nil),
		*equipment.NewActionDescription(action.LH_MIX, "", nil),
		*equipment.NewActionDescription(action.IN_INCUBATE, "", nil),
		*equipment.NewActionDescription(action.IN_INCUBATE_SHAKE, "", nil),
		*equipment.NewActionDescription(action.MLH_CHANGE_TIPS, "", nil),
	}

	for _, b := range expectedBehaviours {
		if anthaManual.Can(b) == false {
			t.Fatal("anthaManual expected behaviour %s not fulfilled", b)
		}
	}
}
