// Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/internal/github.com/disintegration/imaging"
)

// Input parameters for this protocol (data)

// as a proportion of 1 i.e. 0.5 == 50%
//SkipBlackforlowervol bool

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

	var blackvol wunit.Volume

	var maxuint8 uint8 = 255

	var minuint8 uint8

	_output.ShadesofGrey = make([]int, 0)

	chosencolourpalette := image.AvailablePalettes["Gray"]

	if _input.PosterizeImage {
		_, _input.Imagefilename = image.Posterize(_input.Imagefilename, _input.PosterizeLevels)
	}

	image.CheckAllResizealgorithms(_input.Imagefilename, _input.OutPlate, _input.Rotate, imaging.AllResampleFilters)

	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, _input.Rotate, _input.AutoRotate)

	solutions := make([]*wtype.LHComponent, 0)

	counter := 0

	for locationkey, colour := range positiontocolourmap {

		//components := make([]*wtype.LHComponent, 0)

		var solution *wtype.LHComponent
		//var mixedsolution *wtype.LHComponent

		gray := image.ColourtoGrayscale(colour)

		if _input.Negative == false {
			gray.Y = maxuint8 - gray.Y
		}

		minuint8 = uint8(_input.MinimumBlackpercentagethreshold * float64(maxuint8))

		fmt.Println("brand new minuint8", minuint8)

		if gray.Y < minuint8 {
			fmt.Println("skipping well:", locationkey)
			continue

		} else {

			counter = counter + 1

			// check if shade of grey has already been used in image
			greyindexinpalette := chosencolourpalette.Index(colour)

			alreadythere := search.Contains(_output.ShadesofGrey, greyindexinpalette)

			if alreadythere == false {
				_output.ShadesofGrey = append(_output.ShadesofGrey, greyindexinpalette)
			}

			if gray.Y < maxuint8 {
				watervol := wunit.NewVolume((float64(maxuint8-gray.Y) / float64(maxuint8) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				fmt.Println("new well", locationkey, "water vol", watervol.ToString())
				// force hv tip choice
				if _input.OnlyHighVolumetips && watervol.RawValue() < 21 && watervol.Unit().PrefixedSymbol() == "ul" {
					watervol.SetValue(21)
				}
				waterSample := mixer.Sample(_input.Diluent, watervol)
				//components = append(components, waterSample)
				solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, waterSample)
			}
			if gray.Y == maxuint8 {
				blackvol = _input.VolumeForFullcolour
			} else {
				blackvol = wunit.NewVolume((float64(gray.Y) / float64(maxuint8) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
			}

			fmt.Println("new well", locationkey, "black vol", blackvol.ToString())

			//Black.Type = wtype.LiquidTypeFromString("NeedToMix")

			if _input.DontMix {
				_input.Black.Type = wtype.LTDISPENSEABOVE
			} else {
				_input.Black.Type = wtype.LiquidTypeFromString("NeedToMix")
			}

			//fmt.Println("blackvol2",blackvol.ToString())
			if _input.OnlyHighVolumetips && blackvol.RawValue() < 21 && blackvol.Unit().PrefixedSymbol() == "ul" {
				blackvol.SetValue(21)
			}

			blackSample := mixer.Sample(_input.Black, blackvol)
			//components = append(components, blackSample)

			if solution != nil {
				solution = execute.Mix(_ctx, solution, blackSample)
			} else {
				solution = execute.MixTo(_ctx, _input.OutPlate.Type, locationkey, 1, blackSample)
			}
			solutions = append(solutions, solution)

		}
	}

	_output.NumberofShadesofGrey = len(_output.ShadesofGrey)
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
	AutoRotate                      bool
	Black                           *wtype.LHComponent
	Diluent                         *wtype.LHComponent
	DontMix                         bool
	Imagefilename                   string
	MinimumBlackpercentagethreshold float64
	Negative                        bool
	OnlyHighVolumetips              bool
	OutPlate                        *wtype.LHPlate
	PosterizeImage                  bool
	PosterizeLevels                 int
	Rotate                          bool
	VolumeForFullcolour             wunit.Volume
}

type PipetteImage_GrayOutput struct {
	NumberofShadesofGrey int
	Numberofpixels       int
	Pixels               []*wtype.LHComponent
	ShadesofGrey         []int
}

type PipetteImage_GraySOutput struct {
	Data struct {
		NumberofShadesofGrey int
		Numberofpixels       int
		ShadesofGrey         []int
	}
	Outputs struct {
		Pixels []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "PipetteImage_Gray",
		Constructor: PipetteImage_GrayNew,
		Desc: ComponentDesc{
			Desc: "Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes\n",
			Path: "antha/component/an/Liquid_handling/PipetteImage/PipetteImage_Gray.an",
			Params: []ParamDesc{
				{Name: "AutoRotate", Desc: "", Kind: "Parameters"},
				{Name: "Black", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DontMix", Desc: "", Kind: "Parameters"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "MinimumBlackpercentagethreshold", Desc: "as a proportion of 1 i.e. 0.5 == 50%\n", Kind: "Parameters"},
				{Name: "Negative", Desc: "", Kind: "Parameters"},
				{Name: "OnlyHighVolumetips", Desc: "SkipBlackforlowervol bool\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "InPlate *wtype.LHPlate\n", Kind: "Inputs"},
				{Name: "PosterizeImage", Desc: "", Kind: "Parameters"},
				{Name: "PosterizeLevels", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "VolumeForFullcolour", Desc: "", Kind: "Parameters"},
				{Name: "NumberofShadesofGrey", Desc: "", Kind: "Data"},
				{Name: "Numberofpixels", Desc: "", Kind: "Data"},
				{Name: "Pixels", Desc: "", Kind: "Outputs"},
				{Name: "ShadesofGrey", Desc: "", Kind: "Data"},
			},
		},
	})
}
