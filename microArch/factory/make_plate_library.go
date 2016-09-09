// anthalib/factory/make_plate_library.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/logger"
)

//lhPlateParams
//initialising an LHPlate is expensive as it means up to 1536 LHWells
//and large scale malloc and deletes play havoc with GC
//so use lhPlateArgs to only fully initialise lhPlates that are actually requested
type lhPlateParams struct {
	Platetype   string
	Mfr         string
	Nrows       int
	Ncols       int
	Size        wtype.Coordinates
	Welltype    *wtype.LHWell
	WellXOffset float64
	WellYOffset float64
	WellXStart  float64
	WellYStart  float64
	WellZStart  float64
	PostInit    []func(*wtype.LHPlate)
}

func (self *lhPlateParams) Init() *wtype.LHPlate {
	ret := wtype.NewLHPlate(self.Platetype, self.Mfr, self.Nrows, self.Ncols, self.Size, self.Welltype, self.WellXOffset, self.WellYOffset, self.WellXStart, self.WellYStart, self.WellZStart)

	for _, init := range self.PostInit {
		init(ret)
	}
	return ret
}

func makePlateLibrary() map[string]*lhPlateParams {
	plates := make(map[string]*lhPlateParams)

	offset := 0.25
	riserheightinmm := 40.0 - offset
	shallowriserheightinmm := 20.0 - offset
	coolerheight := 15.0
	pcrtuberack496 := 28.0
	incubatorheightinmm := devices.Shaker["3000 T-elm"]["Height"] * 1000

	inhecoincubatorinmm := devices.Shaker["InhecoStaticOnDeck"]["Height"] * 1000

	valueformaxheadtonotintoDSWplatewithp20tips := 4.5
	// deep square well 96
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := &lhPlateParams{"DSW96", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, valueformaxheadtonotintoDSWplatewithp20tips, nil}
	plates[plate.Platetype] = plate

	// deep square well 96 on riser
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DSW96_riser", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, riserheightinmm, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"DSW96_riser40", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, riserheightinmm, nil}
	plates[plate.Platetype] = plate

	// deep square well 96 on q instruments incubator
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DSW96_incubator", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, incubatorheightinmm, nil}
	plates[plate.Platetype] = plate

	// deep square well 96 on inheco incubator
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 1000, 100, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DSW96_inheco", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0.0, 0.0, inhecoincubatorinmm, nil}
	plates[plate.Platetype] = plate

	// 24 well deep square well plate on riser

	bottomtype := wtype.VWellBottom // 0 = flat, 2 = v shaped
	xdim := 16.8
	ydim := 16.8
	zdim := 41.3
	bottomh := 4.7

	wellcapacityinwelltypeunit := 11000.0
	welltypeunit := "ul"
	wellsperrow := 6
	wellspercolumn := 4
	residualvol := 650.0 // assume in ul

	wellxoffset := 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 18.0 //centre of well to centre of neighbouring well in y direction
	xstart := 4.5       // distance from top left side of plate to first well
	ystart := 4.5       // distance from top left side of plate to first well
	zstart := -1.0      // offset of bottom of deck to bottom of well (this includes agar estimate)

	zstart = zstart + riserheightinmm

	heightinmm := 44.1

	squarewell := wtype.NewShape("box", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), welltypeunit, wellcapacityinwelltypeunit, residualvol, squarewell, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"DSW24_riser", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"DSW24_riser40", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// shallow round well flat bottom 96
	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = &lhPlateParams{"SRWFB96", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, 9, 9, 0.0, 0.0, 1.0, nil}
	plates[plate.Platetype] = plate

	// shallow round well flat bottom 96 on riser
	rwshp = wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = &lhPlateParams{"SRWFB96_riser", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, 9, 9, 0.0, 0.0, 40.0, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"SRWFB96_riser40", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, 9, 9, 0.0, 0.0, 40.0, nil}
	plates[plate.Platetype] = plate

	// deep well strip trough 12
	stshp := wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 15000, 1000, stshp, wtype.VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DWST12", "Unknown", 1, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0, 30.0, valueformaxheadtonotintoDSWplatewithp20tips, nil}
	plates[plate.Platetype] = plate

	// deep well strip trough 12 on riser
	stshp = wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 15000, 1000, stshp, wtype.VWellBottom, 8.2, 72, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DWST12_riser", "Unknown", 1, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0, 30.0, riserheightinmm + valueformaxheadtonotintoDSWplatewithp20tips, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"DWST12_riser40", "Unknown", 1, 12, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 0, 30.0, riserheightinmm + valueformaxheadtonotintoDSWplatewithp20tips, nil}
	plates[plate.Platetype] = plate

	// deep well strip trough 8
	stshp = wtype.NewShape("box", "mm", 115.0, 8.2, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 24000, 1000, stshp, wtype.VWellBottom, 115, 8.2, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DWST8", "Unknown", 8, 1, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 49.5, 0.0, 0.0, nil}
	plates[plate.Platetype] = plate

	// deep well reservoir
	rshp := wtype.NewShape("box", "mm", 115.0, 72.0, 41.3)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 300000, 20000, rshp, wtype.VWellBottom, 115, 72, 41.3, 4.7, "mm")
	plate = &lhPlateParams{"DWR1", "Unknown", 1, 1, wtype.Coordinates{127.76, 85.48, 44.1}, welltype, 9, 9, 49.5, 0.0, 0.0, nil}
	plates[plate.Platetype] = plate

	// well area function
	// -- determined empirically since inverse cubic was giving us some numerical issues
	areaf := wutil.Quartic{A: -3.3317851312e-09, B: 0.00000225834467, C: -0.0006305492472, D: 0.1328156706978, E: 0}
	afb, _ := json.Marshal(areaf)
	afs := string(afb)
	// pcr plate with cooler
	cone := wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	welltype.SetAfVFunc(afs)
	//plate = &lhPlateParams{"pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5, nil}
	//plates[plate.Platetype] = plate
	plate = &lhPlateParams{"pcrplate_with_cooler", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, coolerheight + 0.5, nil}
	plates[plate.Platetype] = plate

	// pcr plate with 496rack
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = &lhPlateParams{"pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5, nil}
	//plates[plate.Platetype] = plate
	plate = &lhPlateParams{"pcrplate_with_496rack", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, pcrtuberack496 - 2.5, nil}
	plates[plate.Platetype] = plate

	// pcr plate skirted (on riser)
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	welltype.SetAfVFunc(afs)
	plate = &lhPlateParams{"pcrplate_skirted_riser40", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, riserheightinmm - 1.25, nil}
	plates[plate.Platetype] = plate

	plate = &lhPlateParams{"pcrplate_skirted_riser", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, riserheightinmm - 1.25, nil}
	plates[plate.Platetype] = plate

	// pcr plate skirted (on riser)
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	plate = &lhPlateParams{"pcrplate_skirted_riser20", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, shallowriserheightinmm - 1.25, nil}
	plates[plate.Platetype] = plate

	// pcr plate skirted
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	welltype.SetAfVFunc(afs)
	plate = &lhPlateParams{"pcrplate_skirted", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 0.636, nil}
	plates[plate.Platetype] = plate

	// pcr plate with incubator
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	welltype.SetAfVFunc(afs)
	plate = &lhPlateParams{"pcrplate_with_incubater", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, (15.5 + 44.0), nil}
	plates[plate.Platetype] = plate

	// Block Kombi 2ml
	eppy := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 45)

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 6.0       // offset of bottom of deck to bottom of well

	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 2000, 25, eppy, wtype.VWellBottom, 8.2, 8.2, 45, 4.7, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"Kombi2mlEpp", "Unknown", 4, 2, wtype.Coordinates{127.76, 85.48, 45}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// Eppendorfrack
	eppy = wtype.NewShape("cylinder", "mm", 8.2, 8.2, 45)

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 7.0       // offset of bottom of deck to bottom of well

	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 1500, 25, eppy, wtype.VWellBottom, 8.2, 8.2, 45, 4.7, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"eppendorfrack425_1.5ml", "Unknown", 4, 2, wtype.Coordinates{127.76, 85.48, 45}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// greiner 384 well plate flat bottom

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 2.5      // offset of bottom of deck to bottom of well

	square := wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// greiner 384 well plate flat bottom on riser

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 12.0 // modified from 14
	bottomh = 1.0

	wellxoffset = 4.5               // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5               //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                   // distance from top left side of plate to first well
	ystart = -2.5                   // distance from top left side of plate to first well
	zstart = riserheightinmm + 0.25 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"greiner384_riser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	plate = &lhPlateParams{"greiner384_riser40", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// greiner 384 well plate flat bottom on shallow riser

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 12.0 // modified from 14
	bottomh = 1.0

	wellxoffset = 4.5                     // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5                     //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                         // distance from top left side of plate to first well
	ystart = -2.5                         // distance from top left side of plate to first well
	zstart = shallowriserheightinmm + 0.5 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"greiner384_shallowriser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	plate = &lhPlateParams{"greiner384_riser20", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// NUNC 1536 well plate flat bottom on riser

	bottomtype = wtype.FlatWellBottom
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25           // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.25           //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                // distance from top left side of plate to first well
	ystart = -2.5                // distance from top left side of plate to first well
	zstart = riserheightinmm + 2 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 2, 2, 7)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 13, 2, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	plate = &lhPlateParams{"nunc1536_riser", "Unknown", 32, 48, wtype.Coordinates{127.76, 85.48, 7}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"nunc1536_riser40", "Unknown", 32, 48, wtype.Coordinates{127.76, 85.48, 7}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// 250ml box reservoir (working vol estimated to be 100ml to prevent spillage on moving decks)
	reservoirbox := wtype.NewShape("box", "mm", 71, 107, 38) // 39?
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 100000, 10000, reservoirbox, 0, 107, 71, 38, 3, "mm")
	plate = &lhPlateParams{"reservoir", "unknown", 1, 1, wtype.Coordinates{127.76, 85.48, 45}, welltype, 58, 13, 0, 0, 10, nil}
	plates[plate.Platetype] = plate

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) high res

	bottomtype = wtype.FlatWellBottom
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.25           // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.250          //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                // distance from top left side of plate to first well
	ystart = -2.5                // distance from top left side of plate to first well
	zstart = riserheightinmm + 3 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 2, 2, 7)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 13, 2, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	plate = &lhPlateParams{"Agarplateforpicking1536_riser40", "Unknown", 32, 48, wtype.Coordinates{127.76, 85.48, 7}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// Onewell SBS format Agarplate with colonies on riser (50ml agar) low res

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5              // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5              //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                  // distance from top left side of plate to first well
	ystart = -2.5                  // distance from top left side of plate to first well
	zstart = riserheightinmm + 5.5 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	// greiner one well with 50ml of agar in
	plate = &lhPlateParams{"Agarplateforpicking384_riser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"Agarplateforpicking384_riser40", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// Onewell SBS format Agarplate with colonies on shallowriser (50ml agar) low res

	bottomtype = wtype.FlatWellBottom
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5                     // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5                     //centre of well to centre of neighbouring well in y direction
	xstart = -2.5                         // distance from top left side of plate to first well
	ystart = -2.5                         // distance from top left side of plate to first well
	zstart = shallowriserheightinmm + 5.5 // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 125, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	// greiner one well with 50ml of agar in
	plate = &lhPlateParams{"Agarplateforpicking384_shallowriser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	plate = &lhPlateParams{"Agarplateforpicking384_riser20", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate
	// Onewell SBS format Agarplate with colonies on riser (30ml agar) low res

	zstart = 41 // offset of bottom of deck to bottom of well

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	// greiner one well with 50ml of agar in
	plate = &lhPlateParams{"30mlAgarplateforpicking384_riser", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate
	plate = &lhPlateParams{"30mlAgarplateforpicking384_riser40", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = &lhPlateParams{"pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5, nil}
		//plates[plate.Platetype] = plate
		plate = &lhPlateParams{"pcrplate_with_skirt", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5, nil}
		plates[plate.Platetype] = plate
	*/
	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, rwshp, wtype.FlatWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = &lhPlateParams{"pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5, nil}
	//plates[plate.Platetype] = plate
	plate = &lhPlateParams{"1L_DuranBottle", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5, nil}
	plates[plate.Platetype] = plate

	//forward position

	//	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 15, 0, ep48g, 0, 2, 4, 2, 48, "mm")
	//	plate = &lhPlateParams{"EPAGE48", "Invitrogen", 2, 26, 50, "mm", welltype, 4.5, 34, 0.0, 0.0, 2.0, nil}
	//	plates[plate.Platetype] = plate

	//refactored for reverse position

	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//can't reach all wells; change to 24 wells per row?
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 25, 0, ep48g, wtype.FlatWellBottom, 2, 4, 2, 2, "mm")
	//welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 100, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	//plate = &lhPlateParams{"EPAGE48", "Invitrogen", 2, 26, wtype.Coordinates{127.76, 85.48, 50}, welltype, 4.5, 34, -1.0, 17.25, 49.5, nil}
	plate = &lhPlateParams{"EPAGE48", "Invitrogen", 2, 26, wtype.Coordinates{127.76, 85.48, 48.5}, welltype, 4.5, 33.75, -1.0, 18.0, riserheightinmm + 4.5, nil}
	//plate = &lhPlateParams{"greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}

	plates[plate.Platetype] = plate

	// E-GEL 96 definition

	//same welltype as EPAGE

	// due to staggering of wells: 1 96well gel is set up as two well types

	// 1st type
	//can't reach all wells; change to 12 wells per row?
	plate = &lhPlateParams{"EGEL96_1", "Invitrogen", 4, 13, wtype.Coordinates{127.76, 85.48, 48.5}, welltype, 9, 18.0, 0, -1.0, riserheightinmm + 5.5, nil}
	//plate = &lhPlateParams{"greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// 2nd type
	plate = &lhPlateParams{"EGEL96_2", "Invitrogen", 4, 13, wtype.Coordinates{127.76, 85.48, 48.5}, welltype, 9, 18.0, 4.0, 7.5, riserheightinmm + 5.5, nil}
	//plate = &lhPlateParams{"greiner384", "Unknown", 16, 24, wtype.Coordinates{127.76, 85.48, 14}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}

	plates[plate.Platetype] = plate

	// falcon 6 well plate with Agar flat bottom with 4ml per well

	bottomtype = wtype.FlatWellBottom
	xdim = 37.0
	ydim = 37.0
	zdim = 20.0
	bottomh = 9.0 //(this includes agar estimate)

	wellxoffset = 39.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 39.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 9.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 3
	wellspercolumn = 2
	heightinmm = 20.0

	circle := wtype.NewShape("cylinder", "mm", 37, 37, 20)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"falcon6wellAgar", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	// Nunclon 12 well plate with Agar flat bottom 2ml per well

	bottomtype = wtype.FlatWellBottom
	xdim = 22.5 // diameter
	ydim = 22.5 // diameter
	zdim = 20.0
	bottomh = 9.0 //(this includes agar estimate)

	wellxoffset = 27.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 27.0 //centre of well to centre of neighbouring well in y direction
	xstart = 11.0      // distance from top left side of plate to first well
	ystart = 4.0       // distance from top left side of plate to first well
	zstart = 9.0       // offset of bottom of deck to bottom of well (this includes agar estimate)

	wellsperrow = 4
	wellspercolumn = 3
	heightinmm = 22.0

	circle = wtype.NewShape("cylinder", "mm", xdim, ydim, zdim)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = &lhPlateParams{"Nuncon12wellAgar", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}
	plates[plate.Platetype] = plate

	//	WellXOffset float64
	//	WellYOffset float64
	//	WellXStart  float64
	//	WellYStart  float64
	//	WellZStart  float64

	zstart = 9.0 + incubatorheightinmm // offset of bottom of deck to bottom of well (this includes agar estimate)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	plate = &lhPlateParams{"Nuncon12wellAgar_incubator", "Unknown", wellspercolumn, wellsperrow, wtype.Coordinates{127.76, 85.48, heightinmm}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}

	plate.PostInit = []func(*wtype.LHPlate){
		func(p *wtype.LHPlate) {
			consar := []string{"position_1"}
			p.SetConstrained("Pipetmax", consar)
		},
	}

	plates[plate.Platetype] = plate
	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = &lhPlateParams{"pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5, nil}
		//plates[plate.Platetype] = plate
		plate = &lhPlateParams{"pcrplate_with_skirt", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5, nil}
		plates[plate.Platetype] = plate
	*/

	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 250, 5, rwshp, wtype.FlatWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = &lhPlateParams{"pcrplate", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 6.5, nil}
	//plates[plate.Platetype] = plate
	plate = &lhPlateParams{"1L_DuranBottle", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 25.7}, welltype, 9, 9, 0.0, 0.0, 15.5, nil}
	plates[plate.Platetype] = plate

	plate = MakeGreinerVBottomPlate()
	plates[plate.Platetype] = plate

	plate = MakeGreinerVBottomPlateWithRiser()
	plates[plate.Platetype] = plate

	plate = MakeGreinerVBottomPlateWithRiser()
	plate.Platetype += "40"
	plates[plate.Platetype] = plate

	return plates
}

func MakeGreinerVBottomPlate() *lhPlateParams {
	// greiner V96

	bottomtype := wtype.VWellBottom
	xdim := 6.2
	ydim := 6.2
	zdim := 11.0
	bottomh := 1.0

	wellxoffset := 9.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset := 9.0 //centre of well to centre of neighbouring well in y direction
	xstart := 0.0      // distance from top left side of plate to first well
	ystart := 0.0      // distance from top left side of plate to first well
	zstart := 2.0      // offset of bottom of deck to bottom of well

	//	welltype = wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	rwshp := wtype.NewShape("cylinder", "mm", 6.2, 6.2, 10.0)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype := wtype.NewLHWell(nil, wtype.ZeroWellCoords(), "ul", 500, 1, rwshp, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	//	plate = &lhPlateParams{"SRWFB96", "Unknown", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, 9, 9, 0.0, 0.0, 2.0, nil}
	plate := &lhPlateParams{"GreinerSWVBottom", "Greiner", 8, 12, wtype.Coordinates{127.76, 85.48, 15}, welltype, wellxoffset, wellyoffset, xstart, ystart, zstart, nil}

	return plate
}

func MakeGreinerVBottomPlateWithRiser() *lhPlateParams {
	plate := MakeGreinerVBottomPlate()
	plate.Platetype = "GreinerSWVBottom_riser"
	plate.WellZStart = 42.0
	return plate
}

//	ep48g := wtype.NewShape("box", "mm", 2, 4, 2)
//  welltype := wtype.NewLhWell("EPAGE48", "", "", "ul", 15, 0, ep48g, 0, 2, 4, 2, bottomh, "mm")
//  plate = wtype.LHPlate("EPAGE48", "Invitrogen", 2, 26, height, "mm", welltype, 9, 22, 0.0, 0.0, 50.0)
//	plates[plate.Platetype] = plate

func GetPlateByType(typ string) *wtype.LHPlate {
	plates := makePlateLibrary()
	if p, ok := plates[typ]; ok {
		return p.Init()
	}
	logger.Fatal(fmt.Sprintf("Plate library has no plate of type \"%s\"", typ))
	return nil //keep the compiler happy
}

func GetPlateList() []string {
	plates := makePlateLibrary()

	kz := make([]string, len(plates))
	x := 0
	for name, _ := range plates {
		kz[x] = name
		x += 1
	}
	return kz
}
