package ProtocolName_from_an_file

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

func (e *ProtocolName_from_an_file) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *ProtocolName_from_an_file) setup(p ProtocolName_from_an_fileParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *ProtocolName_from_an_file) steps(p ProtocolName_from_an_fileParamBlock, r *ProtocolName_from_an_fileResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	r.OutputData = make([]string, 0)

	for i := 0; i < p.ParameterVariableAsint; i++ {
		output := p.ParameterVariableAsValuewithunit.ToString() + "of" + p.ParameterVariablestring
		r.OutputData = append(r.OutputData, output)
	}
	sample := mixer.Sample(p.InputVariable, p.ParameterVariableAsValuewithunit)
	r.PhysicalOutput = _wrapper.MixInto(p.OutPlate, sample)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *ProtocolName_from_an_file) analysis(p ProtocolName_from_an_fileParamBlock, r *ProtocolName_from_an_fileResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *ProtocolName_from_an_file) validation(p ProtocolName_from_an_fileParamBlock, r *ProtocolName_from_an_fileResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *ProtocolName_from_an_file) Complete(params interface{}) {
	p := params.(ProtocolName_from_an_fileParamBlock)
	if p.Error {
		e.OutputData <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.PhysicalOutput <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(ProtocolName_from_an_fileResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.OutputData <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.PhysicalOutput <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.OutputData <- execute.ThreadParam{Value: r.OutputData, ID: p.ID, Error: false}

	e.PhysicalOutput <- execute.ThreadParam{Value: r.PhysicalOutput, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *ProtocolName_from_an_file) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *ProtocolName_from_an_file) NewConfig() interface{} {
	return &ProtocolName_from_an_fileConfig{}
}

func (e *ProtocolName_from_an_file) NewParamBlock() interface{} {
	return &ProtocolName_from_an_fileParamBlock{}
}

func NewProtocolName_from_an_file() interface{} { //*ProtocolName_from_an_file {
	e := new(ProtocolName_from_an_file)
	e.init()
	return e
}

// Mapper function
func (e *ProtocolName_from_an_file) Map(m map[string]interface{}) interface{} {
	var res ProtocolName_from_an_fileParamBlock
	res.Error = false || m["InputVariable"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["ParameterVariableAsValuewithunit"].(execute.ThreadParam).Error || m["ParameterVariableAsint"].(execute.ThreadParam).Error || m["ParameterVariablestring"].(execute.ThreadParam).Error

	vInputVariable, is := m["InputVariable"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ProtocolName_from_an_fileJSONBlock
		json.Unmarshal([]byte(vInputVariable.JSONString), &temp)
		res.InputVariable = *temp.InputVariable
	} else {
		res.InputVariable = m["InputVariable"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ProtocolName_from_an_fileJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vParameterVariableAsValuewithunit, is := m["ParameterVariableAsValuewithunit"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ProtocolName_from_an_fileJSONBlock
		json.Unmarshal([]byte(vParameterVariableAsValuewithunit.JSONString), &temp)
		res.ParameterVariableAsValuewithunit = *temp.ParameterVariableAsValuewithunit
	} else {
		res.ParameterVariableAsValuewithunit = m["ParameterVariableAsValuewithunit"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vParameterVariableAsint, is := m["ParameterVariableAsint"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ProtocolName_from_an_fileJSONBlock
		json.Unmarshal([]byte(vParameterVariableAsint.JSONString), &temp)
		res.ParameterVariableAsint = *temp.ParameterVariableAsint
	} else {
		res.ParameterVariableAsint = m["ParameterVariableAsint"].(execute.ThreadParam).Value.(int)
	}

	vParameterVariablestring, is := m["ParameterVariablestring"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ProtocolName_from_an_fileJSONBlock
		json.Unmarshal([]byte(vParameterVariablestring.JSONString), &temp)
		res.ParameterVariablestring = *temp.ParameterVariablestring
	} else {
		res.ParameterVariablestring = m["ParameterVariablestring"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["InputVariable"].(execute.ThreadParam).ID
	res.BlockID = m["InputVariable"].(execute.ThreadParam).BlockID

	return res
}

func (e *ProtocolName_from_an_file) OnInputVariable(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("InputVariable", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *ProtocolName_from_an_file) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
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
func (e *ProtocolName_from_an_file) OnParameterVariableAsValuewithunit(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ParameterVariableAsValuewithunit", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *ProtocolName_from_an_file) OnParameterVariableAsint(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ParameterVariableAsint", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *ProtocolName_from_an_file) OnParameterVariablestring(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ParameterVariablestring", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type ProtocolName_from_an_file struct {
	flow.Component                   // component "superclass" embedded
	lock                             sync.Mutex
	startup                          sync.Once
	params                           map[execute.ThreadID]*execute.AsyncBag
	InputVariable                    <-chan execute.ThreadParam
	OutPlate                         <-chan execute.ThreadParam
	ParameterVariableAsValuewithunit <-chan execute.ThreadParam
	ParameterVariableAsint           <-chan execute.ThreadParam
	ParameterVariablestring          <-chan execute.ThreadParam
	OutputData                       chan<- execute.ThreadParam
	PhysicalOutput                   chan<- execute.ThreadParam
}

type ProtocolName_from_an_fileParamBlock struct {
	ID                               execute.ThreadID
	BlockID                          execute.BlockID
	Error                            bool
	InputVariable                    *wtype.LHComponent
	OutPlate                         *wtype.LHPlate
	ParameterVariableAsValuewithunit wunit.Volume
	ParameterVariableAsint           int
	ParameterVariablestring          string
}

type ProtocolName_from_an_fileConfig struct {
	ID                               execute.ThreadID
	BlockID                          execute.BlockID
	Error                            bool
	InputVariable                    wtype.FromFactory
	OutPlate                         wtype.FromFactory
	ParameterVariableAsValuewithunit wunit.Volume
	ParameterVariableAsint           int
	ParameterVariablestring          string
}

type ProtocolName_from_an_fileResultBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	OutputData     []string
	PhysicalOutput *wtype.LHSolution
}

type ProtocolName_from_an_fileJSONBlock struct {
	ID                               *execute.ThreadID
	BlockID                          *execute.BlockID
	Error                            *bool
	InputVariable                    **wtype.LHComponent
	OutPlate                         **wtype.LHPlate
	ParameterVariableAsValuewithunit *wunit.Volume
	ParameterVariableAsint           *int
	ParameterVariablestring          *string
	OutputData                       *[]string
	PhysicalOutput                   **wtype.LHSolution
}

func (c *ProtocolName_from_an_file) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("InputVariable", "*wtype.LHComponent", "InputVariable", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ParameterVariableAsValuewithunit", "wunit.Volume", "ParameterVariableAsValuewithunit", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ParameterVariableAsint", "int", "ParameterVariableAsint", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ParameterVariablestring", "string", "ParameterVariablestring", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("OutputData", "[]string", "OutputData", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("PhysicalOutput", "*wtype.LHSolution", "PhysicalOutput", true, true, nil, nil))

	ci := execute.NewComponentInfo("ProtocolName_from_an_file", "ProtocolName_from_an_file", "", false, inp, outp)

	return ci
}
