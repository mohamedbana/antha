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

//Partname string
//CellVol wunit.Volume
//DNAVol wunit.Volume

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//Cells *wtype.LHComponent
//DNA *wtype.LHComponent

// Physical outputs from this protocol with types

//TransformedCells *wtype.LHComponent

func _PEG_movementRequirements() {

}

// Conditions to run on startup
func _PEG_movementSetup(_ctx context.Context, _input *PEG_movementInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PEG_movementSteps(_ctx context.Context, _input *PEG_movementInput, _output *PEG_movementOutput) {
	//cellsample:=mixer.Sample(Cells,CellVol)
	//cellsinwell:=MixTo(OutPlatetype, OutWell,1, cellsample)
	//DNASample:=mixer.Sample(DNA, DNAVol)
	//cellsplusdna:=Mix(cellsinwell, DNASample)
	PEGSample := mixer.Sample(_input.PEG, _input.PEGVol)
	_output.PEGinwell = execute.MixTo(_ctx, _input.OutPlatetype, _input.OutWell, 1, PEGSample)
	//TransformedCells=Mix(cellsplusdna,PEGSample)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PEG_movementAnalysis(_ctx context.Context, _input *PEG_movementInput, _output *PEG_movementOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _PEG_movementValidation(_ctx context.Context, _input *PEG_movementInput, _output *PEG_movementOutput) {

}
func _PEG_movementRun(_ctx context.Context, input *PEG_movementInput) *PEG_movementOutput {
	output := &PEG_movementOutput{}
	_PEG_movementSetup(_ctx, input)
	_PEG_movementSteps(_ctx, input, output)
	_PEG_movementAnalysis(_ctx, input, output)
	_PEG_movementValidation(_ctx, input, output)
	return output
}

func PEG_movementRunSteps(_ctx context.Context, input *PEG_movementInput) *PEG_movementSOutput {
	soutput := &PEG_movementSOutput{}
	output := _PEG_movementRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PEG_movementNew() interface{} {
	return &PEG_movementElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PEG_movementInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PEG_movementRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PEG_movementInput{},
			Out: &PEG_movementOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PEG_movementElement struct {
	inject.CheckedRunner
}

type PEG_movementInput struct {
	OutPlatetype string
	OutWell      string
	PEG          *wtype.LHComponent
	PEGVol       wunit.Volume
}

type PEG_movementOutput struct {
	PEGinwell *wtype.LHComponent
}

type PEG_movementSOutput struct {
	Data struct {
	}
	Outputs struct {
		PEGinwell *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "PEG_movement",
		Constructor: PEG_movementNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PEG_movement/PEG_movement.an",
			Params: []ParamDesc{
				{Name: "OutPlatetype", Desc: "", Kind: "Parameters"},
				{Name: "OutWell", Desc: "", Kind: "Parameters"},
				{Name: "PEG", Desc: "Cells *wtype.LHComponent\nDNA *wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "PEGVol", Desc: "", Kind: "Parameters"},
				{Name: "PEGinwell", Desc: "TransformedCells *wtype.LHComponent\n", Kind: "Outputs"},
			},
		},
	})
}
