package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/entrez"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// input seq

// output parts with correct overhangs

func _GeneDesignRequirements() {
}

func _GeneDesignSetup(_ctx context.Context, _input *GeneDesignInput) {
}

func _GeneDesignSteps(_ctx context.Context, _input *GeneDesignInput, _output *GeneDesignOutput) {
	PartDNA := make([]wtype.DNASequence, 4)

	// Retrieve part seqs from entrez
	for i, part := range _input.Parts {
		DNA, _ := entrez.RetrieveSequence(part, "nucleotide")
		PartDNA[i] = DNA
	}

	// look up vector sequence
	VectorSeq, _ := entrez.RetrieveVector(_input.Vector)

	// Look up the restriction enzyme
	EnzymeInf, _ := lookup.TypeIIsLookup(_input.RE)

	// Add overhangs
	_output.PartsWithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(PartDNA, VectorSeq, EnzymeInf)

	// validation
	assembly := enzymes.Assemblyparameters{"MarksConstruct", EnzymeInf, VectorSeq, _output.PartsWithOverhangs}
	Status, _, _, _, _ := enzymes.Assemblysimulator(assembly)
	fmt.Println(Status)

	// check if sequence meets requirements for synthesis
	sequences.ValidateSynthesis(_output.PartsWithOverhangs, _input.Vector, "GenScript")

	// export sequence to fasta
	export.Makefastaserial2("MarksConstruct", _output.PartsWithOverhangs)

}

func _GeneDesignAnalysis(_ctx context.Context, _input *GeneDesignInput, _output *GeneDesignOutput) {

}

func _GeneDesignValidation(_ctx context.Context, _input *GeneDesignInput, _output *GeneDesignOutput) {

}
func _GeneDesignRun(_ctx context.Context, input *GeneDesignInput) *GeneDesignOutput {
	output := &GeneDesignOutput{}
	_GeneDesignSetup(_ctx, input)
	_GeneDesignSteps(_ctx, input, output)
	_GeneDesignAnalysis(_ctx, input, output)
	_GeneDesignValidation(_ctx, input, output)
	return output
}

func GeneDesignRunSteps(_ctx context.Context, input *GeneDesignInput) *GeneDesignSOutput {
	soutput := &GeneDesignSOutput{}
	output := _GeneDesignRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func GeneDesignNew() interface{} {
	return &GeneDesignElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &GeneDesignInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _GeneDesignRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &GeneDesignInput{},
			Out: &GeneDesignOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type GeneDesignElement struct {
	inject.CheckedRunner
}

type GeneDesignInput struct {
	Parts  []string
	RE     string
	Vector string
}

type GeneDesignOutput struct {
	PartsWithOverhangs []wtype.DNASequence
	Sequence           string
}

type GeneDesignSOutput struct {
	Data struct {
		PartsWithOverhangs []wtype.DNASequence
		Sequence           string
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "GeneDesign",
		Constructor: GeneDesignNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/GeneDesign/GeneDesign.an",
			Params: []ParamDesc{
				{Name: "Parts", Desc: "", Kind: "Parameters"},
				{Name: "RE", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "PartsWithOverhangs", Desc: "output parts with correct overhangs\n", Kind: "Data"},
				{Name: "Sequence", Desc: "input seq\n", Kind: "Data"},
			},
		},
	})
}
