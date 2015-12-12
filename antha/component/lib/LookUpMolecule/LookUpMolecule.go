// example of how to look up molecule properties from pubchem
package LookUpMolecule

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Name of compound or array of multiple compounds

// molecule type is returned consisting of name, formula, molecular weight and chemical ID (CID)

// or JSON structure if preferred

// status to be printed out in manual driver console

func _requirements() {
}
func _setup(_ctx context.Context, _input *Input_) {
}
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	// method of making molecule from name
	_output.Compoundprops = pubchem.MakeMolecule(_input.Compound)

	// or returning properties in JSON structure
	_output.Jsonstring = pubchem.Compoundproperties(_input.Compound)

	// method of making a list of compounds from names
	_output.List = pubchem.MakeMolecules(_input.Compoundlist)

	// Print out status
	_output.Status = fmt.Sprintln("Returned data from",
		_input.Compound, "=",
		_output.Compoundprops.Moleculename,
		_output.Compoundprops.MolecularWeight,
		_output.Compoundprops.MolecularFormula,
		_output.Compoundprops.CID,
		"Data in JSON format =", _output.Jsonstring,
		"List=", _output.List)
}
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
	Compound     string
	Compoundlist []string
}

type Output_ struct {
	Compoundprops pubchem.Molecule
	Jsonstring    string
	List          []pubchem.Molecule
	Status        string
}
