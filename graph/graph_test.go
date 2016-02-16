package graph

import (
	"fmt"
	"testing"
)

func sameElements(a, b []string) error {
	makeM := func(xs []string) map[string]int {
		m := make(map[string]int)
		for _, x := range xs {
			m[x] += 1
		}
		return m
	}

	as := makeM(a)
	bs := makeM(b)

	if la, lb := len(as), len(bs); la != lb {
		return fmt.Errorf("expecting %d values found %d", la, lb)
	}

	for k, v := range as {
		if bv := bs[k]; bv != v {
			return fmt.Errorf("expecting %d of %q found %d", v, k, bv)
		}
	}
	return nil
}

func toString(ns []Node) (r []string) {
	for _, n := range ns {
		r = append(r, n.(string))
	}
	return
}

type TestGraph struct {
	Nodes []string
	Outs  map[string][]string
}

func (a *TestGraph) NumNodes() int {
	return len(a.Nodes)
}

func (a *TestGraph) Node(i int) Node {
	return a.Nodes[i]
}

func (a *TestGraph) NumOuts(n Node) int {
	return len(a.Outs[n.(string)])
}

func (a *TestGraph) Out(n Node, i int) Node {
	return a.Outs[n.(string)][i]
}

func MakeTestGraph(m map[string][]string) *TestGraph {
	g := &TestGraph{
		Outs: m,
	}
	ns := make(map[string]bool)
	for k, outs := range m {
		ns[k] = true
		for _, neigh := range outs {
			ns[neigh] = true
		}
	}
	for k := range ns {
		g.Nodes = append(g.Nodes, k)
	}
	return g
}

func TestReverse(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": []string{"b", "c"},
		"d": []string{"c"},
	})
	rg := Reverse(g)
	if l := rg.NumNodes(); l != 4 {
		t.Errorf("expected %d nodes found %d", 4, l)
	} else if l := rg.NumOuts("a"); l != 0 {
		t.Errorf("expected %d neighbors found %d", 0, l)
	} else if l := rg.NumOuts("b"); l != 1 {
		t.Errorf("expected %d neighbors found %d", 1, l)
	} else if l := rg.NumOuts("c"); l != 2 {
		t.Errorf("expected %d neighbors found %d", 2, l)
	} else if l := rg.NumOuts("d"); l != 0 {
		t.Errorf("expected %d neighbors found %d", 0, l)
	}
}

func TestNodeSet(t *testing.T) {
	ns := make(nodeSet)
	ns["a"] = true
	ns["b"] = true
	ns["c"] = true

	var r []Node
	for n := range ns.Range() {
		r = append(r, n)
	}
	if err := sameElements(toString(r), []string{"a", "b", "c"}); err != nil {
		t.Error(err)
	}
}
