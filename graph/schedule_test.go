package graph

import (
	"testing"
)

func TestSchedule(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": []string{"b", "c"},
		"d": []string{"c"},
		"c": []string{"e"},
	})
	expected := [][]string{
		[]string{"a", "d"},
		[]string{"b", "c"},
		[]string{"e"},
	}
	dag := Schedule(g)
	for _, es := range expected {
		if len(dag.Roots) == 0 {
			t.Fatal("expected more nodes")
		}
		if err := sameElements(toString(dag.Roots), es); err != nil {
			t.Error(err)
		}
		var next []Node
		for _, n := range dag.Roots {
			next = append(next, dag.Visit(n)...)
		}
		dag.Roots = next
	}
}
