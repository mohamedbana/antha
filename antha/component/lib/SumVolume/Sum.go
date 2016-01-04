package SumVolume

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"github.com/antha-lang/antha/antha/anthalib/wunit"
// Input parameters for this protocol

//D Concentration
//E float64

// Data which is returned from this protocol

//DmolarConc wunit.MolarConcentration

// Physical inputs to this protocol

// Physical outputs from this protocol

func _requirements() {

}

// Actions to perform before protocol itself
func _setup(_ctx context.Context, _input *Input_) {

}

// Core process of the protocol: steps to be performed for each input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	//var Dmassconc wunit.MassConcentration = D

	/*	molarmass := wunit.NewAmount(E,"M")

		var Dnew = wunit.MoleculeConcentration{D,E}

		mass := wunit.NewMass(1,"g")

		DmolarConc = Dnew.AsMolar(mass)
	*/
	_output.Sum = *(wunit.CopyVolume(&_input.A))
	(&_output.Sum).Add(&_input.B)
	_output.Status = fmt.Sprintln(
		"Sum of", _input.A.ToString(), "and", _input.B.ToString(), "=", _output.Sum.ToString(), "Temp=", _input.C.ToString(),
	) //"D Concentration in g/l", D, "D concentration in M/l", DmolarConc)
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
	A wunit.Volume
	B wunit.Volume
	C wunit.Temperature
}

type Output_ struct {
	Status string
	Sum    wunit.Volume
}
