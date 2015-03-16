// liquidhandling/setupagent.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

import "strconv"

// default setup agent
func BasicSetupAgent(request *LHRequest, params *LHProperties) *LHRequest {
	// this is quite tricky and requires extensive interaction with the liquid handling
	// parameters

	// the principal question is how to define constraints on the system

	// I think this needs to remain tbd for now
	// instead we can rely on the preference system I already use

	plate_lookup := make(map[string]int, 5)
	tip_lookup := make([]*LHTipbox, 0, 5)

	tip_preferences := params.Tip_preferences
	input_preferences := params.Input_preferences
	output_preferences := params.Output_preferences

	// how do we set the below?
	// we don't know how many tips we need until we generate
	// instructions; ditto input or output plates until we've done layout

	// input plates
	input_plates := request.Input_plates

	// output plates
	output_plates := request.Output_plates

	// tips
	tips := request.Tips

	// we put tips on first

	setup := request.Setup

	if len(setup) == 0 {
		setup = NewLHSetup()
	}

	// need preference lists for the lh
	// if these aren't defined we have to do some kind of sensible
	// default thing
	// which is?

	for _, tb := range tips {
		// get the first available position from the preferences
		pos := get_first_available_preference(tip_preferences, setup)
		if pos == -1 {
			RaiseError("No positions left for tipbox")
		}

		position := "position_" + strconv.Itoa(pos)
		setup[position] = tb
		plate_lookup[tb.ID] = pos
		tip_lookup = append(tip_lookup, tb)
	}

	setup["tip_lookup"] = tip_lookup

	// this logic may not transfer well but I expect that outputs are more constrained
	// than inputs for the simple reason that most output takes place to single wells
	// while input takes place from reservoirs

	// outputs

	for _, p := range output_plates {
		pos := get_first_available_preference(output_preferences, setup)
		if pos == -1 {
			RaiseError("No positions left for output")
		}
		position := "position_" + strconv.Itoa(pos)
		setup[position] = p
		plate_lookup[p.ID] = pos
	}

	// inputs

	for _, p := range input_plates {
		pos := get_first_available_preference(input_preferences, setup)
		if pos == -1 {
			RaiseError("No positions left for input")
		}
		position := "position_" + strconv.Itoa(pos)
		setup[position] = p
		plate_lookup[p.ID] = pos
	}

	request.Setup = setup
	request.Plate_lookup = plate_lookup
	return request
}

func get_first_available_preference(prefs []int, setup map[string]interface{}) int {
	for _, pref := range prefs {
		position := "position_" + string(pref)
		_, ok := setup[position]
		if !ok {
			return pref
		}
	}
	return -1
}
