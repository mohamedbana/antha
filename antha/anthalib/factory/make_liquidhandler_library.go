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
// 1 Royal College St, London NW1 0NH UK

package factory

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func makeLiquidhandlerLibrary() map[string]*liquidhandling.LHProperties {
	robots := make(map[string]*liquidhandling.LHProperties, 2)
	robots["CyBioFelix"] = makeCyBio()
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

	lhp.Tip_preferences = []int{2, 3, 4, 5, 6}
	lhp.Input_preferences = []int{7, 8, 9, 10, 11, 12, 13, 14, 15}
	lhp.Output_preferences = []int{16, 17, 18, 19, 20, 21, 22, 23, 24, 25}

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

func makeCyBio() *liquidhandling.LHProperties {
	tips := GetTipList()
	layout := make(map[string]wtype.Coordinates)
	for i := 0; i < 12; i++ {
		posname := fmt.Sprintf("position_%d", i+1)
		// dont know coords for this yet
		var crds wtype.Coordinates
		layout[posname] = crds
	}

	lhp := liquidhandling.NewLHProperties(12, "Felix", "CyBio", "discrete", "disposable", layout)

	for _, tt := range tips {
		tb := GetTipByType(tt)
		if tb.Mnfr == lhp.Mnfr {
			lhp.Tips = append(lhp.Tips, tb.Tips[0][0])
		}
	}

	lhp.Tip_preferences = []int{1, 5, 3}
	lhp.Input_preferences = []int{10, 11, 12}
	lhp.Output_preferences = []int{7, 8, 9, 2, 4}

	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(1000, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("HVconfig", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)

	hvadaptor := wtype.NewLHAdaptor("HVAdaptor", "CyBio", hvconfig)

	minvol = wunit.NewVolume(0.5, "ul")
	maxvol = wunit.NewVolume(50, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("LVconfig", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)
	lvadaptor := wtype.NewLHAdaptor("LVAdaptor", "CyBio", lvconfig)

	minvol = wunit.NewVolume(0.5, "ul")
	maxvol = wunit.NewVolume(1000, "ul")
	headparams := wtype.NewLHChannelParameter("ChoiceHead", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)
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
	tips := GetTipList()

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

	for _, tt := range tips {
		tb := GetTipByType(tt)
		if tb.Mnfr == lhp.Mnfr {
			lhp.Tips = append(lhp.Tips, tb.Tips[0][0])
		}
	}

	lhp.Tip_preferences = []int{2, 3, 4}
	lhp.Input_preferences = []int{4, 5, 6}
	lhp.Output_preferences = []int{7, 8, 9}
	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(250, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("HVconfig", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor

	minvol = wunit.NewVolume(0.5, "ul")
	maxvol = wunit.NewVolume(50, "ul")
	minspd = wunit.NewFlowRate(0.1, "ml/min")
	maxspd = wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("LVconfig", &minvol, &maxvol, &minspd, &maxspd, 8, false, wtype.LHVChannel, 1)
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
