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

// What is the set of UseComps that reach each command
func buildReachingUses(g *ast.Graph) map[ast.Node][]*ast.UseComp {
	// Simple fixpoint:
	//   Value: set of use comps,
	//   Merge: union
	//   Transfer functions:
	//     - Command c -> { }
	//     - UseComp u -> {u}

	values := make(map[ast.Node][]*ast.UseComp)

	merge := func(n ast.Node) []*ast.UseComp {
		var vs []*ast.UseComp
		for i, inum := 0, g.NumOuts(n); i < inum; i += 1 {
			pred := g.Out(n, i).(ast.Node)
			switch pred := pred.(type) {
			case *ast.Command:
				// Kill
			case *ast.UseComp:
				vs = append(vs, values[pred]...)
				vs = append(vs, pred)
			default:
				vs = append(vs, values[pred]...)
			}
		}
		return vs
	}

	dag := graph.Schedule(graph.Reverse(g))

	for len(dag.Roots) > 0 {
		var next []graph.Node
		for _, n := range dag.Roots {
			n := n.(ast.Node)
			seen := make(map[*ast.UseComp]bool)

			for _, v := range merge(n) {
				if seen[v] {
					continue
				}
				seen[v] = true
				values[n] = append(values[n], v)
			}

			next = append(next, dag.Visit(n)...)
		}

		dag.Roots = next
	}

	return values
}

// Build IR
func build(root ast.Node) (*ir, error) {
	g := ast.ToGraph(ast.ToGraphOpt{
		Roots: []ast.Node{root},
	})

	ct := graph.Eliminate(graph.EliminateOpt{
		Graph: g,
		In: func(n graph.Node) bool {
			c, ok := n.(*ast.Command)
			return (ok && c.Output == nil) || n == root
		},
	})

	// TODO: Add back some validity checks like the same UseComp cannot be used
	// multiple times

	return &ir{
		Root:         root,
		Graph:        g,
		Commands:     ct,
		reachingUses: buildReachingUses(g),
	}, nil
}
