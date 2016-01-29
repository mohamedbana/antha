// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"encoding/json"
	"errors"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

var (
	cannotConfigure = errors.New("cannot configure liquid handler")

	defaultMaxPlates            = 4.5
	defaultMaxWells             = 278.0
	defaultResidualVolumeWeight = 1.0
	defaultWellByWell           = false
	DefaultConfig               = Config{
		MaxPlates:            &defaultMaxPlates,
		MaxWells:             &defaultMaxWells,
		ResidualVolumeWeight: &defaultResidualVolumeWeight,
		InputPlateType:       []string{"pcrplate_skirted"},
		OutputPlateType:      []string{"pcrplate_skirted"},
		WellByWell:           &defaultWellByWell,
	}
)

type Options struct {
	WorkflowData []byte         // JSON data describing workflow
	Workflow     *workflow.Desc // Or workflow directly
	ParamData    []byte         // JSON data describing parameters
	Params       *RawParams     // Or parameters directly
	Target       *target.Target // Target machine configuration
	Id           string         // Job Id
	Config       *Config        // Override config data in ParamData
}

// Simple entrypoint for one-shot execution of workflows.
func Run(parent context.Context, opt Options) (*workflow.Workflow, error) {
	w, err := workflow.New(workflow.Options{FromBytes: opt.WorkflowData, FromDesc: opt.Workflow})
	if err != nil {
		return nil, err
	}

	var params *RawParams
	if opt.Params != nil {
		params = opt.Params
	} else if opt.ParamData != nil {
		if err := json.Unmarshal(opt.ParamData, &params); err != nil {
			return nil, err
		}
	}

	cd, err := setParams(parent, params, w)
	if err != nil {
		return nil, err
	}

	ctx := target.WithTarget(WithId(parent, opt.Id), opt.Target)

	lh, err := opt.Target.GetLiquidHandler()
	if err != nil {
		return nil, err
	}
	if err := config(ctx, lh, DefaultConfig.Merge(cd).Merge(opt.Config)); err != nil {
		return nil, err
	}

	if err := w.Run(ctx); err != nil {
		return nil, err
	}

	return w, nil
}

// XXX Move out when equipment config is done
func config(parent context.Context, lh equipment.Equipment, cd Config) error {
	id := getId(parent)
	config := make(map[string]interface{})
	config["BLOCKID"] = wtype.NewBlockID(id)
	config["MAX_N_PLATES"] = *cd.MaxPlates
	config["MAX_N_WELLS"] = *cd.MaxWells
	config["RESIDUAL_VOLUME_WEIGHT"] = *cd.ResidualVolumeWeight
	config["INPUT_PLATETYPE"] = cd.InputPlateType
	config["OUTPUT_PLATETYPE"] = cd.OutputPlateType
	config["WELLBYWELL"] = *cd.WellByWell

	configString, err := json.Marshal(config)
	if err != nil {
		return cannotConfigure
	}
	if err := lh.Do(*equipment.NewActionDescription(action.LH_CONFIG, string(configString), nil)); err != nil {
		return cannotConfigure
	}
	return nil
}
