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

//WellReference used for specifying position within a well
type WellReference int

const (
    BottomReference WellReference = iota //0
    TopReference                         //1
    LiquidReference                      //2
)

//LHDeckObject Provides a unified interface to physical properties
//of items that can be placed on a liquid handler's deck, 
//currently LHPlate, LHTipbox, and LHTipWaste
type LHDeckObject interface {
    //GetSize Return the physical size of the object in mm
    GetSize() Coordinates
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

