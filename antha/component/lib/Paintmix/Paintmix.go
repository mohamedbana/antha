package Paintmix

import (
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Paintmix) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Paintmix) setup(p PaintmixParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Paintmix) steps(p PaintmixParamBlock, r *PaintmixResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	reactions := make([]*wtype.LHSolution, 0)

	for i := 0; i < p.Numberofcopies; i++ {
		eachreaction := make([]*wtype.LHComponent, 0)
		col1Sample := mixer.Sample(p.Colour1, p.Colour1vol)
		eachreaction = append(eachreaction, col1Sample)
		col2Sample := mixer.Sample(p.Colour2, p.Colour2vol)
		eachreaction = append(eachreaction, col2Sample)
		reaction := _wrapper.MixInto(p.OutPlate, eachreaction...)
		reactions = append(reactions, reaction)

	}
	r.NewColours = reactions
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Paintmix) analysis(p PaintmixParamBlock, r *PaintmixResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Paintmix) validation(p PaintmixParamBlock, r *PaintmixResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Paintmix) Complete(params interface{}) {
	p := params.(PaintmixParamBlock)
	if p.Error {
		e.NewColours <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PaintmixResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.NewColours <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.NewColours <- execute.ThreadParam{Value: r.NewColours, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Paintmix) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Paintmix) NewConfig() interface{} {
	return &PaintmixConfig{}
}

func (e *Paintmix) NewParamBlock() interface{} {
	return &PaintmixParamBlock{}
}

func NewPaintmix() interface{} { //*Paintmix {
	e := new(Paintmix)
	e.init()
	return e
}

// Mapper function
func (e *Paintmix) Map(m map[string]interface{}) interface{} {
	var res PaintmixParamBlock
	res.Error = false || m["Colour1"].(execute.ThreadParam).Error || m["Colour1vol"].(execute.ThreadParam).Error || m["Colour2"].(execute.ThreadParam).Error || m["Colour2vol"].(execute.ThreadParam).Error || m["Numberofcopies"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error

	vColour1, is := m["Colour1"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PaintmixJSONBlock
		json.Unmarshal([]byte(vColour1.JSONString), &temp)
		res.Colour1 = *temp.Colour1
	} else {
		res.Colour1 = m["Colour1"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vColour1vol, is := m["Colour1vol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PaintmixJSONBlock
		json.Unmarshal([]byte(vColour1vol.JSONString), &temp)
		res.Colour1vol = *temp.Colour1vol
	} else {
		res.Colour1vol = m["Colour1vol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vColour2, is := m["Colour2"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PaintmixJSONBlock
		json.Unmarshal([]byte(vColour2.JSONString), &temp)
		res.Colour2 = *temp.Colour2
	} else {
		res.Colour2 = m["Colour2"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vColour2vol, is := m["Colour2vol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PaintmixJSONBlock
		json.Unmarshal([]byte(vColour2vol.JSONString), &temp)
		res.Colour2vol = *temp.Colour2vol
	} else {
		res.Colour2vol = m["Colour2vol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vNumberofcopies, is := m["Numberofcopies"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PaintmixJSONBlock
		json.Unmarshal([]byte(vNumberofcopies.JSONString), &temp)
		res.Numberofcopies = *temp.Numberofcopies
	} else {
		res.Numberofcopies = m["Numberofcopies"].(execute.ThreadParam).Value.(int)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PaintmixJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	res.ID = m["Colour1"].(execute.ThreadParam).ID
	res.BlockID = m["Colour1"].(execute.ThreadParam).BlockID

	return res
}

func (e *Paintmix) OnColour1(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Colour1", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Paintmix) OnColour1vol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Colour1vol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Paintmix) OnColour2(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Colour2", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Paintmix) OnColour2vol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Colour2vol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Paintmix) OnNumberofcopies(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Numberofcopies", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Paintmix) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OutPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Paintmix struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	Colour1        <-chan execute.ThreadParam
	Colour1vol     <-chan execute.ThreadParam
	Colour2        <-chan execute.ThreadParam
	Colour2vol     <-chan execute.ThreadParam
	Numberofcopies <-chan execute.ThreadParam
	OutPlate       <-chan execute.ThreadParam
	NewColours     chan<- execute.ThreadParam
	Status         chan<- execute.ThreadParam
}

type PaintmixParamBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	Colour1        *wtype.LHComponent
	Colour1vol     wunit.Volume
	Colour2        *wtype.LHComponent
	Colour2vol     wunit.Volume
	Numberofcopies int
	OutPlate       *wtype.LHPlate
}

type PaintmixConfig struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	Colour1        wtype.FromFactory
	Colour1vol     wunit.Volume
	Colour2        wtype.FromFactory
	Colour2vol     wunit.Volume
	Numberofcopies int
	OutPlate       wtype.FromFactory
}

type PaintmixResultBlock struct {
	ID         execute.ThreadID
	BlockID    execute.BlockID
	Error      bool
	NewColours []*wtype.LHSolution
	Status     string
}

type PaintmixJSONBlock struct {
	ID             *execute.ThreadID
	BlockID        *execute.BlockID
	Error          *bool
	Colour1        **wtype.LHComponent
	Colour1vol     *wunit.Volume
	Colour2        **wtype.LHComponent
	Colour2vol     *wunit.Volume
	Numberofcopies *int
	OutPlate       **wtype.LHPlate
	NewColours     *[]*wtype.LHSolution
	Status         *string
}

func (c *Paintmix) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Colour1", "*wtype.LHComponent", "Colour1", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Colour1vol", "wunit.Volume", "Colour1vol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Colour2", "*wtype.LHComponent", "Colour2", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Colour2vol", "wunit.Volume", "Colour2vol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Numberofcopies", "int", "Numberofcopies", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("NewColours", "[]*wtype.LHSolution", "NewColours", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Paintmix", "Paintmix", "", false, inp, outp)

	return ci
}
