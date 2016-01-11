/* Evaporation calculator based on
http://www.engineeringtoolbox.com/evaporation-water-surface-d_690.html

This engineering function may need to be improved to account for vapour pressure and surface tension

gs = Θ A (xs - x) / 3600         (1)

or

gh = Θ A (xs - x)

where

gs = amount of evaporated water per second (kg/s)

gh = amount of evaporated water per hour (kg/h)

Θ = (25 + 19 v) = evaporation coefficient (kg/m2h)

v = velocity of air above the water surface (m/s)

A = water surface area (m2)

xs = humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)

x = humidity ratio in the air (kg/kg) (kg H2O in kg Dry Air) */

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// ul

// cubesensor streams:
// in pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa. Not yet built in unit so we import it from wunit.
// input in deg C will be converted to Kelvin
// Percentage // density water vapor (kg/m3)

// // velocity of air above water in m/s ; could be calculated or measured

// time

// ul/h
// ul

func _EvaporationrateRequirements() {
}
func _EvaporationrateSetup(_ctx context.Context, _input *EvaporationrateInput) {
}
func _EvaporationrateSteps(_ctx context.Context, _input *EvaporationrateInput, _output *EvaporationrateOutput) {
}
func _EvaporationrateAnalysis(_ctx context.Context, _input *EvaporationrateInput, _output *EvaporationrateOutput) {

	var PWS float64 = eng.Pws(_input.Temp)
	var pw float64 = eng.Pw(_input.Relativehumidity, PWS) // vapour partial pressure in Pascals
	var Gh = (eng.Θ(_input.Liquid, _input.Airvelocity) *
		(labware.Labwaregeometry[_input.Platetype]["Surfacearea"] *
			((eng.Xs(PWS, _input.Pa)) - (eng.X(pw, _input.Pa))))) // Gh is rate of evaporation in kg/h
	evaporatedliquid := (Gh * (_input.Executiontime.SIValue() / 3600))                            // in kg
	evaporatedliquid = (evaporatedliquid * liquidclasses.Liquidclass[_input.Liquid]["ro"]) / 1000 // converted to litres
	_output.Evaporatedliquid = wunit.NewVolume((evaporatedliquid * 1000000), "ul")                // convert to ul

	_output.Evaporationrateestimate = Gh * 1000000 // ul/h if declared in parameters or data it doesn't need declaring again

	estimatedevaporationtime := _input.Volumeperwell.ConvertTo(wunit.ParsePrefixedUnit("ul")) / _output.Evaporationrateestimate
	_output.Estimatedevaporationtime = wunit.NewTime((estimatedevaporationtime * 3600), "s")

	_output.Status = fmt.Sprintln("Well Surface Area=",
		(labware.Labwaregeometry[_input.Platetype]["Surfacearea"])*1000000, "mm2",
		"evaporation rate =", Gh*1000000, "ul/h",
		"total evaporated liquid =", _output.Evaporatedliquid.ToString(), "after", _input.Executiontime.ToString(),
		"estimated evaporation time = ", _output.Estimatedevaporationtime.ToString())

} // works in either analysis or steps sections

func _EvaporationrateValidation(_ctx context.Context, _input *EvaporationrateInput, _output *EvaporationrateOutput) {
	if _output.Evaporatedliquid.SIValue() > _input.Volumeperwell.SIValue() {
		panic("not enough liquid")
	}
}
func _EvaporationrateRun(_ctx context.Context, input *EvaporationrateInput) *EvaporationrateOutput {
	output := &EvaporationrateOutput{}
	_EvaporationrateSetup(_ctx, input)
	_EvaporationrateSteps(_ctx, input, output)
	_EvaporationrateAnalysis(_ctx, input, output)
	_EvaporationrateValidation(_ctx, input, output)
	return output
}

func EvaporationrateRun(_ctx context.Context, input *EvaporationrateInput) *EvaporationrateSOutput {
	soutput := &EvaporationrateSOutput{}
	output := _EvaporationrateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func EvaporationrateNew() interface{} {
	return &EvaporationrateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &EvaporationrateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _EvaporationrateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &EvaporationrateInput{},
			Out: &EvaporationrateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type EvaporationrateElement struct {
	inject.CheckedRunner
}

type EvaporationrateInput struct {
	Airvelocity      wunit.Velocity
	Executiontime    wunit.Time
	Liquid           string
	Pa               wunit.Pressure
	Platetype        string
	Relativehumidity float64
	Temp             wunit.Temperature
	Volumeperwell    wunit.Volume
}

type EvaporationrateOutput struct {
	Estimatedevaporationtime wunit.Time
	Evaporatedliquid         wunit.Volume
	Evaporationrateestimate  float64
	Status                   string
}

type EvaporationrateSOutput struct {
	Data struct {
		Estimatedevaporationtime wunit.Time
		Evaporatedliquid         wunit.Volume
		Evaporationrateestimate  float64
		Status                   string
	}
	Outputs struct {
	}
}

func init() {
	c := Component{Name: "Evaporationrate", Constructor: EvaporationrateNew}
	c.Desc.Desc = ""
	c.Desc.Params = []ParamDesc{
		{Name: "Airvelocity", Desc: "// velocity of air above water in m/s ; could be calculated or measured\n", Kind: "Parameters"},
		{Name: "Executiontime", Desc: "time\n", Kind: "Parameters"},
		{Name: "Liquid", Desc: "", Kind: "Parameters"},
		{Name: "Pa", Desc: "cubesensor streams:\n\nin pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa. Not yet built in unit so we import it from wunit.\n", Kind: "Parameters"},
		{Name: "Platetype", Desc: "", Kind: "Parameters"},
		{Name: "Relativehumidity", Desc: "Percentage // density water vapor (kg/m3)\n", Kind: "Parameters"},
		{Name: "Temp", Desc: "input in deg C will be converted to Kelvin\n", Kind: "Parameters"},
		{Name: "Volumeperwell", Desc: "ul\n", Kind: "Parameters"},
		{Name: "Estimatedevaporationtime", Desc: "", Kind: "Data"},
		{Name: "Evaporatedliquid", Desc: "ul\n", Kind: "Data"},
		{Name: "Evaporationrateestimate", Desc: "ul/h\n", Kind: "Data"},
		{Name: "Status", Desc: "", Kind: "Data"},
	}
	addComponent(c)
}

// Go helper functions:

//Functions for rounding numbers to a specified number of decimal places (places):
/*func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}
*/
/* This function calculates Θ required for the evaporation calculator based on air velocity above the sample;
this will be important in a laminar flow cabinet, fume cabinet and when the plates are mixing:
*/

/*: 0.62198 * pws / (pa - pws), // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
"x":  0.62198 * pw / (pa - pw), */
