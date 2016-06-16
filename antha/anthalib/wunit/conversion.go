// Part of the Antha language
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

// Core Antha package for dealing with units in Antha
package wunit

/*
type

func Splitunit(unit string)(numerators[]string, denominators[]string)

var conversiontable = map[string]map[string]float64{
	"density":map[string]float64{
		"g/L":
	}
}
*/

func MasstoVolume(m Mass, d Density) (v Volume) {
	mass := m.SIValue()
	density := d.SIValue()
	volume := mass / density // in m^3
	volume = volume * 1000   // in l
	v = NewVolume(mass, "l")
	return v
}

func VolumetoMass(v Volume, d Density) (m Mass) {
	//mass := m.SIValue()
	density := d.SIValue()

	volume := v.SIValue() //* 1000 // convert m^3 to l

	mass := volume * density // in m^3

	m = NewMass(mass, "kg")
	return m
}

func VolumeForTargetMass(targetmass Mass, startingconc Concentration) (v Volume) {

	v = NewVolume(float64((targetmass.SIValue()/startingconc.SIValue())*1000000), "ul")

	return v
}
