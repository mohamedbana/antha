package lib

import (
	"fmt"
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
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

//[]string //map[string]string

//NeatSamplewells []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AccuracyTest_2Requirements() {
}

// Conditions to run on startup
func _AccuracyTest_2Setup(_ctx context.Context, _input *AccuracyTest_2Input) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AccuracyTest_2Steps(_ctx context.Context, _input *AccuracyTest_2Input, _output *AccuracyTest_2Output) {

	// declare some global variables for use later
	var rotate = false
	var autorotate = true
	var wellpositionarray = make([]string, 0)
	var alphabet = wutil.MakeAlphabetArray()
	_output.Runtowelllocationmap = make(map[string]string)
	_output.Blankwells = make([]string, 0)
	counter := 0

	var platenum = 1
	var runs = make([]doe.Run, 1)
	var newruns = make([]doe.Run, 0)
	// work out plate layout based on picture or just in order

	if _input.Printasimage {
		chosencolourpalette := image.AvailablePalettes["Palette1"]
		positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, rotate, autorotate)

		//Runtowelllocationmap = make([]string,0)

		for location, colour := range positiontocolourmap {
			R, G, B, A := colour.RGBA()

			if uint8(R) == 242 && uint8(G) == 243 && uint8(B) == 242 && uint8(A) == 255 {
				continue
			} else {
				wellpositionarray = append(wellpositionarray, location)
			}
		}

	} else {

		wellpositionarray = _input.OutPlate.AllWellPositions()

	}
	reactions := make([]*wtype.LHComponent, 0)

	// use first policy as reference to ensure consistent range through map values
	referencepolicy, _ := liquidhandling.GetPolicyByName(_input.LHPolicy)

	referencekeys := make([]string, 0)
	for key := range referencepolicy {
		referencekeys = append(referencekeys, key)
	}

	for l := 0; l < len(_input.TestSolVolumes); l++ {
		for k := 0; k < len(_input.TestSols); k++ {
			for j := 0; j < _input.NumberofReplicates; j++ {
				for i := 0; i < len(runs); i++ {

					var run doe.Run

					if counter == ((_input.OutPlate.WlsX * _input.OutPlate.WlsY) + _input.NumberofBlanks) {
						fmt.Println("plate full, counter = ", counter)
						platenum++
						counter = 0
					}

					var eachreaction []*wtype.LHComponent
					var solution *wtype.LHComponent

					if _input.PipetteOnebyOne {
						eachreaction = make([]*wtype.LHComponent, 0)
					}
					// keep default policy for diluent

					// diluent first

					// change lhpolicy if desired
					if _input.UseLHPolicyDoeforDiluent {
						_input.Diluent.Type = wtype.LiquidTypeFromString(_input.LHPolicy)
					}

					bufferSample := mixer.Sample(_input.Diluent, wunit.NewVolume(_input.TotalVolume.RawValue()-_input.TestSolVolumes[l].RawValue(), _input.TotalVolume.Unit().PrefixedSymbol())) //SampleForTotalVolume(Diluent, TotalVolume)

					if _input.PipetteOnebyOne {
						eachreaction = append(eachreaction, bufferSample)
					} else {

						solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, bufferSample)
					}
					// now test sample

					// change liquid class
					if _input.UseLiquidPolicyForTestSolutions && _input.LHPolicy != "" {
						_input.TestSols[k].Type = wtype.LiquidTypeFromString(_input.LHPolicy)
					}

					if _input.TestSolVolumes[l].RawValue() > 0.0 {
						//sample
						testSample := mixer.Sample(_input.TestSols[k], _input.TestSolVolumes[l])

						if _input.PipetteOnebyOne {
							eachreaction = append(eachreaction, testSample)
							solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, eachreaction...)
						} else {
							// pipette out
							//solution = MixTo(OutPlate.Type,wellpositionarray[counter],platenum, testSample)
							solution = execute.Mix(_ctx, solution, testSample)
						}

					}

					// get annotation info
					doerun := wtype.LiquidTypeName(_input.TestSols[k].Type)

					volume := _input.TestSolVolumes[l].ToString() //strconv.Itoa(wutil.RoundInt(number))+"ul"

					solutionname := _input.TestSols[k].CName

					// add Solution Name
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Solution", solutionname)

					// add Volume
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Volume", volume)

					// add Replicate
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Replicate", strconv.Itoa(j+1))

					// full description
					description := volume + "_" + solutionname + "_replicate" + strconv.Itoa(j+1) + "_platenum" + strconv.Itoa(platenum)
					//setpoints := volume+"_"+solutionname+"_replicate"+strconv.Itoa(j+1)+"_platenum"+strconv.Itoa(platenum)

					// add run to well position lookup table
					_output.Runtowelllocationmap[doerun+"_"+description] = wellpositionarray[counter]

					// add additional info for each run
					fmt.Println("len(runs)", len(runs), "counter", counter, "len(wellpositionarray)", len(wellpositionarray))

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "PlateNumber", strconv.Itoa(platenum))
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Location", wellpositionarray[counter])

					// add run order:
					//runs[counter] = doe.AddAdditionalHeaderandValue(runs[counter],"Additional","runorder", counter)

					// add setpoint printout to double check correct match up:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy", doerun)

					// add plate info:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate Type", _input.OutPlate.Type)

					// add plate ZStart:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellZStart", _input.OutPlate.WellZStart)

					// add plate Height:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate Height", _input.OutPlate.Height)

					// other plate offsets:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellXOffset", _input.OutPlate.WellXOffset)

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellYOffset", _input.OutPlate.WellYOffset)

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellXStart", _input.OutPlate.WellXStart)

					run = doe.AddAdditionalHeaderandValue(run, "Additional", "Plate WellYStart", _input.OutPlate.WellYStart)

					// add LHPolicy setpoint printout to double check correct match up:
					run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy", doerun)

					// print out LHPolicy info
					policy, _ := liquidhandling.GetPolicyByName(doerun)

					for _, key := range referencekeys {
						run = doe.AddAdditionalHeaderandValue(run, "Additional", "LHPolicy"+"_"+key, policy[key])
					}

					reactions = append(reactions, solution)
					newruns = append(newruns, run)

					counter = counter + 1

				}

			}
		}
	}

	// export overall DOE design file showing all well locations for all conditions
	_ = doe.XLSXFileFromRuns(newruns, _input.OutputFilename, _input.DXORJMP)

	// add blanks after

	for n := 0; n < platenum; n++ {
		for m := 0; m < _input.NumberofBlanks; m++ {

			// use defualt policy for blank

			bufferSample := mixer.Sample(_input.Diluent, _input.TotalVolume)
			//eachreaction = append(eachreaction,bufferSample)

			// add blanks to last column of plate
			well := alphabet[_input.OutPlate.WlsY-1-m] + strconv.Itoa(_input.OutPlate.WlsX)

			reaction := execute.MixTo(_ctx, _input.OutPlate.Type, well, n+1, bufferSample)

			_output.Runtowelllocationmap["Blank"+strconv.Itoa(m+1)+" platenum"+strconv.Itoa(n+1)] = well

			_output.Blankwells = append(_output.Blankwells, well)

			reactions = append(reactions, reaction)
			counter = counter + 1

		}

	}

	_output.Reactions = reactions
	_output.Runcount = len(_output.Reactions)
	_output.Pixelcount = len(wellpositionarray)
	_output.Runs = runs
	_output.Wellpositionarray = wellpositionarray
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AccuracyTest_2Analysis(_ctx context.Context, _input *AccuracyTest_2Input, _output *AccuracyTest_2Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AccuracyTest_2Validation(_ctx context.Context, _input *AccuracyTest_2Input, _output *AccuracyTest_2Output) {
}
func _AccuracyTest_2Run(_ctx context.Context, input *AccuracyTest_2Input) *AccuracyTest_2Output {
	output := &AccuracyTest_2Output{}
	_AccuracyTest_2Setup(_ctx, input)
	_AccuracyTest_2Steps(_ctx, input, output)
	_AccuracyTest_2Analysis(_ctx, input, output)
	_AccuracyTest_2Validation(_ctx, input, output)
	return output
}

func AccuracyTest_2RunSteps(_ctx context.Context, input *AccuracyTest_2Input) *AccuracyTest_2SOutput {
	soutput := &AccuracyTest_2SOutput{}
	output := _AccuracyTest_2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AccuracyTest_2New() interface{} {
	return &AccuracyTest_2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AccuracyTest_2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AccuracyTest_2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AccuracyTest_2Input{},
			Out: &AccuracyTest_2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AccuracyTest_2Element struct {
	inject.CheckedRunner
}

type AccuracyTest_2Input struct {
	DXORJMP                         string
	Diluent                         *wtype.LHComponent
	Imagefilename                   string
	LHPolicy                        string
	NumberofBlanks                  int
	NumberofReplicates              int
	OutPlate                        *wtype.LHPlate
	OutputFilename                  string
	PipetteOnebyOne                 bool
	Printasimage                    bool
	TestSolVolumes                  []wunit.Volume
	TestSols                        []*wtype.LHComponent
	TotalVolume                     wunit.Volume
	UseLHPolicyDoeforDiluent        bool
	UseLiquidPolicyForTestSolutions bool
}

type AccuracyTest_2Output struct {
	Blankwells           []string
	Errors               []error
	Pixelcount           int
	Reactions            []*wtype.LHComponent
	Runcount             int
	Runs                 []doe.Run
	Runtowelllocationmap map[string]string
	Wellpositionarray    []string
}

type AccuracyTest_2SOutput struct {
	Data struct {
		Blankwells           []string
		Errors               []error
		Pixelcount           int
		Runcount             int
		Runs                 []doe.Run
		Runtowelllocationmap map[string]string
		Wellpositionarray    []string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "AccuracyTest_2",
		Constructor: AccuracyTest_2New,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Utility/AccuracyTest_2.an",
			Params: []ParamDesc{
				{Name: "DXORJMP", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicy", Desc: "", Kind: "Parameters"},
				{Name: "NumberofBlanks", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputFilename", Desc: "", Kind: "Parameters"},
				{Name: "PipetteOnebyOne", Desc: "", Kind: "Parameters"},
				{Name: "Printasimage", Desc: "", Kind: "Parameters"},
				{Name: "TestSolVolumes", Desc: "", Kind: "Parameters"},
				{Name: "TestSols", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "UseLHPolicyDoeforDiluent", Desc: "", Kind: "Parameters"},
				{Name: "UseLiquidPolicyForTestSolutions", Desc: "", Kind: "Parameters"},
				{Name: "Blankwells", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Pixelcount", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Runcount", Desc: "", Kind: "Data"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "[]string //map[string]string\n", Kind: "Data"},
				{Name: "Wellpositionarray", Desc: "", Kind: "Data"},
			},
		},
	})
}