// Package codegen compiles generic instructions to target-specific ones.
// Target, in this case, is some combination of devices (e.g., two
// ExtendedLiquidHandlers and human plate mover).
package codegen

import (
	"errors"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/target"
)

var (
	tbd = errors.New("not yet implemented")
)

type ir struct {
	Root        ast.Node
	Graph       graph.Graph
	CommandTree graph.Graph // Tree of intrinsic (and potentially BundleExpr root)
}

// Run of a device.
type drun struct {
	Device target.Device
}

type plan struct {
	Assignment map[ast.Node]*drun      // ApplyExprs to device runs
	Output     map[*drun][]target.Inst // Output of device-specific planners
}

// Assign runs of a device to each ApplyExpr. Construct initial plan by
// by maximally coalescing ApplyExprs with the same device into the same
// device run.
func assignDevices(ir *ir, t *target.Target) (*plan, error) {
	colors := make(map[ast.Node][]target.Device)
	for i, inum := 0, ir.CommandTree.NumNodes(); i < inum; i += 1 {
		n := ir.CommandTree.Node(i).(ast.Node)
		var reqs []ast.Request
		if c, ok := n.(ast.Command); ok {
			reqs = c.Requests()
		}
		colors[n] = t.Can(reqs...)
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
		Tree: ir.CommandTree,
		Root: ir.Root,
		Colors: func(n graph.Node) (r []int) {
			for _, d := range colors[n.(ast.Node)] {
				r = append(r, d2c[d])
			}
			return
		},
		EdgeWeight: func(a, b int) int {
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
		for i, inum := 0, ir.CommandTree.NumOuts(n); i < inum; i += 1 {
			kid := ir.CommandTree.Out(n, i).(ast.Node)
			if device[kid] == device[n] {
				return run[kid]
			}
		}
		return nil
	}

	dag := graph.Schedule(graph.Reverse(ir.CommandTree))

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

type instGen struct {
	Ir        *ir
	Plan      *plan
	dg        graph.QGraph
	insts     []target.Inst
	dependsOn map[target.Inst][]target.Inst
	entry     map[graph.Node]target.Inst
	exit      map[graph.Node]target.Inst
}

func (a *instGen) NumNodes() int {
	return len(a.insts)
}

func (a *instGen) Node(i int) graph.Node {
	return a.insts[i]
}

func (a *instGen) NumOuts(n graph.Node) int {
	return len(a.dependsOn[n.(target.Inst)])
}

func (a *instGen) Out(n graph.Node, i int) graph.Node {
	return a.dependsOn[n.(target.Inst)][i]
}

// Find the components that connect from to to
func (a *instGen) findComps(from, to ast.Node) ([]*ast.UseComp, error) {
	findComp := func(u *ast.UseComp, t ast.Node) *ast.UseComp {
		for _, v := range u.From {
			if v == t {
				return u
			}
		}
		return nil
	}

	var comps []*ast.UseComp
	for i, inum := 0, a.Ir.Graph.NumOuts(from); i < inum; i += 1 {
		n := a.Ir.Graph.Out(from, i)
		if u, ok := n.(*ast.UseComp); !ok {
			continue
		} else if c := findComp(u, to); c == nil {
			continue
		} else {
			comps = append(comps, c)
		}
	}
	return comps, nil
}

// Create move of dependencies if neccessary
func (a *instGen) makeMove(run *drun, root graph.Node) (*target.MoveToInst, error) {
	var comps []*wtype.LHComponent
	for i, inum := 0, a.dg.NumOrigs(root); i < inum; i += 1 {
		n := a.dg.Orig(root, i).(ast.Node)
		for j, jnum := 0, a.Ir.CommandTree.NumOuts(n); j < jnum; j += 1 {
			out := a.Ir.CommandTree.Out(n, j).(ast.Node)
			if run == a.Plan.Assignment[out] {
				continue
			}
			// Command has input dependence on a previous run
			if cs, err := a.findComps(n, out); err != nil {
				return nil, err
			} else {
				for _, c := range cs {
					comps = append(comps, c.Value)
				}
			}
		}
	}
	if len(comps) == 0 {
		return nil, nil
	}
	return &target.MoveToInst{
		Comps: comps,
		Dev:   run.Device,
	}, nil
}

func (a *instGen) add(root graph.Node, run *drun, insts []target.Inst) error {
	exit := &target.WaitInst{}
	var entry target.Inst
	if move, err := a.makeMove(run, root); err != nil {
		return err
	} else if move != nil {
		entry = move
	} else {
		entry = &target.WaitInst{}
	}

	a.entry[root] = entry
	a.exit[root] = exit
	a.insts = append(a.insts, entry)
	a.insts = append(a.insts, exit)
	a.dependsOn[exit] = append(a.dependsOn[exit], entry)

	for idx, in := range insts {
		if idx == 0 {
			a.dependsOn[in] = append(a.dependsOn[in], entry)
		}
		if idx == len(insts)-1 {
			a.dependsOn[exit] = append(a.dependsOn[exit], in)
		}
		a.insts = append(a.insts, in)
		for _, v := range in.DependsOn() {
			a.dependsOn[in] = append(a.dependsOn[in], v)
		}
	}

	return nil
}

func (a *instGen) Run() ([]target.Inst, error) {
	a.entry = make(map[graph.Node]target.Inst)
	a.exit = make(map[graph.Node]target.Inst)
	a.dependsOn = make(map[target.Inst][]target.Inst)

	a.dg = graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: a.Ir.CommandTree,
		Colorer: func(n graph.Node) interface{} {
			return a.Plan.Assignment[n.(ast.Node)]
		},
	})

	order, err := graph.TopoSort(graph.TopoSortOpt{
		Graph: a.dg,
	})
	if err != nil {
		return nil, err
	}

	// Add generated instructions along with any required moves
	for _, n := range order {
		if a.dg.NumOrigs(n) == 0 {
			return nil, fmt.Errorf("no instructions for node %q", n)
		}
		someNode := a.dg.Orig(n, 0).(ast.Node)
		run := a.Plan.Assignment[someNode]
		insts := a.Plan.Output[run]
		if err := a.add(n, run, insts); err != nil {
			return nil, err
		}
	}

	// Add tree edges
	for i, inum := 0, a.dg.NumNodes(); i < inum; i += 1 {
		n := a.dg.Node(i)
		nentry := a.entry[n]
		for j, jnum := 0, a.dg.NumOuts(n); j < jnum; j += 1 {
			dst := a.dg.Out(n, j)
			dexit := a.exit[dst]
			a.dependsOn[nentry] = append(a.dependsOn[nentry], dexit)
		}
	}

	// Remove synthetic nodes and redundant edges
	sg, err := graph.TransitiveReduction(graph.Eliminate(graph.EliminateOpt{
		Graph: a,
		In: func(n graph.Node) bool {
			if _, isWait := n.(*target.WaitInst); isWait {
				return false
			}
			return true
		},
	}))
	if err != nil {
		return nil, err
	}

	// Cleanup dependencies
	order, err = graph.TopoSort(graph.TopoSortOpt{
		Graph: sg,
	})
	if err != nil {
		return nil, err
	}

	var insts []target.Inst
	for _, n := range order {
		var depends []target.Inst
		for j, jnum := 0, sg.NumOuts(n); j < jnum; j += 1 {
			depends = append(depends, sg.Out(n, j).(target.Inst))
		}

		in := n.(target.Inst)
		in.SetDependsOn(depends)

		insts = append(insts, in)
	}

	return insts, nil
}

// Lower plan to instructions: add manual moves between devices and update
// original IR nodes for incremental compiles
func genInsts(ir *ir, p *plan) ([]target.Inst, error) {
	gen := &instGen{
		Ir:   ir,
		Plan: p,
	}
	// XXX set outputs
	return gen.Run()
}

// Run plan through device-specific planners. Adjust assignment based on
// planner capabilities and return output.
func tryPlan(ir *ir, p *plan) (*plan, error) {
	dg := graph.MakeQuotient(graph.MakeQuotientOpt{
		Graph: ir.CommandTree,
		Colorer: func(n graph.Node) interface{} {
			return p.Assignment[n.(ast.Node)]
		},
	})
	if err := graph.IsDag(dg); err != nil {
		return nil, fmt.Errorf("invalid assignment: %s", err)
	}

	// TODO: When initial assignment is not feasible for a device (e.g.,
	// capacity constraints), split up run until feasible or give up.

	// TODO: When splitting a mix sequence, adjust LHInstructions to place
	// output samples on the same plate

	cmds := make(map[*drun][]ast.Command)
	for n, d := range p.Assignment {
		if c, ok := n.(ast.Command); !ok {
			continue
		} else {
			cmds[d] = append(cmds[d], c)
		}
	}

	output := make(map[*drun][]target.Inst)
	for d, cs := range cmds {
		if insts, err := d.Device.Compile(cs); err != nil {
			return nil, err
		} else {
			output[d] = insts
		}
	}

	return &plan{
		Assignment: p.Assignment,
		Output:     output,
	}, nil
}

// Compile an expression program into a sequence of instructions for a target
// configuration. This supports incremental compilation, so roots may refer to
// nodes that have already been compiled, in which case, the result may refer
// to previously generated instructions.
func Compile(t *target.Target, roots []ast.Node) ([]target.Inst, error) {
	if len(roots) == 0 {
		return nil, nil
	}

	if root, err := makeRoot(roots); err != nil {
		return nil, fmt.Errorf("invalid program: %s", err)
	} else if ir, err := build(root); err != nil {
		return nil, fmt.Errorf("invalid program: %s", err)
	} else if plan, err := assignDevices(ir, t); err != nil {
		return nil, fmt.Errorf("error assigning devices: %s", err)
	} else if plan, err := tryPlan(ir, plan); err != nil {
		return nil, fmt.Errorf("error planning: %s", err)
	} else if insts, err := genInsts(ir, plan); err != nil {
		return nil, fmt.Errorf("error generating instructions: %s", err)
	} else {
		return insts, nil
	}
}
