package execute

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/logger"
	"golang.org/x/net/context"
)

func Incubate(ctx context.Context, what *wtype.LHSolution, temp wunit.Temperature, time wunit.Time, shaking bool) {
	logger.Debug(fmt.Sprintln("INCUBATE: ", temp.ToString(), " ", time.ToString(), " shaking? ", shaking))
}

func MixInto(ctx context.Context, outplate *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution {
	reaction := mixer.MixInto(outplate, components...)
	reaction.BlockID = wtype.BlockID{ThreadID: wtype.ThreadID(getId(ctx))}
	// XXX(ddn): if needed use arguments on params instead
	//reaction.SName = w.getString("OutputReactionName")
	if reqReaction, err := json.Marshal(reaction); err != nil {
		panic(fmt.Sprintf("error coding reaction data, %v", err))
	} else if err := getLh(ctx).Do(*equipment.NewActionDescription(action.LH_MIX, string(reqReaction), nil)); err != nil {
		panic(fmt.Sprintf("error running liquid handling request: %s", err))
	}
	return reaction
}
