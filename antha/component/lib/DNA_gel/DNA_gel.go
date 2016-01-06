// example protocol for loading a DNAgel

package DNA_gel

import (
	//"LiquidHandler"
	//"Labware"
	//"coldplate"
	//"reagents"
	//"Devices"
	//"strconv"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

//DNAladder Volume // or should this be a concentration?

//DNAgelruntime time.Duration
//DNAgelwellcapacity Volume
//DNAgelnumberofwells int32
//Organism Taxonomy //= http://www.ncbi.nlm.nih.gov/nuccore/49175990?report=genbank
//Organismgenome Genome
//Target_DNA wtype.DNASequence
//Target_DNAsize float64 //Length
//Runvoltage float64
//AgarosePercentage Percentage
// polyerase kit sets key info such as buffer composition, which effects primer melting temperature for example, along with thermocycle parameters

// Data which is returned from this protocol, and data types

//	NumberofBands[] int
//Bandsizes[] Length
//Bandconc[]Concentration
//Pass bool
//PhotoofDNAgel Image

// Physical Inputs to this protocol with types

//WaterSolution
//WaterSolution //Chemspiderlink // not correct link but similar desirable

//Gel

//DNAladder *wtype.LHComponent//NucleicacidSolution
//Water *wtype.LHComponent//WaterSolution

//DNAgelbuffer *wtype.LHComponent//WaterSolution
//DNAgelNucleicacidintercalator *wtype.LHComponent//ToxicSolution // e.g. ethidium bromide, sybrsafe
//QC_sample *wtype.LHComponent//QC // this is a control
//DNASizeladder *wtype.LHComponent//WaterSolution
//Devices.Gelpowerpack Device
// need to calculate which DNASizeladder is required based on target sequence length and required resolution to distinguish from incorrect assembly possibilities

// Physical outputs from this protocol with types

//Gel
//

// No special requirements on inputs
func _requirements() {
	// None
	/* QC if negative result should still show band then include QC which will result in band // in reality this may never happen... the primers should be designed within antha too
	   control blank with no template_DNA */
}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func _setup(_ctx context.Context, _input *Input_) {
	/*control.config.per_DNAgel {
	load DNASizeladder(DNAgelrunvolume) // should run more than one per gel in many cases
	QC := mix (Loadingdye(loadingdyevolume), QC_sample(DNAgelrunvolume-loadingdyevolume))
	load QC(DNAgelrunvolume)
	}*/
}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	if len(_input.Samplenames) != _input.Samplenumber {
		panic(fmt.Sprintln("length of sample names:", len(_input.Samplenames), "is not equal to sample number:", _input.Samplenumber))
	}

	loadedsamples := make([]*wtype.LHSolution, 0)

	var DNAgelloadmix *wtype.LHComponent

	_input.Water.Type = "loadwater"

	for i := 0; i < _input.Samplenumber; i++ {
		// ready to add water to well
		waterSample := mixer.Sample(_input.Water, _input.Watervol)

		// load gel
		if _input.Loadingdyeinsample == false {
			DNAgelloadmixsolution := execute.MixInto(_ctx,

				_input.DNAgel,
				mixer.Sample(_input.Loadingdye, _input.Loadingdyevolume),
				mixer.SampleForTotalVolume(_input.Sampletotest, _input.DNAgelrunvolume),
			)
			DNAgelloadmix = wtype.SolutionToComponent(DNAgelloadmixsolution)
		} else {
			DNAgelloadmix = _input.Sampletotest
		}

		// Ensure  sample will be dispensed appropriately:

		// comment this line out to repeat load of same sample in all wells using first sample name
		DNAgelloadmix.CName = _input.Samplenames[0] //[i] //originalname + strconv.Itoa(i)

		// replacing following line with temporary hard code whilst developing protocol:
		DNAgelloadmix.Type = _input.Mixingpolicy
		//DNAgelloadmix.Type = "loadwater"

		loadedsample := execute.MixInto(_ctx,

			_input.DNAgel,
			waterSample,
			mixer.Sample(DNAgelloadmix, _input.DNAgelrunvolume),
		)

		loadedsamples = append(_output.Loadedsamples, loadedsample)
	}
	_output.Loadedsamples = loadedsamples
	// Then run the gel
	/* DNAgel := electrophoresis.Run(Loadedgel,Runvoltage,DNAgelruntime)

		// then analyse
	   	DNAgel.Visualise()
		PCR_product_length = call(assemblydesign_validation).PCR_product_length
		if DNAgel.Numberofbands() == 1
		&& DNAgel.Bandsize(DNAgel[0]) == PCR_product_length {
			Pass = true
			}

		incorrect_assembly_possibilities := assemblydesign_validation.Otherpossibleassemblysizes()

		for _, incorrect := range incorrect_assembly_possibilities {
			if  PCR_product_length == incorrect {
	    pass == false
		S := "matches size of incorrect assembly possibility"
		}

		//cherrypick(positive_colonies,recoverylocation)*/
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
	// need the control samples to be completed before doing the analysis

	//

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _validation(_ctx context.Context, _input *Input_, _output *Output_) {
	/* 	if calculatedbandsize == expected {
			stop
		}
		if calculatedbandsize != expected {
		if S == "matches size of incorrect assembly possibility" {
			call(assembly_troubleshoot)
			}
		} // loop at beginning should be designed to split labware resource optimally in the event of any failures e.g. if 96well capacity and 4 failures check 96/4 = 12 colonies of each to maximise chance of getting a hit
	    }
	    if repeat > 2
		stop
	    }
	    if (recoverylocation doesn't grow then use backup or repeat
		}
		if sequencingresults do not match expected then use backup or repeat
	    // TODO: */
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

//func cherrypick ()

type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	DNAgel             *wtype.LHPlate
	DNAgelrunvolume    wunit.Volume
	InPlate            *wtype.LHPlate
	Loadingdye         *wtype.LHComponent
	Loadingdyeinsample bool
	Loadingdyevolume   wunit.Volume
	Mixingpolicy       string
	Samplenames        []string
	Samplenumber       int
	Sampletotest       *wtype.LHComponent
	Water              *wtype.LHComponent
	Watervol           wunit.Volume
}

type Output_ struct {
	Loadedsamples []*wtype.LHSolution
}
