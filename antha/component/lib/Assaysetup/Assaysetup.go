package Assaysetup

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

func (e *Assaysetup) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Assaysetup) setup(p AssaysetupParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Assaysetup) steps(p AssaysetupParamBlock, r *AssaysetupResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	reactions := make([]*wtype.LHSolution, 0)

	for i := 0; i < p.NumberofReactions; i++ {
		eachreaction := make([]*wtype.LHComponent, 0)
		bufferSample := mixer.SampleForTotalVolume(p.Buffer, p.TotalVolume)
		eachreaction = append(eachreaction, bufferSample)
		subSample := mixer.Sample(p.Substrate, p.SubstrateVolume)
		eachreaction = append(eachreaction, subSample)
		enzSample := mixer.Sample(p.Enzyme, p.EnzymeVolume)
		eachreaction = append(eachreaction, enzSample)
		reaction := _wrapper.MixInto(p.OutPlate, eachreaction...)
		reactions = append(reactions, reaction)

	}
	r.Reactions = reactions
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Assaysetup) analysis(p AssaysetupParamBlock, r *AssaysetupResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Assaysetup) validation(p AssaysetupParamBlock, r *AssaysetupResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Assaysetup) Complete(params interface{}) {
	p := params.(AssaysetupParamBlock)
	if p.Error {
		e.Reactions <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(AssaysetupResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Reactions <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Reactions <- execute.ThreadParam{Value: r.Reactions, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Assaysetup) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Assaysetup) NewConfig() interface{} {
	return &AssaysetupConfig{}
}

func (e *Assaysetup) NewParamBlock() interface{} {
	return &AssaysetupParamBlock{}
}

func NewAssaysetup() interface{} { //*Assaysetup {
	e := new(Assaysetup)
	e.init()
	return e
}

// Mapper function
func (e *Assaysetup) Map(m map[string]interface{}) interface{} {
	var res AssaysetupParamBlock
	res.Error = false || m["Buffer"].(execute.ThreadParam).Error || m["Enzyme"].(execute.ThreadParam).Error || m["EnzymeVolume"].(execute.ThreadParam).Error || m["NumberofReactions"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Substrate"].(execute.ThreadParam).Error || m["SubstrateVolume"].(execute.ThreadParam).Error || m["TotalVolume"].(execute.ThreadParam).Error

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vEnzyme, is := m["Enzyme"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vEnzyme.JSONString), &temp)
		res.Enzyme = *temp.Enzyme
	} else {
		res.Enzyme = m["Enzyme"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vEnzymeVolume, is := m["EnzymeVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vEnzymeVolume.JSONString), &temp)
		res.EnzymeVolume = *temp.EnzymeVolume
	} else {
		res.EnzymeVolume = m["EnzymeVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vNumberofReactions, is := m["NumberofReactions"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vNumberofReactions.JSONString), &temp)
		res.NumberofReactions = *temp.NumberofReactions
	} else {
		res.NumberofReactions = m["NumberofReactions"].(execute.ThreadParam).Value.(int)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vSubstrate, is := m["Substrate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vSubstrate.JSONString), &temp)
		res.Substrate = *temp.Substrate
	} else {
		res.Substrate = m["Substrate"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vSubstrateVolume, is := m["SubstrateVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vSubstrateVolume.JSONString), &temp)
		res.SubstrateVolume = *temp.SubstrateVolume
	} else {
		res.SubstrateVolume = m["SubstrateVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vTotalVolume, is := m["TotalVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AssaysetupJSONBlock
		json.Unmarshal([]byte(vTotalVolume.JSONString), &temp)
		res.TotalVolume = *temp.TotalVolume
	} else {
		res.TotalVolume = m["TotalVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["Buffer"].(execute.ThreadParam).ID
	res.BlockID = m["Buffer"].(execute.ThreadParam).BlockID

	return res
}

func (e *Assaysetup) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Buffer", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Assaysetup) OnEnzyme(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Enzyme", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Assaysetup) OnEnzymeVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("EnzymeVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Assaysetup) OnNumberofReactions(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("NumberofReactions", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Assaysetup) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *Assaysetup) OnSubstrate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Substrate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Assaysetup) OnSubstrateVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SubstrateVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Assaysetup) OnTotalVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TotalVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Assaysetup struct {
	flow.Component    // component "superclass" embedded
	lock              sync.Mutex
	startup           sync.Once
	params            map[execute.ThreadID]*execute.AsyncBag
	Buffer            <-chan execute.ThreadParam
	Enzyme            <-chan execute.ThreadParam
	EnzymeVolume      <-chan execute.ThreadParam
	NumberofReactions <-chan execute.ThreadParam
	OutPlate          <-chan execute.ThreadParam
	Substrate         <-chan execute.ThreadParam
	SubstrateVolume   <-chan execute.ThreadParam
	TotalVolume       <-chan execute.ThreadParam
	Reactions         chan<- execute.ThreadParam
	Status            chan<- execute.ThreadParam
}

type AssaysetupParamBlock struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Buffer            *wtype.LHComponent
	Enzyme            *wtype.LHComponent
	EnzymeVolume      wunit.Volume
	NumberofReactions int
	OutPlate          *wtype.LHPlate
	Substrate         *wtype.LHComponent
	SubstrateVolume   wunit.Volume
	TotalVolume       wunit.Volume
}

type AssaysetupConfig struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Buffer            wtype.FromFactory
	Enzyme            wtype.FromFactory
	EnzymeVolume      wunit.Volume
	NumberofReactions int
	OutPlate          wtype.FromFactory
	Substrate         wtype.FromFactory
	SubstrateVolume   wunit.Volume
	TotalVolume       wunit.Volume
}

type AssaysetupResultBlock struct {
	ID        execute.ThreadID
	BlockID   execute.BlockID
	Error     bool
	Reactions []*wtype.LHSolution
	Status    string
}

type AssaysetupJSONBlock struct {
	ID                *execute.ThreadID
	BlockID           *execute.BlockID
	Error             *bool
	Buffer            **wtype.LHComponent
	Enzyme            **wtype.LHComponent
	EnzymeVolume      *wunit.Volume
	NumberofReactions *int
	OutPlate          **wtype.LHPlate
	Substrate         **wtype.LHComponent
	SubstrateVolume   *wunit.Volume
	TotalVolume       *wunit.Volume
	Reactions         *[]*wtype.LHSolution
	Status            *string
}

func (c *Assaysetup) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Enzyme", "*wtype.LHComponent", "Enzyme", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("EnzymeVolume", "wunit.Volume", "EnzymeVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("NumberofReactions", "int", "NumberofReactions", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Substrate", "*wtype.LHComponent", "Substrate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SubstrateVolume", "wunit.Volume", "SubstrateVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TotalVolume", "wunit.Volume", "TotalVolume", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reactions", "[]*wtype.LHSolution", "Reactions", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Assaysetup", "Assaysetup", "", false, inp, outp)

	return ci
}
