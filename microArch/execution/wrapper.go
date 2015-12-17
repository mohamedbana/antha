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
// 2 Royal College St, London NW1 0NH UK

// Some wrappers to simplify an->go code generation
package execution

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/internal/github.com/fatih/structs"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipmentManager"
	"github.com/antha-lang/antha/microArch/logger"
)

type Wrapper struct {
	usedMix       bool
	usedIncubate  bool
	liquidHandler equipment.Equipment
	threadID      execute.ThreadID
	blockID       execute.BlockID
	outputCount   int
	paramBlock    interface{}
}

func NewWrapper(threadID execute.ThreadID, blockID execute.BlockID, paramBlock interface{}) *Wrapper {
	//TODO delete the rest of unneeded configuration parameters
	w := &Wrapper{}
	w.threadID = threadID
	w.blockID = blockID
	w.paramBlock = paramBlock
	return w
}

func (w *Wrapper) Incubate(what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) {
	logger.Debug(fmt.Sprintln("INCUBATE: ", temp.ToString(), " ", time.ToString(), " shaking? ", shaking))
}

func (w *Wrapper) getString(x string) string {
	m := structs.Map(w.paramBlock)
	if v, ok := m[x].(string); ok {
		return v
	} else {
		return ""
	}
}

func (w *Wrapper) MixTo(outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHSolution {
	if !w.usedMix {
		if em := equipmentManager.GetEquipmentManager(); em != nil {
			lh := em.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", nil))
			if lh == nil {
				panic("error configuring liquid handling request: could not find equipment that satisfies liquid handler mix instruction")
			}
			w.liquidHandler = lh
		} else {
			panic("equipment manager not configured")
		}

		//We are going to configure the liquid handler for a blockId. BlockId will give us the framework and state handling
		// so, for a certain BlockId config options will be aggregated. Liquid Handler will just regenerate all state per
		// this aggregation layer and that will allow us to run multiple protocols.
		//prepare the values
		config := make(map[string]interface{}) //new(wtype.ConfigItem)
		config["MAX_N_PLATES"] = 4.5
		config["MAX_N_WELLS"] = 278.0
		config["RESIDUAL_VOLUME_WEIGHT"] = 1.0
		config["OUTPUT_COUNT"] = w.outputCount
		config["BLOCKID"] = w.blockID.String()
		// these should come from the paramblock... for now though
		config["INPUT_PLATETYPE"] = "pcrplate_with_cooler"
		config["OUTPUT_PLATETYPE"] = "pcrplate_with_cooler"

		configString, err := json.Marshal(config)
		if err != nil {
			panic(fmt.Sprintf("error configuring liquid handling request: %s", err))
		}
		if w.liquidHandler != nil {
			w.liquidHandler.Do(*equipment.NewActionDescription(action.LH_CONFIG, string(configString), nil))
		}
	}

	reaction := mixer.MixTo(outplate, address, components...)
	reaction.BlockID = w.blockID
	reaction.SName = w.getString("OutputReactionName")
	reqReaction, err := json.Marshal(reaction)
	if err != nil {
		panic(fmt.Sprintf("error coding reaction data, %v", err))
	}
	if w.liquidHandler != nil {
		err = w.liquidHandler.Do(*equipment.NewActionDescription(action.LH_MIX, string(reqReaction), nil))
	}
	if err != nil {
		panic(fmt.Sprintf("error running liquid handling request: %s", err))
	}
	return reaction

}

func (w *Wrapper) MixInto(outplate *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution {
	//func (w *Wrapper) MixInto(outplate *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution {
	// TODO: need better error handling here so we don't take down the monolith
	// when, for example, we're asked to simulate a workflow without having a
	// liquid handler
	if !w.usedMix {
		if em := equipmentManager.GetEquipmentManager(); em != nil {
			lh := em.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", nil))
			if lh == nil {
				panic("error configuring liquid handling request: could not find equipment that satisfies liquid handler mix instruction")
			}
			w.liquidHandler = lh
		} else {
			panic("equipment manager not configured")
		}

		//We are going to configure the liquid handler for a blockId. BlockId will give us the framework and state handling
		// so, for a certain BlockId config options will be aggregated. Liquid Handler will just regenerate all state per
		// this aggregation layer and that will allow us to run multiple protocols.
		//prepare the values
		config := make(map[string]interface{}) //new(wtype.ConfigItem)
		config["MAX_N_PLATES"] = 4.5
		config["MAX_N_WELLS"] = 278.0
		config["RESIDUAL_VOLUME_WEIGHT"] = 1.0
		config["OUTPUT_COUNT"] = w.outputCount
		config["BLOCKID"] = w.blockID.String()
		// these should come from the paramblock... for now though
		//config["INPUT_PLATETYPE"] = "DWR1" // seems to still be hardcoded as pcr plate with cooler somewhere else
		//config["OUTPUT_PLATETYPE"] = "pcrplate_with_cooler"

		configString, err := json.Marshal(config)
		if err != nil {
			panic(fmt.Sprintf("error configuring liquid handling request: %s", err))
		}
		if w.liquidHandler != nil {
			w.liquidHandler.Do(*equipment.NewActionDescription(action.LH_CONFIG, string(configString), nil))
		}
	}

	reaction := mixer.MixInto(outplate, components...)
	reaction.BlockID = w.blockID
	reaction.SName = w.getString("OutputReactionName")
	reqReaction, err := json.Marshal(reaction)
	if err != nil {
		panic(fmt.Sprintf("error coding reaction data, %v", err))
	}
	if w.liquidHandler != nil {
		err = w.liquidHandler.Do(*equipment.NewActionDescription(action.LH_MIX, string(reqReaction), nil))
	}
	if err != nil {
		panic(fmt.Sprintf("error running liquid handling request: %s", err))
	}
	return reaction
}

func (w *Wrapper) WaitToEnd() error {
	if w.liquidHandler != nil {
		w.liquidHandler.Do(*equipment.NewActionDescription(action.LH_END, w.blockID.String(), nil))
	}
	return nil
}
