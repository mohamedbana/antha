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

type TreeVisitOpt struct {
	Graph     Graph
	Root      Node
	PreOrder  Visitor
	PostOrder Visitor
}

// Apply a tree visitor.
func TreeVisit(opt TreeVisitOpt) error {
	apply := func(v Visitor, n Node) error {
		if v == nil {
			return nil
		}
		return v(n)
	}

	type frame struct {
		Node Node
		Post bool
	}

	stack := []frame{frame{Node: opt.Root}}

	for l := len(stack); l > 0; l = len(stack) {
		f := stack[l-1]
		stack = stack[:l-1]

		if f.Post {
			if err := apply(opt.PostOrder, f.Node); err != nil {
				return err
			}
			continue
		} else {
			if err := apply(opt.PreOrder, f.Node); err != nil {
				return err
			}
		}

		stack = append(stack, frame{Node: f.Node, Post: true})

		for i, inum := 0, opt.Graph.NumOuts(f.Node); i < inum; i += 1 {
			stack = append(stack, frame{Node: opt.Graph.Out(f.Node, i)})
		}
	}

	return nil
}
