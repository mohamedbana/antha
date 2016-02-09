package ast

import (
	"fmt"
	"github.com/antha-lang/antha/graph"
)

const (
	AllDeps = iota
	DataDeps
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
type Node interface {
	graph.Node
	NodeString() string
}

// Use an external value
type UseExpr struct {
	Desc string // Description only used during pretty printing
}

func (a *UseExpr) NodeString() string {
	return fmt.Sprintf("%+v", struct {
		Desc interface{}
	}{
		Desc: a.Desc,
	})
}

// Unordered set of nodes
type BundleExpr struct {
	From []Node
}

func (a *BundleExpr) NodeString() string {
	return ""
}

// Ordered set of nodes
type ListExpr struct {
	From []Node
}

func (a *ListExpr) NodeString() string {
	return ""
}

// Create a BundleExpr of Nodes with matching Keys subject to dependencies
type GatherExpr struct {
	Key  interface{} // If nil, gather maximal subject to Gathers with non-nil keys
	From Node
}

func (a *GatherExpr) NodeString() string {
	return fmt.Sprintf("%+v", struct {
		Key interface{}
	}{
		Key: a.Key,
	})
}

// Create a new value.
type NewExpr struct {
	From Node
}

func (a *NewExpr) NodeString() string {
	return ""
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
	From Node        // Input value
	Prev Node        // Optional control dependence
	Near Node        // Optional temporal dependence
	Func string      // Name of function to call
	Gen  interface{} // Generation this application belongs to
}

func (a *ApplyExpr) NodeString() string {
	return fmt.Sprintf("%+v", struct {
		Func, Gen interface{}
	}{
		Func: a.Func,
		Gen:  a.Gen,
	})
}

type Graph struct {
	Nodes     []Node
	whichDeps int
}

func (a *Graph) NumNodes() int {
	return len(a.Nodes)
}

func (a *Graph) Node(i int) graph.Node {
	return a.Nodes[i]
}

// Return subset of nodes that match the predicate
func matching(pred func(Node) bool, nodes ...Node) (r []Node) {
	for _, n := range nodes {
		if !pred(n) {
			continue
		}
		r = append(r, n)
	}
	return
}

func notNil(n Node) bool {
	return n != nil
}

func getOut(n Node, i, deps int) Node {
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
		case AllDeps:
			return matching(notNil, n.From, n.Prev, n.Near)[i]
		case DataDeps:
			return matching(notNil, n.From)[i]
		default:
			panic(fmt.Sprintf("codegen.getOut: unknown dep type %d", deps))
		}
	default:
		panic(fmt.Sprintf("codegen.getOut: unknown node type %T", n))
	}
}

func numOuts(n Node, deps int) int {
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
		case AllDeps:
			return len(matching(notNil, n.From, n.Prev, n.Near))
		case DataDeps:
			return len(matching(notNil, n.From))
		default:
			panic(fmt.Sprintf("codegen.numOuts: unknown dep type %d", deps))
		}
	default:
		panic(fmt.Sprintf("codegen.numOuts: unknown node type %T", n))
	}
}

func (a *Graph) NumOuts(n graph.Node) int {
	return numOuts(n.(Node), a.whichDeps)
}

func (a *Graph) Out(n graph.Node, i int) graph.Node {
	return getOut(n.(Node), i, a.whichDeps)
}

type ToGraphOpt struct {
	Root      Node
	WhichDeps int
}

func ToGraph(opt ToGraphOpt) *Graph {
	g := &Graph{
		whichDeps: opt.WhichDeps,
	}

	// Traverse doesn't use Graph.NumNodes() or Graph.Node(int), so we can pass
	// in our partially constructed graph to extract the reachable nodes in the
	// AST
	results, _ := graph.Visit(graph.VisitOpt{Graph: g, Root: opt.Root})

	for k := range results.Seen.Range() {
		g.Nodes = append(g.Nodes, k.(Node))
	}

	return g
}
