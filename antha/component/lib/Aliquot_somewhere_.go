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

func _Aliquot_somewhereRequirements() {

}

// Conditions to run on startup
func _Aliquot_somewhereSetup(_ctx context.Context, _input *Aliquot_somewhereInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Aliquot_somewhereSteps(_ctx context.Context, _input *Aliquot_somewhereInput, _output *Aliquot_somewhereOutput) {

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
		aliquot := execute.Mix(_ctx, aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Aliquot_somewhereAnalysis(_ctx context.Context, _input *Aliquot_somewhereInput, _output *Aliquot_somewhereOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Aliquot_somewhereValidation(_ctx context.Context, _input *Aliquot_somewhereInput, _output *Aliquot_somewhereOutput) {

}
func _Aliquot_somewhereRun(_ctx context.Context, input *Aliquot_somewhereInput) *Aliquot_somewhereOutput {
	output := &Aliquot_somewhereOutput{}
	_Aliquot_somewhereSetup(_ctx, input)
	_Aliquot_somewhereSteps(_ctx, input, output)
	_Aliquot_somewhereAnalysis(_ctx, input, output)
	_Aliquot_somewhereValidation(_ctx, input, output)
	return output
}

func Aliquot_somewhereRunSteps(_ctx context.Context, input *Aliquot_somewhereInput) *Aliquot_somewhereSOutput {
	soutput := &Aliquot_somewhereSOutput{}
	output := _Aliquot_somewhereRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Aliquot_somewhereNew() interface{} {
	return &Aliquot_somewhereElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Aliquot_somewhereInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Aliquot_somewhereRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Aliquot_somewhereInput{},
			Out: &Aliquot_somewhereOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Aliquot_somewhereElement struct {
	inject.CheckedRunner
}

type Aliquot_somewhereInput struct {
	NumberofAliquots int
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Aliquot_somewhereOutput struct {
	Aliquots []*wtype.LHComponent
}

type Aliquot_somewhereSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Aliquot_somewhere",
		Constructor: Aliquot_somewhereNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Aliquot/Aliquot_somewhereorother.an",
			Params: []ParamDesc{
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
