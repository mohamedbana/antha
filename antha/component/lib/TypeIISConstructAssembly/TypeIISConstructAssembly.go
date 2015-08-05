package TypeIISConstructAssembly

import (
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/execution"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"sync"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *TypeIISConstructAssembly) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *TypeIISConstructAssembly) setup(p TypeIISConstructAssemblyParamBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *TypeIISConstructAssembly) steps(p TypeIISConstructAssemblyParamBlock, r *TypeIISConstructAssemblyResultBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.SampleForTotalVolume(p.Buffer, p.ReactionVolume)
	samples = append(samples, bufferSample)
	atpSample := mixer.Sample(p.Atp, p.AtpVol)
	samples = append(samples, atpSample)
	vectorSample := mixer.SampleForConcentration(p.Vector, p.VectorConc)
	samples = append(samples, vectorSample)

	for _, part := range p.Parts {
		partSample := mixer.SampleForConcentration(part, p.PartConc)
		samples = append(samples, partSample)
	}

	reSample := mixer.Sample(p.RestrictionEnzyme, p.ReVol)
	samples = append(samples, reSample)
	ligSample := mixer.Sample(p.Ligase, p.LigVol)
	samples = append(samples, ligSample)
	reaction := _wrapper.MixInto(p.OutPlate, samples...)

	// incubate the reaction mixture

	_wrapper.Incubate(reaction, p.ReactionTemp, p.ReactionTime, false)

	// inactivate

	_wrapper.Incubate(reaction, p.InactivationTemp, p.InactivationTime, false)

	// all done
	r.Reaction = reaction
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *TypeIISConstructAssembly) analysis(p TypeIISConstructAssemblyParamBlock, r *TypeIISConstructAssemblyResultBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *TypeIISConstructAssembly) validation(p TypeIISConstructAssemblyParamBlock, r *TypeIISConstructAssemblyResultBlock) {
	_wrapper := execution.NewWrapper(p.ID)
	_ = _wrapper

}

// AsyncBag functions
func (e *TypeIISConstructAssembly) Complete(params interface{}) {
	p := params.(TypeIISConstructAssemblyParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TypeIISConstructAssemblyResultBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)
	if r.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}

	e.Reaction <- execute.ThreadParam{Value: r.Reaction, ID: p.ID, Error: false}

	e.analysis(p, r)
	if r.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}

	e.validation(p, r)
	if r.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}

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
	res.Error = false || m["ReactionVolume"].(execute.ThreadParam).Error || m["PartConc"].(execute.ThreadParam).Error || m["VectorConc"].(execute.ThreadParam).Error || m["AtpVol"].(execute.ThreadParam).Error || m["ReVol"].(execute.ThreadParam).Error || m["LigVol"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["Parts"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error || m["RestrictionEnzyme"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["Ligase"].(execute.ThreadParam).Error || m["Atp"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vPartConc, is := m["PartConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vPartConc.JSONString), &temp)
		res.PartConc = *temp.PartConc
	} else {
		res.PartConc = m["PartConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vVectorConc, is := m["VectorConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vVectorConc.JSONString), &temp)
		res.VectorConc = *temp.VectorConc
	} else {
		res.VectorConc = m["VectorConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vAtpVol, is := m["AtpVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vAtpVol.JSONString), &temp)
		res.AtpVol = *temp.AtpVol
	} else {
		res.AtpVol = m["AtpVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vReVol, is := m["ReVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vReVol.JSONString), &temp)
		res.ReVol = *temp.ReVol
	} else {
		res.ReVol = m["ReVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vLigVol, is := m["LigVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vLigVol.JSONString), &temp)
		res.LigVol = *temp.LigVol
	} else {
		res.LigVol = m["LigVol"].(execute.ThreadParam).Value.(wunit.Volume)
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

	vParts, is := m["Parts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vParts.JSONString), &temp)
		res.Parts = *temp.Parts
	} else {
		res.Parts = m["Parts"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vRestrictionEnzyme, is := m["RestrictionEnzyme"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vRestrictionEnzyme.JSONString), &temp)
		res.RestrictionEnzyme = *temp.RestrictionEnzyme
	} else {
		res.RestrictionEnzyme = m["RestrictionEnzyme"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vLigase, is := m["Ligase"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vLigase.JSONString), &temp)
		res.Ligase = *temp.Ligase
	} else {
		res.Ligase = m["Ligase"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vAtp, is := m["Atp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vAtp.JSONString), &temp)
		res.Atp = *temp.Atp
	} else {
		res.Atp = m["Atp"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISConstructAssemblyJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	res.ID = m["ReactionVolume"].(execute.ThreadParam).ID

	return res
}

func (e *TypeIISConstructAssembly) OnAtp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
func (e *TypeIISConstructAssembly) OnInactivationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
func (e *TypeIISConstructAssembly) OnPartConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(17, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PartConc", param)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
		bag.Init(17, e, e)
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
func (e *TypeIISConstructAssembly) OnVectorConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(17, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VectorConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type TypeIISConstructAssembly struct {
	flow.Component    // component "superclass" embedded
	lock              sync.Mutex
	startup           sync.Once
	params            map[execute.ThreadID]*execute.AsyncBag
	ReactionVolume    <-chan execute.ThreadParam
	PartConc          <-chan execute.ThreadParam
	VectorConc        <-chan execute.ThreadParam
	AtpVol            <-chan execute.ThreadParam
	ReVol             <-chan execute.ThreadParam
	LigVol            <-chan execute.ThreadParam
	ReactionTemp      <-chan execute.ThreadParam
	ReactionTime      <-chan execute.ThreadParam
	InactivationTemp  <-chan execute.ThreadParam
	InactivationTime  <-chan execute.ThreadParam
	Parts             <-chan execute.ThreadParam
	Vector            <-chan execute.ThreadParam
	RestrictionEnzyme <-chan execute.ThreadParam
	Buffer            <-chan execute.ThreadParam
	Ligase            <-chan execute.ThreadParam
	Atp               <-chan execute.ThreadParam
	OutPlate          <-chan execute.ThreadParam
	Reaction          chan<- execute.ThreadParam
}

type TypeIISConstructAssemblyParamBlock struct {
	ID                execute.ThreadID
	Error             bool
	ReactionVolume    wunit.Volume
	PartConc          wunit.Concentration
	VectorConc        wunit.Concentration
	AtpVol            wunit.Volume
	ReVol             wunit.Volume
	LigVol            wunit.Volume
	ReactionTemp      wunit.Temperature
	ReactionTime      wunit.Time
	InactivationTemp  wunit.Temperature
	InactivationTime  wunit.Time
	Parts             []*wtype.LHComponent
	Vector            *wtype.LHComponent
	RestrictionEnzyme *wtype.LHComponent
	Buffer            *wtype.LHComponent
	Ligase            *wtype.LHComponent
	Atp               *wtype.LHComponent
	OutPlate          *wtype.LHPlate
}

type TypeIISConstructAssemblyConfig struct {
	ID                execute.ThreadID
	Error             bool
	ReactionVolume    wunit.Volume
	PartConc          wunit.Concentration
	VectorConc        wunit.Concentration
	AtpVol            wunit.Volume
	ReVol             wunit.Volume
	LigVol            wunit.Volume
	ReactionTemp      wunit.Temperature
	ReactionTime      wunit.Time
	InactivationTemp  wunit.Temperature
	InactivationTime  wunit.Time
	Parts             []wtype.FromFactory
	Vector            wtype.FromFactory
	RestrictionEnzyme wtype.FromFactory
	Buffer            wtype.FromFactory
	Ligase            wtype.FromFactory
	Atp               wtype.FromFactory
	OutPlate          wtype.FromFactory
}

type TypeIISConstructAssemblyResultBlock struct {
	ID       execute.ThreadID
	Error    bool
	Reaction *wtype.LHSolution
}

type TypeIISConstructAssemblyJSONBlock struct {
	ID                *execute.ThreadID
	Error             *bool
	ReactionVolume    *wunit.Volume
	PartConc          *wunit.Concentration
	VectorConc        *wunit.Concentration
	AtpVol            *wunit.Volume
	ReVol             *wunit.Volume
	LigVol            *wunit.Volume
	ReactionTemp      *wunit.Temperature
	ReactionTime      *wunit.Time
	InactivationTemp  *wunit.Temperature
	InactivationTime  *wunit.Time
	Parts             *[]*wtype.LHComponent
	Vector            **wtype.LHComponent
	RestrictionEnzyme **wtype.LHComponent
	Buffer            **wtype.LHComponent
	Ligase            **wtype.LHComponent
	Atp               **wtype.LHComponent
	OutPlate          **wtype.LHPlate
	Reaction          **wtype.LHSolution
}

func (c *TypeIISConstructAssembly) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartConc", "wunit.Concentration", "PartConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VectorConc", "wunit.Concentration", "VectorConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("AtpVol", "wunit.Volume", "AtpVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReVol", "wunit.Volume", "ReVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("LigVol", "wunit.Volume", "LigVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTemp", "wunit.Temperature", "ReactionTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTime", "wunit.Time", "ReactionTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTemp", "wunit.Temperature", "InactivationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTime", "wunit.Time", "InactivationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Parts", "[]*wtype.LHComponent", "Parts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "*wtype.LHComponent", "Vector", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RestrictionEnzyme", "*wtype.LHComponent", "RestrictionEnzyme", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Ligase", "*wtype.LHComponent", "Ligase", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Atp", "*wtype.LHComponent", "Atp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))

	ci := execute.NewComponentInfo("TypeIISConstructAssembly", "TypeIISConstructAssembly", "", false, inp, outp)

	return ci
}
