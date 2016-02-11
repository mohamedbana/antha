// example protocol which allows a primitive method for searching the igem registry
// for parts with specified functions or a specified status (e.g. A = available or "Works", or results != none)
// see the igem package ("github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem")
// and igem website for more details about how to make the most of this http://parts.igem.org/Registry_API

package lib

import (
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	"fmt"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	//	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strings"
)

// Input parameters for this protocol (data)

//Constructname 				string
// e.g. promoter
// e.g. arsenic, reporter, alkane, logic gate

//RestrictionsitetoAvoid		[]string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

//Partsfound	[]wtype.DNASequence // map[string]wtype.DNASequence
//map[string][]string

// Input Requirement specification
func _FindPartsthatRequirements() {

}

// Conditions to run on startup
func _FindPartsthatSetup(_ctx context.Context, _input *FindPartsthatInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _FindPartsthatSteps(_ctx context.Context, _input *FindPartsthatInput, _output *FindPartsthatOutput) {
	//var msg string
	// set warnings reported back to user to none initially
	//	warnings := make([]string,0)
	BackupParts := make([]string, 0)
	status := ""
	joinedstatus := make([]string, 0)
	// Look up parts from registry according to properties (this will take a couple of minutes the first time)

	parts := make([][]string, 0)
	subparts := make([]string, 0)

	partstatus := ""

	if _input.OnlyreturnAvailableParts {
		partstatus = "A"
	}

	// first we'll parse the igem registry based on the short description contained in the fasta header for each part sequence
	for _, desc := range _input.Parttypes {

		subparts = igem.FilterRegistry([]string{desc, partstatus})
		status = text.Print(desc+" :", subparts)
		joinedstatus = append(joinedstatus, status)
		parts = append(parts, subparts)
	}

	othercriteria := ""
	if _input.OnlyreturnWorkingparts {
		othercriteria = "WORKS"
	}

	for i, subparts := range parts {

		partdetails := igem.LookUp(subparts)
		// now we can get detailed information of all of those records to interrogate further
		// this can be slow if there are many parts to check (~2 seconds per block of 14 parts)

		for _, subpart := range subparts {

			if strings.Contains(partdetails.Description(subpart), _input.Partdescriptions[i]) &&
				strings.Contains(partdetails.Results(subpart), othercriteria) {
				BackupParts = append(BackupParts, subpart)

			}
			_output.FulllistBackupParts = append(_output.FulllistBackupParts, BackupParts)
		}
	}
	/*
		if len(warnings) != 0 {
		Warnings = fmt.Errorf(strings.Join(warnings,";"))
		}else{Warnings = nil}
	*/

	_output.FulllistBackupParts = parts
	_output.Status = strings.Join(joinedstatus, " ; ")

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			"Warnings:", _output.Warnings.Error(),
			"Back up parts found (Reported to work!)", _input.Parts,
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
	Parts                    [][]string
	Parttypes                []string
}

type FindPartsthatOutput struct {
	FulllistBackupParts [][]string
	Status              string
	Warnings            error
}

type FindPartsthatSOutput struct {
	Data struct {
		FulllistBackupParts [][]string
		Status              string
		Warnings            error
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
				{Name: "Partdescriptions", Desc: "e.g. arsenic, reporter, alkane, logic gate\n", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Parameters"},
				{Name: "Parttypes", Desc: "Constructname \t\t\t\tstring\n\ne.g. promoter\n", Kind: "Parameters"},
				{Name: "FulllistBackupParts", Desc: "Partsfound\t[]wtype.DNASequence // map[string]wtype.DNASequence\n\nmap[string][]string\n", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	})
}
