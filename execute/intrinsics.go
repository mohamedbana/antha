package execute

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/trace"
	"runtime/debug"
)

type incubateInst struct {
	BlockID      string
	Component    *wtype.LHComponent
	Temp         wunit.Temperature
	Time         wunit.Time
	ShakingForce interface{}
}

func Incubate(ctx context.Context, cmp *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *wtype.LHComponent {
	trace.Issue(ctx, &incubateInst{
		BlockID:   getId(ctx),
		Component: cmp,
		Temp:      temp,
		Time:      time,
	})
	return cmp
}

type mixInst struct {
	Args      []*wtype.LHComponent
	OutputIns *wtype.LHInstruction
	OutputCmp *wtype.LHComponent
	Plate     *wtype.LHPlate
	PlateType string
	Address   string
	PlateNum  int
}

func mix(ctx context.Context, inst *mixInst) *wtype.LHComponent {
	cmp := wtype.NewLHComponent()
	cmp.BlockID = wtype.NewBlockID(getId(ctx))
	for _, c := range inst.Args {
		cmp.Mix(c)
	}
	inst.OutputCmp = cmp

	trace.Issue(ctx, inst)

	return cmp
}

func Mix(ctx context.Context, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, &mixInst{Args: components})
}

func MixInto(ctx context.Context, outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, &mixInst{
		Args:    components,
		Plate:   outplate,
		Address: address,
	})
}

func MixTo(ctx context.Context, outplatetype, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, &mixInst{
		Args:      components,
		PlateType: outplatetype,
		Address:   address,
		PlateNum:  platenum,
	})
}

func runMix(ctx context.Context, blockIds map[wtype.BlockID]bool, inst *mixInst) error {
	toSol := func(inst *mixInst) *wtype.LHInstruction {
		r := mixer.GenericMix(mixer.MixOptions{
			Components:  inst.Args,
			Address:     inst.Address,
			Destination: inst.Plate,
			PlateNum:    inst.PlateNum,
			PlateType:   inst.PlateType,
		})
		r.BlockID = inst.OutputCmp.BlockID
		return r
	}

	if t, err := target.GetTarget(ctx); err != nil {
		return err
	} else if lh, err := t.GetLiquidHandler(); err != nil {
		return err
	} else if rbs, err := json.Marshal(toSol(inst)); err != nil {
		return err
	} else if err := lh.Do(*equipment.NewActionDescription(action.LH_MIX, string(rbs), nil)); err != nil {
		return err
	} else {
		id := inst.OutputCmp.BlockID
		blockIds[id] = true
		return nil
	}
}

func runIn(ctx context.Context, blockIds map[wtype.BlockID]bool, in interface{}) (err error) {
	switch inst := in.(type) {
	case nil:
	case *mixInst:
		err = runMix(ctx, blockIds, inst)
	case *incubateInst:
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

type goError struct {
	BaseError interface{}
	Stack     []byte
}

func (a *goError) Error() string {
	return fmt.Sprintf("%s at:\n%s", a.BaseError, string(a.Stack))
}

// Called by trace to resolve blocked instructions
func resolveIntrinsics(ctx context.Context, insts []interface{}) (ret map[int]interface{}, err error) {
	defer func() {
		if res := recover(); res != nil {
			err = &goError{BaseError: res, Stack: debug.Stack()}
		}
	}()
	blockIds := make(map[wtype.BlockID]bool)
	ret = make(map[int]interface{})

	for idx, in := range insts {
		if err = runIn(ctx, blockIds, in); err != nil {
			return
		}
		ret[idx] = nil
	}

	for id := range blockIds {
		if err = runEnd(ctx, id); err != nil {
			return
		}
	}
	return
}
