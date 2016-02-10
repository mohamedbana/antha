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
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

const (
	COLWISE = iota
	ROWWISE
	RANDOM
)

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

func get_assignment(assignments []string, plates *map[string]*wtype.LHPlate, vol wunit.Volume) (string, wunit.Volume, bool) {
	assignment := ""
	ok := false
	prevol := wunit.NewVolume(0.0, "ul")

	for _, assignment = range assignments {
		asstx := strings.Split(assignment, ":")
		plate := (*plates)[asstx[0]]

		crds := asstx[1] + ":" + asstx[2]
		wellidlkp := plate.Wellcoords
		well := wellidlkp[crds]

		currvol := well.CurrVolume()
		currvol.Subtract(well.ResidualVolume())
		if currvol.GreaterThan(vol) || currvol.EqualTo(vol) {
			prevol = well.CurrVolume()
			well.Remove(vol)
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

type IChain struct {
	Parent *IChain
	Child  *IChain
	Values []*wtype.LHInstruction
}

func NewIChain(parent *IChain) *IChain {
	var it IChain
	it.Parent = parent
	it.Values = make([]*wtype.LHInstruction, 0, 1)
	return &it
}

func (it *IChain) Add(ins *wtype.LHInstruction) {
	p := it.FindNodeFor(ins)
	p.Values = append(p.Values, ins)
}

func (it *IChain) GetChild() *IChain {
	if it.Child == nil {
		it.Child = NewIChain(it)
	}
	return it.Child
}

func (it *IChain) FindNodeFor(ins *wtype.LHInstruction) *IChain {
	pstr := ins.ParentString()

	if pstr == "" {
		if it.Parent == nil {
			return it
		} else {
			// should not be here!
			logger.Fatal("Improper use of IChain")
		}
	} else {
		for _, v := range it.Values {
			// true if any component used by ins is *this*
			if ins.HasParent(v.ProductID) {
				return it.GetChild()
			}
		}

		return it.Child.FindNodeFor(ins)
	}
	// unreachable: pstr either is or isn't ""
	return nil
}

func (it *IChain) Flatten() []string {
	var ret []string
	for _, v := range it.Values {
		ret = append(ret, v.ID)
	}

	ret = append(ret, it.Child.Flatten()...)

	return ret
}

func set_output_order(rq *LHRequest) {
	// gather things into groups with dependency relationships
	// TODO -- implement time constraints and anything else

	it := NewIChain(nil)

	for _, v := range rq.Order_instructions_added {
		it.Add(rq.LHInstructions[v])
	}

	rq.Output_order = it.Flatten()
}

func ConvertInstruction(insIn *wtype.LHInstruction, robot *driver.LHProperties) (insOut *driver.TransferInstruction) {
	wh := make([]string, len(insIn.Components))        // component types
	va := make([]*wunit.Volume, len(insIn.Components)) // volumes

	// four parameters applying to the destination

	pt := make([]string, len(insIn.Components)) // dest plate positions
	wt := make([]string, len(insIn.Components)) // dest wells
	ptwx := make([]int, len(insIn.Components))  // dimensions of plate pipetting to (X)
	ptwy := make([]int, len(insIn.Components))  // dimensions of plate pipetting to (Y)

	// four parameters applying to the source

	fromPlateID, fromWells := robot.GetComponents(insIn.Components)

	pf := make([]string, len(insIn.Components))
	wf := make([]string, len(insIn.Components))
	pfwx := make([]int, len(insIn.Components))
	pfwy := make([]int, len(insIn.Components))

	for i, v := range insIn.Components {
		wh[i] = v.TypeName()
		v2 := wunit.NewVolume(v.Vol, v.Vunit)
		va[i] = &v2
		pt[i] = robot.PlateIDLookup[insIn.PlateID]
		wt[i] = insIn.Welladdress
		ptwx[i] = robot.Plates[insIn.PlateID].WellsX()
		ptwy[i] = robot.Plates[insIn.PlateID].WellsY()
		pf[i] = robot.PlateIDLookup[fromPlateID[i]]
		wf[i] = fromWells[i]
		pfwx[i] = robot.Plates[fromPlateID[i]].WellsX()
		pfwy[i] = robot.Plates[fromPlateID[i]].WellsY()
	}

	ti := driver.TransferInstruction{Type: driver.TFR, What: wh, Volume: va, PltTo: pt, WellTo: wt, TPlateWX: ptwx, TPlateWY: ptwy, PltFrom: pf, WellFrom: wf, FPlateWX: pfwx, FPlateWY: pfwy}
	return &ti
}
