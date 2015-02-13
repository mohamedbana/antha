// wunit/wdimension.go: Part of the Antha language
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

package wunit

import (
)

// length
type Length struct{
	ConcreteMeasurement
}

// make a length 
func NewLength(v float64, unit string) Length{
	// need to enforce consistency at runtime
	// meanwhile let's bodge it

	if unit!="m"{
		panic("Can't make lengths which aren't metres")
	}

	l:=Length{NewMeasurement(v, "", unit)}

	return l
}

// area
type Area struct{
	ConcreteMeasurement
}

// make an area unit
func NewArea(v float64, unit string) Area{
	if unit!="m^2"{
		panic("Can't make areas which aren't square metres")
	}

	a:=Area{NewMeasurement(v, "", unit)}
	return a
}

// volume -- strictly speaking of course this is length^3
type Volume struct{
	ConcreteMeasurement
}

// make a volume
func NewVolume(v float64, unit string) Volume{
	if unit!="L" && unit!="l" && unit!="M^3"{
		panic("Can't make volumes which aren't Litres or cubic metres")
	}
	o:=Volume{NewMeasurement(v, "", unit)}
	return o
}

// temperature
type Temperature struct{
	ConcreteMeasurement
}

// make a temperature
func NewTemperature(v float64, unit string)Temperature{
	if unit!= "˚C"{
		panic("Can't make temperatures which aren't in degrees C")
	}
	t:=Temperature{NewMeasurement(v, "", unit)}
	return t
}

// time
type Time struct{
	ConcreteMeasurement
}

// make a time unit
func NewTime(v float64, unit string)Time{
	if unit!="s"{
		panic("Can't make temperatures which aren't in seconds")
	}

	t:=Time{NewMeasurement(v, "", unit)}
	return t
}

// mass
type Mass struct{
	ConcreteMeasurement
}

// make a mass unit
func NewMass(v float64, unit string)Mass{
	if unit!="kg" && unit!="g"{
		panic("Can't make masses which aren't in grammes or kilograms")
	}

	var t Mass

	if(unit=="kg"){
		t=Mass{NewMeasurement(v, "k", "g")}
	} else {
		t=Mass{NewMeasurement(v, "", "g")}
	}
	return t
}

// defines mass to be a SubstanceQuantity
func (m *Mass)Quantity() Measurement{
	return m
}

// mole
type Amount struct{
	ConcreteMeasurement
}

// generate a new Amount in moles
func NewAmount(v float64, unit string)Amount{
	if unit!="M"{
		panic ("Can't make amounts which aren't in moles")
	}

	m:=Amount{NewMeasurement(v, "", unit)}
	return m
}

// defines Amount to be a SubstanceQuantity
func (a *Amount)Quantity()Measurement{
	return a
}

// angle
type Angle struct{
	ConcreteMeasurement
}

// generate a new angle unit
func NewAngle(v float64, unit string)Angle{
	if unit!="radians"{
		panic("Can't make angles which aren't in radians")
	}

	a:=Angle{NewMeasurement(v, "", unit)}
	return a
}

// this is really Mass(Length/Time)^2
type Energy struct{
	ConcreteMeasurement
}

// make a new energy unit
func NewEnergy(v float64, unit string)Energy{
	if unit!="J"{
		panic("Can't make energies which aren't in Joules")
	}

	e:=Energy{NewMeasurement(v, "", unit)}
	return e
}

// mass or mole
type SubstanceQuantity interface{
	Quantity() Measurement
}

// a Force
type Force struct{
	ConcreteMeasurement
}

// a new force in Newtons
func NewForce(v float64, unit string) Force{
	if unit!="N"{
		panic("Can't make forces which aren't in Newtons")
	}

	f:=Force{NewMeasurement(v, "", unit)}
	return f
}

// a Pressure structure
type Pressure struct{
	ConcreteMeasurement
}

// make a new pressure in Pascals
func NewPressure(v float64, unit string)Pressure{
	if(unit!="Pa"){
		panic("Can't make pressures which aren't in Pascals")
	}

	p:=Pressure{NewMeasurement(v, "", unit)}

	return p
}

// defines a concentration unit
type Concentration struct{
	ConcreteMeasurement
}

// make a new concentration in SI units... either M/l or kg/l
func NewConcentration(v float64, unit string)Concentration{
	if(unit!="g/l" && unit !="M/l"){
		// this should never be seen by users
		panic ("Can't make concentrations which aren't either Mol/l or g/l")
	}

	c:=Concentration{NewMeasurement(v, "", unit)}
	return c
}

// a structure which defines a specific heat capacity
type SpecificHeatCapacity struct{
	ConcreteMeasurement
}

// make a new specific heat capacity structure in SI units
func NewSpecificHeatCapacity(v float64, unit string)SpecificHeatCapacity{
	if(unit!="J/kg"){
		panic("Can't make specific heat capacities which aren't in J/kg")
	}

	s:=SpecificHeatCapacity{NewMeasurement(v, "", unit)}
	return s
}

// a structure which defines a density
type Density struct{
	ConcreteMeasurement
}

// make a new density structure in SI units
func NewDensity(v float64, unit string)Density{
	if(unit!="kg/m^3"){
		panic("Can't make densities which aren't in kg/m^3")
	}

	d:=Density{NewMeasurement(v, "", unit)}
	return d
}
