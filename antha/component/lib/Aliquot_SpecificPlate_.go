package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _Aliquot_SpecificPlateRequirements() {

}

// Conditions to run on startup
func _Aliquot_SpecificPlateSetup(_ctx context.Context, _input *Aliquot_SpecificPlateInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Aliquot_SpecificPlateSteps(_ctx context.Context, _input *Aliquot_SpecificPlateInput, _output *Aliquot_SpecificPlateOutput) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		panic("Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofAliquots; i++ {
		if _input.Solution.TypeName() == "dna" {
			_input.Solution.Type = wtype.LTDoNotMix
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
		aliquot := execute.MixTo(_ctx, _input.OutPlate, "", 1, aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Aliquot_SpecificPlateAnalysis(_ctx context.Context, _input *Aliquot_SpecificPlateInput, _output *Aliquot_SpecificPlateOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Aliquot_SpecificPlateValidation(_ctx context.Context, _input *Aliquot_SpecificPlateInput, _output *Aliquot_SpecificPlateOutput) {

}
func _Aliquot_SpecificPlateRun(_ctx context.Context, input *Aliquot_SpecificPlateInput) *Aliquot_SpecificPlateOutput {
	output := &Aliquot_SpecificPlateOutput{}
	_Aliquot_SpecificPlateSetup(_ctx, input)
	_Aliquot_SpecificPlateSteps(_ctx, input, output)
	_Aliquot_SpecificPlateAnalysis(_ctx, input, output)
	_Aliquot_SpecificPlateValidation(_ctx, input, output)
	return output
}

func Aliquot_SpecificPlateRunSteps(_ctx context.Context, input *Aliquot_SpecificPlateInput) *Aliquot_SpecificPlateSOutput {
	soutput := &Aliquot_SpecificPlateSOutput{}
	output := _Aliquot_SpecificPlateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Aliquot_SpecificPlateNew() interface{} {
	return &Aliquot_SpecificPlateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Aliquot_SpecificPlateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Aliquot_SpecificPlateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Aliquot_SpecificPlateInput{},
			Out: &Aliquot_SpecificPlateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Aliquot_SpecificPlateElement struct {
	inject.CheckedRunner
}

type Aliquot_SpecificPlateInput struct {
	NumberofAliquots int
	OutPlate         string
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Aliquot_SpecificPlateOutput struct {
	Aliquots []*wtype.LHComponent
}

type Aliquot_SpecificPlateSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Aliquot_SpecificPlate",
		Constructor: Aliquot_SpecificPlateNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Aliquot/AliquotTo_PlateType.an",
			Params: []ParamDesc{
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
