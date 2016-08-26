package execute

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/target"
)

// Converts execute instructions to their ast equivalents
type resolver struct {
	insts   []target.Inst
	nodes   []ast.Node
	useComp map[string]*ast.UseComp
}

func (a *resolver) makeComp(c *wtype.LHComponent) *ast.UseComp {
	if a.useComp == nil {
		a.useComp = make(map[string]*ast.UseComp)
	}

	n, ok := a.useComp[c.ID]
	if !ok {
		n = &ast.UseComp{
			Value: c,
		}
		a.useComp[c.ID] = n
	}
	return n
}

func (a *resolver) addCommand(in *commandInst) {
	for _, arg := range in.Args {
		in.Command.From = append(in.Command.From, a.makeComp(arg))
	}

	out := a.makeComp(in.Comp)
	out.From = append(out.From, in.Command)

	a.nodes = append(a.nodes, out)
}

// Called by trace to resolve blocked instructions
func (a *resolver) resolve(ctx context.Context, insts []interface{}) (map[int]interface{}, error) {
	ret := make(map[int]interface{})
	for idx, in := range insts {
		switch inst := in.(type) {
		case nil:
		case *commandInst:
			a.addCommand(inst)
		default:
			return nil, fmt.Errorf("invalid instruction: %T", inst)
		}

		// TODO(ddn): Under prefix execution, we need to notify someone to
		// execute instructions to unblock the resolver. Now, since we don't do
		// prefix execution, we can just mark as done.
		ret[idx] = nil
	}

	if t, err := target.GetTarget(ctx); err != nil {
		return nil, err
	} else if insts, err := codegen.Compile(t, a.nodes); err != nil {
		return nil, err
	} else {
		a.insts = append(a.insts, insts...)
		return ret, nil
	}
}
