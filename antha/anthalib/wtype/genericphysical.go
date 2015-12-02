// wtype/genericphysical.go: Part of the Antha language
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

// GenericPhysical structure: holds data items required to define a physical object
type GenericPhysical struct {
	GenericMatter
	Myname string
	Mymass wunit.Mass
	Myvol  wunit.Volume
	Mytemp wunit.Temperature
}

func (gp *GenericPhysical) Name() string {
	return gp.Myname
}

func (gp *GenericPhysical) SetName(s string) string {
	oldMyname := gp.Myname
	gp.Myname = s
	return oldMyname
}

func NewGenericPhysical(mattertype string) GenericPhysical {
	gp := GenericPhysical{MatterByName(mattertype), mattertype, wunit.NewMass(0.0, "g"), wunit.NewVolume(0.0, "L"), wunit.NewTemperature(0.0, "˚C")}
	return gp
}

func (gp *GenericPhysical) Clone() GenericPhysical {
	return GenericPhysical{gp.GenericMatter.Clone(), gp.Name(), gp.Mymass, gp.Myvol, gp.Mytemp}
}

func (gp *GenericPhysical) Mass() wunit.Mass {
	return gp.Mymass
}

func (gp *GenericPhysical) SetMass(m wunit.Mass) wunit.Mass {
	om := gp.Mymass
	gp.Mymass = m
	return om
}

func (gp *GenericPhysical) Volume() wunit.Volume {
	return gp.Myvol
}

func (gp *GenericPhysical) SetVolume(v wunit.Volume) wunit.Volume {
	ov := gp.Myvol
	gp.Myvol = v
	return ov
}

func (gp *GenericPhysical) Density() wunit.Density {
	densgperl := gp.Mymass.SIValue() / gp.Myvol.SIValue()
	// g/l == kg/m^3
	return wunit.NewDensity(densgperl, "kg/m^3")
}

/*
func (gp *GenericPhysical)Location() coordinates{
	return gp.loc
}

func (gp *GenericPhysical)SetLocation(c coordinates){
	gp.loc = c
}
*/

func (gp *GenericPhysical) Temperature() wunit.Temperature {
	return gp.Mytemp
}

func (gp *GenericPhysical) SetTemperature(t wunit.Temperature) {
	gp.Mytemp = t
}

func (gp *GenericPhysical) Density() wunit.Density {
	m := gp.Mass()
	v := gp.Volume()
	return wunit.NewDensity(m.SIValue()/v.SIValue(), "kg/m^3")
}
