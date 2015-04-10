// wtype/geometry.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wtype

import (
	"github.com/Synthace/vectormath"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/anthalib/wutil"
	"strconv"
	"strings"
)

// alias coordinate structure to spate's vectormath
type coordinates vectormath.Vector3

// interface to 3D geometry
type Geometry interface {
	Height() wunit.Length
	Width() wunit.Length
	Depth() wunit.Length
}

// defines a shape
type Shape interface {
	ShapeName() string
	IsShape()
	MinEnclosingBox() Geometry
}

// convenience structure for handling well coordinates
type WellCoords struct {
	X int
	Y int
}

// make well coordinates in the "A1" convention
func MakeWellCoordsA1(a1 string) WellCoords {
	// only handles 96 well plates
	return WellCoords{wutil.ParseInt(a1[1:len(a1)]) - 1, AlphaToNum(string(a1[0])) - 1}
}

// make well coordinates in the "1A" convention
func MakeWellCoords1A(a1 string) WellCoords {
	// only handles 96 well plates
	return WellCoords{AlphaToNum(string(a1[0])) - 1, wutil.ParseInt(a1[1:len(a1)]) - 1}
}

// make well coordinates in a manner compatble with "X1,Y1" etc.
func MakeWellCoordsXYsep(x, y string) WellCoords {
	return WellCoords{wutil.ParseInt(y[1:len(y)]) - 1, wutil.ParseInt(x[1:len(x)]) - 1}
}

func MakeWellCoordsXY(xy string) WellCoords {
	tx := strings.Split(xy, "Y")
	x := wutil.ParseInt(tx[0][1:len(tx[0])]) - 1
	y := wutil.ParseInt(tx[1]) - 1
	return WellCoords{x, y}
}

// return well coordinates in "X1Y1" format
func (wc *WellCoords) FormatXY() string {
	return "X" + strconv.Itoa(wc.X+1) + "Y" + strconv.Itoa(wc.Y+1)
}
func (wc *WellCoords) FormatA1() string {
	return strconv.Itoa(wc.X+1) + NumToAlpha(wc.Y+1)
}
func (wc *WellCoords) Format1A() string {
	return NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}
func (wc *WellCoords) WellNumber() int {
	return (8*(wc.X-1) + wc.Y)
}
