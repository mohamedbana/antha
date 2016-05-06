package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/buffers"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _DiluteBufferRequirements() {

}

// Conditions to run on startup
func _DiluteBufferSetup(_ctx context.Context, _input *DiluteBufferInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _DiluteBufferSteps(_ctx context.Context, _input *DiluteBufferInput, _output *DiluteBufferOutput) {
	//Bufferstockvolume := wunit.NewVolume((FinalVolume.SIValue() * FinalConcentration.SIValue()/Bufferstockconc.SIValue()),"l")

	_output.FinalConcentration = buffers.Dilute(_input.Buffername, _input.Bufferstockconc, _input.BufferVolumeAdded, _input.Diluentname, _input.DiluentVolume)

	_output.Buffer = execute.MixInto(_ctx, _input.OutPlate, "",
		mixer.Sample(_input.Bufferstock, _input.BufferVolumeAdded),
		mixer.Sample(_input.Diluent, _input.DiluentVolume))

	_output.Status = fmt.Sprintln("Buffer stock volume = ", _input.BufferVolumeAdded.ToString(), "of", _input.Bufferstock.CName,
		"was added to ", _input.DiluentVolume.ToString(), "of", _input.Diluent.CName,
		"to make ", _input.BufferVolumeAdded.SIValue()+_input.DiluentVolume.SIValue(), "L", "of", _input.Buffername,
		"Buffer stock conc =", _output.FinalConcentration.ToString())

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _DiluteBufferAnalysis(_ctx context.Context, _input *DiluteBufferInput, _output *DiluteBufferOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _DiluteBufferValidation(_ctx context.Context, _input *DiluteBufferInput, _output *DiluteBufferOutput) {
}
func _DiluteBufferRun(_ctx context.Context, input *DiluteBufferInput) *DiluteBufferOutput {
	output := &DiluteBufferOutput{}
	_DiluteBufferSetup(_ctx, input)
	_DiluteBufferSteps(_ctx, input, output)
	_DiluteBufferAnalysis(_ctx, input, output)
	_DiluteBufferValidation(_ctx, input, output)
	return output
}

func DiluteBufferRunSteps(_ctx context.Context, input *DiluteBufferInput) *DiluteBufferSOutput {
	soutput := &DiluteBufferSOutput{}
	output := _DiluteBufferRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DiluteBufferNew() interface{} {
	return &DiluteBufferElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DiluteBufferInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DiluteBufferRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DiluteBufferInput{},
			Out: &DiluteBufferOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DiluteBufferElement struct {
	inject.CheckedRunner
}

type DiluteBufferInput struct {
	BufferVolumeAdded wunit.Volume
	Buffername        string
	Bufferstock       *wtype.LHComponent
	Bufferstockconc   wunit.Concentration
	Diluent           *wtype.LHComponent
	DiluentVolume     wunit.Volume
	Diluentname       string
	InPlate           *wtype.LHPlate
	OutPlate          *wtype.LHPlate
}

type DiluteBufferOutput struct {
	Buffer             *wtype.LHComponent
	DiluentVolume      wunit.Volume
	FinalConcentration wunit.Concentration
	Status             string
}

type DiluteBufferSOutput struct {
	Data struct {
		DiluentVolume      wunit.Volume
		FinalConcentration wunit.Concentration
		Status             string
	}
	Outputs struct {
		Buffer *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "DiluteBuffer",
		Constructor: DiluteBufferNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/MakeBuffer/DiluteBuffer.an",
			Params: []ParamDesc{
				{Name: "BufferVolumeAdded", Desc: "", Kind: "Parameters"},
				{Name: "Buffername", Desc: "", Kind: "Parameters"},
				{Name: "Bufferstock", Desc: "", Kind: "Inputs"},
				{Name: "Bufferstockconc", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DiluentVolume", Desc: "", Kind: "Parameters"},
				{Name: "Diluentname", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Buffer", Desc: "", Kind: "Outputs"},
				{Name: "DiluentVolume", Desc: "", Kind: "Data"},
				{Name: "FinalConcentration", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	})
}

/*
type Mole struct {
	number float64
}*/
