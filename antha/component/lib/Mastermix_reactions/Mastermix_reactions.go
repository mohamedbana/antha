package Mastermix_reactions

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// if buffer is being added

// add as many as possible option

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// optional if nil this is ignored

// Physical outputs from this protocol with types

func _requirements() {
}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	// work out volume to top up to in each case (per reaction) in l:
	topupVolumeperreacttion := _input.TotalVolumeperreaction.SIValue() - _input.VolumetoLeaveforDNAperreaction.SIValue()

	// multiply by number of reactions per mastermix
	topupVolume := wunit.NewVolume(float64(_input.Reactionspermastermix)*topupVolumeperreacttion, "l")

	if len(_input.Components) != len(_input.ComponentVolumesperReaction) {
		panic("len(Components) != len(OtherComponentVolumes)")
	}

	mastermixes := make([]*wtype.LHSolution, 0)

	if _input.AliquotbyRow {
		panic("MixTo based method coming soon!")
	} else {
		for i := 0; i < _input.NumberofMastermixes; i++ {

			eachmastermix := make([]*wtype.LHComponent, 0)

			if _input.TopUpBuffer != nil {
				bufferSample := mixer.SampleForTotalVolume(_input.TopUpBuffer, topupVolume)
				eachmastermix = append(eachmastermix, bufferSample)
			}

			for k, component := range _input.Components {
				if k == len(_input.Components) {
					component.Type = "NeedToMix"
				}

				// multiply volume of each component by number of reactions per mastermix
				adjustedvol := wunit.NewVolume(float64(_input.Reactionspermastermix)*_input.ComponentVolumesperReaction[k].SIValue(), "l")

				componentSample := mixer.Sample(component, adjustedvol)
				eachmastermix = append(eachmastermix, componentSample)
			}

			mastermix := execute.MixInto(_ctx,

				_input.OutPlate, eachmastermix...)
			mastermixes = append(mastermixes, mastermix)

		}

	}
	_output.Mastermixes = mastermixes

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
	AliquotbyRow                   bool
	ComponentVolumesperReaction    []wunit.Volume
	Components                     []*wtype.LHComponent
	Inplate                        *wtype.LHPlate
	NumberofMastermixes            int
	OutPlate                       *wtype.LHPlate
	Reactionspermastermix          int
	TopUpBuffer                    *wtype.LHComponent
	TotalVolumeperreaction         wunit.Volume
	VolumetoLeaveforDNAperreaction wunit.Volume
}

type Output_ struct {
	Mastermixes []*wtype.LHSolution
	Status      string
}
