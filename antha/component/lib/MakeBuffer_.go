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

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _MakeBufferRequirements() {

}

// Conditions to run on startup
func _MakeBufferSetup(_ctx context.Context, _input *MakeBufferInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeBufferSteps(_ctx context.Context, _input *MakeBufferInput, _output *MakeBufferOutput) {
	//Bufferstockvolume := wunit.NewVolume((FinalVolume.SIValue() * FinalConcentration.SIValue()/Bufferstockconc.SIValue()),"l")

	_output.Buffer = execute.MixInto(_ctx, _input.OutPlate, "",
		mixer.Sample(_input.Bufferstock, _input.Bufferstockvolume),
		mixer.Sample(_input.Diluent, _input.Diluentvolume))

	_output.Status = fmt.Sprintln("Buffer stock volume = ", _input.Bufferstockvolume.ToString(), "of", _input.Bufferstock.CName,
		"was added to ", _input.Diluentvolume.ToString(), "of", _input.Diluent.CName,
		"to make ", _input.FinalVolume.ToString(), "of", _input.Buffername,
		"Buffer stock conc =", _input.Bufferstockconc.ToString())

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeBufferAnalysis(_ctx context.Context, _input *MakeBufferInput, _output *MakeBufferOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeBufferValidation(_ctx context.Context, _input *MakeBufferInput, _output *MakeBufferOutput) {
}
func _MakeBufferRun(_ctx context.Context, input *MakeBufferInput) *MakeBufferOutput {
	output := &MakeBufferOutput{}
	_MakeBufferSetup(_ctx, input)
	_MakeBufferSteps(_ctx, input, output)
	_MakeBufferAnalysis(_ctx, input, output)
	_MakeBufferValidation(_ctx, input, output)
	return output
}

func MakeBufferRunSteps(_ctx context.Context, input *MakeBufferInput) *MakeBufferSOutput {
	soutput := &MakeBufferSOutput{}
	output := _MakeBufferRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeBufferNew() interface{} {
	return &MakeBufferElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeBufferInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeBufferRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeBufferInput{},
			Out: &MakeBufferOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakeBufferElement struct {
	inject.CheckedRunner
}

type MakeBufferInput struct {
	Buffername         string
	Bufferstock        *wtype.LHComponent
	Bufferstockconc    wunit.Concentration
	Bufferstockvolume  wunit.Volume
	Diluent            *wtype.LHComponent
	Diluentname        string
	Diluentvolume      wunit.Volume
	FinalConcentration wunit.Concentration
	FinalVolume        wunit.Volume
	InPlate            *wtype.LHPlate
	OutPlate           *wtype.LHPlate
}

type MakeBufferOutput struct {
	Buffer *wtype.LHSolution
	Status string
}

type MakeBufferSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Buffer *wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "MakeBuffer",
		Constructor: MakeBufferNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/MakeBuffer/Makebuffer.an",
			Params: []ParamDesc{
				{Name: "Buffername", Desc: "", Kind: "Parameters"},
				{Name: "Bufferstock", Desc: "", Kind: "Inputs"},
				{Name: "Bufferstockconc", Desc: "", Kind: "Parameters"},
				{Name: "Bufferstockvolume", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Diluentname", Desc: "", Kind: "Parameters"},
				{Name: "Diluentvolume", Desc: "", Kind: "Parameters"},
				{Name: "FinalConcentration", Desc: "", Kind: "Parameters"},
				{Name: "FinalVolume", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Buffer", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	})
}

/*
type Mole struct {
	number float64
}*/
