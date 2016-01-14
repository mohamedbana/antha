// example protocol which allows a primitive method for searching the igem registry
// for parts with specified functions or a specified status (e.g. A = available or "Works", or results != none)
// see the igem package ("github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem")
// and igem website for more details about how to make the most of this http://parts.igem.org/Registry_API

package FindPartsthat

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
func _requirements() {

}

// Conditions to run on startup
func _setup(_ctx context.Context, _input *Input_) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
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
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _validation(_ctx context.Context, _input *Input_, _output *Output_) {
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
	OnlyreturnAvailableParts bool
	OnlyreturnWorkingparts   bool
	Partdescriptions         []string
	Parts                    [][]string
	Parttypes                []string
}

type Output_ struct {
	FulllistBackupParts [][]string
	Status              string
	Warnings            error
}
