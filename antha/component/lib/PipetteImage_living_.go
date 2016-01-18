// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

/*AntibioticVolume Volume
InducerVolume Volume
RepressorVolume Volume*/

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate
//Media *wtype.LHComponent
/*Antibiotic *wtype.LHComponent
Inducer *wtype.LHComponent
Repressor *wtype.LHComponent*/

// Physical outputs from this protocol with types

func _PipetteImage_livingRequirements() {

}

// Conditions to run on startup
func _PipetteImage_livingSetup(_ctx context.Context, _input *PipetteImage_livingInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_livingSteps(_ctx context.Context, _input *PipetteImage_livingInput, _output *PipetteImage_livingOutput) {

	_output.UniqueComponents = make([]string, 0)

	chosencolourpalette := image.AvailablePalettes[_input.Palettename]
	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, chosencolourpalette)

	if _input.UVimage {
		uvmap := image.AvailableComponentmaps[_input.Palettename]
		visiblemap := image.Visibleequivalentmaps[_input.Palettename]

		image.PrintFPImagePreview(_input.Imagefilename, _input.OutPlate, visiblemap, uvmap)
	}

	// get components from factory
	componentmap := make(map[string]*wtype.LHComponent, 0)

	colourtostringmap := image.AvailableComponentmaps[_input.Palettename]

	for _, colourname := range chosencolourpalette {

		componentname := colourtostringmap[colourname]

		componentmap[componentname] = factory.GetComponentByType(componentname)

	}
	fmt.Println(componentmap)

	solutions := make([]*wtype.LHSolution, 0)

	counter := 0

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		component := componentmap[colourtostringmap[colour]]

		if component.Type == "dna" {
			component.Type = "DoNotMix"
		}
		fmt.Println(image.Colourcomponentmap[colour])

		if _input.OnlythisColour != "" {

			if image.Colourcomponentmap[colour] == _input.OnlythisColour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				counter = counter + 1
				fmt.Println("wells", counter)
				//mediaSample := mixer.SampleForTotalVolume(Media, VolumePerWell)
				//components = append(components,mediaSample)
				/*antibioticSample := mixer.Sample(Antibiotic, AntibioticVolume)
				components = append(components,antibioticSample)
				repressorSample := mixer.Sample(Repressor, RepressorVolume)
				components = append(components,repressorSample)
				inducerSample := mixer.Sample(Inducer, InducerVolume)
				components = append(components,inducerSample)*/
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				components = append(components, pixelSample)
				solution := execute.MixTo(_ctx, _input.OutPlate, locationkey, components...)
				execute.Incubate(_ctx, solution, _input.IncTemp, _input.IncTime, true)
				solutions = append(solutions, solution)
			}

		} else {
			if component.CName != _input.Notthiscolour {

				_output.UniqueComponents = append(_output.UniqueComponents, component.CName)

				counter = counter + 1
				fmt.Println("wells", counter)
				//mediaSample := mixer.SampleForTotalVolume(Media, VolumePerWell)
				//components = append(components,mediaSample)
				/*antibioticSample := mixer.Sample(Antibiotic, AntibioticVolume)
				components = append(components,antibioticSample)
				repressorSample := mixer.Sample(Repressor, RepressorVolume)
				components = append(components,repressorSample)
				inducerSample := mixer.Sample(Inducer, InducerVolume)
				components = append(components,inducerSample)*/
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				components = append(components, pixelSample)
				solution := execute.MixTo(_ctx, _input.OutPlate, locationkey, components...)

				execute.Incubate(_ctx, solution, _input.IncTemp, _input.IncTime, true)
				solutions = append(solutions, solution)
			}
		}
	}

	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

	_output.UniqueComponents = search.RemoveDuplicates(_output.UniqueComponents)
	text.Print("Unique Components:", _output.UniqueComponents)
	fmt.Println("number of unique components", len(_output.UniqueComponents))
	_output.Pixels = solutions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_livingAnalysis(_ctx context.Context, _input *PipetteImage_livingInput, _output *PipetteImage_livingOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_livingValidation(_ctx context.Context, _input *PipetteImage_livingInput, _output *PipetteImage_livingOutput) {

}
func _PipetteImage_livingRun(_ctx context.Context, input *PipetteImage_livingInput) *PipetteImage_livingOutput {
	output := &PipetteImage_livingOutput{}
	_PipetteImage_livingSetup(_ctx, input)
	_PipetteImage_livingSteps(_ctx, input, output)
	_PipetteImage_livingAnalysis(_ctx, input, output)
	_PipetteImage_livingValidation(_ctx, input, output)
	return output
}

func PipetteImage_livingRunSteps(_ctx context.Context, input *PipetteImage_livingInput) *PipetteImage_livingSOutput {
	soutput := &PipetteImage_livingSOutput{}
	output := _PipetteImage_livingRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_livingNew() interface{} {
	return &PipetteImage_livingElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_livingInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_livingRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_livingInput{},
			Out: &PipetteImage_livingOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_livingElement struct {
	inject.CheckedRunner
}

type PipetteImage_livingInput struct {
	Imagefilename     string
	IncTemp           wunit.Temperature
	IncTime           wunit.Time
	InoculationVolume wunit.Volume
	Notthiscolour     string
	OnlythisColour    string
	OutPlate          *wtype.LHPlate
	Palettename       string
	UVimage           bool
	VolumePerWell     wunit.Volume
}

type PipetteImage_livingOutput struct {
	Numberofpixels   int
	Pixels           []*wtype.LHSolution
	UniqueComponents []string
}

type PipetteImage_livingSOutput struct {
	Data struct {
		Numberofpixels   int
		UniqueComponents []string
	}
	Outputs struct {
		Pixels []*wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "PipetteImage_living",
		Constructor: PipetteImage_livingNew,
		Desc: ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteLivingimage.an",
			Params: []ParamDesc{
				{Name: "Imagefilename", Desc: "AntibioticVolume Volume\n\tInducerVolume Volume\n\tRepressorVolume Volume\n", Kind: "Parameters"},
				{Name: "IncTemp", Desc: "", Kind: "Parameters"},
				{Name: "IncTime", Desc: "", Kind: "Parameters"},
				{Name: "InoculationVolume", Desc: "", Kind: "Parameters"},
				{Name: "Notthiscolour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\nMedia *wtype.LHComponent\nAntibiotic *wtype.LHComponent\n\tInducer *wtype.LHComponent\n\tRepressor *wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "Palettename", Desc: "", Kind: "Parameters"},
				{Name: "UVimage", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
				{Name: "UniqueComponents", Desc: "", Kind: "Data"},
			},
		},
	})
}
