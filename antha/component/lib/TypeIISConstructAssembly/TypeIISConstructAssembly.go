package TypeIISConstructAssembly

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _requirements() {}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	bufferSample := mixer.Sample(_input.Buffer, _input.BufferVol)
	samples = append(samples, bufferSample)

	atpSample := mixer.Sample(_input.Atp, _input.AtpVol)
	samples = append(samples, atpSample)

	//vectorSample := mixer.Sample(Vector, VectorVol)
	vectorSample := mixer.Sample(_input.Vector, _input.VectorVol)
	samples = append(samples, vectorSample)

	for k, part := range _input.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartNames[k], " vol ", _input.PartVols[k])
		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	reSample := mixer.Sample(_input.RestrictionEnzyme, _input.ReVol)
	samples = append(samples, reSample)

	ligSample := mixer.Sample(_input.Ligase, _input.LigVol)
	samples = append(samples, ligSample)

	_output.Reaction = execute.MixInto(_ctx,

		_input.OutPlate, samples...)

	// incubate the reaction mixture
	execute.Incubate(_ctx,

		_output.Reaction, _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	execute.Incubate(_ctx,

		_output.Reaction, _input.InactivationTemp, _input.InactivationTime, false)
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
	Atp                *wtype.LHComponent
	AtpVol             wunit.Volume
	Buffer             *wtype.LHComponent
	BufferVol          wunit.Volume
	InPlate            *wtype.LHPlate
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LigVol             wunit.Volume
	Ligase             *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReVol              wunit.Volume
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	RestrictionEnzyme  *wtype.LHComponent
	Vector             *wtype.LHComponent
	VectorVol          wunit.Volume
	Water              *wtype.LHComponent
}

type Output_ struct {
	Reaction *wtype.LHSolution
}
