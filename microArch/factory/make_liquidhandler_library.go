// anthalib/factory/make_liquidhandler_library.go: Part of the Antha language
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

package factory

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
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

func makeLiquidhandlerLibrary() map[string]*liquidhandling.LHProperties {
	robots := make(map[string]*liquidhandling.LHProperties, 2)
	robots["CyBioFelix"] = makeFelix()
	robots["CyBioGeneTheatre"] = makeGeneTheatre()
	robots["GilsonPipetmax"] = makeGilson()
	robots["Manual"] = makeManual()
	return robots
}

func makeManual() *liquidhandling.LHProperties {
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

	minvol := wunit.NewVolume(200, "ul")
	maxvol := wunit.NewVolume(1000, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("P1000Config", &minvol, &maxvol, &minspd, &maxspd, 1, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("P1000", "Gilson", hvconfig)

	minvol = wunit.NewVolume(50, "ul")
	maxvol = wunit.NewVolume(200, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	mvconfig := wtype.NewLHChannelParameter("P200Config", &minvol, &maxvol, &minspd, &maxspd, 1, false, wtype.LHVChannel, 0)
	mvadaptor := wtype.NewLHAdaptor("P200", "Gilson", mvconfig)

	minvol = wunit.NewVolume(2, "ul")
	maxvol = wunit.NewVolume(20, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lmvconfig := wtype.NewLHChannelParameter("P20Config", &minvol, &maxvol, &minspd, &maxspd, 1, false, wtype.LHVChannel, 0)
	lmvadaptor := wtype.NewLHAdaptor("P20", "Gilson", lmvconfig)

	minvol = wunit.NewVolume(1, "ul")
	maxvol = wunit.NewVolume(10, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("P10Config", &minvol, &maxvol, &minspd, &maxspd, 1, false, wtype.LHVChannel, 0)
	lvadaptor := wtype.NewLHAdaptor("P10", "Gilson", lvconfig)

	minvol = wunit.NewVolume(0.2, "ul")
	maxvol = wunit.NewVolume(2, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	vlvconfig := wtype.NewLHChannelParameter("P2Config", &minvol, &maxvol, &minspd, &maxspd, 1, false, wtype.LHVChannel, 0)
	vlvadaptor := wtype.NewLHAdaptor("P2", "Gilson", vlvconfig)

	minvol = wunit.NewVolume(0.2, "ul")
	maxvol = wunit.NewVolume(5000, "ul")
	headparams := wtype.NewLHChannelParameter("LabHand", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)
	head := wtype.NewLHHead("LabHand", "MotherNature", headparams)
	head.Adaptor = hvadaptor

	lhp.Adaptors = append(lhp.Adaptors, hvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, mvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, lmvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, lvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, vlvadaptor)
	lhp.Heads = append(lhp.Heads, head)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, head)
	return lhp
}

func makeGeneTheatre() *liquidhandling.LHProperties {
	layout := make(map[string]wtype.Coordinates)
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			posname := fmt.Sprintf("%s%d", wutil.NumToAlphaCountFromZero(j), i+1)
			var crds wtype.Coordinates
			layout[posname] = crds
		}
	}
	lhp := liquidhandling.NewLHProperties(12, "GeneTheatre", "CyBio", "discrete", "disposable", layout)

	// crucial constraint info
	lhp.Tip_preferences = []string{"A1", "A2", "A3"}
	lhp.Input_preferences = []string{"D3", "D2", "C3", "C2", "B1", "B2"}
	lhp.Output_preferences = []string{"B3", "B2", "B1", "C2", "C3", "D2"}
	lhp.Wash_preferences = []string{"C2"}
	lhp.Waste_preferences = []string{"C1"}
	lhp.Tipwaste_preferences = []string{"D1", "C1"}

	// there will be many potential configs here but in the first instance we only have
	// single-channel low-volume

	minvol := wunit.NewVolume(0.5, "ul")
	maxvol := wunit.NewVolume(25, "ul")
	minspd := wunit.NewFlowRate(0.05, "ml/min")
	maxspd := wunit.NewFlowRate(0.5, "ml/min")

	config := wtype.NewLHChannelParameter("LVconfigsingle", &minvol, &maxvol, &minspd, &maxspd, 1, false, wtype.LHVChannel, 0)
	adaptor := wtype.NewLHAdaptor("LVSingleAdaptor", "CyBio", config)
	head := wtype.NewLHHead("LVSingleHead", "CyBio", config)
	head.Adaptor = adaptor
	lhp.Adaptors = append(lhp.Adaptors, adaptor)
	lhp.Heads = append(lhp.Heads, head)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, head)

	return lhp
}

func makeFelix() *liquidhandling.LHProperties {
	layout := make(map[string]wtype.Coordinates)
	for i := 0; i < 12; i++ {
		posname := fmt.Sprintf("position_%d", i+1)
		// dont know coords for this yet
		var crds wtype.Coordinates
		layout[posname] = crds
	}

	lhp := liquidhandling.NewLHProperties(12, "Felix", "CyBio", "discrete", "disposable", layout)

	// get tips permissible from the factory
	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"position_1", "position_5", "position_3"}
	lhp.Input_preferences = []string{"position_10", "position_11", "position_12"}
	lhp.Output_preferences = []string{"position_7", "position_8", "position_9", "position_2", "position_4"}
	lhp.Wash_preferences = []string{"position_4"}
	lhp.Waste_preferences = []string{"position_6"}
	lhp.Tipwaste_preferences = []string{"position_2"} // not really used

	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(1000, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("HVconfig", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("HVAdaptor", "CyBio", hvconfig)

	newminvol := wunit.NewVolume(0.5, "ul")
	newmaxvol := wunit.NewVolume(50, "ul")
	newminspd := wunit.NewFlowRate(0.1, "ml/min")
	newmaxspd := wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("LVconfig", &newminvol, &newmaxvol, &newminspd, &newmaxspd, 8, false, wtype.LHVChannel, 0)
	lvadaptor := wtype.NewLHAdaptor("LVAdaptor", "CyBio", lvconfig)

	minvol3 := wunit.NewVolume(0.5, "ul")
	maxvol3 := wunit.NewVolume(1000, "ul")
	minspd3 := wunit.NewFlowRate(0.1, "ml/min")
	maxspd3 := wunit.NewFlowRate(0.5, "ml/min")
	headparams := wtype.NewLHChannelParameter("ChoiceHead", &minvol3, &maxvol3, &minspd3, &maxspd3, 8, false, wtype.LHVChannel, 0)
	head := wtype.NewLHHead("ChoiceHead", "CyBio", headparams)
	head.Adaptor = hvadaptor

	lhp.Adaptors = append(lhp.Adaptors, hvadaptor)
	lhp.Adaptors = append(lhp.Adaptors, lvadaptor)
	lhp.Heads = append(lhp.Heads, head)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, head)

	return lhp
}

func makeGilson() *liquidhandling.LHProperties {
	// gilson pipetmax

	layout := make(map[string]wtype.Coordinates)
	i := 0
	x0 := 3.886
	y0 := 3.513
	z0 := -82.035
	xi := 149.86
	yi := 95.25
	xp := x0
	yp := y0
	zp := z0
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			posname := fmt.Sprintf("position_%d", i+1)
			crds := wtype.Coordinates{xp, yp, zp}
			layout[posname] = crds
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := liquidhandling.NewLHProperties(9, "Pipetmax", "Gilson", "discrete", "disposable", layout)
	// get tips permissible from the factory
	SetUpTipsFor(lhp)

	lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"}
	lhp.Input_preferences = []string{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"}
	lhp.Output_preferences = []string{"position_7", "position_8", "position_9", "position_6", "position_5", "position_3"}
	lhp.Wash_preferences = []string{"position_8"}
	lhp.Tipwaste_preferences = []string{"position_1"}
	lhp.Waste_preferences = []string{"position_9"}
	//	lhp.Tip_preferences = []int{2, 3, 6, 9, 5, 8, 4, 7}
	//	lhp.Input_preferences = []int{24, 25, 26, 29, 28, 23}
	//	lhp.Output_preferences = []int{10, 11, 12, 13, 14, 15}
	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(250, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("HVconfig", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor
	newminvol := wunit.NewVolume(0.5, "ul")
	newmaxvol := wunit.NewVolume(20, "ul")
	newminspd := wunit.NewFlowRate(0.1, "ml/min")
	newmaxspd := wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("LVconfig", &newminvol, &newmaxvol, &newminspd, &newmaxspd, 8, false, wtype.LHVChannel, 1)
	lvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", lvconfig)
	lvhead := wtype.NewLHHead("LVHead", "Gilson", lvconfig)
	lvhead.Adaptor = lvadaptor

	lhp.Heads = append(lhp.Heads, hvhead)
	lhp.Heads = append(lhp.Heads, lvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, hvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, lvhead)

	return lhp
}

func GetLiquidhandlerByType(typ string) *liquidhandling.LHProperties {
	liquidhandlers := makeLiquidhandlerLibrary()
	t := liquidhandlers[typ]
	return t.Dup()
}

func LiquidhandlerList() []string {
	liquidhandlers := makeLiquidhandlerLibrary()
	kz := make([]string, len(liquidhandlers))
	x := 0
	for name, _ := range liquidhandlers {
		kz[x] = name
		x += 1
	}
	return kz
}
