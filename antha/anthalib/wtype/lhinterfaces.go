// liquidhandling/lhinterfaces.go: Part of the Antha language
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

//An LHObject that can hold other LHObjects
type LHSlot interface {
	//GetChild get the contained object, nil if none
	GetChild() LHObject
	//SetChild set the contained object, error if it cannot
	SetChild(LHObject) error
	//Accepts can the slot accept the given object? (Can return true even if the slot is full)
	Accepts(LHObject) bool
	//GetChildPosition get the (absolute) position of the child object
	GetChildPosition() Coordinates
}

//WellReference used for specifying position within a well
type WellReference int

const (
	BottomReference WellReference = iota //0
	TopReference                         //1
	LiquidReference                      //2
)

//LHObject Provides a unified interface to physical size
//of items that can be placed on a liquid handler's deck
type LHObject interface {
	//GetBounds Return the absolute coordinates of the bounding box of the object
	GetBounds() BBox
	//SetOffset set the offset of the object relative to its parent (global if parent is nil)
	SetOffset(Coordinates)
	//SetParent Store the offset of the object
	SetParent(LHObject)
	//GetParent
	GetParent() LHObject
}

//Helper functions for objects as most are named and typed

//GetObjectName
func GetObjectName(o LHObject) string {
	if on, ok := o.(Named); ok {
		return on.GetName()
	}
	return "<unnamed>"
}

//GetObjectType
func GetObjectType(o LHObject) string {
	if ot, ok := o.(Typed); ok {
		return ot.GetType()
	}
	return "<untyped>"
}

//Addressable unifies the interface to objects which have
//sub-components that can be addressed by WellCoords (e.g. "A1")
//for example tip-boxes, plates, etc
type Addressable interface {
	//HasLocation Do the given coordinates exist in the object?
	HasLocation(WellCoords) bool
	//GetCoords Returns the object at the given well coords
	//nil if empty or position doesn't exist
	GetLocation(WellCoords) LHObject
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
