//Some examples functions
// Calculate rate of reaction, V, of enzyme displaying Micahelis-Menten kinetics with Vmax, Km and [S] declared
// Calculating [S] and V from g/l concentration and looking up molecular weight of named substrate
// Calculating [S] and V from g/l concentration of DNA of known sequence
// Calculating [S] and V from g/l concentration of Protein product of DNA of known sequence

package Datacrunch

import (
	"fmt"
	//"math"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol

//Amount
// i.e. Moles, M

//Amount

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _requirements() {

}

// Actions to perform before protocol itself
func _setup(_ctx context.Context, _input *Input_) {

}

// Core process of the protocol: steps to be performed for each input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	// Work out rate of reaction, V of enzyme with Michaelis-Menten kinetics and [S], Km and Vmax declared
	//Using declared values for S and unit of S
	km := wunit.NewAmount(_input.Km, _input.Kmunit) //.SIValue()
	s := wunit.NewAmount(_input.S, _input.Sunit)    //.SIValue()

	_output.V = ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate based on conc and looking up molecular weight in pubchem

	// Look up properties
	substrate_mw := pubchem.MakeMolecule(_input.Substrate_name)

	// calculate moles
	submoles := sequences.Moles(_input.SubstrateConc, substrate_mw.MolecularWeight, _input.SubstrateVol)
	// calculate molar concentration
	submolarconc := sequences.GtoMolarConc(_input.SubstrateConc, substrate_mw.MolecularWeight)

	// make a new amount
	s = wunit.NewAmount(submolarconc, "M")

	// use michaelis menton equation
	v_substrate_name := ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate from DNA Sequence
	// calculate molar concentration
	dna_mw := sequences.MassDNA(_input.DNA_seq, false, false)
	dnamolarconc := sequences.GtoMolarConc(_input.DNAConc, dna_mw)

	// make a new amount
	s = wunit.NewAmount(dnamolarconc, "M")

	// use michaelis menton equation
	v_dna := ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate from Protein product of dna Sequence

	// translate
	orf, orftrue := sequences.FindORF(_input.DNA_seq)
	var protein_mw float64
	if orftrue == true {
		protein_mw_kDA := sequences.Molecularweight(orf)
		protein_mw = protein_mw_kDA * 1000
		_output.Orftrue = orftrue
	}

	// calculate molar concentration
	proteinmolarconc := sequences.GtoMolarConc(_input.ProteinConc, protein_mw)

	// make a new amount
	s = wunit.NewAmount(submolarconc, "M")

	// use michaelis menton equation
	v_protein := ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// print report
	_output.Status = fmt.Sprintln(
		"Rate, V of enzyme at substrate conc", _input.S, _input.Sunit,
		"of enzyme with Km", km.ToString(),
		"and Vmax", _input.Vmax, _input.Vmaxunit,
		"=", _output.V, _input.Vunit, ".",
		"Substrate =", _input.Substrate_name, ". We have", _input.SubstrateVol.ToString(), "of", _input.Substrate_name, "at concentration of", _input.SubstrateConc.ToString(),
		"Therefore... Moles of", _input.Substrate_name, "=", submoles, "Moles.",
		"Molar Concentration of", _input.Substrate_name, "=", submolarconc, "Mol/L.",
		"Rate, V = ", v_substrate_name, _input.Vmaxunit,
		"Substrate =", "DNA Sequence of", _input.Gene_name, "We have", "concentration of", _input.DNAConc.ToString(),
		"Therefore... Molar conc", "=", dnamolarconc, "Mol/L",
		"Rate, V = ", v_dna, _input.Vmaxunit,
		"Substrate =", "protein from DNA sequence", _input.Gene_name, ".",
		"We have", "concentration of", _input.ProteinConc.ToString(),
		"Therefore... Molar conc", "=", proteinmolarconc, "Mol/L",
		"Rate, V = ", v_protein, _input.Vmaxunit)
}

// Actions to perform after steps block to analyze data
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {

}

func _validation(_ctx context.Context, _input *Input_, _output *Output_) {

}

func _run(_ctx context.Context, value inject.Value) (inject.Value, error) {
	input := &Input_{}
	output := &Output_{}
	if err := inject.Assign(value, input); err != nil {
		return nil, err
	}
	_setup(_ctx, input)
	_steps(_ctx, input, output)
	_analysis(_ctx, input, output)
	_validation(_ctx, input, output)
	return inject.MakeValue(output), nil
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

func New() interface{} {
	return &Element_{
		inject.CheckedRunner{
			RunFunc: _run,
			In:      &Input_{},
			Out:     &Output_{},
		},
	}
}

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	DNAConc        wunit.Concentration
	DNA_seq        string
	Gene_name      string
	Km             float64
	Kmunit         string
	ProteinConc    wunit.Concentration
	S              float64
	SubstrateConc  wunit.Concentration
	SubstrateVol   wunit.Volume
	Substrate_name string
	Sunit          string
	Vmax           float64
	Vmaxunit       string
	Vunit          string
}

type Output_ struct {
	Orftrue bool
	Status  string
	V       float64
}
