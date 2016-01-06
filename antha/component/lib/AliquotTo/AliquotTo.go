// variant of aliquot.an whereby the low level MixTo command is used to pipette by row

package AliquotTo

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _requirements() {

}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		panic("Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHSolution, 0)

	// work out well coordinates for any plate
	wellpositionarray := make([]string, 0)

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//k := 0
	for j := 0; j < _input.OutPlate.WlsY; j++ {
		for i := 0; i < _input.OutPlate.WlsX; i++ { //countingfrom1iswhatmakesushuman := j + 1
			//k = k + 1
			wellposition := string(alphabet[j]) + strconv.Itoa(i+1)
			//fmt.Println(wellposition, k)
			wellpositionarray = append(wellpositionarray, wellposition)
		}

	}

	for k := 0; k < _input.NumberofAliquots; k++ {
		if _input.Solution.Type == "dna" {
			_input.Solution.Type = "DoNotMix"
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
		aliquot := execute.MixTo(_ctx,

			_input.OutPlate, wellpositionarray[k], aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	_output.Aliquots = aliquots
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
	InPlate          *wtype.LHPlate
	NumberofAliquots int
	OutPlate         *wtype.LHPlate
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Output_ struct {
	Aliquots []*wtype.LHSolution
}
