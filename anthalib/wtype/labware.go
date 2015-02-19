// wtype/labware.go: Part of the Antha language
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
//	"github.com/antha-lang/antha/anthalib/wunit"
)

// general interface applicable to all labware
type Labware interface {
	Entity
	Manufacturer() string
	LabwareType() string
}

// defines microplates. Microplates have wells.
type Plate interface {
	Labware
	Wells() [][]Well
	WellAt(crds WellCoords) Well
	WellsX() int
	WellsY() int
}

// A generic to define an SBS format plate
type GenericSBSFormatPlate struct {
	GenericEntity
	Manufr  string
	LType   string
	WellArr [][]Well
}

// find the first empty well and add this to it
func (gl *GenericSBSFormatPlate) Add(p Physical) {
	gl.FirstEmptyWell().Add(p)
}

/*
func (gl *GenericSBSFormatPlate)Material() Matter{
	return &(gl.GenericMatter)
}
*/

func (gl *GenericSBSFormatPlate) Manufacturer() string {
	return gl.Manufr
}

func (gl *GenericSBSFormatPlate) LabwareType() string {
	return gl.LType
}

func (gl *GenericSBSFormatPlate) Wells() [][]Well {
	return gl.WellArr
}

func (gl *GenericSBSFormatPlate) WellAt(crds WellCoords) Well {
	return gl.WellArr[crds.X][crds.Y]
}

// find the first empty well in the plate
func (gl *GenericSBSFormatPlate) FirstEmptyWell() Well {
	for _, wa := range gl.WellArr {
		for _, w := range wa {
			if w.Empty() {
				return w
			}
		}
	}
	return nil
}
