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

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

func AdvancedExecutionPlanner2(request *LHRequest, parameters *liquidhandling.LHProperties) *LHRequest {
	// in the first instance we assume this is done component-wise
	// we also need to identify dependencies, i.e. if certain components
	// are only available after other actions
	// this will only work if the components to be added are to go in the same ordero

	// IT'S THAT HIDEOUS HACK AGAIN
	volume_correction := 0.5

	// get the layout groups

	minorlayoutgroups := request.Output_minor_group_layouts
	ass := request.Output_assignments

	// sort them, we might want to record the acutal order somewhere... also to allow user configuration of this

	minorlayoutgroups, ass = sortOutputOrder(minorlayoutgroups, ass, COLWISE)

	inass := request.Input_assignments
	output_solutions := request.Output_solutions
	input_plates := copyplates(request.Input_plates)
	output_plate_layout := request.Output_plate_layout
	output_plates := copyplates(request.Output_plates)
	plate_lookup := request.Plate_lookup

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

	cnt := 1
	logger.Debug("OUTORDER INFO STARTS HERE")
	for n, g := range minorlayoutgroups {

		grp := []string(g)

		// get the group assignment string
		assignment := ass[n]
		for _, solID := range grp {
			sol := output_solutions[solID]
			whats := make([]string, len(sol.Components))
			pltfrom := make([]string, len(sol.Components))
			pltto := make([]string, len(sol.Components))
			plttypefrom := make([]string, len(sol.Components))
			plttypeto := make([]string, len(sol.Components))
			wellfrom := make([]string, len(sol.Components))
			wellto := make([]string, len(sol.Components))
			vols := make([]*wunit.Volume, len(sol.Components))
			fvols := make([]*wunit.Volume, len(sol.Components))
			tvols := make([]*wunit.Volume, len(sol.Components))
			fpwx := make([]int, len(sol.Components))
			fpwy := make([]int, len(sol.Components))
			tpwx := make([]int, len(sol.Components))
			tpwy := make([]int, len(sol.Components))

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
			logger.Debug(fmt.Sprintf("OUTORDER:%d:%d:%s:%d", cnt, toplatenum, asstx[1], col))

			i := 0
			for ordinal, name := range order {
				logger.Debug(fmt.Sprintf("Component %s EXECUTE OUTORDER %d", name, ordinal+1))

				// we need to get the relevant component out
				smpl := get_aggregate_component(sol, name)
				if smpl == nil {
					//	row += incrow
					//	col += inccol
					continue
				}

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

				outwell.Add(smpl)

				// update the output solution with its location

				sol.Plateaddress = outplate.PlateName
				sol.PlateID = outplate.ID
				sol.Welladdress = wellto[i]
				cnt += 1

				// if we get here without finding any components of this type in this group we don't make an instruction

				i += 1
			}
			ins := liquidhandling.NewTransferInstruction(whats, pltfrom, pltto, wellfrom, wellto, plttypefrom, plttypeto, vols, fvols, tvols, fpwx, fpwy, tpwx, tpwy)

			instructions.Add(ins)
			row += incrow
			col += inccol
		}
	}

	logger.Debug("OUTORDER info ends here")
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
