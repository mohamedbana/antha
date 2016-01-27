// liquidhandling/layoutagent.go: Part of the Antha language
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
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"math"
	"strconv"

	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

func BasicLayoutAgent(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest {
	plate := request.Output_platetypes[0]
	solutions := request.Output_solutions

	// get the incoming group IDs
	// the purpose of this check is to determine whether there
	// already exist assignments...this is quite tricky
	MajorLayoutGroupIDs, _ := getLayoutGroupsAndPlates(solutions)

	// check we have enough and assign more if necessary
	// needs to be done per plate
	n_plates_required := int(math.Floor(float64(len(solutions))/float64(plate.Nwells))) + 1

	if len(MajorLayoutGroupIDs) < n_plates_required {
		lmg := len(MajorLayoutGroupIDs)
		mx, ok := wutil.Max(MajorLayoutGroupIDs)
		if !ok {
			mx = -1
		}

		for i := mx + 1; i < (n_plates_required - lmg); i++ {
			MajorLayoutGroupIDs = append(MajorLayoutGroupIDs, i)
		}
	}

	// then tidy them up

	MajorLayoutGroupRanks := wutil.MakeRankedList(MajorLayoutGroupIDs)

	// now we need to map solutions to groups

	MajorLayoutGroups := make([][]string, len(MajorLayoutGroupIDs))

	// make the receptacles

	for _, i := range MajorLayoutGroupIDs {
		MajorLayoutGroups[MajorLayoutGroupRanks[i]] = make([]string, 0, 4)
	}

	// helper var

	max_major_group_size := plate.Nwells

	for _, soln := range solutions {
		id := soln.ID

		lg := soln.Majorlayoutgroup

		// -1 means unassigned
		if lg == -1 {
			lg = choose_major_layout_group(MajorLayoutGroups, max_major_group_size)
		}

		MajorLayoutGroups[MajorLayoutGroupRanks[lg]] = append(MajorLayoutGroups[MajorLayoutGroupRanks[lg]], id)
	}

	plateLayouts := do_major_layouts(request, MajorLayoutGroups)
	request.Output_plate_layout = plateLayouts

	// now we need to set the minor layout groups attribute in request
	// in this instance this is just mapping everything to columns
	// this needs to work when outputs are pre-assigned

	minor_group_layouts := make([][]string, 0, len(solutions))
	assignments := make([]string, 0, len(solutions))

	for i, grp := range MajorLayoutGroups {
		dplate := plateLayouts[i]

		plate_minor_groups, plate_assignments := assign_minor_layouts(grp, plate, dplate, solutions)

		minor_group_layouts = append(minor_group_layouts, plate_minor_groups...)
		for _, as := range plate_assignments {
			assignments = append(assignments, as)
		}
	}

	request.Output_minor_group_layouts = minor_group_layouts
	request.Output_major_group_layouts = MajorLayoutGroups
	request.Output_assignments = assignments
	return request
}

func make_get_assignations(solutions map[string]*wtype.LHSolution, WlsX, WlsY int) ([][]string, bool) {
	assigned := make([][]string, WlsX)
	for i := 0; i < WlsX; i++ {
		assigned[i] = make([]string, WlsY)
	}

	anyassignments := false

	for _, sol := range solutions {
		ass := sol.Welladdress
		if ass != "" {
			anyassignments = true
			wc := wtype.MakeWellCoordsA1(ass)
			assigned[wc.X][wc.Y] = sol.ID
		}
	}

	return assigned, anyassignments
}

func assign_minor_layouts(group []string, plate *wtype.LHPlate, plateID string, solutions map[string]*wtype.LHSolution) (mgrps [][]string, masss []string) {
	// in this version we just use the number of wells in a column
	colsize := plate.WlsY
	rowsize := plate.WlsX

	// we need to check for existing assignments and make sure we don't clobber them

	assigned, anyassignments := make_get_assignations(solutions, rowsize, colsize)

	if anyassignments {
		return make_minor_layouts_hidebound(group, assigned, plateID, solutions, rowsize, colsize)
	} else {
		return make_minor_layouts_anew(group, plateID, rowsize, colsize)
	}
}

func make_minor_layouts_hidebound(group []string, assigned [][]string, plateID string, solutions map[string]*wtype.LHSolution, rowsize, colsize int) (mgrps [][]string, masss []string) {
	masss = make([]string, 0, rowsize)
	mgrps = make([][]string, 0, 10)

	row := 0
	col := 0

	for i := 0; i < len(group); i += 1 {
		// a layout group is now just a single entity
		grp := make([]string, 1)
		sol := solutions[group[i]]
		ass := sol.Welladdress
		grp[0] = sol.ID

		mass := ""

		if ass != "" {
			wc := wtype.MakeWellCoordsA1(ass)
			mass = plateID + ":" + wc.RowLettString() + ":" + wc.ColNumString() + ":0:0"
		} else {
			// put it in the next free slot
			tru := false
			for ; col < rowsize; row++ {
				if row == colsize {
					row = 0
					col += 1
					continue
				}

				if assigned[row][col] == "" {
					tru = true
					break
				}
				row += 1
			}

			if !tru {
				RaiseError("Inconsistent plate layout ... uArch/sched/lh/layoutagent")
			}

			wc := wtype.WellCoords{col, row}
			mass = plateID + ":" + wc.RowLettString() + ":" + wc.ColNumString() + ":0:0"
		}

		masss = append(masss, mass)
		mgrps = append(mgrps, grp)

		if row == colsize {
			row = 0
			col += 1
		}
	}
	return mgrps, masss
}

func make_minor_layouts_anew(group []string, plateID string, rowsize, colsize int) (mgrps [][]string, masss []string) {
	masss = make([]string, rowsize)
	mgrps = make([][]string, 0, 10)

	row := 1
	col := 1
	for i := 0; i < len(group); i += colsize {
		// make a layout group

		grp := make([]string, 0, colsize)

		grpsize := colsize

		if len(group)-i < grpsize {
			grpsize = len(group) - i
		}

		for j := 0; j < grpsize; j++ {
			if i+j >= len(group) {
				break
			}
			grp = append(grp, group[i+j])
		}

		// get its assignment

		ass := plateID + ":" + wutil.NumToAlpha(row) + ":" + strconv.Itoa(col) + ":" + strconv.Itoa(1) + ":" + strconv.Itoa(0)

		mgrps = append(mgrps, grp)
		masss[col-1] = ass
		col += 1
	}
	return mgrps, masss
}

func choose_major_layout_group(groups [][]string, mx int) int {
	g := 0
	for x, ar := range groups {
		if len(ar) < mx {
			g = x
			break
		}
	}
	return g
}

func do_major_layouts(request *LHRequest, majorlayoutgroups [][]string) []string {
	// we assign layout groups to plates
	plateLayouts := request.Output_plate_layout
	// we assign each mlg to a plate... since we don't have any plates yet we just give them numbers

	if len(plateLayouts) == 0 {
		plateLayouts = make([]string, len(majorlayoutgroups))

		platenum := 0
		for k, _ := range majorlayoutgroups {
			plateLayouts[k] = strconv.Itoa(platenum)
			platenum += 1
		}
	}
	return plateLayouts
}

func getLayoutGroups(solutions map[string]*wtype.LHSolution) ([]int, []int) {
	// determine which groups exist
	// we define major and minor layout groupings

	MajorLayoutGroupIDs := wutil.NewIntSet(4)
	MinorLayoutGroupIDs := wutil.NewIntSet(4)

	for _, s := range solutions {
		logger.Debug(fmt.Sprint("SOLUTION: ", s.SName, " PLAATE: ", s.Platetype))
		Mlg := s.Majorlayoutgroup

		if Mlg != -1 {
			MajorLayoutGroupIDs.Add(Mlg)
		}

		mlg := s.Minorlayoutgroup

		if mlg != -1 {
			MinorLayoutGroupIDs.Add(mlg)
		}
	}

	return MajorLayoutGroupIDs.AsSlice(), MinorLayoutGroupIDs.AsSlice()
}
