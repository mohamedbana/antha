package Recovery

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

//= 2 (hours)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Recovery) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Recovery) setup(p RecoveryParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Recovery) steps(p RecoveryParamBlock, r *RecoveryResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	recoverymix := make([]*wtype.LHComponent, 0)

	transformedcellsComp := p.Transformedcells

	recoverymixture := mixer.Sample(p.Recoverymedium, p.Recoveryvolume)

	recoverymix = append(recoverymix, transformedcellsComp, recoverymixture)
	recoverymix2 := _wrapper.MixInto(p.OutPlate, recoverymix...)

	_wrapper.Incubate(recoverymix2, p.Recoverytemp, p.Recoverytime, true)

	r.RecoveredCells = recoverymix2
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Recovery) analysis(p RecoveryParamBlock, r *RecoveryResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Recovery) validation(p RecoveryParamBlock, r *RecoveryResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Recovery) Complete(params interface{}) {
	p := params.(RecoveryParamBlock)
	if p.Error {
		e.RecoveredCells <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(RecoveryResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.RecoveredCells <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.RecoveredCells <- execute.ThreadParam{Value: r.RecoveredCells, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Recovery) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Recovery) NewConfig() interface{} {
	return &RecoveryConfig{}
}

func (e *Recovery) NewParamBlock() interface{} {
	return &RecoveryParamBlock{}
}

func NewRecovery() interface{} { //*Recovery {
	e := new(Recovery)
	e.init()
	return e
}

// Mapper function
func (e *Recovery) Map(m map[string]interface{}) interface{} {
	var res RecoveryParamBlock
	res.Error = false || m["AgarPlate"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Recoverymedium"].(execute.ThreadParam).Error || m["Recoverytemp"].(execute.ThreadParam).Error || m["Recoverytime"].(execute.ThreadParam).Error || m["Recoveryvolume"].(execute.ThreadParam).Error || m["Transformedcells"].(execute.ThreadParam).Error

	vAgarPlate, is := m["AgarPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vAgarPlate.JSONString), &temp)
		res.AgarPlate = *temp.AgarPlate
	} else {
		res.AgarPlate = m["AgarPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vRecoverymedium, is := m["Recoverymedium"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vRecoverymedium.JSONString), &temp)
		res.Recoverymedium = *temp.Recoverymedium
	} else {
		res.Recoverymedium = m["Recoverymedium"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vRecoverytemp, is := m["Recoverytemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vRecoverytemp.JSONString), &temp)
		res.Recoverytemp = *temp.Recoverytemp
	} else {
		res.Recoverytemp = m["Recoverytemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vRecoverytime, is := m["Recoverytime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vRecoverytime.JSONString), &temp)
		res.Recoverytime = *temp.Recoverytime
	} else {
		res.Recoverytime = m["Recoverytime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vRecoveryvolume, is := m["Recoveryvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vRecoveryvolume.JSONString), &temp)
		res.Recoveryvolume = *temp.Recoveryvolume
	} else {
		res.Recoveryvolume = m["Recoveryvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vTransformedcells, is := m["Transformedcells"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RecoveryJSONBlock
		json.Unmarshal([]byte(vTransformedcells.JSONString), &temp)
		res.Transformedcells = *temp.Transformedcells
	} else {
		res.Transformedcells = m["Transformedcells"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["AgarPlate"].(execute.ThreadParam).ID
	res.BlockID = m["AgarPlate"].(execute.ThreadParam).BlockID

	return res
}

func (e *Recovery) OnAgarPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AgarPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Recovery) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *Recovery) OnRecoverymedium(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoverymedium", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Recovery) OnRecoverytemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoverytemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Recovery) OnRecoverytime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoverytime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Recovery) OnRecoveryvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoveryvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Recovery) OnTransformedcells(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Transformedcells", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Recovery struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	AgarPlate        <-chan execute.ThreadParam
	OutPlate         <-chan execute.ThreadParam
	Recoverymedium   <-chan execute.ThreadParam
	Recoverytemp     <-chan execute.ThreadParam
	Recoverytime     <-chan execute.ThreadParam
	Recoveryvolume   <-chan execute.ThreadParam
	Transformedcells <-chan execute.ThreadParam
	RecoveredCells   chan<- execute.ThreadParam
}

type RecoveryParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	AgarPlate        *wtype.LHPlate
	OutPlate         *wtype.LHPlate
	Recoverymedium   *wtype.LHComponent
	Recoverytemp     wunit.Temperature
	Recoverytime     wunit.Time
	Recoveryvolume   wunit.Volume
	Transformedcells *wtype.LHComponent
}

type RecoveryConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	AgarPlate        wtype.FromFactory
	OutPlate         wtype.FromFactory
	Recoverymedium   wtype.FromFactory
	Recoverytemp     wunit.Temperature
	Recoverytime     wunit.Time
	Recoveryvolume   wunit.Volume
	Transformedcells wtype.FromFactory
}

type RecoveryResultBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	RecoveredCells *wtype.LHSolution
}

type RecoveryJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	AgarPlate        **wtype.LHPlate
	OutPlate         **wtype.LHPlate
	Recoverymedium   **wtype.LHComponent
	Recoverytemp     *wunit.Temperature
	Recoverytime     *wunit.Time
	Recoveryvolume   *wunit.Volume
	Transformedcells **wtype.LHComponent
	RecoveredCells   **wtype.LHSolution
}

func (c *Recovery) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AgarPlate", "*wtype.LHPlate", "AgarPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoverymedium", "*wtype.LHComponent", "Recoverymedium", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoverytemp", "wunit.Temperature", "Recoverytemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoverytime", "wunit.Time", "Recoverytime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoveryvolume", "wunit.Volume", "Recoveryvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Transformedcells", "*wtype.LHComponent", "Transformedcells", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("RecoveredCells", "*wtype.LHSolution", "RecoveredCells", true, true, nil, nil))

	ci := execute.NewComponentInfo("Recovery", "Recovery", "", false, inp, outp)

	return ci
}
