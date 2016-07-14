package codegen

import (
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
)

// Build rooted graph
func makeRoot(nodes []ast.Node) (ast.Node, error) {
	someNode := func(g graph.Graph, m map[graph.Node]bool) graph.Node {
		for i, inum := 0, g.NumNodes(); i < inum; i += 1 {
			n := g.Node(i)
			if !m[n] {
				return n
			}
		}
		return nil
	}

	g := ast.ToGraph(ast.ToGraphOpt{
		Roots: nodes,
	})

	roots := graph.Schedule(g).Roots
	seen := make(map[graph.Node]bool)
	for _, root := range roots {
		results, _ := graph.Visit(graph.VisitOpt{
			Graph: g,
			Root:  root,
			Visitor: func(n graph.Node) error {
				if seen[n] {
					return graph.NextNode
				}
				return nil
			},
		})
		for _, k := range results.Seen.Values() {
			seen[k] = true
		}
	}

	// If some nodes are not reachable from roots, there must be a cycle
	if len(seen) != g.NumNodes() {
		return nil, fmt.Errorf("cycle containing %q", someNode(g, seen))
	}

	ret := &ast.Bundle{}
	for _, r := range roots {
		ret.From = append(ret.From, r.(ast.Node))
	}
	return ret, nil
}

// Build IR
func build(root ast.Node) (*ir, error) {
	g := ast.ToGraph(ast.ToGraphOpt{
		Roots: []ast.Node{root},
	})

	// Check that data dependencies are single-consumer (because they model
	// destructive update of values)
	if err := graph.IsTree(g, root); err != nil {
		return nil, err
	}

	for i, inum := 0, g.NumNodes(); i < inum; i += 1 {
		switch n := g.Node(i).(type) {
		case *ast.UseComp:
			if len(n.From) > 1 {
				return nil, fmt.Errorf("component %q created multiple times", n)
			}
		default:
		}
	}

	ct := graph.Eliminate(graph.EliminateOpt{
		Graph: g,
		In: func(n graph.Node) bool {
			c, ok := n.(ast.Command)
			return (ok && c.Output() == nil) || n == root
		},
	})
	if err := graph.IsTree(ct, root); err != nil {
		return nil, err
	}

	return &ir{
		Root:        root,
		Graph:       g,
		CommandTree: ct,
	}, nil
}
