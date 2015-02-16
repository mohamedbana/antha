package lhreference

import (
	"bytes"
	"encoding/json"
	"github.com/Synthace/antha/execute"
	"github.com/Synthace/goflow"
	"io"
	"log"
	"sync"
	"time"
	"github.com/Synthace/antha/anthalib/wtype"
	"github.com/Synthace/antha/anthalib/wunit"
	"github.com/Synthace/antha/anthalib/liquidhandling"
	"github.com/Synthace/antha/anthalib/mixer"
)

var params = [...]string{
	"A_vol",
	"B_vol"}

// channel interfaces
// with threadID grouped types
type LHExample struct {
	flow.Component                            // component "superclass" embedded
	A_vol		<-chan execute.ThreadParam // volume of "A" to put in
	B_vol		<-chan execute.ThreadParam // volume of "B" to put in
	Dest		<-chan execute.ThreadParam // Microplate to mix into 
	S	       chan<- execute.ThreadParam // output solution
	lock           sync.Mutex
	params         map[execute.ThreadID]*execute.AsyncBag
}

// single execution thread variables
// with concrete types
type ParamBlock struct {
	A_vol	wunit.Volume
	B_vol	wunit.Volume
	Dest	wtype.LiquidContainer
	ID	execute.ThreadID
}

type JSONBlock struct {
	A_vol	*wunit.Volume
	B_vol	*wunit.Volume
	Dest	*wtype.LiquidContainer
	ID	*execute.ThreadID
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

// could handle mapping in the threadID better...
func (e *LHExample) Map(m map[string]interface{}) interface{} {
	var res ParamBlock
	res.A_vol = m["A_vol"].(execute.ThreadParam).Value.(wunit.Volume)
	res.B_vol = m["B_vol"].(execute.ThreadParam).Value.(wunit.Volume)
	res.Dest = m["Dest"].(execute.ThreadParam).Value.(wtype.LiquidContainer)
	m["ID"]=m["A_vol"].(execute.ThreadParam).Value.(execute.ThreadID)
	return res
}

func (e *LHExample) OnA_vol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("A_vol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

func (e *LHExample) OnB_vol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("B_vol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

func (e *LHExample) OnDest(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Dest", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

// execute.AsyncBag functions
func (e *LHExample) Complete(params interface{}) {
	p := params.(ParamBlock)
	e.steps(p)
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

// setup function 

func (e *LHExample) setup(p ParamBlock){


}

// main function for use in goroutines
func (e *LHExample) steps(p ParamBlock) {
	//time.Sleep(p.SleepTime)
	//e.WellColor <- execute.ThreadParam{p.Color, p.ID}

}
