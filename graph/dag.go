package graph

import (
	"fmt"
)

type topoOrder struct {
	Graph Graph
	Order []Node  // Order in topographical sort
	Cycle []Node  // If not DAG, which nodes participate in a cycle
	black nodeSet // Fully processed nodes
	gray  nodeSet // Nodes currently being processed
}

// TODO(ddn): For performance consider avoiding recursion

// Visit leaves, then nodes next to leaves, etc.
func (a *topoOrder) visit(n Node) {
	if a.gray[n] {
		if !a.black[n] {
			a.black[n] = true
			a.Cycle = append(a.Cycle, n)
		}
		return
	}
	if a.black[n] {
		return
	}
	a.gray[n] = true
	for i, num := 0, a.Graph.NumOuts(n); i < num; i += 1 {
		a.visit(a.Graph.Out(n, i))
	}
	delete(a.gray, n)
	if !a.black[n] {
		a.black[n] = true
		a.Order = append(a.Order, n)
	}
}

func (a *topoOrder) cycleError() error {
	if len(a.Cycle) == 0 {
		return nil
	} else {
		return fmt.Errorf("cycle containing %q", a.Cycle[0])
	}
}

// Run topographic sort
func topoSort(opt TopoSortOpt) *topoOrder {
	g := opt.Graph
	if opt.NodeOrder != nil {
		g = makeSortedGraph(opt.Graph, opt.NodeOrder)
	}
	to := &topoOrder{
		Graph: g,
		black: make(nodeSet),
		gray:  make(nodeSet),
	}
	for i, num := 0, g.NumNodes(); i < num; i += 1 {
		n := g.Node(i)
		if _, seen := to.black[n]; !seen {
			to.visit(n)
		}
	}
	return to
}

// Return nil if graph is acyclic. If graph contains a cycle, return error.
func IsDag(g Graph) error {
	return topoSort(TopoSortOpt{Graph: g}).cycleError()
}

type TopoSortOpt struct {
	Graph     Graph
	NodeOrder LessThan // Optional argument to ensure deterministic output
}

// Return topological sort of graph. Returns an error if graph contains a
// cycle.
func TopoSort(opt TopoSortOpt) ([]Node, error) {
	to := topoSort(opt)
	if err := to.cycleError(); err != nil {
		return nil, err
	} else {
		return to.Order, nil
	}
}
