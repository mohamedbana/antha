package lib

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

func _Aliquot_SamePositionsMultipleplatesRequirements() {

}

// Conditions to run on startup
func _Aliquot_SamePositionsMultipleplatesSetup(_ctx context.Context, _input *Aliquot_SamePositionsMultipleplatesInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Aliquot_SamePositionsMultipleplatesSteps(_ctx context.Context, _input *Aliquot_SamePositionsMultipleplatesInput, _output *Aliquot_SamePositionsMultipleplatesOutput) {

	aliquots := make([]*wtype.LHComponent, 0)

	for i := 1; i < _input.NumberofPlates+1; i++ {

		for _, position := range _input.Positions {
			if _input.Solution.TypeName() == "dna" {
				_input.Solution.Type = wtype.LTDoNotMix
			}
			aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
			aliquot := execute.MixTo(_ctx, _input.OutPlate, position, i, aliquotSample)
			aliquots = append(aliquots, aliquot)
		}
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Aliquot_SamePositionsMultipleplatesAnalysis(_ctx context.Context, _input *Aliquot_SamePositionsMultipleplatesInput, _output *Aliquot_SamePositionsMultipleplatesOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Aliquot_SamePositionsMultipleplatesValidation(_ctx context.Context, _input *Aliquot_SamePositionsMultipleplatesInput, _output *Aliquot_SamePositionsMultipleplatesOutput) {

}
func _Aliquot_SamePositionsMultipleplatesRun(_ctx context.Context, input *Aliquot_SamePositionsMultipleplatesInput) *Aliquot_SamePositionsMultipleplatesOutput {
	output := &Aliquot_SamePositionsMultipleplatesOutput{}
	_Aliquot_SamePositionsMultipleplatesSetup(_ctx, input)
	_Aliquot_SamePositionsMultipleplatesSteps(_ctx, input, output)
	_Aliquot_SamePositionsMultipleplatesAnalysis(_ctx, input, output)
	_Aliquot_SamePositionsMultipleplatesValidation(_ctx, input, output)
	return output
}

func Aliquot_SamePositionsMultipleplatesRunSteps(_ctx context.Context, input *Aliquot_SamePositionsMultipleplatesInput) *Aliquot_SamePositionsMultipleplatesSOutput {
	soutput := &Aliquot_SamePositionsMultipleplatesSOutput{}
	output := _Aliquot_SamePositionsMultipleplatesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Aliquot_SamePositionsMultipleplatesNew() interface{} {
	return &Aliquot_SamePositionsMultipleplatesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Aliquot_SamePositionsMultipleplatesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Aliquot_SamePositionsMultipleplatesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Aliquot_SamePositionsMultipleplatesInput{},
			Out: &Aliquot_SamePositionsMultipleplatesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Aliquot_SamePositionsMultipleplatesElement struct {
	inject.CheckedRunner
}

type Aliquot_SamePositionsMultipleplatesInput struct {
	NumberofPlates   int
	OutPlate         string
	Positions        []string
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Aliquot_SamePositionsMultipleplatesOutput struct {
	Aliquots []*wtype.LHComponent
}

type Aliquot_SamePositionsMultipleplatesSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Aliquot_SamePositionsMultipleplates",
		Constructor: Aliquot_SamePositionsMultipleplatesNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Aliquot/AliquotTo_samepositionmultipleplates.an",
			Params: []ParamDesc{
				{Name: "NumberofPlates", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Parameters"},
				{Name: "Positions", Desc: "", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
