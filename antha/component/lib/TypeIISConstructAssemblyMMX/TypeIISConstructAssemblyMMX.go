package TypeIISConstructAssemblyMMX

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

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func (e *TypeIISConstructAssemblyMMX) requirements() { _ = wunit.Make_units }

// Conditions to run on startup
func (e *TypeIISConstructAssemblyMMX) setup(p TypeIISConstructAssemblyMMXParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *TypeIISConstructAssemblyMMX) steps(p TypeIISConstructAssemblyMMXParamBlock, r *TypeIISConstructAssemblyMMXResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	mmxSample := mixer.Sample(p.MasterMix, p.MMXVol)
	samples = append(samples, mmxSample)

	waterSample := mixer.SampleForTotalVolume(p.Water, p.ReactionVolume)
	samples = append(samples, waterSample)

	vectorSample := mixer.Sample(p.Vector, p.VectorVol)
	samples = append(samples, vectorSample)

	for k, part := range p.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", p.PartNames[k], " vol ", p.PartVols[k])
		partSample := mixer.Sample(part, p.PartVols[k])
		partSample.CName = p.PartNames[k]
		samples = append(samples, partSample)
	}

	r.Reaction = _wrapper.MixInto(p.OutPlate, samples...)

	// incubate the reaction mixture
	_wrapper.Incubate(r.Reaction, p.ReactionTemp, p.ReactionTime, false)
	// inactivate
	_wrapper.Incubate(r.Reaction, p.InactivationTemp, p.InactivationTime, false)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *TypeIISConstructAssemblyMMX) analysis(p TypeIISConstructAssemblyMMXParamBlock, r *TypeIISConstructAssemblyMMXResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *TypeIISConstructAssemblyMMX) validation(p TypeIISConstructAssemblyMMXParamBlock, r *TypeIISConstructAssemblyMMXResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *TypeIISConstructAssemblyMMX) Complete(params interface{}) {
	p := params.(TypeIISConstructAssemblyMMXParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TypeIISConstructAssemblyMMXResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Reaction <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Reaction <- execute.ThreadParam{Value: r.Reaction, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *TypeIISConstructAssemblyMMX) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *TypeIISConstructAssemblyMMX) NewConfig() interface{} {
	return &TypeIISConstructAssemblyMMXConfig{}
}

func (e *TypeIISConstructAssemblyMMX) NewParamBlock() interface{} {
	return &TypeIISConstructAssemblyMMXParamBlock{}
}

func NewTypeIISConstructAssemblyMMX() interface{} { //*TypeIISConstructAssemblyMMX {
	e := new(TypeIISConstructAssemblyMMX)
	e.init()
	return e
}

// Mapper function
func (e *TypeIISConstructAssemblyMMX) Map(m map[string]interface{}) interface{} {
	var res TypeIISConstructAssemblyMMXParamBlock
	res.Error = false || m["InPlate"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["MMXVol"].(execute.ThreadParam).Error || m["MasterMix"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["OutputReactionName"].(execute.ThreadParam).Error || m["PartNames"].(execute.ThreadParam).Error || m["PartVols"].(execute.ThreadParam).Error || m["Parts"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error || m["VectorVol"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInactivationTemp, is := m["InactivationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vInactivationTemp.JSONString), &temp)
		res.InactivationTemp = *temp.InactivationTemp
	} else {
		res.InactivationTemp = m["InactivationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vInactivationTime, is := m["InactivationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vInactivationTime.JSONString), &temp)
		res.InactivationTime = *temp.InactivationTime
	} else {
		res.InactivationTime = m["InactivationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vMMXVol, is := m["MMXVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vMMXVol.JSONString), &temp)
		res.MMXVol = *temp.MMXVol
	} else {
		res.MMXVol = m["MMXVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vMasterMix, is := m["MasterMix"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vMasterMix.JSONString), &temp)
		res.MasterMix = *temp.MasterMix
	} else {
		res.MasterMix = m["MasterMix"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vOutputReactionName, is := m["OutputReactionName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vOutputReactionName.JSONString), &temp)
		res.OutputReactionName = *temp.OutputReactionName
	} else {
		res.OutputReactionName = m["OutputReactionName"].(execute.ThreadParam).Value.(string)
	}

	vPartNames, is := m["PartNames"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vPartNames.JSONString), &temp)
		res.PartNames = *temp.PartNames
	} else {
		res.PartNames = m["PartNames"].(execute.ThreadParam).Value.([]string)
	}

	vPartVols, is := m["PartVols"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vPartVols.JSONString), &temp)
		res.PartVols = *temp.PartVols
	} else {
		res.PartVols = m["PartVols"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vParts, is := m["Parts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vParts.JSONString), &temp)
		res.Parts = *temp.Parts
	} else {
		res.Parts = m["Parts"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vReactionTemp, is := m["ReactionTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vReactionTemp.JSONString), &temp)
		res.ReactionTemp = *temp.ReactionTemp
	} else {
		res.ReactionTemp = m["ReactionTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vReactionTime, is := m["ReactionTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vReactionTime.JSONString), &temp)
		res.ReactionTime = *temp.ReactionTime
	} else {
		res.ReactionTime = m["ReactionTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVectorVol, is := m["VectorVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vVectorVol.JSONString), &temp)
		res.VectorVol = *temp.VectorVol
	} else {
		res.VectorVol = m["VectorVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyMMXJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["InPlate"].(execute.ThreadParam).ID
	res.BlockID = m["InPlate"].(execute.ThreadParam).BlockID

	return res
}

func (e *TypeIISConstructAssemblyMMX) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
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
func (e *TypeIISConstructAssemblyMMX) OnInactivationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("InactivationTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnInactivationTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("InactivationTime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnMMXVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("MMXVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnMasterMix(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("MasterMix", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
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
func (e *TypeIISConstructAssemblyMMX) OnOutputReactionName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OutputReactionName", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnPartNames(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PartNames", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnPartVols(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PartVols", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnParts(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Parts", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnReactionTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ReactionTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnReactionTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ReactionTime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ReactionVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnVector(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vector", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnVectorVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VectorVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssemblyMMX) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(16, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Water", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type TypeIISConstructAssemblyMMX struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	InPlate            <-chan execute.ThreadParam
	InactivationTemp   <-chan execute.ThreadParam
	InactivationTime   <-chan execute.ThreadParam
	MMXVol             <-chan execute.ThreadParam
	MasterMix          <-chan execute.ThreadParam
	OutPlate           <-chan execute.ThreadParam
	OutputReactionName <-chan execute.ThreadParam
	PartNames          <-chan execute.ThreadParam
	PartVols           <-chan execute.ThreadParam
	Parts              <-chan execute.ThreadParam
	ReactionTemp       <-chan execute.ThreadParam
	ReactionTime       <-chan execute.ThreadParam
	ReactionVolume     <-chan execute.ThreadParam
	Vector             <-chan execute.ThreadParam
	VectorVol          <-chan execute.ThreadParam
	Water              <-chan execute.ThreadParam
	Reaction           chan<- execute.ThreadParam
}

type TypeIISConstructAssemblyMMXParamBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	InPlate            *wtype.LHPlate
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	MMXVol             wunit.Volume
	MasterMix          *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	Vector             *wtype.LHComponent
	VectorVol          wunit.Volume
	Water              *wtype.LHComponent
}

type TypeIISConstructAssemblyMMXConfig struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	InPlate            wtype.FromFactory
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	MMXVol             wunit.Volume
	MasterMix          wtype.FromFactory
	OutPlate           wtype.FromFactory
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []wtype.FromFactory
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	Vector             wtype.FromFactory
	VectorVol          wunit.Volume
	Water              wtype.FromFactory
}

type TypeIISConstructAssemblyMMXResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Reaction *wtype.LHSolution
}

type TypeIISConstructAssemblyMMXJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	InPlate            **wtype.LHPlate
	InactivationTemp   *wunit.Temperature
	InactivationTime   *wunit.Time
	MMXVol             *wunit.Volume
	MasterMix          **wtype.LHComponent
	OutPlate           **wtype.LHPlate
	OutputReactionName *string
	PartNames          *[]string
	PartVols           *[]wunit.Volume
	Parts              *[]*wtype.LHComponent
	ReactionTemp       *wunit.Temperature
	ReactionTime       *wunit.Time
	ReactionVolume     *wunit.Volume
	Vector             **wtype.LHComponent
	VectorVol          *wunit.Volume
	Water              **wtype.LHComponent
	Reaction           **wtype.LHSolution
}

func (c *TypeIISConstructAssemblyMMX) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTemp", "wunit.Temperature", "InactivationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTime", "wunit.Time", "InactivationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("MMXVol", "wunit.Volume", "MMXVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("MasterMix", "*wtype.LHComponent", "MasterMix", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutputReactionName", "string", "OutputReactionName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartNames", "[]string", "PartNames", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartVols", "[]wunit.Volume", "PartVols", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Parts", "[]*wtype.LHComponent", "Parts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTemp", "wunit.Temperature", "ReactionTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTime", "wunit.Time", "ReactionTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "*wtype.LHComponent", "Vector", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VectorVol", "wunit.Volume", "VectorVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))

	ci := execute.NewComponentInfo("TypeIISConstructAssemblyMMX", "TypeIISConstructAssemblyMMX", "", false, inp, outp)

	return ci
}
