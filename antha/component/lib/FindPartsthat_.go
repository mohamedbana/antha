// example protocol which allows a primitive method for searching the igem registry
// for parts with specified functions or a specified status (e.g. A = available or "Works", or results != none)
// see the igem package ("github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem")
// and igem website for more details about how to make the most of this http://parts.igem.org/Registry_API

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	//	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

// e.g. rbs, reporter
// e.g. strong, arsenic, fluorescent, alkane, logic gate

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _FindPartsthatRequirements() {

}

// Conditions to run on startup
func _FindPartsthatSetup(_ctx context.Context, _input *FindPartsthatInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _FindPartsthatSteps(_ctx context.Context, _input *FindPartsthatInput, _output *FindPartsthatOutput) {

	BackupParts := make([]string, 0)
	status := ""
	joinedstatus := make([]string, 0)
	// Look up parts from registry according to properties (this will take a couple of minutes the first time)

	parts := make([][]string, 0)
	_output.PartMap = make(map[string][]string)
	_output.BiobrickDescriptions = make(map[string]string)
	subparts := make([]string, 0)
	var highestrating int

	partstatus := ""

	if _input.OnlyreturnAvailableParts {
		partstatus = "A"
	}

	// first we'll parse the igem registry based on the short description contained in the fasta header for each part sequence
	for _, desc := range _input.Parttypes {

		subparts = igem.FilterRegistry(desc, []string{desc, partstatus})
		status = text.Print(desc+" :", subparts)
		joinedstatus = append(joinedstatus, status)
		parts = append(parts, subparts)
		_output.PartMap[desc] = subparts
	}

	othercriteria := ""
	if _input.OnlyreturnWorkingparts {
		othercriteria = "WORKS"
	}

	var i int

	for desc, subparts := range _output.PartMap {

		partdetails := igem.LookUp(subparts)

		// now we can get detailed information of all of those records to interrogate further
		// this can be slow if there are many parts to check (~2 seconds per block of 14 parts)

		for _, subpart := range subparts {

			if strings.Contains(strings.ToUpper(partdetails.Description(subpart)), strings.ToUpper(_input.Partdescriptions[i])) &&
				strings.Contains(partdetails.Results(subpart), othercriteria) {
				BackupParts = append(BackupParts, subpart)
				_output.BiobrickDescriptions[subpart] = partdetails.Description(subpart)

				rating, err := strconv.Atoi(partdetails.Rating(subpart))

				if err == nil && rating > highestrating {
					_output.HighestRatedMatch = subpart

					seq := partdetails.Sequence(_output.HighestRatedMatch)

					_output.HighestRatedMatchDNASequence = wtype.MakeLinearDNASequence(_output.HighestRatedMatch, seq)
				}
			}

			delete(_output.PartMap, desc)
			_output.PartMap[desc] = BackupParts

			_output.FulllistBackupParts = append(_output.FulllistBackupParts, BackupParts)
		}
		i = i + 1
	}

	_output.FulllistBackupParts = parts
	_output.Status = strings.Join(joinedstatus, " ; ")

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			"Warnings:", _output.Warnings.Error(),
			"Back up parts found (Reported to work!)", _output.FulllistBackupParts,
		)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _FindPartsthatAnalysis(_ctx context.Context, _input *FindPartsthatInput, _output *FindPartsthatOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _FindPartsthatValidation(_ctx context.Context, _input *FindPartsthatInput, _output *FindPartsthatOutput) {
}
func _FindPartsthatRun(_ctx context.Context, input *FindPartsthatInput) *FindPartsthatOutput {
	output := &FindPartsthatOutput{}
	_FindPartsthatSetup(_ctx, input)
	_FindPartsthatSteps(_ctx, input, output)
	_FindPartsthatAnalysis(_ctx, input, output)
	_FindPartsthatValidation(_ctx, input, output)
	return output
}

func FindPartsthatRunSteps(_ctx context.Context, input *FindPartsthatInput) *FindPartsthatSOutput {
	soutput := &FindPartsthatSOutput{}
	output := _FindPartsthatRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func FindPartsthatNew() interface{} {
	return &FindPartsthatElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &FindPartsthatInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _FindPartsthatRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &FindPartsthatInput{},
			Out: &FindPartsthatOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type FindPartsthatElement struct {
	inject.CheckedRunner
}

type FindPartsthatInput struct {
	OnlyreturnAvailableParts bool
	OnlyreturnWorkingparts   bool
	Partdescriptions         []string
	Parttypes                []string
}

type FindPartsthatOutput struct {
	BiobrickDescriptions         map[string]string
	FulllistBackupParts          [][]string
	HighestRatedMatch            string
	HighestRatedMatchDNASequence wtype.DNASequence
	PartMap                      map[string][]string
	Status                       string
	Warnings                     error
}

type FindPartsthatSOutput struct {
	Data struct {
		BiobrickDescriptions         map[string]string
		FulllistBackupParts          [][]string
		HighestRatedMatch            string
		HighestRatedMatchDNASequence wtype.DNASequence
		PartMap                      map[string][]string
		Status                       string
		Warnings                     error
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "FindPartsthat",
		Constructor: FindPartsthatNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/FindPartsthat/Findpartsthat.an",
			Params: []ParamDesc{
				{Name: "OnlyreturnAvailableParts", Desc: "", Kind: "Parameters"},
				{Name: "OnlyreturnWorkingparts", Desc: "", Kind: "Parameters"},
				{Name: "Partdescriptions", Desc: "e.g. strong, arsenic, fluorescent, alkane, logic gate\n", Kind: "Parameters"},
				{Name: "Parttypes", Desc: "e.g. rbs, reporter\n", Kind: "Parameters"},
				{Name: "BiobrickDescriptions", Desc: "", Kind: "Data"},
				{Name: "FulllistBackupParts", Desc: "", Kind: "Data"},
				{Name: "HighestRatedMatch", Desc: "", Kind: "Data"},
				{Name: "HighestRatedMatchDNASequence", Desc: "", Kind: "Data"},
				{Name: "PartMap", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	})
}
