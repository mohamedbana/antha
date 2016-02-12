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
// 2 Royal College St, London NW1 0NH UK

package wtype

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

type Coordinates struct {
	X float64
	Y float64
	Z float64
}

// interface to 3D geometry
type Geometry interface {
	Height() wunit.Length
	Width() wunit.Length
	Depth() wunit.Length
}

type Shape struct {
	ShapeName  string
	LengthUnit string
	H          float64
	W          float64
	D          float64
}

// let shape implement geometry

func (sh *Shape) Height() wunit.Length { // y?
	return wunit.NewLength(sh.H, sh.LengthUnit)
}
func (sh *Shape) Width() wunit.Length { // X?
	return wunit.NewLength(sh.W, sh.LengthUnit)
}
func (sh *Shape) Depth() wunit.Length { // Z?
	return wunit.NewLength(sh.D, sh.LengthUnit)
}

func (sh *Shape) Dup() *Shape {
	return &(Shape{sh.ShapeName, sh.LengthUnit, sh.H, sh.W, sh.D})
}

func (sh *Shape) CrossSectionalArea() (area wunit.Area, err error) {

	shapename := strings.ToLower(sh.ShapeName)
	areaunit := sh.LengthUnit + `^` + strconv.Itoa(2)

	var circular bool
	var boxlike bool

	if shapename == "circle" || shapename == "cylinder" || shapename == "round" || shapename == "sphere" {
		circular = true
	} else if shapename == "square" || shapename == "rectangle" || shapename == "box" {
		boxlike = true
	}

	if circular && sh.Height() == sh.Width() {
		area = wunit.NewArea(math.Pi*sh.H*sh.H, areaunit)
	} else if boxlike {
		area = wunit.NewArea(sh.H*sh.W, areaunit)
	} else {
		err = fmt.Errorf("No method to work out cross sectional area for shape", sh.ShapeName, "yet")
	}
	return
}

func (sh *Shape) Volume() (volume wunit.Volume, err error) {

	shapename := strings.ToLower(sh.ShapeName)
	var volumeunit string
	if sh.LengthUnit == "mm" {
		volumeunit = "ul"
	} else {
		err = fmt.Errorf("can't handle conversion of", sh.LengthUnit, "to volume unit yet")
	}

	var cylinder bool
	var boxlike bool

	if shapename == "cylinder" {
		cylinder = true
	} else if shapename == "square" || shapename == "rectangle" || shapename == "box" {
		boxlike = true
	}

	if cylinder && sh.Height() == sh.Width() {
		volume = wunit.NewVolume(math.Pi*sh.H*sh.H*sh.D, volumeunit)
	} else if boxlike {
		volume = wunit.NewVolume(sh.H*sh.W*sh.D, volumeunit)
	} else {
		err = fmt.Errorf("No method to work out volume for shape", sh.ShapeName, "yet")
	}
	return
}

func NewShape(name, lengthunit string, h, w, d float64) *Shape {
	sh := Shape{name, lengthunit, h, w, d}
	return &sh
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
func (wc *WellCoords) Format1A() string {
	return strconv.Itoa(wc.X+1) + NumToAlpha(wc.Y+1)
}
func (wc *WellCoords) FormatA1() string {
	return NumToAlpha(wc.Y+1) + strconv.Itoa(wc.X+1)
}
func (wc *WellCoords) WellNumber() int {
	return (8*(wc.X-1) + wc.Y)
}

func (wc *WellCoords) ColNumString() string {
	return strconv.Itoa(wc.X + 1)
}
func (wc *WellCoords) RowLettString() string {
	return NumToAlpha(wc.Y + 1)
}
