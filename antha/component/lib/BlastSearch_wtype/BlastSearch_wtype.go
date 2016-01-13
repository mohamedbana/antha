// Example element demonstrating how to perform a BLAST search using the megablast algorithm

package BlastSearch_wtype

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	biogo "github.com/antha-lang/antha/internal/github.com/biogo/ncbi/blast"
)

// Input parameters for this protocol

// Data which is returned from this protocol; output data

// Physical inputs to this protocol

// Physical outputs from this protocol

func _requirements() {

}

// Actions to perform before protocol itself
func _setup(_ctx context.Context, _input *Input_) {

}

// Core process of the protocol: steps to be performed for each input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {

	var err error
	var hits []biogo.Hit
	/*
		if Querytype == "PROTEIN" {
		hits, err = blast.MegaBlastP(Query)
		if err != nil {
			fmt.Println(err.Error())
		}

		Hits = fmt.Sprintln(blast.HitSummary(hits))


		} else if Querytype == "DNA" {
		hits, err = blast.MegaBlastN(Query)
		if err != nil {
			fmt.Println(err.Error())
		}

		Hits = fmt.Sprintln(blast.HitSummary(hits))
		}
	*/
	_output.AnthaSeq = _input.DNA

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

	} //else {

	_output.Hits = fmt.Sprintln(blast.HitSummary(hits))

	// Rename Sequence with ID of top blast hit
	_output.AnthaSeq.Nm = hits[0].Id
	//}

}

// Actions to perform after steps block to analyze data
func _analysis(_ctx context.Context, _input *Input_, _output *Output_) {

}

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
	DNA wtype.DNASequence
}

type Output_ struct {
	AnthaSeq wtype.DNASequence
	Hits     string
}
