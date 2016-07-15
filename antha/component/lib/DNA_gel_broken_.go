// example protocol for loading a DNAgel

package lib

import (
	//"LiquidHandler"
	//"Labware"
	//"coldplate"
	//"reagents"
	//"Devices"
	//"strconv"
	//"fmt"
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

// gel
// plate to mix samples if required

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
func _DNA_gel_brokenRequirements() {
	// None
	/* QC if negative result should still show band then include QC which will result in band // in reality this may never happen... the primers should be designed within antha too
	   control blank with no template_DNA */
}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func _DNA_gel_brokenSetup(_ctx context.Context, _input *DNA_gel_brokenInput) {
	/*control.config.per_DNAgel {
	load DNASizeladder(DNAgelrunvolume) // should run more than one per gel in many cases
	QC := mix (Loadingdye(loadingdyevolume), QC_sample(DNAgelrunvolume-loadingdyevolume))
	load QC(DNAgelrunvolume)
	}*/
}

// The core process for this protocol, with the steps to be performed
// for every input
func _DNA_gel_brokenSteps(_ctx context.Context, _input *DNA_gel_brokenInput, _output *DNA_gel_brokenOutput) {

	loadedsamples := make([]*wtype.LHComponent, 0)
	wells := make([]string, 0)
	volumes := make([]wunit.Volume, 0)

	var DNAgelloadmix *wtype.LHComponent

	_input.Water.Type = wtype.LTloadwater

	var counter int

	for j := 0; j < _input.Samplenumber; j++ {
		for i := 0; i < len(_input.Samplenames); i++ {

			// ready to add water to well
			waterSample := mixer.Sample(_input.Water, _input.Watervol)

			// get position, ensuring the list is by row rather than by column
			position := _input.DNAgel.AllWellPositions(wtype.BYROW)[counter]

			// work out sample volume

			// copy volume
			samplevolume := (wunit.CopyVolume(_input.DNAgelrunvolume))

			// subtract volume of water
			samplevolume.Subtract(_input.Watervol)

			//get well coordinates
			//wellcoords := wtype.MakeWellCoordsA1(position)
			/*
				// if first column add ladder sample
				if wellcoords.X == 1 {

					loadedsample := MixTo(
				DNAgel.Type,
				position,
				1,
				waterSample,
				mixer.Sample(Ladder, samplevolume),
				)

				loadedsamples = append(Loadedsamples,loadedsample)
				wells = append(wells,position)
				volumes = append(volumes,loadedsample.Volume())
				counter++

				}else {
			*/

			_input.Sampletotest.CName = _input.Samplenames[i]

			//if strings.Contains(position)

			// load gel
			if _input.Loadingdyeinsample == false {

				_input.Loadingdye.Type, _ = wtype.LiquidTypeFromString("NeedToMix")

				DNAgelloadmixsolution := execute.MixTo(_ctx, _input.MixPlate.Type,
					position,
					1,
					mixer.Sample(_input.Sampletotest, samplevolume),
					mixer.Sample(_input.Loadingdye, _input.Loadingdyevolume),
				)
				DNAgelloadmix = DNAgelloadmixsolution
			} else {
				DNAgelloadmix = _input.Sampletotest
			}

			// Ensure  sample will be dispensed appropriately:

			// comment this line out to repeat load of same sample in all wells using first sample name
			DNAgelloadmix.CName = _input.Samplenames[i] //[i] //originalname + strconv.Itoa(i)

			// replacing following line with temporary hard code whilst developing protocol:
			DNAgelloadmix.Type, _ = wtype.LiquidTypeFromString(_input.Mixingpolicy)
			//DNAgelloadmix.Type = "loadwater"

			loadedsample := execute.MixInto(_ctx, _input.DNAgel,
				position,
				waterSample,
				mixer.Sample(DNAgelloadmix, samplevolume),
			)

			loadedsamples = append(_output.Loadedsamples, loadedsample)
			wells = append(wells, position)
			volumes = append(volumes, loadedsample.Volume())
			counter++
		}
	}
	_output.Loadedsamples = loadedsamples
	//fmt.Println(ProjectName+".csv", DNAgel, ProjectName+"gelouput", wells, Loadedsamples, volumes)
	// export to file
	//wtype.AutoExportPlateCSV(ProjectName+".csv",DNAgel)
	_output.Error = wtype.ExportPlateCSV(_input.ProjectName+".csv", _input.DNAgel, _input.ProjectName+"gelouput", wells, _output.Loadedsamples, volumes)
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
func _DNA_gel_brokenAnalysis(_ctx context.Context, _input *DNA_gel_brokenInput, _output *DNA_gel_brokenOutput) {
	// need the control samples to be completed before doing the analysis

	//

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _DNA_gel_brokenValidation(_ctx context.Context, _input *DNA_gel_brokenInput, _output *DNA_gel_brokenOutput) {
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
func _DNA_gel_brokenRun(_ctx context.Context, input *DNA_gel_brokenInput) *DNA_gel_brokenOutput {
	output := &DNA_gel_brokenOutput{}
	_DNA_gel_brokenSetup(_ctx, input)
	_DNA_gel_brokenSteps(_ctx, input, output)
	_DNA_gel_brokenAnalysis(_ctx, input, output)
	_DNA_gel_brokenValidation(_ctx, input, output)
	return output
}

func DNA_gel_brokenRunSteps(_ctx context.Context, input *DNA_gel_brokenInput) *DNA_gel_brokenSOutput {
	soutput := &DNA_gel_brokenSOutput{}
	output := _DNA_gel_brokenRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DNA_gel_brokenNew() interface{} {
	return &DNA_gel_brokenElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DNA_gel_brokenInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DNA_gel_brokenRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DNA_gel_brokenInput{},
			Out: &DNA_gel_brokenOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DNA_gel_brokenElement struct {
	inject.CheckedRunner
}

type DNA_gel_brokenInput struct {
	DNAgel             *wtype.LHPlate
	DNAgelrunvolume    wunit.Volume
	InPlate            *wtype.LHPlate
	Ladder             *wtype.LHComponent
	Loadingdye         *wtype.LHComponent
	Loadingdyeinsample bool
	Loadingdyevolume   wunit.Volume
	MixPlate           *wtype.LHPlate
	Mixingpolicy       string
	ProjectName        string
	Samplenames        []string
	Samplenumber       int
	Sampletotest       *wtype.LHComponent
	Water              *wtype.LHComponent
	Watervol           wunit.Volume
}

type DNA_gel_brokenOutput struct {
	Error         error
	Loadedsamples []*wtype.LHComponent
}

type DNA_gel_brokenSOutput struct {
	Data struct {
		Error error
	}
	Outputs struct {
		Loadedsamples []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(Component{Name: "DNA_gel_broken",
		Constructor: DNA_gel_brokenNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/DNA_gel/DNA_gel_broken.an",
			Params: []ParamDesc{
				{Name: "DNAgel", Desc: "gel\n", Kind: "Inputs"},
				{Name: "DNAgelrunvolume", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "Ladder", Desc: "", Kind: "Inputs"},
				{Name: "Loadingdye", Desc: "WaterSolution //Chemspiderlink // not correct link but similar desirable\n", Kind: "Inputs"},
				{Name: "Loadingdyeinsample", Desc: "", Kind: "Parameters"},
				{Name: "Loadingdyevolume", Desc: "", Kind: "Parameters"},
				{Name: "MixPlate", Desc: "plate to mix samples if required\n", Kind: "Inputs"},
				{Name: "Mixingpolicy", Desc: "wtype.LiquidType\n", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "", Kind: "Parameters"},
				{Name: "Samplenames", Desc: "", Kind: "Parameters"},
				{Name: "Samplenumber", Desc: "", Kind: "Parameters"},
				{Name: "Sampletotest", Desc: "WaterSolution\n", Kind: "Inputs"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Watervol", Desc: "", Kind: "Parameters"},
				{Name: "Error", Desc: "\tNumberofBands[] int\nBandsizes[] Length\nBandconc[]Concentration\nPass bool\nPhotoofDNAgel Image\n", Kind: "Data"},
				{Name: "Loadedsamples", Desc: "Gel\n", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}

//func cherrypick ()