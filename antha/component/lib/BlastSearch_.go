// Example element demonstrating how to perform a BLAST search using the megablast algorithm

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	biogo "github.com/biogo/ncbi/blast"
)

// Input parameters for this protocol

// Data which is returned from this protocol; output data

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
	var hitsummary string

	// Convert the sequence to an anthatype
	_output.AnthaSeq = wtype.MakeLinearDNASequence(_input.Name, _input.DNA)

	// look for orfs
	orf, orftrue := sequences.FindORF(_output.AnthaSeq.Seq)

	if orftrue == true && len(orf.DNASeq) == len(_output.AnthaSeq.Seq) {
		// if open reading frame is detected, we'll perform a blastP search'
		fmt.Println("ORF detected:", "full sequence length: ", len(_output.AnthaSeq.Seq), "ORF length: ", len(orf.DNASeq))
		hits, err = blast.MegaBlastP(orf.ProtSeq)
	} else {
		// otherwise we'll blast the nucleotide sequence
		hits, err = _output.AnthaSeq.Blast()
	}
	if err != nil {
		fmt.Println(err.Error())

	}
	hitsummary, err = blast.HitSummary(hits, 10, 10)

	_output.Hits = hits
	_output.Hitssummary = hitsummary
	fmt.Println(_output.Hitssummary)
	// Rename Sequence with ID of top blast hit
	_output.AnthaSeq.Nm = hits[0].Id

	_output.Warning = err

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
	AnthaSeq    wtype.DNASequence
	Hits        []biogo.Hit
	Hitssummary string
	Warning     error
}

type BlastSearchSOutput struct {
	Data struct {
		AnthaSeq    wtype.DNASequence
		Hits        []biogo.Hit
		Hitssummary string
		Warning     error
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "BlastSearch",
		Constructor: BlastSearchNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/BlastSearch/BlastSearch.an",
			Params: []ParamDesc{
				{Name: "DNA", Desc: "", Kind: "Parameters"},
				{Name: "Name", Desc: "", Kind: "Parameters"},
				{Name: "AnthaSeq", Desc: "", Kind: "Data"},
				{Name: "Hits", Desc: "", Kind: "Data"},
				{Name: "Hitssummary", Desc: "", Kind: "Data"},
				{Name: "Warning", Desc: "", Kind: "Data"},
			},
		},
	})
}
