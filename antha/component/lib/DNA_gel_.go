// example protocol for loading a DNAgel

package lib

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

//wtype.LiquidType

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
func _DNA_gelRequirements() {
	// None
	/* QC if negative result should still show band then include QC which will result in band // in reality this may never happen... the primers should be designed within antha too
	   control blank with no template_DNA */
}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func _DNA_gelSetup(_ctx context.Context, _input *DNA_gelInput) {
	/*control.config.per_DNAgel {
	load DNASizeladder(DNAgelrunvolume) // should run more than one per gel in many cases
	QC := mix (Loadingdye(loadingdyevolume), QC_sample(DNAgelrunvolume-loadingdyevolume))
	load QC(DNAgelrunvolume)
	}*/
}

// The core process for this protocol, with the steps to be performed
// for every input
func _DNA_gelSteps(_ctx context.Context, _input *DNA_gelInput, _output *DNA_gelOutput) {

	if len(_input.Samplenames) != _input.Samplenumber {
		panic(fmt.Sprintln("length of sample names:", len(_input.Samplenames), "is not equal to sample number:", _input.Samplenumber))
	}

	loadedsamples := make([]*wtype.LHComponent, 0)

	var DNAgelloadmix *wtype.LHComponent

	_input.Water.Type = wtype.LTloadwater

	for i := 0; i < _input.Samplenumber; i++ {
		// ready to add water to well
		waterSample := mixer.Sample(_input.Water, _input.Watervol)

		// load gel
		if _input.Loadingdyeinsample == false {
			DNAgelloadmixsolution := execute.MixInto(_ctx, _input.DNAgel,
				"",
				mixer.Sample(_input.Loadingdye, _input.Loadingdyevolume),
				mixer.SampleForTotalVolume(_input.Sampletotest, _input.DNAgelrunvolume),
			)
			DNAgelloadmix = DNAgelloadmixsolution
		} else {
			DNAgelloadmix = _input.Sampletotest
		}

		// Ensure  sample will be dispensed appropriately:

		// comment this line out to repeat load of same sample in all wells using first sample name
		DNAgelloadmix.CName = _input.Samplenames[0] //[i] //originalname + strconv.Itoa(i)

		// replacing following line with temporary hard code whilst developing protocol:
		DNAgelloadmix.Type, _ = wtype.LiquidTypeFromString(_input.Mixingpolicy)
		//DNAgelloadmix.Type = "loadwater"

		loadedsample := execute.MixInto(_ctx, _input.DNAgel,
			"",
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
func _DNA_gelAnalysis(_ctx context.Context, _input *DNA_gelInput, _output *DNA_gelOutput) {
	// need the control samples to be completed before doing the analysis

	//

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _DNA_gelValidation(_ctx context.Context, _input *DNA_gelInput, _output *DNA_gelOutput) {
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
func _DNA_gelRun(_ctx context.Context, input *DNA_gelInput) *DNA_gelOutput {
	output := &DNA_gelOutput{}
	_DNA_gelSetup(_ctx, input)
	_DNA_gelSteps(_ctx, input, output)
	_DNA_gelAnalysis(_ctx, input, output)
	_DNA_gelValidation(_ctx, input, output)
	return output
}

func DNA_gelRunSteps(_ctx context.Context, input *DNA_gelInput) *DNA_gelSOutput {
	soutput := &DNA_gelSOutput{}
	output := _DNA_gelRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DNA_gelNew() interface{} {
	return &DNA_gelElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DNA_gelInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DNA_gelRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DNA_gelInput{},
			Out: &DNA_gelOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DNA_gelElement struct {
	inject.CheckedRunner
}

type DNA_gelInput struct {
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

type DNA_gelOutput struct {
	Loadedsamples []*wtype.LHComponent
}

type DNA_gelSOutput struct {
	Data struct {
	}
	Outputs struct {
		Loadedsamples []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "DNA_gel",
		Constructor: DNA_gelNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/DNA_gel/DNA_gel.an",
			Params: []ParamDesc{
				{Name: "DNAgel", Desc: "Gel\n", Kind: "Inputs"},
				{Name: "DNAgelrunvolume", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "Loadingdye", Desc: "WaterSolution //Chemspiderlink // not correct link but similar desirable\n", Kind: "Inputs"},
				{Name: "Loadingdyeinsample", Desc: "", Kind: "Parameters"},
				{Name: "Loadingdyevolume", Desc: "", Kind: "Parameters"},
				{Name: "Mixingpolicy", Desc: "wtype.LiquidType\n", Kind: "Parameters"},
				{Name: "Samplenames", Desc: "", Kind: "Parameters"},
				{Name: "Samplenumber", Desc: "", Kind: "Parameters"},
				{Name: "Sampletotest", Desc: "WaterSolution\n", Kind: "Inputs"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Watervol", Desc: "", Kind: "Parameters"},
				{Name: "Loadedsamples", Desc: "Gel\n", Kind: "Outputs"},
			},
		},
	})
}

//func cherrypick ()
