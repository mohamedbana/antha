// wtype/genericsolid.go: Part of the Antha language
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

// defines a generic solid structure
type GenericSolid struct {
	GenericPhysical
	shape Shape
}

func (gs *GenericSolid) Shape() Shape {
	return gs.shape
}

// function for returning a blank generic solid - it has the basic type info
// but no mass, dimensions etc. etc.
// this will eventually be taken from database entries
func NewGenericSolid(mattertype string, shapetype string) *GenericSolid {
	gp := NewGenericPhysical(mattertype)
	shape := NewShape(shapetype)
	gs := GenericSolid{gp, shape}
	return &gs
}
