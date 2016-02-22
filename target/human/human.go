package human

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
)

const (
	HumanByHumanCost = 50  // Cost of manually moving from another human device
	HumanByXCost     = 100 // Cost of manually moving from any non-human device
)

var (
	_ target.Device = &Human{}
)

type Human struct {
	opt Opt
}

func (a *Human) Can(req ast.Request) bool {
	if !a.opt.CanMix && req.MixVol != nil {
		return false
	}
	return true
}

func (a *Human) MoveCost(from target.Device) int {
	if _, ok := from.(*Human); ok {
		return HumanByHumanCost
	}
	return HumanByXCost
}

func (a *Human) String() string {
	return "Human"
}

func (a *Human) Compile(cmds []ast.Command) ([]target.Inst, error) {
	extract := func(m *ast.Move) (r []*wtype.LHComponent) {
		for _, u := range m.From {
			r = append(r, u.Value)
		}
		return
	}

	// TODO bulk up moves and vectorize mixes

	var ret []target.Inst
	for _, c := range cmds {
		switch c := c.(type) {
		case *ast.Move:
			ret = append(ret, &target.Move{
				Comps: extract(c),
			})
		default:
			ret = append(ret, &target.Manual{
				Details: fmt.Sprintf("%s", c),
			})
		}
	}

	for i, n := 1, len(ret); i < n; i += 1 {
		ret[i].SetDependsOn([]target.Inst{ret[i-1]})
	}

	return ret, nil
}

type Opt struct {
	CanMix bool
}

func New(opt Opt) *Human {
	return &Human{opt}
}
