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
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//var commonwelltypes

func makePlateLibrary() map[string]*wtype.LHPlate {
	plates := make(map[string]*wtype.LHPlate)

	riserheightinmm := 40.0

	// deep square well 96
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("DSW96", "", "", "ul", 1000, 100, swshp, 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate("DSW96", "Unknown", 8, 12, 44.1, "mm", welltype, 9, 9, 0.0, 0.0, 0.0)
	plates[plate.Type] = plate

	// deep square well 96 on riser
	swshp = wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype = wtype.NewLHWell("DSW96", "", "", "ul", 1000, 100, swshp, 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DSW96_riser", "Unknown", 8, 12, 44.1, "mm", welltype, 9, 9, 0.0, 0.0, 40.0)
	plates[plate.Type] = plate

	// 24 well deep square well plate on riser

	bottomtype := 3 // 0 = flat, 3 = v shaped?
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
	welltype = wtype.NewLHWell("24DSW", "", "", welltypeunit, wellcapacityinwelltypeunit, residualvol, squarewell, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("DSW24_riser", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// shallow round well flat bottom 96
	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell("SRWFB96", "", "", "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96", "Unknown", 8, 12, 15, "mm", welltype, 9, 9, 0.0, 0.0, 2.0)
	plates[plate.Type] = plate

	// shallow round well flat bottom 96 on riser
	rwshp = wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell("SRWFB96", "", "", "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96_riser", "Unknown", 8, 12, 15, "mm", welltype, 9, 9, 0.0, 0.0, 41.0)
	plates[plate.Type] = plate

	// deep well strip trough 12
	stshp := wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell("DWST12", "", "", "ul", 15000, 1000, stshp, 3, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12", "Unknown", 1, 12, 44.1, "mm", welltype, 9, 9, 0, 0, 0.0)
	plates[plate.Type] = plate

	// deep well strip trough 12 on riser
	stshp = wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell("DWST12", "", "", "ul", 15000, 1000, stshp, 3, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12_riser", "Unknown", 1, 12, 44.1, "mm", welltype, 9, 9, 0, 0, 41.0)
	plates[plate.Type] = plate

	// deep well strip trough 8
	stshp = wtype.NewShape("box", "mm", 115.0, 8.2, 41.3)
	welltype = wtype.NewLHWell("DWST8", "", "", "ul", 24000, 1000, stshp, 3, 115, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST8", "Unknown", 8, 1, 44.1, "mm", welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	// deep well reservoir
	rshp := wtype.NewShape("box", "mm", 115.0, 72.0, 41.3)
	welltype = wtype.NewLHWell("DWR1", "", "", "ul", 300000, 20000, rshp, 3, 115, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWR1", "Unknown", 1, 1, 44.1, "mm", welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	// pcr plate with cooler
	cone := wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, cone, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("pcrplate_with_cooler", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	// pcr plate with incubator
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, cone, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	plate = wtype.NewLHPlate("pcrplate_with_incubater", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, (15.5 + 44.0))
	plates[plate.Type] = plate
	// pcr plate skirted
	cone = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 200, 5, cone, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	plate = wtype.NewLHPlate("pcrplate_skirted", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 37.5)
	plates[plate.Type] = plate

	// Block Kombi 2ml
	eppy := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 45)

	wellxoffset = 18.0 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 18.0 //centre of well to centre of neighbouring well in y direction
	xstart = 5.0       // distance from top left side of plate to first well
	ystart = 5.0       // distance from top left side of plate to first well
	zstart = 6.0       // offset of bottom of deck to bottom of well

	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("2mlEpp", "", "", "ul", 2000, 25, eppy, 3, 8.2, 8.2, 45, 4.7, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("Kombi2mlEpp", "Unknown", 4, 2, 45, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// greiner 384 well plate flat bottom

	bottomtype = 0
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 2        // offset of bottom of deck to bottom of well

	square := wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("384flat", "", "", "ul", 100, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, 14, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// greiner 384 well plate flat bottom on riser

	bottomtype = 0
	xdim = 4.0
	ydim = 4.0
	zdim = 14.0
	bottomh = 1.0

	wellxoffset = 4.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 4.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 43       // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 4, 4, 14)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("384flat", "", "", "ul", 100, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("greiner384_riser", "Unknown", 16, 24, 14, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// NUNC 1536 well plate flat bottom on riser

	bottomtype = 0
	xdim = 2.0 // of well
	ydim = 2.0
	zdim = 7.0
	bottomh = 0.5

	wellxoffset = 2.5 // centre of well to centre of neighbouring well in x direction
	wellyoffset = 2.5 //centre of well to centre of neighbouring well in y direction
	xstart = -2.5     // distance from top left side of plate to first well
	ystart = -2.5     // distance from top left side of plate to first well
	zstart = 42       // offset of bottom of deck to bottom of well

	square = wtype.NewShape("box", "mm", 2, 2, 7)
	//func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	welltype = wtype.NewLHWell("1536flat", "", "", "ul", 13, 2, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("nunc1536_riser", "Unknown", 32, 48, 7, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// 250ml box reservoir (working vol estimated to be 100ml to prevent spillage on moving decks)
	reservoirbox := wtype.NewShape("box", "mm", 71, 107, 38) // 39?
	welltype = wtype.NewLHWell("Reservoir", "", "", "ul", 100000, 10000, reservoirbox, 0, 107, 71, 38, 3, "mm")
	plate = wtype.NewLHPlate("reservoir", "unknown", 1, 1, 45, "mm", welltype, 58, 13, 0, 0, 10)
	plates[plate.Type] = plate
	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
		//plates[plate.Type] = plate
		plate = wtype.NewLHPlate("pcrplate_with_skirt", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
		plates[plate.Type] = plate
	*/
	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	//forward position

	//	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//	welltype = wtype.NewLHWell("EPAGE48", "", "", "ul", 15, 0, ep48g, 0, 2, 4, 2, 48, "mm")
	//	plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, 50, "mm", welltype, 4.5, 34, 0.0, 0.0, 2.0)
	//	plates[plate.Type] = plate

	//refactored for reverse position

	ep48g := wtype.NewShape("trap", "mm", 2, 4, 2)
	//can't reach all wells; change to 24 wells per row?
	welltype = wtype.NewLHWell("EPAGE48", "", "", "ul", 25, 0, ep48g, 0, 2, 4, 2, 2, "mm")
	//welltype = wtype.NewLHWell("384flat", "", "", "ul", 100, 10, square, bottomtype, xdim, ydim, zdim, bottomh, "mm")
	//plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, 50, "mm", welltype, 4.5, 34, -1.0, 17.25, 49.5)
	plate = wtype.NewLHPlate("EPAGE48", "Invitrogen", 2, 26, 48.5, "mm", welltype, 4.5, 33.75, -1.0, 17.25, 47.5)
	//plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, 14, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates[plate.Type] = plate

	// E-GEL 96 definition

	//same welltype as EPAGE

	// due to staggering of wells: 1 96well gel is set up as two well types

	// 1st type
	//can't reach all wells; change to 12 wells per row?
	plate = wtype.NewLHPlate("EGEL96_1", "Invitrogen", 4, 13, 48.5, "mm", welltype, 9, 18.0, 0, -1.0, 47.5)
	//plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, 14, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	// 2nd type
	plate = wtype.NewLHPlate("EGEL96_2", "Invitrogen", 4, 13, 48.5, "mm", welltype, 9, 18.0, 4.0, 7.5, 47.5)
	//plate = wtype.NewLHPlate("greiner384", "Unknown", 16, 24, 14, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)

	plates[plate.Type] = plate

	// falcon 6 well plate with Agar flat bottom

	bottomtype = 0
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
	welltype = wtype.NewLHWell("falcon6well", "", "", "ul", 100, 10, circle, bottomtype, xdim, ydim, zdim, bottomh, "mm")

	//func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	plate = wtype.NewLHPlate("falcon6wellAgar", "Unknown", wellspercolumn, wellsperrow, heightinmm, "mm", welltype, wellxoffset, wellyoffset, xstart, ystart, zstart)
	plates[plate.Type] = plate

	//	WellXOffset float64
	//	WellYOffset float64
	//	WellXStart  float64
	//	WellYStart  float64
	//	WellZStart  float64

	/*
		rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
		welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
		//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
		//plates[plate.Type] = plate
		plate = wtype.NewLHPlate("pcrplate_with_skirt", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
		plates[plate.Type] = plate
	*/
	/// placeholder for non plate container for testing
	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("1L_DuranBottle", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate

	//NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64)

	return plates
}

//	ep48g := wtype.NewShape("box", "mm", 2, 4, 2)
//  welltype := wtype.NewLhWell("EPAGE48", "", "", "ul", 15, 0, ep48g, 0, 2, 4, 2, bottomh, "mm")
//  plate = wtype.LHPlate("EPAGE48", "Invitrogen", 2, 26, height, "mm", welltype, 9, 22, 0.0, 0.0, 50.0)
//	plates[plate.Type] = plate

func GetPlateByType(typ string) *wtype.LHPlate {
	plates := makePlateLibrary()
	p := plates[typ]
	return p.Dup()
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
