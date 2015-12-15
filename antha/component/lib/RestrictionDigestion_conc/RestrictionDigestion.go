package RestrictionDigestion_conc

import (
	"encoding/json"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
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

//DNAVol						Volume

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func (e *RestrictionDigestion_conc) requirements() { _ = wunit.Make_units }

// Conditions to run on startup
func (e *RestrictionDigestion_conc) setup(p RestrictionDigestion_concParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *RestrictionDigestion_conc) steps(p RestrictionDigestion_concParamBlock, r *RestrictionDigestion_concResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(p.Water, p.ReactionVolume)
	samples = append(samples, waterSample)

	// workout volume of buffer to add in SI units
	BufferVol := wunit.NewVolume(1000000*float64(p.ReactionVolume.SIValue()/float64(p.BufferConcX)), "ul")

	bufferSample := mixer.Sample(p.Buffer, BufferVol)
	samples = append(samples, bufferSample)

	if p.BSAvol.Mvalue != 0 {
		bsaSample := mixer.Sample(p.BSAoptional, p.BSAvol)
		samples = append(samples, bsaSample)
	}

	p.DNASolution.CName = p.DNAName

	// work out necessary volume to add
	DNAVol := wunit.NewVolume(float64(1000000*(p.DNAMassperReaction.SIValue()/p.DNAConc.SIValue())), "ul")
	text.Print("DNAVOL", DNAVol.ToString())
	dnaSample := mixer.Sample(p.DNASolution, DNAVol)
	samples = append(samples, dnaSample)

	for k, enzyme := range p.EnzSolutions {

		/*
			e.g.
			DesiredUinreaction = 1  // U
			StockReConcinUperml = 10000 // U/ml
			ReactionVolume = 20ul
		*/
		stockconcinUperul := p.StockReConcinUperml[k] / 1000
		enzvoltoaddinul := p.DesiredConcinUperml[k] / stockconcinUperul

		var enzvoltoadd wunit.Volume

		if enzvoltoaddinul < 1 {
			enzvoltoadd = wunit.NewVolume(float64(1), "ul")
		} else {
			enzvoltoadd = wunit.NewVolume(float64(enzvoltoaddinul), "ul")
		}
		enzyme.CName = p.EnzymeNames[k]
		text.Print("adding enzyme"+p.EnzymeNames[k], "to"+p.DNAName)
		enzSample := mixer.Sample(enzyme, enzvoltoadd)
		enzSample.CName = p.EnzymeNames[k]
		samples = append(samples, enzSample)
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
func (e *RestrictionDigestion_conc) analysis(p RestrictionDigestion_concParamBlock, r *RestrictionDigestion_concResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *RestrictionDigestion_conc) validation(p RestrictionDigestion_concParamBlock, r *RestrictionDigestion_concResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *RestrictionDigestion_conc) Complete(params interface{}) {
	p := params.(RestrictionDigestion_concParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(RestrictionDigestion_concResultBlock)
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
func (e *RestrictionDigestion_conc) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *RestrictionDigestion_conc) NewConfig() interface{} {
	return &RestrictionDigestion_concConfig{}
}

func (e *RestrictionDigestion_conc) NewParamBlock() interface{} {
	return &RestrictionDigestion_concParamBlock{}
}

func NewRestrictionDigestion_conc() interface{} { //*RestrictionDigestion_conc {
	e := new(RestrictionDigestion_conc)
	e.init()
	return e
}

// Mapper function
func (e *RestrictionDigestion_conc) Map(m map[string]interface{}) interface{} {
	var res RestrictionDigestion_concParamBlock
	res.Error = false || m["BSAoptional"].(execute.ThreadParam).Error || m["BSAvol"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["BufferConcX"].(execute.ThreadParam).Error || m["DNAConc"].(execute.ThreadParam).Error || m["DNAMassperReaction"].(execute.ThreadParam).Error || m["DNAName"].(execute.ThreadParam).Error || m["DNASolution"].(execute.ThreadParam).Error || m["DesiredConcinUperml"].(execute.ThreadParam).Error || m["EnzSolutions"].(execute.ThreadParam).Error || m["EnzymeNames"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["StockReConcinUperml"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vBSAoptional, is := m["BSAoptional"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vBSAoptional.JSONString), &temp)
		res.BSAoptional = *temp.BSAoptional
	} else {
		res.BSAoptional = m["BSAoptional"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBSAvol, is := m["BSAvol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vBSAvol.JSONString), &temp)
		res.BSAvol = *temp.BSAvol
	} else {
		res.BSAvol = m["BSAvol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferConcX, is := m["BufferConcX"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vBufferConcX.JSONString), &temp)
		res.BufferConcX = *temp.BufferConcX
	} else {
		res.BufferConcX = m["BufferConcX"].(execute.ThreadParam).Value.(int)
	}

	vDNAConc, is := m["DNAConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vDNAConc.JSONString), &temp)
		res.DNAConc = *temp.DNAConc
	} else {
		res.DNAConc = m["DNAConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vDNAMassperReaction, is := m["DNAMassperReaction"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vDNAMassperReaction.JSONString), &temp)
		res.DNAMassperReaction = *temp.DNAMassperReaction
	} else {
		res.DNAMassperReaction = m["DNAMassperReaction"].(execute.ThreadParam).Value.(wunit.Mass)
	}

	vDNAName, is := m["DNAName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vDNAName.JSONString), &temp)
		res.DNAName = *temp.DNAName
	} else {
		res.DNAName = m["DNAName"].(execute.ThreadParam).Value.(string)
	}

	vDNASolution, is := m["DNASolution"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vDNASolution.JSONString), &temp)
		res.DNASolution = *temp.DNASolution
	} else {
		res.DNASolution = m["DNASolution"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDesiredConcinUperml, is := m["DesiredConcinUperml"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vDesiredConcinUperml.JSONString), &temp)
		res.DesiredConcinUperml = *temp.DesiredConcinUperml
	} else {
		res.DesiredConcinUperml = m["DesiredConcinUperml"].(execute.ThreadParam).Value.([]int)
	}

	vEnzSolutions, is := m["EnzSolutions"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vEnzSolutions.JSONString), &temp)
		res.EnzSolutions = *temp.EnzSolutions
	} else {
		res.EnzSolutions = m["EnzSolutions"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vEnzymeNames, is := m["EnzymeNames"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vEnzymeNames.JSONString), &temp)
		res.EnzymeNames = *temp.EnzymeNames
	} else {
		res.EnzymeNames = m["EnzymeNames"].(execute.ThreadParam).Value.([]string)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInactivationTemp, is := m["InactivationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vInactivationTemp.JSONString), &temp)
		res.InactivationTemp = *temp.InactivationTemp
	} else {
		res.InactivationTemp = m["InactivationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vInactivationTime, is := m["InactivationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vInactivationTime.JSONString), &temp)
		res.InactivationTime = *temp.InactivationTime
	} else {
		res.InactivationTime = m["InactivationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vReactionTemp, is := m["ReactionTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vReactionTemp.JSONString), &temp)
		res.ReactionTemp = *temp.ReactionTemp
	} else {
		res.ReactionTemp = m["ReactionTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vReactionTime, is := m["ReactionTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vReactionTime.JSONString), &temp)
		res.ReactionTime = *temp.ReactionTime
	} else {
		res.ReactionTime = m["ReactionTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vStockReConcinUperml, is := m["StockReConcinUperml"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vStockReConcinUperml.JSONString), &temp)
		res.StockReConcinUperml = *temp.StockReConcinUperml
	} else {
		res.StockReConcinUperml = m["StockReConcinUperml"].(execute.ThreadParam).Value.([]int)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestion_concJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["BSAoptional"].(execute.ThreadParam).ID
	res.BlockID = m["BSAoptional"].(execute.ThreadParam).BlockID

	return res
}

func (e *RestrictionDigestion_conc) OnBSAoptional(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BSAoptional", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnBSAvol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BSAvol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnBufferConcX(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BufferConcX", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnDNAConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnDNAMassperReaction(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAMassperReaction", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnDNAName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAName", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnDNASolution(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNASolution", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnDesiredConcinUperml(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DesiredConcinUperml", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnEnzSolutions(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("EnzSolutions", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnEnzymeNames(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("EnzymeNames", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnInactivationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnInactivationTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnReactionTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnReactionTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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
func (e *RestrictionDigestion_conc) OnStockReConcinUperml(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("StockReConcinUperml", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion_conc) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(20, e, e)
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

type RestrictionDigestion_conc struct {
	flow.Component      // component "superclass" embedded
	lock                sync.Mutex
	startup             sync.Once
	params              map[execute.ThreadID]*execute.AsyncBag
	BSAoptional         <-chan execute.ThreadParam
	BSAvol              <-chan execute.ThreadParam
	Buffer              <-chan execute.ThreadParam
	BufferConcX         <-chan execute.ThreadParam
	DNAConc             <-chan execute.ThreadParam
	DNAMassperReaction  <-chan execute.ThreadParam
	DNAName             <-chan execute.ThreadParam
	DNASolution         <-chan execute.ThreadParam
	DesiredConcinUperml <-chan execute.ThreadParam
	EnzSolutions        <-chan execute.ThreadParam
	EnzymeNames         <-chan execute.ThreadParam
	InPlate             <-chan execute.ThreadParam
	InactivationTemp    <-chan execute.ThreadParam
	InactivationTime    <-chan execute.ThreadParam
	OutPlate            <-chan execute.ThreadParam
	ReactionTemp        <-chan execute.ThreadParam
	ReactionTime        <-chan execute.ThreadParam
	ReactionVolume      <-chan execute.ThreadParam
	StockReConcinUperml <-chan execute.ThreadParam
	Water               <-chan execute.ThreadParam
	Reaction            chan<- execute.ThreadParam
}

type RestrictionDigestion_concParamBlock struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	BSAoptional         *wtype.LHComponent
	BSAvol              wunit.Volume
	Buffer              *wtype.LHComponent
	BufferConcX         int
	DNAConc             wunit.Concentration
	DNAMassperReaction  wunit.Mass
	DNAName             string
	DNASolution         *wtype.LHComponent
	DesiredConcinUperml []int
	EnzSolutions        []*wtype.LHComponent
	EnzymeNames         []string
	InPlate             *wtype.LHPlate
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	OutPlate            *wtype.LHPlate
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	StockReConcinUperml []int
	Water               *wtype.LHComponent
}

type RestrictionDigestion_concConfig struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	BSAoptional         wtype.FromFactory
	BSAvol              wunit.Volume
	Buffer              wtype.FromFactory
	BufferConcX         int
	DNAConc             wunit.Concentration
	DNAMassperReaction  wunit.Mass
	DNAName             string
	DNASolution         wtype.FromFactory
	DesiredConcinUperml []int
	EnzSolutions        []wtype.FromFactory
	EnzymeNames         []string
	InPlate             wtype.FromFactory
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	OutPlate            wtype.FromFactory
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	StockReConcinUperml []int
	Water               wtype.FromFactory
}

type RestrictionDigestion_concResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Reaction *wtype.LHSolution
}

type RestrictionDigestion_concJSONBlock struct {
	ID                  *execute.ThreadID
	BlockID             *execute.BlockID
	Error               *bool
	BSAoptional         **wtype.LHComponent
	BSAvol              *wunit.Volume
	Buffer              **wtype.LHComponent
	BufferConcX         *int
	DNAConc             *wunit.Concentration
	DNAMassperReaction  *wunit.Mass
	DNAName             *string
	DNASolution         **wtype.LHComponent
	DesiredConcinUperml *[]int
	EnzSolutions        *[]*wtype.LHComponent
	EnzymeNames         *[]string
	InPlate             **wtype.LHPlate
	InactivationTemp    *wunit.Temperature
	InactivationTime    *wunit.Time
	OutPlate            **wtype.LHPlate
	ReactionTemp        *wunit.Temperature
	ReactionTime        *wunit.Time
	ReactionVolume      *wunit.Volume
	StockReConcinUperml *[]int
	Water               **wtype.LHComponent
	Reaction            **wtype.LHSolution
}

func (c *RestrictionDigestion_conc) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("BSAoptional", "*wtype.LHComponent", "BSAoptional", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BSAvol", "wunit.Volume", "BSAvol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferConcX", "int", "BufferConcX", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAConc", "wunit.Concentration", "DNAConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAMassperReaction", "wunit.Mass", "DNAMassperReaction", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAName", "string", "DNAName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNASolution", "*wtype.LHComponent", "DNASolution", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DesiredConcinUperml", "[]int", "DesiredConcinUperml", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("EnzSolutions", "[]*wtype.LHComponent", "EnzSolutions", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("EnzymeNames", "[]string", "EnzymeNames", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTemp", "wunit.Temperature", "InactivationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTime", "wunit.Time", "InactivationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTemp", "wunit.Temperature", "ReactionTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTime", "wunit.Time", "ReactionTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("StockReConcinUperml", "[]int", "StockReConcinUperml", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))

	ci := execute.NewComponentInfo("RestrictionDigestion_conc", "RestrictionDigestion_conc", "", false, inp, outp)

	return ci
}
