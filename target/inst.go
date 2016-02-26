package target

import (
	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

type Inst interface {
	Device() Device
	DependsOn() []Inst
	SetDependsOn([]Inst)
}

type Graph struct {
	Insts []Inst
}

func (a *Graph) NumNodes() int {
	return len(a.Insts)
}

func (a *Graph) Node(i int) graph.Node {
	return a.Insts[i]
}

func (a *Graph) NumOuts(n graph.Node) int {
	return len(n.(Inst).DependsOn())
}

func (a *Graph) Out(n graph.Node, i int) graph.Node {
	return n.(Inst).DependsOn()[i]
}

type Mix struct {
	Dev        Device
	Depends    []Inst
	Request    *lh.LHRequest
	Properties liquidhandling.LHProperties
	Files      []byte
}

func (a *Mix) Device() Device {
	return a.Dev
}

func (a *Mix) DependsOn() []Inst {
	return a.Depends
}

func (a *Mix) SetDependsOn(x []Inst) {
	a.Depends = x
}

type Manual struct {
	Dev     Device
	Label   string
	Details string
	Depends []Inst
}

func (a *Manual) DependsOn() []Inst {
	return a.Depends
}

func (a *Manual) Device() Device {
	return a.Dev
}

func (a *Manual) SetDependsOn(x []Inst) {
	a.Depends = x
}

// Virtual instruction to hang dependencies on
type Wait struct {
	Depends []Inst
}

func (a *Wait) Device() Device {
	return nil
}

func (a *Wait) DependsOn() []Inst {
	return a.Depends
}

func (a *Wait) SetDependsOn(x []Inst) {
	a.Depends = x
}
