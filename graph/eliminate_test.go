package graph

import "testing"

func TestTreeEliminate(t *testing.T) {
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

	gnext := Eliminate(EliminateOpt{
		Graph: g,
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

func TestGraphEliminate(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": []string{"c"},
		"b": []string{"c"},
		"c": []string{"d"},
		"d": []string{"e"},
		"e": []string{"f"},
		"f": []string{"g", "h"},
	})

	out := map[string]bool{
		"c": true,
		"d": true,
		"e": true,
		"f": true,
	}

	gnext := Eliminate(EliminateOpt{
		Graph: g,
		In: func(n Node) bool {
			return !out[n.(string)]
		},
	})

	if l := gnext.NumNodes(); l != 4 {
		t.Errorf("expected %d nodes found %d", 4, l)
	} else if l := gnext.NumOuts("a"); l != 2 {
		t.Errorf("expected %d nodes found %d", 2, l)
	} else if l := gnext.NumOuts("b"); l != 2 {
		t.Errorf("expected %d nodes found %d", 2, l)
	}
}
