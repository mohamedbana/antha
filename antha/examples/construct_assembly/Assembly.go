// /examples/construct_assembly/Assembly.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/antha-lang/antha/antha/anthalib/execution"
	"github.com/antha-lang/antha/antha/anthalib/factory"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/internal/github.com/nu7hatch/gouuid"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/manual"
	"github.com/antha-lang/antha/microArch/equipmentManager"
	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/microArch/logger/middleware"
)

//r, params.Reactiontemp, params.Reactiontime, false
func incubate(what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) {
	fmt.Println("INCUBATE: ", temp.ToString(), " ", time.ToString(), " shaking? ", shaking)
}

func main() {
	eid, _ := uuid.NewV4()
	em := equipmentManager.NewAnthaEquipmentManager(eid.String())
	defer em.Shutdown()
	eem := equipmentManager.EquipmentManager(em)
	equipmentManager.SetEquipmentManager(&eem)
	//manual driver equipment
	mid, _ := uuid.NewV4()
	var mde equipment.Equipment
	var amd manual.AnthaManual
	amd = *manual.NewAnthaManual(mid.String())
	mde = amd
	em.RegisterEquipment(&mde)

	//cui logger middleware
	cmw := middleware.NewLogToCui(&amd.Cui)
	log_id, _ := uuid.NewV4()
	l := logger.NewAnthaFileLogger(log_id.String())
	l.RegisterMiddleware(cmw)
	var params Parameters
	var inputs Inputs

	// give this thing an arbitrary ID for testing

	id := execute.ThreadID(fmt.Sprintf("EXPERIMENT_1_%s", string(eid.String()[1:5])))

	fmt.Println(id)

	// set up parameters and inputs

	params.Reactionvolume = wunit.NewVolume(20, "ul")
	params.Partconc = wunit.NewConcentration(0.0001, "g/l")
	params.Vectorconc = wunit.NewConcentration(0.001, "g/l")
	params.Atpvol = wunit.NewVolume(1, "ul")
	params.Revol = wunit.NewVolume(1, "ul")
	params.Ligvol = wunit.NewVolume(1, "ul")
	params.Reactiontemp = wunit.NewTemperature(25, "C")
	params.Reactiontime = wunit.NewTime(1800, "s")
	params.Inactivationtemp = wunit.NewTemperature(40, "C")
	params.Inactivationtime = wunit.NewTime(60, "s")
	params.BlockID = id

	inputs.Parts = make([]*wtype.LHComponent, 4)

	for i := 0; i < 4; i++ {
		inputs.Parts[i] = factory.GetComponentByType("dna_part")
		inputs.Parts[i].CName = inputs.Parts[i].CName + "_" + strconv.Itoa(i+1)
	}

	inputs.Vector = factory.GetComponentByType("standard_cloning_vector_mark_1")
	inputs.RestrictionEnzyme = factory.GetComponentByType("SapI")
	inputs.Ligase = factory.GetComponentByType("T4Ligase")
	inputs.Buffer = factory.GetComponentByType("CutsmartBuffer")
	inputs.ATP = factory.GetComponentByType("ATP")
	inputs.Outplate = factory.GetPlateByType("pcrplate")
	inputs.TipType = factory.GetTipboxByType("Gilson50")

	ctx := execution.GetContext()
	conf := make(map[string]interface{})
	conf["MAX_N_PLATES"] = 1.5
	conf["MAX_N_WELLS"] = 12.0
	conf["RESIDUAL_VOLUME_WEIGHT"] = 1.0
	conf["SQLITE_FILE_IN"] = "/Users/msadowski/synthace/protocol_language/checkout/synthace-antha/anthalib/driver/liquidhandling/pm_driver/default.sqlite"
	conf["SQLITE_FILE_OUT"] = "/tmp/output_file.sqlite"
	ctx.ConfigService.SetConfig(id, conf)
	outputs := Steps(params, inputs)
	fmt.Println(outputs.Reaction)
}

type Parameters struct {
	Reactionvolume   wunit.Volume
	Partconc         wunit.Concentration
	Vectorconc       wunit.Concentration
	Atpvol           wunit.Volume
	Revol            wunit.Volume
	Ligvol           wunit.Volume
	Reactiontemp     wunit.Temperature
	Reactiontime     wunit.Time
	Inactivationtemp wunit.Temperature
	Inactivationtime wunit.Time
	BlockID          execute.ThreadID
}

type Inputs struct {
	Parts             []*wtype.LHComponent
	Vector            *wtype.LHComponent
	RestrictionEnzyme *wtype.LHComponent
	Buffer            *wtype.LHComponent
	Ligase            *wtype.LHComponent
	ATP               *wtype.LHComponent
	Outplate          *wtype.LHPlate
	TipType           *wtype.LHTipbox
}

type Outputs struct {
	Reaction *wtype.LHSolution
}

func Steps(params Parameters, inputs Inputs) Outputs {
	var outputs Outputs

	ctx := execution.GetContext()

	em := ctx.EquipmentManager
	rqout := em.MakeDeviceRequest("liquidhandler", "Manual")
	response := <-rqout

	// handle response problems
	if response["status"] == "FAIL" {
		log.Fatal("Error requesting liquid handler service")
	}

	liquidhandler := response["devicequeue"].(*execution.LiquidHandlingService)

	samples := make([]*wtype.LHComponent, 0)

	buffersample := mixer.SampleForTotalVolume(inputs.Buffer, params.Reactionvolume)
	samples = append(samples, buffersample)

	atpsample := mixer.Sample(inputs.ATP, params.Atpvol)
	samples = append(samples, atpsample)

	vectorsample := mixer.SampleForConcentration(inputs.Vector, params.Vectorconc)
	samples = append(samples, vectorsample)

	for _, part := range inputs.Parts {
		partsample := mixer.SampleForConcentration(part, params.Partconc)
		samples = append(samples, partsample)
	}

	resample := mixer.Sample(inputs.RestrictionEnzyme, params.Revol)
	samples = append(samples, resample)
	ligsample := mixer.Sample(inputs.Ligase, params.Ligvol)
	samples = append(samples, ligsample)
	reaction := mixer.MixInto(inputs.Outplate, samples...)
	// set the block ID
	reaction.BlockID = string(params.BlockID)

	rq2 := liquidhandler.MakeMixRequest(reaction)
	if rq2 == nil {
		log.Fatal("Error running liquid handling request")
	}

	rq2.Tip_Type = inputs.TipType
	liquidhandler.Run()

	// incubate the reaction mixtures

	incubate(reaction, params.Reactiontemp, params.Reactiontime, false)

	// inactivate

	incubate(reaction, params.Inactivationtemp, params.Inactivationtime, false)

	// all done
	outputs.Reaction = reaction
	return outputs
}
