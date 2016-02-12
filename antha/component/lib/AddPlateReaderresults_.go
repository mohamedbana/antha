package lib

import (
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wutil"
	//"github.com/antha-lang/antha/antha/anthalib/mixer"
	//"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	//"path/filepath"
	//antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"strconv"

// Input parameters for this protocol (data)

//= "AbsMV"

//Wavelength            int    = 440
// = "Abs Spectrum"

//= []string{"P9"}

//= []string{"P24"}

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AddPlateReaderresultsRequirements() {
}

// Conditions to run on startup
func _AddPlateReaderresultsSetup(_ctx context.Context, _input *AddPlateReaderresultsInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input

func _AddPlateReaderresultsSteps(_ctx context.Context, _input *AddPlateReaderresultsInput, _output *AddPlateReaderresultsOutput) {

	var marsdata parser.MarsData
	var err error

	if _input.DOEFiletype == "DX" || _input.DOEFiletype == "Design Expert" {

		marsdata, err = parser.ParseMarsXLSXOutput(_input.MarsResultsFileXLSX, _input.SheetNumber)
		if err != nil {
			panic(err)
		}

	}

	runs, err := doe.RunsFromDXDesign(_input.DOEFilewithwelllocationsadded, []string{"Pre_MIX", "POST_MIX"})

	// find optimal wavlength from scan of positive control and blank
	optimalwavelength := marsdata.FindOptimalWavelength(_input.ManualControls[0], _input.Blanks[0], "Raw Data")

	// range through pairing up wells from mars output and doe design

	measuredoptimalwavelengths := make([]int, 0)

	//for i, additional := range AdditionalFactors {
	for _, run := range runs {

		well, err := run.GetAdditionalInfo("Well ID")
		if err != nil {
			panic(err)
		}

		//if run.CheckAdditionalInfo(additional, AdditionalValues[i]) && additional == Additionalfactortoresponsepair[0] {
		/*
			average, err := marsdata.AbsorbanceReading(well.(string), Wavelength, ReadingTypeinMarsFile)
			if err != nil {
				panic(err)
			}
		*/
		// check optimal difference for each well
		meassuredoptwavelength := marsdata.FindOptimalWavelength(well.(string), _input.Blanks[0], "Raw Data")
		measuredoptimalwavelengths = append(measuredoptimalwavelengths, meassuredoptwavelength)
		// blank correct

		samples := []string{well.(string)}

		blankcorrected, err := marsdata.BlankCorrect(samples, _input.Blanks, optimalwavelength, _input.ReadingTypeinMarsFile)
		run.AddResponseValue(_input.Responsecolumntofill, blankcorrected)
		//	}
	}
	//	}

	//runs, err := doe.RunsFromDXDesign(xlsxwithresultsadded string, []string{"Pre_MIX", "POST_MIX"})

	_ = doe.DXXLSXFilefromRuns(runs, _output.OutputFilename)
	//OutputFilename = doe.XLfileFromRuns(runs)

	fmt.Println("Optimal wavelength from manual", optimalwavelength)
	fmt.Println("Optimal wavelength of each sample", measuredoptimalwavelengths)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AddPlateReaderresultsAnalysis(_ctx context.Context, _input *AddPlateReaderresultsInput, _output *AddPlateReaderresultsOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AddPlateReaderresultsValidation(_ctx context.Context, _input *AddPlateReaderresultsInput, _output *AddPlateReaderresultsOutput) {
}
func _AddPlateReaderresultsRun(_ctx context.Context, input *AddPlateReaderresultsInput) *AddPlateReaderresultsOutput {
	output := &AddPlateReaderresultsOutput{}
	_AddPlateReaderresultsSetup(_ctx, input)
	_AddPlateReaderresultsSteps(_ctx, input, output)
	_AddPlateReaderresultsAnalysis(_ctx, input, output)
	_AddPlateReaderresultsValidation(_ctx, input, output)
	return output
}

func AddPlateReaderresultsRunSteps(_ctx context.Context, input *AddPlateReaderresultsInput) *AddPlateReaderresultsSOutput {
	soutput := &AddPlateReaderresultsSOutput{}
	output := _AddPlateReaderresultsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AddPlateReaderresultsNew() interface{} {
	return &AddPlateReaderresultsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AddPlateReaderresultsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AddPlateReaderresultsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AddPlateReaderresultsInput{},
			Out: &AddPlateReaderresultsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AddPlateReaderresultsElement struct {
	inject.CheckedRunner
}

type AddPlateReaderresultsInput struct {
	Blanks                        []string
	DOEFiletype                   string
	DOEFilewithwelllocationsadded string
	ManualControls                []string
	MarsResultsFileXLSX           string
	ReadingTypeinMarsFile         string
	Responsecolumntofill          string
	SheetNumber                   int
}

type AddPlateReaderresultsOutput struct {
	OutputFilename string
}

type AddPlateReaderresultsSOutput struct {
	Data struct {
		OutputFilename string
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "AddPlateReaderresults",
		Constructor: AddPlateReaderresultsNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/platereader/ReadPlatereaderOutput.an",
			Params: []ParamDesc{
				{Name: "Blanks", Desc: "= []string{\"P9\"}\n", Kind: "Parameters"},
				{Name: "DOEFiletype", Desc: "", Kind: "Parameters"},
				{Name: "DOEFilewithwelllocationsadded", Desc: "", Kind: "Parameters"},
				{Name: "ManualControls", Desc: "= []string{\"P24\"}\n", Kind: "Parameters"},
				{Name: "MarsResultsFileXLSX", Desc: "", Kind: "Parameters"},
				{Name: "ReadingTypeinMarsFile", Desc: "Wavelength            int    = 440\n\n= \"Abs Spectrum\"\n", Kind: "Parameters"},
				{Name: "Responsecolumntofill", Desc: "= \"AbsMV\"\n", Kind: "Parameters"},
				{Name: "SheetNumber", Desc: "", Kind: "Parameters"},
				{Name: "OutputFilename", Desc: "", Kind: "Data"},
			},
		},
	})
}
