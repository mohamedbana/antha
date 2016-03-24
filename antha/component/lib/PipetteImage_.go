// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

//AvailableColours []string

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _PipetteImageRequirements() {

}

// Conditions to run on startup
func _PipetteImageSetup(_ctx context.Context, _input *PipetteImageInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImageSteps(_ctx context.Context, _input *PipetteImageInput, _output *PipetteImageOutput) {

	availableColours := make([]string, 0)

	for _, component := range _input.Colourcomponents {
		availableColours = append(availableColours, component.CName)
	}

	//chosencolourpalette := image.AvailableComponentmaps[Palettename]
	//

	subpalette := image.MakeSubPallette(_input.Palettename, availableColours)
	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &subpalette, _input.Rotate)

	//Pixels = image.PipetteImagetoPlate(OutPlate, positiontocolourmap, AvailableColours, Colourcomponents, VolumePerWell)
	/*
		componentmap, err := image.MakestringtoComponentMap(AvailableColours, Colourcomponents)
		if err != nil {
			panic(err)
		}
	*/

	// get components from factory
	componentmap := make(map[string]*wtype.LHComponent, 0)

	colourtostringmap := image.AvailableComponentmaps[_input.Palettename]

	submap := image.MakeSubMapfromMap(colourtostringmap, availableColours)

	for colourname := range submap {

		componentname := colourtostringmap[colourname]

		componentmap[componentname] = factory.GetComponentByType(componentname)

	}

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0
	// currently set up to only pipette if yellow (to make visualisation easier in trilution simulator
	for locationkey, colour := range positiontocolourmap {

		component := componentmap[image.Colourcomponentmap[colour]]

		if component.TypeName() == "dna" {
			component.Type = wtype.LTDoNotMix //"DoNotMix"
		}
		fmt.Println(image.Colourcomponentmap[colour])

		if _input.OnlythisColour != "" {

			if image.Colourcomponentmap[colour] == _input.OnlythisColour {
				counter = counter + 1
				fmt.Println("wells", counter)
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
				solutions = append(solutions, solution)
			}

		} else {
			if component.CName != _input.NotthisColour {
				counter = counter + 1
				fmt.Println("wells", counter)
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				solution := execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, pixelSample)
				solutions = append(solutions, solution)
			}
		}
	}

	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)
	_output.Pixels = solutions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImageAnalysis(_ctx context.Context, _input *PipetteImageInput, _output *PipetteImageOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImageValidation(_ctx context.Context, _input *PipetteImageInput, _output *PipetteImageOutput) {

}
func _PipetteImageRun(_ctx context.Context, input *PipetteImageInput) *PipetteImageOutput {
	output := &PipetteImageOutput{}
	_PipetteImageSetup(_ctx, input)
	_PipetteImageSteps(_ctx, input, output)
	_PipetteImageAnalysis(_ctx, input, output)
	_PipetteImageValidation(_ctx, input, output)
	return output
}

func PipetteImageRunSteps(_ctx context.Context, input *PipetteImageInput) *PipetteImageSOutput {
	soutput := &PipetteImageSOutput{}
	output := _PipetteImageRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImageNew() interface{} {
	return &PipetteImageElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImageInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImageRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImageInput{},
			Out: &PipetteImageOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImageElement struct {
	inject.CheckedRunner
}

type PipetteImageInput struct {
	Colourcomponents []*wtype.LHComponent
	Imagefilename    string
	NotthisColour    string
	OnlythisColour   string
	OutPlate         *wtype.LHPlate
	Palettename      string
	Rotate           bool
	VolumePerWell    wunit.Volume
}

type PipetteImageOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHComponent
}

type PipetteImageSOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "PipetteImage",
		Constructor: PipetteImageNew,
		Desc: ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteImage.an",
			Params: []ParamDesc{
				{Name: "Colourcomponents", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "NotthisColour", Desc: "", Kind: "Parameters"},
				{Name: "OnlythisColour", Desc: "AvailableColours []string\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "Palettename", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerWell", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
