package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

//"strconv"

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ScreenLHPoliciesRequirements() {
}

// Conditions to run on startup
func _ScreenLHPoliciesSetup(_ctx context.Context, _input *ScreenLHPoliciesInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ScreenLHPoliciesSteps(_ctx context.Context, _input *ScreenLHPoliciesInput, _output *ScreenLHPoliciesOutput) {

	reactions := make([]*wtype.LHSolution, 0)

	//policies, names := liquidhandling.PolicyMaker(liquidhandling.Allpairs, "DOE_run",false)

	//intfactors := []string{"Pre_MIX","POST_MIX"}
	policies, names, err := liquidhandling.PolicyMakerfromDesign("LHPolicydesign.xlsx", "DOE_run")
	if err != nil {
		panic(err)
	}

	for k := 0; k < len(_input.TestSols); k++ {
		for j := 0; j < _input.NumberofReplicates; j++ {
			for i := 0; i < len(policies); i++ {

				eachreaction := make([]*wtype.LHComponent, 0)

				_input.Diluent.Type = names[i]
				fmt.Println(_input.Diluent.Type)

				bufferSample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
				eachreaction = append(eachreaction, bufferSample)
				testSample := mixer.Sample(_input.TestSols[k], _input.TestSolVolume)
				eachreaction = append(eachreaction, testSample)
				reaction := execute.MixInto(_ctx, _input.OutPlate, eachreaction...)
				reactions = append(reactions, reaction)

			}
		}
	}
	_output.Reactions = reactions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ScreenLHPoliciesAnalysis(_ctx context.Context, _input *ScreenLHPoliciesInput, _output *ScreenLHPoliciesOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ScreenLHPoliciesValidation(_ctx context.Context, _input *ScreenLHPoliciesInput, _output *ScreenLHPoliciesOutput) {
}
func _ScreenLHPoliciesRun(_ctx context.Context, input *ScreenLHPoliciesInput) *ScreenLHPoliciesOutput {
	output := &ScreenLHPoliciesOutput{}
	_ScreenLHPoliciesSetup(_ctx, input)
	_ScreenLHPoliciesSteps(_ctx, input, output)
	_ScreenLHPoliciesAnalysis(_ctx, input, output)
	_ScreenLHPoliciesValidation(_ctx, input, output)
	return output
}

func ScreenLHPoliciesRunSteps(_ctx context.Context, input *ScreenLHPoliciesInput) *ScreenLHPoliciesSOutput {
	soutput := &ScreenLHPoliciesSOutput{}
	output := _ScreenLHPoliciesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ScreenLHPoliciesNew() interface{} {
	return &ScreenLHPoliciesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ScreenLHPoliciesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ScreenLHPoliciesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ScreenLHPoliciesInput{},
			Out: &ScreenLHPoliciesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ScreenLHPoliciesElement struct {
	inject.CheckedRunner
}

type ScreenLHPoliciesInput struct {
	Diluent            *wtype.LHComponent
	NumberofReplicates int
	OutPlate           *wtype.LHPlate
	TestSolVolume      wunit.Volume
	TestSols           []*wtype.LHComponent
	TotalVolume        wunit.Volume
}

type ScreenLHPoliciesOutput struct {
	Reactions []*wtype.LHSolution
	Status    string
}

type ScreenLHPoliciesSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Reactions []*wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "ScreenLHPolicies",
		Constructor: ScreenLHPoliciesNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/FindbestLHPolicy/ScreenLHPolicies.an",
			Params: []ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "NumberofReplicates", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "TestSolVolume", Desc: "", Kind: "Parameters"},
				{Name: "TestSols", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	})
}
