package trace

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/target"
	"runtime/debug"
)

func configMix(inst *MixInst) (*wtype.LHComponent, error) {
	var comps []*wtype.LHComponent
	for _, arg := range inst.Values {
		v := arg.Get()
		if comp, ok := v.(*wtype.LHComponent); !ok {
			return nil, fmt.Errorf("invalid argument to mix %q %T", v, v)
		} else {
			comps = append(comps, comp)
		}
	}
	r := mixer.GenericMix(mixer.MixOptions{
		Components:  comps,
		Instruction: inst.Opt.OutputIns,
		Address:     inst.Opt.Address,
		Destination: inst.Opt.OutPlate,
		PlateNum:    inst.Opt.PlateNum,
		PlateType:   inst.Opt.PlateType,
	})
	return r, nil
}

func runMix(ctx context.Context, blockIds map[wtype.BlockID]bool, inst *MixInst) error {
	if t, err := target.GetTarget(ctx); err != nil {
		return err
	} else if lh, err := t.GetLiquidHandler(); err != nil {
		return err
	} else if r, err := configMix(inst); err != nil {
		return err
	} else if rbs, err := json.Marshal(r); err != nil {
		return err
	} else if err := lh.Do(*equipment.NewActionDescription(action.LH_MIX, string(rbs), nil)); err != nil {
		return err
	} else {
		id := inst.Opt.OutputIns.BlockID
		blockIds[id] = true
		return nil
	}
}

func runIn(ctx context.Context, blockIds map[wtype.BlockID]bool, in instp) (err error) {
	switch inst := in.inst.(type) {
	case nil:
	case *MixInst:
		err = runMix(ctx, blockIds, inst)
	case *IncubateInst:
	case *DebugInst:
	case *NoopInst:
	default:
		err = fmt.Errorf("invalid instruction: %T", inst)
	}
	return err
}

func runEnd(ctx context.Context, blockId wtype.BlockID) error {
	if t, err := target.GetTarget(ctx); err != nil {
		return err
	} else if lh, err := t.GetLiquidHandler(); err != nil {
		return err
	} else {
		return lh.Do(*equipment.NewActionDescription(action.LH_END, blockId.String(), nil))
	}
}

// Run trace instructions.
//
// XXX(ddn): place after code generation.
func run(ctx context.Context, instps []instp) (err error) {
	defer func() {
		if res := recover(); res != nil {
			err = &goError{BaseError: res, Stack: debug.Stack()}
		}
	}()
	blockIds := make(map[wtype.BlockID]bool)
	for _, in := range instps {
		if err = runIn(ctx, blockIds, in); err != nil {
			return
		}
	}

	for id := range blockIds {
		if err = runEnd(ctx, id); err != nil {
			return
		}
	}
	return
}
