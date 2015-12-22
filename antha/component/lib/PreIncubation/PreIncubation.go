package PreIncubation

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

//= 50.(uL)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *PreIncubation) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *PreIncubation) setup(p PreIncubationParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *PreIncubation) steps(p PreIncubationParamBlock, r *PreIncubationResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	competentcells := make([]*wtype.LHComponent, 0)
	competentcellsample := mixer.Sample(p.CompetentCells, p.CompetentCellvolumeperassembly)
	competentcells = append(competentcells, competentcellsample)
	readycompetentcells := _wrapper.MixInto(p.OutPlate, competentcells...)
	_wrapper.Incubate(readycompetentcells, p.Preplasmidtemp, p.Preplasmidtime, false)

	r.ReadyCompCells = readycompetentcells
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *PreIncubation) analysis(p PreIncubationParamBlock, r *PreIncubationResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *PreIncubation) validation(p PreIncubationParamBlock, r *PreIncubationResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *PreIncubation) Complete(params interface{}) {
	p := params.(PreIncubationParamBlock)
	if p.Error {
		e.ReadyCompCells <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PreIncubationResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.ReadyCompCells <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.ReadyCompCells <- execute.ThreadParam{Value: r.ReadyCompCells, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *PreIncubation) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *PreIncubation) NewConfig() interface{} {
	return &PreIncubationConfig{}
}

func (e *PreIncubation) NewParamBlock() interface{} {
	return &PreIncubationParamBlock{}
}

func NewPreIncubation() interface{} { //*PreIncubation {
	e := new(PreIncubation)
	e.init()
	return e
}

// Mapper function
func (e *PreIncubation) Map(m map[string]interface{}) interface{} {
	var res PreIncubationParamBlock
	res.Error = false || m["CompetentCells"].(execute.ThreadParam).Error || m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Preplasmidtemp"].(execute.ThreadParam).Error || m["Preplasmidtime"].(execute.ThreadParam).Error

	vCompetentCells, is := m["CompetentCells"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PreIncubationJSONBlock
		json.Unmarshal([]byte(vCompetentCells.JSONString), &temp)
		res.CompetentCells = *temp.CompetentCells
	} else {
		res.CompetentCells = m["CompetentCells"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vCompetentCellvolumeperassembly, is := m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PreIncubationJSONBlock
		json.Unmarshal([]byte(vCompetentCellvolumeperassembly.JSONString), &temp)
		res.CompetentCellvolumeperassembly = *temp.CompetentCellvolumeperassembly
	} else {
		res.CompetentCellvolumeperassembly = m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PreIncubationJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPreplasmidtemp, is := m["Preplasmidtemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PreIncubationJSONBlock
		json.Unmarshal([]byte(vPreplasmidtemp.JSONString), &temp)
		res.Preplasmidtemp = *temp.Preplasmidtemp
	} else {
		res.Preplasmidtemp = m["Preplasmidtemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vPreplasmidtime, is := m["Preplasmidtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PreIncubationJSONBlock
		json.Unmarshal([]byte(vPreplasmidtime.JSONString), &temp)
		res.Preplasmidtime = *temp.Preplasmidtime
	} else {
		res.Preplasmidtime = m["Preplasmidtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	res.ID = m["CompetentCells"].(execute.ThreadParam).ID
	res.BlockID = m["CompetentCells"].(execute.ThreadParam).BlockID

	return res
}

func (e *PreIncubation) OnCompetentCells(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("CompetentCells", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PreIncubation) OnCompetentCellvolumeperassembly(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("CompetentCellvolumeperassembly", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PreIncubation) OnOutPlate(param execute.ThreadParam) {
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
func (e *PreIncubation) OnPreplasmidtemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Preplasmidtemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PreIncubation) OnPreplasmidtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Preplasmidtime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type PreIncubation struct {
	flow.Component                 // component "superclass" embedded
	lock                           sync.Mutex
	startup                        sync.Once
	params                         map[execute.ThreadID]*execute.AsyncBag
	CompetentCells                 <-chan execute.ThreadParam
	CompetentCellvolumeperassembly <-chan execute.ThreadParam
	OutPlate                       <-chan execute.ThreadParam
	Preplasmidtemp                 <-chan execute.ThreadParam
	Preplasmidtime                 <-chan execute.ThreadParam
	ReadyCompCells                 chan<- execute.ThreadParam
}

type PreIncubationParamBlock struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	CompetentCells                 *wtype.LHComponent
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       *wtype.LHPlate
	Preplasmidtemp                 wunit.Temperature
	Preplasmidtime                 wunit.Time
}

type PreIncubationConfig struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	CompetentCells                 wtype.FromFactory
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       wtype.FromFactory
	Preplasmidtemp                 wunit.Temperature
	Preplasmidtime                 wunit.Time
}

type PreIncubationResultBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	ReadyCompCells *wtype.LHSolution
}

type PreIncubationJSONBlock struct {
	ID                             *execute.ThreadID
	BlockID                        *execute.BlockID
	Error                          *bool
	CompetentCells                 **wtype.LHComponent
	CompetentCellvolumeperassembly *wunit.Volume
	OutPlate                       **wtype.LHPlate
	Preplasmidtemp                 *wunit.Temperature
	Preplasmidtime                 *wunit.Time
	ReadyCompCells                 **wtype.LHSolution
}

func (c *PreIncubation) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("CompetentCells", "*wtype.LHComponent", "CompetentCells", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("CompetentCellvolumeperassembly", "wunit.Volume", "CompetentCellvolumeperassembly", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Preplasmidtemp", "wunit.Temperature", "Preplasmidtemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Preplasmidtime", "wunit.Time", "Preplasmidtime", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("ReadyCompCells", "*wtype.LHSolution", "ReadyCompCells", true, true, nil, nil))

	ci := execute.NewComponentInfo("PreIncubation", "PreIncubation", "", false, inp, outp)

	return ci
}
