package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// if buffer is being added

//NumberofMastermixes int // add as many as possible option e.g. if == -1

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// optional if nil this is ignored

// Physical outputs from this protocol with types

func _Mastermix_oneRequirements() {
}

// Conditions to run on startup
func _Mastermix_oneSetup(_ctx context.Context, _input *Mastermix_oneInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Mastermix_oneSteps(_ctx context.Context, _input *Mastermix_oneInput, _output *Mastermix_oneOutput) {

	var mastermix *wtype.LHComponent

	// work out volume to top up to in each case (per reaction) in l:
	topupVolumeperreacttion := _input.TotalVolumeperreaction.SIValue() - _input.VolumetoLeaveforDNAperreaction.SIValue()

	// multiply by number of reactions per mastermix
	topupVolume := wunit.NewVolume(float64(_input.Reactionspermastermix)*topupVolumeperreacttion, "l")

	if len(_input.Components) != len(_input.ComponentVolumesperReaction) {
		panic("len(Components) != len(OtherComponentVolumes)")
	}

	if _input.AliquotbyRow {
		panic("MixTo based method coming soon!")
	} else {

		eachmastermix := make([]*wtype.LHComponent, 0)

		if _input.TopUpBuffer != nil {
			bufferSample := mixer.SampleForTotalVolume(_input.TopUpBuffer, topupVolume)
			eachmastermix = append(eachmastermix, bufferSample)
		}

		for k, component := range _input.Components {
			if k == len(_input.Components) {
				component.Type = wtype.LTNeedToMix //"NeedToMix"
			}

			// multiply volume of each component by number of reactions per mastermix
			adjustedvol := wunit.NewVolume(float64(_input.Reactionspermastermix)*_input.ComponentVolumesperReaction[k].SIValue(), "l")

			componentSample := mixer.Sample(component, adjustedvol)
			eachmastermix = append(eachmastermix, componentSample)
		}

		mastermix = execute.MixInto(_ctx, _input.OutPlate, "", eachmastermix...)

	}
	_output.Mastermix = mastermix

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Mastermix_oneAnalysis(_ctx context.Context, _input *Mastermix_oneInput, _output *Mastermix_oneOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Mastermix_oneValidation(_ctx context.Context, _input *Mastermix_oneInput, _output *Mastermix_oneOutput) {
}
func _Mastermix_oneRun(_ctx context.Context, input *Mastermix_oneInput) *Mastermix_oneOutput {
	output := &Mastermix_oneOutput{}
	_Mastermix_oneSetup(_ctx, input)
	_Mastermix_oneSteps(_ctx, input, output)
	_Mastermix_oneAnalysis(_ctx, input, output)
	_Mastermix_oneValidation(_ctx, input, output)
	return output
}

func Mastermix_oneRunSteps(_ctx context.Context, input *Mastermix_oneInput) *Mastermix_oneSOutput {
	soutput := &Mastermix_oneSOutput{}
	output := _Mastermix_oneRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Mastermix_oneNew() interface{} {
	return &Mastermix_oneElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Mastermix_oneInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Mastermix_oneRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Mastermix_oneInput{},
			Out: &Mastermix_oneOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Mastermix_oneElement struct {
	inject.CheckedRunner
}

type Mastermix_oneInput struct {
	AliquotbyRow                   bool
	ComponentVolumesperReaction    []wunit.Volume
	Components                     []*wtype.LHComponent
	Inplate                        *wtype.LHPlate
	OutPlate                       *wtype.LHPlate
	Reactionspermastermix          int
	TopUpBuffer                    *wtype.LHComponent
	TotalVolumeperreaction         wunit.Volume
	VolumetoLeaveforDNAperreaction wunit.Volume
}

type Mastermix_oneOutput struct {
	Mastermix *wtype.LHComponent
	Status    string
}

type Mastermix_oneSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Mastermix *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Mastermix_one",
		Constructor: Mastermix_oneNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/MakeMastermix/Mastermix_reactions.an",
			Params: []ParamDesc{
				{Name: "AliquotbyRow", Desc: "NumberofMastermixes int // add as many as possible option e.g. if == -1\n", Kind: "Parameters"},
				{Name: "ComponentVolumesperReaction", Desc: "", Kind: "Parameters"},
				{Name: "Components", Desc: "", Kind: "Inputs"},
				{Name: "Inplate", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Reactionspermastermix", Desc: "", Kind: "Parameters"},
				{Name: "TopUpBuffer", Desc: "optional if nil this is ignored\n", Kind: "Inputs"},
				{Name: "TotalVolumeperreaction", Desc: "if buffer is being added\n", Kind: "Parameters"},
				{Name: "VolumetoLeaveforDNAperreaction", Desc: "", Kind: "Parameters"},
				{Name: "Mastermix", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	})
}
