package lhreference

import (
	"bytes"
	"encoding/json"
	"github.com/antha-lang/antha/anthalib/execution"
	"github.com/antha-lang/antha/anthalib/liquidhandling"
	"github.com/antha-lang/antha/anthalib/mixer"
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/goflow"
	"io"
	"log"
	"sync"
)

var params = [...]string{
	"A_vol",
	"B_vol"}

// channel interfaces
// with threadID grouped types
type LHExample struct {
	flow.Component                            // component "superclass" embedded
	A_vol          <-chan execute.ThreadParam // volume of "A" to put in
	B_vol          <-chan execute.ThreadParam // volume of "B" to put in
	A              <-chan execute.ThreadParam // component A
	B              <-chan execute.ThreadParam //component B
	Dest           <-chan execute.ThreadParam // Microplate to mix into
	Mixture        chan<- execute.ThreadParam // output solution
	lock           sync.Mutex
	params         map[execute.ThreadID]*execute.AsyncBag
	inputs         map[execute.ThreadID]*execute.AsyncBag
	outputs        map[execute.ThreadID]*execute.AsyncBag
}

// make LHExample implement AsyncBag

type BlockSet struct {
	Params ParamBlock
	Inputs InputBlock
}

func (e *LHExample) Map(m map[string]interface{}) interface{} {
	var b BlockSet
	b.Params = m["params"].(ParamBlock)
	b.Inputs = m["inputs"].(InputBlock)
	return b
}

// single execution thread variables
// with concrete types
type ParamBlock struct {
	A_vol wunit.Volume
	B_vol wunit.Volume
	ID    execute.ThreadID
}

type JSONParamBlock struct {
	A_vol *wunit.Volume
	B_vol *wunit.Volume
	ID    *execute.ThreadID
}

// support function for wire format
func (p *ParamBlock) ToJSON() (b bytes.Buffer) {
	enc := json.NewEncoder(&b)
	if err := enc.Encode(p); err != nil {
		log.Fatalln(err) // currently fatal error
	}
	return
}

// helper generator function
func ParamsFromJSON(r io.Reader) (p *ParamBlock) {
	p = new(ParamBlock)
	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		log.Fatalln(err)
	}
	return
}
func InputsFromJSON(r io.Reader) (p *InputBlock) {
	p = new(InputBlock)
	dec := json.NewDecoder(r)
	if err := dec.Decode(p); err != nil {
		log.Fatalln(err)
	}
	return
}

// could handle mapping in the threadID better...
func (p *ParamBlock) Map(m map[string]interface{}) interface{} {
	p.A_vol = m["A_vol"].(execute.ThreadParam).Value.(wunit.Volume)
	p.B_vol = m["B_vol"].(execute.ThreadParam).Value.(wunit.Volume)
	m["ID"] = m["A_vol"].(execute.ThreadParam).Value.(execute.ThreadID)
	p.ID = m["ID"]
	return p
}

func (p *ParamBlock) Complete(interface{}) {

}

func (p *ParamBlock) OnA_vol(param execute.ThreadParam) {
	p.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, p, p)
		p.params[param.ID] = bag
	}
	p.lock.Unlock()

	fired := bag.AddValue("A_vol", param)
	if fired {
		p.lock.Lock()
		delete(p.params, param.ID)
		p.lock.Unlock()
	}
}

func (p *ParamBlock) OnB_vol(param execute.ThreadParam) {
	p.lock.Lock()
	var bag *execute.AsyncBag = p.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, p, p)
		p.params[param.ID] = bag
	}
	p.lock.Unlock()

	fired := bag.AddValue("B_vol", param)
	if fired {
		p.lock.Lock()
		delete(p.params, param.ID)
		p.lock.Unlock()
	}
}

// execute.AsyncBag functions
func (e *LHExample) Complete(blocks interface{}) {
	blx := blocks.(BlockSet)
	p := blx.ParamBlock
	i := blx.InputBlock
	e.setup(p, i)
	e.steps(p, i)
}

// generic typing for interface support
func (e *LHExample) anthaElement() {}

// init function, read characterization info from seperate file to validate ranges?
func (e *LHExample) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func NewLHExample() *LHExample {
	e := new(LHExample)
	e.init()
	return e
}

// copying the above block type structure for inputs and outputs

// we use separate input and output structures
// to allow each to have its own AsyncBag
type InputBlock struct {
	A    wtype.Liquid
	B    wtype.Liquid
	Dest wtype.Labware
}

type JSONInputBlock struct {
	A    *wtype.Liquid
	B    *wtype.Liquid
	Dest *wtype.Labware
}

// InputBlock needs to implement AsyncBag

func (ib *InputBlock) Complete() {

}

func (ib *InputBlock) Map(m map[string]interface{}) interface{} {
}

func (ib *InputBlock) OnDest(param execute.ThreadParam) {
	ib.lock.Lock()
	var bag *execute.AsyncBag = ib.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, ib, ib)
		ib.params[param.ID] = bag
	}
	ib.lock.Unlock()

	fired := bag.AddValue("Dest", param)
	if fired {
		ib.lock.Lock()
		delete(ib.params, param.ID)
		ib.lock.Unlock()
	}
}

func (ib *InputBlock) OnA(param execute.ThreadParam) {
	ib.lock.Lock()
	var bag *execute.AsyncBag = ib.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, ib, ib)
		ib.params[param.ID] = bag
	}
	ib.lock.Unlock()

	fired := bag.AddValue("A", param)
	if fired {
		ib.lock.Lock()
		delete(ib.params, param.ID)
		ib.lock.Unlock()
	}

}

func (ib *InputBlock) OnB(param execute.ThreadParam) {
	ib.lock.Lock()
	var bag *execute.AsyncBag = ib.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, ib, ib)
		ib.params[param.ID] = bag
	}
	ib.lock.Unlock()

	fired := bag.AddValue("B", param)
	if fired {
		ib.lock.Lock()
		delete(ib.params, param.ID)
		ib.lock.Unlock()
	}

}

type OutputBlock struct {
	SolOut wtype.Solution
}

// setup function

func (e *LHExample) setup(pb ParamBlock, ib InputBlock) {

}

// main function for use in goroutines
func (e *LHExample) steps(pb ParamBlock, ib InputBlock) {
	// get the execution context

	ctx := execution.GetContext()

	sample_a := mixer.Sample(ib.A, p.A_vol)
	sample_b := mixer.Sample(ib.B, p.B_vol)

	ret := MixInto(ib.Dest, sample_a, sample_b)

	ob.SolOut = ret
}
