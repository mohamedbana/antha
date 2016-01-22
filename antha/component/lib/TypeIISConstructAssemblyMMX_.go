package lib

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

//MMXVol					Volume

// Physical Inputs to this protocol with types

//	Water			*wtype.LHComponent

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _TypeIISConstructAssemblyMMXRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblyMMXSetup(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblyMMXSteps(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput, _output *TypeIISConstructAssemblyMMXOutput) {
	samples := make([]*wtype.LHComponent, 0)
	fmt.Println("FIRST")
	//	waterSample := mixer.SampleForTotalVolume(Water, ReactionVolume)
	mmxSample := mixer.SampleForTotalVolume(_input.MasterMix, _input.ReactionVolume)
	samples = append(samples, mmxSample)

	fmt.Println("SECOND")
	for k, part := range _input.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartNames[k], " vol ", _input.PartVols[k])
		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	_output.Reaction = execute.MixTo(_ctx, _input.OutPlate, _input.OutputLocation, samples...)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	//Incubate(Reaction, ReactionTemp, ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyMMXAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput, _output *TypeIISConstructAssemblyMMXOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyMMXValidation(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput, _output *TypeIISConstructAssemblyMMXOutput) {
}
func _TypeIISConstructAssemblyMMXRun(_ctx context.Context, input *TypeIISConstructAssemblyMMXInput) *TypeIISConstructAssemblyMMXOutput {
	output := &TypeIISConstructAssemblyMMXOutput{}
	_TypeIISConstructAssemblyMMXSetup(_ctx, input)
	_TypeIISConstructAssemblyMMXSteps(_ctx, input, output)
	_TypeIISConstructAssemblyMMXAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyMMXValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyMMXRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyMMXInput) *TypeIISConstructAssemblyMMXSOutput {
	soutput := &TypeIISConstructAssemblyMMXSOutput{}
	output := _TypeIISConstructAssemblyMMXRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyMMXNew() interface{} {
	return &TypeIISConstructAssemblyMMXElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyMMXInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyMMXRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyMMXInput{},
			Out: &TypeIISConstructAssemblyMMXOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyMMXElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyMMXInput struct {
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	MasterMix          *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputLocation     string
	OutputPlateNum     string
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
}

type TypeIISConstructAssemblyMMXOutput struct {
	Reaction *wtype.LHSolution
}

type TypeIISConstructAssemblyMMXSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "TypeIISConstructAssemblyMMX",
		Constructor: TypeIISConstructAssemblyMMXNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX/TypeIISConstructAssemblyMMX.an",
			Params: []ParamDesc{
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "\tWater\t\t\t*wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputLocation", Desc: "", Kind: "Parameters"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "", Kind: "Parameters"},
				{Name: "PartNames", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "MMXVol\t\t\t\t\tVolume\n", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
