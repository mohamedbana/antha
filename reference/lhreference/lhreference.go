package lhreference

import (
	// some things
	// goflow most likely
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/anthalib/liquidhandling"
	"github.com/antha-lang/antha/anthalib/mixer"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/goflow"
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
	A_vol <-chan execute.ThreadParam
	B_vol <-chan execute.ThreadParam

	// these are materials

	A    <-chan execute.ThreadParam
	B    <-chan execute.ThreadParam
	Dest <-chan execute.ThreadParam

	// this is the output

	Mixture chan<- execute.ThreadParam

	// holders for the blocks

	ParamBlocks map[execute.ThreadID]*execute.AsyncBag
	InputBlocks map[execute.ThreadID]*execute.AsyncBag
	PIBlocks    map[execute.ThreadID]*execute.AsyncBag

	// sync structure

	lock sync.Mutex
}

// constructor

func (lhr *LHReference) init() {
	lhr.ParamBlocks = make(map[execute.ThreadID]*execute.AsyncBag)
	lhr.InputBlocks = make(map[execute.ThreadID]*execute.AsyncBag)
	lhr.PIBlocks = make(map[execute.ThreadID]*execute.AsyncBag)
}

func NewLHReference() *LHReference {
	lhr := new(LHReference)
	lhr.init()
	return lhr
}

// complete function for LHReference

func (lh *LHReference) Complete(val interface{}) {
	switch val.(type) {
	case *ParamBlock:
		fmt.Println("Complete ParamBlock")
		var pib PIBlock
		v := val.(*ParamBlock)
		tp := execute.ThreadParam{v, v.ID}
		AddFeature("Params", tp, &pib, lh, &lh.PIBlocks, 2, lh.lock)
	case *InputBlock:
		fmt.Println("Complete InputBlock")
		var pib PIBlock
		v := val.(*InputBlock)
		tp := execute.ThreadParam{v, v.ID}
		AddFeature("Inputs", tp, &pib, lh, &lh.PIBlocks, 2, lh.lock)
	case *PIBlock:
		fmt.Println("Complete PIBlock")
		// we have everything we need so just do the steps
		pib := val.(*PIBlock)
		lh.Setup(pib)
		lh.Steps(pib)
	}
}

// candidate for refactoring out into execute
func AddFeature(name string, param execute.ThreadParam, mapper execute.AsyncMapper, completer execute.AsyncCompleter, blocks *map[execute.ThreadID]*execute.AsyncBag, nkeys int, lock sync.Mutex) {
	var bag *execute.AsyncBag = (*blocks)[param.ID]

	if bag == nil {
		lock.Lock()
		bag = new(execute.AsyncBag)
		bag.Init(nkeys, completer, mapper)
		(*blocks)[param.ID] = bag
		lock.Unlock()
	}

	fired := bag.AddValue(name, param)

	if fired {
		lock.Lock()
		delete(*blocks, param.ID)
		lock.Unlock()
	}
}

// ports for wiring into the network
func (lh *LHReference) OnA_vol(param execute.ThreadParam) {
	var p ParamBlock
	AddFeature("A_vol", param, &p, lh, &(lh.ParamBlocks), 2, lh.lock)
}
func (lh *LHReference) OnB_vol(param execute.ThreadParam) {
	var p ParamBlock
	AddFeature("B_vol", param, &p, lh, &(lh.ParamBlocks), 2, lh.lock)
}
func (lh *LHReference) OnA(param execute.ThreadParam) {
	var i InputBlock
	AddFeature("A", param, &i, lh, &(lh.InputBlocks), 3, lh.lock)
}
func (lh *LHReference) OnB(param execute.ThreadParam) {
	var i InputBlock
	AddFeature("B", param, &i, lh, &(lh.InputBlocks), 3, lh.lock)
}
func (lh *LHReference) OnDest(param execute.ThreadParam) {
	var i InputBlock
	AddFeature("Dest", param, &i, lh, &(lh.InputBlocks), 3, lh.lock)
}

// we need a two-level asyncbag structure

// the top level is the PIblock

type PIBlock struct {
	flow.Component
	Params *ParamBlock
	Inputs *InputBlock
	ID     execute.ThreadID
}

// the next levels down are the paramblock and input block structs

type ParamBlock struct {
	A_vol wunit.Volume
	B_vol wunit.Volume
	ID    execute.ThreadID
}

type InputBlock struct {
	A    *liquidhandling.LHComponent
	B    *liquidhandling.LHComponent
	Dest *liquidhandling.LHWell
	ID   execute.ThreadID
}

// JSON blocks are also required... not quite sure why though
// I'm sure we can serialize the paramblock OK anyway
type JSONBlock struct {
	A_vol *wunit.Volume
	B_vol *wunit.Volume
	A     *liquidhandling.LHComponent
	B     *liquidhandling.LHComponent
	Dest  *liquidhandling.LHWell
	ID    *execute.ThreadID
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
	p.A_vol = *(m["A_vol"].(execute.ThreadParam).Value.(*wunit.Volume))
	p.B_vol = *(m["B_vol"].(execute.ThreadParam).Value.(*wunit.Volume))
	p.ID = m["A_vol"].(execute.ThreadParam).ID
	return p
}

func (i *InputBlock) Map(m map[string]interface{}) interface{} {
	i.A = m["A"].(execute.ThreadParam).Value.(*liquidhandling.LHComponent)
	i.B = m["B"].(execute.ThreadParam).Value.(*liquidhandling.LHComponent)
	i.Dest = m["Dest"].(execute.ThreadParam).Value.(*liquidhandling.LHWell)
	i.ID = m["A"].(execute.ThreadParam).ID
	return i
}

func (pi *PIBlock) Map(m map[string]interface{}) interface{} {
	pi.Params = m["Params"].(execute.ThreadParam).Value.(*ParamBlock)
	pi.Inputs = m["Inputs"].(execute.ThreadParam).Value.(*InputBlock)
	pi.ID = m["Params"].(execute.ThreadParam).ID
	return pi
}

// and definitions of the setup and steps methods

func (lh *LHReference) Setup(v interface{}) {
	// Only needed where there are controls
}

func (lh *LHReference) Steps(v interface{}) {
	pib := v.(*PIBlock)
	params := pib.Params
	inputs := pib.Inputs
	// needs an overhaul
	s := mixer.Sample(inputs.A, params.A_vol)
	s2 := mixer.Sample(inputs.B, params.B_vol)
	lhr := mixer.MixInto(inputs.Dest, s, s2)
	liquidhandler := liquidhandling.Init(lhp)
	liquidhandler.MakeSolutions(&lhr)
}
