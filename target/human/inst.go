package human

import (
	"fmt"
	"strings"

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

func (a *Human) makeFromMove(c *ast.Move) target.Inst {
	from := extractFromUseNodes(c.From...)
	return &target.Manual{
		Dev:     a,
		Label:   "Move",
		Details: fmt.Sprintf("Move %q to %s", from, c.ToLoc),
	}
}

func (a *Human) makeFromMix(c *ast.Mix) target.Inst {
	from := extractFromNodes(c.From...)
	return &target.Manual{
		Dev:     a,
		Label:   "Mix",
		Details: fmt.Sprintf("Mix %q", from),
	}
}

func (a *Human) makeFromIncubate(c *ast.Incubate) target.Inst {
	from := extractFromNodes(c.From...)
	return &target.Manual{
		Dev:     a,
		Label:   "Incubate",
		Details: fmt.Sprintf("Incubate %q at %s for %s", from, c.Temp.ToString(), c.Time.ToString()),
	}
}

func (a *Human) makeInst(cmd ast.Command) (target.Inst, error) {
	switch cmd := cmd.(type) {
	case *ast.Move:
		return a.makeFromMove(cmd), nil
	case *ast.Mix:
		return a.makeFromMix(cmd), nil
	case *ast.Incubate:
		return a.makeFromIncubate(cmd), nil
	default:
		return nil, fmt.Errorf("unknown command %T", cmd)
	}
}
