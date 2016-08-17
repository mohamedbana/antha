// wtype/shape.go: Part of the Antha language
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
	"math"
)

//Shape a more detailed representation of an object's shape
//Shapes should also implement LHObject
type Shape interface {
	Volume() wunit.Volume
	MaxCrossSectionalArea() wunit.Area
}

type Direction int

const (
	Upwards Direction = iota
	Downwards
	Leftwards
	Rightwards
	Forwards
	Backwards
)

func lengthInMM(length float64, units string) float64 {
	l := wunit.NewLength(length, units)
	return l.ConvertToString("mm")
}

//###########################################################
//						BoxShape
//###########################################################

type BoxShape struct {
	bounds BBox
	parent LHObject
}

//Constructor
func NewBoxShape(size_x, size_y, size_z float64, units string) *BoxShape {
	size := Coordinates{
		lengthInMM(size_x, units),
		lengthInMM(size_y, units),
		lengthInMM(size_z, units),
	}
	r := BoxShape{BBox{Coordinates{}, size}, nil}
	return &r
}

//@implement LHObject
func (self *BoxShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *BoxShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *BoxShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *BoxShape) GetPointIntersections(point Coordinates) []LHObject {
	//conveniently this shape is a bounding box :)
	if self.bounds.IntersectsPoint(point) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *BoxShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *BoxShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *BoxShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *BoxShape) Volume() wunit.Volume {
	return wunit.NewVolume(self.bounds.Size().X*self.bounds.Size().Y*self.bounds.Size().Z, "mm^3")
}

//@implement Shape
func (self *BoxShape) MaxCrossSectionalArea() wunit.Area {
	return wunit.NewArea(self.bounds.Size().X*self.bounds.Size().Y, "mm^3")
}

//###########################################################
//						CylinderShape
//###########################################################

type CylinderShape struct {
	bounds BBox
	parent LHObject
}

//Constructor
func NewCylinderShape(size_x, size_y, size_z float64, units string) *CylinderShape {
	size := Coordinates{
		lengthInMM(size_x, units),
		lengthInMM(size_y, units),
		lengthInMM(size_z, units),
	}
	r := CylinderShape{BBox{Coordinates{}, size}, nil}
	return &r
}

//@implement LHObject
func (self *CylinderShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *CylinderShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *CylinderShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *CylinderShape) GetPointIntersections(point Coordinates) []LHObject {
	//coarse intersection check
	if !self.bounds.IntersectsPoint(point) {
		return nil
	}

	rx := self.GetSize().X / 2.
	ry := self.GetSize().Y / 2
	//point relative to centeral X-axis
	point = point.Subtract(self.GetPosition().Add(Coordinates{rx, ry, 0}))

	//general case for any ellipse
	if (point.X/rx)*(point.X/rx)+(point.Y/ry)*(point.Y/ry) < 1 {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *CylinderShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *CylinderShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *CylinderShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *CylinderShape) Volume() wunit.Volume {
	return wunit.NewVolume(0.25*math.Pi*self.bounds.Size().X*self.bounds.Size().Y*self.bounds.Size().Z, "mm^3")
}

//@implement Shape
func (self *CylinderShape) MaxCrossSectionalArea() wunit.Area {
	return wunit.NewArea(0.25*math.Pi*self.bounds.Size().X*self.bounds.Size().Y, "mm^2")
}

//###########################################################
//						SphereShape
//###########################################################

type SphereShape struct {
	bounds BBox
	parent LHObject
}

//Constructor
func NewSphereShape(size_x, size_y, size_z float64, units string) *SphereShape {
	size := Coordinates{
		lengthInMM(size_x, units),
		lengthInMM(size_y, units),
		lengthInMM(size_z, units),
	}
	r := SphereShape{BBox{Coordinates{}, size}, nil}
	return &r
}

//@implement LHObject
func (self *SphereShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *SphereShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *SphereShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *SphereShape) GetPointIntersections(point Coordinates) []LHObject {
	//coarse intersection check
	if !self.bounds.IntersectsPoint(point) {
		return nil
	}

	rx := self.bounds.GetSize().X / 2.
	ry := self.bounds.GetSize().Y / 2.
	rz := self.bounds.GetSize().Z / 2.

	//point relative to sphere center
	point = point.Subtract(Coordinates{rx, ry, rz}).Subtract(self.GetPosition())

	//general to any ellipsoid
	if (point.X/rx)*(point.X/rx)+(point.Y/ry)*(point.Y/ry)+(point.Z/rz)*(point.Z/rz) < 1 {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *SphereShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *SphereShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *SphereShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *SphereShape) Volume() wunit.Volume {
	return wunit.NewVolume((4.*math.Pi/3.)*(self.GetSize().X/2.)*(self.GetSize().Y/2.)*(self.GetSize().Z/2.), "mm^3")
}

//@implement Shape
func (self *SphereShape) MaxCrossSectionalArea() wunit.Area {
	return wunit.NewArea(0.25*math.Pi*self.bounds.Size().X*self.bounds.Size().Y, "mm^2")
}

//###########################################################
//						HSphereShape
//###########################################################

type HSphereShape struct {
	face   Direction
	bounds BBox
	parent LHObject
}

//Constructor a hemi-sphere, face sets the direction of the flat face
func NewHSphereShape(size_x, size_y, size_z float64, units string, face Direction) *HSphereShape {
	size := Coordinates{
		lengthInMM(size_x, units),
		lengthInMM(size_y, units),
		lengthInMM(size_z, units),
	}
	r := HSphereShape{face, BBox{Coordinates{}, size}, nil}
	return &r
}

//@implement LHObject
func (self *HSphereShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *HSphereShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *HSphereShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *HSphereShape) GetPointIntersections(point Coordinates) []LHObject {
	//coarse intersection check
	if !self.bounds.IntersectsPoint(point) {
		return nil
	}

	//find the center
	var center Coordinates
	s := self.bounds.GetSize()
	switch self.face {
	case Upwards:
		center = Coordinates{s.X / 2., s.Y / 2., 0}
	case Downwards:
		center = Coordinates{s.X / 2., s.Y / 2., s.Z}
	case Leftwards:
		center = Coordinates{0, s.Y / 2., s.Z / 2.}
	case Rightwards:
		center = Coordinates{s.X, s.Y / 2., s.Z / 2.}
	case Forwards:
		center = Coordinates{s.X / 2., 0, s.Z / 2.}
	case Backwards:
		center = Coordinates{s.X / 2., s.Y, s.Z / 2.}
	}
	center = center.Add(self.GetPosition())

	//get point relative to center
	point = point.Subtract(center)

	if (point.X*2./s.X)*(point.X*2./s.X)+(point.Y*2./s.Y)*(point.Y*2./s.Y)+(point.Z*2./s.Z)*(point.Z*2./s.Z) < 1 {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *HSphereShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *HSphereShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *HSphereShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *HSphereShape) Volume() wunit.Volume {
	return wunit.NewVolume((2.*math.Pi/3.)*(self.GetSize().X/2.)*(self.GetSize().Y/2.)*(self.GetSize().Z/2.), "mm^3")
}

//@implement Shape
func (self *HSphereShape) MaxCrossSectionalArea() wunit.Area {
	//it's possible this is meant to be max cross section in XY plane, in which case this is wrong
	//As it just calculates the area of the flat face
	//but it seems unlikely to be a problem for the time being
	var a float64
	switch self.face {
	case Upwards:
		a = 0.25 * math.Pi * self.bounds.Size().X * self.bounds.Size().Y
	case Downwards:
		a = 0.25 * math.Pi * self.bounds.Size().X * self.bounds.Size().Y
	case Leftwards:
		a = 0.25 * math.Pi * self.bounds.Size().Z * self.bounds.Size().Y
	case Rightwards:
		a = 0.25 * math.Pi * self.bounds.Size().Z * self.bounds.Size().Y
	case Forwards:
		a = 0.25 * math.Pi * self.bounds.Size().X * self.bounds.Size().Z
	case Backwards:
		a = 0.25 * math.Pi * self.bounds.Size().X * self.bounds.Size().Z
	}
	return wunit.NewArea(a, "mm^2")
}

//###########################################################
//						ConeShape
//###########################################################

type ConeShape struct {
	face   Direction
	bounds BBox
	parent LHObject
}

//Constructor
func NewConeShape(size_x, size_y, size_z float64, units string, face Direction) *ConeShape {
	size := Coordinates{
		lengthInMM(size_x, units),
		lengthInMM(size_y, units),
		lengthInMM(size_z, units),
	}
	r := ConeShape{face, BBox{Coordinates{}, size}, nil}
	return &r
}

//@implement LHObject
func (self *ConeShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *ConeShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *ConeShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *ConeShape) GetPointIntersections(point Coordinates) []LHObject {
	//coarse intersection check
	if !self.bounds.IntersectsPoint(point) {
		return nil
	}

	//point relative to center
	s := self.GetSize()
	point = point.Subtract(self.GetPosition()).Subtract(s.Divide(2.))

	//convert to the same axes
	var a, b, h, A, B, H float64
	switch self.face {
	case Upwards:
		a = point.X
		b = point.Y
		h = 0.5*s.Z - point.Z
		A = s.X / 2
		B = s.Y / 2
		H = s.Z
	case Downwards:
		a = point.X
		b = point.Y
		h = point.Z - 0.5*s.Z
		A = s.X / 2
		B = s.Y / 2
		H = s.Z
	case Leftwards:
		a = point.Z
		b = point.Y
		h = 0.5*s.X - point.X
		A = s.Z / 2
		B = s.Y / 2
		H = s.X
	case Rightwards:
		a = point.Z
		b = point.Y
		h = point.X - 0.5*s.X
		A = s.Z / 2
		B = s.Y / 2
		H = s.X
	case Forwards:
		a = point.X
		b = point.Z
		h = 0.5*s.Y - point.Y
		A = s.X / 2
		B = s.Z / 2
		H = s.Y
	case Backwards:
		a = point.X
		b = point.Z
		h = point.Y - 0.5*s.Y
		A = s.X / 2
		B = s.Z / 2
		H = s.Y
	}

	//find radii at height
	A = A * (1 - h/H)
	B = B * (1 - h/H)

	//is the point in the circle?
	if (a/A)*(a/A)+(b/B)*(b/B) < 1 {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *ConeShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *ConeShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *ConeShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *ConeShape) Volume() wunit.Volume {
	s := self.GetSize()
	return wunit.NewVolume(math.Pi*s.X*s.Y*s.Z/(3*4), "mm^3")
}

//@implement Shape
func (self *ConeShape) MaxCrossSectionalArea() wunit.Area {
	s := self.GetSize()
	var a, b float64
	switch self.face {
	case Upwards:
		a = s.X / 2
		b = s.Y / 2
	case Downwards:
		a = s.X / 2
		b = s.Y / 2
	case Leftwards:
		a = s.Z / 2
		b = s.Y / 2
	case Rightwards:
		a = s.Z / 2
		b = s.Y / 2
	case Forwards:
		a = s.X / 2
		b = s.Z / 2
	case Backwards:
		a = s.X / 2
		b = s.Z / 2
	}

	return wunit.NewArea(math.Pi*a*b, "mm^2")
}

//###########################################################
//						SqPyrShape
//###########################################################

type SqPyrShape struct {
	face   Direction
	bounds BBox
	parent LHObject
}

//Constructor
func NewSqPyrShape(size_x, size_y, size_z float64, units string, face Direction) *SqPyrShape {
	size := Coordinates{
		lengthInMM(size_x, units),
		lengthInMM(size_y, units),
		lengthInMM(size_z, units),
	}
	r := SqPyrShape{face, BBox{Coordinates{}, size}, nil}
	return &r
}

//@implement LHObject
func (self *SqPyrShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *SqPyrShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *SqPyrShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *SqPyrShape) GetPointIntersections(point Coordinates) []LHObject {
	//coarse intersection check
	if !self.bounds.IntersectsPoint(point) {
		return nil
	}

	//point relative to center
	s := self.GetSize()
	point = point.Subtract(self.GetPosition()).Subtract(s.Divide(2.))

	//convert to the same axes
	var a, b, h, A, B, H float64
	switch self.face {
	case Upwards:
		a = point.X
		b = point.Y
		h = 0.5*s.Z - point.Z
		A = s.X / 2
		B = s.Y / 2
		H = s.Z
	case Downwards:
		a = point.X
		b = point.Y
		h = point.Z - 0.5*s.Z
		A = s.X / 2
		B = s.Y / 2
		H = s.Z
	case Leftwards:
		a = point.Z
		b = point.Y
		h = 0.5*s.X - point.X
		A = s.Z / 2
		B = s.Y / 2
		H = s.X
	case Rightwards:
		a = point.Z
		b = point.Y
		h = point.X - 0.5*s.X
		A = s.Z / 2
		B = s.Y / 2
		H = s.X
	case Forwards:
		a = point.X
		b = point.Z
		h = 0.5*s.Y - point.Y
		A = s.X / 2
		B = s.Z / 2
		H = s.Y
	case Backwards:
		a = point.X
		b = point.Z
		h = point.Y - 0.5*s.Y
		A = s.X / 2
		B = s.Z / 2
		H = s.Y
	}

	//find widths at height
	A = A * (1 - h/H)
	B = B * (1 - h/H)

	//is the point in the square
	if math.Abs(a) < A && math.Abs(b) < B {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *SqPyrShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *SqPyrShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *SqPyrShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *SqPyrShape) Volume() wunit.Volume {
	s := self.GetSize()
	return wunit.NewVolume(s.X*s.Y*s.Z/3, "mm^3")
}

//@implement Shape
func (self *SqPyrShape) MaxCrossSectionalArea() wunit.Area {
	s := self.GetSize()
	var a, b float64
	switch self.face {
	case Upwards:
		a = s.X
		b = s.Y
	case Downwards:
		a = s.X
		b = s.Y
	case Leftwards:
		a = s.Z
		b = s.Y
	case Rightwards:
		a = s.Z
		b = s.Y
	case Forwards:
		a = s.X
		b = s.Z
	case Backwards:
		a = s.X
		b = s.Z
	}

	return wunit.NewArea(a*b, "mm^2")
}

//###########################################################
//						CompositeShape
//###########################################################

type CompositeShape struct {
	children []Shape
	bounds   BBox
	parent   LHObject
}

//Constructor
func NewCompositeShape(children []Shape) *CompositeShape {
	bounds := BBox{}
	for _, child := range children {
		ch := child.(*LHObject)
		bounds = bounds.Merge(BBox{ch.GetPosition(), ch.GetSize()})
	}
	r := CompositeShape{children, bounds, nil}
	return &r
}

//@implement LHObject
func (self *CompositeShape) GetPosition() Coordinates {
	if self.parent != nil {
		return self.bounds.GetPosition().Add(self.parent.GetPosition())
	}
	return self.bounds.GetPosition()
}

//@implement LHObject
func (self *CompositeShape) GetSize() Coordinates {
	return self.bounds.GetSize()
}

//@implement LHObject
func (self *CompositeShape) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	if self.parent {
		box.SetRelativeTo(self.parent.GetPosition())
	}
	if box.IntersectsBox(self.bounds) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *CompositeShape) GetPointIntersections(point Coordinates) []LHObject {
	//coarse intersection check
	if !self.bounds.IntersectsPoint(point) {
		return nil
	}

	for _, ch := range self.children {
		if len(ch.GetPointIntersections(point)) > 0 {
			return []LHObject{self}
		}
	}
	return nil
}

//@implement LHObject
func (self *CompositeShape) SetOffset(c Coordinates) error {
	self.bounds.SetPosition(c)
	return nil
}

//@implement LHObject
func (self *CompositeShape) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *CompositeShape) GetParent() LHObject {
	return self.parent
}

//@implement Shape
func (self *CompositeShape) Volume() wunit.Volume {
	v = wunit.NewVolume(0, "mm^3")
	for _, ch := range self.children {
		r.Add(ch.Volume())
	}
	return r
}

//@implement Shape
func (self *CompositeShape) MaxCrossSectionalArea() wunit.Area {
	a = wunit.NewArea(0, "mm^2")
	for _, ch := range self.children {
		if na := ch.MaxCrossSectionalArea(); na.GreaterThan(a) {
			a = na
		}
	}
	return a
}
