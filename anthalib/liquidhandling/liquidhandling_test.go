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
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	//"strings"
	"encoding/json"
	"github.com/antha-lang/antha/anthalib/wtype"
	"math/rand"
	"testing"
)

func TestTwo(*testing.T) {
	//	for i:=0;i<20;i++{
	//		ExampleTwo()
	//	}
}

func TestOne(*testing.T) {
	ExampleOne()
}

func ExampleOne() {
	var lhr LHRequest

	sarr := make(map[string]*LHSolution, 1)

	welltype := NewLHWell("ACMEMicroPlatesDW96ConicalBottom", "", "", 2000, 25, wtype.NewShape("box"), 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := NewLHPlate("ACMEMicroPlatesDW96ConicalBottom", "ACMEMicroPlates", 8, 12, 44.1, "mm", welltype)

	lhr.Output_platetype = plate
	lhr.Input_platetype = plate

	for i := 0; i < 10; i++ {
		s := NewLHSolution()
		s.SName = fmt.Sprintf("Solution%02d", i)
		cmp := make([]*LHComponent, 0, 3)
		for j := 0; j < 3; j++ {
			c := NewLHComponent()
			c.CName = fmt.Sprintf("Component%02d", j)
			c.Type = "water"
			c.Vol = 40.0
			c.Vunit = "ul"
			cmp = append(cmp, c)
		}
		s.Components = cmp
		sarr[s.ID] = s
	}

	lhr.Output_solutions = sarr

	// make tips

	tipboxes := make([]*LHTipbox, 2, 2)
	tip := NewLHTip("ACMEliquidhandlers", "ACMEliquidhandlers250", 20.0, 250.0)
	for i := 0; i < 2; i++ {
		tb := NewLHTipbox(8, 12, "ACMEliquidhandlers", tip)
		tipboxes[i] = tb
	}

	lhr.Tips = tipboxes

	// make a liquid handling structure

	lhp := NewLHProperties(12, "ALiquidHandler", "ACMEliquidhandlers", "discrete", "disposable", []string{"plate"})

	// I suspect this might need to be in the constructor
	// or at least wrapped into a factory method

	lhp.Tip_preferences = []int{1, 5, 3}
	lhp.Input_preferences = []int{10, 11, 12}
	lhp.Output_preferences = []int{7, 8, 9, 2, 4}

	// need to add some configs

	hvconfig := NewLHParameter("HVconfig", 10, 250, "ul")

	cnfvol := lhp.Cnfvol
	cnfvol[0] = hvconfig

	lhp.CurrConf = hvconfig

	// these depend on the tip

	liquidhandler := Init(lhp)
	liquidhandler.MakeSolutions(&lhr)
}

func ExampleTwo() {
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

	for i, _ := range minrequired {
		var v float64
		v, ok := cncs[i]

		if !ok {
			v = -1.0
		}
		fmt.Printf("Concentration of %10s = %8.1f, volume High: %-8.1f volume Low: %-8.1f min required: %-8.1f Max required: %-8.1f Smax: %-8.1f T: %-6.1f\n", i, v, T[i]*maxrequired[i]/v, T[i]*minrequired[i]/v, minrequired[i], maxrequired[i], Smax[i], T[i])
	}

	fmt.Println()
}

func TestThree(*testing.T) {
	var lhr LHRequest
	sarr := make(map[string]*LHSolution, 1)

	welltype := NewLHWell("ACMEMicroPlatesDW96ConicalBottom", "", "", 2000, 25, wtype.NewShape("box"), 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := NewLHPlate("ACMEMicroPlatesDW96ConicalBottom", "ACMEMicroPlates", 8, 12, 44.1, "mm", welltype)

	m, _ := json.Marshal(plate)

	var p2 LHPlate

	u := json.Unmarshal(m, &p2)
	fmt.Println("U: ", u)
	/*
		fmt.Println()
		fmt.Println(string(m))
		fmt.Println()
	*/

	plate = &p2

	lhr.Output_platetype = plate
	lhr.Input_platetype = plate

	for i := 0; i < 10; i++ {
		s := NewLHSolution()
		s.SName = fmt.Sprintf("Solution%02d", i)
		cmp := make([]*LHComponent, 0, 3)
		for j := 0; j < 3; j++ {
			c := NewLHComponent()
			c.CName = fmt.Sprintf("Component%02d", j)
			c.Type = "water"
			c.Vol = 40.0
			c.Vunit = "ul"
			cmp = append(cmp, c)
		}
		s.Components = cmp
		sarr[s.ID] = s
	}

	lhr.Output_solutions = sarr

	// make tips

	tipboxes := make([]*LHTipbox, 2, 2)
	tip := NewLHTip("ACMEliquidhandlers", "250", 20.0, 250.0)
	for i := 0; i < 2; i++ {
		tb := NewLHTipbox(8, 12, "ACMEliquidhandlers", tip)
		tipboxes[i] = tb
	}

	lhr.Tips = tipboxes

	// make a liquid handling structure

	lhp := NewLHProperties(12, "ALiquidHandler", "ACMEliquidhandlers", "discrete", "disposable", []string{"plate"})

	// I suspect this might need to be in the constructor
	// or at least wrapped into a factory method

	lhp.Tip_preferences = []int{1, 5, 3}
	lhp.Input_preferences = []int{10, 11, 12}
	lhp.Output_preferences = []int{7, 8, 9, 2, 4}

	// need to add some configs

	hvconfig := NewLHParameter("HVconfig", 10, 250, "ul")

	cnfvol := lhp.Cnfvol
	cnfvol[0] = hvconfig
	lhp.Cnfvol = cnfvol

	lhp.CurrConf = hvconfig

	liquidhandler := Init(lhp)
	liquidhandler.MakeSolutions(&lhr)

	m, err := json.Marshal(lhr)
	fmt.Println(err)
	fmt.Println(string(m))
}
