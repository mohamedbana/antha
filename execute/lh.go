package execute

import (
	"encoding/json"
	"errors"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipmentManager"
	"github.com/antha-lang/antha/target"
)

var (
	cannotConfigLh = errors.New("cannot configure liquid handler")
	noLh           = errors.New("no liquid handler found")
)

func getLhFromEm(em equipmentManager.EquipmentManager) (equipment.Equipment, error) {
	lh := em.GetActionCandidate(*equipment.NewActionDescription(action.LH_MIX, "", nil))
	if lh == nil {
		return nil, noLh
	}
	return lh, nil
}

func getNumOrDef(x, def float64) float64 {
	var defv float64
	if x == defv {
		return def
	}
	return x
}

func getStrOrDef(x, def string) string {
	var defv string
	if x == defv {
		return def
	}
	return x
}

func newLHContext(parent context.Context, lh equipment.Equipment, cdata *ConfigData) (context.Context, func(), error) {
	// We are going to configure the liquid handler for a blockId. BlockId will
	// give us the framework and state handling so, for a certain BlockId
	// config options will be aggregated. Liquid Handler will just regenerate
	// all state per this aggregation layer and that will allow us to run
	// multiple protocols.
	id := getId(parent)
	// XXX: move to trace/run.go
	config := make(map[string]interface{})
	config["BLOCKID"] = wtype.BlockID{ThreadID: wtype.ThreadID(id)}
	config["MAX_N_PLATES"] = getNumOrDef(cdata.MaxPlates, 4.5)
	config["MAX_N_WELLS"] = getNumOrDef(cdata.MaxWells, 278.0)
	config["RESIDUAL_VOLUME_WEIGHT"] = getNumOrDef(cdata.ResidualVolumeWeight, 1.0)
	config["INPUT_PLATETYPE"] = getStrOrDef(cdata.InputPlateType, "pcrplate_with_cooler")
	config["OUTPUT_PLATETYPE"] = getStrOrDef(cdata.OutputPlateType, "pcrplate_with_cooler")

	configString, err := json.Marshal(config)
	if err != nil {
		return nil, nil, cannotConfigLh
	}
	if err := lh.Do(*equipment.NewActionDescription(action.LH_CONFIG, string(configString), nil)); err != nil {
		return nil, nil, err
	}

	t := &target.Target{}
	t.AddLiquidHandler(lh)

	return target.WithTarget(parent, t), func() {}, nil
}
