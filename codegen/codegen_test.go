package codegen

import (
	"testing"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

type incubateInst struct {
	Depends []target.Inst
}

func (a *incubateInst) Device() target.Device {
	return nil
}

func (a *incubateInst) DependsOn() []target.Inst {
	return a.Depends
}

func (a *incubateInst) SetDependsOn(xs []target.Inst) {
	a.Depends = xs
}

func (a *incubateInst) GetTimeEstimate() float64 {
	return 0.0
}

type incubator struct{}

func (a *incubator) CanCompile(req ast.Request) bool {
	if req.MixVol != nil {
		return false
	}
	return req.Time != nil || req.Temp != nil
}

func (a *incubator) Compile(insts []ast.Command) ([]target.Inst, error) {
	return []target.Inst{&incubateInst{}}, nil
}

func (a *incubator) MoveCost(from target.Device) int {
	if a == from {
		return 0
	}
	return human.HumanByXCost - 1
}

func (a *incubator) String() string {
	return "Incubator"
}

func TestWellFormed(t *testing.T) {
	var nodes []ast.Node
	for idx := 0; idx < 4; idx += 1 {
		m := &ast.Mix{
			Reqs: []ast.Request{
				ast.Request{
					MixVol: ast.NewInterval(0.1, 1.0),
				},
			},
			From: []ast.Node{
				&ast.UseComp{},
				&ast.UseComp{},
				&ast.UseComp{},
			},
		}
		u := &ast.UseComp{}
		u.From = append(u.From, m)

		i := &ast.Incubate{
			Reqs: []ast.Request{
				ast.Request{
					Temp: ast.NewPoint(25),
					Time: ast.NewPoint(60 * 60),
				},
			},
			From: []ast.Node{u},
		}

		nodes = append(nodes, i)
	}

	machine := target.New()
	if err := machine.AddDevice(human.New(human.Opt{CanMix: true})); err != nil {
		t.Fatal(err)
	}
	if err := machine.AddDevice(&incubator{}); err != nil {
		t.Fatal(err)
	}

	if insts, err := Compile(machine, nodes); err != nil {
		t.Fatal(err)
	} else if l := len(insts); l == 0 {
		t.Errorf("expected > %d instructions found %d", 0, l)
	} else if last, ok := insts[l-1].(*incubateInst); !ok {
		t.Errorf("expected incubateInst found %T", insts[l-1])
	} else if n := len(last.Depends); n != 1 {
		t.Errorf("expected %d dependencies found %d", 1, n)
	}
}
