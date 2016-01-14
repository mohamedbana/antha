//
package Phytip_miniprep

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

func _requirements() {
}
func _setup(_ctx context.Context, _input *Input_) {
}
func _steps(_ctx context.Context, _input *Input_, _output *Output_) {
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

type Output_ struct {
	PlasmidDNAsolution *wtype.LHComponent
}
