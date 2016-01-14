package SDSprep

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//Input parameters for this protocol. Single instance of an SDS-PAGE sample preperation step.
//Mix 10ul of 4x stock buffer with 30ul of proteinX sample to create 40ul sample for loading.

//ProteinX
//30uL

//SDSBuffer
//10ul
//100g/L

//25g/L
//40uL

//5min
//95oC

//Biologicals

//Purified protein or cell lysate...

//Chemicals

//Consumables

//Contains protein and buffer
//Final plate with mixed components

//Biologicals

func _setup(_ctx context.Context, _input *Input_) {
}

func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	//Method 1. Mix two things. DOES NOT WORK as recognises protein to be 1 single entity and won't handle as seperate components. ie end result is 5 things created all
	//from the same well. Check typeIIs workflow for hints.
	//
	//	Step1a
	//	LoadSample = MixInto(OutPlate,
	//	mixer.Sample(Protein, SampleVolume),
	//	mixer.Sample(Buffer, BufferVolume))
	//Try something else. Outputs are an array taking in a single (not array) of protein and buffer. Do this 12 times.

	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.Sample(_input.Buffer, _input.BufferVolume)
	bufferSample.CName = _input.BufferName
	samples = append(samples, bufferSample)

	proteinSample := mixer.Sample(_input.Protein, _input.SampleVolume)
	proteinSample.CName = _input.SampleName
	samples = append(samples, proteinSample)
	fmt.Println("This is a sample list ", samples)
	_output.LoadSample = execute.MixInto(_ctx,

		_input.OutPlate, samples...)

	//Methods 2.Make a sample of two things creating a list
	//	Step 1b

	//	sample	    := make([]wtype.LHComponent, 0)

	//	bufferPart  := mixer.Sample(Buffer, BufferVolume)
	//	sample	     = append([]samples, bufferSample)

	//	proteinPart := mixer.Sample(Protein, SampleVolume)
	//	sample      = append([]samples, proteinSample)

	//	LoadSample   = MixInto(OutPlate, sample...)

	//Denature the load mixture at specified temperature and time ie 95oC for 5min
	//	Step2
	execute.Incubate(_ctx,

		_output.LoadSample, _input.DenatureTemp, _input.DenatureTime, false)

	//Load the water in EPAGE gel wells
	//	Step3

	//	var water water volume
	//	waterLoad := mixer.Sample(Water, WaterLoadVolume)
	//
	//Load the LoadSample into EPAGE gel
	//
	//	Loader = MixInto(EPAGE48, LoadSample)
	//
	//
	//

	//	Status = fmtSprintln(BufferVolume.ToString() "uL of", BufferName,"mixed with", SampleVolume.ToString(), "uL of", SampleName, "Total load sample available is", ReactionVolume.ToString())
}

func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
}

func _validation(_ctx context.Context, _input *Input_, _output *Output_) {
}

func _run(_ctx context.Context, value inject.Value) (inject.Value, error) {
	input := &Input_{}
	output := &Output_{}
	if err := inject.Assign(value, input); err != nil {
		return nil, err
	}
	_setup(_ctx, input)
	_steps(_ctx, input, output)
	_analysis(_ctx, input, output)
	_validation(_ctx, input, output)
	return inject.MakeValue(output), nil
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

func New() interface{} {
	return &Element_{
		inject.CheckedRunner{
			RunFunc: _run,
			In:      &Input_{},
			Out:     &Output_{},
		},
	}
}

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	Buffer             *wtype.LHComponent
	BufferName         string
	BufferStockConc    wunit.Concentration
	BufferVolume       wunit.Volume
	DenatureTemp       wunit.Temperature
	DenatureTime       wunit.Time
	FinalConcentration wunit.Concentration
	InPlate            *wtype.LHPlate
	OutPlate           *wtype.LHPlate
	Protein            *wtype.LHComponent
	ReactionVolume     wunit.Volume
	SampleName         string
	SampleVolume       wunit.Volume
}

type Output_ struct {
	LoadSample *wtype.LHSolution
	Status     string
}
