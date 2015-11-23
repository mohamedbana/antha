package TypeIISConstructAssembly_sim

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
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

// Input Requirement specification
func (e *TypeIISConstructAssembly_sim) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *TypeIISConstructAssembly_sim) setup(p TypeIISConstructAssembly_simParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *TypeIISConstructAssembly_sim) steps(p TypeIISConstructAssembly_simParamBlock, r *TypeIISConstructAssembly_simResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	// Check that assembly is feasible by simulating assembly of the sequences with the chosen enzyme
	partsinorder := make([]wtype.DNASequence, 0)

	for _, part := range p.Partsinorder {
		partDNA := Inventory.Partslist[part]
		partsinorder = append(partsinorder, partDNA)
	}

	vectordata := Inventory.Partslist[p.Vectordata]
	assembly := enzymes.Assemblyparameters{p.Constructname, p.RestrictionEnzyme.CName, vectordata, partsinorder}
	status, numberofassemblies, sitesfound, newDNASequence, _ := enzymes.Assemblysimulator(assembly)

	r.NewDNASequence = newDNASequence
	r.Sitesfound = sitesfound

	if status == "Yay! this should work" && numberofassemblies == 1 {

		r.Simulationpass = true
	}
	// Monitor molar ratios of parts for possible troubleshooting / success correlation

	molesofeachdnaelement := make([]float64, 0)
	molarratios := make([]float64, 0)

	vector_mw := sequences.MassDNA(vectordata.Seq, false, true)
	vector_moles := sequences.Moles(p.VectorConcentration, vector_mw, p.VectorVol)
	molesofeachdnaelement = append(molesofeachdnaelement, vector_moles)

	molarratios = append(molarratios, (vector_moles / vector_moles))

	var part_mw float64
	var part_moles float64
	for i := 0; i < len(p.Partsinorder); i++ {

		part_mw = sequences.MassDNA(partsinorder[i].Seq, false, true)
		part_moles = sequences.Moles(p.PartConcs[i], part_mw, p.PartVols[i])

		molesofeachdnaelement = append(molesofeachdnaelement, part_moles)
		molarratios = append(molarratios, (part_moles / vector_moles))
	}

	r.Molesperpart = molesofeachdnaelement
	r.MolarratiotoVector = molarratios

	// Print status
	r.Status = fmt.Sprintln(
		"Simulationpass=", r.Simulationpass,
		"Molesperpart", r.Molesperpart,
		"MolarratiotoVector", r.MolarratiotoVector,
		"NewDNASequence", r.NewDNASequence,
		"Sitesfound", r.Sitesfound,
	)

	// Now Perform the physical assembly
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
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", p.Partsinorder[k], " vol ", p.PartVols[k])
		partSample := mixer.Sample(part, p.PartVols[k])
		partSample.CName = p.Partsinorder[k]
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
func (e *TypeIISConstructAssembly_sim) analysis(p TypeIISConstructAssembly_simParamBlock, r *TypeIISConstructAssembly_simResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *TypeIISConstructAssembly_sim) validation(p TypeIISConstructAssembly_simParamBlock, r *TypeIISConstructAssembly_simResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *TypeIISConstructAssembly_sim) Complete(params interface{}) {
	p := params.(TypeIISConstructAssembly_simParamBlock)
	if p.Error {
		e.MolarratiotoVector <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Molesperpart <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.NewDNASequence <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Simulationpass <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Sitesfound <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TypeIISConstructAssembly_simResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.MolarratiotoVector <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Molesperpart <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.NewDNASequence <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Reaction <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Simulationpass <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Sitesfound <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.MolarratiotoVector <- execute.ThreadParam{Value: r.MolarratiotoVector, ID: p.ID, Error: false}

	e.Molesperpart <- execute.ThreadParam{Value: r.Molesperpart, ID: p.ID, Error: false}

	e.NewDNASequence <- execute.ThreadParam{Value: r.NewDNASequence, ID: p.ID, Error: false}

	e.Reaction <- execute.ThreadParam{Value: r.Reaction, ID: p.ID, Error: false}

	e.Simulationpass <- execute.ThreadParam{Value: r.Simulationpass, ID: p.ID, Error: false}

	e.Sitesfound <- execute.ThreadParam{Value: r.Sitesfound, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *TypeIISConstructAssembly_sim) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *TypeIISConstructAssembly_sim) NewConfig() interface{} {
	return &TypeIISConstructAssembly_simConfig{}
}

func (e *TypeIISConstructAssembly_sim) NewParamBlock() interface{} {
	return &TypeIISConstructAssembly_simParamBlock{}
}

func NewTypeIISConstructAssembly_sim() interface{} { //*TypeIISConstructAssembly_sim {
	e := new(TypeIISConstructAssembly_sim)
	e.init()
	return e
}

// Mapper function
func (e *TypeIISConstructAssembly_sim) Map(m map[string]interface{}) interface{} {
	var res TypeIISConstructAssembly_simParamBlock
	res.Error = false || m["Atp"].(execute.ThreadParam).Error || m["AtpVol"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["BufferVol"].(execute.ThreadParam).Error || m["Constructname"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["LigVol"].(execute.ThreadParam).Error || m["Ligase"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["PartConcs"].(execute.ThreadParam).Error || m["PartVols"].(execute.ThreadParam).Error || m["Parts"].(execute.ThreadParam).Error || m["Partsinorder"].(execute.ThreadParam).Error || m["ReVol"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["RestrictionEnzyme"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error || m["VectorConcentration"].(execute.ThreadParam).Error || m["VectorVol"].(execute.ThreadParam).Error || m["Vectordata"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vAtp, is := m["Atp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vAtp.JSONString), &temp)
		res.Atp = *temp.Atp
	} else {
		res.Atp = m["Atp"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vAtpVol, is := m["AtpVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vAtpVol.JSONString), &temp)
		res.AtpVol = *temp.AtpVol
	} else {
		res.AtpVol = m["AtpVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferVol, is := m["BufferVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vBufferVol.JSONString), &temp)
		res.BufferVol = *temp.BufferVol
	} else {
		res.BufferVol = m["BufferVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vConstructname, is := m["Constructname"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vConstructname.JSONString), &temp)
		res.Constructname = *temp.Constructname
	} else {
		res.Constructname = m["Constructname"].(execute.ThreadParam).Value.(string)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInactivationTemp, is := m["InactivationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vInactivationTemp.JSONString), &temp)
		res.InactivationTemp = *temp.InactivationTemp
	} else {
		res.InactivationTemp = m["InactivationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vInactivationTime, is := m["InactivationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vInactivationTime.JSONString), &temp)
		res.InactivationTime = *temp.InactivationTime
	} else {
		res.InactivationTime = m["InactivationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vLigVol, is := m["LigVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vLigVol.JSONString), &temp)
		res.LigVol = *temp.LigVol
	} else {
		res.LigVol = m["LigVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vLigase, is := m["Ligase"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vLigase.JSONString), &temp)
		res.Ligase = *temp.Ligase
	} else {
		res.Ligase = m["Ligase"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPartConcs, is := m["PartConcs"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vPartConcs.JSONString), &temp)
		res.PartConcs = *temp.PartConcs
	} else {
		res.PartConcs = m["PartConcs"].(execute.ThreadParam).Value.([]wunit.Concentration)
	}

	vPartVols, is := m["PartVols"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vPartVols.JSONString), &temp)
		res.PartVols = *temp.PartVols
	} else {
		res.PartVols = m["PartVols"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vParts, is := m["Parts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vParts.JSONString), &temp)
		res.Parts = *temp.Parts
	} else {
		res.Parts = m["Parts"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vPartsinorder, is := m["Partsinorder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vPartsinorder.JSONString), &temp)
		res.Partsinorder = *temp.Partsinorder
	} else {
		res.Partsinorder = m["Partsinorder"].(execute.ThreadParam).Value.([]string)
	}

	vReVol, is := m["ReVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vReVol.JSONString), &temp)
		res.ReVol = *temp.ReVol
	} else {
		res.ReVol = m["ReVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vReactionTemp, is := m["ReactionTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vReactionTemp.JSONString), &temp)
		res.ReactionTemp = *temp.ReactionTemp
	} else {
		res.ReactionTemp = m["ReactionTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vReactionTime, is := m["ReactionTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vReactionTime.JSONString), &temp)
		res.ReactionTime = *temp.ReactionTime
	} else {
		res.ReactionTime = m["ReactionTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vRestrictionEnzyme, is := m["RestrictionEnzyme"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vRestrictionEnzyme.JSONString), &temp)
		res.RestrictionEnzyme = *temp.RestrictionEnzyme
	} else {
		res.RestrictionEnzyme = m["RestrictionEnzyme"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vVectorConcentration, is := m["VectorConcentration"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vVectorConcentration.JSONString), &temp)
		res.VectorConcentration = *temp.VectorConcentration
	} else {
		res.VectorConcentration = m["VectorConcentration"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vVectorVol, is := m["VectorVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vVectorVol.JSONString), &temp)
		res.VectorVol = *temp.VectorVol
	} else {
		res.VectorVol = m["VectorVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vVectordata, is := m["Vectordata"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vVectordata.JSONString), &temp)
		res.Vectordata = *temp.Vectordata
	} else {
		res.Vectordata = m["Vectordata"].(execute.ThreadParam).Value.(string)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssembly_simJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["Atp"].(execute.ThreadParam).ID
	res.BlockID = m["Atp"].(execute.ThreadParam).BlockID

	return res
}

/*
type Mole struct {
	number float64
}*/

func (e *TypeIISConstructAssembly_sim) OnAtp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnAtpVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnBufferVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnConstructname(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Constructname", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly_sim) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnInactivationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnInactivationTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnLigVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnLigase(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnPartConcs(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnPartVols(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnParts(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnPartsinorder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Partsinorder", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly_sim) OnReVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnReactionTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnReactionTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnRestrictionEnzyme(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnVector(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnVectorConcentration(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VectorConcentration", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly_sim) OnVectorVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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
func (e *TypeIISConstructAssembly_sim) OnVectordata(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vectordata", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISConstructAssembly_sim) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(25, e, e)
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

type TypeIISConstructAssembly_sim struct {
	flow.Component      // component "superclass" embedded
	lock                sync.Mutex
	startup             sync.Once
	params              map[execute.ThreadID]*execute.AsyncBag
	Atp                 <-chan execute.ThreadParam
	AtpVol              <-chan execute.ThreadParam
	Buffer              <-chan execute.ThreadParam
	BufferVol           <-chan execute.ThreadParam
	Constructname       <-chan execute.ThreadParam
	InPlate             <-chan execute.ThreadParam
	InactivationTemp    <-chan execute.ThreadParam
	InactivationTime    <-chan execute.ThreadParam
	LigVol              <-chan execute.ThreadParam
	Ligase              <-chan execute.ThreadParam
	OutPlate            <-chan execute.ThreadParam
	PartConcs           <-chan execute.ThreadParam
	PartVols            <-chan execute.ThreadParam
	Parts               <-chan execute.ThreadParam
	Partsinorder        <-chan execute.ThreadParam
	ReVol               <-chan execute.ThreadParam
	ReactionTemp        <-chan execute.ThreadParam
	ReactionTime        <-chan execute.ThreadParam
	ReactionVolume      <-chan execute.ThreadParam
	RestrictionEnzyme   <-chan execute.ThreadParam
	Vector              <-chan execute.ThreadParam
	VectorConcentration <-chan execute.ThreadParam
	VectorVol           <-chan execute.ThreadParam
	Vectordata          <-chan execute.ThreadParam
	Water               <-chan execute.ThreadParam
	MolarratiotoVector  chan<- execute.ThreadParam
	Molesperpart        chan<- execute.ThreadParam
	NewDNASequence      chan<- execute.ThreadParam
	Reaction            chan<- execute.ThreadParam
	Simulationpass      chan<- execute.ThreadParam
	Sitesfound          chan<- execute.ThreadParam
	Status              chan<- execute.ThreadParam
}

type TypeIISConstructAssembly_simParamBlock struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	Atp                 *wtype.LHComponent
	AtpVol              wunit.Volume
	Buffer              *wtype.LHComponent
	BufferVol           wunit.Volume
	Constructname       string
	InPlate             *wtype.LHPlate
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	LigVol              wunit.Volume
	Ligase              *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	PartConcs           []wunit.Concentration
	PartVols            []wunit.Volume
	Parts               []*wtype.LHComponent
	Partsinorder        []string
	ReVol               wunit.Volume
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	RestrictionEnzyme   *wtype.LHComponent
	Vector              *wtype.LHComponent
	VectorConcentration wunit.Concentration
	VectorVol           wunit.Volume
	Vectordata          string
	Water               *wtype.LHComponent
}

type TypeIISConstructAssembly_simConfig struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	Atp                 wtype.FromFactory
	AtpVol              wunit.Volume
	Buffer              wtype.FromFactory
	BufferVol           wunit.Volume
	Constructname       string
	InPlate             wtype.FromFactory
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	LigVol              wunit.Volume
	Ligase              wtype.FromFactory
	OutPlate            wtype.FromFactory
	PartConcs           []wunit.Concentration
	PartVols            []wunit.Volume
	Parts               []wtype.FromFactory
	Partsinorder        []string
	ReVol               wunit.Volume
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	RestrictionEnzyme   wtype.FromFactory
	Vector              wtype.FromFactory
	VectorConcentration wunit.Concentration
	VectorVol           wunit.Volume
	Vectordata          string
	Water               wtype.FromFactory
}

type TypeIISConstructAssembly_simResultBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	MolarratiotoVector []float64
	Molesperpart       []float64
	NewDNASequence     wtype.DNASequence
	Reaction           *wtype.LHSolution
	Simulationpass     bool
	Sitesfound         []enzymes.Restrictionsites
	Status             string
}

type TypeIISConstructAssembly_simJSONBlock struct {
	ID                  *execute.ThreadID
	BlockID             *execute.BlockID
	Error               *bool
	Atp                 **wtype.LHComponent
	AtpVol              *wunit.Volume
	Buffer              **wtype.LHComponent
	BufferVol           *wunit.Volume
	Constructname       *string
	InPlate             **wtype.LHPlate
	InactivationTemp    *wunit.Temperature
	InactivationTime    *wunit.Time
	LigVol              *wunit.Volume
	Ligase              **wtype.LHComponent
	OutPlate            **wtype.LHPlate
	PartConcs           *[]wunit.Concentration
	PartVols            *[]wunit.Volume
	Parts               *[]*wtype.LHComponent
	Partsinorder        *[]string
	ReVol               *wunit.Volume
	ReactionTemp        *wunit.Temperature
	ReactionTime        *wunit.Time
	ReactionVolume      *wunit.Volume
	RestrictionEnzyme   **wtype.LHComponent
	Vector              **wtype.LHComponent
	VectorConcentration *wunit.Concentration
	VectorVol           *wunit.Volume
	Vectordata          *string
	Water               **wtype.LHComponent
	MolarratiotoVector  *[]float64
	Molesperpart        *[]float64
	NewDNASequence      *wtype.DNASequence
	Reaction            **wtype.LHSolution
	Simulationpass      *bool
	Sitesfound          *[]enzymes.Restrictionsites
	Status              *string
}

func (c *TypeIISConstructAssembly_sim) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Atp", "*wtype.LHComponent", "Atp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("AtpVol", "wunit.Volume", "AtpVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferVol", "wunit.Volume", "BufferVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Constructname", "string", "Constructname", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTemp", "wunit.Temperature", "InactivationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTime", "wunit.Time", "InactivationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("LigVol", "wunit.Volume", "LigVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Ligase", "*wtype.LHComponent", "Ligase", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartConcs", "[]wunit.Concentration", "PartConcs", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartVols", "[]wunit.Volume", "PartVols", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Parts", "[]*wtype.LHComponent", "Parts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Partsinorder", "[]string", "Partsinorder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReVol", "wunit.Volume", "ReVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTemp", "wunit.Temperature", "ReactionTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTime", "wunit.Time", "ReactionTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RestrictionEnzyme", "*wtype.LHComponent", "RestrictionEnzyme", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "*wtype.LHComponent", "Vector", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VectorConcentration", "wunit.Concentration", "VectorConcentration", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VectorVol", "wunit.Volume", "VectorVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vectordata", "string", "Vectordata", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("MolarratiotoVector", "[]float64", "MolarratiotoVector", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Molesperpart", "[]float64", "Molesperpart", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("NewDNASequence", "wtype.DNASequence", "NewDNASequence", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Simulationpass", "bool", "Simulationpass", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Sitesfound", "[]enzymes.Restrictionsites", "Sitesfound", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("TypeIISConstructAssembly_sim", "TypeIISConstructAssembly_sim", "", false, inp, outp)

	return ci
}
