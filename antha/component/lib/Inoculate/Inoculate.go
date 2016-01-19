// Example inoculation protocol.
// Inoculates seed culture into fresh media (and logs conditions?)
// TODO: in progress from edited bradford protocol
package Inoculate

import (
	// "liquid handler"
	//"labware"
	//"OD"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// we do comments like this

// Input parameters for this protocol (data)

//= uL(25)
//= uL(475)
//= mgperml (100)
//= mgperml  (0.1)
//= 0 // Note: 1 replicate means experiment is in duplicate, etc.

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//suspension which contains living cells

// type buffer which could have a concentration automatically?

// Physical outputs from this protocol with types

func _requirements() {
	// None
}

func _setup(_ctx context.Context, _input *Input_) {
	//none
}

func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	//antibiotic_volume  := wunit.NewVolume(Media_volume.SIValue() * (Desiredantibioticconcentration.SIValue()/Antibioticstockconc.SIValue()),"l")

	media_with_antibiotic := _input.Media
	//media_with_antibiotic := mixer.Mix(mixer.Sample(Antibiotic,antibiotic_volume), mixer.Sample(Media,Media_volume))
	_output.Inoculated_culture = execute.MixInto(_ctx,

		_input.OutPlate, mixer.Sample(_input.Seed, _input.Seed_volume), mixer.Sample(media_with_antibiotic, _input.Media_volume))
}

//should the transfer to thermomixer/incubator command be included in this protocol or in a separate protocol
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
	//OD_at_inoculation = OD.Inoculated_culture // need to know signatures of protocol_OD I,O,Q - function signature
}

func _validation(_ctx context.Context, _input *Input_, _output *Output_) {
	/*
		if OD.sample_absorbance > 1 {
		panic("Sample likely needs further dilution")
		}
		if OD.sample_absorbance < 0.02 {
		warn("low inoculation OD")
		//could add visual (i.e. manual or camera based) validation
		// TODO: add test of replicate variance
		}*/
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
	Antibiotic                     *wtype.LHComponent
	Antibioticstockconc            wunit.Concentration
	Desiredantibioticconcentration wunit.Concentration
	Media                          *wtype.LHComponent
	Media_volume                   wunit.Volume
	OutPlate                       *wtype.LHPlate
	Replicate_count                int
	Seed                           *wtype.LHComponent
	Seed_volume                    wunit.Volume
}

type Output_ struct {
	Inoculated_culture *wtype.LHSolution
	OD_at_inoculation  float64
}
