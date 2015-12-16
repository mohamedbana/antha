package RestrictionDigestion

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

//	StockReConcinUperml 		[]int
//	DesiredConcinUperml	 		[]int

//OutputReactionName			string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func (e *RestrictionDigestion) requirements() { _ = wunit.Make_units }

// Conditions to run on startup
func (e *RestrictionDigestion) setup(p RestrictionDigestionParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *RestrictionDigestion) steps(p RestrictionDigestionParamBlock, r *RestrictionDigestionResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(p.Water, p.ReactionVolume)
	samples = append(samples, waterSample)

	bufferSample := mixer.Sample(p.Buffer, p.BufferVol)
	samples = append(samples, bufferSample)

	if p.BSAvol.Mvalue != 0 {
		bsaSample := mixer.Sample(p.BSAoptional, p.BSAvol)
		samples = append(samples, bsaSample)
	}

	// change to fixing concentration(or mass) of dna per reaction
	p.DNASolution.CName = p.DNAName
	dnaSample := mixer.Sample(p.DNASolution, p.DNAVol)
	samples = append(samples, dnaSample)

	for k, enzyme := range p.EnzSolutions {

		// work out volume to add in L

		// e.g. 1 U / (10000 * 1000) * 0.000002
		//volinL := DesiredUinreaction/(StockReConcinUperml*1000) * ReactionVolume.SIValue()
		//volumetoadd := wunit.NewVolume(volinL,"L")
		enzyme.CName = p.EnzymeNames[k]
		text.Print("adding enzyme"+p.EnzymeNames[k], "to"+p.DNAName)
		enzSample := mixer.Sample(enzyme, p.EnzVolumestoadd[k])
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
func (e *RestrictionDigestion) analysis(p RestrictionDigestionParamBlock, r *RestrictionDigestionResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *RestrictionDigestion) validation(p RestrictionDigestionParamBlock, r *RestrictionDigestionResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *RestrictionDigestion) Complete(params interface{}) {
	p := params.(RestrictionDigestionParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(RestrictionDigestionResultBlock)
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
func (e *RestrictionDigestion) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *RestrictionDigestion) NewConfig() interface{} {
	return &RestrictionDigestionConfig{}
}

func (e *RestrictionDigestion) NewParamBlock() interface{} {
	return &RestrictionDigestionParamBlock{}
}

func NewRestrictionDigestion() interface{} { //*RestrictionDigestion {
	e := new(RestrictionDigestion)
	e.init()
	return e
}

// Mapper function
func (e *RestrictionDigestion) Map(m map[string]interface{}) interface{} {
	var res RestrictionDigestionParamBlock
	res.Error = false || m["BSAoptional"].(execute.ThreadParam).Error || m["BSAvol"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["BufferVol"].(execute.ThreadParam).Error || m["DNAName"].(execute.ThreadParam).Error || m["DNASolution"].(execute.ThreadParam).Error || m["DNAVol"].(execute.ThreadParam).Error || m["EnzSolutions"].(execute.ThreadParam).Error || m["EnzVolumestoadd"].(execute.ThreadParam).Error || m["EnzymeNames"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["InactivationTemp"].(execute.ThreadParam).Error || m["InactivationTime"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["ReactionTemp"].(execute.ThreadParam).Error || m["ReactionTime"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vBSAoptional, is := m["BSAoptional"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vBSAoptional.JSONString), &temp)
		res.BSAoptional = *temp.BSAoptional
	} else {
		res.BSAoptional = m["BSAoptional"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBSAvol, is := m["BSAvol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vBSAvol.JSONString), &temp)
		res.BSAvol = *temp.BSAvol
	} else {
		res.BSAvol = m["BSAvol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferVol, is := m["BufferVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vBufferVol.JSONString), &temp)
		res.BufferVol = *temp.BufferVol
	} else {
		res.BufferVol = m["BufferVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vDNAName, is := m["DNAName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vDNAName.JSONString), &temp)
		res.DNAName = *temp.DNAName
	} else {
		res.DNAName = m["DNAName"].(execute.ThreadParam).Value.(string)
	}

	vDNASolution, is := m["DNASolution"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vDNASolution.JSONString), &temp)
		res.DNASolution = *temp.DNASolution
	} else {
		res.DNASolution = m["DNASolution"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDNAVol, is := m["DNAVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vDNAVol.JSONString), &temp)
		res.DNAVol = *temp.DNAVol
	} else {
		res.DNAVol = m["DNAVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vEnzSolutions, is := m["EnzSolutions"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vEnzSolutions.JSONString), &temp)
		res.EnzSolutions = *temp.EnzSolutions
	} else {
		res.EnzSolutions = m["EnzSolutions"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vEnzVolumestoadd, is := m["EnzVolumestoadd"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vEnzVolumestoadd.JSONString), &temp)
		res.EnzVolumestoadd = *temp.EnzVolumestoadd
	} else {
		res.EnzVolumestoadd = m["EnzVolumestoadd"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vEnzymeNames, is := m["EnzymeNames"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vEnzymeNames.JSONString), &temp)
		res.EnzymeNames = *temp.EnzymeNames
	} else {
		res.EnzymeNames = m["EnzymeNames"].(execute.ThreadParam).Value.([]string)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInactivationTemp, is := m["InactivationTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vInactivationTemp.JSONString), &temp)
		res.InactivationTemp = *temp.InactivationTemp
	} else {
		res.InactivationTemp = m["InactivationTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vInactivationTime, is := m["InactivationTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vInactivationTime.JSONString), &temp)
		res.InactivationTime = *temp.InactivationTime
	} else {
		res.InactivationTime = m["InactivationTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vReactionTemp, is := m["ReactionTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vReactionTemp.JSONString), &temp)
		res.ReactionTemp = *temp.ReactionTemp
	} else {
		res.ReactionTemp = m["ReactionTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vReactionTime, is := m["ReactionTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vReactionTime.JSONString), &temp)
		res.ReactionTime = *temp.ReactionTime
	} else {
		res.ReactionTime = m["ReactionTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp RestrictionDigestionJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["BSAoptional"].(execute.ThreadParam).ID
	res.BlockID = m["BSAoptional"].(execute.ThreadParam).BlockID

	return res
}

func (e *RestrictionDigestion) OnBSAoptional(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnBSAvol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnBufferVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnDNAName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnDNASolution(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnDNAVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion) OnEnzSolutions(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnEnzVolumestoadd(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("EnzVolumestoadd", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *RestrictionDigestion) OnEnzymeNames(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnInactivationTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnInactivationTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnReactionTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnReactionTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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
func (e *RestrictionDigestion) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(18, e, e)
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

type RestrictionDigestion struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	BSAoptional      <-chan execute.ThreadParam
	BSAvol           <-chan execute.ThreadParam
	Buffer           <-chan execute.ThreadParam
	BufferVol        <-chan execute.ThreadParam
	DNAName          <-chan execute.ThreadParam
	DNASolution      <-chan execute.ThreadParam
	DNAVol           <-chan execute.ThreadParam
	EnzSolutions     <-chan execute.ThreadParam
	EnzVolumestoadd  <-chan execute.ThreadParam
	EnzymeNames      <-chan execute.ThreadParam
	InPlate          <-chan execute.ThreadParam
	InactivationTemp <-chan execute.ThreadParam
	InactivationTime <-chan execute.ThreadParam
	OutPlate         <-chan execute.ThreadParam
	ReactionTemp     <-chan execute.ThreadParam
	ReactionTime     <-chan execute.ThreadParam
	ReactionVolume   <-chan execute.ThreadParam
	Water            <-chan execute.ThreadParam
	Reaction         chan<- execute.ThreadParam
}

type RestrictionDigestionParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	BSAoptional      *wtype.LHComponent
	BSAvol           wunit.Volume
	Buffer           *wtype.LHComponent
	BufferVol        wunit.Volume
	DNAName          string
	DNASolution      *wtype.LHComponent
	DNAVol           wunit.Volume
	EnzSolutions     []*wtype.LHComponent
	EnzVolumestoadd  []wunit.Volume
	EnzymeNames      []string
	InPlate          *wtype.LHPlate
	InactivationTemp wunit.Temperature
	InactivationTime wunit.Time
	OutPlate         *wtype.LHPlate
	ReactionTemp     wunit.Temperature
	ReactionTime     wunit.Time
	ReactionVolume   wunit.Volume
	Water            *wtype.LHComponent
}

type RestrictionDigestionConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	BSAoptional      wtype.FromFactory
	BSAvol           wunit.Volume
	Buffer           wtype.FromFactory
	BufferVol        wunit.Volume
	DNAName          string
	DNASolution      wtype.FromFactory
	DNAVol           wunit.Volume
	EnzSolutions     []wtype.FromFactory
	EnzVolumestoadd  []wunit.Volume
	EnzymeNames      []string
	InPlate          wtype.FromFactory
	InactivationTemp wunit.Temperature
	InactivationTime wunit.Time
	OutPlate         wtype.FromFactory
	ReactionTemp     wunit.Temperature
	ReactionTime     wunit.Time
	ReactionVolume   wunit.Volume
	Water            wtype.FromFactory
}

type RestrictionDigestionResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Reaction *wtype.LHSolution
}

type RestrictionDigestionJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	BSAoptional      **wtype.LHComponent
	BSAvol           *wunit.Volume
	Buffer           **wtype.LHComponent
	BufferVol        *wunit.Volume
	DNAName          *string
	DNASolution      **wtype.LHComponent
	DNAVol           *wunit.Volume
	EnzSolutions     *[]*wtype.LHComponent
	EnzVolumestoadd  *[]wunit.Volume
	EnzymeNames      *[]string
	InPlate          **wtype.LHPlate
	InactivationTemp *wunit.Temperature
	InactivationTime *wunit.Time
	OutPlate         **wtype.LHPlate
	ReactionTemp     *wunit.Temperature
	ReactionTime     *wunit.Time
	ReactionVolume   *wunit.Volume
	Water            **wtype.LHComponent
	Reaction         **wtype.LHSolution
}

func (c *RestrictionDigestion) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("BSAoptional", "*wtype.LHComponent", "BSAoptional", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BSAvol", "wunit.Volume", "BSAvol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferVol", "wunit.Volume", "BufferVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAName", "string", "DNAName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNASolution", "*wtype.LHComponent", "DNASolution", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAVol", "wunit.Volume", "DNAVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("EnzSolutions", "[]*wtype.LHComponent", "EnzSolutions", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("EnzVolumestoadd", "[]wunit.Volume", "EnzVolumestoadd", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("EnzymeNames", "[]string", "EnzymeNames", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTemp", "wunit.Temperature", "InactivationTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InactivationTime", "wunit.Time", "InactivationTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTemp", "wunit.Temperature", "ReactionTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionTime", "wunit.Time", "ReactionTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))

	ci := execute.NewComponentInfo("RestrictionDigestion", "RestrictionDigestion", "", false, inp, outp)

	return ci
}
