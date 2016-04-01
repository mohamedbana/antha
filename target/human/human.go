package human

import (
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

func (a *Human) CanCompile(req ast.Request) bool {
	canMove := true
	mov := len(req.Move) > 0
	mix := req.MixVol != nil
	inc := req.Temp != nil || req.Time != nil

	switch {
	case !canMove && mov:
		return false
	case !a.opt.CanMix && mix:
		return false
	case !a.opt.CanIncubate && inc:
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
	addDep := func(in, dep target.Inst) {
		in.SetDependsOn(append(in.DependsOn(), dep))
	}

	var nodes []ast.Node
	for _, c := range cmds {
		nodes = append(nodes, c)
	}
	g := ast.Deps(nodes)

	entry := &target.Wait{}
	exit := &target.Wait{}
	var insts []target.Inst
	inst := make(map[ast.Node]target.Inst)

	insts = append(insts, entry)
	for i, inum := 0, g.NumNodes(); i < inum; i += 1 {
		n := g.Node(i).(ast.Node)
		in, err := a.makeInst(n.(ast.Command))
		if err != nil {
			return nil, err
		}
		insts = append(insts, in)
		inst[n] = in
	}
	insts = append(insts, exit)

	for i, inum := 0, g.NumNodes(); i < inum; i += 1 {
		n := g.Node(i).(ast.Node)
		in := inst[n]
		for j, jnum := 0, g.NumOuts(n); j < jnum; j += 1 {
			kid := g.Out(n, j).(ast.Node)
			kidIn := inst[kid]
			addDep(in, kidIn)
		}
		addDep(in, entry)
		addDep(exit, in)
	}

	return insts, nil
}

type Opt struct {
	CanMix      bool
	CanIncubate bool
}

func New(opt Opt) *Human {
	return &Human{opt}
}
