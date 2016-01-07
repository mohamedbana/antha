package LoadGel

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//    RunVoltage      Int
//    RunLength       Time

//preload well with 10uL of water
//protein samples for running
//96 well plate with water, marker and samples
//Gel to load ie OutPlate

//Run length in cm, and protein band height and pixed density after digital scanning

func _setup(_ctx context.Context, _input *Input_) {
}

func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.Sample(_input.Water, _input.WaterVolume)
	waterSample.CName = _input.WaterName
	samples = append(samples, waterSample)

	loadSample := mixer.Sample(_input.Protein, _input.LoadVolume)
	loadSample.CName = _input.SampleName
	samples = append(samples, loadSample)
	fmt.Println("This is a list of samples for loading:", samples)

	_output.RunSolution = execute.MixInto(_ctx,

		_input.GelPlate, samples...)
}

func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
}

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
	GelPlate    *wtype.LHPlate
	InPlate     *wtype.LHPlate
	LoadVolume  wunit.Volume
	Protein     *wtype.LHComponent
	SampleName  string
	Water       *wtype.LHComponent
	WaterName   string
	WaterVolume wunit.Volume
}

type Output_ struct {
	RunSolution *wtype.LHSolution
	Status      string
}
