// enzymebuffers.go Part of the Antha language
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

// Package for dealing with manipulation of buffers
package buffers

//"fmt"

type pH struct {
	pH        float64
	tempindeg float64 //wunit.Temperature
}

type SimpleBuffer struct {
	Components []string
	pH
}

type Buffer struct {
	Components []buffercomponent
	pH         pH
}

type buffercomponent struct {
	Molecule Molecule
	Molarity float64 //wunit.Moles g/Mol

}

type Molecule struct {
	Moleculename     string
	Molecularformula string
	Molecularweight  float64
	PubchemCID       int
}

var PotassiumAcetate50mM = buffercomponent{
	PotassiumAcetate,
	0.05,
}

var PotassiumAcetate = Molecule{
	"Potassium Acetate",
	"C2H3KO2",
	98.14232,
	31371,
}

var TrisAcetate = Molecule{
	"Potassium Acetate",
	"C4H11NO3",
	121.13504,
	6503,
}

var Cutsmartbuffer = SimpleBuffer{
	[]string{"PotassiumAcetate50mM",
		"TrisAcetate20mM",
		"MagnesiumAcetate10mM",
		"BSA100μgperml"},
	pH{7.9, 25.0},
}

var SapIstoragebuffer = SimpleBuffer{
	[]string{"300 mM NaCl", "10 mM Tris-HCl (pH 7.4)", "1 mM DTT", "0.1 mM EDTA", "50% Glycerol", "500 µg/ml BSA"},
	pH{7.4, 25.0},
}

func Newbuffer(buffer Buffer, diluent Buffer, dilution float64) (newbuffer Buffer) {

	newbuffer.Components = make([]buffercomponent, 0)
	for i := 0; i < len(buffer.Components); i++ {
		newbuffer.Components[i].Molecule = buffer.Components[i].Molecule
		newbuffer.Components[i].Molarity = (buffer.Components[i].Molarity * dilution)
		newbuffer.Components = append(newbuffer.Components, newbuffer.Components[i])
		for j := 0; j < len(diluent.Components); j++ {
			for _, newcomponent := range newbuffer.Components {
				if newcomponent.Molecule == diluent.Components[j].Molecule {
					newcomponent.Molarity = (newcomponent.Molarity + diluent.Components[j].Molarity*(1-dilution))
				}
				if newcomponent.Molecule != diluent.Components[j].Molecule {
					newbuffer.Components[i].Molecule = diluent.Components[j].Molecule
					newbuffer.Components[i].Molarity = (diluent.Components[j].Molarity * (1 - dilution))
					newbuffer.Components = append(newbuffer.Components, diluent.Components[j])
				}

			}
		}
	}
	newbuffer.pH = buffer.pH // this is incorrect and needs changing!!!
	return newbuffer
}
