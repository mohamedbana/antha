package NewDNASequence

import (
	"fmt"
	//"math"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _requirements() {

}

// Actions to perform before protocol itself
func _setup(_ctx context.Context, _input *Input_) {

}

// Core process of the protocol: steps to be performed for each input
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
	fmt.Println("In steps!")
	if _input.Plasmid != _input.Linear {
		if _input.Plasmid {
			_output.DNA = wtype.MakePlasmidDNASequence(_input.Gene_name, _input.DNA_seq)
		} else if _input.Linear {
			_output.DNA = wtype.MakeLinearDNASequence(_input.Gene_name, _input.DNA_seq)
		} else if _input.SingleStranded {
			_output.DNA = wtype.MakeSingleStrandedDNASequence(_input.Gene_name, _input.DNA_seq)
		}

		orfs := sequences.FindallORFs(_output.DNA.Seq)
		features := sequences.ORFs2Features(orfs)

		_output.DNAwithORFs = sequences.Annotate(_output.DNA, features)

		_output.Status = fmt.Sprintln(
			text.Print("DNA_Seq: ", _input.DNA_seq),
			text.Print("ORFs: ", _output.DNAwithORFs.Features),
		)

	} else {
		_output.Status = fmt.Sprintln("correct conditions not met")
	}

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
	DNA_seq        string
	Gene_name      string
	Linear         bool
	Plasmid        bool
	SingleStranded bool
}

type Output_ struct {
	DNA         wtype.DNASequence
	DNAwithORFs sequences.AnnotatedSeq
	Status      string
}
