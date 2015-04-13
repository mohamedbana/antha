// wtype/interface.go: Part of the Antha language
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
	"encoding/json"
	"github.com/antha-lang/antha/anthalib/wunit"
	"io/ioutil"
)

// map of matter types
var MatterLib map[string]GenericMatter

// Functions for dealing with matter
func MatterByName(name string) GenericMatter {
	if MatterLib == nil {
		MatterLib = MakeMatterLib()
	}

	return MatterLib[name]
}

// deserializes matter library from a JSON map structure
func GetMatterLib(fn string) (*(map[string]GenericMatter), error) {
	f, err := ioutil.ReadFile(fn)

	if err != nil {
		return nil, err
	}

	matter := make(map[string]GenericMatter, 100)
	e2 := json.Unmarshal(f, &matter)

	if e2 != nil {
		panic(e2)
	}
	return &matter, err
}

// check errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// make the initial matter library. This will eventually be deprecated.
func MakeMatterLib() map[string]GenericMatter {
	mtypes := []string{"polypropylene", "polycarbonate", "ptfe", "glass", "steel", "water", "glycerol", "ethanol", "surfactant", ""}
	mps := []float64{160.0, 292.0, 326.8, 1000.0, 1370.0, 0, 17.8, -114.0, 0.0, 0.0}
	bps := []float64{100000.0, 100000.0, 100000.0, 100000.0, 1000000.0, 100.0, 290.0, 78.37, 100.0, 0.0}
	shcs := []float64{900.0, 1200.0, 1300.0, 840.0, 420.0, 4181.0, 221.9, 2.46, 4181.0, 1.0}

	matter_map := make(map[string]GenericMatter, len(mtypes))

	// the following is necessary since we need an immutable pointer
	degreesc := "˚C"

	for i, t := range mtypes {
		mp := wunit.NewTemperature(mps[i], degreesc)
		bp := wunit.NewTemperature(bps[i], degreesc)
		shc := wunit.NewSpecificHeatCapacity(shcs[i], "J/kg")
		gm := GenericMatter{mtypes[i], mp, bp, shc}
		matter_map[t] = gm
	}

	return matter_map
}
