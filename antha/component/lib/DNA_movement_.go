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

//Cells *wtype.LHComponent

//PEG *wtype.LHComponent

// Physical outputs from this protocol with types

//TransformedCells *wtype.LHComponent

func _DNA_movementRequirements() {

}

// Conditions to run on startup
func _DNA_movementSetup(_ctx context.Context, _input *DNA_movementInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _DNA_movementSteps(_ctx context.Context, _input *DNA_movementInput, _output *DNA_movementOutput) {
	//cellsample:=mixer.Sample(Cells,CellVol)
	//cellsinwell:=MixTo(OutPlatetype, OutWell,1, cellsample)
	//DNASample:=mixer.Sample(DNA, DNAVol)
	//cellsplusdna:=Mix(cellsinwell, DNASample)
	DNASample := mixer.Sample(_input.DNA, _input.DNAVol)
	_output.DNAinwell = execute.MixTo(_ctx, _input.OutPlatetype, _input.OutWell, 1, DNASample)
	//TransformedCells=Mix(cellsplusdna,PEGSample)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _DNA_movementAnalysis(_ctx context.Context, _input *DNA_movementInput, _output *DNA_movementOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _DNA_movementValidation(_ctx context.Context, _input *DNA_movementInput, _output *DNA_movementOutput) {

}
func _DNA_movementRun(_ctx context.Context, input *DNA_movementInput) *DNA_movementOutput {
	output := &DNA_movementOutput{}
	_DNA_movementSetup(_ctx, input)
	_DNA_movementSteps(_ctx, input, output)
	_DNA_movementAnalysis(_ctx, input, output)
	_DNA_movementValidation(_ctx, input, output)
	return output
}

func DNA_movementRunSteps(_ctx context.Context, input *DNA_movementInput) *DNA_movementSOutput {
	soutput := &DNA_movementSOutput{}
	output := _DNA_movementRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DNA_movementNew() interface{} {
	return &DNA_movementElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DNA_movementInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DNA_movementRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DNA_movementInput{},
			Out: &DNA_movementOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DNA_movementElement struct {
	inject.CheckedRunner
}

type DNA_movementInput struct {
	DNA          *wtype.LHComponent
	DNAVol       wunit.Volume
	OutPlatetype string
	OutWell      string
	Partname     string
}

type DNA_movementOutput struct {
	DNAinwell *wtype.LHComponent
}

type DNA_movementSOutput struct {
	Data struct {
	}
	Outputs struct {
		DNAinwell *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "DNA_movement",
		Constructor: DNA_movementNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/DNA_movement/DNA_movement.an",
			Params: []ParamDesc{
				{Name: "DNA", Desc: "Cells *wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "DNAVol", Desc: "", Kind: "Parameters"},
				{Name: "OutPlatetype", Desc: "", Kind: "Parameters"},
				{Name: "OutWell", Desc: "", Kind: "Parameters"},
				{Name: "Partname", Desc: "", Kind: "Parameters"},
				{Name: "DNAinwell", Desc: "TransformedCells *wtype.LHComponent\n", Kind: "Outputs"},
			},
		},
	})
}
