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

func _Protoplast_movement_overlayRequirements() {

}

// Conditions to run on startup
func _Protoplast_movement_overlaySetup(_ctx context.Context, _input *Protoplast_movement_overlayInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Protoplast_movement_overlaySteps(_ctx context.Context, _input *Protoplast_movement_overlayInput, _output *Protoplast_movement_overlayOutput) {
	cellsample := mixer.Sample(_input.Cells, _input.CellVol)
	cellsinwell := execute.MixTo(_ctx, _input.OutPlatetype, _input.OutWell, 1, cellsample)

	DNASample := mixer.Sample(_input.DNA, _input.DNAVol)
	cellsplusdna := execute.Mix(_ctx, cellsinwell, DNASample)

	PEGSample := mixer.Sample(_input.PEG, _input.PEGVol)
	_output.TransformedCells = execute.Mix(_ctx, cellsplusdna, PEGSample)

	CellsAgar := mixer.Sample(_output.TransformedCells, _input.TransformedCellsVol)
	CellsAgarinWell := execute.MixTo(_ctx, _input.OutPlatetype2, _input.OutWell2, 1, CellsAgar)

	//CellsAgarinWellFinal=CellsAgarinWell

	plateout := mixer.Sample(CellsAgarinWell, _input.Plateoutvolume)
	platedculture := execute.MixTo(_ctx, _input.AgarPlate, _input.OutWell3, 1, plateout)

	_output.Platedculture = platedculture

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Protoplast_movement_overlayAnalysis(_ctx context.Context, _input *Protoplast_movement_overlayInput, _output *Protoplast_movement_overlayOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Protoplast_movement_overlayValidation(_ctx context.Context, _input *Protoplast_movement_overlayInput, _output *Protoplast_movement_overlayOutput) {

}
func _Protoplast_movement_overlayRun(_ctx context.Context, input *Protoplast_movement_overlayInput) *Protoplast_movement_overlayOutput {
	output := &Protoplast_movement_overlayOutput{}
	_Protoplast_movement_overlaySetup(_ctx, input)
	_Protoplast_movement_overlaySteps(_ctx, input, output)
	_Protoplast_movement_overlayAnalysis(_ctx, input, output)
	_Protoplast_movement_overlayValidation(_ctx, input, output)
	return output
}

func Protoplast_movement_overlayRunSteps(_ctx context.Context, input *Protoplast_movement_overlayInput) *Protoplast_movement_overlaySOutput {
	soutput := &Protoplast_movement_overlaySOutput{}
	output := _Protoplast_movement_overlayRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Protoplast_movement_overlayNew() interface{} {
	return &Protoplast_movement_overlayElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Protoplast_movement_overlayInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Protoplast_movement_overlayRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Protoplast_movement_overlayInput{},
			Out: &Protoplast_movement_overlayOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Protoplast_movement_overlayElement struct {
	inject.CheckedRunner
}

type Protoplast_movement_overlayInput struct {
	AgarPlate           string
	CellVol             wunit.Volume
	Cells               *wtype.LHComponent
	DNA                 *wtype.LHComponent
	DNAVol              wunit.Volume
	OutPlatetype        string
	OutPlatetype2       string
	OutWell             string
	OutWell2            string
	OutWell3            string
	PEG                 *wtype.LHComponent
	PEGVol              wunit.Volume
	Partname            string
	Plateoutvolume      wunit.Volume
	TransformedCellsVol wunit.Volume
}

type Protoplast_movement_overlayOutput struct {
	CellsAgarinWellFinal *wtype.LHComponent
	Platedculture        *wtype.LHComponent
	TransformedCells     *wtype.LHComponent
}

type Protoplast_movement_overlaySOutput struct {
	Data struct {
	}
	Outputs struct {
		CellsAgarinWellFinal *wtype.LHComponent
		Platedculture        *wtype.LHComponent
		TransformedCells     *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Protoplast_movement_overlay",
		Constructor: Protoplast_movement_overlayNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Protoplast_movement_overlay/Protoplast_movement_overlay.an",
			Params: []ParamDesc{
				{Name: "AgarPlate", Desc: "", Kind: "Parameters"},
				{Name: "CellVol", Desc: "", Kind: "Parameters"},
				{Name: "Cells", Desc: "", Kind: "Inputs"},
				{Name: "DNA", Desc: "", Kind: "Inputs"},
				{Name: "DNAVol", Desc: "", Kind: "Parameters"},
				{Name: "OutPlatetype", Desc: "", Kind: "Parameters"},
				{Name: "OutPlatetype2", Desc: "", Kind: "Parameters"},
				{Name: "OutWell", Desc: "", Kind: "Parameters"},
				{Name: "OutWell2", Desc: "", Kind: "Parameters"},
				{Name: "OutWell3", Desc: "", Kind: "Parameters"},
				{Name: "PEG", Desc: "", Kind: "Inputs"},
				{Name: "PEGVol", Desc: "", Kind: "Parameters"},
				{Name: "Partname", Desc: "", Kind: "Parameters"},
				{Name: "Plateoutvolume", Desc: "", Kind: "Parameters"},
				{Name: "TransformedCellsVol", Desc: "", Kind: "Parameters"},
				{Name: "CellsAgarinWellFinal", Desc: "", Kind: "Outputs"},
				{Name: "Platedculture", Desc: "", Kind: "Outputs"},
				{Name: "TransformedCells", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
