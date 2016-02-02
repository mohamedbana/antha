package execute

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/trace"
)

func Incubate(ctx context.Context, what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) *wtype.LHSolution {
	sol := wtype.NewLHSolution()
	trace.Issue(ctx, trace.MakeIncubate(trace.IncubateOpt{
		BlockID:   getId(ctx),
		OutputSol: sol,
		Temp:      temp,
		Time:      time,
	}, trace.MakeValue(ctx, "", what)))
	return sol
}

func mix(ctx context.Context, opt trace.MixOpt, components []*wtype.LHComponent) *wtype.LHSolution {
	var values []trace.Value
	for _, c := range components {
		values = append(values, trace.MakeValue(ctx, "", c))
	}

	sol := wtype.NewLHSolution()
	sol.BlockID = wtype.NewBlockID(getId(ctx))
	opt.OutputSol = sol

	trace.Issue(ctx, trace.MakeMix(opt, values))

	return sol
}

func Mix(ctx context.Context, components ...*wtype.LHComponent) *wtype.LHSolution {
	return mix(ctx, trace.MixOpt{}, components)
}

func MixInto(ctx context.Context, outplate *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution {
	return mix(ctx, trace.MixOpt{
		OutPlate: outplate,
	}, components)
}

func MixTo(ctx context.Context, outplate *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHSolution {
	return mix(ctx, trace.MixOpt{
		OutPlate: outplate,
		Address:  address,
	}, components)
}
