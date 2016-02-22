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

	// TODO: parallelize and vectorize
	entry := &target.Wait{}
	exit := &target.Wait{}
	var insts []target.Inst

	insts = append(insts, entry)
	exit.Depends = append(exit.Depends, entry)

	for _, c := range cmds {
		switch c := c.(type) {
		case *ast.Move:
			insts = append(insts, &target.Manual{
				Depends: []target.Inst{entry},
				Details: fmt.Sprintf("MOVE %q", extract(c)),
			})
		case *ast.Mix:
			insts = append(insts, &target.Manual{
				Depends: []target.Inst{entry},
				Details: fmt.Sprintf("MIX %q", c.Inst),
			})
		case *ast.Incubate:
			insts = append(insts, &target.Manual{
				Depends: []target.Inst{entry},
				Details: fmt.Sprintf("INCUBATE %q", c),
			})
		default:
			return nil, fmt.Errorf("unknown command %T", c)
		}
		exit.Depends = append(exit.Depends, insts[len(insts)-1])
	}

	insts = append(insts, exit)

	return insts, nil
}

type Opt struct {
	CanMix bool
}

func New(opt Opt) *Human {
	return &Human{opt}
}
