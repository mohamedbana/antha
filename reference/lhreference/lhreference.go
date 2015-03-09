package lhreference

import (
	// some things
	// goflow most likely
	"bytes"
	"encoding/json"
	"github.com/Synthace/Goflow"
	"github.com/antha-lang/antha/anthalib/mixer"
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/execute"
	"io"
	"log"
	"sync"
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
	Dest <-chan wtype.LiquidContainer

	// this is the output

	Mixture chan<- wtype.Solution

	// holders for the blocks

	ParamBlocks map[execute.ThreadID]*execute.AsyncBag
	InputBlocks map[execute.ThreadID]*execute.AsyncBag
	PIBlocks    map[execute.ThreadID]*execute.AsyncBag

	// sync structure

	lock sync.Mutex
}

// complete function for LHReference

func (lh *LHReference) Complete(val interface{}) {
	switch val.(type) {
	case ParamBlock:
		v := val.(ParamBlock)
		tp := execute.ThreadParam{v, v.ID}
		lh.lock.Lock()
		AddFeature("Params", tp, &v, lh, &lh.PIBlocks, 2)
		lh.lock.Unlock()
	case InputBlock:
		v := val.(InputBlock)
		tp := execute.ThreadParam{v, v.ID}
		lh.lock.Lock()
		AddFeature("Params", tp, &v, lh, &lh.PIBlocks, 2)
		lh.lock.Unlock()
	case PIBlock:
		// we have everything we need so just do the steps
		pib := val.(PIBlock)
		lh.Setup(pib)
		lh.Steps(pib)
	}
}

// candidate for refactoring out into execute
func AddFeature(name string, param execute.ThreadParam, mapper execute.AsyncMapper, completer execute.AsyncCompleter, blocks *map[execute.ThreadID]*execute.AsyncBag, nkeys int) {
	var bag *execute.AsyncBag = (*blocks)[param.ID]

	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(nkeys, completer, mapper)
	}

	fired := bag.AddValue(name, param)

	if fired {
		delete(*blocks, param.ID)
	}
}

// ports for wiring into the network
func (lh *LHReference) OnA_vol(param execute.ThreadParam) {
	var p *ParamBlock
	lh.lock.Lock()
	AddFeature("A_vol", param, p, lh, &(lh.ParamBlocks), 2)
	lh.lock.Unlock()
}
func (lh *LHReference) OnB_vol(param execute.ThreadParam) {
	var p *ParamBlock
	lh.lock.Lock()
	AddFeature("B_vol", param, p, lh, &(lh.ParamBlocks), 2)
	lh.lock.Unlock()
}
func (lh *LHReference) OnA(param execute.ThreadParam) {
	var i *InputBlock
	lh.lock.Lock()
	AddFeature("A", param, i, lh, &(lh.InputBlocks), 3)
	lh.lock.Unlock()
}
func (lh *LHReference) OnB(param execute.ThreadParam) {
	var i *InputBlock
	lh.lock.Lock()
	AddFeature("B", param, i, lh, &(lh.InputBlocks), 3)
	lh.lock.Unlock()
}
func (lh *LHReference) OnDest(param execute.ThreadParam) {
	var i *InputBlock
	lh.lock.Lock()
	AddFeature("Dest", param, i, lh, &(lh.InputBlocks), 3)
	lh.lock.Unlock()
}

// we need a two-level asyncbag structure

// the top level is the PIblock

type PIBlock struct {
	flow.Component
	Params ParamBlock
	Inputs InputBlock
	ID     execute.ThreadID
}

// the next levels down are the paramblock and input block structs

type ParamBlock struct {
	A_vol wunit.Volume
	B_vol wunit.Volume
	ID    execute.ThreadID
}

type InputBlock struct {
	A    wtype.Liquid
	B    wtype.Liquid
	Dest wtype.LiquidContainer
	ID   execute.ThreadID
}

// JSON blocks are also required... not quite sure why though
// I'm sure we can serialize the paramblock OK anyway

type JSONPBlock struct {
	A_vol *wunit.Volume
	B_vol *wunit.Volume
	ID    *execute.ThreadID
}

type JSONIBlock struct {
	A    *wtype.Liquid
	B    *wtype.Liquid
	Dest *wtype.LiquidContainer
	ID   *execute.ThreadID
}

func (p *ParamBlock) ToJSON() (b bytes.Buffer) {
	enc := json.NewEncoder(&b)
	if err := enc.Encode(p); err != nil {
		log.Fatalln(err)
	}
	return
}

func (i *InputBlock) ToJSON(b bytes.Buffer) {
	enc := json.NewEncoder(&b)
	if err := enc.Encode(i); err != nil {
		log.Fatalln(err)
	}
	return
}

func ParamsFromJSON(r io.Reader) (p *ParamBlock) {
	p = new(ParamBlock)
	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		log.Fatalln(err)
	}
	return
}

func InputsFromJSON(r io.Reader) (i *InputBlock) {
	i = new(InputBlock)
	dec := json.NewDecoder(r)
	if err := dec.Decode(i); err != nil {
		log.Fatalln(err)
	}
	return
}

// each block type needs its mapper

func (p *ParamBlock) Map(m map[string]interface{}) interface{} {
	p.A_vol = m["A_vol"].(execute.ThreadParam).Value.(wunit.Volume)
	p.B_vol = m["B_vol"].(execute.ThreadParam).Value.(wunit.Volume)
	p.ID = m["A_vol"].(execute.ThreadParam).ID
	return p
}

func (i *InputBlock) Map(m map[string]interface{}) interface{} {
	i.A = m["A"].(execute.ThreadParam).Value.(wtype.Liquid)
	i.B = m["B"].(execute.ThreadParam).Value.(wtype.Liquid)
	i.Dest = m["Dest"].(execute.ThreadParam).Value.(wtype.LiquidContainer)
	i.ID = m["A"].(execute.ThreadParam).ID
	return i
}

func (pi *PIBlock) Map(m map[string]interface{}) interface{} {
	pi.Params = m["Params"].(execute.ThreadParam).Value.(ParamBlock)
	pi.Inputs = m["Inputs"].(execute.ThreadParam).Value.(InputBlock)
	pi.ID = m["ParamBlock"].(execute.ThreadParam).ID
	return pi
}

// and definitions of the setup and steps methods

func (lh *LHReference) Setup(v interface{}) {
	// Only needed where there are controls
}

func (lh *LHReference) Steps(v interface{}) {
	pib := v.(PIBlock)
	params := pib.Params
	inputs := pib.Inputs
	// I'm not so keen on this mechanism at the moment
	// it probably needs redoing to make it easier to auto-generate
	s := mixer.Sample(inputs.A, params.A_vol)
	s2 := mixer.Sample(inputs.B, params.B_vol)
	mixer.MixInto(inputs.Dest, s, s2)
}
