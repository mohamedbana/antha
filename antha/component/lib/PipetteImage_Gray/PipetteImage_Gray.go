// Generates instructions to pipette out a defined image onto a defined plate by blending cyan magenta yellow and black dyes
package PipetteImage_Gray

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

func _requirements() {

}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

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
	Diluent             *wtype.LHComponent
	Imagefilename       string
	OutPlate            *wtype.LHPlate
	VolumeForFullcolour wunit.Volume
}

type Output_ struct {
	Numberofpixels int
	Pixels         []*wtype.LHSolution
}
