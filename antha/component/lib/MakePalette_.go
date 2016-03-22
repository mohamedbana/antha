// Generates instructions to make a pallette of all colours in an image
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

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
//"image/color"

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _MakePaletteRequirements() {

}

// Conditions to run on startup
func _MakePaletteSetup(_ctx context.Context, _input *MakePaletteInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePaletteSteps(_ctx context.Context, _input *MakePaletteInput, _output *MakePaletteOutput) {

	//var chosencolourpalette color.Palette

	// make pallette of colours from image
	chosencolourpalette := image.MakeSmallPalleteFromImage(_input.Imagefilename, _input.OutPlate)

	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette)

	// remove duplicates
	//positiontocolourmap = image.RemoveDuplicatesValuesfromMap(positiontocolourmap)

	fmt.Println("positions", positiontocolourmap)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	//solutions := image.PipetteImagebyBlending(OutPlate, positiontocolourmap,Cyan, Magenta, Yellow,Black, VolumeForFullcolour)

	for _, colour := range positiontocolourmap {

		if colour != nil {
			components := make([]*wtype.LHComponent, 0)

			cmyk := image.ColourtoCMYK(colour)

			var maxuint8 uint8 = 255

			if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

				continue

			} else {

				counter = counter + 1

				if cmyk.C > 0 {

					cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if cyanvol.RawValue() < 10 && cyanvol.Unit().PrefixedSymbol() == "ul" {
						cyanvol.SetValue(10)
					}

					cyanSample := mixer.Sample(_input.Cyan, cyanvol)
					components = append(components, cyanSample)
				}

				if cmyk.Y > 0 {
					yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if yellowvol.RawValue() < 10 && yellowvol.Unit().PrefixedSymbol() == "ul" {
						yellowvol.SetValue(10)
					}

					yellowSample := mixer.Sample(_input.Yellow, yellowvol)
					components = append(components, yellowSample)
				}

				if cmyk.M > 0 {
					magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if magentavol.RawValue() < 10 && magentavol.Unit().PrefixedSymbol() == "ul" {
						magentavol.SetValue(10)
					}

					magentaSample := mixer.Sample(_input.Magenta, magentavol)
					components = append(components, magentaSample)
				}

				if cmyk.K > 0 {
					blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())

					if blackvol.RawValue() < 10 && blackvol.Unit().PrefixedSymbol() == "ul" {
						blackvol.SetValue(10)
					}

					blackSample := mixer.Sample(_input.Black, blackvol)
					components = append(components, blackSample)
				}

				solution := execute.MixInto(_ctx, _input.OutPlate, "", components...)
				solutions = append(solutions, solution)

			}

		}
	}

	_output.Colours = solutions
	_output.Numberofcolours = len(_output.Colours)
	fmt.Println("Unique Colours =", _output.Numberofcolours)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakePaletteAnalysis(_ctx context.Context, _input *MakePaletteInput, _output *MakePaletteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePaletteValidation(_ctx context.Context, _input *MakePaletteInput, _output *MakePaletteOutput) {

}
func _MakePaletteRun(_ctx context.Context, input *MakePaletteInput) *MakePaletteOutput {
	output := &MakePaletteOutput{}
	_MakePaletteSetup(_ctx, input)
	_MakePaletteSteps(_ctx, input, output)
	_MakePaletteAnalysis(_ctx, input, output)
	_MakePaletteValidation(_ctx, input, output)
	return output
}

func MakePaletteRunSteps(_ctx context.Context, input *MakePaletteInput) *MakePaletteSOutput {
	soutput := &MakePaletteSOutput{}
	output := _MakePaletteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePaletteNew() interface{} {
	return &MakePaletteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePaletteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePaletteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePaletteInput{},
			Out: &MakePaletteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakePaletteElement struct {
	inject.CheckedRunner
}

type MakePaletteInput struct {
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	Magenta             *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	VolumeForFullcolour wunit.Volume
	Yellow              *wtype.LHComponent
}

type MakePaletteOutput struct {
	Colours         []*wtype.LHComponent
	Numberofcolours int
}

type MakePaletteSOutput struct {
	Data struct {
		Numberofcolours int
	}
	Outputs struct {
		Colours []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "MakePalette",
		Constructor: MakePaletteNew,
		Desc: ComponentDesc{
			Desc: "Generates instructions to make a pallette of all colours in an image\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/MakePallete.an",
			Params: []ParamDesc{
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Cyan", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "Magenta", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "Yellow", Desc: "", Kind: "Inputs"},
				{Name: "Colours", Desc: "", Kind: "Outputs"},
				{Name: "Numberofcolours", Desc: "", Kind: "Data"},
			},
		},
	})
}
