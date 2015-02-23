package wtype

import (
	"fmt"
	"github.com/antha-lang/antha/anthalib/wunit"
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
