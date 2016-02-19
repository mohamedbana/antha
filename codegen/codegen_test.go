package codegen

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

type incubator struct{}

func (a *incubator) Can(req ast.Request) bool {
	if req.MixVol != nil {
		return false
	}
	return req.Time != nil || req.Temp != nil
}

func (a *incubator) Compile(insts []ast.Command) ([]target.Inst, error) {
	return []target.Inst{
		&target.ManualInst{
			Details: "hello",
		},
	}, nil
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
	machine.AddDevice(human.New())
	machine.AddDevice(&incubator{})

	if insts, err := Compile(machine, nodes); err != nil {
		t.Fatal(err)
	} else {
		// XXX
		for _, in := range insts {
			fmt.Printf("  %T %+v\n", in, in)
		}
	}
}
