package PlateOut

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

func (e *PlateOut) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *PlateOut) setup(p PlateOutParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *PlateOut) steps(p PlateOutParamBlock, r *PlateOutResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	plateout := make([]*wtype.LHComponent, 0)

	if p.Diluent != nil && p.DilutionX > 1 {
		diluentsample := mixer.SampleForTotalVolume(p.Diluent, p.Plateoutvolume)
		plateout = append(plateout, diluentsample)
		// redeclare Plateoutvolume for adjusted volume to add of recovery mixture based on dilution ratio
		p.Plateoutvolume = wunit.NewVolume(p.Plateoutvolume.SIValue()/float64(p.DilutionX), "l")
	}
	plateoutsample := mixer.Sample(p.RecoveredCells, p.Plateoutvolume)
	plateout = append(plateout, plateoutsample)
	platedculture := _wrapper.MixInto(p.AgarPlate, plateout...)
	_wrapper.Incubate(platedculture, p.IncubationTemp, p.IncubationTime, false)
	r.Platedculture = platedculture
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *PlateOut) analysis(p PlateOutParamBlock, r *PlateOutResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *PlateOut) validation(p PlateOutParamBlock, r *PlateOutResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *PlateOut) Complete(params interface{}) {
	p := params.(PlateOutParamBlock)
	if p.Error {
		e.Platedculture <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PlateOutResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Platedculture <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Platedculture <- execute.ThreadParam{Value: r.Platedculture, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *PlateOut) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *PlateOut) NewConfig() interface{} {
	return &PlateOutConfig{}
}

func (e *PlateOut) NewParamBlock() interface{} {
	return &PlateOutParamBlock{}
}

func NewPlateOut() interface{} { //*PlateOut {
	e := new(PlateOut)
	e.init()
	return e
}

// Mapper function
func (e *PlateOut) Map(m map[string]interface{}) interface{} {
	var res PlateOutParamBlock
	res.Error = false || m["AgarPlate"].(execute.ThreadParam).Error || m["Diluent"].(execute.ThreadParam).Error || m["DilutionX"].(execute.ThreadParam).Error || m["IncubationTemp"].(execute.ThreadParam).Error || m["IncubationTime"].(execute.ThreadParam).Error || m["Plateoutvolume"].(execute.ThreadParam).Error || m["RecoveredCells"].(execute.ThreadParam).Error

	vAgarPlate, is := m["AgarPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vAgarPlate.JSONString), &temp)
		res.AgarPlate = *temp.AgarPlate
	} else {
		res.AgarPlate = m["AgarPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vDiluent, is := m["Diluent"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vDiluent.JSONString), &temp)
		res.Diluent = *temp.Diluent
	} else {
		res.Diluent = m["Diluent"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDilutionX, is := m["DilutionX"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vDilutionX.JSONString), &temp)
		res.DilutionX = *temp.DilutionX
	} else {
		res.DilutionX = m["DilutionX"].(execute.ThreadParam).Value.(int)
	}

	vIncubationTemp, is := m["IncubationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vIncubationTemp.JSONString), &temp)
		res.IncubationTemp = *temp.IncubationTemp
	} else {
		res.IncubationTemp = m["IncubationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vIncubationTime, is := m["IncubationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vIncubationTime.JSONString), &temp)
		res.IncubationTime = *temp.IncubationTime
	} else {
		res.IncubationTime = m["IncubationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vPlateoutvolume, is := m["Plateoutvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vPlateoutvolume.JSONString), &temp)
		res.Plateoutvolume = *temp.Plateoutvolume
	} else {
		res.Plateoutvolume = m["Plateoutvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vRecoveredCells, is := m["RecoveredCells"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlateOutJSONBlock
		json.Unmarshal([]byte(vRecoveredCells.JSONString), &temp)
		res.RecoveredCells = *temp.RecoveredCells
	} else {
		res.RecoveredCells = m["RecoveredCells"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["AgarPlate"].(execute.ThreadParam).ID
	res.BlockID = m["AgarPlate"].(execute.ThreadParam).BlockID

	return res
}

func (e *PlateOut) OnAgarPlate(param execute.ThreadParam) {
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
func (e *PlateOut) OnDiluent(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Diluent", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PlateOut) OnDilutionX(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DilutionX", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PlateOut) OnIncubationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("IncubationTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PlateOut) OnIncubationTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("IncubationTime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PlateOut) OnPlateoutvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Plateoutvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PlateOut) OnRecoveredCells(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RecoveredCells", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type PlateOut struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	AgarPlate      <-chan execute.ThreadParam
	Diluent        <-chan execute.ThreadParam
	DilutionX      <-chan execute.ThreadParam
	IncubationTemp <-chan execute.ThreadParam
	IncubationTime <-chan execute.ThreadParam
	Plateoutvolume <-chan execute.ThreadParam
	RecoveredCells <-chan execute.ThreadParam
	Platedculture  chan<- execute.ThreadParam
}

type PlateOutParamBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	AgarPlate      *wtype.LHPlate
	Diluent        *wtype.LHComponent
	DilutionX      int
	IncubationTemp wunit.Temperature
	IncubationTime wunit.Time
	Plateoutvolume wunit.Volume
	RecoveredCells *wtype.LHComponent
}

type PlateOutConfig struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	AgarPlate      wtype.FromFactory
	Diluent        wtype.FromFactory
	DilutionX      int
	IncubationTemp wunit.Temperature
	IncubationTime wunit.Time
	Plateoutvolume wunit.Volume
	RecoveredCells wtype.FromFactory
}

type PlateOutResultBlock struct {
	ID            execute.ThreadID
	BlockID       execute.BlockID
	Error         bool
	Platedculture *wtype.LHSolution
}

type PlateOutJSONBlock struct {
	ID             *execute.ThreadID
	BlockID        *execute.BlockID
	Error          *bool
	AgarPlate      **wtype.LHPlate
	Diluent        **wtype.LHComponent
	DilutionX      *int
	IncubationTemp *wunit.Temperature
	IncubationTime *wunit.Time
	Plateoutvolume *wunit.Volume
	RecoveredCells **wtype.LHComponent
	Platedculture  **wtype.LHSolution
}

func (c *PlateOut) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AgarPlate", "*wtype.LHPlate", "AgarPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Diluent", "*wtype.LHComponent", "Diluent", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DilutionX", "int", "DilutionX", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("IncubationTemp", "wunit.Temperature", "IncubationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("IncubationTime", "wunit.Time", "IncubationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Plateoutvolume", "wunit.Volume", "Plateoutvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RecoveredCells", "*wtype.LHComponent", "RecoveredCells", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Platedculture", "*wtype.LHSolution", "Platedculture", true, true, nil, nil))

	ci := execute.NewComponentInfo("PlateOut", "PlateOut", "", false, inp, outp)

	return ci
}
