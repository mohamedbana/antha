// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"strconv"
	"strings"
)

// Input parameters for this protocol (data)

// enter each as amino acid sequence

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func _Scarfree_designRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _Scarfree_designSetup(_ctx context.Context, _input *Scarfree_designInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _Scarfree_designSteps(_ctx context.Context, _input *Scarfree_designInput, _output *Scarfree_designOutput) {
	//var msg string
	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	var warning string
	var err error
	// make an empty array of DNA Sequences ready to fill
	partsinorder := make([]wtype.DNASequence, 0)

	var partDNA wtype.DNASequence
	var vectordata wtype.DNASequence

	for i, part := range _input.Seqsinorder {
		if strings.Contains(part, ".gb") && strings.Contains(part, "Feature:") {

			split := strings.SplitAfter(part, ".gb")
			file := split[0]

			split2 := strings.Split(split[1], ":")
			feature := split2[1]

			partDNA, _ = parser.GenbankFeaturetoDNASequence(file, feature)
		} else if strings.Contains(part, ".gb") {

			/*annotated,_ := parser.GenbanktoAnnotatedSeq(part)
			partDNA = annotated.DNASequence */

			partDNA, _ = parser.GenbanktoAnnotatedSeq(part)
		} else {

			if strings.Contains(part, "BBa_") {
				part = igem.GetSequence(part)
			}
			partDNA = wtype.MakeLinearDNASequence("Part "+strconv.Itoa(i), part)
		}
		partsinorder = append(partsinorder, partDNA)
	}

	// make vector into an antha type DNASequence

	if strings.Contains(_input.Vector, ".gb") {

		vectordata, _ = parser.GenbanktoAnnotatedSeq(_input.Vector)
		vectordata.Plasmid = true
	} else {

		if strings.Contains(_input.Vector, "BBa_") {
			_input.Vector = igem.GetSequence(_input.Vector)

		}
		vectordata = wtype.MakePlasmidDNASequence("Vector", _input.Vector)
	}

	//lookup restriction enzyme
	restrictionenzyme, err := lookup.TypeIIsLookup(_input.Enzymename)
	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
	}

	//  Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard

	//PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

	if _input.EndsAlreadyadded {
		_output.PartswithOverhangs = partsinorder
	} else {
		_output.PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)
	}

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{_input.Constructname, restrictionenzyme.Name, vectordata, _output.PartswithOverhangs}
	status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

	endreport := "Endreport only run in the event of assembly simulation failure"
	//sites := "Restriction mapper only run in the event of assembly simulation failure"
	_output.NewDNASequence = newDNASequence
	if err == nil && numberofassemblies == 1 {

		_output.Simulationpass = true
	} else {

		warnings = append(warnings, status)
		// perform mock digest to test fragement overhangs (fragments are hidden by using _, )
		_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

		allends := make([]string, 0)
		ends := ""

		ends = text.Print(vectordata.Nm+" 5 Prime end: ", stickyends5)
		allends = append(allends, ends)
		ends = text.Print(vectordata.Nm+" 3 Prime end: ", stickyends3)
		allends = append(allends, ends)

		for _, part := range _output.PartswithOverhangs {
			_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(part, restrictionenzyme)
			ends = text.Print(part.Nm+" 5 Prime end: ", stickyends5)
			allends = append(allends, ends)
			ends = text.Print(part.Nm+" 3 Prime end: ", stickyends3)
			allends = append(allends, ends)
		}
		endreport = strings.Join(allends, " ")
		warnings = append(warnings, endreport)
	}

	// check number of sites per part !

	sites := make([]int, 0)
	multiple := make([]string, 0)

	enz := lookup.EnzymeLookup(_input.Enzymename)
	for _, part := range _output.PartswithOverhangs {

		info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})

		sitepositions := enzymes.SitepositionString(info[0])

		sites = append(sites, info[0].Numberofsites)
		sitepositions = text.Print(part.Nm+" "+_input.Enzymename+" positions:", sitepositions)
		multiple = append(multiple, sitepositions)
	}

	for _, orf := range _input.ORFstoConfirm {
		if sequences.LookforSpecificORF(_output.NewDNASequence.Seq, orf) == false {
			warning = text.Print("orf not present: ", orf)
			warnings = append(warnings, warning)
			_output.ORFmissing = true
		}
	}

	if len(warnings) == 0 {
		warnings = append(warnings, "none")
	}
	_output.Warnings = fmt.Errorf(strings.Join(warnings, ";"))

	partsummary := make([]string, 0)
	for _, part := range _output.PartswithOverhangs {
		partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
	}

	partstoorder := text.Print("PartswithOverhangs: ", partsummary)

	// Print status
	if _output.Status != "all parts available" {
		_output.Status = fmt.Sprintln(_output.Status)
	} else {
		_output.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Sites per part for "+_input.Enzymename, sites),
			text.Print("Positions: ", multiple),
			text.Print("Warnings:", _output.Warnings.Error()),
			text.Print("Simulationpass=", _output.Simulationpass),
			text.Print("NewDNASequence: ", _output.NewDNASequence),
			text.Print("Any Orfs to confirm missing from new DNA sequence:", _output.ORFmissing),
			partstoorder,
		)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Scarfree_designAnalysis(_ctx context.Context, _input *Scarfree_designInput, _output *Scarfree_designOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Scarfree_designValidation(_ctx context.Context, _input *Scarfree_designInput, _output *Scarfree_designOutput) {
}
func _Scarfree_designRun(_ctx context.Context, input *Scarfree_designInput) *Scarfree_designOutput {
	output := &Scarfree_designOutput{}
	_Scarfree_designSetup(_ctx, input)
	_Scarfree_designSteps(_ctx, input, output)
	_Scarfree_designAnalysis(_ctx, input, output)
	_Scarfree_designValidation(_ctx, input, output)
	return output
}

func Scarfree_designRunSteps(_ctx context.Context, input *Scarfree_designInput) *Scarfree_designSOutput {
	soutput := &Scarfree_designSOutput{}
	output := _Scarfree_designRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Scarfree_designNew() interface{} {
	return &Scarfree_designElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Scarfree_designInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Scarfree_designRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Scarfree_designInput{},
			Out: &Scarfree_designOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Scarfree_designElement struct {
	inject.CheckedRunner
}

type Scarfree_designInput struct {
	Constructname    string
	EndsAlreadyadded bool
	Enzymename       string
	ORFstoConfirm    []string
	Seqsinorder      []string
	Vector           string
}

type Scarfree_designOutput struct {
	NewDNASequence     wtype.DNASequence
	ORFmissing         bool
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           error
}

type Scarfree_designSOutput struct {
	Data struct {
		NewDNASequence     wtype.DNASequence
		ORFmissing         bool
		PartswithOverhangs []wtype.DNASequence
		Simulationpass     bool
		Status             string
		Warnings           error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(Component{Name: "Scarfree_design",
		Constructor: Scarfree_designNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/TypeIISAssembly_design/Scarfree_design.an",
			Params: []ParamDesc{
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyadded", Desc: "", Kind: "Parameters"},
				{Name: "Enzymename", Desc: "", Kind: "Parameters"},
				{Name: "ORFstoConfirm", Desc: "enter each as amino acid sequence\n", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "NewDNASequence", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "ORFmissing", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangs", Desc: "parts to order\n", Kind: "Data"},
				{Name: "Simulationpass", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
