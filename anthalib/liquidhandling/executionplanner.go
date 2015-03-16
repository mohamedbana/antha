// anthalib//liquidhandling/executionplanner.go: Part of the Antha language
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
	"errors"
	"fmt"
	"github.com/antha-lang/antha/anthalib/wutil"
	"strconv"
	"strings"
)

func AdvancedExecutionPlanner(request *LHRequest, parameters *LHProperties) *LHRequest {
	// in the first instance we assume this is done component-wise
	// we also need to identify dependencies, i.e. if certain components
	// are only available after other actions

	// get the layout groups

	minorlayoutgroups := request.Output_minor_group_layouts
	ass := request.Output_assignments
	inass := request.Input_assignments
	output_solutions := request.Output_solutions
	input_plates := request.Input_plates

	instructions := make([]RobotInstruction, 0, 10)
	// need to deal with solutions

	cmps := request.Input_solutions

	for name, _ := range cmps {
		fmt.Println(name)

		// cmpa holds a list of inputs required per destination
		// i.e. if destination X requires 15 ul of h2o this will be listed separately
		// at this point all of the requests must be for volumes,

		// we need a mapping from these to where they belong
		// do we?
		/*
			cmpa:=value.([]map[string]interface{})
			cmp_map:=make(map[string]interface{}, len(cmpa))
			for _,cmp:=range cmpa{
				srcid:=cmp.Srcid
				cmp_map[srcid]=cmp
			}
		*/

		for n, g := range minorlayoutgroups {
			grp := []string(g)
			// get the group assignment string

			assignment := ass[n]

			// the assignment has the format plateID:row:column:incrow:inccol
			// where inc defines how the next one is to be calculated
			// e.g. {GUID}:A:1:1:0
			// 	{GUID}:A:1:0:1

			asstx := strings.Split(assignment, ":")

			plate := asstx[0]
			row := wutil.AlphaToNum(asstx[1])
			col := wutil.ParseInt(asstx[2])
			incrow := wutil.ParseInt(asstx[3])
			inccol := wutil.ParseInt(asstx[4])

			whats := make([]string, len(grp))
			pltfrom := make([]string, len(grp))
			pltto := make([]string, len(grp))
			wellfrom := make([]string, len(grp))
			wellto := make([]string, len(grp))
			vols := make([]float64, len(grp))
			volunits := make([]string, len(grp))

			for i, solID := range grp {
				sol := output_solutions[solID]

				// we need to get the relevant component out
				smpl := get_aggregate_component(sol, name)

				// we need to know where this component was assigned to
				inassignmentar := []string(inass[name])
				inassignment, ok := get_assignment(inassignmentar, &input_plates, smpl.Vol)

				if !ok {
					wutil.Error(errors.New(fmt.Sprintf("No input assignment for %s with vol %-4.1f", name, smpl.Vol)))
				}

				inasstx := strings.Split(inassignment, ":")

				inplt := inasstx[0]
				inrow := string(inasstx[1])
				incol := wutil.ParseInt(inasstx[2])

				// we can fill the structure now

				whats[i] = name
				pltfrom[i] = string(inplt)
				pltto[i] = string(plate)
				wellfrom[i] = inrow + strconv.Itoa(incol)
				wellto[i] = wutil.NumToAlpha(row) + strconv.Itoa(col)
				vols[i] = smpl.Vol
				volunits[i] = smpl.Vunit
				row += incrow
				col += inccol
			}

			ins := Transfer(whats, pltfrom, pltto, wellfrom, wellto, vols, volunits, parameters.Cnfvol[0])
			instructions = append(instructions, ins)
		}
	}

	request.Instructions = instructions

	return request
}

func get_aggregate_component(sol *LHSolution, name string) *LHComponent {
	components := sol.Components

	ret := NewLHComponent()

	ret.CName = name

	vol := 0.0

	for _, component := range components {
		nm := component.CName

		if nm == name {
			ret.Type = component.Type
			vol += component.Vol
			ret.Vunit = component.Vunit
			ret.Loc = component.Loc
			ret.Order = component.Order
		}
	}

	ret.Vol = vol

	return ret
}

func get_assignment(assignments []string, plates *map[string]*LHPlate, vol float64) (string, bool) {
	assignment := ""
	ok := false

	for _, assignment = range assignments {
		asstx := strings.Split(assignment, ":")
		plate := (*plates)[asstx[0]]

		crds := asstx[1] + ":" + asstx[2]
		wellidlkp := plate.Wellcoords
		well := wellidlkp[crds]
		currvol := well.Currvol
		maxvol := well.Vol - well.Rvol
		if currvol+vol < maxvol {
			vol += well.Currvol
			well.Currvol = vol
			plate.HWells[well.ID] = well
			(*plates)[asstx[0]] = plate
			ok = true
			break
		}
	}

	return assignment, ok
}
