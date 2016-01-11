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

func _AliquotRequirements() {

}

// Conditions to run on startup
func _AliquotSetup(_ctx context.Context, _input *AliquotInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotSteps(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		panic("Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHSolution, 0)

	for i := 0; i < _input.NumberofAliquots; i++ {
		if _input.Solution.Type == "dna" {
			_input.Solution.Type = "DoNotMix"
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
		aliquot := execute.MixInto(_ctx, _input.OutPlate, aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotAnalysis(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotValidation(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {

}
func _AliquotRun(_ctx context.Context, input *AliquotInput) *AliquotOutput {
	output := &AliquotOutput{}
	_AliquotSetup(_ctx, input)
	_AliquotSteps(_ctx, input, output)
	_AliquotAnalysis(_ctx, input, output)
	_AliquotValidation(_ctx, input, output)
	return output
}

func AliquotRun(_ctx context.Context, input *AliquotInput) *AliquotSOutput {
	soutput := &AliquotSOutput{}
	output := _AliquotRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotNew() interface{} {
	return &AliquotElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotInput{},
			Out: &AliquotOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AliquotElement struct {
	inject.CheckedRunner
}

type AliquotInput struct {
	InPlate          *wtype.LHPlate
	NumberofAliquots int
	OutPlate         *wtype.LHPlate
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type AliquotOutput struct {
	Aliquots []*wtype.LHSolution
}

type AliquotSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHSolution
	}
}

func init() {
	c := Component{Name: "Aliquot", Constructor: AliquotNew}
	c.Desc.Desc = ""
	c.Desc.Params = []ParamDesc{
		{Name: "InPlate", Desc: "", Kind: "Inputs"},
		{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
		{Name: "OutPlate", Desc: "", Kind: "Inputs"},
		{Name: "Solution", Desc: "", Kind: "Inputs"},
		{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
		{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
		{Name: "Aliquots", Desc: "", Kind: "Outputs"},
	}
	addComponent(c)
}
