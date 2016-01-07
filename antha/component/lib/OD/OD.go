// Example OD measurement protocol.
// Computes the OD and dry cell weight estimate from absorbance reading
// TODO: implement replicates from parameters
package OD

import (
	//"liquid handler"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"standard_labware"
// Input parameters for this protocol (data)

//= uL(100)
//= uL(0)
//Total_volume Volume//= ul (sample_volume+diluent_volume)
//Wavelength //= nm(600)
//Diluent_type //= (PBS)
//= (0.25)
//Replicate_count uint32 //= 1 // Note: 1 replicate means experiment is in duplicate, etc.
// calculate path length? - takes place under plate reader since this will only be necessary for plate reader protocols? labware?
// Data which is returned from this protocol, and data types
//= 0.0533
//WellCrosssectionalArea float64// should be calculated from plate and well type automatically

//Absorbance
//Absorbance
//(pathlength corrected)

//R_squared float32
//Control_absorbance [control_curve_points+1]float64//Absorbance
//Control_concentrations [control_curve_points+1]float64

// Physical Inputs to this protocol with types

//Culture

// Physical outputs from this protocol with types

// None

func _requirements() {
	// sufficient sample volume available to sacrifice
}
func _setup(_ctx context.Context, _input *Input_) {
	/*control.Config(config.per_plate)
	var control_blank[total_volume]WaterSolution

	blank_absorbance = platereader.Read(ODplate,control_blank, wavelength)*/
}
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	var product *wtype.LHSolution //WaterSolution

	for {
		product = execute.MixInto(_ctx,

			_input.ODplate, mixer.Sample(_input.Sampletotest, _input.Sample_volume), mixer.Sample(_input.Diluent, _input.Diluent_volume))
		/*Is it necessary to include platetype in Read function?
		or is the info on volume, opacity, pathlength etc implied in LHSolution?*/
		_output.Sample_absorbance = platereader.ReadAbsorbance(*_input.ODplate, *product, _input.Wlength)

		if _output.Sample_absorbance.Reading < 1 {
			break
		}
		_input.Diluent_volume.Mvalue += 1 //diluent_volume = diluent_volume + 1

	}
} // serial dilution or could write element for finding optimum dilution or search historical data
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
	// Need to substract blank from measurement; normalise to path length of 1cm for OD value; apply conversion factor to estimate dry cell weight

	_output.Blankcorrected_absorbance = platereader.Blankcorrect(_output.Sample_absorbance, _input.Blank_absorbance)
	volumetopathlengthconversionfactor := wunit.NewLength(_input.Heightof100ulinm, "m")                               //WellCrosssectionalArea
	_output.OD = platereader.PathlengthCorrect(volumetopathlengthconversionfactor, _output.Blankcorrected_absorbance) // 0.0533 could be written as function of labware and liquid volume (or measureed height)
	_output.Estimateddrycellweight_conc = wunit.NewConcentration(_output.OD.Reading*_input.ODtoDCWconversionfactor, "g/L")
}
func _validation(_ctx context.Context, _input *Input_, _output *Output_) { /*
		if Sample_absorbance > 1 {
		panic("Sample likely needs further dilution")
		}
		if Sample_absorbance < 0.1 {
		warn("Low OD, sample likely needs increased volume")
		}
		}*/
	// TODO: add test of replicate variance
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
	Blank_absorbance        wtype.Absorbance
	Diluent                 *wtype.LHComponent
	Diluent_volume          wunit.Volume
	Heightof100ulinm        float64
	ODplate                 *wtype.LHPlate
	ODtoDCWconversionfactor float64
	Sample_volume           wunit.Volume
	Sampletotest            *wtype.LHComponent
	Wlength                 float64
}

type Output_ struct {
	Blankcorrected_absorbance   wtype.Absorbance
	Estimateddrycellweight_conc wunit.Concentration
	OD                          wtype.Absorbance
	Sample_absorbance           wtype.Absorbance
}
