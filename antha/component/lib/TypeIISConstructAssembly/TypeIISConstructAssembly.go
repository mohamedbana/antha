package TypeIISConstructAssembly

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func (e *TypeIISConstructAssembly) requirements() { _ = wunit.Make_units }

// Conditions to run on startup
func (e *TypeIISConstructAssembly) setup(p TypeIISConstructAssemblyParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *TypeIISConstructAssembly) steps(p TypeIISConstructAssemblyParamBlock, r *TypeIISConstructAssemblyResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(p.Water, p.ReactionVolume)
	samples = append(samples, waterSample)

	bufferSample := mixer.Sample(p.Buffer, p.BufferVol)
	samples = append(samples, bufferSample)

	atpSample := mixer.Sample(p.Atp, p.AtpVol)
	samples = append(samples, atpSample)

	//vectorSample := mixer.Sample(Vector, VectorVol)
	vectorSample := mixer.Sample(p.Vector, p.VectorVol)
	samples = append(samples, vectorSample)

	for k, part := range p.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", p.PartNames[k], " vol ", p.PartVols[k])
		partSample := mixer.Sample(part, p.PartVols[k])
		partSample.CName = p.PartNames[k]
		samples = append(samples, partSample)
	}

	reSample := mixer.Sample(p.RestrictionEnzyme, p.ReVol)
	samples = append(samples, reSample)

	ligSample := mixer.Sample(p.Ligase, p.LigVol)
	samples = append(samples, ligSample)

	r.Reaction = _wrapper.MixInto(p.OutPlate, samples...)

	// incubate the reaction mixture
	_wrapper.Incubate(r.Reaction, p.ReactionTemp, p.ReactionTime, false)
	// inactivate
	_wrapper.Incubate(r.Reaction, p.InactivationTemp, p.InactivationTime, false)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *TypeIISConstructAssembly) analysis(p TypeIISConstructAssemblyParamBlock, r *TypeIISConstructAssemblyResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *TypeIISConstructAssembly) validation(p TypeIISConstructAssemblyParamBlock, r *TypeIISConstructAssemblyResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *TypeIISConstructAssembly) Complete(params interface{}) {
	p := params.(TypeIISConstructAssemblyParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TypeIISConstructAssemblyResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Reaction <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
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
func (e *TypeIISConstructAssembly) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *TypeIISConstructAssembly) NewConfig() interface{} {
	return &TypeIISConstructAssemblyConfig{}
}

func (e *TypeIISConstructAssembly) NewParamBlock() interface{} {
	return &TypeIISConstructAssemblyParamBlock{}
}

func NewTypeIISConstructAssembly() interface{} { //*TypeIISConstructAssembly {
	e := new(TypeIISConstructAssembly)
	e.init()
	return e
}

// Mapper function
func (e *TypeIISConstructAssembly) Map(m map[string]interface{}) interface{} {
	var res TypeIISConstructAssemblyParamBlock
	res.Error = false || m["Atp"].(execute.ThreadParam).Error || m["AtpVol"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["BufferVol"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["LigVol"].(execute.ThreadParam).Error || m["Ligase"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["OutputReactionName"].(execute.ThreadParam).Error || m["PartNames"].(execute.ThreadParam).Error || m["PartVols"].(execute.ThreadParam).Error || m["Parts"].(execute.ThreadParam).Error || m["ReVol"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["RestrictionEnzyme"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error || m["VectorVol"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vAtp, is := m["Atp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vAtp.JSONString), &temp)
		res.Atp = *temp.Atp
	} else {
		res.Atp = m["Atp"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vAtpVol, is := m["AtpVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vAtpVol.JSONString), &temp)
		res.AtpVol = *temp.AtpVol
	} else {
		res.AtpVol = m["AtpVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferVol, is := m["BufferVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vBufferVol.JSONString), &temp)
		res.BufferVol = *temp.BufferVol
	} else {
		res.BufferVol = m["BufferVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInactivationTemp, is := m["InactivationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vInactivationTemp.JSONString), &temp)
		res.InactivationTemp = *temp.InactivationTemp
	} else {
		res.InactivationTemp = m["InactivationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vInactivationTime, is := m["InactivationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vInactivationTime.JSONString), &temp)
		res.InactivationTime = *temp.InactivationTime
	} else {
		res.InactivationTime = m["InactivationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vLigVol, is := m["LigVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vLigVol.JSONString), &temp)
		res.LigVol = *temp.LigVol
	} else {
		res.LigVol = m["LigVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vLigase, is := m["Ligase"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vLigase.JSONString), &temp)
		res.Ligase = *temp.Ligase
	} else {
		res.Ligase = m["Ligase"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vOutputReactionName, is := m["OutputReactionName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vOutputReactionName.JSONString), &temp)
		res.OutputReactionName = *temp.OutputReactionName
	} else {
		res.OutputReactionName = m["OutputReactionName"].(execute.ThreadParam).Value.(string)
	}

	vPartNames, is := m["PartNames"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vPartNames.JSONString), &temp)
		res.PartNames = *temp.PartNames
	} else {
		res.PartNames = m["PartNames"].(execute.ThreadParam).Value.([]string)
	}

	vPartVols, is := m["PartVols"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vPartVols.JSONString), &temp)
		res.PartVols = *temp.PartVols
	} else {
		res.PartVols = m["PartVols"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vParts, is := m["Parts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vParts.JSONString), &temp)
		res.Parts = *temp.Parts
	} else {
		res.Parts = m["Parts"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vReVol, is := m["ReVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vReVol.JSONString), &temp)
		res.ReVol = *temp.ReVol
	} else {
		res.ReVol = m["ReVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vReactionTemp, is := m["ReactionTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vReactionTemp.JSONString), &temp)
		res.ReactionTemp = *temp.ReactionTemp
	} else {
		res.ReactionTemp = m["ReactionTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vReactionTime, is := m["ReactionTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vReactionTime.JSONString), &temp)
		res.ReactionTime = *temp.ReactionTime
	} else {
		res.ReactionTime = m["ReactionTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vRestrictionEnzyme, is := m["RestrictionEnzyme"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vRestrictionEnzyme.JSONString), &temp)
		res.RestrictionEnzyme = *temp.RestrictionEnzyme
	} else {
		res.RestrictionEnzyme = m["RestrictionEnzyme"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVectorVol, is := m["VectorVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vVectorVol.JSONString), &temp)
		res.VectorVol = *temp.VectorVol
	} else {
		res.VectorVol = m["VectorVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["Atp"].(execute.ThreadParam).ID
	res.BlockID = m["Atp"].(execute.ThreadParam).BlockID

	return res
}

func (e *TypeIISConstructAssembly) OnAtp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Atp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnAtpVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AtpVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnBufferVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BufferVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnInactivationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnInactivationTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnLigVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("LigVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnLigase(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Ligase", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnOutputReactionName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnPartNames(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnPartVols(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnParts(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnReVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ReVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnReactionTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnReactionTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnRestrictionEnzyme(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RestrictionEnzyme", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly) OnVector(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnVectorVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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
func (e *TypeIISConstructAssembly) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
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

type TypeIISConstructAssembly struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	Atp                <-chan execute.ThreadParam
	AtpVol             <-chan execute.ThreadParam
	Buffer             <-chan execute.ThreadParam
	BufferVol          <-chan execute.ThreadParam
	InPlate            <-chan execute.ThreadParam
	InactivationTemp   <-chan execute.ThreadParam
	InactivationTime   <-chan execute.ThreadParam
	LigVol             <-chan execute.ThreadParam
	Ligase             <-chan execute.ThreadParam
	OutPlate           <-chan execute.ThreadParam
	OutputReactionName <-chan execute.ThreadParam
	PartNames          <-chan execute.ThreadParam
	PartVols           <-chan execute.ThreadParam
	Parts              <-chan execute.ThreadParam
	ReVol              <-chan execute.ThreadParam
	ReactionTemp       <-chan execute.ThreadParam
	ReactionTime       <-chan execute.ThreadParam
	ReactionVolume     <-chan execute.ThreadParam
	RestrictionEnzyme  <-chan execute.ThreadParam
	Vector             <-chan execute.ThreadParam
	VectorVol          <-chan execute.ThreadParam
	Water              <-chan execute.ThreadParam
	Reaction           chan<- execute.ThreadParam
}

type TypeIISConstructAssemblyParamBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	Atp                *wtype.LHComponent
	AtpVol             wunit.Volume
	Buffer             *wtype.LHComponent
	BufferVol          wunit.Volume
	InPlate            *wtype.LHPlate
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LigVol             wunit.Volume
	Ligase             *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReVol              wunit.Volume
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	RestrictionEnzyme  *wtype.LHComponent
	Vector             *wtype.LHComponent
	VectorVol          wunit.Volume
	Water              *wtype.LHComponent
}

type TypeIISConstructAssemblyConfig struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	Atp                wtype.FromFactory
	AtpVol             wunit.Volume
	Buffer             wtype.FromFactory
	BufferVol          wunit.Volume
	InPlate            wtype.FromFactory
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LigVol             wunit.Volume
	Ligase             wtype.FromFactory
	OutPlate           wtype.FromFactory
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []wtype.FromFactory
	ReVol              wunit.Volume
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	RestrictionEnzyme  wtype.FromFactory
	Vector             wtype.FromFactory
	VectorVol          wunit.Volume
	Water              wtype.FromFactory
}

type TypeIISConstructAssemblyResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Reaction *wtype.LHSolution
}

type TypeIISConstructAssemblyJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	Atp                **wtype.LHComponent
	AtpVol             *wunit.Volume
	Buffer             **wtype.LHComponent
	BufferVol          *wunit.Volume
	InPlate            **wtype.LHPlate
	InactivationTemp   *wunit.Temperature
	InactivationTime   *wunit.Time
	LigVol             *wunit.Volume
	Ligase             **wtype.LHComponent
	OutPlate           **wtype.LHPlate
	OutputReactionName *string
	PartNames          *[]string
	PartVols           *[]wunit.Volume
	Parts              *[]*wtype.LHComponent
	ReVol              *wunit.Volume
	ReactionTemp       *wunit.Temperature
	ReactionTime       *wunit.Time
	ReactionVolume     *wunit.Volume
	RestrictionEnzyme  **wtype.LHComponent
	Vector             **wtype.LHComponent
	VectorVol          *wunit.Volume
	Water              **wtype.LHComponent
	Reaction           **wtype.LHSolution
}

func (c *TypeIISConstructAssembly) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Atp", "*wtype.LHComponent", "Atp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("AtpVol", "wunit.Volume", "AtpVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferVol", "wunit.Volume", "BufferVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTemp", "wunit.Temperature", "InactivationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTime", "wunit.Time", "InactivationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("LigVol", "wunit.Volume", "LigVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Ligase", "*wtype.LHComponent", "Ligase", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutputReactionName", "string", "OutputReactionName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartNames", "[]string", "PartNames", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartVols", "[]wunit.Volume", "PartVols", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Parts", "[]*wtype.LHComponent", "Parts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReVol", "wunit.Volume", "ReVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTemp", "wunit.Temperature", "ReactionTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTime", "wunit.Time", "ReactionTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RestrictionEnzyme", "*wtype.LHComponent", "RestrictionEnzyme", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "*wtype.LHComponent", "Vector", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VectorVol", "wunit.Volume", "VectorVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))

	ci := execute.NewComponentInfo("TypeIISConstructAssembly", "TypeIISConstructAssembly", "", false, inp, outp)

	return ci
}
