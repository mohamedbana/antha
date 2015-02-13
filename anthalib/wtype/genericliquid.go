// wtype/genericliquid.go: Part of the Antha language
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

import(
	"github.com/antha-lang/antha/anthalib/wunit"
)

// Structure which defines a generic liquid 
type GenericLiquid struct{
	GenericPhysical
	viscosity float64
}

// factory method for creating a new generic liquid
func NewGenericLiquid(name string, mattertype string, volume wunit.Volume) *GenericLiquid{
	gp:=NewGenericPhysical(mattertype)
	gp.SetVolume(volume)
	gp.SetName(name)			// bit of a fudge to allow us to extend the basic types 
	gl:=GenericLiquid{gp, 0.000894}
	return &gl
}


func (gl *GenericLiquid)Clone()GenericLiquid{
	return GenericLiquid{gl.GenericPhysical.Clone(), gl.Viscosity()}
}


func (gl *GenericLiquid) Viscosity() float64{
	return gl.viscosity
}

// sample method for a generic liquid 
func (gl *GenericLiquid)Sample(v wunit.Volume) Liquid{
	// check we have enough

	// this mechanism is necessary to allow us to take the address of the receiver
	glv:=gl.Volume()

	if glv.SIValue() < v.SIValue(){
		panic("ERROR taking sample -- insufficient volume")
	}
	
	// we need to agree on the underlying representation
	// work out the value of 'v' in the same units as gl's volume

	sv:=v.SIValue() / glv.Munit.BaseSIConversionFactor()
	nv:=wunit.Volume{wunit.NewMeasurement((glv.RawValue() - sv), glv.Unit().Prefix().Name, glv.Unit().Symbol())}
	gl.SetVolume(nv)

	// make the new liquid

	ngl:=gl.Clone()
	ngl.SetVolume(v)
	return &ngl
}
