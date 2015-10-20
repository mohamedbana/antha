// antha/component/lib/TypeIISConstructAssembly_alt/TypeIISConstructAssembly_alt.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package TypeIISConstructAssembly_alt

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"strings"
	"sync"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func (e *TypeIISConstructAssembly_alt) requirements() { _ = wunit.Make_units }

// Conditions to run on startup
func (e *TypeIISConstructAssembly_alt) setup(p TypeIISConstructAssembly_altParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *TypeIISConstructAssembly_alt) steps(p TypeIISConstructAssembly_altParamBlock, r *TypeIISConstructAssembly_altResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
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

	s := ""
	comments := make([]string, 0)
	var partSample *wtype.LHComponent

	for k, part := range p.Parts {
		if p.PartConcs[k].SIValue() <= 0.1 {
			s = fmt.Sprintln("creating dna part num ", k, " comp ", part.CName, " renamed to ", p.PartNames[k], " vol ", p.PartConcs[k].ToString())
			partSample = mixer.SampleForConcentration(part, p.PartConcs[k])
		} else {
			s = fmt.Sprintln("Conc too low so minimum volume used", "creating dna part num ", k, " comp ", part.CName, " renamed to ", p.PartNames[k], " vol ", p.PartMinVol.ToString())
			partSample = mixer.Sample(part, p.PartMinVol)
		}
		partSample.CName = p.PartNames[k]
		samples = append(samples, partSample)
		comments = append(comments, s)

	}
	r.S = strings.Join(comments, "")

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
func (e *TypeIISConstructAssembly_alt) analysis(p TypeIISConstructAssembly_altParamBlock, r *TypeIISConstructAssembly_altResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *TypeIISConstructAssembly_alt) validation(p TypeIISConstructAssembly_altParamBlock, r *TypeIISConstructAssembly_altResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *TypeIISConstructAssembly_alt) Complete(params interface{}) {
	p := params.(TypeIISConstructAssembly_altParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.S <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TypeIISConstructAssembly_altResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Reaction <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.S <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Reaction <- execute.ThreadParam{Value: r.Reaction, ID: p.ID, Error: false}

	e.S <- execute.ThreadParam{Value: r.S, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *TypeIISConstructAssembly_alt) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *TypeIISConstructAssembly_alt) NewConfig() interface{} {
	return &TypeIISConstructAssembly_altConfig{}
}

func (e *TypeIISConstructAssembly_alt) NewParamBlock() interface{} {
	return &TypeIISConstructAssembly_altParamBlock{}
}

func NewTypeIISConstructAssembly_alt() interface{} { //*TypeIISConstructAssembly_alt {
	e := new(TypeIISConstructAssembly_alt)
	e.init()
	return e
}

// Mapper function
func (e *TypeIISConstructAssembly_alt) Map(m map[string]interface{}) interface{} {
	var res TypeIISConstructAssembly_altParamBlock
	res.Error = false || m["Atp"].(execute.ThreadParam).Error || m["AtpVol"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["BufferVol"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["LigVol"].(execute.ThreadParam).Error || m["Ligase"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["PartConcs"].(execute.ThreadParam).Error || m["PartMinVol"].(execute.ThreadParam).Error || m["PartNames"].(execute.ThreadParam).Error || m["Parts"].(execute.ThreadParam).Error || m["ReVol"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["RestrictionEnzyme"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error || m["VectorVol"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vAtp, is := m["Atp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vAtp.JSONString), &temp)
		res.Atp = *temp.Atp
	} else {
		res.Atp = m["Atp"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vAtpVol, is := m["AtpVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vAtpVol.JSONString), &temp)
		res.AtpVol = *temp.AtpVol
	} else {
		res.AtpVol = m["AtpVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferVol, is := m["BufferVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vBufferVol.JSONString), &temp)
		res.BufferVol = *temp.BufferVol
	} else {
		res.BufferVol = m["BufferVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInactivationTemp, is := m["InactivationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vInactivationTemp.JSONString), &temp)
		res.InactivationTemp = *temp.InactivationTemp
	} else {
		res.InactivationTemp = m["InactivationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vInactivationTime, is := m["InactivationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vInactivationTime.JSONString), &temp)
		res.InactivationTime = *temp.InactivationTime
	} else {
		res.InactivationTime = m["InactivationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vLigVol, is := m["LigVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vLigVol.JSONString), &temp)
		res.LigVol = *temp.LigVol
	} else {
		res.LigVol = m["LigVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vLigase, is := m["Ligase"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vLigase.JSONString), &temp)
		res.Ligase = *temp.Ligase
	} else {
		res.Ligase = m["Ligase"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPartConcs, is := m["PartConcs"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vPartConcs.JSONString), &temp)
		res.PartConcs = *temp.PartConcs
	} else {
		res.PartConcs = m["PartConcs"].(execute.ThreadParam).Value.([]wunit.Concentration)
	}

	vPartMinVol, is := m["PartMinVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vPartMinVol.JSONString), &temp)
		res.PartMinVol = *temp.PartMinVol
	} else {
		res.PartMinVol = m["PartMinVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vPartNames, is := m["PartNames"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vPartNames.JSONString), &temp)
		res.PartNames = *temp.PartNames
	} else {
		res.PartNames = m["PartNames"].(execute.ThreadParam).Value.([]string)
	}

	vParts, is := m["Parts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vParts.JSONString), &temp)
		res.Parts = *temp.Parts
	} else {
		res.Parts = m["Parts"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vReVol, is := m["ReVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vReVol.JSONString), &temp)
		res.ReVol = *temp.ReVol
	} else {
		res.ReVol = m["ReVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vReactionTemp, is := m["ReactionTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vReactionTemp.JSONString), &temp)
		res.ReactionTemp = *temp.ReactionTemp
	} else {
		res.ReactionTemp = m["ReactionTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vReactionTime, is := m["ReactionTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vReactionTime.JSONString), &temp)
		res.ReactionTime = *temp.ReactionTime
	} else {
		res.ReactionTime = m["ReactionTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vRestrictionEnzyme, is := m["RestrictionEnzyme"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vRestrictionEnzyme.JSONString), &temp)
		res.RestrictionEnzyme = *temp.RestrictionEnzyme
	} else {
		res.RestrictionEnzyme = m["RestrictionEnzyme"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVectorVol, is := m["VectorVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vVectorVol.JSONString), &temp)
		res.VectorVol = *temp.VectorVol
	} else {
		res.VectorVol = m["VectorVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_altJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["Atp"].(execute.ThreadParam).ID
	res.BlockID = m["Atp"].(execute.ThreadParam).BlockID

	return res
}

func (e *TypeIISConstructAssembly_alt) OnAtp(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnAtpVol(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnBuffer(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnBufferVol(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnInPlate(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnInactivationTemp(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnInactivationTime(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnLigVol(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnLigase(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnOutPlate(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnPartConcs(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PartConcs", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly_alt) OnPartMinVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(22, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PartMinVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly_alt) OnPartNames(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnParts(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnReVol(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnReactionTemp(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnReactionTime(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnReactionVolume(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnRestrictionEnzyme(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnVector(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnVectorVol(param execute.ThreadParam) {
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
func (e *TypeIISConstructAssembly_alt) OnWater(param execute.ThreadParam) {
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

type TypeIISConstructAssembly_alt struct {
	flow.Component    // component "superclass" embedded
	lock              sync.Mutex
	startup           sync.Once
	params            map[execute.ThreadID]*execute.AsyncBag
	Atp               <-chan execute.ThreadParam
	AtpVol            <-chan execute.ThreadParam
	Buffer            <-chan execute.ThreadParam
	BufferVol         <-chan execute.ThreadParam
	InPlate           <-chan execute.ThreadParam
	InactivationTemp  <-chan execute.ThreadParam
	InactivationTime  <-chan execute.ThreadParam
	LigVol            <-chan execute.ThreadParam
	Ligase            <-chan execute.ThreadParam
	OutPlate          <-chan execute.ThreadParam
	PartConcs         <-chan execute.ThreadParam
	PartMinVol        <-chan execute.ThreadParam
	PartNames         <-chan execute.ThreadParam
	Parts             <-chan execute.ThreadParam
	ReVol             <-chan execute.ThreadParam
	ReactionTemp      <-chan execute.ThreadParam
	ReactionTime      <-chan execute.ThreadParam
	ReactionVolume    <-chan execute.ThreadParam
	RestrictionEnzyme <-chan execute.ThreadParam
	Vector            <-chan execute.ThreadParam
	VectorVol         <-chan execute.ThreadParam
	Water             <-chan execute.ThreadParam
	Reaction          chan<- execute.ThreadParam
	S                 chan<- execute.ThreadParam
}

type TypeIISConstructAssembly_altParamBlock struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Atp               *wtype.LHComponent
	AtpVol            wunit.Volume
	Buffer            *wtype.LHComponent
	BufferVol         wunit.Volume
	InPlate           *wtype.LHPlate
	InactivationTemp  wunit.Temperature
	InactivationTime  wunit.Time
	LigVol            wunit.Volume
	Ligase            *wtype.LHComponent
	OutPlate          *wtype.LHPlate
	PartConcs         []wunit.Concentration
	PartMinVol        wunit.Volume
	PartNames         []string
	Parts             []*wtype.LHComponent
	ReVol             wunit.Volume
	ReactionTemp      wunit.Temperature
	ReactionTime      wunit.Time
	ReactionVolume    wunit.Volume
	RestrictionEnzyme *wtype.LHComponent
	Vector            *wtype.LHComponent
	VectorVol         wunit.Volume
	Water             *wtype.LHComponent
}

type TypeIISConstructAssembly_altConfig struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Atp               wtype.FromFactory
	AtpVol            wunit.Volume
	Buffer            wtype.FromFactory
	BufferVol         wunit.Volume
	InPlate           wtype.FromFactory
	InactivationTemp  wunit.Temperature
	InactivationTime  wunit.Time
	LigVol            wunit.Volume
	Ligase            wtype.FromFactory
	OutPlate          wtype.FromFactory
	PartConcs         []wunit.Concentration
	PartMinVol        wunit.Volume
	PartNames         []string
	Parts             []wtype.FromFactory
	ReVol             wunit.Volume
	ReactionTemp      wunit.Temperature
	ReactionTime      wunit.Time
	ReactionVolume    wunit.Volume
	RestrictionEnzyme wtype.FromFactory
	Vector            wtype.FromFactory
	VectorVol         wunit.Volume
	Water             wtype.FromFactory
}

type TypeIISConstructAssembly_altResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Reaction *wtype.LHSolution
	S        string
}

type TypeIISConstructAssembly_altJSONBlock struct {
	ID                *execute.ThreadID
	BlockID           *execute.BlockID
	Error             *bool
	Atp               **wtype.LHComponent
	AtpVol            *wunit.Volume
	Buffer            **wtype.LHComponent
	BufferVol         *wunit.Volume
	InPlate           **wtype.LHPlate
	InactivationTemp  *wunit.Temperature
	InactivationTime  *wunit.Time
	LigVol            *wunit.Volume
	Ligase            **wtype.LHComponent
	OutPlate          **wtype.LHPlate
	PartConcs         *[]wunit.Concentration
	PartMinVol        *wunit.Volume
	PartNames         *[]string
	Parts             *[]*wtype.LHComponent
	ReVol             *wunit.Volume
	ReactionTemp      *wunit.Temperature
	ReactionTime      *wunit.Time
	ReactionVolume    *wunit.Volume
	RestrictionEnzyme **wtype.LHComponent
	Vector            **wtype.LHComponent
	VectorVol         *wunit.Volume
	Water             **wtype.LHComponent
	Reaction          **wtype.LHSolution
	S                 *string
}

func (c *TypeIISConstructAssembly_alt) ComponentInfo() *execute.ComponentInfo {
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
	inp = append(inp, *execute.NewPortInfo("PartConcs", "[]wunit.Concentration", "PartConcs", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartMinVol", "wunit.Volume", "PartMinVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartNames", "[]string", "PartNames", true, true, nil, nil))
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
	outp = append(outp, *execute.NewPortInfo("S", "string", "S", true, true, nil, nil))

	ci := execute.NewComponentInfo("TypeIISConstructAssembly_alt", "TypeIISConstructAssembly_alt", "", false, inp, outp)

	return ci
}
