package RestrictionDigestion

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

//	StockReConcinUperml 		[]int
//	DesiredConcinUperml	 		[]int

//OutputReactionName			string

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

	bufferSample := mixer.Sample(_input.Buffer, _input.BufferVol)
	samples = append(samples, bufferSample)

	if _input.BSAvol.Mvalue != 0 {
		bsaSample := mixer.Sample(_input.BSAoptional, _input.BSAvol)
		samples = append(samples, bsaSample)
	}

	// change to fixing concentration(or mass) of dna per reaction
	_input.DNASolution.CName = _input.DNAName
	dnaSample := mixer.Sample(_input.DNASolution, _input.DNAVol)
	samples = append(samples, dnaSample)

	for k, enzyme := range _input.EnzSolutions {

		// work out volume to add in L

		// e.g. 1 U / (10000 * 1000) * 0.000002
		//volinL := DesiredUinreaction/(StockReConcinUperml*1000) * ReactionVolume.SIValue()
		//volumetoadd := wunit.NewVolume(volinL,"L")
		enzyme.CName = _input.EnzymeNames[k]
		text.Print("adding enzyme"+_input.EnzymeNames[k], "to"+_input.DNAName)
		enzSample := mixer.Sample(enzyme, _input.EnzVolumestoadd[k])
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
	BSAoptional      *wtype.LHComponent
	BSAvol           wunit.Volume
	Buffer           *wtype.LHComponent
	BufferVol        wunit.Volume
	DNAName          string
	DNASolution      *wtype.LHComponent
	DNAVol           wunit.Volume
	EnzSolutions     []*wtype.LHComponent
	EnzVolumestoadd  []wunit.Volume
	EnzymeNames      []string
	InPlate          *wtype.LHPlate
	InactivationTemp wunit.Temperature
	InactivationTime wunit.Time
	OutPlate         *wtype.LHPlate
	ReactionTemp     wunit.Temperature
	ReactionTime     wunit.Time
	ReactionVolume   wunit.Volume
	Water            *wtype.LHComponent
}

type Output_ struct {
	Reaction *wtype.LHSolution
}
