package Mastermix

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

//ComponentNames []string

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Mastermix) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Mastermix) setup(p MastermixParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Mastermix) steps(p MastermixParamBlock, r *MastermixResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	if len(p.OtherComponents) != len(p.OtherComponentVolumes) {
		panic("len(OtherComponents) != len(OtherComponentVolumes)")
	}

	mastermixes := make([]*wtype.LHSolution, 0)
	//var mastermix  *wtype.LHSolution

	if p.AliquotbyRow {
		panic("Add MixTo based method coming soon!")
	} else {
		for i := 0; i < p.NumberofMastermixes; i++ {

			eachmastermix := make([]*wtype.LHComponent, 0)

			if p.Buffer != nil {
				bufferSample := mixer.SampleForTotalVolume(p.Buffer, p.TotalVolumeperMastermix)
				eachmastermix = append(eachmastermix, bufferSample)
			}

			for k, component := range p.OtherComponents {
				if k == len(p.OtherComponents) {
					component.Type = "NeedToMix"
				}
				componentSample := mixer.Sample(component, p.OtherComponentVolumes[k])
				eachmastermix = append(eachmastermix, componentSample)
			}

			mastermix := _wrapper.MixInto(p.OutPlate, eachmastermix...)
			mastermixes = append(mastermixes, mastermix)

		}

	}
	r.Mastermixes = mastermixes
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Mastermix) analysis(p MastermixParamBlock, r *MastermixResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Mastermix) validation(p MastermixParamBlock, r *MastermixResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Mastermix) Complete(params interface{}) {
	p := params.(MastermixParamBlock)
	if p.Error {
		e.Mastermixes <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(MastermixResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Mastermixes <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Mastermixes <- execute.ThreadParam{Value: r.Mastermixes, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Mastermix) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Mastermix) NewConfig() interface{} {
	return &MastermixConfig{}
}

func (e *Mastermix) NewParamBlock() interface{} {
	return &MastermixParamBlock{}
}

func NewMastermix() interface{} { //*Mastermix {
	e := new(Mastermix)
	e.init()
	return e
}

// Mapper function
func (e *Mastermix) Map(m map[string]interface{}) interface{} {
	var res MastermixParamBlock
	res.Error = false || m["AliquotbyRow"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["Inplate"].(execute.ThreadParam).Error || m["NumberofMastermixes"].(execute.ThreadParam).Error || m["OtherComponentVolumes"].(execute.ThreadParam).Error || m["OtherComponents"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["TotalVolumeperMastermix"].(execute.ThreadParam).Error

	vAliquotbyRow, is := m["AliquotbyRow"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vAliquotbyRow.JSONString), &temp)
		res.AliquotbyRow = *temp.AliquotbyRow
	} else {
		res.AliquotbyRow = m["AliquotbyRow"].(execute.ThreadParam).Value.(bool)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vInplate, is := m["Inplate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vInplate.JSONString), &temp)
		res.Inplate = *temp.Inplate
	} else {
		res.Inplate = m["Inplate"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vNumberofMastermixes, is := m["NumberofMastermixes"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vNumberofMastermixes.JSONString), &temp)
		res.NumberofMastermixes = *temp.NumberofMastermixes
	} else {
		res.NumberofMastermixes = m["NumberofMastermixes"].(execute.ThreadParam).Value.(int)
	}

	vOtherComponentVolumes, is := m["OtherComponentVolumes"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vOtherComponentVolumes.JSONString), &temp)
		res.OtherComponentVolumes = *temp.OtherComponentVolumes
	} else {
		res.OtherComponentVolumes = m["OtherComponentVolumes"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vOtherComponents, is := m["OtherComponents"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vOtherComponents.JSONString), &temp)
		res.OtherComponents = *temp.OtherComponents
	} else {
		res.OtherComponents = m["OtherComponents"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vTotalVolumeperMastermix, is := m["TotalVolumeperMastermix"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MastermixJSONBlock
		json.Unmarshal([]byte(vTotalVolumeperMastermix.JSONString), &temp)
		res.TotalVolumeperMastermix = *temp.TotalVolumeperMastermix
	} else {
		res.TotalVolumeperMastermix = m["TotalVolumeperMastermix"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["AliquotbyRow"].(execute.ThreadParam).ID
	res.BlockID = m["AliquotbyRow"].(execute.ThreadParam).BlockID

	return res
}

func (e *Mastermix) OnAliquotbyRow(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AliquotbyRow", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix) OnBuffer(param execute.ThreadParam) {
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
func (e *Mastermix) OnInplate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Inplate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix) OnNumberofMastermixes(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("NumberofMastermixes", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix) OnOtherComponentVolumes(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OtherComponentVolumes", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix) OnOtherComponents(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OtherComponents", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix) OnOutPlate(param execute.ThreadParam) {
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
func (e *Mastermix) OnTotalVolumeperMastermix(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TotalVolumeperMastermix", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Mastermix struct {
	flow.Component          // component "superclass" embedded
	lock                    sync.Mutex
	startup                 sync.Once
	params                  map[execute.ThreadID]*execute.AsyncBag
	AliquotbyRow            <-chan execute.ThreadParam
	Buffer                  <-chan execute.ThreadParam
	Inplate                 <-chan execute.ThreadParam
	NumberofMastermixes     <-chan execute.ThreadParam
	OtherComponentVolumes   <-chan execute.ThreadParam
	OtherComponents         <-chan execute.ThreadParam
	OutPlate                <-chan execute.ThreadParam
	TotalVolumeperMastermix <-chan execute.ThreadParam
	Mastermixes             chan<- execute.ThreadParam
	Status                  chan<- execute.ThreadParam
}

type MastermixParamBlock struct {
	ID                      execute.ThreadID
	BlockID                 execute.BlockID
	Error                   bool
	AliquotbyRow            bool
	Buffer                  *wtype.LHComponent
	Inplate                 *wtype.LHComponent
	NumberofMastermixes     int
	OtherComponentVolumes   []wunit.Volume
	OtherComponents         []*wtype.LHComponent
	OutPlate                *wtype.LHPlate
	TotalVolumeperMastermix wunit.Volume
}

type MastermixConfig struct {
	ID                      execute.ThreadID
	BlockID                 execute.BlockID
	Error                   bool
	AliquotbyRow            bool
	Buffer                  wtype.FromFactory
	Inplate                 wtype.FromFactory
	NumberofMastermixes     int
	OtherComponentVolumes   []wunit.Volume
	OtherComponents         []wtype.FromFactory
	OutPlate                wtype.FromFactory
	TotalVolumeperMastermix wunit.Volume
}

type MastermixResultBlock struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	Mastermixes []*wtype.LHSolution
	Status      string
}

type MastermixJSONBlock struct {
	ID                      *execute.ThreadID
	BlockID                 *execute.BlockID
	Error                   *bool
	AliquotbyRow            *bool
	Buffer                  **wtype.LHComponent
	Inplate                 **wtype.LHComponent
	NumberofMastermixes     *int
	OtherComponentVolumes   *[]wunit.Volume
	OtherComponents         *[]*wtype.LHComponent
	OutPlate                **wtype.LHPlate
	TotalVolumeperMastermix *wunit.Volume
	Mastermixes             *[]*wtype.LHSolution
	Status                  *string
}

func (c *Mastermix) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AliquotbyRow", "bool", "AliquotbyRow", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Inplate", "*wtype.LHComponent", "Inplate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("NumberofMastermixes", "int", "NumberofMastermixes", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OtherComponentVolumes", "[]wunit.Volume", "OtherComponentVolumes", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OtherComponents", "[]*wtype.LHComponent", "OtherComponents", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TotalVolumeperMastermix", "wunit.Volume", "TotalVolumeperMastermix", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Mastermixes", "[]*wtype.LHSolution", "Mastermixes", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Mastermix", "Mastermix", "", false, inp, outp)

	return ci
}
