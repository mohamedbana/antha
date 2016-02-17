package execute

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
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

	toInst := func(inst *mixInst) *wtype.LHInstruction {
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

	ret = make(map[int]interface{})

	var mixes []*wtype.LHInstruction
	for idx, in := range insts {
		switch inst := in.(type) {
		case nil:
		case *incubateInst:
		case *mixInst:
			mixes = append(mixes, toInst(inst))
		default:
			err = fmt.Errorf("invalid instruction: %T", inst)
			return
		}

		ret[idx] = nil
	}

	var t *target.Target
	if t, err = target.GetTarget(ctx); err != nil {
		return
	} else if err = t.Mix(mixes); err != nil {
		return
	}

	return
}
