package graph

import (
	"fmt"
)

// Return nil if graph is a rooted tree. If not, return error containing
// a counterexample.
func IsTree(g Graph, root Node) error {
	_, err := Visit(VisitOpt{
		Root:  root,
		Graph: g,
		Seen: func(n Node) error {
			return fmt.Errorf("not tree: at least two paths to %s", n)
		},
	})
	return err
}

type FilterTreeOpt struct {
	Tree Graph
	Root Node
	In   func(Node) bool // Should node be included
}

// Return tree (generally forest) that results from removing some nodes of a
// tree while keeping the same reachability relation (e.g., if a is an ancestor
// of b in G, a' will be an ancestor of b' in G').
func FilterTree(opt FilterTreeOpt) Graph {
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}
	inParent := make(map[Node]Node)
	VisitTree(VisitTreeOpt{
		Tree: opt.Tree,
		Root: opt.Root,
		PreOrder: func(n, parent Node, err error) error {
			var p Node
			if !opt.In(n) {
				// Not in
				p = inParent[parent]
			} else {
				p = n
				ret.Nodes = append(ret.Nodes, n)
				if pp := inParent[parent]; pp != nil {
					ret.Outs[pp] = append(ret.Outs[pp], n)
				}
			}
			inParent[n] = p
			return nil
		},
	})
	return ret
}

type TreeVisitor func(n, parent Node, err error) error

type VisitTreeOpt struct {
	Tree      Graph
	Root      Node
	PreOrder  TreeVisitor
	PostOrder TreeVisitor
}

// Apply a tree visitor.
func VisitTree(opt VisitTreeOpt) error {
	apply := func(v TreeVisitor, n, parent Node, err error) error {
		if v == nil {
			return nil
		}
		return v(n, parent, err)
	}

	type frame struct {
		Parent, Node Node
		Post         bool
	}

	stack := []frame{frame{Node: opt.Root}}

	var lastError error
	for l := len(stack); l > 0; l = len(stack) {
		f := stack[l-1]
		stack = stack[:l-1]

		if f.Post {
			if err := apply(opt.PostOrder, f.Node, f.Parent, lastError); err != nil {
				lastError = err
			}
			continue
		} else {
			if err := apply(opt.PreOrder, f.Node, f.Parent, lastError); err != nil {
				lastError = err
				continue
			}
		}

		stack = append(stack, frame{Node: f.Node, Parent: f.Parent, Post: true})

		for i, inum := 0, opt.Tree.NumOuts(f.Node); i < inum; i += 1 {
			stack = append(stack, frame{Node: opt.Tree.Out(f.Node, i), Parent: f.Node})
		}
	}

	return lastError
}
