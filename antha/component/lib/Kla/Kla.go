/* Islam, R. S., Tisi, D., Levy, M. S. & Lye, G. J. Scale-up of Escherichia coli growth and recombinant protein expression conditions from microwell to laboratory and pilot scale based on matched kLa. Biotechnol. Bioeng. 99, 1128–1139 (2008).

equation (6)

func kLa_squaremicrowell = (3.94 x 10E-4) * (D/dv)* ai * RE^1.91 * exp ^ (a * Fr^b) // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case

kla = dimensionless
	var D = diffusion coefficient, m2 􏰀 s􏰁1
	var dv = microwell vessel diameter, m
	var ai = initial specific surface area, m􏰁1
	var RE = Reynolds number, (ro * n * dv * 2/mu), dimensionless
		var	ro	= density, kg 􏰀/ m􏰁3
		var	n 	= shaking frequency, s􏰁1
		var	mu	= viscosity, kg 􏰀/ m􏰁 /􏰀 s
	const exp = Eulers number, 2.718281828

	var Fr = Froude number = dt(2 * math.Pi * n)^2 /(2 * g), (dimensionless)
		var dt = shaking amplitude, m
		const g = acceleration due to gravity, m 􏰀/ s􏰁2
	const	a = constant
	const	b = constant
*/
// make type /time and units of /hour and per second
// check accuracy against literature and experimental values
package Kla

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Setpoints"
//"github.com/montanaflynn/stats"

//float64

//diffusion coefficient, m2 􏰀 s􏰁1 // from wikipedia: Oxygen (dis) - Water (l) 	@25 degrees C 	2.10x10−5 cm2/s // should call from elsewhere really
// add temp etc?

//float64

//float64

func _requirements() {
}
func _setup(_ctx context.Context, _input *Input_) {
}
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	dv := labware.Labwaregeometry[_input.Platetype]["dv"] // microwell vessel diameter, m 0.017 //
	ai := labware.Labwaregeometry[_input.Platetype]["ai"] // initial specific surface area, /m 96.0

	ro := liquidclasses.Liquidclass[_input.Liquid]["ro"] //density, kg 􏰀/ m􏰁3 999.7 // environment dependent
	mu := liquidclasses.Liquidclass[_input.Liquid]["mu"] //0.001           environment dependent                        //liquidclasses.Liquidclass[liquid]["mu"] viscosity, kg 􏰀/ m􏰁 /􏰀 s

	var n float64 //shaking frequency per second

	if _input.Rpm.Unit().RawSymbol() == `/s` {
		n = _input.Rpm.RawValue()
	} else if _input.Rpm.Unit().RawSymbol() == `/min` {
		n = _input.Rpm.RawValue() / 60
	}

	//n = Rpm / 60 //shaking frequency, s􏰁1
	//var RE = Reynolds number, (ro * n * dv * 2/mu), dimensionless
	//const exp = Eulers number, 2.718281828
	//Fr = Froude number = dt(2 * math.Pi * n)^2 /(2 * g), (dimensionless)

	dt := devices.Shaker[_input.Shakertype]["dt"] //0.008                                  //shaking amplitude, m // move to shaker package

	a := labware.Labwaregeometry[_input.Platetype]["a"] //0.88   //
	b := labware.Labwaregeometry[_input.Platetype]["b"] //1.24

	Fr := eng.Froude(dt, n, eng.G)
	Re := eng.RE(ro, n, mu, dv)
	_output.Necessaryshakerspeed = eng.Shakerspeed(_input.TargetRE, ro, mu, dv)

	Vl := _input.Fillvolume.SIValue()
	Sigma := liquidclasses.Liquidclass[_input.Liquid]["sigma"]

	// Check Ncrit! original paper used this to calculate speed in shallow round well plates... double check paper

	// add loop to use correct formula dependent on Platetype etc...
	// currently only one plate type supported
	//Criticalshakerspeed := "error"
	if labware.Labwaregeometry[_input.Platetype]["numberofwellsides"] == 4.0 {
		_output.Ncrit = eng.Ncrit_srw(Sigma, dv, Vl, ro, dt)
	} /*else{Criticalshakerspeed := "error: kla estimation for this plate type not yet implemented"}
	/*if i == 4.0 {
		Criticalshakerspeed := "error"
	}
	*/
	//Criticalshakerspeed := stats.Round(eng.Ncrit_srw(Sigma, dv, Vl , ro , dt ),3)

	if Re > 5E3 {
		_output.Flowstate = fmt.Sprintln("Flowstate = Turbulent flow")
	}

	//klainputs :=fmt.Sprintln("D",D,"dv", dv,"ai", ai,"Re", Re,"a", a,"Fr", Fr,"b", b)

	_output.CalculatedKla = eng.KLa_squaremicrowell(_input.D, dv, ai, Re, a, Fr, b)

	_output.Status = fmt.Sprintln("TargetRE = ", _input.TargetRE, "Calculated Reynolds number = ", Re, "shakerspeedrequired for targetRE= ", _output.Necessaryshakerspeed.ToString() /* *60 */, "Froude number = ", Fr, "kla =", _output.CalculatedKla, "/h", "Ncrit	=", _output.Ncrit.ToString() /*,"/S"*/)
	//CalculatedKla = setpoints.CalculateKlasquaremicrowell(Platetype, Liquid, Rpm, Shakertype, TargetRE, D)

}
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {

} // works in either analysis or steps sections

func _validation(_ctx context.Context, _input *Input_, _output *Output_) {
	/*if Evaporatedliquid > Volumeperwell {
	panic("not enough liquid") */
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
	D          float64
	Fillvolume wunit.Volume
	Liquid     string
	Platetype  string
	Rpm        wunit.Rate
	Shakertype string
	TargetRE   float64
}

type Output_ struct {
	CalculatedKla        float64
	Flowstate            string
	Ncrit                wunit.Rate
	Necessaryshakerspeed wunit.Rate
	Status               string
}
