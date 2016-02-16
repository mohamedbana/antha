package codegen

import (
	"fmt"
	"github.com/antha-lang/antha/graph"
)

const (
	allDeps = iota
	dataDeps
)

// Input to code generation. An abstract syntax tree generated via execution of
// an Antha protocol.
//
// The basic design philosophy is to capture the semantics of the Antha
// language while reducing the cases for code generation. A secondary goal is
// to ease the creation of the AST at runtime (e.g., online, incremental
// generation of nodes).
//
// Conveniently, a tree naturally expresses the single-use (i.e., destructive
// update) aspect of physical things, so the code generation keeps this
// representation longer than a traditional compiler flow would.
type AstNode interface {
	graph.Node
}

// Use an external value
type UseExpr struct {
	Desc string // Description only used during pretty printing
}

// Unordered set of nodes
type BundleExpr struct {
	From []AstNode
}

// Ordered set of nodes
type ListExpr struct {
	From []AstNode
}

// Create a BundleExpr of AstNodes with matching Keys subject to dependencies
type GatherExpr struct {
	Key  interface{} // If nil, gather maximal subject to Gathers with non-nil keys
	From AstNode
}

// Create a new value.
type NewExpr struct {
	From AstNode
}

// Apply some intrinsic function.
//
// Sequencing between applications can be expressed in four ways:
//   1. Data dependency on From node
//   2. Control dependency on Prev node
//   3. Temporal dependency on Near node
//   4. Generation color: all the applications from one generation are done
//      before any subsequent generations
type ApplyExpr struct {
	Opt  interface{} // Function-specific options to pass to code gen
	From AstNode     // Input value
	Prev AstNode     // Optional control dependence
	Near AstNode     // Optional temporal dependence
	Func string      // Name of function to call
	Gen  interface{} // Generation this application belongs to
}

type AstGraph struct {
	Nodes     []AstNode
	whichDeps int
}

func (a *AstGraph) NumNodes() int {
	return len(a.Nodes)
}

func (a *AstGraph) Node(i int) graph.Node {
	return a.Nodes[i]
}

// Return subset of nodes that match the predicate
func matching(pred func(AstNode) bool, nodes ...AstNode) (r []AstNode) {
	for _, n := range nodes {
		if !pred(n) {
			continue
		}
		r = append(r, n)
	}
	return
}

func notNil(n AstNode) bool {
	return n != nil
}

func getOut(n AstNode, i, deps int) AstNode {
	switch n := n.(type) {
	case *UseExpr:
		return nil
	case *BundleExpr:
		return n.From[i]
	case *ListExpr:
		return n.From[i]
	case *GatherExpr:
		return n.From
	case *NewExpr:
		return n.From
	case *ApplyExpr:
		switch deps {
		case allDeps:
			return matching(notNil, n.From, n.Prev, n.Near)[i]
		case dataDeps:
			return matching(notNil, n.From)[i]
		default:
			panic(fmt.Sprintf("codegen.getOut: unknown dep type %d", deps))
		}
	default:
		panic(fmt.Sprintf("codegen.getOut: unknown node type %T", n))
	}
}

func numOuts(n AstNode, deps int) int {
	switch n := n.(type) {
	case *UseExpr:
		return 0
	case *BundleExpr:
		return len(n.From)
	case *ListExpr:
		return len(n.From)
	case *GatherExpr:
		return len(matching(notNil, n.From))
	case *NewExpr:
		return len(matching(notNil, n.From))
	case *ApplyExpr:
		switch deps {
		case allDeps:
			return len(matching(notNil, n.From, n.Prev, n.Near))
		case dataDeps:
			return len(matching(notNil, n.From))
		default:
			panic(fmt.Sprintf("codegen.numOuts: unknown dep type %d", deps))
		}
	default:
		panic(fmt.Sprintf("codegen.numOuts: unknown node type %T", n))
	}
}

func (a *AstGraph) NumOuts(n graph.Node) int {
	return numOuts(n.(AstNode), a.whichDeps)
}

func (a *AstGraph) Out(n graph.Node, i int) graph.Node {
	return getOut(n.(AstNode), i, a.whichDeps)
}

type toGraphOpt struct {
	Root      AstNode
	WhichDeps int
}

func toGraph(opt toGraphOpt) *AstGraph {
	g := &AstGraph{
		whichDeps: opt.WhichDeps,
	}

	// Traverse doesn't use Graph.NumNodes() or Graph.Node(int), so we can pass
	// in our partially constructed graph to extract the reachable nodes in the
	// AST
	results, _ := graph.Visit(graph.VisitOpt{Graph: g, Root: opt.Root})

	for k := range results.Seen.Range() {
		g.Nodes = append(g.Nodes, k)
	}

	return g
}
