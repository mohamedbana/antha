// Example element demonstrating how to perform a BLAST search using the megablast algorithm

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	biogo "github.com/biogo/ncbi/blast"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// Data which is returned from this protocol; output data

//AnthaSeq wtype.DNASequence

// Physical inputs to this protocol

// Physical outputs from this protocol

func _BlastSearchRequirements() {

}

// Actions to perform before protocol itself
func _BlastSearchSetup(_ctx context.Context, _input *BlastSearchInput) {

}

// Core process of the protocol: steps to be performed for each input
func _BlastSearchSteps(_ctx context.Context, _input *BlastSearchInput, _output *BlastSearchOutput) {

	var err error
	var hits []biogo.Hit

	// Convert the sequence to an anthatype
	AnthaSeq := wtype.MakeLinearDNASequence(_input.Name, _input.DNA)

	// look for orfs
	orf, orftrue := sequences.FindORF(AnthaSeq.Seq)

	if orftrue == true && len(orf.DNASeq) == len(AnthaSeq.Seq) {
		// if open reading frame is detected, we'll perform a blastP search'
		fmt.Println("ORF detected:", "full sequence length: ", len(AnthaSeq.Seq), "ORF length: ", len(orf.DNASeq))
		hits, err = blast.MegaBlastP(orf.ProtSeq)
	} else {
		// otherwise we'll blast the nucleotide sequence
		hits, err = AnthaSeq.Blast()
	}
	if err != nil {
		fmt.Println(err.Error())

	} //else {

	//Hits = fmt.Sprintln(blast.HitSummary(hits))

	// Rename Sequence with ID of top blast hit
	AnthaSeq.Nm = hits[0].Id
	//}

}

// Actions to perform after steps block to analyze data
func _BlastSearchAnalysis(_ctx context.Context, _input *BlastSearchInput, _output *BlastSearchOutput) {

}

func _BlastSearchValidation(_ctx context.Context, _input *BlastSearchInput, _output *BlastSearchOutput) {

}
func _BlastSearchRun(_ctx context.Context, input *BlastSearchInput) *BlastSearchOutput {
	output := &BlastSearchOutput{}
	_BlastSearchSetup(_ctx, input)
	_BlastSearchSteps(_ctx, input, output)
	_BlastSearchAnalysis(_ctx, input, output)
	_BlastSearchValidation(_ctx, input, output)
	return output
}

func BlastSearchRunSteps(_ctx context.Context, input *BlastSearchInput) *BlastSearchSOutput {
	soutput := &BlastSearchSOutput{}
	output := _BlastSearchRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func BlastSearchNew() interface{} {
	return &BlastSearchElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &BlastSearchInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _BlastSearchRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &BlastSearchInput{},
			Out: &BlastSearchOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type BlastSearchElement struct {
	inject.CheckedRunner
}

type BlastSearchInput struct {
	DNA  string
	Name string
}

type BlastSearchOutput struct {
	Hits string
}

type BlastSearchSOutput struct {
	Data struct {
		Hits string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "BlastSearch",
		Constructor: BlastSearchNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/BlastSearch/BlastSearch.an",
			Params: []component.ParamDesc{
				{Name: "DNA", Desc: "", Kind: "Parameters"},
				{Name: "Name", Desc: "", Kind: "Parameters"},
				{Name: "Hits", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
