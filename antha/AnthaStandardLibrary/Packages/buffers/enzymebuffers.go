// enzymebuffers.go
package buffers

import (
//"fmt"
)

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
