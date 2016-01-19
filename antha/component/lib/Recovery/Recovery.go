package Recovery

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

//Recoveryvolume wunit.Volume
//= 2 (hours)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Recoverymedium *wtype.LHComponent

// Physical outputs from this protocol with types

func _requirements() {
}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	recoverymix := make([]*wtype.LHComponent, 0)

	transformedcellsComp := mixer.Sample(_input.Transformedcells, _input.TransformedcellVolume)

	//recoverymixture := mixer.Sample(Recoverymedium, Recoveryvolume)

	recoverymix = append(recoverymix, transformedcellsComp)
	//recoverymix = append(recoverymix,recoverymixture)

	recoverymix2 := execute.MixInto(_ctx,

		_input.OutPlate, recoverymix...)

	execute.Incubate(_ctx,

		recoverymix2, _input.Recoverytemp, _input.Recoverytime, true)

	_output.RecoveredCells = recoverymix2

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

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	OutPlate              *wtype.LHPlate
	Recoverytemp          wunit.Temperature
	Recoverytime          wunit.Time
	TransformedcellVolume wunit.Volume
	Transformedcells      *wtype.LHComponent
}

type Output_ struct {
	RecoveredCells *wtype.LHSolution
}
