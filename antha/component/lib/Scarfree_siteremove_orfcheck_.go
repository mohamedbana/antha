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
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"

// Input parameters for this protocol (data)

// enter each as amino acid sequence

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func _Scarfree_siteremove_orfcheckRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _Scarfree_siteremove_orfcheckSetup(_ctx context.Context, _input *Scarfree_siteremove_orfcheckInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Scarfree_siteremove_orfcheckSteps(_ctx context.Context, _input *Scarfree_siteremove_orfcheckInput, _output *Scarfree_siteremove_orfcheckOutput) {

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// declare some temporary variables to be used later
	var warning string
	var err error

	// make an empty array of DNA Sequences ready to fill
	partsinorder := make([]wtype.DNASequence, 0)

	var partDNA wtype.DNASequence
	var vectordata wtype.DNASequence

	_output.Status = "all parts available"
	for i, part := range _input.Seqsinorder {
		if strings.Contains(part, ".gb") && strings.Contains(part, "Feature:") {

			split := strings.SplitAfter(part, ".gb")
			file := split[0]

			split2 := strings.Split(split[1], ":")
			feature := split2[1]

			partDNA, _ = parser.GenbankFeaturetoDNASequence(file, feature)
		} else if strings.Contains(part, ".gb") {

			partDNA, _ = parser.GenbanktoAnnotatedSeq(part)
		} else {
			nm := "Part " + strconv.Itoa(i)
			if strings.Contains(part, "BBa_") {
				nm = nm + "_" + part
				part = igem.GetSequence(part)
			}
			partDNA = wtype.MakeLinearDNASequence(nm, part)
		}
		partsinorder = append(partsinorder, partDNA)
	}
	// check parts for restriction sites first and remove if the user has chosen to
	enz := lookup.EnzymeLookup(_input.Enzymename)

	warning = text.Print("RemoveproblemRestrictionSites =", _input.RemoveproblemRestrictionSites)
	warnings = append(warnings, warning)
	if _input.RemoveproblemRestrictionSites {
		newparts := make([]wtype.DNASequence, 0)
		warning = "Starting process or removing restrictionsite"
		warnings = append(warnings, warning)

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})

			for _, anysites := range info {
				if anysites.Sitefound {
					warning = "problem site found in " + part.Nm
					warnings = append(warnings, warning)
					orf, orftrue := sequences.FindBiggestORF(part.Seq)
					warning = fmt.Sprintln("site found in orf ", part.Nm, " ", orftrue, " site positions ", anysites.Positions("ALL"), "orf between", orf.StartPosition, " and ", orf.EndPosition /*orf.DNASeq[orf.StartPosition:orf.EndPosition]*/)
					warnings = append(warnings, warning)
					if orftrue /* && len(orf.ProtSeq) > 20 */ {
						allsitestoavoid := make([]string, 0)
						allsitestoavoid = append(allsitestoavoid, anysites.Recognitionsequence, sequences.RevComp(anysites.Recognitionsequence))
						orfcoordinates := sequences.MakeStartendPair(orf.StartPosition, orf.EndPosition)
						for _, position := range anysites.Positions("ALL") {
							if orf.StartPosition < position && position < orf.EndPosition {
								originalcodon := ""
								codonoption := ""
								part, originalcodon, codonoption, err = sequences.ReplaceCodoninORF(part, orfcoordinates, position, allsitestoavoid)
								warning = fmt.Sprintln("sites to avoid: ", allsitestoavoid[0], allsitestoavoid[1])
								warnings = append(warnings, warning)
								warnings = append(warnings, "Paaaaerrttseq: "+part.Seq+"position: "+strconv.Itoa(position)+" original: "+originalcodon+" replacementcodon: "+codonoption)
								if err != nil {
									warning := text.Print("removal of site from orf "+orf.DNASeq, " failed! improve your algorithm! "+err.Error())
									warnings = append(warnings, warning)
								}
							} else {
								allsitestoavoid := make([]string, 0)
								part, err = sequences.RemoveSite(part, anysites.Enzyme, allsitestoavoid)
								if err != nil {

									warning = text.Print("position found to be outside of orf: "+orf.DNASeq, " failed! improve your algorithm! "+err.Error())
									warnings = append(warnings, warning)
								}
							}
						}
					} else {
						allsitestoavoid := make([]string, 0)
						temppart, err := sequences.RemoveSite(part, anysites.Enzyme, allsitestoavoid)
						fmt.Println("part= ", part)
						fmt.Println("temppart= ", temppart)
						if err != nil {
							warning := text.Print("removal of site failed! improve your algorithm!", err.Error())
							warnings = append(warnings, warning)

						}
						warning = fmt.Sprintln("modified "+temppart.Nm+"new seq: ", temppart.Seq, "original seq: ", part.Seq)
						warnings = append(warnings, warning)
						part = temppart

						//	}
					}
				}

				newparts = append(newparts, part)

				partsinorder = newparts
			}
		}
	}

	// make vector into an antha type DNASequence

	if strings.Contains(_input.Vector, ".gb") {

		vectordata, _ = parser.GenbanktoAnnotatedSeq(_input.Vector)
		vectordata.Plasmid = true
	} else {
		vectornm := "Vector"
		if strings.Contains(_input.Vector, "BBa_") || strings.Contains(_input.Vector, "pSB") {
			vectornm = _input.Vector
			_input.Vector = igem.GetSequence(_input.Vector)

		}
		vectordata = wtype.MakePlasmidDNASequence(vectornm, _input.Vector)
	}

	//lookup restriction enzyme
	restrictionenzyme, err := lookup.TypeIIsLookup(_input.Enzymename)
	if err != nil {
		warnings = append(warnings, text.Print("Error", err.Error()))
	}

	//  Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	fmt.Println("warnings:", warnings)

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
	partsummary = append(partsummary, text.Print("Vector:"+vectordata.Nm, vectordata.Seq))
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
		// export data to file
		//anthapath.ExporttoFile("Report"+"_"+Constructname+".txt",[]byte(Status))
		//anthapath.ExportTextFile("Report"+"_"+Constructname+".txt",Status)
		fmt.Println(_output.Status)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Scarfree_siteremove_orfcheckAnalysis(_ctx context.Context, _input *Scarfree_siteremove_orfcheckInput, _output *Scarfree_siteremove_orfcheckOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Scarfree_siteremove_orfcheckValidation(_ctx context.Context, _input *Scarfree_siteremove_orfcheckInput, _output *Scarfree_siteremove_orfcheckOutput) {
}
func _Scarfree_siteremove_orfcheckRun(_ctx context.Context, input *Scarfree_siteremove_orfcheckInput) *Scarfree_siteremove_orfcheckOutput {
	output := &Scarfree_siteremove_orfcheckOutput{}
	_Scarfree_siteremove_orfcheckSetup(_ctx, input)
	_Scarfree_siteremove_orfcheckSteps(_ctx, input, output)
	_Scarfree_siteremove_orfcheckAnalysis(_ctx, input, output)
	_Scarfree_siteremove_orfcheckValidation(_ctx, input, output)
	return output
}

func Scarfree_siteremove_orfcheckRunSteps(_ctx context.Context, input *Scarfree_siteremove_orfcheckInput) *Scarfree_siteremove_orfcheckSOutput {
	soutput := &Scarfree_siteremove_orfcheckSOutput{}
	output := _Scarfree_siteremove_orfcheckRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Scarfree_siteremove_orfcheckNew() interface{} {
	return &Scarfree_siteremove_orfcheckElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Scarfree_siteremove_orfcheckInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Scarfree_siteremove_orfcheckRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Scarfree_siteremove_orfcheckInput{},
			Out: &Scarfree_siteremove_orfcheckOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Scarfree_siteremove_orfcheckElement struct {
	inject.CheckedRunner
}

type Scarfree_siteremove_orfcheckInput struct {
	Constructname                 string
	EndsAlreadyadded              bool
	Enzymename                    string
	ORFstoConfirm                 []string
	RemoveproblemRestrictionSites bool
	Seqsinorder                   []string
	Vector                        string
}

type Scarfree_siteremove_orfcheckOutput struct {
	NewDNASequence     wtype.DNASequence
	ORFmissing         bool
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           error
}

type Scarfree_siteremove_orfcheckSOutput struct {
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
	if err := addComponent(component.Component{Name: "Scarfree_siteremove_orfcheck",
		Constructor: Scarfree_siteremove_orfcheckNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/TypeIISAssembly_design/Scarfree_removesites_checkorfs.an",
			Params: []component.ParamDesc{
				{Name: "Constructname", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyadded", Desc: "", Kind: "Parameters"},
				{Name: "Enzymename", Desc: "", Kind: "Parameters"},
				{Name: "ORFstoConfirm", Desc: "enter each as amino acid sequence\n", Kind: "Parameters"},
				{Name: "RemoveproblemRestrictionSites", Desc: "", Kind: "Parameters"},
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
