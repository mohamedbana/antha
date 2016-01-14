package PCR

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

/*type Polymerase struct {
	wtype.LHComponent
	Rate_BPpers float64
	Fidelity_errorrate float64 // could dictate how many colonies are checked in validation!
	Extensiontemp Temperature
	Hotstart bool
	StockConcentration Concentration // this is normally in U?
	TargetConcentration Concentration
	// this is also a glycerol solution rather than a watersolution!
}
*/

// Input parameters for this protocol (data)

// PCRprep parameters:

/*
	// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume

	//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...
	//FullTemplatesequence string // better to use Sid's type system here after proof of concept
	//FullTemplatelength int	// clearly could be calculated from the sequence... Sid will have a method to do this already so check!
	//TargetTemplatesequence string // better to use Sid's type system here after proof of concept
	//TargetTemplatelengthinBP int
*/
// Reaction parameters: (could be a entered as thermocycle parameters type possibly?)

//Denaturationtemp Temperature

// Should be calculated from primer and template binding
// should be calculated from template length and polymerase rate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// e.g. DMSO

// Physical outputs from this protocol with types

func _requirements() {
}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	// Mix components
	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.SampleForTotalVolume(_input.Buffer, _input.ReactionVolume)
	samples = append(samples, bufferSample)
	templateSample := mixer.Sample(_input.Template, _input.Templatevolume)
	samples = append(samples, templateSample)
	dntpSample := mixer.SampleForConcentration(_input.DNTPS, _input.DNTPconc)
	samples = append(samples, dntpSample)
	FwdPrimerSample := mixer.SampleForConcentration(_input.FwdPrimer, _input.FwdPrimerConc)
	samples = append(samples, FwdPrimerSample)
	RevPrimerSample := mixer.SampleForConcentration(_input.RevPrimer, _input.RevPrimerConc)
	samples = append(samples, RevPrimerSample)

	for _, additive := range _input.Additives {
		additiveSample := mixer.SampleForConcentration(additive, _input.Additiveconc)
		samples = append(samples, additiveSample)
	}

	polySample := mixer.SampleForConcentration(_input.PCRPolymerase, _input.TargetpolymeraseConcentration)
	samples = append(samples, polySample)
	reaction := execute.MixInto(_ctx,

		_input.OutPlate, samples...)

	// thermocycle parameters called from enzyme lookup:

	polymerase := _input.PCRPolymerase.CName

	extensionTemp := enzymes.DNApolymerasetemps[polymerase]["extensiontemp"]
	meltingTemp := enzymes.DNApolymerasetemps[polymerase]["meltingtemp"]

	// initial Denaturation

	execute.Incubate(_ctx,

		reaction, meltingTemp, _input.InitDenaturationtime, false)

	for i := 0; i < _input.Numberofcycles; i++ {

		// Denature

		execute.Incubate(_ctx,

			reaction, meltingTemp, _input.Denaturationtime, false)

		// Anneal
		execute.Incubate(_ctx,

			reaction, _input.AnnealingTemp, _input.Annealingtime, false)

		//extensiontime := TargetTemplatelengthinBP/PCRPolymerase.RateBPpers // we'll get type issues here so leave it out for now

		// Extend
		execute.Incubate(_ctx,

			reaction, extensionTemp, _input.Extensiontime, false)

	}
	// Final Extension
	execute.Incubate(_ctx,

		reaction, extensionTemp, _input.Finalextensiontime, false)

	// all done
	_output.Reaction = reaction
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
	Additiveconc                  wunit.Concentration
	Additives                     []*wtype.LHComponent
	AnnealingTemp                 wunit.Temperature
	Annealingtime                 wunit.Time
	Buffer                        *wtype.LHComponent
	DNTPS                         *wtype.LHComponent
	DNTPconc                      wunit.Concentration
	Denaturationtime              wunit.Time
	Extensiontemp                 wunit.Temperature
	Extensiontime                 wunit.Time
	Finalextensiontime            wunit.Time
	FwdPrimer                     *wtype.LHComponent
	FwdPrimerConc                 wunit.Concentration
	InitDenaturationtime          wunit.Time
	Numberofcycles                int
	OutPlate                      *wtype.LHPlate
	PCRPolymerase                 *wtype.LHComponent
	ReactionVolume                wunit.Volume
	RevPrimer                     *wtype.LHComponent
	RevPrimerConc                 wunit.Concentration
	TargetpolymeraseConcentration wunit.Concentration
	Template                      *wtype.LHComponent
	Templatevolume                wunit.Volume
}

type Output_ struct {
	Reaction *wtype.LHSolution
}
