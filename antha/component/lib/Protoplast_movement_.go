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

func _Protoplast_movementRequirements() {

}

// Conditions to run on startup
func _Protoplast_movementSetup(_ctx context.Context, _input *Protoplast_movementInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Protoplast_movementSteps(_ctx context.Context, _input *Protoplast_movementInput, _output *Protoplast_movementOutput) {
	cellsample := mixer.Sample(_input.Cells, _input.CellVol)
	cellsinwell := execute.MixTo(_ctx, _input.OutPlatetype, _input.OutWell, 1, cellsample)
	DNASample := mixer.Sample(_input.DNA, _input.DNAVol)
	cellsplusdna := execute.Mix(_ctx, cellsinwell, DNASample)
	PEGSample := mixer.Sample(_input.PEG, _input.PEGVol)
	_output.TransformedCells = execute.Mix(_ctx, cellsplusdna, PEGSample)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Protoplast_movementAnalysis(_ctx context.Context, _input *Protoplast_movementInput, _output *Protoplast_movementOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Protoplast_movementValidation(_ctx context.Context, _input *Protoplast_movementInput, _output *Protoplast_movementOutput) {

}
func _Protoplast_movementRun(_ctx context.Context, input *Protoplast_movementInput) *Protoplast_movementOutput {
	output := &Protoplast_movementOutput{}
	_Protoplast_movementSetup(_ctx, input)
	_Protoplast_movementSteps(_ctx, input, output)
	_Protoplast_movementAnalysis(_ctx, input, output)
	_Protoplast_movementValidation(_ctx, input, output)
	return output
}

func Protoplast_movementRunSteps(_ctx context.Context, input *Protoplast_movementInput) *Protoplast_movementSOutput {
	soutput := &Protoplast_movementSOutput{}
	output := _Protoplast_movementRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Protoplast_movementNew() interface{} {
	return &Protoplast_movementElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Protoplast_movementInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Protoplast_movementRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Protoplast_movementInput{},
			Out: &Protoplast_movementOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Protoplast_movementElement struct {
	inject.CheckedRunner
}

type Protoplast_movementInput struct {
	CellVol      wunit.Volume
	Cells        *wtype.LHComponent
	DNA          *wtype.LHComponent
	DNAVol       wunit.Volume
	OutPlatetype string
	OutWell      string
	PEG          *wtype.LHComponent
	PEGVol       wunit.Volume
	Partname     string
}

type Protoplast_movementOutput struct {
	TransformedCells *wtype.LHComponent
}

type Protoplast_movementSOutput struct {
	Data struct {
	}
	Outputs struct {
		TransformedCells *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Protoplast_movement",
		Constructor: Protoplast_movementNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Protoplast_movement/Protoplast_movement.an",
			Params: []ParamDesc{
				{Name: "CellVol", Desc: "", Kind: "Parameters"},
				{Name: "Cells", Desc: "", Kind: "Inputs"},
				{Name: "DNA", Desc: "", Kind: "Inputs"},
				{Name: "DNAVol", Desc: "", Kind: "Parameters"},
				{Name: "OutPlatetype", Desc: "", Kind: "Parameters"},
				{Name: "OutWell", Desc: "", Kind: "Parameters"},
				{Name: "PEG", Desc: "", Kind: "Inputs"},
				{Name: "PEGVol", Desc: "", Kind: "Parameters"},
				{Name: "Partname", Desc: "", Kind: "Parameters"},
				{Name: "TransformedCells", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
