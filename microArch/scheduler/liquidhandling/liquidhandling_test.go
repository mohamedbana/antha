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
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/equipment/manual/grpc"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
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

// need to test marshalling components

func _TestMarshal(*testing.T) {
	lhc := wtype.NewLHComponent()

	lhc.Vol = 34
	lhc.Vunit = "ul"

	b, _ := json.Marshal(lhc)
	lhc2 := wtype.NewLHComponent()

	json.Unmarshal(b, &lhc2)
}

func _TestIPLinear(*testing.T) {
	// get component library
	ctypes := factory.GetComponentList()

	// make components
	cmps := make(map[string]wunit.Volume)

	for _, cmpn := range ctypes {
		vf := rand.Float64() * 4000.0
		vol := wunit.NewVolume(vf, "ul")
		cmps[cmpn] = vol
	}

	// get plate library
	plist := factory.GetPlateList()

	// no need to subselect just stick em all in

	plates := make([]*wtype.LHPlate, 0)

	for _, p := range plist {
		if p == "pcrplate_with_cooler" || p == "DSW96" {
			plates = append(plates, factory.GetPlateByType(p))
		}
	}

	// we need a map between components and volumes
	// an array of plates
	// and a map of weights and constraints

	wtc := make(map[string]float64, 3)
	wtc["MAX_N_PLATES"] = 4.0
	wtc["MAX_N_WELLS"] = 270.0
	wtc["RESIDUAL_VOLUME_WEIGHT"] = 1.0
	ass := choose_plate_assignments(cmps, plates, wtc)
	ass = ass
	cnt := 0
	for component, cmap := range ass {
		for plt, nw := range cmap {
			volreq := cmps[component]
			fmt.Printf(fmt.Sprintln("\t", nw, " wells of ", plt.Type, " total volume ", float64(nw)*(plt.Welltype.Vol-plt.Welltype.Rvol), " residual volume ", float64(nw)*plt.Welltype.Rvol, " volume required: ", volreq.RawValue()))
			cnt += nw
		}

	}
	logger.Debug(fmt.Sprintf("%d Wells total", cnt))
}

func GetNewRequest() *LHRequest {
	rq := NewLHRequest()
	isw := make(map[string]float64, 3)
	isw["MAX_N_PLATES"] = 2.5
	isw["MAX_N_WELLS"] = 96
	isw["RESIDUAL_VOLUME_WEIGHT"] = 1.0
	rq.Input_Setup_Weights = isw
	rq.Policies = liquidhandling.GetLHPolicyForTest()
	pt := factory.GetPlateByType("pcrplate_with_cooler")
	inpt := factory.GetPlateByType("pcrplate_with_cooler")
	rq.Input_platetypes = append(rq.Input_platetypes, inpt)
	rq.Output_platetype = pt

	ctypes := []string{"water", "DNAsolution", "restrictionenzyme", "tartrazine", "SapI"}

	// make components
	for i := 0; i < 10; i++ {
		samples := make([]*wtype.LHComponent, 0, len(ctypes))
		for _, cmpn := range ctypes {
			vf := 20.0
			vol := wunit.NewVolume(vf, "ul")
			samples = append(samples, mixer.Sample(factory.GetComponentByType(cmpn), vol))
		}

		soln := mixer.Mix(samples...)

		rq.Output_solutions[soln.ID] = soln
	}
	return rq
}
func GetNewDevice() *liquidhandling.LHProperties {
	lhp := factory.GetLiquidhandlerByType("CyBioGeneTheatre")
	lhp.Driver = grpc.NewDriver("localhost:50051")
	return lhp
}

func checkInputAssignments(r1, r2 *LHRequest) bool {
	for k, v := range r1.Input_assignments {
		v2 := r2.Input_assignments[k]

		if len(v) != len(v2) {
			return false
		}

		for i := 0; i < len(v); i++ {
			t1 := strings.Split(v[i], ":")
			t2 := strings.Split(v2[i], ":")

			if t1[1] != t2[1] && t1[2] != t2[2] {
				return false
			}
		}
	}
	return true
}

func TestSetupDeterminism(t *testing.T) {
	fmt.Println("Testing setup determinism")

	var lastrq *LHRequest

	for i := 0; i < 10; i++ {
		lhp := GetNewDevice()
		rq := GetNewRequest()
		lh := Init(lhp)
		lh.MakeSolutions(rq)

		if i != 0 {
			if !checkInputAssignments(rq, lastrq) {
				t.Error("Input assignments not identical")
			}
		}

		lastrq = rq
	}

}
