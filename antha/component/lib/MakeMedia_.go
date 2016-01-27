package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

//Volume //Mass // Should be Mass

//  +/- x  e.g. 7.0 +/- 0.2

//LiqComponentkeys	[]string
//Solidcomponentkeys	[]string // name or barcode id
//Acidkey string
//Basekey string

// Physical Inputs to this protocol with types

// should be new type or field indicating solid and mass
/*Acid				*wtype.LHComponent
Base 				*wtype.LHComponent
*/

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _MakeMediaRequirements() {

}

// Conditions to run on startup
func _MakeMediaSetup(_ctx context.Context, _input *MakeMediaInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeMediaSteps(_ctx context.Context, _input *MakeMediaInput, _output *MakeMediaOutput) {
	recipestring := make([]string, 0)
	var step string
	stepcounter := 1 // counting from 1 is what makes us human
	liquids := make([]*wtype.LHComponent, 0)
	step = text.Print("Recipe for: ", _input.Name)
	recipestring = append(recipestring, step)

	for i, liq := range _input.LiqComponents {
		liqsamp := mixer.Sample(liq, _input.LiqComponentVolumes[i])
		liquids = append(liquids, liqsamp)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add "+_input.LiqComponentVolumes[i].ToString()+" of "+liq.CName)
		recipestring = append(recipestring, step)
		stepcounter++
	}

	//solids := make([]*wtype.LHComponent,0)

	for k, sol := range _input.SolidComponents {
		solsamp := mixer.SampleSolidtoLiquid(sol, _input.SolidComponentMasses[k], _input.SolidComponentDensities[k])
		liquids = append(liquids, solsamp)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add "+_input.SolidComponentMasses[k].ToString()+" of "+sol.CName)
		recipestring = append(recipestring, step)
		stepcounter = stepcounter + k
	}

	watersample := mixer.SampleForTotalVolume(_input.Water, _input.TotalVolume)
	liquids = append(liquids, watersample)
	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add up to "+_input.TotalVolume.ToString()+" of "+_input.Water.CName)
	recipestring = append(recipestring, step)
	stepcounter++

	// Add pH handling functions and driver calls etc...

	description := fmt.Sprint("adjust pH to ", _input.PH_setPoint, " +/-", _input.PH_tolerance, " for temp ", _input.PH_setPointTemp.ToString(), "C")
	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", description)
	recipestring = append(recipestring, step)
	stepcounter++

	/*
		prepH := MixInto(Vessel,liquids...)

		pHactual := prepH.Measure("pH")

		step = text.Print("pH measured = ", pHactual)
		recipestring = append(recipestring,step)

		//pHactual = wutil.Roundto(pHactual,PH_tolerance)

		pHmax := PH_setpoint + PH_tolerance
		pHmin := PH_setpoint - PH_tolerance

		if pHactual < pHmax || pHactual < pHmin {
			// basically just a series of sample, stir, wait and recheck pH
		Media, newph, componentadded = prepH.AdjustpH(PH_setPoint, pHactual, PH_setPointTemp,Acid,Base)

		step = text.Print("Adjusted pH = ", newpH)
		recipestring = append(recipestring,step)

		step = text.Print("Component added = ", componentadded.Vol + componentadded.Vunit + " of " + componentadded.Conc + componentadded.Cunit + " " + componentadded.CName + )
		recipestring = append(recipestring,step)
		}
	*/
	_output.Media = execute.MixInto(_ctx, _input.Vessel, "", liquids...)

	_output.Status = fmt.Sprintln(recipestring)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeMediaAnalysis(_ctx context.Context, _input *MakeMediaInput, _output *MakeMediaOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeMediaValidation(_ctx context.Context, _input *MakeMediaInput, _output *MakeMediaOutput) {
}
func _MakeMediaRun(_ctx context.Context, input *MakeMediaInput) *MakeMediaOutput {
	output := &MakeMediaOutput{}
	_MakeMediaSetup(_ctx, input)
	_MakeMediaSteps(_ctx, input, output)
	_MakeMediaAnalysis(_ctx, input, output)
	_MakeMediaValidation(_ctx, input, output)
	return output
}

func MakeMediaRunSteps(_ctx context.Context, input *MakeMediaInput) *MakeMediaSOutput {
	soutput := &MakeMediaSOutput{}
	output := _MakeMediaRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeMediaNew() interface{} {
	return &MakeMediaElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeMediaInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeMediaRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeMediaInput{},
			Out: &MakeMediaOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakeMediaElement struct {
	inject.CheckedRunner
}

type MakeMediaInput struct {
	LiqComponentVolumes     []wunit.Volume
	LiqComponents           []*wtype.LHComponent
	Name                    string
	PH_setPoint             float64
	PH_setPointTemp         wunit.Temperature
	PH_tolerance            float64
	SolidComponentDensities []wunit.Density
	SolidComponentMasses    []wunit.Mass
	SolidComponents         []*wtype.LHComponent
	TotalVolume             wunit.Volume
	Vessel                  *wtype.LHPlate
	Water                   *wtype.LHComponent
}

type MakeMediaOutput struct {
	Media  *wtype.LHSolution
	Status string
}

type MakeMediaSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Media *wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "MakeMedia",
		Constructor: MakeMediaNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/MakeMedia/MakeMedia.an",
			Params: []ParamDesc{
				{Name: "LiqComponentVolumes", Desc: "", Kind: "Parameters"},
				{Name: "LiqComponents", Desc: "", Kind: "Inputs"},
				{Name: "Name", Desc: "", Kind: "Parameters"},
				{Name: "PH_setPoint", Desc: "", Kind: "Parameters"},
				{Name: "PH_setPointTemp", Desc: "", Kind: "Parameters"},
				{Name: "PH_tolerance", Desc: " +/- x  e.g. 7.0 +/- 0.2\n", Kind: "Parameters"},
				{Name: "SolidComponentDensities", Desc: "", Kind: "Parameters"},
				{Name: "SolidComponentMasses", Desc: "Volume //Mass // Should be Mass\n", Kind: "Parameters"},
				{Name: "SolidComponents", Desc: "should be new type or field indicating solid and mass\n", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Vessel", Desc: "", Kind: "Inputs"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Media", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	})
}

/*
type Mole struct {
	number float64
}*/
