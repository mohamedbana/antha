// wunit/interface.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// units mapped by string
var unitMap map[string]GenericUnit

// deserialize JSON prefix library
func GetPrefixLib(fn string) (*(map[string]SIPrefix), error) {
	f, err := ioutil.ReadFile(fn)

	if err != nil {
		return nil, err
	}

	prefices := make(map[string]SIPrefix, 20)
	json.Unmarshal(f, &prefices)
	return &prefices, err
}

// deserialize JSON unit library
func GetUnitLib(fn string) (*(map[string]GenericUnit), error) {
	f, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	units := make(map[string]GenericUnit, 20)
	e2 := json.Unmarshal(f, &units)

	if e2 != nil {
		panic(e2)
	}

	for k, v := range units {
		fmt.Println(k, " ", v)
	}

	return &units, err
}

// helper function to make it easier to
// make a new unit with prefix directly
func NewPrefixedUnit(prefix string, unit string) PrefixedUnit {
	u := UnitBySymbol(unit)
	p := SIPrefixBySymbol(prefix)

	/*
		if p==nil{
			panic(fmt.Sprintf("Can't instantiate this prefix: %s", prefix))
		}else if u==nil{
			panic(fmt.Sprintf("Can't instantiate this unit: %s", unit))
		}
	*/

	gpu := GenericPrefixedUnit{u, p}
	return &gpu
}

// get a unit from a string

func ParsePrefixedUnit(unit string) PrefixedUnit {
	parser := &SIPrefixedUnitGrammar{Buffer: unit}
	parser.Init()
	parser.SIPrefixedUnit.Init([]byte(unit))

	if err := parser.Parse(); err != nil {
		panic(err)
	}

	parser.Execute()

	return NewPrefixedUnit(parser.TreeTop.Children[0].Value.(string), parser.TreeTop.Children[1].Value.(string))
}

// look up unit by symbol
func UnitBySymbol(sym string) GenericUnit {
	if unitMap == nil {
		unitMap = Make_units()
	}

	return unitMap[sym]
}

// generate an initial unit library
func Make_units() map[string]GenericUnit {
	units := []string{"M", "m", "l", "L", "g", "V", "J", "A", "N", "s", "radians", "degrees", "rads", "Hz", "rpm", "˚C", "M/l", "g/l", "J/kg"}
	unitnames := []string{"mole", "minute", "litre", "litre", "Gramme", "Volt", "Joule", "Ampere", "Newton", "second", "radian", "degree", "radian", "Herz", "revolutions per minute", "Celsius", "Mol/litre", "g/litre", "Joule/kilogram"}
	//unitdimensions:=[]string{"amount", "time", "length^3", "length^3", "mass", "mass*length/time^2*charge", "mass*length^2/time^2", "charge/time", "charge", "mass*length/time^2", "time", "angle", "angle", "angle", "time^-1", "angle/time", "temperature"}

	unitbaseconvs := []float64{1, 0.1666666666666666667, 1, 1, 0.001, 1, 1, 1, 1, 1, 1, 0.01745329251994, 1, 1, 1, 1, 1, 1, 1}

	unit_map := make(map[string]GenericUnit, len(units))

	for i, u := range units {
		baseunit := u

		if u == "g" {
			baseunit = "kg"
		}
		gu := GenericUnit{unitnames[i], u, unitbaseconvs[i], baseunit}
		unit_map[u] = gu
	}

	return unit_map
}
