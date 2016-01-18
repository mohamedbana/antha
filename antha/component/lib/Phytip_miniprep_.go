//
package lib

import (
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/UnitOperations"
	//"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"time"
)

//Cellpelletmass Mass

//Torr

// cubesensor streams to work out drying time:
/*Pa float64 // in pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa
Temp float64 // in Kelvin
Relativehumidity float64 // Percentage // density water vapor (kg/m3)
*/
//Time time.Duration //float64// time

/*
	Parameters before refactoring into Chromstep structs

	RBvolume Volume // 150ul
	RBflowrate Rate
	RBpause Time // seconds
	RBcycles int

	LBvolume Volume
	LBflowrate Rate
	LBpause Time
	LBcycles int

	PBvolume Volume
	PBflowrate Rate
	PBpause Time
	PBcycles int

	Equilibrationvolume Volume
	Equilibrationflowrate Rate
	Equilibrationpause Time
	Equilibrationcycles int

	Airdispensevolume Volume
	Airdispenseflowrate Rate
	Airdispensepause Time
	Airdispensecycles int



	Airaspiratevolume Volume
	Airaspirateflowrate Rate
	Airaspiratepause Time
	Airaspiratecylces int

	Capturevoume Volume
	Captureflowrate Rate
	Capturepause Time
	Capturecycles int

	Washbuffervolume [] Volume
	Washbufferflowrate [] Rate
	Washbufferpause [] Time
	Washbuffercycles [] int



	Elutionbuffervolume Volume
	Elutionflowrate Rate
	Elutionpause Time
	Elutioncycles int

*/
//or

/* PlasmidConc Concentration
Storagelocation Location
Storageconditions StorageHistory
Plasmidbuffer Composition */ // is this all inferred from a PLasmid solution  type anyway?

//
// wtype.LHTip
//UnitOperations.Pellet // wrong type?

//RB *wtype.LHComponent //Watersolution
//LB *wtype.LHComponent //Watersolution
//PB *wtype.LHComponent //Watersolution
//Water *wtype.LHComponent //Watersolution // equilibration buffer
//Air *wtype.LHComponent //Gas
//Washbuffer []*wtype.LHComponent //Watersolution
//Elutionbuffer *wtype.LHComponent //Watersolution

//Solution //PlasmidSolution

func _Phytip_miniprepRequirements() {
}
func _Phytip_miniprepSetup(_ctx context.Context, _input *Phytip_miniprepInput) {
}
func _Phytip_miniprepSteps(_ctx context.Context, _input *Phytip_miniprepInput, _output *Phytip_miniprepOutput) {
	resuspension, _ := UnitOperations.Resuspend(_input.Cellpellet, _input.Resuspensionstep, _input.Tips)
	lysate, _ := UnitOperations.Chromatography(resuspension, _input.Lysisstep, _input.Tips)
	precipitate, _ := UnitOperations.Chromatography(lysate, _input.Precipitationstep, _input.Tips)

	_, columnready := UnitOperations.Chromatography(_input.Equilibrationstep.Buffer, _input.Equilibrationstep, _input.Phytips)

	_, readyforcapture := UnitOperations.Chromatography(_input.Airstep.Buffer, _input.Airstep, columnready)
	capture, readyforcapture := UnitOperations.Chromatography(precipitate, _input.Capturestep, readyforcapture)

	for _, washstep := range _input.Washsteps {
		_, readyforcapture = UnitOperations.Chromatography(capture, washstep, readyforcapture)
	}
	readyfordrying := UnitOperations.Blot(readyforcapture, _input.Blotcycles, _input.Blottime)

	/*if Vacuum == true {
		drytips := UnitOperations.Dry(Tips,Drytime,Vacuumstrength)


		//parameters required for evaporation calculator
		Liquid := Washsteps[0].Pipetstep.Name //ethanol?
		// lookup properties via liquidclasses package to workout evaporation time using Evaporationrate element?


		//Platetype := Phytips.tip //.surfacearea? labware.phytip.surfacearea?
		Volumeperwell := (Washsteps[0].Pipetstep.Volume.SIValue() / 10) // assume max 10% residual volume for now??

		drytimerequired := Evaporation.Estimatedevaporationtime(Airvelocity, Liquid, Platetype, Volumeperwell)


	} else {*/drytips := UnitOperations.Dry(readyfordrying, _input.Drytime, _input.Vacuumstrength) //}

	_output.PlasmidDNAsolution, _ = UnitOperations.Chromatography(_input.Elutionstep.Buffer, _input.Elutionstep, drytips)

}
func _Phytip_miniprepAnalysis(_ctx context.Context, _input *Phytip_miniprepInput, _output *Phytip_miniprepOutput) {
}
func _Phytip_miniprepValidation(_ctx context.Context, _input *Phytip_miniprepInput, _output *Phytip_miniprepOutput) {
}
func _Phytip_miniprepRun(_ctx context.Context, input *Phytip_miniprepInput) *Phytip_miniprepOutput {
	output := &Phytip_miniprepOutput{}
	_Phytip_miniprepSetup(_ctx, input)
	_Phytip_miniprepSteps(_ctx, input, output)
	_Phytip_miniprepAnalysis(_ctx, input, output)
	_Phytip_miniprepValidation(_ctx, input, output)
	return output
}

func Phytip_miniprepRunSteps(_ctx context.Context, input *Phytip_miniprepInput) *Phytip_miniprepSOutput {
	soutput := &Phytip_miniprepSOutput{}
	output := _Phytip_miniprepRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Phytip_miniprepNew() interface{} {
	return &Phytip_miniprepElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Phytip_miniprepInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Phytip_miniprepRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Phytip_miniprepInput{},
			Out: &Phytip_miniprepOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Phytip_miniprepElement struct {
	inject.CheckedRunner
}

type Phytip_miniprepInput struct {
	Airstep           UnitOperations.Chromstep
	Blotcycles        int
	Blottime          time.Duration
	Capturestep       UnitOperations.Chromstep
	Cellpellet        *wtype.Physical
	Drytime           time.Duration
	Elutionstep       UnitOperations.Chromstep
	Equilibrationstep UnitOperations.Chromstep
	Lysisstep         UnitOperations.Chromstep
	Phytips           UnitOperations.Column
	Precipitationstep UnitOperations.Chromstep
	Resuspensionstep  UnitOperations.Chromstep
	Tips              UnitOperations.Column
	Vacuum            bool
	Vacuumstrength    float64
	Washsteps         []UnitOperations.Chromstep
}

type Phytip_miniprepOutput struct {
	PlasmidDNAsolution *wtype.LHComponent
}

type Phytip_miniprepSOutput struct {
	Data struct {
	}
	Outputs struct {
		PlasmidDNAsolution *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "Phytip_miniprep",
		Constructor: Phytip_miniprepNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/Phytip_miniprep/Phytip_miniprep.an",
			Params: []ParamDesc{
				{Name: "Airstep", Desc: "", Kind: "Parameters"},
				{Name: "Blotcycles", Desc: "", Kind: "Parameters"},
				{Name: "Blottime", Desc: "", Kind: "Parameters"},
				{Name: "Capturestep", Desc: "", Kind: "Parameters"},
				{Name: "Cellpellet", Desc: "UnitOperations.Pellet // wrong type?\n", Kind: "Inputs"},
				{Name: "Drytime", Desc: "", Kind: "Parameters"},
				{Name: "Elutionstep", Desc: "", Kind: "Parameters"},
				{Name: "Equilibrationstep", Desc: "", Kind: "Parameters"},
				{Name: "Lysisstep", Desc: "", Kind: "Parameters"},
				{Name: "Phytips", Desc: "", Kind: "Inputs"},
				{Name: "Precipitationstep", Desc: "", Kind: "Parameters"},
				{Name: "Resuspensionstep", Desc: "", Kind: "Parameters"},
				{Name: "Tips", Desc: "wtype.LHTip\n", Kind: "Inputs"},
				{Name: "Vacuum", Desc: "", Kind: "Parameters"},
				{Name: "Vacuumstrength", Desc: "Torr\n", Kind: "Parameters"},
				{Name: "Washsteps", Desc: "", Kind: "Parameters"},
				{Name: "PlasmidDNAsolution", Desc: "Solution //PlasmidSolution\n", Kind: "Outputs"},
			},
		},
	})
}
