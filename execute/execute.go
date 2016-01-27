// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

type Options struct {
	WorkflowData []byte
	ParamData    []byte
	Target       *target.Target
	Id           string
}

// Simple entrypoint for one-shot execution of workflows.
func Run(parent context.Context, opt Options) (*workflow.Workflow, error) {
	w, err := workflow.New(workflow.Options{FromBytes: opt.WorkflowData})
	if err != nil {
		return nil, err
	}

	cd, err := setParams(parent, opt.ParamData, w)
	if err != nil {
		return nil, fmt.Errorf("cannot set initial parameters: %s", err)
	}

	ctx := target.WithTarget(WithId(parent, opt.Id), opt.Target)

	lh, err := opt.Target.GetLiquidHandler()
	if err != nil {
		return nil, err
	}
	if err := config(ctx, lh, cd); err != nil {
		return nil, err
	}

	if err := w.Run(ctx); err != nil {
		return nil, err
	}

	return w, nil
}

// XXX Move out when equipment config is done
func config(parent context.Context, lh equipment.Equipment, cd *ConfigData) error {
	// XXX: move to trace/run.go
	type cvalue struct {
		Key  string
		UseA bool
		A    interface{}
		B    interface{}
	}

	id := getId(parent)
	cvalues := []cvalue{
		cvalue{Key: "BLOCKID", A: wtype.NewBlockID(id), UseA: true},
		cvalue{Key: "MAX_N_PLATES", A: 4.5, B: cd.MaxPlates, UseA: cd.MaxPlates == 0.0},
		cvalue{Key: "MAX_N_WELLS", A: 278.0, B: cd.MaxWells, UseA: cd.MaxWells == 0.0},
		cvalue{Key: "RESIDUAL_VOLUME_WEIGHT", A: 1.0, B: cd.ResidualVolumeWeight, UseA: cd.ResidualVolumeWeight == 0.0},
		cvalue{Key: "INPUT_PLATETYPE", A: []string{"pcrplate_skirted"}, B: cd.InputPlateType, UseA: len(cd.InputPlateType) == 0},
		cvalue{Key: "OUTPUT_PLATETYPE", A: []string{"pcrplate_skirted"}, B: cd.OutputPlateType, UseA: len(cd.OutputPlateType) == 0},
		cvalue{Key: "WELLBYWELL", A: false, B: cd.WellByWell, UseA: false},
	}

	config := make(map[string]interface{})
	for _, cv := range cvalues {
		if cv.UseA {
			config[cv.Key] = cv.A

			if cv.Key == "INPUT_PLATETYPE" {
				logger.Info(fmt.Sprint("WARNING: No input plate types specified, reverting to default: ", cv.A))
			} else if cv.Key == "OUTPUT_PLATETYPE" {
				logger.Info(fmt.Sprint("WARNING: No output plate types specified, reverting to default: ", cv.A))
			}

		} else {
			config[cv.Key] = cv.B
		}
	}

	configString, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("cannot configure")
	}
	if err := lh.Do(*equipment.NewActionDescription(action.LH_CONFIG, string(configString), nil)); err != nil {
		return fmt.Errorf("cannot configure")
	}
	return nil
}
