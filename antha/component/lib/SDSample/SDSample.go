package SDSample

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

//Input parameters for this protocol. Single instance of an SDS-PAGE sample preperation step.
//Mix 10ul of 4x stock buffer with 30ul of proteinX sample to create 40ul sample for loading.

//ProteinX
//30uL

//SDSBuffer
//10ul
//100g/L

//25g/L
//40uL

//5min
//95oC

//Biologicals

//Purified protein or cell lysate...

//Chemicals

//Consumables

//Contains protein and buffer
//Final plate with mixed components

//Biologicals

func (e *SDSample) setup(p SDSampleParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *SDSample) steps(p SDSampleParamBlock, r *SDSampleResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	//Method 1. Mix two things. DOES NOT WORK as recognises protein to be 1 single entity and wont handle as seperate components. ie end result is 5 things created all
	//from the same well. Check typeIIs workflow for hints.
	//
	//	Step1a
	//	LoadSample = MixInto(OutPlate,
	//	mixer.Sample(Protein, SampleVolume),
	//	mixer.Sample(Buffer, BufferVolume))
	//Try something else. Outputs are an array taking in a single (not array) of protein and buffer.

	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.Sample(p.Buffer, p.BufferVolume)
	bufferSample.CName = p.BufferName
	samples = append(samples, bufferSample)

	proteinSample := mixer.Sample(p.Protein, p.SampleVolume)
	proteinSample.CName = p.SampleName
	samples = append(samples, proteinSample)
	fmt.Println("This is a sample list ", samples)
	r.LoadSample = _wrapper.MixInto(p.OutPlate, samples...)

	//Methods 2.Make a sample of two things creating a list
	//	Step 1b

	//	sample	    := make([]wtype.LHComponent, 0)

	//	bufferPart  := mixer.Sample(Buffer, BufferVolume)
	//	sample	     = append([]samples, bufferSample)

	//	proteinPart := mixer.Sample(Protein, SampleVolume)
	//	sample      = append([]samples, proteinSample)

	//	LoadSample   = MixInto(OutPlate, sample...)

	//Denature the load mixture at specified temperature and time ie 95oC for 5min
	//	Step2
	_wrapper.Incubate(r.LoadSample, p.DenatureTemp, p.DenatureTime, false)
	_ = _wrapper.WaitToEnd()

	//Load the water in EPAGE gel wells
	//	Step3

	//	var water water volume
	//	waterLoad := mixer.Sample(Water, WaterLoadVolume)
	//
	//Load the LoadSample into EPAGE gel
	//
	//	Loader = MixInto(EPAGE48, LoadSample)
	//
	//
	//

	//	Status = fmtSprintln(BufferVolume.ToString() "uL of", BufferName,"mixed with", SampleVolume.ToString(), "uL of", SampleName, "Total load sample available is", ReactionVolume.ToString())
}

func (e *SDSample) analysis(p SDSampleParamBlock, r *SDSampleResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *SDSample) validation(p SDSampleParamBlock, r *SDSampleResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *SDSample) Complete(params interface{}) {
	p := params.(SDSampleParamBlock)
	if p.Error {
		e.LoadSample <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(SDSampleResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.LoadSample <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.LoadSample <- execute.ThreadParam{Value: r.LoadSample, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *SDSample) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *SDSample) NewConfig() interface{} {
	return &SDSampleConfig{}
}

func (e *SDSample) NewParamBlock() interface{} {
	return &SDSampleParamBlock{}
}

func NewSDSample() interface{} { //*SDSample {
	e := new(SDSample)
	e.init()
	return e
}

// Mapper function
func (e *SDSample) Map(m map[string]interface{}) interface{} {
	var res SDSampleParamBlock
	res.Error = false || m["Buffer"].(execute.ThreadParam).Error || m["BufferName"].(execute.ThreadParam).Error || m["BufferStockConc"].(execute.ThreadParam).Error || m["BufferVolume"].(execute.ThreadParam).Error || m["DenatureTemp"].(execute.ThreadParam).Error || m["DenatureTime"].(execute.ThreadParam).Error || m["FinalConcentration"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Protein"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["SampleName"].(execute.ThreadParam).Error || m["SampleVolume"].(execute.ThreadParam).Error

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vBufferName, is := m["BufferName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vBufferName.JSONString), &temp)
		res.BufferName = *temp.BufferName
	} else {
		res.BufferName = m["BufferName"].(execute.ThreadParam).Value.(string)
	}

	vBufferStockConc, is := m["BufferStockConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vBufferStockConc.JSONString), &temp)
		res.BufferStockConc = *temp.BufferStockConc
	} else {
		res.BufferStockConc = m["BufferStockConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vBufferVolume, is := m["BufferVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vBufferVolume.JSONString), &temp)
		res.BufferVolume = *temp.BufferVolume
	} else {
		res.BufferVolume = m["BufferVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vDenatureTemp, is := m["DenatureTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vDenatureTemp.JSONString), &temp)
		res.DenatureTemp = *temp.DenatureTemp
	} else {
		res.DenatureTemp = m["DenatureTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vDenatureTime, is := m["DenatureTime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vDenatureTime.JSONString), &temp)
		res.DenatureTime = *temp.DenatureTime
	} else {
		res.DenatureTime = m["DenatureTime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vFinalConcentration, is := m["FinalConcentration"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vFinalConcentration.JSONString), &temp)
		res.FinalConcentration = *temp.FinalConcentration
	} else {
		res.FinalConcentration = m["FinalConcentration"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vProtein, is := m["Protein"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vProtein.JSONString), &temp)
		res.Protein = *temp.Protein
	} else {
		res.Protein = m["Protein"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vSampleName, is := m["SampleName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vSampleName.JSONString), &temp)
		res.SampleName = *temp.SampleName
	} else {
		res.SampleName = m["SampleName"].(execute.ThreadParam).Value.(string)
	}

	vSampleVolume, is := m["SampleVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SDSampleJSONBlock
		json.Unmarshal([]byte(vSampleVolume.JSONString), &temp)
		res.SampleVolume = *temp.SampleVolume
	} else {
		res.SampleVolume = m["SampleVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["Buffer"].(execute.ThreadParam).ID
	res.BlockID = m["Buffer"].(execute.ThreadParam).BlockID

	return res
}

func (e *SDSample) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
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
func (e *SDSample) OnBufferName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BufferName", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnBufferStockConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BufferStockConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnBufferVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BufferVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnDenatureTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DenatureTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnDenatureTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DenatureTime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnFinalConcentration(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
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
func (e *SDSample) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
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
func (e *SDSample) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
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
func (e *SDSample) OnProtein(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Protein", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
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
func (e *SDSample) OnSampleName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SampleName", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SDSample) OnSampleVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(13, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SampleVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type SDSample struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	Buffer             <-chan execute.ThreadParam
	BufferName         <-chan execute.ThreadParam
	BufferStockConc    <-chan execute.ThreadParam
	BufferVolume       <-chan execute.ThreadParam
	DenatureTemp       <-chan execute.ThreadParam
	DenatureTime       <-chan execute.ThreadParam
	FinalConcentration <-chan execute.ThreadParam
	InPlate            <-chan execute.ThreadParam
	OutPlate           <-chan execute.ThreadParam
	Protein            <-chan execute.ThreadParam
	ReactionVolume     <-chan execute.ThreadParam
	SampleName         <-chan execute.ThreadParam
	SampleVolume       <-chan execute.ThreadParam
	LoadSample         chan<- execute.ThreadParam
	Status             chan<- execute.ThreadParam
}

type SDSampleParamBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	Buffer             *wtype.LHComponent
	BufferName         string
	BufferStockConc    wunit.Concentration
	BufferVolume       wunit.Volume
	DenatureTemp       wunit.Temperature
	DenatureTime       wunit.Time
	FinalConcentration wunit.Concentration
	InPlate            *wtype.LHPlate
	OutPlate           *wtype.LHPlate
	Protein            *wtype.LHComponent
	ReactionVolume     wunit.Volume
	SampleName         string
	SampleVolume       wunit.Volume
}

type SDSampleConfig struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	Buffer             wtype.FromFactory
	BufferName         string
	BufferStockConc    wunit.Concentration
	BufferVolume       wunit.Volume
	DenatureTemp       wunit.Temperature
	DenatureTime       wunit.Time
	FinalConcentration wunit.Concentration
	InPlate            wtype.FromFactory
	OutPlate           wtype.FromFactory
	Protein            wtype.FromFactory
	ReactionVolume     wunit.Volume
	SampleName         string
	SampleVolume       wunit.Volume
}

type SDSampleResultBlock struct {
	ID         execute.ThreadID
	BlockID    execute.BlockID
	Error      bool
	LoadSample *wtype.LHSolution
	Status     string
}

type SDSampleJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	Buffer             **wtype.LHComponent
	BufferName         *string
	BufferStockConc    *wunit.Concentration
	BufferVolume       *wunit.Volume
	DenatureTemp       *wunit.Temperature
	DenatureTime       *wunit.Time
	FinalConcentration *wunit.Concentration
	InPlate            **wtype.LHPlate
	OutPlate           **wtype.LHPlate
	Protein            **wtype.LHComponent
	ReactionVolume     *wunit.Volume
	SampleName         *string
	SampleVolume       *wunit.Volume
	LoadSample         **wtype.LHSolution
	Status             *string
}

func (c *SDSample) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferName", "string", "BufferName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferStockConc", "wunit.Concentration", "BufferStockConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BufferVolume", "wunit.Volume", "BufferVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DenatureTemp", "wunit.Temperature", "DenatureTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DenatureTime", "wunit.Time", "DenatureTime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("FinalConcentration", "wunit.Concentration", "FinalConcentration", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Protein", "*wtype.LHComponent", "Protein", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SampleName", "string", "SampleName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SampleVolume", "wunit.Volume", "SampleVolume", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("LoadSample", "*wtype.LHSolution", "LoadSample", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("SDSample", "SDSample", "", false, inp, outp)

	return ci
}
