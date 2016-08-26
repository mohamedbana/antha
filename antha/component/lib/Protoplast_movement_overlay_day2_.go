package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

/*CellVol wunit.Volume
DNAVol wunit.Volume
PEGVol wunit.Volume
TransformedCellsVol wunit.Volume*/

/*OutPlatetype string
OutPlatetype2 string*/

/*OutWell2 string
OutWell3 string
Partname string*/

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//OverlayAgar *wtype.LHComponent
/*DNA *wtype.LHComponent
PEG *wtype.LHComponent*/

// Physical outputs from this protocol with types

/*TransformedCells *wtype.LHComponent
CellsAgarinWellFinal *wtype.LHComponent*/

func _Protoplast_movement_overlay_day2Requirements() {

}

// Conditions to run on startup
func _Protoplast_movement_overlay_day2Setup(_ctx context.Context, _input *Protoplast_movement_overlay_day2Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Protoplast_movement_overlay_day2Steps(_ctx context.Context, _input *Protoplast_movement_overlay_day2Input, _output *Protoplast_movement_overlay_day2Output) {
	/*cellsample:=mixer.Sample(Cells,CellVol)
	cellsinwell:=MixTo(OutPlatetype, OutWell,1, cellsample)

	DNASample:=mixer.Sample(DNA, DNAVol)
	cellsplusdna:=Mix(cellsinwell, DNASample)

	PEGSample:=mixer.Sample(PEG, PEGVol)
	TransformedCells=Mix(cellsplusdna, PEGSample)*/

	/*CellsAgar:=mixer.Sample(TransformedCells,TransformedCellsVol)
	CellsAgarinWell:=MixTo(OutPlatetype2, OutWell2, 1, CellsAgar)*/

	//CellsAgarinWellFinal=CellsAgarinWell

	overlayAgar := factory.GetComponentByType("protoplasts")

	overlayAgar.CName = _input.OverlayAgar

	plateout := mixer.Sample(overlayAgar, _input.Plateoutvolume)
	platedculture := execute.MixTo(_ctx, _input.AgarPlate, _input.OutWell, _input.Platenumber, plateout)

	_output.Platedculture = platedculture

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Protoplast_movement_overlay_day2Analysis(_ctx context.Context, _input *Protoplast_movement_overlay_day2Input, _output *Protoplast_movement_overlay_day2Output) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Protoplast_movement_overlay_day2Validation(_ctx context.Context, _input *Protoplast_movement_overlay_day2Input, _output *Protoplast_movement_overlay_day2Output) {

}
func _Protoplast_movement_overlay_day2Run(_ctx context.Context, input *Protoplast_movement_overlay_day2Input) *Protoplast_movement_overlay_day2Output {
	output := &Protoplast_movement_overlay_day2Output{}
	_Protoplast_movement_overlay_day2Setup(_ctx, input)
	_Protoplast_movement_overlay_day2Steps(_ctx, input, output)
	_Protoplast_movement_overlay_day2Analysis(_ctx, input, output)
	_Protoplast_movement_overlay_day2Validation(_ctx, input, output)
	return output
}

func Protoplast_movement_overlay_day2RunSteps(_ctx context.Context, input *Protoplast_movement_overlay_day2Input) *Protoplast_movement_overlay_day2SOutput {
	soutput := &Protoplast_movement_overlay_day2SOutput{}
	output := _Protoplast_movement_overlay_day2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Protoplast_movement_overlay_day2New() interface{} {
	return &Protoplast_movement_overlay_day2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Protoplast_movement_overlay_day2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Protoplast_movement_overlay_day2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Protoplast_movement_overlay_day2Input{},
			Out: &Protoplast_movement_overlay_day2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Protoplast_movement_overlay_day2Element struct {
	inject.CheckedRunner
}

type Protoplast_movement_overlay_day2Input struct {
	AgarPlate      string
	OutWell        string
	OverlayAgar    string
	Platenumber    int
	Plateoutvolume wunit.Volume
}

type Protoplast_movement_overlay_day2Output struct {
	Platedculture *wtype.LHComponent
}

type Protoplast_movement_overlay_day2SOutput struct {
	Data struct {
	}
	Outputs struct {
		Platedculture *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Protoplast_movement_overlay_day2",
		Constructor: Protoplast_movement_overlay_day2New,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Protoplast_movement_overlay/Protoplast_movement_overlay2.an",
			Params: []ParamDesc{
				{Name: "AgarPlate", Desc: "OutPlatetype string\n\tOutPlatetype2 string\n", Kind: "Parameters"},
				{Name: "OutWell", Desc: "", Kind: "Parameters"},
				{Name: "OverlayAgar", Desc: "", Kind: "Parameters"},
				{Name: "Platenumber", Desc: "", Kind: "Parameters"},
				{Name: "Plateoutvolume", Desc: "CellVol wunit.Volume\n\tDNAVol wunit.Volume\n\tPEGVol wunit.Volume\n\tTransformedCellsVol wunit.Volume\n", Kind: "Parameters"},
				{Name: "Platedculture", Desc: "TransformedCells *wtype.LHComponent\n\tCellsAgarinWellFinal *wtype.LHComponent\n", Kind: "Outputs"},
			},
		},
	})
}
