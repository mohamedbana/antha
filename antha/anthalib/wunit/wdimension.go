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
// 2 Royal College St, London NW1 0NH UK

package wunit

import (
	"fmt"
	"github.com/antha-lang/antha/microArch/logger"
	"strings"
)

// length
type Length struct {
	ConcreteMeasurement
}

func EZLength(v float64) Length {
	return NewLength(v, "m")
}

// make a length
func NewLength(v float64, unit string) Length {
	l := Length{NewPMeasurement(v, unit)}

	// check

	if l.Unit().RawSymbol() != "m" {
		logger.Fatal("Base unit for lengths must be meters")
		panic("Base unit for lengths must be meters")
	}

	return l
}

// area
type Area struct {
	ConcreteMeasurement
}

// make an area unit
func NewArea(v float64, unit string) Area {
	if unit != "m^2" {
		logger.Fatal("Can't make areas which aren't square metres")
		panic("Can't make areas which aren't square metres")
	}

	a := Area{NewMeasurement(v, "", unit)}
	return a
}

// volume -- strictly speaking of course this is length^3
type Volume struct {
	ConcreteMeasurement
}

// make a volume
func NewVolume(v float64, unit string) Volume {
	if len(strings.TrimSpace(unit)) == 0 {
		logger.Fatal("Can't make Volumes without unit")
		panic("Can't make Volumes without unit")
	}
	o := Volume{NewPMeasurement(v, unit)}
	return o
}

func CopyVolume(v *Volume) *Volume {
	ret := NewVolume(v.RawValue(), v.Unit().PrefixedSymbol())
	return &ret
}

// temperature
type Temperature struct {
	ConcreteMeasurement
}

// make a temperature
func NewTemperature(v float64, unit string) Temperature {
	if unit != "˚C" && unit != "C" {
		logger.Fatal("Can't make temperatures which aren't in degrees C")
		panic("Can't make temperatures which aren't in degrees C")
	}
	t := Temperature{NewMeasurement(v, "", unit)}
	return t
}

// time
type Time struct {
	ConcreteMeasurement
}

// make a time unit
func NewTime(v float64, unit string) Time {
	if unit != "s" {
		logger.Fatal("Can't make temperatures which aren't in seconds")
		panic("Can't make temperatures which aren't in seconds")
	}

	t := Time{NewMeasurement(v, "", unit)}
	return t
}

// mass
type Mass struct {
	ConcreteMeasurement
}

// make a mass unit

func NewMass(v float64, unit string) Mass {
	return Mass{NewPMeasurement(v, unit)}
}

/*
func NewMass(v float64, unit string) Mass {
	if unit != "kg" && unit != "g" {
		panic("Can't make masses which aren't in grammes or kilograms")
	}

	var t Mass

	if unit == "kg" {
		t = Mass{NewMeasurement(v, "k", "g")}
	} else {
		t = Mass{NewMeasurement(v, "", "g")}
	}
	return t
}
*/
// defines mass to be a SubstanceQuantity
func (m *Mass) Quantity() Measurement {
	return m
}

// mole
type Amount struct {
	ConcreteMeasurement
}

// generate a new Amount in moles
func NewAmount(v float64, unit string) Amount {
	if unit != "M" {
		logger.Fatal("Can't make amounts which aren't in moles")
		panic("Can't make amounts which aren't in moles")
	}

	m := Amount{NewMeasurement(v, "", unit)}
	return m
}

// defines Amount to be a SubstanceQuantity
func (a *Amount) Quantity() Measurement {
	return a
}

// angle
type Angle struct {
	ConcreteMeasurement
}

// generate a new angle unit
func NewAngle(v float64, unit string) Angle {
	if unit != "radians" {
		logger.Fatal("Can't make angles which aren't in radians")
		panic("Can't make angles which aren't in radians")
	}

	a := Angle{NewMeasurement(v, "", unit)}
	return a
}

// this is really Mass(Length/Time)^2
type Energy struct {
	ConcreteMeasurement
}

// make a new energy unit
func NewEnergy(v float64, unit string) Energy {
	if unit != "J" {
		logger.Fatal("Can't make energies which aren't in Joules")
		panic("Can't make energies which aren't in Joules")
	}

	e := Energy{NewMeasurement(v, "", unit)}
	return e
}

// a Force
type Force struct {
	ConcreteMeasurement
}

// a new force in Newtons
func NewForce(v float64, unit string) Force {
	if unit != "N" {
		logger.Fatal("Can't make forces which aren't in Newtons")
		panic("Can't make forces which aren't in Newtons")
	}

	f := Force{NewMeasurement(v, "", unit)}
	return f
}

// a Pressure structure
type Pressure struct {
	ConcreteMeasurement
}

// make a new pressure in Pascals
func NewPressure(v float64, unit string) Pressure {
	if unit != "Pa" {
		logger.Fatal("Can't make pressures which aren't in Pascals")
		panic("Can't make pressures which aren't in Pascals")
	}

	p := Pressure{NewMeasurement(v, "", unit)}

	return p
}

// defines a concentration unit
type Concentration struct {
	ConcreteMeasurement
	//MolecularWeight *float64
}

// make a new concentration in SI units... either M/l or kg/l
func NewConcentration(v float64, unit string) Concentration {
	if unit != "g/l" && unit != "M/l" {
		// this should never be seen by users
		logger.Fatal("Can't make concentrations which aren't either Mol/l or g/l")
		panic("Can't make concentrations which aren't either Mol/l or g/l")
	}

	c := Concentration{NewMeasurement(v, "", unit)}
	return c
}

// mass or mole
type SubstanceQuantity interface {
	Quantity() Measurement
}

/*
type Protein interface {
	Molecule
	AASequence() string
}

type Enzyme struct {
	Class    string
	Synonyms []string
}
*/
/*
type Molecule interface {
	MolecularWeight() wtype.Mass
}

func (p *DNASequence) MolecularWeight() wtype.Mass {

}

func (d *ProteinSequence) MolecularWeight() wtype.Mass {

}
/*
func (e *Enzyme) AASequence() string {

}
*/
/*
// Sid's stuff

type Conc interface {
	AsMolar(mass Mass) MolarConcentration
	AsMass(mass Mass) MassConcentration
}

type MolarConcentration struct {
	Moles Amount
	Vol   Volume
}

func (m MolarConcentration) AsMolar(actualmass Mass) MolarConcentration {
	return m
}

// "M" and "g" need to be prefixed units to work!
func (m MolarConcentration) AsMass(mass Mass) MassConcentration {
	return MassConcentration{NewMass(m.Moles.ConvertTo(ParsePrefixedUnit("M"))*mass.ConvertTo(ParsePrefixedUnit("g")), "g"), m.Vol}
}

type MassConcentration struct {
	Mass Mass
	Vol  Volume
}

func (m MassConcentration) AsMolar(mass Mass) MolarConcentration {
	return MolarConcentration{NewAmount(m.Mass.SIValue()/mass.SIValue(), "M"), m.Vol}
}

func (m MassConcentration) AsMass(mass Mass) MassConcentration {
	return m
}
*/
/*
func (conc *Concentration)AddMolecularweight(molecularweight float64){
	conc.MolecularWeight = molecularweight
}
*/
func (conc *Concentration) GramPerL(molecularweight float64) (conc_g Concentration) {
	if conc.Munit.BaseSISymbol() == "g/l" {
		conc_g = *conc
	}
	if conc.Munit.BaseSISymbol() == "M/l" {
		conc_g = NewConcentration((conc.SIValue() * molecularweight), "M/l")
	}
	return conc_g
}

func (conc *Concentration) MolPerL(molecularweight float64) (conc_M Concentration) {
	if conc.Munit.BaseSISymbol() == "g/l" {
		conc_M = NewConcentration((conc.SIValue() / molecularweight), "g/l")
	}
	if conc.Munit.BaseSISymbol() == "M/l" {
		conc_M = *conc
	}
	return conc_M
}

/*
type Conc interface {
	AsMolar(gperl Concentration) MoleculeConcentration
	AsMass(Mperl float64) MoleculeConcentration
}

type MoleculeConcentration struct {
	Conc_gperl      Concentration
	Molecularweight float64 // Really a g/mol
	//Vol   Volume
}

const Avogadro = 6.0221417930 * 1E23 // number of molecules in a mol

func (m MoleculeConcentration) AsMolar(gperl Concentration) MoleculeConcentration {
	return m
}

func (m MoleculeConcentration) AsMass(mperl Concentration) MoleculeConcentration {
	return MoleculeConcentration{NewMass(m.Moles.ConvertTo("M")*mass.ConvertTo("g"), "g"), m.Vol}
}

func (m MassConcentration) AsMolar(mass Mass) MolarConcentration {
	return MolarConcentration{NewAmount(m.Mass.SIValue()/mass.SIValue(), "M"), m.Vol}
}

func (m MassConcentration) AsMass(mass Mass) MassConcentration {
	return m
}
*/
// a structure which defines a specific heat capacity
type SpecificHeatCapacity struct {
	ConcreteMeasurement
}

// make a new specific heat capacity structure in SI units
func NewSpecificHeatCapacity(v float64, unit string) SpecificHeatCapacity {
	if unit != "J/kg" {
		logger.Fatal("Can't make specific heat capacities which aren't in J/kg")
		panic("Can't make specific heat capacities which aren't in J/kg")
	}

	s := SpecificHeatCapacity{NewMeasurement(v, "", unit)}
	return s
}

// a structure which defines a density
type Density struct {
	ConcreteMeasurement
}

// make a new density structure in SI units
func NewDensity(v float64, unit string) Density {
	if unit != "kg/m^3" {
		logger.Fatal("Can't make densities which aren't in kg/m^3")
		panic("Can't make densities which aren't in kg/m^3")
	}

	d := Density{NewMeasurement(v, "", unit)}
	return d
}

type FlowRate struct {
	ConcreteMeasurement
}

// new flow rate in ml/min

func NewFlowRate(v float64, unit string) FlowRate {
	if unit != "ml/min" {
		logger.Fatal("Can't make flow rate which aren't in ml/min")
		panic("Can't make flow rate which aren't in ml/min")
	}
	fr := FlowRate{NewMeasurement(v, "", unit)}

	return fr
}

type Rate struct {
	ConcreteMeasurement
	Timeunit string //time.Duration
}

func (cm *Rate) ToString() string {
	return fmt.Sprintf("%-6.3f%s", cm.RawValue(), cm.Unit().PrefixedSymbol(), cm.Timeunit)
}
func NewRate(v float64, unit string, timeunit string) (r Rate, err error) {
	if unit != `/` {
		err = fmt.Errorf("Can't make flow rate which aren't in per")
		logger.Fatal(err.Error())
		panic(err.Error())
	}
	concrete := NewMeasurement(v, "", unit)

	approvedtimeunits := []string{"ns", "us", "µs", "ms", "s", "m", "h"}
	//Mvalue float64
	// the relevant units

	for _, approvedunit := range approvedtimeunits {
		if timeunit == approvedunit {
			r.Mvalue = concrete.Mvalue
			r.Munit = concrete.Munit
			r.Timeunit = timeunit
			return
		}
	}
	err = fmt.Errorf(timeunit, " Not approved time unit. Approved units time are: ", approvedtimeunits)
	return r, err
}
