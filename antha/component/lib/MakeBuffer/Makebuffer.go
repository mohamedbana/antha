package MakeBuffer

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func (e *MakeBuffer) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *MakeBuffer) setup(p MakeBufferParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *MakeBuffer) steps(p MakeBufferParamBlock, r *MakeBufferResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	//Bufferstockvolume := wunit.NewVolume((FinalVolume.SIValue() * FinalConcentration.SIValue()/Bufferstockconc.SIValue()),"l")

	r.Buffer = _wrapper.MixInto(p.OutPlate,
		mixer.Sample(p.Bufferstock, p.Bufferstockvolume),
		mixer.Sample(p.Diluent, p.Diluentvolume))

	r.Status = fmt.Sprintln("Buffer stock volume = ", p.Bufferstockvolume.ToString(), "of", p.Bufferstock.CName,
		"was added to ", p.Diluentvolume.ToString(), "of", p.Diluent.CName,
		"to make ", p.FinalVolume.ToString(), "of", p.Buffername,
		"Buffer stock conc =", p.Bufferstockconc.ToString())
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *MakeBuffer) analysis(p MakeBufferParamBlock, r *MakeBufferResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *MakeBuffer) validation(p MakeBufferParamBlock, r *MakeBufferResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *MakeBuffer) Complete(params interface{}) {
	p := params.(MakeBufferParamBlock)
	if p.Error {
		e.Buffer <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(MakeBufferResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Buffer <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Buffer <- execute.ThreadParam{Value: r.Buffer, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *MakeBuffer) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *MakeBuffer) NewConfig() interface{} {
	return &MakeBufferConfig{}
}

func (e *MakeBuffer) NewParamBlock() interface{} {
	return &MakeBufferParamBlock{}
}

func NewMakeBuffer() interface{} { //*MakeBuffer {
	e := new(MakeBuffer)
	e.init()
	return e
}

// Mapper function
func (e *MakeBuffer) Map(m map[string]interface{}) interface{} {
	var res MakeBufferParamBlock
	res.Error = false || m["Buffername"].(execute.ThreadParam).Error || m["Bufferstock"].(execute.ThreadParam).Error || m["Bufferstockconc"].(execute.ThreadParam).Error || m["Bufferstockvolume"].(execute.ThreadParam).Error || m["Diluent"].(execute.ThreadParam).Error || m["Diluentname"].(execute.ThreadParam).Error || m["Diluentvolume"].(execute.ThreadParam).Error || m["FinalConcentration"].(execute.ThreadParam).Error || m["FinalVolume"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error

	vBuffername, is := m["Buffername"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vBuffername.JSONString), &temp)
		res.Buffername = *temp.Buffername
	} else {
		res.Buffername = m["Buffername"].(execute.ThreadParam).Value.(string)
	}

	vBufferstock, is := m["Bufferstock"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vBufferstock.JSONString), &temp)
		res.Bufferstock = *temp.Bufferstock
	} else {
		res.Bufferstock = m["Bufferstock"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferstockconc, is := m["Bufferstockconc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vBufferstockconc.JSONString), &temp)
		res.Bufferstockconc = *temp.Bufferstockconc
	} else {
		res.Bufferstockconc = m["Bufferstockconc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vBufferstockvolume, is := m["Bufferstockvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vBufferstockvolume.JSONString), &temp)
		res.Bufferstockvolume = *temp.Bufferstockvolume
	} else {
		res.Bufferstockvolume = m["Bufferstockvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vDiluent, is := m["Diluent"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vDiluent.JSONString), &temp)
		res.Diluent = *temp.Diluent
	} else {
		res.Diluent = m["Diluent"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDiluentname, is := m["Diluentname"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vDiluentname.JSONString), &temp)
		res.Diluentname = *temp.Diluentname
	} else {
		res.Diluentname = m["Diluentname"].(execute.ThreadParam).Value.(string)
	}

	vDiluentvolume, is := m["Diluentvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vDiluentvolume.JSONString), &temp)
		res.Diluentvolume = *temp.Diluentvolume
	} else {
		res.Diluentvolume = m["Diluentvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vFinalConcentration, is := m["FinalConcentration"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vFinalConcentration.JSONString), &temp)
		res.FinalConcentration = *temp.FinalConcentration
	} else {
		res.FinalConcentration = m["FinalConcentration"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vFinalVolume, is := m["FinalVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vFinalVolume.JSONString), &temp)
		res.FinalVolume = *temp.FinalVolume
	} else {
		res.FinalVolume = m["FinalVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeBufferJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	res.ID = m["Buffername"].(execute.ThreadParam).ID
	res.BlockID = m["Buffername"].(execute.ThreadParam).BlockID

	return res
}

/*
type Mole struct {
	number float64
}*/

func (e *MakeBuffer) OnBuffername(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Buffername", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnBufferstock(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Bufferstock", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnBufferstockconc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Bufferstockconc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnBufferstockvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Bufferstockvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnDiluent(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
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
func (e *MakeBuffer) OnDiluentname(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Diluentname", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnDiluentvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Diluentvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnFinalConcentration(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("FinalConcentration", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnFinalVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("FinalVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("InPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeBuffer) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(11, e, e)
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

type MakeBuffer struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	Buffername         <-chan execute.ThreadParam
	Bufferstock        <-chan execute.ThreadParam
	Bufferstockconc    <-chan execute.ThreadParam
	Bufferstockvolume  <-chan execute.ThreadParam
	Diluent            <-chan execute.ThreadParam
	Diluentname        <-chan execute.ThreadParam
	Diluentvolume      <-chan execute.ThreadParam
	FinalConcentration <-chan execute.ThreadParam
	FinalVolume        <-chan execute.ThreadParam
	InPlate            <-chan execute.ThreadParam
	OutPlate           <-chan execute.ThreadParam
	Buffer             chan<- execute.ThreadParam
	Status             chan<- execute.ThreadParam
}

type MakeBufferParamBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	Buffername         string
	Bufferstock        *wtype.LHComponent
	Bufferstockconc    wunit.Concentration
	Bufferstockvolume  wunit.Volume
	Diluent            *wtype.LHComponent
	Diluentname        string
	Diluentvolume      wunit.Volume
	FinalConcentration wunit.Concentration
	FinalVolume        wunit.Volume
	InPlate            *wtype.LHPlate
	OutPlate           *wtype.LHPlate
}

type MakeBufferConfig struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	Buffername         string
	Bufferstock        wtype.FromFactory
	Bufferstockconc    wunit.Concentration
	Bufferstockvolume  wunit.Volume
	Diluent            wtype.FromFactory
	Diluentname        string
	Diluentvolume      wunit.Volume
	FinalConcentration wunit.Concentration
	FinalVolume        wunit.Volume
	InPlate            wtype.FromFactory
	OutPlate           wtype.FromFactory
}

type MakeBufferResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Buffer  *wtype.LHSolution
	Status  string
}

type MakeBufferJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	Buffername         *string
	Bufferstock        **wtype.LHComponent
	Bufferstockconc    *wunit.Concentration
	Bufferstockvolume  *wunit.Volume
	Diluent            **wtype.LHComponent
	Diluentname        *string
	Diluentvolume      *wunit.Volume
	FinalConcentration *wunit.Concentration
	FinalVolume        *wunit.Volume
	InPlate            **wtype.LHPlate
	OutPlate           **wtype.LHPlate
	Buffer             **wtype.LHSolution
	Status             *string
}

func (c *MakeBuffer) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Buffername", "string", "Buffername", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Bufferstock", "*wtype.LHComponent", "Bufferstock", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Bufferstockconc", "wunit.Concentration", "Bufferstockconc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Bufferstockvolume", "wunit.Volume", "Bufferstockvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Diluent", "*wtype.LHComponent", "Diluent", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Diluentname", "string", "Diluentname", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Diluentvolume", "wunit.Volume", "Diluentvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("FinalConcentration", "wunit.Concentration", "FinalConcentration", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("FinalVolume", "wunit.Volume", "FinalVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Buffer", "*wtype.LHSolution", "Buffer", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("MakeBuffer", "MakeBuffer", "", false, inp, outp)

	return ci
}
