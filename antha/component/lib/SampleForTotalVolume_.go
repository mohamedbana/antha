// example protocol demonstrating the use of the SampleForTotalVolume function
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

// e.g. 2ul
// e.g. 20ul

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _SampleForTotalVolumeRequirements() {

}

// Conditions to run on startup
func _SampleForTotalVolumeSetup(_ctx context.Context, _input *SampleForTotalVolumeInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleForTotalVolumeSteps(_ctx context.Context, _input *SampleForTotalVolumeInput, _output *SampleForTotalVolumeOutput) {

	// make empty slice of components ready to sequentially add all samples to
	allsamples := make([]*wtype.LHComponent, 0)

	// SampleForTotalVolume will "top up" solution to the TotalVolume with Diluent.
	// in this case it will still add diluent first but calculates the volume to add by substracting the volumes of subsequent components
	diluentsample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume) // i.e. if TotalVolume == 20ul and SolutionVolume == 2ul then 18ul of Diluent will be sampled here

	// append will add the diliuent sample to the allsamples slice
	allsamples = append(allsamples, diluentsample)

	solutionsample := mixer.Sample(_input.Solution, _input.SolutionVolume)

	allsamples = append(allsamples, solutionsample)

	// finally use Mix to combine samples into a new component
	_output.DilutedSample = execute.Mix(_ctx, allsamples...)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleForTotalVolumeAnalysis(_ctx context.Context, _input *SampleForTotalVolumeInput, _output *SampleForTotalVolumeOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleForTotalVolumeValidation(_ctx context.Context, _input *SampleForTotalVolumeInput, _output *SampleForTotalVolumeOutput) {

}
func _SampleForTotalVolumeRun(_ctx context.Context, input *SampleForTotalVolumeInput) *SampleForTotalVolumeOutput {
	output := &SampleForTotalVolumeOutput{}
	_SampleForTotalVolumeSetup(_ctx, input)
	_SampleForTotalVolumeSteps(_ctx, input, output)
	_SampleForTotalVolumeAnalysis(_ctx, input, output)
	_SampleForTotalVolumeValidation(_ctx, input, output)
	return output
}

func SampleForTotalVolumeRunSteps(_ctx context.Context, input *SampleForTotalVolumeInput) *SampleForTotalVolumeSOutput {
	soutput := &SampleForTotalVolumeSOutput{}
	output := _SampleForTotalVolumeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleForTotalVolumeNew() interface{} {
	return &SampleForTotalVolumeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleForTotalVolumeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleForTotalVolumeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleForTotalVolumeInput{},
			Out: &SampleForTotalVolumeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SampleForTotalVolumeElement struct {
	inject.CheckedRunner
}

type SampleForTotalVolumeInput struct {
	Diluent        *wtype.LHComponent
	Solution       *wtype.LHComponent
	SolutionVolume wunit.Volume
	TotalVolume    wunit.Volume
}

type SampleForTotalVolumeOutput struct {
	DilutedSample *wtype.LHComponent
}

type SampleForTotalVolumeSOutput struct {
	Data struct {
	}
	Outputs struct {
		DilutedSample *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "SampleForTotalVolume",
		Constructor: SampleForTotalVolumeNew,
		Desc: ComponentDesc{
			Desc: "example protocol demonstrating the use of the SampleForTotalVolume function\n",
			Path: "antha/component/an/AnthaAcademy/Lesson1_Sample/SampleForTotalVolume.an",
			Params: []ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "e.g. 2ul\n", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "e.g. 20ul\n", Kind: "Parameters"},
				{Name: "DilutedSample", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
