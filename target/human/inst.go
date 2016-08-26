package human

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
)

func extractFromUseNodes(nodes ...*ast.UseComp) string {
	var vs []string
	for _, n := range nodes {
		if n.Value == nil {
			vs = append(vs, "<nil>")
		} else {
			vs = append(vs, n.Value.Name())
		}
	}
	return strings.Join(vs, ",")
}

func extractFromNodes(nodes ...ast.Node) string {
	var vs []string
	for _, n := range nodes {
		switch c := n.(type) {
		case *ast.UseComp:
			vs = append(vs, extractFromUseNodes(c))
		case *ast.Move:
			vs = append(vs, extractFromUseNodes(c.From...))
		default:
			panic(fmt.Sprintf("human.Human: unknown node %T", c))
		}
	}
	return strings.Join(vs, ",")
}

func (a *Human) makeFromMove(c *ast.Move) (*target.Manual, error) {
	from := extractFromUseNodes(c.From...)
	return &target.Manual{
		Dev:     a,
		Label:   "Move",
		Details: fmt.Sprintf("Move %q to %s", from, c.ToLoc),
	}, nil
}

func (a *Human) makeFromCommand(c *ast.Command) (*target.Manual, error) {
	from := extractFromNodes(c.From...)
	switch inst := c.Inst.(type) {
	case *wtype.LHInstruction:
		return &target.Manual{
			Dev:     a,
			Label:   "Mix",
			Details: fmt.Sprintf("Mix %q", from),
		}, nil
	case *ast.IncubateInst:
		return &target.Manual{
			Dev:     a,
			Label:   "Incubate",
			Details: fmt.Sprintf("Incubate %q at %s for %s", from, inst.Temp.ToString(), inst.Time.ToString()),
			Time:    inst.Time.Seconds(),
		}, nil
	case *ast.HandleInst:
		return &target.Manual{
			Dev:     a,
			Label:   "Handle",
			Details: inst.Group,
		}, nil
	default:
		return nil, fmt.Errorf("unknown inst %T", inst)
	}
}

func (a *Human) makeFromManual(ms []*target.Manual) target.Inst {
	m := ms[0]
	if len(ms) == 1 {
		return m
	}
	var details []string
	var maxSec float64

	for _, m := range ms {
		if t := m.GetTimeEstimate(); maxSec < t {
			maxSec = t
		}
		details = append(details, m.Details)
	}

	return &target.Manual{
		Dev:     m.Dev,
		Label:   m.Label,
		Details: strings.Join(details, "\n"),
		Time:    maxSec,
	}
}

func (a *Human) makeInst(n ast.Node) (*target.Manual, error) {
	switch n := n.(type) {
	case *ast.Move:
		return a.makeFromMove(n)
	case *ast.Command:
		return a.makeFromCommand(n)
	default:
		return nil, fmt.Errorf("unknown node %T", n)
	}
}
