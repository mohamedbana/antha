// Package codegen compiles generic instructions to target-specific ones.
// Target, in this case, is some combination of devices (e.g., two
// ExtendedLiquidHandlers and human plate mover).
package codegen

import (
	"errors"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

var (
	tbd = errors.New("not yet implemented")
)

type ir struct {
	Root      ast.Node
	Graph     graph.Graph
	DataTree  graph.Graph // Data dependencies
	ApplyTree graph.Graph // Tree of ApplyExprs (and potentially BundleExpr root)
}

func check(root ast.Node) (*ir, error) {
	if r, err := normalize(root); err != nil {
		return nil, err
	} else {
		root = r
	}

	// Check that data dependencies are single-consumer (because they model
	// destructive update of values)
	dg := ast.ToGraph(ast.ToGraphOpt{
		Root:      root,
		WhichDeps: ast.DataDeps,
	})
	if err := graph.IsTree(dg, root); err != nil {
		return nil, err
	}

	g := ast.ToGraph(ast.ToGraphOpt{Root: root})
	if err := graph.IsDag(coalesceByGen(g)); err != nil {
		return nil, err
	}

	return &ir{Root: root,
		Graph:    g,
		DataTree: dg,
		ApplyTree: graph.FilterTree(graph.FilterTreeOpt{
			Tree: dg,
			Root: root,
			In: func(n graph.Node) bool {
				_, ok := n.(*ast.ApplyExpr)
				return ok || n == root
			},
		}),
	}, nil
}

// Return graph with apply expressions in the same generation as a single node
func coalesceByGen(g *ast.Graph) graph.QGraph {
	return graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: g,
		Colorer: func(n graph.Node) interface{} {
			return n.(*ast.ApplyExpr).Gen
		},
		HasColor: func(n graph.Node) bool {
			switch n := n.(ast.Node).(type) {
			case *ast.ApplyExpr:
				return n.Gen != nil
			}
			return false
		},
	})
}

// TODO(ddn): if/when we support pausing the current run of machine a, taking
// something out, running it on machine b, and then returning it to machine a
// and continuing the run, drun will nolonger represent an actual run on a
// machine but rather a step in a larger run...

// Run of a device.
type drun struct {
	Device target.Device
}

type plan struct {
	Assignment map[ast.Node]*drun    // ApplyExprs to device runs
	Output     map[*drun]interface{} // Output of device-specific planners
}

// Assign runs of a device to each ApplyExpr. Construct initial plan by
// by maximally coalescing ApplyExprs with the same device into the same
// device run.
func assignDevices(ir *ir, t *target.Target) (*plan, error) {
	colors := make(map[ast.Node][]target.Device)
	for i, inum := 0, ir.ApplyTree.NumNodes(); i < inum; i += 1 {
		n := ir.ApplyTree.Node(i).(ast.Node)
		colors[n] = t.Can() // XXX
	}

	var devices []target.Device
	d2c := make(map[target.Device]int)
	for _, ds := range colors {
		for _, d := range ds {
			if _, seen := d2c[d]; !seen {
				d2c[d] = len(devices)
				devices = append(devices, d)
			}
		}
	}

	r, err := graph.PartitionTree(graph.PartitionTreeOpt{
		Tree: ir.ApplyTree,
		Root: ir.Root,
		Colors: func(n graph.Node) (r []int) {
			for _, d := range colors[n.(ast.Node)] {
				r = append(r, d2c[d])
			}
			return
		},
		Weights: func(a, b int) int {
			return devices[a].MoveCost(devices[b])
		},
	})
	if err != nil {
		return nil, err
	}

	ret := make(map[ast.Node]target.Device)
	for n, idx := range r.Parts {
		ret[n.(ast.Node)] = devices[idx]
	}
	return coalesceDevices(ir, ret), nil
}

// Coalesce adjacent devices into the same run of a device
func coalesceDevices(ir *ir, device map[ast.Node]target.Device) *plan {
	run := make(map[ast.Node]*drun)

	kidRun := func(n ast.Node) *drun {
		for i, inum := 0, ir.ApplyTree.NumOuts(n); i < inum; i += 1 {
			kid := ir.ApplyTree.Out(n, i).(ast.Node)
			if device[kid] == device[n] {
				return run[kid]
			}
		}
		return nil
	}

	dag := graph.Schedule(graph.Reverse(ir.ApplyTree))

	for len(dag.Roots) > 0 {
		var next []graph.Node
		newRuns := make(map[target.Device]*drun)
		for _, n := range dag.Roots {
			n := n.(ast.Node)

			myRun := kidRun(n)
			if myRun == nil {
				d := device[n]
				if r, seen := newRuns[d]; seen {
					myRun = r
				} else {
					myRun = &drun{Device: d}
					newRuns[d] = myRun
				}
			}
			run[n] = myRun
			next = append(next, dag.Visit(n)...)
		}

		dag.Roots = next
	}

	return &plan{Assignment: run}
}

// Run plan past device-specific planners. When initial assignment is not be
// feasible for a device (e.g., capacity constraints), split up run until
// feasible or give up. Return updated assignment and output from running
// device-specific planners on the assignment.
func tryPlan(ir *ir, plan *plan) (*plan, error) {
	return nil, tbd
}

// Compile an expression program into a sequence of instructions for a target
// configuration.
func Compile(t *target.Target, root ast.Node) ([]Inst, error) {
	if root == nil {
		return nil, nil
	}

	if ir, err := check(root); err != nil {
		return nil, err
	} else if plan, err := assignDevices(ir, t); err != nil {
		return nil, err
	} else if err := graph.IsDag(graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: ir.ApplyTree,
		Colorer: func(n graph.Node) interface{} {
			return plan.Assignment[n.(ast.Node)]
		},
	})); err != nil {
		return nil, err
	} else if plan, err := tryPlan(ir, plan); err != nil {
		return nil, err
	} else {
		fmt.Println(plan)
		return nil, tbd
	}
	// 0. Add ManualMover, "ThermoMix", ManualMixer
	// 1. Update ApplyExpr to have target.Constraints
	// 2. Decorate ApplyExprs with possible Devices
	// 3. Assign each ApplyExpr to a target.Device
	// 4. Coarsen wrt to dependencies (bottom-up) "maximal"
	// 5. Check size constraints (top-down): plan: splitting Devices
	// 6. Take plan and add manual moves
	// 7. Output instructions
}
