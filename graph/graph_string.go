package graph

type StringGraph struct {
	Nodes []string
	Outs  map[string][]string
}

func (a *StringGraph) NumNodes() int {
	return len(a.Nodes)
}

func (a *StringGraph) Node(i int) Node {
	return a.Nodes[i]
}

func (a *StringGraph) NumOuts(n Node) int {
	return len(a.Outs[n.(string)])
}

func (a *StringGraph) Out(n Node, i int) Node {
	return a.Outs[n.(string)][i]
}
