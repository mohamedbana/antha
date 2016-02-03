package graph

type DisjointSet struct {
	parent map[Node]Node
	nodes  []Node
}

func NewDisjointSet() *DisjointSet {
	return &DisjointSet{
		parent: make(map[Node]Node),
	}
}

func (a *DisjointSet) NumNodes() int {
	return len(a.nodes)
}

func (a *DisjointSet) Node(i int) Node {
	return a.nodes[i]
}

func (a *DisjointSet) NumOuts(n Node) int {
	nr := a.Find(n)
	if nr == n {
		return 0
	} else {
		return 1
	}
}

func (a *DisjointSet) Out(n Node, i int) Node {
	return a.Find(n)
}

func (a *DisjointSet) Union(x, y Node) {
	xr := a.Find(x)
	yr := a.Find(y)

	a.parent[xr] = yr
}

func (a *DisjointSet) Find(n Node) Node {
	p := a.parent[n]
	if p == nil {
		a.parent[n] = n
		a.nodes = append(a.nodes, n)
		p = n
	}

	if p != n {
		a.parent[n] = a.Find(p)
	}
	return a.parent[n]
}
