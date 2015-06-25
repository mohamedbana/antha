// /anthalib/wtype/concreteshape.go: Part of the Antha language
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

package wtype

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// concrete shapes for use

// concrete Geometry object

type G3D struct {
	H wunit.Length
	W wunit.Length
	D wunit.Length
}

func (g *G3D) Height() wunit.Length {
	return g.H
}
func (g *G3D) Width() wunit.Length {
	return g.W
}
func (g *G3D) Depth() wunit.Length {
	return g.D
}

func (g *G3D) IsShape() {
}
func (g *G3D) ShapeName() string {
	return "box"
}
func (g *G3D) MinEnclosingBox() Geometry {
	return g
}

func EZG3D(h, w, d float64) *G3D {
	// danger danger danger - defaults to si units which are in METRES!
	// your boxes may be a touch on the large side if you're not careful!

	return NewG3D(wunit.EZLength(h), wunit.EZLength(w), wunit.EZLength(d))
}
func NewG3D(h, w, d wunit.Length) *G3D {
	if h.RawValue() < 0.0 || w.RawValue() < 0.0 || d.RawValue() < 0.0 {
		panic(fmt.Sprintf("Nonsensical dimensions requested for G3D object: %f %f %f\n", h, w, d))
	}

	// TODO -- enforce required relationships here: h >=w
	// 	-- we don't screw around with d since there may be a defined top face

	g3d := G3D{h, w, d}
	return &g3d
}

// easy access to a few basic shapetypes
// these really need to be better wrapped
func NewShape(shapetype string) Shape {
	var s Shape
	switch shapetype {
	case "box":
		b := Box{EZG3D(0.0, 0.0, 0.0)}
		s = &b
	case "cylinder":
		c := Cylinder{}
		s = &c
	default:
		b := Box{EZG3D(0.0, 0.0, 0.0)}
		s = &b
	}
	return s
}

type Box struct {
	*G3D
}

func (b *Box) ShapeName() string {
	return "box"
}

func (b *Box) IsShape() {
}

func (b *Box) MinEnclosingBox() Geometry {
	return b.G3D
}

type Cylinder struct {
	R wunit.Length
	H wunit.Length
}

func (c *Cylinder) ShapeName() string {
	return "cylinder"
}
func (c *Cylinder) IsShape() {
}
func (c *Cylinder) MinEnclosingBox() Geometry {
	return EZG3D(c.H.SIValue(), c.R.SIValue(), c.R.SIValue())
}
