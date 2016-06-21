// anthalib//liquidhandling/liquidhandling_test.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math/rand"
	"testing"
)

func TestStockConcs(*testing.T) {
	names := []string{"tea", "milk", "sugar"}

	minrequired := make(map[string]float64, len(names))
	maxrequired := make(map[string]float64, len(names))
	Smax := make(map[string]float64, len(names))
	T := make(map[string]float64, len(names))
	vmin := 10.0

	for _, name := range names {
		r := rand.Float64() + 1.0
		r2 := rand.Float64() + 1.0
		r3 := rand.Float64() + 1.0

		minrequired[name] = r * r2 * 20.0
		maxrequired[name] = r * r2 * 30.0
		Smax[name] = r * r2 * r3 * 70.0
		T[name] = 100.0
	}

	cncs := choose_stock_concentrations(minrequired, maxrequired, Smax, vmin, T)
	cncs = cncs
	/*for k, v := range cncs {
		logger.Debug(fmt.Sprintln(k, " ", minrequired[k], " ", maxrequired[k], " ", T[k], " ", v))
	}*/
}

func TestPlateReuse(t *testing.T) {
	lh := GetLiquidHandlerForTest()
	rq := GetLHRequestForTest()

	water := GetComponentForTest("water", wunit.NewVolume(100.0, "ul"))
	mmx := GetComponentForTest("mastermix_sapI", wunit.NewVolume(100.0, "ul"))
	part := GetComponentForTest("dna", wunit.NewVolume(50.0, "ul"))

	for k := 0; k < 9; k++ {
		ins := wtype.NewLHInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(8.0, "ul"))
		mmxs := mixer.Sample(mmx, wunit.NewVolume(8.0, "ul"))
		ps := mixer.Sample(part, wunit.NewVolume(1.0, "ul"))

		ins.AddComponent(ws)
		ins.AddComponent(mmxs)
		ins.AddComponent(ps)
		ins.AddProduct(GetComponentForTest("water", wunit.NewVolume(17.0, "ul")))
		rq.Add_instruction(ins)
	}

	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	fmt.Println("PLAN 1")
	rq.ConfigureYourself()
	err := lh.Plan(rq)

	// this should test whether reuse is functioning

	fmt.Println("PLAN 2")
	instrx := rq.LHInstructions
	// reset the request
	rq = GetLHRequestForTest()

	for _, ins := range instrx {
		rq.Add_instruction(ins)
	}

	for _, plateid := range lh.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		thing := lh.Properties.PlateLookup[plateid]

		plate, ok := thing.(*wtype.LHPlate)
		if !ok {
			continue
		}

		// this should kill it

		for _, v := range plate.HWells {
			if !v.Empty() {
				v.Remove(wunit.NewVolume(5.0, "ul"))
			}
		}

		rq.Input_plates[plateid] = plate
	}
	rq.ConfigureYourself()

	lh = GetLiquidHandlerForTest()
	err = lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got error resimulating: ", err))
	}

}
