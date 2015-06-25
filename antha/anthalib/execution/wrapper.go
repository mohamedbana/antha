// /anthalib/execution/wrapper.go: Part of the Antha language
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

// Some wrappers to simplify an->go code generation
package execution

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/factory"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"log"
)

type Wrapper struct {
	usedMix       bool
	usedIncubate  bool
	liquidHandler *LiquidHandlingService
	tipType       *wtype.LHTipbox
	threadID      execute.ThreadID
}

func NewWrapper(threadID execute.ThreadID) *Wrapper {
	w := &Wrapper{}
	w.tipType = factory.GetTipboxByType("Gilson50") // TODO use real configuration
	w.threadID = threadID
	return w
}

func (w *Wrapper) Incubate(what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) {
	fmt.Println("INCUBATE: ", temp.ToString(), " ", time.ToString(), " shaking? ", shaking)
}

func (w *Wrapper) MixInto(outplate *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution {
	if !w.usedMix {
		ctx := GetContext()
		em := ctx.EquipmentManager
		/*
			cfg := ctx.ConfigService.GetConfig(w.threadID)
			devname := ""
			_, ok := cfg["LIQUIDHANDLER"]
			if ok {
				devname = cfg["LIQUIDHANDLER"].(string)
			}
			rqOut := em.MakeDeviceRequest("liquidhandler", devname)
		*/
		rqOut := em.MakeDeviceRequest("liquidhandler", "Manual")
		response := <-rqOut
		if response["status"] == "FAIL" {
			log.Fatal("Error requesting liquid handler service")
		}
		w.liquidHandler = response["devicequeue"].(*LiquidHandlingService)
		w.usedMix = true
	}

	reaction := mixer.MixInto(outplate, components...)
	reaction.BlockID = string(w.threadID)

	req := w.liquidHandler.MakeMixRequest(reaction)
	if req == nil {
		log.Fatal("Error running liquid handling request")
	}
	req.Tip_Type = w.tipType
	w.liquidHandler.Run()

	return reaction
}
