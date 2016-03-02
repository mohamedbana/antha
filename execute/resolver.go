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
	insts []target.Inst
	nodes []ast.Node
	comp  map[*wtype.LHComponent]*ast.UseComp
}

func (a *resolver) makeComp(c *wtype.LHComponent) *ast.UseComp {
	if a.comp == nil {
		a.comp = make(map[*wtype.LHComponent]*ast.UseComp)
	}

	n, ok := a.comp[c]
	if !ok {
		n = &ast.UseComp{
			Value: c,
		}
		a.comp[c] = n
	}
	return n
}

func (a *resolver) addIncubate(in *incubateInst) {
	in.Node.From = append(in.Node.From, a.makeComp(in.Arg))

	out := a.makeComp(in.Comp)
	out.From = append(out.From, in.Node)

	a.nodes = append(a.nodes, out)
}

func (a *resolver) addMix(in *mixInst) {
	for _, arg := range in.Args {
		in.Node.From = append(in.Node.From, a.makeComp(arg))
	}

	out := a.makeComp(in.Comp)
	out.From = append(out.From, in.Node)

	a.nodes = append(a.nodes, out)
}

// Called by trace to resolve blocked instructions
func (a *resolver) resolve(ctx context.Context, insts []interface{}) (map[int]interface{}, error) {
	ret := make(map[int]interface{})
	for idx, in := range insts {
		switch inst := in.(type) {
		case nil:
		case *incubateInst:
			a.addIncubate(inst)
		case *mixInst:
			a.addMix(inst)
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
