package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

func summarize(inst target.Inst) (string, error) {
	switch inst := inst.(type) {
	case *target.Manual:
		return fmt.Sprintf("Manual: %s", inst.Details), nil
	case *target.Mix:
		return fmt.Sprintf("Run %p", inst.Dev), nil
	case *target.Wait:
		return "", nil
	default:
		return "", fmt.Errorf("unknown inst %T", inst)

	}
}

func printTimeline(out io.Writer, result *execute.Result) error {
	g := &target.Graph{
		Insts: result.Insts,
	}

	dag := graph.Schedule(graph.Reverse(g))
	for round := 1; len(dag.Roots) != 0; round += 1 {
		if _, err := fmt.Fprintf(out, "== Round %2d:\n", round); err != nil {
			return err
		}
		var next []graph.Node
		for _, n := range dag.Roots {
			inst := n.(target.Inst)
			if s, err := summarize(inst); err != nil {
				return err
			} else if _, err := fmt.Fprintf(out, "    * %s\n", s); err != nil {
				return err
			}
			next = append(next, dag.Visit(n)...)
		}

		dag.Roots = next
	}

	if _, err := fmt.Fprint(out, "== Workflow Outputs:\n"); err != nil {
		return err
	}

	for k, v := range result.Workflow.Outputs {
		bs, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "    - %s: %s\n", k, string(bs)); err != nil {
			return err
		}
	}

	return nil
}
