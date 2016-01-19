package PlateOut

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _requirements() {
}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	plateout := make([]*wtype.LHComponent, 0)

	if _input.Diluent != nil && _input.DilutionX > 1 {
		diluentsample := mixer.SampleForTotalVolume(_input.Diluent, _input.Plateoutvolume)
		plateout = append(plateout, diluentsample)
		// redeclare Plateoutvolume for adjusted volume to add of recovery mixture based on dilution ratio
		_input.Plateoutvolume = wunit.NewVolume(_input.Plateoutvolume.RawValue()/float64(_input.DilutionX), _input.Plateoutvolume.Unit().PrefixedSymbol())

	}
	plateoutsample := mixer.Sample(_input.RecoveredCells, _input.Plateoutvolume)
	plateout = append(plateout, plateoutsample)
	platedculture := execute.MixInto(_ctx,

		_input.AgarPlate, plateout...)
	execute.Incubate(_ctx,

		platedculture, _input.IncubationTemp, _input.IncubationTime, false)
	_output.Platedculture = platedculture

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
	AgarPlate      *wtype.LHPlate
	Diluent        *wtype.LHComponent
	DilutionX      int
	IncubationTemp wunit.Temperature
	IncubationTime wunit.Time
	Plateoutvolume wunit.Volume
	RecoveredCells *wtype.LHComponent
}

type Output_ struct {
	Platedculture *wtype.LHSolution
}
