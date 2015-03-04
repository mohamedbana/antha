package lhreference

import (
	// some things
	// goflow most likely
	"github.com/Synthace/Goflow"
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
)

// struct defining the antha element as a flow component

type LHReference struct {
	flow.Component

	// the element is the receiver plugged into the network
	// it holds channels for receipt of data

	// these are data items
	A_vol <-chan wunit.Volume
	B_vol <-chan wunit.Volume

	// these are materials

	A    <-chan wtype.Liquid
	B    <-chan wtype.Liquid
	Dest <-chan wtype.Labware

	// this is the output

	Mixture chan<- wtype.Solution

	// holders for the blocks

	ParamBlocks map[execute.ThreadID]*execute.AsyncBag
	InputBlocks map[execute.ThreadID]*execute.AsyncBag
	PIBlocks    map[execute.ThreadID]*execute.AsyncBag

	// sync structure

	lock sync.Lock
}

// complete function for LHReference

func (lh *LHReference) Complete(val interface{}) {
	switch val.(type) {
	case ParamBlock:
		v := val.(ParamBlock)

	case InputBlock:

	case PIBlock:

	}

}

// these generic functions should be refactored out of this class into
// the anthaelement package
func (lh *LHReference) AddParameter(name string, param execute.ThreadParam, mapper execute.AsyncMapper, completer execute.AsyncCompleter) {
	lh.lock.Lock()
	var bag *execute.AsyncBag = lh.ParamBlocks[param.ID]

	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, completer, mapper)
	}

	lh.lock.Unlock()

	fired := bag.AddValue(name, param.Value())

	if fired {
		lh.lock.Lock()
		delete(lh.ParamBlocks, param.ID)
		lh.lock.Unlock()
	}
}

func (lh *LHReference) AddInput(name string, param execute.ThreadParam, mapper execute.AsyncMapper, completer execute.AsyncCompleter) {
	lh.lock.Lock()
	var bag *execute.AsyncBag = lh.InputBlocks[param.ID]

	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.init(3, completer, mapper)
	}

	lh.lock.Unlock()

	fired := bag.AddValue(name, param.Value())

	if fired {
		lh.lock.Lock()
		delete(lh.InputBlocks, param.ID)
		lh.lock.Unlock()
	}
}

// ports for wiring into the network
func (lh *LHReference) OnA_vol(param execute.ThreadParam) {
	lh.AddParameter("A_vol", param, lh, lh)
}
func (lh *LHReference) OnB_vol(param execute.ThreadParam) {
	lh.AddParameter("B_vol", param, lh, lh)
}
func (lh *LHReference) OnA(param execute.ThreadParam) {
	lh.AddInput("A", param, lh, lh)
}
func (lh *LHReference) OnB(param execute.ThreadParam) {
	lh.AddInput("B", param, lh, lh)
}
func (lh *LHReference) OnDest(param execute.ThreadParam) {
	lh.AddInput("Dest", param, lh, lh)
}

// we need a two-level asyncbag structure

// the top level is the PIblock

type PIBlock struct {
	flow.Component
	Params Paramblock
	Inputs Inputblock
	ID     *execute.ThreadID
}

// the next levels down are the paramblock and input block structs

type ParamBlock struct {
	A_vol wunit.Volume
	B_vol wunit.Volume
	ID    *execute.ThreadID
}

type InputBlock struct {
	A    wtype.Liquid
	B    wtype.Liquid
	Dest wtype.Labware
	ID   *execute.ThreadID
}
