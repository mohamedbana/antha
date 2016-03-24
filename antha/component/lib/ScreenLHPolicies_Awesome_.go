package lib

import (
	"fmt"
	antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"path/filepath"
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

//map[string]string

//NeatSamplewells []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ScreenLHPolicies_AwesomeRequirements() {
}

// Conditions to run on startup
func _ScreenLHPolicies_AwesomeSetup(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ScreenLHPolicies_AwesomeSteps(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput, _output *ScreenLHPolicies_AwesomeOutput) {

	var rotate = false

	chosencolourpalette := image.AvailablePalettes["Palette1"]
	positiontocolourmap, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, rotate)

	_output.Runtowelllocationmap = make([]string, 0)
	perconditionuntowelllocationmap := make([]string, 0)

	//Runtowelllocationmap = make(map[string]string)

	// work out well coordinates for any plate
	wellpositionarray := make([]string, 0)

	for location, colour := range positiontocolourmap {
		R, G, B, A := colour.RGBA()

		if uint8(R) == 242 && uint8(G) == 243 && uint8(B) == 242 && uint8(A) == 255 {
			continue
		} else {
			wellpositionarray = append(wellpositionarray, location)
		}
	}

	/*
		//alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
			"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
			"Y", "Z", "AA", "BB", "CC", "DD", "EE", "FF"}
		//k := 0
		for j := 0; j < OutPlate.WlsY; j++ {
			for i := 0; i < OutPlate.WlsX; i++ { //countingfrom1iswhatmakesushuman := j + 1
				//k = k + 1
				wellposition := string(alphabet[j]) + strconv.Itoa(i+1)
				//fmt.Println(wellposition, k)
				wellpositionarray = append(wellpositionarray, wellposition)
			}

		}
	*/
	reactions := make([]*wtype.LHComponent, 0)

	//policies, names := liquidhandling.PolicyMaker(liquidhandling.Allpairs, "DOE_run",false)

	//intfactors := []string{"Pre_MIX","POST_MIX"}
	policies, names, err := liquidhandling.PolicyMakerfromDesign("ScreenLHPolicyDOE2.xlsx", "DOE_run")
	if err != nil {
		panic(err)
	}

	counter := 0
	for l := 0; l < len(_input.TestSolVolumes); l++ {
		for k := 0; k < len(_input.TestSols); k++ {
			for j := 0; j < _input.NumberofReplicates; j++ {
				for i := 0; i < len(policies); i++ {

					eachreaction := make([]*wtype.LHComponent, 0)

					// keep default policy for diluent
					//Diluent.Type = names[i]
					//fmt.Println(Diluent.Type)

					bufferSample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
					eachreaction = append(eachreaction, bufferSample)
					testSample := mixer.Sample(_input.TestSols[k], _input.TestSolVolumes[l])

					_input.TestSols[k].Type = wtype.LiquidTypeFromString(names[i])

					eachreaction = append(eachreaction, testSample)
					reaction := execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], 0, eachreaction...)
					fmt.Println("where am I?", wellpositionarray[counter])
					_output.Runtowelllocationmap = append(_output.Runtowelllocationmap, wtype.LiquidTypeName(_input.TestSols[k].Type)+":"+wellpositionarray[counter])
					perconditionuntowelllocationmap = append(perconditionuntowelllocationmap, wtype.LiquidTypeName(_input.TestSols[k].Type)+":"+wellpositionarray[counter])
					//Runtowelllocationmap[Diluent.Type]= wellpositionarray[counter]
					reactions = append(reactions, reaction)
					counter = counter + 1
				}

				outputsandwich := strconv.Itoa(wutil.RoundInt(_input.TestSolVolumes[l].RawValue())) + _input.TestSols[k].CName + strconv.Itoa(j+1)

				outputfilename := filepath.Join(antha.Dirpath(), "DOE2"+"_"+outputsandwich+".xlsx")
				_output.Errors = append(_output.Errors, doe.AddWelllocations(filepath.Join(antha.Dirpath(), "ScreenLHPolicyDOE2.xlsx"), 0, perconditionuntowelllocationmap, "DOE_run", outputfilename, []string{"Volume", "Solution", "Replicate"}, []interface{}{_input.TestSolVolumes[l].ToString(), _input.TestSols[k].CName, string(j)}))

				// other things to add to check for covariance
				// order in which wells were pippetted
				// plate ID
				// row
				// column
				// ambient temp

				// empty
				perconditionuntowelllocationmap = make([]string, 0)
			}
		}
	}

	// add blanks

	_output.Blankwells = make([]string, 0)

	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "AA", "BB", "CC", "DD", "EE", "FF"}

	for m := 0; m < _input.NumberofBlanks; m++ {
		//eachreaction := make([]*wtype.LHComponent, 0)

		// use defualt policy for blank

		bufferSample := mixer.Sample(_input.Diluent, _input.TotalVolume)
		//eachreaction = append(eachreaction,bufferSample)

		// add blanks to last column of plate
		well := alphabet[_input.OutPlate.WlsY-1-m] + strconv.Itoa(_input.OutPlate.WlsX)
		fmt.Println("blankwell", well)
		reaction := execute.MixTo(_ctx, _input.OutPlate.Type, well, 0, bufferSample)
		//fmt.Println("where am I?",wellpositionarray[counter])
		//Runtowelllocationmap= append(Runtowelllocationmap,"Blank"+ strconv.Itoa(m+1) +":" + well)
		_output.Blankwells = append(_output.Blankwells, well)

		reactions = append(reactions, reaction)
		counter = counter + 1

	}

	_output.Reactions = reactions
	_output.Runcount = len(_output.Reactions)
	_output.Pixelcount = len(wellpositionarray)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ScreenLHPolicies_AwesomeAnalysis(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput, _output *ScreenLHPolicies_AwesomeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ScreenLHPolicies_AwesomeValidation(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput, _output *ScreenLHPolicies_AwesomeOutput) {
}
func _ScreenLHPolicies_AwesomeRun(_ctx context.Context, input *ScreenLHPolicies_AwesomeInput) *ScreenLHPolicies_AwesomeOutput {
	output := &ScreenLHPolicies_AwesomeOutput{}
	_ScreenLHPolicies_AwesomeSetup(_ctx, input)
	_ScreenLHPolicies_AwesomeSteps(_ctx, input, output)
	_ScreenLHPolicies_AwesomeAnalysis(_ctx, input, output)
	_ScreenLHPolicies_AwesomeValidation(_ctx, input, output)
	return output
}

func ScreenLHPolicies_AwesomeRunSteps(_ctx context.Context, input *ScreenLHPolicies_AwesomeInput) *ScreenLHPolicies_AwesomeSOutput {
	soutput := &ScreenLHPolicies_AwesomeSOutput{}
	output := _ScreenLHPolicies_AwesomeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ScreenLHPolicies_AwesomeNew() interface{} {
	return &ScreenLHPolicies_AwesomeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ScreenLHPolicies_AwesomeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ScreenLHPolicies_AwesomeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ScreenLHPolicies_AwesomeInput{},
			Out: &ScreenLHPolicies_AwesomeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ScreenLHPolicies_AwesomeElement struct {
	inject.CheckedRunner
}

type ScreenLHPolicies_AwesomeInput struct {
	Diluent            *wtype.LHComponent
	Imagefilename      string
	NumberofBlanks     int
	NumberofReplicates int
	OutPlate           *wtype.LHPlate
	TestSolVolumes     []wunit.Volume
	TestSols           []*wtype.LHComponent
	TotalVolume        wunit.Volume
}

type ScreenLHPolicies_AwesomeOutput struct {
	Blankwells           []string
	Errors               []error
	Pixelcount           int
	Reactions            []*wtype.LHComponent
	Runcount             int
	Runtowelllocationmap []string
}

type ScreenLHPolicies_AwesomeSOutput struct {
	Data struct {
		Blankwells           []string
		Errors               []error
		Pixelcount           int
		Runcount             int
		Runtowelllocationmap []string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "ScreenLHPolicies_Awesome",
		Constructor: ScreenLHPolicies_AwesomeNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/FindbestLHPolicy/ScreenLHPolicies_Awesome.an",
			Params: []ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "NumberofBlanks", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "TestSolVolumes", Desc: "", Kind: "Parameters"},
				{Name: "TestSols", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Blankwells", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Pixelcount", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Runcount", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "map[string]string\n", Kind: "Data"},
			},
		},
	})
}
