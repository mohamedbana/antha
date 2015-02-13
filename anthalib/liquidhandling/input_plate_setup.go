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

package liquidhandling // import "github.com/antha-lang/antha/anthalib/liquidhandling"

import (
	"github.com/antha-lang/antha/anthalib/wutil"
	"fmt"
	"errors"
)

//  TASK: 	Map inputs to input plates 
// INPUT: 	"input_platetype", "inputs"
//OUTPUT: 	"input_plates"      -- these each have components in wells
//		"input_assignments" -- map with arrays of assignment strings, i.e. {tea: [plate1:A:1, plate1:A:2...] }etc.
func input_plate_setup(request *LHRequest)(map[string]*LHPlate, map[string][]string){
	input_platetype:=(*request).Input_platetype
	if(input_platetype.ID==""){
		wutil.Error(errors.New("plate_setup: No input plate type defined"))
	}
	input_plates:=(*request).Input_plates

	if(len(input_plates)==0){
		input_plates=make(map[string]*LHPlate, 3)
	}


	// need to fill each plate type

	var curr_plate *LHPlate

	curr_plate=new_plate(input_platetype)
	inputs:=(*request).Input_solutions
	input_assignments:=make(map[string][]string, len(inputs))

	for name,input:=range inputs{
		fmt.Println("Plate_setup - component", name, ":")
		var curr_well *LHWell
		var ok bool
		ass:=make([]string, 0, 3)

		for _,component:=range input{
			// find somewhere to put it
			curr_well,ok=get_next_well(curr_plate, component, curr_well)

			if(!ok){
				curr_plate=new_plate(input_platetype)
				curr_well,ok=get_next_well(curr_plate, component, nil)
			}

			// now put it there

			contents:=curr_well.Contents
			if(len(contents)==0){
				contents=make([]*LHComponent, 0, 4)
			}

			location:=curr_plate.ID+":"+curr_well.Coords
			ass=append(ass, location)

			component.Loc=location

			contents=append(contents, component)

			curr_well.Contents=contents
			input_plates[curr_plate.ID]=curr_plate
		}
		input_assignments[name]=ass
	}

	return input_plates, input_assignments
}