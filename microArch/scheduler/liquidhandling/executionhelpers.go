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
	"fmt"
	"sort"
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

func insSliceFromMap(m map[string]*wtype.LHInstruction) []*wtype.LHInstruction {
	ret := make([]*wtype.LHInstruction, 0, len(m))

	for _, v := range m {
		ret = append(ret, v)
	}

	return ret
}

type ByGeneration []*wtype.LHInstruction

func (bg ByGeneration) Len() int      { return len(bg) }
func (bg ByGeneration) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByGeneration) Less(i, j int) bool {
	if bg[i].Generation() == bg[j].Generation() {
		strings.Compare(bg[i].Welladdress, bg[j].Welladdress)
	}

	return bg[i].Generation() == bg[j].Generation()
}

func set_output_order(rq *LHRequest) {
	// sort into equivalence classes by generation
	// nb that this basically means the below is a bit pointless
	// however for now it will be maintained to test whether
	// parent-child relationships are working OK

	sorted := insSliceFromMap(rq.LHInstructions)

	sort.Sort(ByGeneration(sorted))

	it := NewIChain(nil)

	for _, v := range sorted {
		it.Add(v)
	}

	rq.Output_order = it.Flatten()

	// wha
	it.Print()

	rq.InstructionChain = it
}

func optimize_runs(rq *LHRequest, chain *IChain, newoutputorder []string) {
	// go through instructions on the same level and see if any might be candidates for
	// aggregation

	// this will replace both the instructions and the order, since the instructions now have new IDs

	// might as well make this recursive

	if chain == nil {
		rq.Output_order = newoutputorder
		return
	}

	arrIns := groupByComponents(chain.Values)

	for _, ins := range arrIns {
		newoutputorder = append(newoutputorder, ins.ID)
		rq.LHInstructions[ins.ID] = ins
	}

	optimize_runs(rq, chain.Child, newoutputorder)
}

func groupByComponents(instructions []*wtype.LHInstruction) []*wtype.LHInstruction {
	/*
		hsh := make(map[string][]*LHInstruction, len(instructions))

		for _, ins := range instructions {
			hsh[ins.Result.CName] = append(hsh[ins.Result.CName], ins)
		}

		// component ordering needs deciding here... as a general rule it's
		// best to stick with higher volumes first
	*/
	r := make([]*wtype.LHInstruction, 0, 1)
	return r
}

func ConvertInstruction(insIn *wtype.LHInstruction, robot *driver.LHProperties) (insOut *driver.TransferInstruction) {

	cmps := insIn.Components

	lenToMake := len(insIn.Components)

	if insIn.IsMixInPlace() {
		lenToMake = lenToMake - 1
		cmps = cmps[1:len(cmps)]
	}

	wh := make([]string, lenToMake)       // component types
	va := make([]wunit.Volume, lenToMake) // volumes

	// six parameters applying to the source

	fromPlateID, fromWells := robot.GetComponents(cmps)

	pf := make([]string, lenToMake)
	wf := make([]string, lenToMake)
	pfwx := make([]int, lenToMake)
	pfwy := make([]int, lenToMake)
	vf := make([]wunit.Volume, lenToMake)
	ptt := make([]string, lenToMake)

	// six parameters applying to the destination

	pt := make([]string, lenToMake)       // dest plate positions
	wt := make([]string, lenToMake)       // dest wells
	ptwx := make([]int, lenToMake)        // dimensions of plate pipetting to (X)
	ptwy := make([]int, lenToMake)        // dimensions of plate pipetting to (Y)
	vt := make([]wunit.Volume, lenToMake) // volume in well to
	ptf := make([]string, lenToMake)      // plate types

	ix := 0

	for i, v := range insIn.Components {
		if insIn.IsMixInPlace() && i == 0 {
			continue
		}

		// get dem big ole plates out
		// TODO -- pass them in instead of all this nonsense

		flhp := robot.PlateLookup[fromPlateID[ix]].(*wtype.LHPlate)
		tlhp := robot.PlateLookup[insIn.PlateID].(*wtype.LHPlate)

		wlt, ok := tlhp.WellAtString(insIn.Welladdress)

		if !ok {
			logger.Fatal(fmt.Sprint("Well ", insIn.Welladdress, " not found on dest plate ", insIn.PlateID))
		}

		v2 := wunit.NewVolume(v.Vol, v.Vunit)
		vt[ix] = wlt.CurrVolume()
		wh[ix] = v.TypeName()
		va[ix] = v2
		pt[ix] = robot.PlateIDLookup[insIn.PlateID]
		wt[ix] = insIn.Welladdress
		ptwx[ix] = tlhp.WellsX()
		ptwy[ix] = tlhp.WellsY()
		ptt[ix] = tlhp.Type

		wlf, ok := flhp.WellAtString(fromWells[ix])

		if !ok {
			logger.Fatal(fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateID[ix]))
		}

		vf[ix] = wlf.CurrVolume()
		//wlf.Remove(va[ix])

		pf[ix] = robot.PlateIDLookup[fromPlateID[ix]]
		wf[ix] = fromWells[ix]
		pfwx[ix] = flhp.WellsX()
		pfwy[ix] = flhp.WellsY()
		ptf[ix] = flhp.Type

		//fmt.Println("HERE GOES: ", i, wh[i], vf[i].ToString(), vt[i].ToString(), va[i].ToString(), pt[i], wt[i], pf[i], wf[i], pfwx[i], pfwy[i], ptwx[i], ptwy[i])

		ix += 1
	}

	ti := driver.TransferInstruction{Type: driver.TFR, What: wh, Volume: va, PltTo: pt, WellTo: wt, TPlateWX: ptwx, TPlateWY: ptwy, PltFrom: pf, WellFrom: wf, FPlateWX: pfwx, FPlateWY: pfwy, FVolume: vf, TVolume: vt, FPlateType: ptf, TPlateType: ptt}
	return &ti
}
