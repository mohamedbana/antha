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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"github.com/antha-lang/antha/anthalib/wutil"
	"strconv"
)

// default layout: requests fill plates in column order
func BasicLayoutAgent(request *LHRequest, params *LHProperties) *LHRequest {
	// we limit this to the case where all outputs go to the same plate type

	plate := request.Output_platetype
	solutions := request.Output_solutions

	// get the incoming group IDs

	MajorLayoutGroupIDs, _ := getLayoutGroups(solutions)

	// check we have enough and assign more if necessary

	n_plates_required := wutil.RoundInt(float64(len(solutions)) / float64(plate.Nwells))

	if len(MajorLayoutGroupIDs) < n_plates_required {
		for i := wutil.Max(MajorLayoutGroupIDs) + 1; i < n_plates_required-len(MajorLayoutGroupIDs); i++ {
			MajorLayoutGroupIDs = append(MajorLayoutGroupIDs, i)
		}
	}

	// then tidy them up

	MajorLayoutGroupRanks := wutil.MakeRankedList(MajorLayoutGroupIDs)

	// now we need to map solutions to groups

	MajorLayoutGroups := make(map[int][]string, len(MajorLayoutGroupIDs))

	// make the receptacles

	for _, i := range MajorLayoutGroupIDs {
		MajorLayoutGroups[MajorLayoutGroupRanks[i]] = make([]string, 0, 4)
	}

	// helper var

	max_major_group_size := plate.Nwells

	for _, soln := range solutions {
		id := soln.ID

		lg := soln.Majorlayoutgroup

		if lg == 0 {
			lg = choose_major_layout_group(MajorLayoutGroups, max_major_group_size)
		}

		/*
			_,ok=MajorLayoutGroups[MajorLayoutGroupRanks[lg]]

			if !ok{
				MajorLayoutGroups[MajorLayoutGroupRanks[lg]]=make([]string, 0, 4)
			}

		*/
		MajorLayoutGroups[MajorLayoutGroupRanks[lg]] = append(MajorLayoutGroups[MajorLayoutGroupRanks[lg]], id)
	}

	plateLayouts := do_major_layouts(request, MajorLayoutGroups)
	request.Output_plate_layout = plateLayouts

	// now we need to set the minor layout groups attribute in request
	// in this instance this is just mapping everything to columns

	minor_group_layouts := make([][]string, 0, len(solutions))
	assignments := make([]string, len(solutions))

	for i, grp := range MajorLayoutGroups {
		dplate := plateLayouts[i]

		plate_minor_groups, plate_assignments := assign_minor_layouts(grp, plate, dplate)

		minor_group_layouts = append(minor_group_layouts, plate_minor_groups...)
		for j, as := range plate_assignments {
			assignments[j] = as
		}
	}
	request.Output_minor_group_layouts = minor_group_layouts
	request.Output_major_group_layouts = MajorLayoutGroups
	request.Output_assignments = assignments
	return request
}

func assign_minor_layouts(group []string, plate *LHPlate, plateID string) (mgrps [][]string, masss map[int]string) {
	mgrps = make([][]string, 0, 10)
	masss = make(map[int]string, 10)

	// in this version we just use the number of wells in a column

	colsize := plate.WlsY

	row := 1
	col := 1

	for i := 0; i < len(group); i += colsize {
		// make a layout group

		grp := make([]string, 0, 8)
		grpsize := len(group) - i

		for j := 0; j < grpsize; j++ {
			if i+j >= len(group) {
				break
			}
			grp = append(grp, group[i+j])
		}

		// get its assignment

		ass := plateID + ":" + wutil.NumToAlpha(row) + ":" + strconv.Itoa(col) + ":" + strconv.Itoa(0) + ":" + strconv.Itoa(1)

		mgrps = append(mgrps, grp)
		masss[col-1] = ass
		col += 1
	}
	return mgrps, masss
}

func choose_major_layout_group(groups map[int][]string, mx int) int {
	g := 0
	for x, ar := range groups {
		if len(ar) < mx {
			g = x
			break
		}
	}
	return g
}

func do_major_layouts(request *LHRequest, majorlayoutgroups map[int][]string) map[int]string {
	// we assign layout groups to plates
	plateLayouts := request.Output_plate_layout
	// we assign each mlg to a plate... since we don't have any plates yet we just give them numbers

	if len(plateLayouts) == 0 {
		// ERROR HERE
		plateLayouts = make(map[int]string, 10)

		platenum := 0
		for k, _ := range majorlayoutgroups {
			plateLayouts[k] = strconv.Itoa(platenum)
			platenum += 1
		}
	}
	return plateLayouts
}

func getLayoutGroups(solutions map[string]*LHSolution) ([]int, []int) {
	// determine which groups exist
	// we define major and minor layout groupings

	MajorLayoutGroupIDs := make([]int, 0, 4)
	MinorLayoutGroupIDs := make([]int, 0, 4)

	for _, s := range solutions {
		Mlg := 0
		Mlg = s.Majorlayoutgroup
		MajorLayoutGroupIDs = append(MajorLayoutGroupIDs, Mlg)
		mlg := 0
		mlg = s.Minorlayoutgroup
		MinorLayoutGroupIDs = append(MinorLayoutGroupIDs, mlg)
	}

	return MajorLayoutGroupIDs, MinorLayoutGroupIDs
}

// looks up where a plate is mounted on a liquid handler as expressed in a request
func PlateLookup(rq LHRequest, id string) int {
	lookupmap := rq.Plate_lookup

	if len(lookupmap) == 0 {
		raiseError("Cannot find plate lookup")
	}

	return lookupmap[id]
}
