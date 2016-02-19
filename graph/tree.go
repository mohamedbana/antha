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

type TreeVisitor func(n, parent Node, err error) error

type VisitTreeOpt struct {
	Tree      Graph
	Root      Node
	PreOrder  TreeVisitor // if err != TraversalDone, propagate error
	PostOrder TreeVisitor // if err != TraversalDone, propagate error
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
				if err == TraversalDone {
					break
				}
			}
			continue
		} else {
			if err := apply(opt.PreOrder, f.Node, f.Parent, lastError); err != nil {
				lastError = err
				if err == TraversalDone {
					break
				}
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
