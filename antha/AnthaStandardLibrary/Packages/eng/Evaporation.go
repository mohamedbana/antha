// antha/AnthaStandardLibrary/Packages/eng/Evaporation.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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

// Package containing formulae for the estimation of evaporation times based upon thermodynamics and empirical equations

package eng

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math"
)

func Θ(liquid string, airvelocity wunit.Velocity) float64 {

	return (liquidclasses.Liquidclass[liquid]["c"]) + ((liquidclasses.Liquidclass[liquid]["d"]) * airvelocity.SIValue())
}

//Some functions to calculate evaporation
func Pws(Temp wunit.Temperature) float64 {
	var tempinKelvin float64
	if Temp.Unit().RawSymbol() == "K" {
		tempinKelvin = Temp.RawValue()
	} else if Temp.Unit().RawSymbol() == "C" || Temp.Unit().RawSymbol() == "˚C" {
		tempinKelvin = (Temp.SIValue() + 273.15)
	}
	return (math.Pow(math.E, (77.3450+(0.0057*tempinKelvin)-7235/tempinKelvin)) / math.Pow(tempinKelvin, 8.2))
}

func Pw(Relativehumidity float64, PWS float64) float64 {
	return (Relativehumidity * PWS)
}

func Xs(pws float64, Pa wunit.Pressure) float64 {
	return (0.62198 * pws / (Pa.SIValue() - pws))
}

func X(pw float64, Pa wunit.Pressure) float64 {
	return (0.62198 * pw / (Pa.SIValue() - pw))
}
