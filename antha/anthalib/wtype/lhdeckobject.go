// liquidhandling/lhdeckobject.go: Part of the Antha language
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
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// defines types for dealing with liquid handling requests
package wtype

import (
    "math"
)

//WellReference used for specifying position within a well
type WellReference int

const (
    BottomReference WellReference = iota //0
    TopReference                         //1
    LiquidReference                      //2
)

type Dimension int
const (
    XDim Dimension = iota
    YDim
    ZDim
)

//LHObject Provides a unified interface to physical size
//of items that can be placed on a liquid handler's deck
type LHObject interface {
    //GetSize Return the physical size of the object's bounding box in mm
    GetSize() Coordinates
    //GetPosition Return the absolute offset of the object
    GetPosition() Coordinates
    //SetOffset Store the offset of the object
    SetPosition(Coordinates)
    //GetTop maximum x-position
    GetMax(Dimension) float64
    //GetTop minimum x-position
    GetMin(Dimension) float64
}

//Addressable unifies the interface to objects which have
//sub-components that can be adressed by WellCoords (e.g. "A1")
//for example tip-boxes, plates, etc
type Addressable interface {
    //GetSize return the number of rows and columns
    GetSize() WellCoords
    //HasCoords Do the given coordinates exist in the object?
    HasCoords(WellCoords) bool
    //GetCoords Return the object at the given well coords
    //bool is false if the coordinate doesn't exist
    //Currently this is either an LHWell or LHTip
    GetCoords(WellCoords) (interface{}, bool)
    //CoordsToWellCoords Convert Real world coordinates 
    //(relative to the object origin) to WellCoords. 
    //The returned WellCoords should be the closest 
    //addressable location to the coorinates, and shold only be 
    //invalid if the object has no adressable locations (e.g. wells
    //or tips). The second return value gives the offset from the top
    //of the center of the well/tip to the given coordinate
    //(this leaves the caller to ascertain whether any mis-alignment
    //is acceptable)
    CoordsToWellCoords(Coordinates) (WellCoords, Coordinates)
    //WellCoordsToCoords Get the physical location of an addressable
    //position relative to the object origin.
    //WellCoords should be valid in the object, or the bool will 
    //return false and Coordinates are undefined.
    //WellReference is the position within a well.
    //Requesting LiquidReference on a LHTipbox will return false
    WellCoordsToCoords(WellCoords, WellReference) (Coordinates, bool)
}

//Intersects just checks for bounding box intersection
func Intersects(lhs LHObject, rhs LHObject) bool {
    if lhs == nil || rhs == nil {
        return false
    }
    //test a single dimension. 
    //(a,b) are the start and end of the first position
    //(c,d) are the start and end of the second pos
    // assert(a > b  and  d > c)
    f := func(a,b,c,d float64) bool {
        return !(c >= b || d <= a)
    }

    lo := lhs.GetPosition()
    ls := lhs.GetSize()
    ro := rhs.GetPosition()
    rs := rhs.GetSize()

    return (f(lo.X, lo.X + ls.X, ro.X,  ro.X +  rs.X) &&
            f(lo.Y, lo.Y + ls.Y, ro.Y,  ro.Y +  rs.Y) &&
            f(lo.Z, lo.Z + ls.Z, ro.Z,  ro.Z +  rs.Z))
}

//BBox is a simple LHObject representing a bounding box, 
//useful for checking if there's stuff in the way 
type BBox struct {
    position    Coordinates
    size        Coordinates
}

func NewBBox(position, size Coordinates) *BBox {
    r := BBox{position, size}
    return &r
}

func NewBBox6f(pos_x, pos_y, pos_z, size_x, size_y, size_z float64) *BBox {
    return NewBBox(Coordinates{ pos_x,  pos_y,  pos_z}, 
                   Coordinates{size_x, size_y, size_z})
}

func (self *BBox) GetSize() Coordinates {
    return self.size
}

func (self *BBox) GetPosition() Coordinates {
    return self.position
}

func (self *BBox) SetPosition(c Coordinates) {
    self.position = c
}

func (self *BBox) GetMax(d Dimension) float64 {
    switch d {
    case XDim:
        return math.Max(self.position.X, self.position.X + self.size.X)
    case YDim:
        return math.Max(self.position.Y, self.position.Y + self.size.Y)
    }
    //case ZDim:
    return math.Max(self.position.Z, self.position.Z + self.size.Z)
}

func (self *BBox) GetMin(d Dimension) float64 {
    switch d {
    case XDim:
        return math.Min(self.position.X, self.position.X + self.size.X)
    case YDim:
        return math.Min(self.position.Y, self.position.Y + self.size.Y)
    }
    //case ZDim:
    return math.Min(self.position.Z, self.position.Z + self.size.Z)
}

//XBox is a BBox which extends infinitely in the X direction
type XBox struct {
    BBox
}

func NewXBox(position, size Coordinates) *XBox {
    size.X = math.MaxFloat64
    position.X = -0.5 * math.MaxFloat64
    r := XBox{BBox{position, size}}
    return &r
}

func NewXBox4f(pos_y, pos_z, size_y, size_z float64) *XBox {
    return NewXBox(Coordinates{0.,  pos_y,  pos_z}, 
                   Coordinates{0., size_y, size_z})
}

func (self *XBox) SetPosition(c Coordinates) {
    c.X = -0.5 * math.MaxFloat64
    self.position = c
}

//YBox is a BBox which extends infinitely in the Y direction
type YBox struct {
    BBox
}

func NewYBox(position, size Coordinates) *YBox {
    size.Y = math.MaxFloat64
    position.Y = -0.5 * math.MaxFloat64
    r := YBox{BBox{position, size}}
    return &r
}

func NewYBox4f(pos_x, pos_z, size_x, size_z float64) *YBox {
    return NewYBox(Coordinates{ pos_x, 0.,  pos_z}, 
                   Coordinates{size_x, 0., size_z})
}

func (self *YBox) SetPosition(c Coordinates) {
    c.Y = -0.5 * math.MaxFloat64
    self.position = c
}

//ZBox is a BBox which extends infinitely in the Z direction
type ZBox struct {
    BBox
}

func NewZBox(position, size Coordinates) *ZBox {
    size.Z = math.MaxFloat64
    position.Z = -0.5 * math.MaxFloat64
    r := ZBox{BBox{position, size}}
    return &r
}

func NewZBox4f(pos_x, pos_y, size_x, size_y float64) *ZBox {
    return NewZBox(Coordinates{ pos_x,  pos_y, 0.}, 
                   Coordinates{size_x, size_y, 0.})
}

func (self *ZBox) SetPosition(c Coordinates) {
    c.Z = -0.5 * math.MaxFloat64
    self.position = c
}
