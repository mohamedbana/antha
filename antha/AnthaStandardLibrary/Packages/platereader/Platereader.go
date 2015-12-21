// antha/AnthaStandardLibrary/Packages/Platereader/Platereader.go: Part of the Antha language
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

//Package containing functions for manipulating absorbance readings
package platereader

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func ReadAbsorbance(plate wtype.LHPlate, solution wtype.LHSolution, wavelength float64) (abs wtype.Absorbance) {
	abs.Reading = 0.0 // obviously placeholder
	abs.Wavelength = wavelength
	// add calculation to work out pathlength from volume and well geometry abs.Pathlength

	return abs
}

func Blankcorrect(blank wtype.Absorbance, sample wtype.Absorbance) (blankcorrected wtype.Absorbance) {

	if sample.Wavelength == blank.Wavelength &&
		sample.Pathlength == blank.Pathlength &&
		sample.Reader == blank.Reader {
		blankcorrected.Reading = sample.Reading - blank.Reading

		//currentstatus = make([]string,0)

		for _, status := range sample.Status {
			blankcorrected.Status = append(blankcorrected.Status, status)
		}
		blankcorrected.Status = append(blankcorrected.Status, "Blank Corrected")
	}
	return
}

func PathlengthCorrect(pathlength wunit.Length, reading wtype.Absorbance) (pathlengthcorrected wtype.Absorbance) {

	referencepathlength := wunit.NewLength(0.01, "m")

	pathlengthcorrected.Reading = reading.Reading * referencepathlength.SIValue() / pathlength.SIValue()
	return
}

//example

/*
func OD(Platetype wtype.LHPLate,wellvolume wtype.Volume,reading wtype.Absorbance) (od wtype.Absorbance){
volumetopathlengthconversionfactor := 0.0533//WellCrosssectionalArea
OD = (Blankcorrected_absorbance * 10/(total_volume*volumetopathlengthconversionfactor)// 0.0533 could be written as function of labware and liquid volume (or measureed height)
}

DCW = OD * ODtoDCWconversionfactor

*/
/*
type Absorbance struct {
	Reading    float64
	Wavelength float64
	Pathlength *wtype.Length
	Status     *[]string
	Reader     *string
}
*/
