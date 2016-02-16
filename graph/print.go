package graph

import (
	"fmt"
	"strings"
)

type Labeler func(interface{}) string

// Only support dot output
type PrintOpt struct {
	Graph        Graph
	NodeLabelers []Labeler
}

func typeLabeler(n interface{}) string {
	return fmt.Sprintf("%T", n)
}

func defaultLabeler(n interface{}) string {
	return fmt.Sprintf("%+v", n)
}

func Print(opt PrintOpt) string {
	var lines []string

	var labelers []Labeler
	labelers = append(labelers, typeLabeler)
	if len(opt.NodeLabelers) == 0 {
		labelers = append(labelers, defaultLabeler)
	}
	labelers = append(labelers, opt.NodeLabelers...)

	nodes := make(map[Node]string)
	lines = append(lines, "digraph {")
	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i += 1 {
		n := opt.Graph.Node(i)
		name := fmt.Sprintf("v%d", i)
		nodes[n] = name

		var labels []string
		for _, ler := range labelers {
			labels = append(labels, ler(n))
		}

		lines = append(lines, fmt.Sprintf("%s [label=%q];", name, strings.Join(labels, "\n")))
	}

	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i += 1 {
		src := opt.Graph.Node(i)
		sname := nodes[src]
		for j, jnum := 0, opt.Graph.NumOuts(src); j < jnum; j += 1 {
			dst := opt.Graph.Out(src, j)
			dname := nodes[dst]
			lines = append(lines, fmt.Sprintf("%s->%s;", sname, dname))
		}
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}
