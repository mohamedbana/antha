// Assemble multiple assemblies using TypeIIs construct assembly
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

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

// variable
// Volumes corresponding to input parts // coupled with PartsArray and should be equal in length
// Names corresonding to input parts

//fixed

// Output plate

// Variable
// Input parts, one per assembly

// List of assembled parts

func _ScreenAssemblyConditions_specificSetup(_ctx context.Context, _input *ScreenAssemblyConditions_specificInput) {
}

func _ScreenAssemblyConditions_specificSteps(_ctx context.Context, _input *ScreenAssemblyConditions_specificInput, _output *ScreenAssemblyConditions_specificOutput) {

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

		for k := range _input.PartVolsArray {

			for j := range _input.PartNamesArray {

				for i := 0; i < len(runs); i++ {

					if counter == (_input.OutPlate.WlsX * _input.OutPlate.WlsY) {
						fmt.Println("plate full, counter = ", counter)
						platenum++
						counter = 0
					}

					fmt.Println("counter:", counter)
					fmt.Println("WellPositionarray", wellpositionarray, "OutPlate.WlsX", _input.OutPlate.WlsX)

					result := TypeIISConstructAssemblyMMXRunSteps(_ctx, &TypeIISConstructAssemblyMMXInput{ReactionVolume: _input.ReactionVolume,
						PartVols:           _input.PartVolsArray[k],
						PartNames:          _input.PartNamesArray[j],
						ReactionTemp:       _input.ReactionTemp,
						ReactionTime:       _input.ReactionTime,
						OutputReactionName: fmt.Sprintf("%s%d", _input.OutputReactionName, counter),
						OutputLocation:     wellpositionarray[counter],
						OutputPlateNum:     platenum,
						LHPolicyName:       names[i],

						Parts:     _input.PartsArray[j],
						MasterMix: _input.Mastermix,
						OutPlate:  _input.OutPlate},
					)
					_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)

					// get annotation info
					doerun := names[i]

					partvols := make([]string, 0)

					for _, volume := range _input.PartVolsArray[k] {
						partvols = append(partvols, volume.ToString())
					} //strconv.Itoa(wutil.RoundInt(number))+"ul"

					solutionnames := _input.PartNamesArray[j]

					description := strings.Join(partvols, ":") + "_" + strings.Join(solutionnames, ":") + "_replicate" + strconv.Itoa(l+1) + "_platenum" + strconv.Itoa(platenum)
					//setpoints := volume+"_"+solutionname+"_replicate"+strconv.Itoa(j+1)+"_platenum"+strconv.Itoa(platenum)

					// add run to well position lookup table
					_output.Runtowelllocationmap[doerun+"_"+description] = wellpositionarray[counter]

					// replace responses with relevant ones
					runs[i] = doe.DeleteAllResponses(runs[i])

					runs[i] = doe.AddNewResponseField(runs[i], "Number of Colonies")

					// add additional info for each run
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "Location", wellpositionarray[counter])

					// add run order:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "runorder", counter)

					// add setpoint printout to double check correct match up:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "doerun", doerun)

					// add description:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "description", description)
					//runs[i].AddAdditionalValue("Replicate", strconv.Itoa(j+1))
					//runs[i].AddAdditionalValue("Solution name", TestSols[k].CName)
					//runs[i].AddAdditionalValue("Volume", strconv.Itoa(wutil.RoundInt(TestSolVolumes[l].RawValue()))+"ul)

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

func _ScreenAssemblyConditions_specificAnalysis(_ctx context.Context, _input *ScreenAssemblyConditions_specificInput, _output *ScreenAssemblyConditions_specificOutput) {
}

func _ScreenAssemblyConditions_specificValidation(_ctx context.Context, _input *ScreenAssemblyConditions_specificInput, _output *ScreenAssemblyConditions_specificOutput) {
}
func _ScreenAssemblyConditions_specificRun(_ctx context.Context, input *ScreenAssemblyConditions_specificInput) *ScreenAssemblyConditions_specificOutput {
	output := &ScreenAssemblyConditions_specificOutput{}
	_ScreenAssemblyConditions_specificSetup(_ctx, input)
	_ScreenAssemblyConditions_specificSteps(_ctx, input, output)
	_ScreenAssemblyConditions_specificAnalysis(_ctx, input, output)
	_ScreenAssemblyConditions_specificValidation(_ctx, input, output)
	return output
}

func ScreenAssemblyConditions_specificRunSteps(_ctx context.Context, input *ScreenAssemblyConditions_specificInput) *ScreenAssemblyConditions_specificSOutput {
	soutput := &ScreenAssemblyConditions_specificSOutput{}
	output := _ScreenAssemblyConditions_specificRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ScreenAssemblyConditions_specificNew() interface{} {
	return &ScreenAssemblyConditions_specificElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ScreenAssemblyConditions_specificInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ScreenAssemblyConditions_specificRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ScreenAssemblyConditions_specificInput{},
			Out: &ScreenAssemblyConditions_specificOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ScreenAssemblyConditions_specificElement struct {
	inject.CheckedRunner
}

type ScreenAssemblyConditions_specificInput struct {
	DXORJMP              string
	LHDOEFile            string
	Mastermix            *wtype.LHComponent
	OutPlate             *wtype.LHPlate
	OutputDesignFilename string
	OutputReactionName   string
	PartNamesArray       [][]string
	PartVolsArray        [][]wunit.Volume
	PartsArray           [][]*wtype.LHComponent
	ReactionTemp         wunit.Temperature
	ReactionTime         wunit.Time
	ReactionVolume       wunit.Volume
	Replicates           int
}

type ScreenAssemblyConditions_specificOutput struct {
	NumberofReactions    int
	Reactions            []*wtype.LHComponent
	Runs                 []doe.Run
	Runtowelllocationmap map[string]string
}

type ScreenAssemblyConditions_specificSOutput struct {
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
	addComponent(Component{Name: "ScreenAssemblyConditions_specific",
		Constructor: ScreenAssemblyConditions_specificNew,
		Desc: ComponentDesc{
			Desc: "Assemble multiple assemblies using TypeIIs construct assembly\n",
			Path: "antha/component/an/Liquid_handling/TypeIIsAssembly/ScreenAssemblyConditions/ScreenAssemblyConditions_specific.an",
			Params: []ParamDesc{
				{Name: "DXORJMP", Desc: "", Kind: "Parameters"},
				{Name: "LHDOEFile", Desc: "file containing design for liquid handling DOE\n", Kind: "Parameters"},
				{Name: "Mastermix", Desc: "fixed\n", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "Output plate\n", Kind: "Inputs"},
				{Name: "OutputDesignFilename", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "Prefix for reaction names\n", Kind: "Parameters"},
				{Name: "PartNamesArray", Desc: "Names corresonding to input parts\n", Kind: "Parameters"},
				{Name: "PartVolsArray", Desc: "variable\n\nVolumes corresponding to input parts // coupled with PartsArray and should be equal in length\n", Kind: "Parameters"},
				{Name: "PartsArray", Desc: "Variable\n\nInput parts, one per assembly\n", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "Reaction temperature\n", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "Reaction time\n", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "Reaction volume\n", Kind: "Parameters"},
				{Name: "Replicates", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReactions", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "List of assembled parts\n", Kind: "Outputs"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "", Kind: "Data"},
			},
		},
	})
}
