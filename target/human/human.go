package human

import (
	"fmt"

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
	var ret []target.Inst
	// XXX
	for _, c := range cmds {
		ret = append(ret, &target.ManualInst{
			Details: fmt.Sprintf("%s", c),
		})
	}

	return ret, nil
}

type Opt struct {
	CanMix bool
}

func New(opt Opt) *Human {
	return &Human{opt}
}
