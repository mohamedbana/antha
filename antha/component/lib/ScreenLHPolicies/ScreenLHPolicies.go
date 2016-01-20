package ScreenLHPolicies

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

//"strconv"

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

	reactions := make([]*wtype.LHSolution, 0)

	policies, names := liquidhandling.PolicyMaker(liquidhandling.Allpairs, "DOE_run")

	for k := 0; k < len(_input.TestSols); k++ {
		for j := 0; j < _input.NumberofReplicates; j++ {
			for i := 0; i < len(policies); i++ {

				eachreaction := make([]*wtype.LHComponent, 0)

				_input.Diluent.Type = names[i]

				bufferSample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
				eachreaction = append(eachreaction, bufferSample)
				testSample := mixer.Sample(_input.TestSols[k], _input.TestSolVolume)
				eachreaction = append(eachreaction, testSample)
				reaction := execute.MixInto(_ctx,

					_input.OutPlate, eachreaction...)
				reactions = append(reactions, reaction)

			}
		}
	}
	_output.Reactions = reactions

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
	Diluent            *wtype.LHComponent
	NumberofReplicates int
	OutPlate           *wtype.LHPlate
	TestSolVolume      wunit.Volume
	TestSols           []*wtype.LHComponent
	TotalVolume        wunit.Volume
}

type Output_ struct {
	Reactions []*wtype.LHSolution
	Status    string
}
