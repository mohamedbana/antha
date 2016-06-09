// example protocol demonstrating the use of the SampleAll function
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

func _SampleAllRequirements() {

}

// Conditions to run on startup
func _SampleAllSetup(_ctx context.Context, _input *SampleAllInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleAllSteps(_ctx context.Context, _input *SampleAllInput, _output *SampleAllOutput) {
	// the SampleAll function samples the entire contents of the LHComponent
	if _input.Sampleall {
		_output.Sample = mixer.SampleAll(_input.Solution)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleAllAnalysis(_ctx context.Context, _input *SampleAllInput, _output *SampleAllOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleAllValidation(_ctx context.Context, _input *SampleAllInput, _output *SampleAllOutput) {

}
func _SampleAllRun(_ctx context.Context, input *SampleAllInput) *SampleAllOutput {
	output := &SampleAllOutput{}
	_SampleAllSetup(_ctx, input)
	_SampleAllSteps(_ctx, input, output)
	_SampleAllAnalysis(_ctx, input, output)
	_SampleAllValidation(_ctx, input, output)
	return output
}

func SampleAllRunSteps(_ctx context.Context, input *SampleAllInput) *SampleAllSOutput {
	soutput := &SampleAllSOutput{}
	output := _SampleAllRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleAllNew() interface{} {
	return &SampleAllElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleAllInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleAllRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleAllInput{},
			Out: &SampleAllOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SampleAllElement struct {
	inject.CheckedRunner
}

type SampleAllInput struct {
	Sampleall bool
	Solution  *wtype.LHComponent
}

type SampleAllOutput struct {
	Sample *wtype.LHComponent
}

type SampleAllSOutput struct {
	Data struct {
	}
	Outputs struct {
		Sample *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "SampleAll",
		Constructor: SampleAllNew,
		Desc: ComponentDesc{
			Desc: "example protocol demonstrating the use of the SampleAll function\n",
			Path: "antha/component/an/AnthaAcademy/Lesson1_Sample/SampleAll.an",
			Params: []ParamDesc{
				{Name: "Sampleall", Desc: "", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "Sample", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
