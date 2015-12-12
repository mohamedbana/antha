package execute

import (
	"encoding/json"
	"errors"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipmentManager"
)

var (
	noEquipmentManager = errors.New("no equipment manager configured")
	cannotConfigLh     = errors.New("cannot configure liquid handler")
)

type lhKey int

const theLhKey lhKey = 0

func getLh(ctx context.Context) equipment.Equipment {
	v, ok := ctx.Value(theLhKey).(equipment.Equipment)
	if !ok {
		return nil
	}
	return v
}

func getNumOrDef(x, def float64) float64 {
	if x == 0.0 {
		return def
	}
	return x
}

func newLHContext(parent context.Context, cdata *ConfigData) (context.Context, func(), error) {
	em := equipmentManager.GetEquipmentManager()
	if em == nil {
		return nil, nil, noEquipmentManager
	}
	lh := em.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", nil))
	if lh == nil {
		return nil, nil, noEquipmentManager
	}

	//We are going to configure the liquid handler for a blockId. BlockId will give us the framework and state handling
	// so, for a certain BlockId config options will be aggregated. Liquid Handler will just regenerate all state per
	// this aggregation layer and that will allow us to run multiple protocols.
	//prepare the values
	id := getId(parent)
	config := make(map[string]interface{})
	config["BLOCKID"] = wtype.BlockID{ThreadID: wtype.ThreadID(id)}
	config["MAX_N_PLATES"] = getNumOrDef(cdata.MaxPlates, 4.5)
	config["MAX_N_WELLS"] = getNumOrDef(cdata.MaxWells, 278.0)
	config["RESIDUAL_VOLUME_WEIGHT"] = getNumOrDef(cdata.ResidualVolumeWeight, 1.0)
	// these should come from the paramblock... for now though
	config["INPUT_PLATETYPE"] = "pcrplate_with_cooler"
	config["OUTPUT_PLATETYPE"] = "pcrplate_with_cooler"

	configString, err := json.Marshal(config)
	if err != nil {
		return nil, nil, cannotConfigLh
	}
	if err := lh.Do(*equipment.NewActionDescription(action.LH_CONFIG, string(configString), nil)); err != nil {
		return nil, nil, err
	}

	return context.WithValue(parent, theLhKey, lh),
		func() {
			lh.Do(*equipment.NewActionDescription(action.LH_END, id, nil))
		}, nil
}
