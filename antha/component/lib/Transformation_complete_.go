package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

//= 50.(uL)

//= 2 (hours)

//Shakerspeed float64 // correct type?

//Plateoutdilution float64

/*ReactionVolume wunit.Volume
PartConc wunit.Concentration
VectorConc wunit.Concentration
AtpVol wunit.Volume
ReVol wunit.Volume
LigVol wunit.Volume
ReactionTemp wunit.Temperature
ReactionTime wunit.Time
InactivationTemp wunit.Temperature
InactivationTime wunit.Time
*/

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _Transformation_completeRequirements() {
}

// Conditions to run on startup
func _Transformation_completeSetup(_ctx context.Context, _input *Transformation_completeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Transformation_completeSteps(_ctx context.Context, _input *Transformation_completeInput, _output *Transformation_completeOutput) {
	competentcells := make([]*wtype.LHComponent, 0)
	competentcells = append(competentcells, _input.CompetentCells)
	readycompetentcells := execute.MixInto(_ctx, _input.OutPlate, "", competentcells...)             // readycompetentcells IS now a LHSolution
	execute.Incubate(_ctx, readycompetentcells, _input.Preplasmidtemp, _input.Preplasmidtime, false) // we can incubate an LHSolution so this is fine

	readycompetentcellsComp := wtype.SolutionToComponent(readycompetentcells)

	competetentcellmix := mixer.Sample(readycompetentcellsComp, _input.CompetentCellvolumeperassembly) // ERROR! mixer.Sample needs a liquid, not an LHSolution! however, the typeIIs method worked with a *wtype.LHComponent from inputs!
	transformationmix := make([]*wtype.LHComponent, 0)
	transformationmix = append(transformationmix, competetentcellmix)
	DNAsample := mixer.Sample(_input.Reaction, _input.Reactionvolume)
	transformationmix = append(transformationmix, DNAsample)

	transformedcells := execute.MixInto(_ctx, _input.OutPlate, "", transformationmix...)

	execute.Incubate(_ctx, transformedcells, _input.Postplasmidtemp, _input.Postplasmidtime, false)

	transformedcellsComp := wtype.SolutionToComponent(transformedcells)

	recoverymix := make([]*wtype.LHComponent, 0)
	recoverymixture := mixer.Sample(_input.Recoverymedium, _input.Recoveryvolume)

	recoverymix = append(recoverymix, transformedcellsComp) // ERROR! transformedcells is now an LHSolution, not a liquid, so can't be used here
	recoverymix = append(recoverymix, recoverymixture)
	recoverymix2 := execute.MixInto(_ctx, _input.OutPlate, "", recoverymix...)

	execute.Incubate(_ctx, recoverymix2, _input.Recoverytemp, _input.Recoverytime, true)

	recoverymix2Comp := wtype.SolutionToComponent(recoverymix2)

	plateout := mixer.Sample(recoverymix2Comp, _input.Plateoutvolume) // ERROR! recoverymix2 is now an LHSolution, not a liquid, so can't be used here
	platedculture := execute.MixInto(_ctx, _input.AgarPlate, "", plateout)

	_output.Platedculture = platedculture

	/*atpSample := mixer.Sample(Atp, AtpVol)
	samples = append(samples, atpSample)
	vectorSample := mixer.SampleForConcentration(Vector, VectorConc)
	samples = append(samples, vectorSample)

	for _, part := range Parts {
		partSample := mixer.SampleForConcentration(part, PartConc)
		samples = append(samples, partSample)
	}

	reSample := mixer.Sample(RestrictionEnzyme, ReVol)
	samples = append(samples, reSample)
	ligSample := mixer.Sample(Ligase, LigVol)
	samples = append(samples, ligSample)


	// incubate the reaction mixture

	Incubate(reaction, ReactionTemp, ReactionTime, false)

	// inactivate

	Incubate(reaction, InactivationTemp, InactivationTime, false)

	// all done
	Reaction = reaction

	readycompetentcells := Incubate (CompetentCells,Preplasmidtemp, Preplasmidtime, false)


	product := Mix (Reaction(ReactionVolume), readycompetentcells(CompetentCellvolumeperassembly))
	transformedcells := Incubate (product, Postplasmidtime,Postplasmidtemp,false)
	recoverymixture := Mix (transformedcells, Recoverymedium (Recoveryvolume)) // or alternative recovery medium
	Incubate (recoverymixture, Recoverytime, Recoverytemp, Shakerspeed)
	platedculture := MixInto(AgarPlate, Plateoutvolume)

	Platedculture = platedculture

	*/
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Transformation_completeAnalysis(_ctx context.Context, _input *Transformation_completeInput, _output *Transformation_completeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Transformation_completeValidation(_ctx context.Context, _input *Transformation_completeInput, _output *Transformation_completeOutput) {
}
func _Transformation_completeRun(_ctx context.Context, input *Transformation_completeInput) *Transformation_completeOutput {
	output := &Transformation_completeOutput{}
	_Transformation_completeSetup(_ctx, input)
	_Transformation_completeSteps(_ctx, input, output)
	_Transformation_completeAnalysis(_ctx, input, output)
	_Transformation_completeValidation(_ctx, input, output)
	return output
}

func Transformation_completeRunSteps(_ctx context.Context, input *Transformation_completeInput) *Transformation_completeSOutput {
	soutput := &Transformation_completeSOutput{}
	output := _Transformation_completeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Transformation_completeNew() interface{} {
	return &Transformation_completeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Transformation_completeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Transformation_completeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Transformation_completeInput{},
			Out: &Transformation_completeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Transformation_completeElement struct {
	inject.CheckedRunner
}

type Transformation_completeInput struct {
	AgarPlate                      *wtype.LHPlate
	CompetentCells                 *wtype.LHComponent
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       *wtype.LHPlate
	Plateoutvolume                 wunit.Volume
	Postplasmidtemp                wunit.Temperature
	Postplasmidtime                wunit.Time
	Preplasmidtemp                 wunit.Temperature
	Preplasmidtime                 wunit.Time
	Reaction                       *wtype.LHComponent
	Reactionvolume                 wunit.Volume
	Recoverymedium                 *wtype.LHComponent
	Recoverytemp                   wunit.Temperature
	Recoverytime                   wunit.Time
	Recoveryvolume                 wunit.Volume
}

type Transformation_completeOutput struct {
	Platedculture *wtype.LHSolution
}

type Transformation_completeSOutput struct {
	Data struct {
	}
	Outputs struct {
		Platedculture *wtype.LHSolution
	}
}

func init() {
	addComponent(Component{Name: "Transformation_complete",
		Constructor: Transformation_completeNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Transformation/Transformation_almostworking.an",
			Params: []ParamDesc{
				{Name: "AgarPlate", Desc: "", Kind: "Inputs"},
				{Name: "CompetentCells", Desc: "", Kind: "Inputs"},
				{Name: "CompetentCellvolumeperassembly", Desc: "= 50.(uL)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Plateoutvolume", Desc: "", Kind: "Parameters"},
				{Name: "Postplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Postplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "Preplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Preplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Inputs"},
				{Name: "Reactionvolume", Desc: "", Kind: "Parameters"},
				{Name: "Recoverymedium", Desc: "", Kind: "Inputs"},
				{Name: "Recoverytemp", Desc: "", Kind: "Parameters"},
				{Name: "Recoverytime", Desc: "= 2 (hours)\n", Kind: "Parameters"},
				{Name: "Recoveryvolume", Desc: "", Kind: "Parameters"},
				{Name: "Platedculture", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
