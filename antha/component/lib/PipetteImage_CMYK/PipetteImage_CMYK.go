// Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes
package PipetteImage_CMYK

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

//"image/color"

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _requirements() {

}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	//var chosencolourpalette color.Palette
	chosencolourpalette := image.AvailablePalettes["Plan9"]
	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette)

	solutions := make([]*wtype.LHSolution, 0)

	counter := 0

	//solutions := image.PipetteImagebyBlending(OutPlate, positiontocolourmap,Cyan, Magenta, Yellow,Black, VolumeForFullcolour)

	for locationkey, colour := range positiontocolourmap {

		components := make([]*wtype.LHComponent, 0)

		cmyk := image.ColourtoCMYK(colour)

		var maxuint8 uint8 = 255

		if cmyk.C == 0 && cmyk.Y == 0 && cmyk.M == 0 && cmyk.K == 0 {

			continue

		} else {

			counter = counter + 1

			if cmyk.C > 0 {

				cyanvol := wunit.NewVolume(((float64(cmyk.C) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				cyanSample := mixer.Sample(_input.Cyan, cyanvol)
				components = append(components, cyanSample)
			}

			if cmyk.Y > 0 {
				yellowvol := wunit.NewVolume(((float64(cmyk.Y) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				yellowSample := mixer.Sample(_input.Yellow, yellowvol)
				components = append(components, yellowSample)
			}

			if cmyk.M > 0 {
				magentavol := wunit.NewVolume(((float64(cmyk.M) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				magentaSample := mixer.Sample(_input.Magenta, magentavol)
				components = append(components, magentaSample)
			}

			if cmyk.K > 0 {
				blackvol := wunit.NewVolume(((float64(cmyk.K) / float64(maxuint8)) * _input.VolumeForFullcolour.RawValue()), _input.VolumeForFullcolour.Unit().PrefixedSymbol())
				blackSample := mixer.Sample(_input.Black, blackvol)
				components = append(components, blackSample)
			}

			solution := execute.MixTo(_ctx,

				_input.OutPlate, locationkey, components...)
			solutions = append(solutions, solution)

		}
	}

	_output.Pixels = solutions
	_output.Numberofpixels = len(_output.Pixels)
	fmt.Println("Pixels =", _output.Numberofpixels)

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
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	Magenta             *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	VolumeForFullcolour wunit.Volume
	Yellow              *wtype.LHComponent
}

type Output_ struct {
	Numberofpixels int
	Pixels         []*wtype.LHSolution
}
