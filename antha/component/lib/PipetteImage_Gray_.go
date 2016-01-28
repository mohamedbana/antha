// Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes
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
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _PipetteImage_GrayRequirements() {

}

// Conditions to run on startup
func _PipetteImage_GraySetup(_ctx context.Context, _input *PipetteImage_GrayInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PipetteImage_GraySteps(_ctx context.Context, _input *PipetteImage_GrayInput, _output *PipetteImage_GrayOutput) {

	chosencolourpalette := image.AvailablePalettes["Gray"]

	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette)

	solutions := make([]*wtype.LHSolution, 0)

	counter := 0

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		gray := image.ColourtoGrayscale(colour)

		var maxuint8 uint8 = 255

		if gray.Y == 0 {

			continue

		} else {

			counter = counter + 1

			if gray.Y < maxuint8 {
				watervol := wunit.NewVolume((float64(maxuint8-gray.Y) / float64(maxuint8) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				fmt.Println(watervol)
				if watervol.RawValue() < 10 && watervol.Unit().PrefixedSymbol() == "ul" {
					watervol.SetValue(10)
				}
				waterSample := mixer.Sample(_input.Diluent, watervol)
				components = append(components, waterSample)
			}
			blackvol := wunit.NewVolume((float64(gray.Y/maxuint8) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
			fmt.Println("blackvol", blackvol)
			if blackvol.RawValue() < 10 && blackvol.Unit().PrefixedSymbol() == "ul" {
				blackvol.SetValue(10)
			}
			blackSample := mixer.Sample(_input.Black, blackvol)
			components = append(components, blackSample)

			solution := execute.MixTo(_ctx, _input.OutPlate, locationkey, components...)
			solutions = append(solutions, solution)

		}
	}

	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PipetteImage_GrayAnalysis(_ctx context.Context, _input *PipetteImage_GrayInput, _output *PipetteImage_GrayOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PipetteImage_GrayValidation(_ctx context.Context, _input *PipetteImage_GrayInput, _output *PipetteImage_GrayOutput) {

}
func _PipetteImage_GrayRun(_ctx context.Context, input *PipetteImage_GrayInput) *PipetteImage_GrayOutput {
	output := &PipetteImage_GrayOutput{}
	_PipetteImage_GraySetup(_ctx, input)
	_PipetteImage_GraySteps(_ctx, input, output)
	_PipetteImage_GrayAnalysis(_ctx, input, output)
	_PipetteImage_GrayValidation(_ctx, input, output)
	return output
}

func PipetteImage_GrayRunSteps(_ctx context.Context, input *PipetteImage_GrayInput) *PipetteImage_GraySOutput {
	soutput := &PipetteImage_GraySOutput{}
	output := _PipetteImage_GrayRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PipetteImage_GrayNew() interface{} {
	return &PipetteImage_GrayElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PipetteImage_GrayInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PipetteImage_GrayRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PipetteImage_GrayInput{},
			Out: &PipetteImage_GrayOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PipetteImage_GrayElement struct {
	inject.CheckedRunner
}

type PipetteImage_GrayInput struct {
	Black               *wtype.LHComponent
	Diluent             *wtype.LHComponent
	Imagefilename       string
	OutPlate            *wtype.LHPlate
	VolumeForFullcolour wunit.Volume
}

type PipetteImage_GrayOutput struct {
	Numberofpixels int
	Pixels         []*wtype.LHSolution
}

type PipetteImage_GraySOutput struct {
	Data struct {
		Numberofpixels int
	}
	Outputs struct {
		Pixels []*wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "PipetteImage_Gray",
		Constructor: PipetteImage_GrayNew,
		Desc: ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteImage_Gray.an",
			Params: []ParamDesc{
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
