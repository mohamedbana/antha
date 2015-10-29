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
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"errors"
	"fmt"

	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

func AdvancedExecutionPlanner(request *LHRequest, parameters *liquidhandling.LHProperties) *LHRequest {
	// in the first instance we assume this is done component-wise
	// we also need to identify dependencies, i.e. if certain components
	// are only available after other actions
	// this will only work if the components to be added are to go in the same ordero

	// IT'S THAT HIDEOUS HACK AGAIN
	volume_correction := 0.5

	// get the layout groups

	minorlayoutgroups := request.Output_minor_group_layouts
	ass := request.Output_assignments
	inass := request.Input_assignments
	output_solutions := request.Output_solutions
	input_plates := copyplates(request.Input_plates)
	output_plate_layout := request.Output_plate_layout
	output_plates := copyplates(request.Output_plates)
	plate_lookup := request.Plate_lookup

	// highly inelegant... we should swap Tip Box Setup
	// around to take place after, then this whole thing
	// is a non-issue
	// more to the point this won't work in general and needs a big fix
	// we need to allow for more than one type of tip to be available
	//	tt := make([]*wtype.LHTip, 1)
	//	tt[0] = request.Tip_Type.Tiptype
	//	parameters.Tips = tt
	//TODO this has been removed because it is already handled on the equipmentManager microArch repo code.
	instructions := liquidhandling.NewRobotInstructionSet(nil)
	order := request.Input_order

	// this whole bit is just to get input well volumes sorted out

	for _, name := range order {
		for _, g := range minorlayoutgroups {
			grp := []string(g)
			for _, solID := range grp {
				sol := output_solutions[solID]

				// we need to get the relevant component out
				smpl := get_aggregate_component(sol, name)
				if smpl == nil {
					continue
				}

				// just for the side-effects, eesh
				inassignmentar := []string(inass[name])
				vol := smpl.Vol + volume_correction
				_, _, ok := get_assignment(inassignmentar, &input_plates, vol)

				if !ok {
					wutil.Error(errors.New(fmt.Sprintf("No input assignment for %s with vol %-4.1f", name, smpl.Vol)))
				}
			}
		}
	}
	ip := request.Input_plates

	for k, v := range ip {
		// p2 is the version we have just taken from
		p2 := input_plates[k]

		for i, row := range p2.Rows {
			for j, well := range row {
				well2 := v.Rows[i][j]

				if well.Currvol != well2.Currvol {
					// we add everything we took away and the rv

					vol := well2.Currvol - well.Currvol
					well2.Currvol = roundup(vol + well.Rvol)
				}
			}
		}

	}

	// now we generate instructions

	input_plates = copyplates(ip)

	// write out what's in the input plates now

	for _, p := range input_plates {
		for i := 0; i < p.WlsY; i++ {
			for j := 0; j < p.WlsX; j++ {
				if p.Rows[i][j].Currvol != 0.0 {
					s := fmt.Sprintf("Plate ", p.PlateName, " TYPE ", p.Type, " ID ", p.ID, " WELL ", p.Rows[i][j].Crds, " ")

					for _, cmp := range p.Rows[i][j].WContents {
						s += fmt.Sprintf(cmp.ID, " ", cmp.CName, " ", cmp.Vol, " ", cmp.Vunit)
					}

					logger.Info(s)
				}
			}
		}

	}

	for _, name := range order {
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
			fpwx := make([]int, len(grp))
			fpwy := make([]int, len(grp))
			tpwx := make([]int, len(grp))
			tpwy := make([]int, len(grp))

			compingroup := false

			for i, solID := range grp {
				sol := output_solutions[solID]

				// we need to get the relevant component out
				smpl := get_aggregate_component(sol, name)
				if smpl == nil {
					row += incrow
					col += inccol
					continue
				}
				// important: is there at least one component in the group?
				// yes if we are here
				compingroup = true

				// we need to know where this component was assigned to
				inassignmentar := []string(inass[name])

				vol := smpl.Vol + volume_correction

				inassignment, wvol, ok := get_assignment(inassignmentar, &input_plates, vol)

				if !ok {
					wutil.Error(errors.New(fmt.Sprintf("No input assignment for %s with vol %-4.1f", name, smpl.Vol)))
				}

				inasstx := strings.Split(inassignment, ":")

				inplt := inasstx[0]
				inrow := string(inasstx[1])
				incol := wutil.ParseInt(inasstx[2])

				// we can fill the structure now

				whats[i] = smpl.Type
				pltfrom[i] = plate_lookup[string(inplt)]
				plttypefrom[i] = input_plates[string(inplt)].Type
				pltto[i] = plate_lookup[output_plate_layout[toplatenum]]
				plttypeto[i] = output_plates[output_plate_layout[toplatenum]].Type
				wellfrom[i] = inrow + strconv.Itoa(incol)
				wellto[i] = wutil.NumToAlpha(row) + strconv.Itoa(col)

				outplate := output_plates[output_plate_layout[toplatenum]]

				outwell := outplate.Rows[row-1][col-1]
				v := wunit.NewVolume(smpl.Vol, smpl.Vunit)
				tpwx[i] = outplate.WellsX()
				tpwy[i] = outplate.WellsY()

				vt := wunit.NewVolume(outwell.Currvol, "ul")
				vf := wunit.NewVolume(wvol, "ul")
				vols[i] = &v
				fvols[i] = &vf
				tvols[i] = &vt

				inplate := input_plates[string(inplt)]
				fpwx[i] = inplate.WellsX()
				fpwy[i] = inplate.WellsY()

				row += incrow
				col += inccol
				outwell.Add(smpl)

				// update the output solution with its location

				sol.Plateaddress = outplate.PlateName
				sol.PlateID = outplate.ID
				sol.Welladdress = wellto[i]
			}

			// if we get here without finding any components of this type in this group we don't make an instruction

			if !compingroup {
				continue
			}

			ins := liquidhandling.NewTransferInstruction(whats, pltfrom, pltto, wellfrom, wellto, plttypefrom, plttypeto, vols, fvols, tvols, fpwx, fpwy, tpwx, tpwy)

			instructions.Add(ins)
		}
	}

	inx := instructions.Generate(request.Policies, parameters)
	instrx := make([]liquidhandling.TerminalRobotInstruction, len(inx))
	for i := 0; i < len(inx); i++ {
		instrx[i] = inx[i].(liquidhandling.TerminalRobotInstruction)
	}
	request.Instructions = instrx

	// write output destinations to the log

	for _, sol := range output_solutions {
		s := fmt.Sprintf("SOLUTION: %s (ID %s) in block %s mapped to output plate %s (ID %s) well %s", sol.SName, sol.ID, sol.BlockID, sol.Plateaddress, sol.PlateID, sol.Welladdress)
		logger.Info(s)
	}

	return request
}

func roundup(f float64) float64 {
	return float64(int(f) + 1)
}

func get_aggregate_component(sol *wtype.LHSolution, name string) *wtype.LHComponent {
	components := sol.Components

	ret := wtype.NewLHComponent()

	ret.CName = name

	vol := 0.0
	found := false

	for _, component := range components {
		nm := component.CName

		if nm == name {
			ret.Type = component.Type
			vol += component.Vol
			ret.Vunit = component.Vunit
			ret.Loc = component.Loc
			ret.Order = component.Order
			found = true
		}
	}
	if !found {
		return nil
	}
	ret.Vol = vol
	return ret
}

func get_assignment(assignments []string, plates *map[string]*wtype.LHPlate, vol float64) (string, float64, bool) {
	assignment := ""
	ok := false
	prevol := 0.0

	for _, assignment = range assignments {
		asstx := strings.Split(assignment, ":")
		plate := (*plates)[asstx[0]]

		crds := asstx[1] + ":" + asstx[2]
		wellidlkp := plate.Wellcoords
		well := wellidlkp[crds]

		currvol := well.Currvol - well.Rvol
		if currvol >= vol {
			prevol = well.Currvol
			well.Currvol -= vol
			plate.HWells[well.ID] = well
			(*plates)[asstx[0]] = plate
			ok = true
			break
		}
	}

	return assignment, prevol, ok
}

func copyplates(plts map[string]*wtype.LHPlate) map[string]*wtype.LHPlate {
	ret := make(map[string]*wtype.LHPlate, len(plts))

	for k, v := range plts {
		ret[k] = v.Dup()
	}

	return ret
}
