// anthalib//liquidhandling/input_plate_setup.go: Part of the Antha language
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
	//"fmt"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/factory"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

//  TASK: 	Map inputs to input plates
// INPUT: 	"input_platetype", "inputs"
//OUTPUT: 	"input_plates"      -- these each have components in wells
//		"input_assignments" -- map with arrays of assignment strings, i.e. {tea: [plate1:A:1, plate1:A:2...] }etc.
func input_plate_setup(request *LHRequest) *LHRequest {
	input_platetypes := (*request).Input_platetypes
	if input_platetypes == nil || len(input_platetypes) == 0 {
		// this configuration needs to happen outside but for now...
		list := factory.GetPlateList()
		input_platetypes = make([]*wtype.LHPlate, len(list))
		for i, platetype := range list {
			input_platetypes[i] = factory.GetPlateByType(platetype)
		}
		(*request).Input_platetypes = input_platetypes
	}
	input_plates := (*request).Input_plates

	if len(input_plates) == 0 {
		input_plates = make(map[string]*wtype.LHPlate, 3)
	}

	// need to fill each plate type

	var curr_plate *wtype.LHPlate

	inputs := (*request).Input_solutions

	input_volumes := make(map[string]wunit.Volume, len(inputs))

	// aggregate the volumes for the inputs
	for k, v := range inputs {
		v2 := v[0].Volume()
		vol := &v2
		for i := 1; i < len(v); i++ {
			vv := v[i].Volume()
			vol.Add(&vv)
		}
		input_volumes[k] = *vol
		//fmt.Println("TOTAL Volume for ", k, " : ", vol.ToString())
	}

	weights_constraints := request.Input_Setup_Weights

	// get the assignments

	well_count_assignments := choose_plate_assignments(input_volumes, input_platetypes, weights_constraints)

	input_assignments := make(map[string][]string, len(well_count_assignments))

	plates_in_play := make(map[string]*wtype.LHPlate)

	curplaten := 1
	for cname, volume := range input_volumes {
		component := inputs[cname][0]
		//fmt.Println("Plate_setup - component", cname, ":")

		well_assignments := well_count_assignments[cname]

		//fmt.Println("Well assignments: ", well_assignments)

		var curr_well *wtype.LHWell
		var ok bool
		ass := make([]string, 0, 3)

		for platetype, nwells := range well_assignments {
			for i := 0; i < nwells; i++ {
				curr_plate = plates_in_play[platetype.Type]

				if curr_plate == nil {
					plates_in_play[platetype.Type] = factory.GetPlateByType(platetype.Type)
					curr_plate = plates_in_play[platetype.Type]
					platename := fmt.Sprintf("Input_plate_%d", curplaten)
					curr_plate.PlateName = platename
					curplaten += 1
				}

				// find somewhere to put it
				curr_well, ok = wtype.Get_Next_Well(curr_plate, component, curr_well)

				if !ok {
					plates_in_play[platetype.Type] = factory.GetPlateByType(platetype.Type)
					curr_well, ok = wtype.Get_Next_Well(curr_plate, component, nil)
				}

				// now put it there

				contents := curr_well.WContents
				if len(contents) == 0 {
					contents = make([]*wtype.LHComponent, 0, 4)
				}

				location := curr_plate.ID + ":" + curr_well.Crds
				ass = append(ass, location)

				// make a duplicate of this component to stick in the well

				newcomponent := component.Dup()
				newcomponent.Vol = curr_well.Vol
				newcomponent.Loc = location
				volume.Subtract(curr_well.WorkingVolume())

				contents = append(contents, newcomponent)

				curr_well.WContents = contents
				curr_well.Currvol = newcomponent.Vol
				input_plates[curr_plate.ID] = curr_plate
			}
		}
		input_assignments[cname] = ass
		//fmt.Println("ASSIGNMENT: ", cname, " ", ass)
	}

	(*request).Input_plates = input_plates
	(*request).Input_assignments = input_assignments
	//return input_plates, input_assignments
	return request
}
