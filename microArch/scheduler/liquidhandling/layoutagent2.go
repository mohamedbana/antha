// liquidhandling/layoutagent2.go: Part of the Antha language
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
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
)

func ImprovedLayoutAgent(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest {
	// do this multiply based on the order in the chain

	logger.Debug("IMPROVED LAY OUT AGENT YEAH")

	ch := request.InstructionChain
	pc := make([]PlateChoice, 0, 3)
	mp := make(map[int]string)
	for {
		if ch == nil {
			break
		}
		request, pc, mp = LayoutStage(request, params, ch, pc, mp)
		ch = ch.Child
	}

	return request
}
func LayoutStage(request *LHRequest, params *liquidhandling.LHProperties, chain *IChain, plate_choices []PlateChoice, mapchoices map[int]string) (*LHRequest, []PlateChoice, map[int]string) {
	// we have three kinds of solution
	// 1- ones going to a specific plate
	// 2- ones going to a specific plate type
	// 3- ones going to a plate of our choosing

	logger.Debug("LAY OUT STAGE YEAH")

	for _, lhi := range chain.Values {
		logger.Debug("Instruction: ", lhi.ID, " 1st component: ", lhi.Components[0].ID, " Result: ", lhi.ProductID)
	}

	// find existing assignments

	plate_choices, mapchoices = get_and_complete_assignments(request, chain.ValueIDs(), plate_choices, mapchoices)

	logger.Debug(fmt.Sprint("PLATE CHOICE LENGTH 1: ", len(plate_choices)))

	// now we know what remains unassigned, we assign it

	plate_choices = choose_plates(request, plate_choices, chain.ValueIDs())

	logger.Debug(fmt.Sprint("PLATE CHOICE LENGTH 2: ", len(plate_choices)))

	// now we have plates of type 1 & 2

	// make specific plates... this may mean splitting stuff out into multiple plates

	remap := make_plates(request, chain.ValueIDs())

	// give them names

	for _, v := range request.Output_plates {
		// we need to ensure this has a name
		if v.Name() == "" {
			v.PlateName = fmt.Sprintf("Output_plate_%s", v.ID[0:6])
		}
	}

	// now we have solutions of type 1 only -- we just need to
	// say where on each plate they will go
	// this needs to set Output_assignments
	make_layouts(request, plate_choices)

	lkp := make(map[string][]*wtype.LHComponent)
	lk2 := make(map[string]string)
	// fix the output locations correctly

	for _, v := range request.LHInstructions {
		lkp[v.ID] = make([]*wtype.LHComponent, 0, 1) //v.Result
		lk2[v.Result.ID] = v.ID
	}

	for _, v := range request.LHInstructions {
		for _, c := range v.Components {
			// if this component has the same ID
			// as the result of another instruction
			// we map it in
			iID, ok := lk2[c.ID]

			if ok {
				// iID is an instruction ID
				lkp[iID] = append(lkp[iID], c)
			}
		}

		// now we put the actual result in
		lkp[v.ID] = append(lkp[v.ID], v.Result)
	}

	// now map the output assignments in
	for k, v := range request.Output_assignments {
		for _, id := range v {
			l := lkp[id]
			for _, x := range l {
				// x.Loc = k
				// also need to remap the plate id
				tx := strings.Split(k, ":")
				x.Loc = remap[tx[0]] + ":" + tx[1]
				logger.Debug(fmt.Sprint("REMAPPING HERE: ", id, " ", x.ID, " ", x.Loc))
				logger.Track(fmt.Sprintf("OUTPUT ASSIGNMENT I=%s R=%s A=%s", id, x.ID, x.Loc))
			}
		}
	}

	return request, plate_choices, mapchoices
}

type PlateChoice struct {
	Platetype string
	Assigned  []string
	ID        string
	Wells     []string
}

func get_and_complete_assignments(request *LHRequest, order []string, s []PlateChoice, m map[int]string) ([]PlateChoice, map[int]string) {
	fmt.Println("GET AND COMPLETE ASSIGNMENTS")
	//s := make([]PlateChoice, 0, 3)
	//m := make(map[int]string)

	// inconsistent plate types will be assigned randomly!
	//	for k, v := range request.LHInstructions {
	//for _, k := range request.Output_order {

	for _, k := range order {
		v := request.LHInstructions[k]
		if v.PlateID != "" {
			i := defined(v.PlateID, s)

			if i == -1 {
				s = append(s, PlateChoice{v.Platetype, []string{v.ID}, v.PlateID, []string{v.Welladdress}})
			} else {
				s[i].Assigned = append(s[i].Assigned, v.ID)
				s[i].Wells = append(s[i].Wells, v.Welladdress)
			}
		} else if v.Majorlayoutgroup != -1 {
			id, ok := m[v.Majorlayoutgroup]
			if !ok {
				id = wtype.NewUUID()
				m[v.Majorlayoutgroup] = id
			}

			//  fix the plate id to this temporary one
			request.LHInstructions[k].PlateID = id

			i := defined(id, s)

			if i == -1 {
				s = append(s, PlateChoice{v.Platetype, []string{v.ID}, id, []string{v.Welladdress}})
			} else {
				s[i].Assigned = append(s[i].Assigned, v.ID)
				s[i].Wells = append(s[i].Wells, v.Welladdress)
			}
		} else if v.IsMixInPlace() {
			// the first component sets the destination
			// and now it should indeed be set

			addr := v.Components[0].Loc
			logger.Debug(fmt.Sprint("ID: ", v.ID, " ", v.Components[0].ID, " THIS SHOULD NOT BE NIL: ", addr))
			tx := strings.Split(addr, ":")
			request.LHInstructions[k].Plateaddress = tx[0]
			request.LHInstructions[k].Welladdress = tx[1]
		}
	}

	// make sure the plate choices all have defined types

	for i, _ := range s {
		if s[i].Platetype == "" {
			s[i].Platetype = request.Output_platetypes[0].Type
		}
	}

	return s, m
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

func choose_plates(request *LHRequest, pc []PlateChoice, order []string) []PlateChoice {
	for _, k := range order {
		v := request.LHInstructions[k]
		// this id may be temporary, only things without it still are not assigned to a
		// plate, even a virtual one
		if v.PlateID == "" {
			pt := v.Platetype

			ass := assignmentWithType(pt, pc)

			if ass == -1 {
				// make a new plate
				ass = len(pc)
				pc = append(pc, PlateChoice{chooseAPlate(request, v), []string{v.ID}, wtype.GetUUID(), []string{""}})
			}

			pc[ass].Assigned = append(pc[ass].Assigned, v.ID)
		}
	}

	// make sure the plate isn't too full

	pc2 := make([]PlateChoice, 0, len(pc))

	for _, v := range pc {
		plate := factory.GetPlateByType(v.Platetype)

		// chop the assignments up

		pc2 = append(pc2, modpc(v, plate.Nwells)...)
	}

	// copy the choices in

	for _, c := range pc2 {
		for _, i := range c.Assigned {
			request.LHInstructions[i].PlateID = c.ID
			request.LHInstructions[i].Platetype = c.Platetype
		}
	}
	return pc2
}

// chop the assignments up modulo plate size
func modpc(choice PlateChoice, nwell int) []PlateChoice {
	r := make([]PlateChoice, 0, 1)

	for s := 0; s < len(choice.Assigned); s += nwell {
		e := s + nwell
		if e > len(choice.Assigned) {
			e = len(choice.Assigned)
		}
		logger.Debug("S:", s, " E:", e, " L: ", len(choice.Assigned), " LW: ", len(choice.Wells))
		r = append(r, PlateChoice{choice.Platetype, choice.Assigned[s:e], wtype.GetUUID(), choice.Wells[s:e]})
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
		if pt == v.Platetype {
			r = i
			break
		}
	}

	return r
}

func chooseAPlate(request *LHRequest, ins *wtype.LHInstruction) string {
	// for now we ignore ins and just choose the First Output Platetype
	return request.Output_platetypes[0].Type
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

func plateidarray(arr []*wtype.LHPlate) []string {
	ret := make([]string, 0, 3)

	for _, v := range arr {
		ret = append(ret, v.ID)
	}
	return ret
}

// we have potentially added extra theoretical plates above
// now we make real plates and swap them in

func make_plates(request *LHRequest, order []string) map[string]string {
	remap := make(map[string]string)
	//for k, v := range request.LHInstructions {
	for _, k := range order {
		v := request.LHInstructions[k]
		_, skip := remap[v.PlateID]

		if skip {
			request.LHInstructions[k].PlateID = remap[v.PlateID]
			continue
		}
		_, ok := request.Output_plates[v.PlateID]

		if !ok {
			plate := factory.GetPlateByType(v.Platetype)
			request.Output_plates[plate.ID] = plate
			remap[v.PlateID] = plate.ID
			request.LHInstructions[k].PlateID = remap[v.PlateID]
		}

	}

	return remap
}

func make_layouts(request *LHRequest, pc []PlateChoice) {
	// we need to fill in the platechoice structure then
	// transfer the info across to the solutions

	opa := request.Output_assignments

	for _, c := range pc {
		// make a temporary plate to hold info

		plat := factory.GetPlateByType(c.Platetype)

		// make an iterator for it

		it := request.OutputIteratorFactory(plat)

		//seed in the existing assignments

		for _, w := range c.Wells {
			if w != "" {
				wc := wtype.MakeWellCoords(w)
				//plat.Cols[wc.X][wc.Y].Currvol += 100.0
				dummycmp := wtype.NewLHComponent()
				dummycmp.SetVolume(wunit.NewVolume(100.0, "ul"))
				plat.Cols[wc.X][wc.Y].Add(dummycmp)
			}
		}

		for i, _ := range c.Assigned {
			sID := c.Assigned[i]
			well := c.Wells[i]

			var assignment string

			if well == "" {
				wc := plat.NextEmptyWell(it)

				if wc.IsZero() {
					// something very bad has happened
					logger.Fatal("DIRE WARNING: The unthinkable has happened... output plate has too many assignments!")
				}

				//plat.Cols[wc.X][wc.Y].Currvol += 100.0
				dummycmp := wtype.NewLHComponent()
				dummycmp.SetVolume(wunit.NewVolume(100.0, "ul"))
				plat.Cols[wc.X][wc.Y].Add(dummycmp)
				request.LHInstructions[sID].Welladdress = wc.FormatA1()
				assignment = c.ID + ":" + wc.FormatA1()
			} else {
				assignment = c.ID + ":" + well
			}

			opa[assignment] = append(opa[assignment], sID)
		}
	}

	request.Output_assignments = opa
}
