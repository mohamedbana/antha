package execute

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/trace"
)

func Incubate(ctx context.Context, cmp *wtype.LHComponent, temp wunit.Temperature, time wunit.Time, shaking bool) *wtype.LHComponent {
	trace.Issue(ctx, trace.MakeIncubate(trace.IncubateOpt{
		BlockID:   getId(ctx),
		Component: cmp,
		Temp:      temp,
		Time:      time,
	}, trace.MakeValue(ctx, "", cmp)))
	return cmp
}

func mix(ctx context.Context, opt trace.MixOpt, components []*wtype.LHComponent) *wtype.LHComponent {
	cmp := wtype.NewLHComponent()
	cmp.BlockID = wtype.NewBlockID(getId(ctx))

	var values []trace.Value
	for _, c := range components {
		values = append(values, trace.MakeValue(ctx, "", c))
		cmp.Mix(c)
	}

	opt.OutputCmp = cmp

	trace.Issue(ctx, trace.MakeMix(opt, values))

	return cmp
}

func Mix(ctx context.Context, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, trace.MixOpt{}, components)
}

func MixInto(ctx context.Context, outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, trace.MixOpt{
		OutPlate: outplate,
		Address:  address,
	}, components)
}

func MixTo(ctx context.Context, outplatetype string, address string, platenum int, components ...*wtype.LHComponent) *wtype.LHComponent {
	return mix(ctx, trace.MixOpt{
		PlateType: outplatetype,
		Address:   address,
		PlateNum:  platenum,
	}, components)
}
