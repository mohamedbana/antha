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
	"github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"strconv"
	"strings"
)

func AdvancedExecutionPlanner(request *LHRequest, parameters *liquidhandling.LHProperties) *LHRequest {
	// in the first instance we assume this is done component-wise
	// we also need to identify dependencies, i.e. if certain components
	// are only available after other actions
	// this will only work if the components to be added are to go in the same order

	// get the layout groups

	minorlayoutgroups := request.Output_minor_group_layouts
	ass := request.Output_assignments
	inass := request.Input_assignments
	output_solutions := request.Output_solutions
	input_plates := request.Input_plates
	output_plate_layout := request.Output_plate_layout

	plate_lookup := request.Plate_lookup

	// highly inelegant... we should swap Tip Box Setup
	// around to take place after, then this whole thing
	// is a non-issue
	tt := make([]*wtype.LHTip, 1)
	tt[0] = request.Tip_Type.Tiptype
	parameters.Tips = tt

	instructions := liquidhandling.NewRobotInstructionSet(nil)
	// need to deal with solutions

	order := request.Input_order

	for _, name := range order {
		//fmt.Println(name)

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
			toplatenum := wutil.ParseInt(plate)
			row := wutil.AlphaToNum(asstx[1])
			col := wutil.ParseInt(asstx[2])
			incrow := wutil.ParseInt(asstx[3])
			inccol := wutil.ParseInt(asstx[4])

			whats := make([]string, len(grp))
			pltfrom := make([]string, len(grp))
			pltto := make([]string, len(grp))
			plttypefrom := make([]string, len(grp))
			plttypeto := make([]string, len(grp))
			wellfrom := make([]string, len(grp))
			wellto := make([]string, len(grp))
			vols := make([]*wunit.Volume, len(grp))
			fvols := make([]*wunit.Volume, len(grp))
			tvols := make([]*wunit.Volume, len(grp))
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
				pltfrom[i] = plate_lookup[string(inplt)]
				pltto[i] = plate_lookup[output_plate_layout[toplatenum]]
				wellfrom[i] = inrow + strconv.Itoa(incol)
				wellto[i] = wutil.NumToAlpha(row) + strconv.Itoa(col)
				v := wunit.NewVolume(smpl.Vol, smpl.Vunit)
				v2 := wunit.NewVolume(0.0, "ul")
				vols[i] = &v
				// TODO Get the proper volumes here
				fvols[i] = &v2
				tvols[i] = &v2
				row += incrow
				col += inccol
			}

			ins := liquidhandling.NewTransferInstruction(whats, pltfrom, pltto, wellfrom, wellto, plttypefrom, plttypeto, vols, fvols, tvols /*, parameters.Cnfvol*/)
			instructions.Add(ins)
		}
	}

	inx := instructions.Generate(request.Policies, parameters)
	instrx := make([]liquidhandling.TerminalRobotInstruction, len(inx))
	for i := 0; i < len(inx); i++ {
		instrx[i] = inx[i].(liquidhandling.TerminalRobotInstruction)
	}
	request.Instructions = instrx

	return request
}

func get_aggregate_component(sol *wtype.LHSolution, name string) *wtype.LHComponent {
	components := sol.Components

	ret := wtype.NewLHComponent()

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

func get_assignment(assignments []string, plates *map[string]*wtype.LHPlate, vol float64) (string, bool) {
	assignment := ""
	ok := false

	for _, assignment = range assignments {
		asstx := strings.Split(assignment, ":")
		plate := (*plates)[asstx[0]]

		crds := asstx[1] + ":" + asstx[2]
		wellidlkp := plate.Wellcoords
		well := wellidlkp[crds]

		currvol := well.Currvol - well.Rvol
		if currvol >= vol {
			vol += well.Currvol
			well.Currvol -= vol
			plate.HWells[well.ID] = well
			(*plates)[asstx[0]] = plate
			ok = true
			break
		}
	}

	return assignment, ok
}
