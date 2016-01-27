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

//= 2 (hours)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _RecoveryRequirements() {
}

// Conditions to run on startup
func _RecoverySetup(_ctx context.Context, _input *RecoveryInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _RecoverySteps(_ctx context.Context, _input *RecoveryInput, _output *RecoveryOutput) {

	recoverymix := make([]*wtype.LHComponent, 0)

	transformedcellsComp := _input.Transformedcells

	recoverymixture := mixer.Sample(_input.Recoverymedium, _input.Recoveryvolume)

	recoverymix = append(recoverymix, transformedcellsComp, recoverymixture)
	recoverymix2 := execute.MixInto(_ctx, _input.OutPlate, "", recoverymix...)

	execute.Incubate(_ctx, recoverymix2, _input.Recoverytemp, _input.Recoverytime, true)

	_output.RecoveredCells = recoverymix2

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _RecoveryAnalysis(_ctx context.Context, _input *RecoveryInput, _output *RecoveryOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _RecoveryValidation(_ctx context.Context, _input *RecoveryInput, _output *RecoveryOutput) {
}
func _RecoveryRun(_ctx context.Context, input *RecoveryInput) *RecoveryOutput {
	output := &RecoveryOutput{}
	_RecoverySetup(_ctx, input)
	_RecoverySteps(_ctx, input, output)
	_RecoveryAnalysis(_ctx, input, output)
	_RecoveryValidation(_ctx, input, output)
	return output
}

func RecoveryRunSteps(_ctx context.Context, input *RecoveryInput) *RecoverySOutput {
	soutput := &RecoverySOutput{}
	output := _RecoveryRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func RecoveryNew() interface{} {
	return &RecoveryElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &RecoveryInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _RecoveryRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &RecoveryInput{},
			Out: &RecoveryOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type RecoveryElement struct {
	inject.CheckedRunner
}

type RecoveryInput struct {
	AgarPlate        *wtype.LHPlate
	OutPlate         *wtype.LHPlate
	Recoverymedium   *wtype.LHComponent
	Recoverytemp     wunit.Temperature
	Recoverytime     wunit.Time
	Recoveryvolume   wunit.Volume
	Transformedcells *wtype.LHComponent
}

type RecoveryOutput struct {
	RecoveredCells *wtype.LHSolution
}

type RecoverySOutput struct {
	Data struct {
	}
	Outputs struct {
		RecoveredCells *wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "Recovery",
		Constructor: RecoveryNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Transformation/Recovery.an",
			Params: []ParamDesc{
				{Name: "AgarPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Recoverymedium", Desc: "", Kind: "Inputs"},
				{Name: "Recoverytemp", Desc: "", Kind: "Parameters"},
				{Name: "Recoverytime", Desc: "= 2 (hours)\n", Kind: "Parameters"},
				{Name: "Recoveryvolume", Desc: "", Kind: "Parameters"},
				{Name: "Transformedcells", Desc: "", Kind: "Inputs"},
				{Name: "RecoveredCells", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
