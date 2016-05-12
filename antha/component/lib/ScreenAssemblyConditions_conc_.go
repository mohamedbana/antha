// Assemble multiple assemblies using TypeIIs construct assembly
package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wutil"

	"github.com/antha-lang/antha/microArch/driver/liquidhandling"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"

	antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
	"strings"
)

// file containing design for liquid handling DOE

// fixed factors

// Reaction temperature
// Reaction time
// Prefix for reaction names
// Reaction volume

//variable

// variable but coupled
// Volumes corresponding to input parts // coupled with PartsArray and should be equal in length
// Names corresonding to input parts

//fixed

// Output plate

// Variable
// Input parts, one per assembly

// List of assembled parts

func _ScreenAssemblyConditions_concSetup(_ctx context.Context, _input *ScreenAssemblyConditions_concInput) {
}

func _ScreenAssemblyConditions_concSteps(_ctx context.Context, _input *ScreenAssemblyConditions_concInput, _output *ScreenAssemblyConditions_concOutput) {

	// validate presence of doe design file in anthapath
	if antha.Anthafileexists(_input.LHDOEFile) == false {
		fmt.Println("This DOE file ", _input.LHDOEFile, " was not found in anthapath ~.antha. Please move it there, change file name and type in antha-lang/antha/microarch/driver/makelhpolicy.go and recompile antha to use this liquidhandling doe design")
		fmt.Println("currently set to ", liquidhandling.DOEliquidhandlingFile, " type ", liquidhandling.DXORJMP)
	} else {
		fmt.Println("found lhpolicy doe file", _input.LHDOEFile)
	}

	// declare some global variables for use later
	var wellpositionarray = make([]string, 0)
	var alphabet = wutil.MakeAlphabetArray()
	_output.Runtowelllocationmap = make(map[string]string)
	counter := 0
	var platenum = 1

	// range through well coordinates
	for j := 0; j < _input.OutPlate.WlsX; j++ {
		for i := 0; i < _input.OutPlate.WlsY; i++ { //countingfrom1iswhatmakesushuman := j + 1
			//k = k + 1
			wellposition := string(alphabet[i]) + strconv.Itoa(j+1)
			//fmt.Println(wellposition, k)
			wellpositionarray = append(wellpositionarray, wellposition)
		}

	}

	_, names, runs, err := liquidhandling.PolicyMakerfromDesign(_input.DXORJMP, _input.LHDOEFile, "DOE_run")
	if err != nil {
		panic(err)
	}

	var newRuns = make([]doe.Run, 0)

	for l := 0; l < _input.Replicates; l++ {

		for k := range _input.DNAMassesPerReaction {

			for j := range _input.PartNamesArray {

				for i := 0; i < len(runs); i++ {

					if counter == (_input.OutPlate.WlsX * _input.OutPlate.WlsY) {
						fmt.Println("plate full, counter = ", counter)
						platenum++
						counter = 0
					}

					var PartVolsArray = make([]wunit.Volume, 0)

					for _, conc := range _input.PartConcArray[j] {

						PartVolsArray = append(PartVolsArray, wunit.NewVolume(float64((_input.DNAMassesPerReaction[k].SIValue()/conc.SIValue())*1000000), "ul"))
					}

					fmt.Println("counter:", counter)
					fmt.Println("WellPositionarray", wellpositionarray, "OutPlate.WlsX", _input.OutPlate.WlsX)

					result := TypeIISConstructAssemblyMMX_forscreenRunSteps(_ctx, &TypeIISConstructAssemblyMMX_forscreenInput{ReactionVolume: _input.ReactionVolume,
						PartVols:           PartVolsArray,
						PartNames:          _input.PartNamesArray[j],
						MasterMixVolume:    _input.MastermixVolume,
						ReactionTemp:       _input.ReactionTemp,
						ReactionTime:       _input.ReactionTime,
						OutputReactionName: fmt.Sprintf("%s%d", _input.OutputReactionName, counter),
						OutputLocation:     wellpositionarray[counter],
						OutputPlateNum:     platenum,
						LHPolicyName:       names[i],

						Parts:     _input.PartsArray[j],
						MasterMix: _input.Mastermix,
						Water:     _input.Water,
						OutPlate:  _input.OutPlate},
					)
					_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)

					// get annotation info
					doerun := names[i]

					partvols := make([]string, 0)

					for _, volume := range PartVolsArray {
						partvols = append(partvols, volume.ToString())
					} //strconv.Itoa(wutil.RoundInt(number))+"ul"

					solutionnames := _input.PartNamesArray[j]

					description := strings.Join(partvols, ":") + "_" + strings.Join(solutionnames, ":") + "_replicate" + strconv.Itoa(l+1) + "_platenum" + strconv.Itoa(platenum)
					//setpoints := volume+"_"+solutionname+"_replicate"+strconv.Itoa(j+1)+"_platenum"+strconv.Itoa(platenum)

					// add run to well position lookup table
					_output.Runtowelllocationmap[doerun+"_"+description] = wellpositionarray[counter]

					// replace responses with relevant ones
					runs[i] = doe.DeleteAllResponses(runs[i])

					runs[i] = doe.AddNewResponseField(runs[i], "Colonies")

					// add additional info for each run
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "Location", wellpositionarray[counter])

					// add run order:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "runorder", counter)

					// add setpoint printout to double check correct match up:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "doerun", doerun)

					// add description:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "description", description)

					// add dna Mass set point
					// add run order:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "DNA mass per Part", _input.DNAMassesPerReaction[k].ToString())

					newRuns = append(newRuns, runs[i])
					counter++

				}
			}
		}

		// export overall DOE design file showing all well locations for all conditions
		_ = doe.JMPXLSXFilefromRuns(newRuns, _input.OutputDesignFilename)

		_output.Runs = newRuns
		_output.NumberofReactions = len(_output.Runs) //counter //len(Reactions)
	}

}

func _ScreenAssemblyConditions_concAnalysis(_ctx context.Context, _input *ScreenAssemblyConditions_concInput, _output *ScreenAssemblyConditions_concOutput) {
}

func _ScreenAssemblyConditions_concValidation(_ctx context.Context, _input *ScreenAssemblyConditions_concInput, _output *ScreenAssemblyConditions_concOutput) {
}
func _ScreenAssemblyConditions_concRun(_ctx context.Context, input *ScreenAssemblyConditions_concInput) *ScreenAssemblyConditions_concOutput {
	output := &ScreenAssemblyConditions_concOutput{}
	_ScreenAssemblyConditions_concSetup(_ctx, input)
	_ScreenAssemblyConditions_concSteps(_ctx, input, output)
	_ScreenAssemblyConditions_concAnalysis(_ctx, input, output)
	_ScreenAssemblyConditions_concValidation(_ctx, input, output)
	return output
}

func ScreenAssemblyConditions_concRunSteps(_ctx context.Context, input *ScreenAssemblyConditions_concInput) *ScreenAssemblyConditions_concSOutput {
	soutput := &ScreenAssemblyConditions_concSOutput{}
	output := _ScreenAssemblyConditions_concRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ScreenAssemblyConditions_concNew() interface{} {
	return &ScreenAssemblyConditions_concElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ScreenAssemblyConditions_concInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ScreenAssemblyConditions_concRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ScreenAssemblyConditions_concInput{},
			Out: &ScreenAssemblyConditions_concOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ScreenAssemblyConditions_concElement struct {
	inject.CheckedRunner
}

type ScreenAssemblyConditions_concInput struct {
	DNAMassesPerReaction []wunit.Mass
	DXORJMP              string
	LHDOEFile            string
	Mastermix            *wtype.LHComponent
	MastermixVolume      wunit.Volume
	OutPlate             *wtype.LHPlate
	OutputDesignFilename string
	OutputReactionName   string
	PartConcArray        [][]wunit.Concentration
	PartNamesArray       [][]string
	PartsArray           [][]*wtype.LHComponent
	ReactionTemp         wunit.Temperature
	ReactionTime         wunit.Time
	ReactionVolume       wunit.Volume
	Replicates           int
	Water                *wtype.LHComponent
}

type ScreenAssemblyConditions_concOutput struct {
	NumberofReactions    int
	Reactions            []*wtype.LHComponent
	Runs                 []doe.Run
	Runtowelllocationmap map[string]string
}

type ScreenAssemblyConditions_concSOutput struct {
	Data struct {
		NumberofReactions    int
		Runs                 []doe.Run
		Runtowelllocationmap map[string]string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "ScreenAssemblyConditions_conc",
		Constructor: ScreenAssemblyConditions_concNew,
		Desc: ComponentDesc{
			Desc: "Assemble multiple assemblies using TypeIIs construct assembly\n",
			Path: "antha/component/an/Liquid_handling/TypeIIsAssembly/ScreenAssemblyConditions/ScreenAssemblyConditions_conc.an",
			Params: []ParamDesc{
				{Name: "DNAMassesPerReaction", Desc: "variable\n", Kind: "Parameters"},
				{Name: "DXORJMP", Desc: "", Kind: "Parameters"},
				{Name: "LHDOEFile", Desc: "file containing design for liquid handling DOE\n", Kind: "Parameters"},
				{Name: "Mastermix", Desc: "fixed\n", Kind: "Inputs"},
				{Name: "MastermixVolume", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "Output plate\n", Kind: "Inputs"},
				{Name: "OutputDesignFilename", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "Prefix for reaction names\n", Kind: "Parameters"},
				{Name: "PartConcArray", Desc: "variable but coupled\n\nVolumes corresponding to input parts // coupled with PartsArray and should be equal in length\n", Kind: "Parameters"},
				{Name: "PartNamesArray", Desc: "Names corresonding to input parts\n", Kind: "Parameters"},
				{Name: "PartsArray", Desc: "Variable\n\nInput parts, one per assembly\n", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "Reaction temperature\n", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "Reaction time\n", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "Reaction volume\n", Kind: "Parameters"},
				{Name: "Replicates", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "NumberofReactions", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "List of assembled parts\n", Kind: "Outputs"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "", Kind: "Data"},
			},
		},
	})
}
