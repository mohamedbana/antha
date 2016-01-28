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
	"github.com/antha-lang/antha/antha/anthalib/wunit"
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

func NewShape(name, lengthunit string, h, w, d float64) *Shape {
	sh := Shape{name, lengthunit, h, w, d}
	return &sh
}
