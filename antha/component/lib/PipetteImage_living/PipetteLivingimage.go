// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package PipetteImage_living

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

	chosencolourpalette := image.AvailablePalettes[_input.Palettename]
	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, chosencolourpalette)

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

		component := componentmap[colourtostringmap[colour]]

		if component.Type == "dna" {
			component.Type = "DoNotMix"
		}
		fmt.Println(image.Colourcomponentmap[colour])

		if _input.OnlythisColour != "" {

			if image.Colourcomponentmap[colour] == _input.OnlythisColour {
				counter = counter + 1
				fmt.Println("wells", counter)
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				solution := execute.MixTo(_ctx,

					_input.OutPlate, locationkey, pixelSample)
				solutions = append(solutions, solution)
			}

		} else {
			if component.CName != "white" {
				counter = counter + 1
				fmt.Println("wells", counter)
				pixelSample := mixer.Sample(component, _input.VolumePerWell)
				solution := execute.MixTo(_ctx,

					_input.OutPlate, locationkey, pixelSample)
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
	Imagefilename  string
	OnlythisColour string
	OutPlate       *wtype.LHPlate
	Palettename    string
	VolumePerWell  wunit.Volume
}

type Output_ struct {
	Numberofpixels int
	Pixels         []*wtype.LHSolution
}
