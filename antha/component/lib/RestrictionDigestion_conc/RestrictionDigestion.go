package RestrictionDigestion_conc

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _requirements() {}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	// workout volume of buffer to add in SI units
	BufferVol := wunit.NewVolume(float64(_input.ReactionVolume.SIValue()/float64(_input.BufferConcX)), "l")

	bufferSample := mixer.Sample(_input.Buffer, BufferVol)
	samples = append(samples, bufferSample)

	if _input.BSAvol.Mvalue != 0 {
		bsaSample := mixer.Sample(_input.BSAoptional, _input.BSAvol)
		samples = append(samples, bsaSample)
	}

	_input.DNASolution.CName = _input.DNAName

	// work out necessary volume to add
	DNAVol := wunit.NewVolume(float64((_input.DNAMassperReaction.SIValue() / _input.DNAConc.SIValue())), "l")
	text.Print("DNAVOL", DNAVol.ToString())
	dnaSample := mixer.Sample(_input.DNASolution, DNAVol)
	samples = append(samples, dnaSample)

	for k, enzyme := range _input.EnzSolutions {

		/*
			e.g.
			DesiredUinreaction = 1  // U
			StockReConcinUperml = 10000 // U/ml
			ReactionVolume = 20ul
		*/
		stockconcinUperul := _input.StockReConcinUperml[k] / 1000
		enzvoltoaddinul := _input.DesiredConcinUperml[k] / stockconcinUperul

		var enzvoltoadd wunit.Volume

		if enzvoltoaddinul < 1 {
			enzvoltoadd = wunit.NewVolume(float64(1), "ul")
		} else {
			enzvoltoadd = wunit.NewVolume(float64(enzvoltoaddinul), "ul")
		}
		enzyme.CName = _input.EnzymeNames[k]
		text.Print("adding enzyme"+_input.EnzymeNames[k], "to"+_input.DNAName)
		enzSample := mixer.Sample(enzyme, enzvoltoadd)
		enzSample.CName = _input.EnzymeNames[k]
		samples = append(samples, enzSample)
	}

	_output.Reaction = execute.MixInto(_ctx,

		_input.OutPlate, samples...)

	// incubate the reaction mixture
	execute.Incubate(_ctx,

		_output.Reaction, _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	execute.Incubate(_ctx,

		_output.Reaction, _input.InactivationTemp, _input.InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
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
	BSAoptional         *wtype.LHComponent
	BSAvol              wunit.Volume
	Buffer              *wtype.LHComponent
	BufferConcX         int
	DNAConc             wunit.Concentration
	DNAMassperReaction  wunit.Mass
	DNAName             string
	DNASolution         *wtype.LHComponent
	DesiredConcinUperml []int
	EnzSolutions        []*wtype.LHComponent
	EnzymeNames         []string
	InPlate             *wtype.LHPlate
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	OutPlate            *wtype.LHPlate
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	StockReConcinUperml []int
	Water               *wtype.LHComponent
}

type Output_ struct {
	Reaction *wtype.LHSolution
}
