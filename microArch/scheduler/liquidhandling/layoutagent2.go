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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	//"github.com/antha-lang/antha/microArch/logger"
)

func ImprovedLayoutAgent(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest {
	// we have three kinds of solution
	// 1- ones going to a specific plate
	// 2- ones going to a specific plate type
	// 3- ones going to a plate of our choosing

	// find existing assignments

	plate_choices := get_assignments(request)

	// now we know what remains unassigned, we assign it

	choose_plates(request, plate_choices)

	// now we have plates of type 1 & 2

	// make specific plates... this may mean splitting stuff out into multiple plates

	make_plates(request)

	// now we have plates of type 1 only -- we just need to
	// say where on each plate they will go
	// this needs to set Output_assignments
	make_layouts(request)

	return request
}

type PlateChoice struct {
	Platetype string
	Assigned  []string
	ID        string
}

func get_assignments(request *LHRequest) []PlateChoice {
	s := make([]PlateChoice, 0, 3)
	m := make(map[int]string)

	// inconsistent plate types will be assigned randomly!
	for k, v := range request.Output_solutions {
		if v.PlateID != "" {
			i := defined(v.PlateID, s)

			if i == -1 {
				s = append(s, PlateChoice{v.Platetype, []string{v.ID}, v.PlateID})
			} else {
				s[i].Assigned = append(s[i].Assigned, v.ID)
			}
		} else if v.Majorlayoutgroup != -1 {
			id, ok := m[v.Majorlayoutgroup]
			if !ok {
				id := wtype.NewUUID()
				m[v.Majorlayoutgroup] = id
				request.Output_solutions[k].PlateID = id
			}

			i := defined(id, s)

			if i == -1 {
				s = append(s, PlateChoice{v.Platetype, []string{v.ID}, v.PlateID})
			} else {
				s[i].Assigned = append(s[i].Assigned, v.ID)
			}
		}
	}

	// make sure the plate choices all have defined types

	for i, _ := range s {
		if s[i].Platetype == "" {
			s[i].Platetype = request.Output_platetypes[0].Type
		}
	}

	return s
}

func defined(s string, pc []PlateChoice) int {
	r := -1

	for i, v := range pc {
		if v.ID == s {
			r = i
			break
		}
	}
	return r
}

func choose_plates(request *LHRequest, pc []PlateChoice) {
	for k, v := range request.Output_solutions {
		if v.PlateID == "" {
			pt := v.PlateType

			ass := assignmentWithType(pt, pc)

			if ass == -1 {
				// make a new plate
				ass = len(pc)
				pc = append(pc, PlateChoice{chooseAPlate(request, v), []string{v.ID}, wutil.GetUUID()})
			}

			pc[ass].Assigned = append(pc[ass].Assigned, v.ID)
		}
	}

	// make sure the plate isn't too full

	pc2 := make([]PlateChoice, 0, len(pc))

	for i, v := range pc {
		plate := factory.GetPlateByType(v.Platetype)

		// chop the assignments up

		pc2 = append(pc2, modpc(v, plate.NWells)...)
	}

	// copy the choices in

	for _, c := range pc2 {
		for _, i := range pc.Assigned {
			request.Output_solutions[i].PlateID = c.ID
			request.Output_solutions[i].Platetype = c.Platetype
		}
	}
}

func modpc(choice PlateChoice, nwell int) []PlateChoice {
	r := make([]Platechoice, 0, 1)

	for s := 0; s < len(choice.Assigned); s += nwell {
		e := s + nwell
		if e > len(choice.Assigned) {
			e = len(choice.Assigned)
		}
		r = append(r, PlateChoice{choice.PlateType, choice.Assigned[s:e], wutil.NewUUID()})
	}
	return r
}

func assignmentWithType(pt string, pc []PlateChoice) int {
	r := -1

	if pt == "" {
		if len(pc) != 0 {
			r = 0
		}
		return r
	}

	for i, v := range pc {
		if pt == pc.PlateType {
			r = i
			break
		}
	}

	return r
}

func chooseAPlate(request *LHRequest, sol *wtype.LHSolution) string {
	// for now we ignore sol and just choose the First Output Platetype
	return request.Output_platetype[0]
}
func stringinarray(s string, array []string) int {
	r := -1

	for i, k := range array {
		if k == s {
			r = i
			break
		}
	}

	return r
}

func plateidarray(arr []*wtype.LHPLate) []string {
	ret := make([]string, 0, 3)

	for _, v := range arr {
		ret = append(ret, v.ID)
	}
	return ret
}

func make_plates(request *LHRequest) {
	pids := plateidarray(request.Output_plates)
	remap := make(map[string]string)
	for k, v := range request.Output_solutions {
		i := stringinarray(v.PlateID, pids)

		if i == -1 {
			plate := factory.GetPlateByType(v.PlateType)
			request.Output_plates = append(request.Output_plates, plate)
			pids = append(pids, v.PlateID)
			remap[v.PlateID] = plate.ID
		}

		rm, ok := remap[v.PlateID]

		if ok {
			request.Output_Solutions[k].PlateID = remap[v.PlateID]
		}
	}
}

func make_assignments(request *LHRequest) {
	// finally what are the relevant assignments?
	// TODO
	// -- major concern: we need to get Output_assignments filled in here
	//    need to revise the way this is used
	//    basically this layout might still not be that useful -- different
	//    planning might be preferred

	// this is quite a big job potentially. Also possible conflict with the execution plan
}
