// /anthalib/simulator/liquidhandling/simulator_test.go: Part of the Antha language
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
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	. "github.com/antha-lang/antha/microArch/factory"
)

func SetUpTipsFor(lhp *liquidhandling.LHProperties) *liquidhandling.LHProperties {
	tips := GetTipList()
	for _, tt := range tips {
		tb := GetTipByType(tt)
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			lhp.Tips = append(lhp.Tips, tb.Tips[0][0])
		}
	}
	return lhp
}


func getLHProperties() *liquidhandling.LHProperties {

	//	tips := GetTipList()

	// dummy layout of 25 positions... arbitrary limitation

	x := 0.0
	y := 0.0
	z := 0.0
	xinc := 100.0
	yinc := 100.0

	i := 0
	layout := make(map[string]wtype.Coordinates)
	for xi := 0; xi < 5; xi++ {
		for yi := 0; yi < 5; yi++ {
			posname := fmt.Sprintf("position_%d", i+1)
			crds := wtype.Coordinates{x, y, z}
			layout[posname] = crds
			i += 1
			y += yinc
		}
		x += xinc
	}
	lhp := liquidhandling.NewLHProperties(25, "Human", "MotherNature", "discrete", "disposable", layout)

	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"tips1", "tips2", "tips3", "tips4"}
	lhp.Input_preferences = []string{"in1", "in2", "in3", "in4"}
	lhp.Output_preferences = []string{"out1", "out2", "out3", "out4"}
	lhp.Tipwaste_preferences = []string{"tip_waste"}
	lhp.Wash_preferences = []string{"tip_wash"}
	lhp.Waste_preferences = []string{"liquid_waste"}

	adaptors := make([]*wtype.LHAdaptor, 0, 1)

	minvol := wunit.NewVolume(200, "ul")
	maxvol := wunit.NewVolume(1000, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("P1000Config", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("P1000", "Gilson", hvconfig)

	adaptors = append(adaptors, hvadaptor)

	minvol = wunit.NewVolume(20, "ul")
	maxvol = wunit.NewVolume(200, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	mvconfig := wtype.NewLHChannelParameter("P200Config", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	mvadaptor := wtype.NewLHAdaptor("P200", "Gilson", mvconfig)

	adaptors = append(adaptors, mvadaptor)

	minvol = wunit.NewVolume(2, "ul")
	maxvol = wunit.NewVolume(20, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lmvconfig := wtype.NewLHChannelParameter("P20Config", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	lmvadaptor := wtype.NewLHAdaptor("P20", "Gilson", lmvconfig)
	adaptors = append(adaptors, lmvadaptor)

	minvol = wunit.NewVolume(1, "ul")
	maxvol = wunit.NewVolume(10, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("P10Config", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	lvadaptor := wtype.NewLHAdaptor("P10", "Gilson", lvconfig)
	adaptors = append(adaptors, lvadaptor)

	minvol = wunit.NewVolume(0.2, "ul")
	maxvol = wunit.NewVolume(2, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	vlvconfig := wtype.NewLHChannelParameter("P2Config", minvol, maxvol, minspd, maxspd, 1, false, wtype.LHVChannel, 0)
	vlvadaptor := wtype.NewLHAdaptor("P2", "Gilson", vlvconfig)
	adaptors = append(adaptors, vlvadaptor)

	for i, adaptor := range adaptors {
		maxvol = adaptor.Params.Maxvol
		minvol = adaptor.Params.Minvol
		maxspd = adaptor.Params.Maxspd
		minspd = adaptor.Params.Minspd
		headparams := wtype.NewLHChannelParameter(fmt.Sprintf("LabHand_%d", i+1), minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
		head := wtype.NewLHHead(fmt.Sprintf("LabHand_%d", i+1), "MotherNature", headparams)
		head.Adaptor = adaptor
		lhp.Heads = append(lhp.Heads, head)
		lhp.HeadsLoaded = append(lhp.HeadsLoaded, head)
	}
	return lhp
}
