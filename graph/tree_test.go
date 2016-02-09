package graph

import (
	"testing"
)

func TestIsTree(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
	})
	if err := IsTree(g, "root"); err != nil {
		t.Fatal(err)
	}
}

func TestIsNotTree(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
		"a":    []string{"c", "d"},
		"b":    []string{"d"},
	})
	if err := IsTree(g, "root"); err == nil {
		t.Fatal("found tree but did not expect one")
	}
}

func TestIsNotTreeWhenHasCycle(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
		"a":    []string{"c", "root"},
	})
	if err := IsTree(g, "root"); err == nil {
		t.Fatal("found tree but did not expect one")
	}
}

func TestFilterTree(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
		"a":    []string{"c", "d"},
		"b":    []string{"e"},
		"e":    []string{"f", "g"},
	})

	in := map[string]bool{
		"root": true,
		"b":    true,
		"f":    true,
		"g":    true,
	}

	gnext := FilterTree(FilterTreeOpt{
		Tree: g,
		Root: "root",
		In: func(n Node) bool {
			return in[n.(string)]
		},
	})

	if l := gnext.NumNodes(); l != 4 {
		t.Errorf("expected %d nodes found %d", 4, l)
	} else if l := gnext.NumOuts("root"); l != 1 {
		t.Errorf("expected %d nodes found %d", 1, l)
	} else if n := gnext.Out("root", 0).(string); n != "b" {
		t.Errorf("expected %q found %q", "b", n)
	} else if l := gnext.NumOuts("b"); l != 2 {
		t.Errorf("expected %d nodes found %d", 2, l)
	} else if n := gnext.Out("b", 0).(string); n != "f" && n != "g" {
		t.Errorf("expected %q or %q found %q", "f", "g", n)
	} else if n := gnext.Out("b", 1).(string); n != "f" && n != "g" {
		t.Errorf("expected %q or %q found %q", "f", "g", n)
	}
}

func TestTreeVisit(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
		"a":    []string{"c", "d"},
		"b":    []string{"e"},
	})
	v := map[string]int{
		"c": 4,
		"d": 5,
		"e": 1,
	}
	if err := VisitTree(VisitTreeOpt{
		Tree: g,
		Root: "root",
		PostOrder: func(n, parent Node, err error) error {
			if g.NumOuts(n) == 0 {
				return nil
			}
			sum := 0
			for i, inum := 0, g.NumOuts(n); i < inum; i += 1 {
				k := g.Out(n, i).(string)
				sum += v[k]
			}
			if sum > 0 {
				v[n.(string)] = sum
			}
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	} else if vv := v["c"]; vv != 4 {
		t.Errorf("leaf value changed: expected %d found %d", 4, vv)
	} else if vv := v["d"]; vv != 5 {
		t.Errorf("leaf value changed: expected %d found %d", 5, vv)
	} else if vv := v["e"]; vv != 1 {
		t.Errorf("leaf value changed: expected %d found %d", 1, vv)
	} else if vv := v["a"]; vv != 9 {
		t.Errorf("expected %d found %d", 9, vv)
	} else if vv := v["b"]; vv != 1 {
		t.Errorf("expected %d found %d", 1, vv)
	} else if vv := v["root"]; vv != 10 {
		t.Errorf("expected %d found %d", 10, vv)
	}
}
