// Package codegen compiles generic instructions to target-specific ones.
// Target, in this case, is some combination of devices (e.g., two
// ExtendedLiquidHandlers and human plate mover).
package codegen

import (
	"errors"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

var (
	tbd           = errors.New("not yet implemented")
	notInjective  = errors.New("relation not injective")
	notFunctional = errors.New("relation not a function")
)

type ir struct {
	Root     AstNode
	Graph    graph.Graph
	DataTree graph.Graph // Data dependencies
}

func check(root AstNode) (*ir, error) {
	if r, err := normalize(root); err != nil {
		return nil, err
	} else {
		root = r
	}

	// Check that data dependencies are single-consumer (because they model
	// destructive update of values)
	dg := toGraph(toGraphOpt{
		Root:      root,
		WhichDeps: dataDeps,
	})
	if err := graph.IsTree(dg, root); err != nil {
		return nil, err
	}

	g := toGraph(toGraphOpt{Root: root})
	if err := graph.IsDag(coalesceByGen(g)); err != nil {
		return nil, err
	}

	return &ir{Root: root, Graph: g, DataTree: dg}, nil
}

// Return graph with apply expressions in the same generation as a single node
func coalesceByGen(g *AstGraph) graph.QGraph {
	return graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: g,
		Colorer: func(n graph.Node) interface{} {
			return n.(*ApplyExpr).Gen
		},
		HasColor: func(n graph.Node) bool {
			switch n := n.(AstNode).(type) {
			case *ApplyExpr:
				return n.Gen != nil
			}
			return false
		},
	})
}

func mergeGathers(ir *ir) (int, error) {
	return 0, tbd
}

// Compile an expression program into a sequence of instructions for a target
// configuration.
func Compile(t *target.Target, root AstNode) ([]Inst, error) {
	if root == nil {
		return nil, nil
	}
	if ir, err := check(root); err != nil {
		return nil, err
	} else if _, err := mergeGathers(ir); err != nil {
		return nil, err
	} else {
		return nil, tbd
	}
}
