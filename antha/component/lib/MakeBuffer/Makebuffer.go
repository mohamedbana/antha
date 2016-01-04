package MakeBuffer

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _requirements() {

}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	//Bufferstockvolume := wunit.NewVolume((FinalVolume.SIValue() * FinalConcentration.SIValue()/Bufferstockconc.SIValue()),"l")

	_output.Buffer = execute.MixInto(_ctx,

		_input.OutPlate,
		mixer.Sample(_input.Bufferstock, _input.Bufferstockvolume),
		mixer.Sample(_input.Diluent, _input.Diluentvolume))

	_output.Status = fmt.Sprintln("Buffer stock volume = ", _input.Bufferstockvolume.ToString(), "of", _input.Bufferstock.CName,
		"was added to ", _input.Diluentvolume.ToString(), "of", _input.Diluent.CName,
		"to make ", _input.FinalVolume.ToString(), "of", _input.Buffername,
		"Buffer stock conc =", _input.Bufferstockconc.ToString())

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
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

/*
type Mole struct {
	number float64
}*/

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	Buffername         string
	Bufferstock        *wtype.LHComponent
	Bufferstockconc    wunit.Concentration
	Bufferstockvolume  wunit.Volume
	Diluent            *wtype.LHComponent
	Diluentname        string
	Diluentvolume      wunit.Volume
	FinalConcentration wunit.Concentration
	FinalVolume        wunit.Volume
	InPlate            *wtype.LHPlate
	OutPlate           *wtype.LHPlate
}

type Output_ struct {
	Buffer *wtype.LHSolution
	Status string
}
