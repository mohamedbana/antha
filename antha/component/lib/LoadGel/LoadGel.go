package LoadGel

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

//    RunVoltage      Int
//    RunLength       Time

//preload well with 10uL of water
//protein samples for running
//96 well plate with water, marker and samples
//Gel to load ie OutPlate

//Run length in cm, and protein band height and pixed density after digital scanning

func (e *LoadGel) setup(p LoadGelParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *LoadGel) steps(p LoadGelParamBlock, r *LoadGelResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.Sample(p.Water, p.WaterVolume)
	waterSample.CName = p.WaterName
	samples = append(samples, waterSample)

	loadSample := mixer.Sample(p.Protein, p.LoadVolume)
	loadSample.CName = p.SampleName
	samples = append(samples, loadSample)
	fmt.Println("This is a list of samples for loading:", samples)

	r.RunSolution = _wrapper.MixInto(p.GelPlate, samples...)
	_ = _wrapper.WaitToEnd()

}

func (e *LoadGel) analysis(p LoadGelParamBlock, r *LoadGelResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *LoadGel) validation(p LoadGelParamBlock, r *LoadGelResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *LoadGel) Complete(params interface{}) {
	p := params.(LoadGelParamBlock)
	if p.Error {
		e.RunSolution <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(LoadGelResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.RunSolution <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.RunSolution <- execute.ThreadParam{Value: r.RunSolution, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *LoadGel) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *LoadGel) NewConfig() interface{} {
	return &LoadGelConfig{}
}

func (e *LoadGel) NewParamBlock() interface{} {
	return &LoadGelParamBlock{}
}

func NewLoadGel() interface{} { //*LoadGel {
	e := new(LoadGel)
	e.init()
	return e
}

// Mapper function
func (e *LoadGel) Map(m map[string]interface{}) interface{} {
	var res LoadGelParamBlock
	res.Error = false || m["GelPlate"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["LoadVolume"].(execute.ThreadParam).Error || m["Protein"].(execute.ThreadParam).Error || m["SampleName"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error || m["WaterName"].(execute.ThreadParam).Error || m["WaterVolume"].(execute.ThreadParam).Error

	vGelPlate, is := m["GelPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vGelPlate.JSONString), &temp)
		res.GelPlate = *temp.GelPlate
	} else {
		res.GelPlate = m["GelPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vLoadVolume, is := m["LoadVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vLoadVolume.JSONString), &temp)
		res.LoadVolume = *temp.LoadVolume
	} else {
		res.LoadVolume = m["LoadVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vProtein, is := m["Protein"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vProtein.JSONString), &temp)
		res.Protein = *temp.Protein
	} else {
		res.Protein = m["Protein"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vSampleName, is := m["SampleName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vSampleName.JSONString), &temp)
		res.SampleName = *temp.SampleName
	} else {
		res.SampleName = m["SampleName"].(execute.ThreadParam).Value.(string)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vWaterName, is := m["WaterName"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vWaterName.JSONString), &temp)
		res.WaterName = *temp.WaterName
	} else {
		res.WaterName = m["WaterName"].(execute.ThreadParam).Value.(string)
	}

	vWaterVolume, is := m["WaterVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LoadGelJSONBlock
		json.Unmarshal([]byte(vWaterVolume.JSONString), &temp)
		res.WaterVolume = *temp.WaterVolume
	} else {
		res.WaterVolume = m["WaterVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["GelPlate"].(execute.ThreadParam).ID
	res.BlockID = m["GelPlate"].(execute.ThreadParam).BlockID

	return res
}

func (e *LoadGel) OnGelPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("GelPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *LoadGel) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *LoadGel) OnLoadVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("LoadVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *LoadGel) OnProtein(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *LoadGel) OnSampleName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *LoadGel) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *LoadGel) OnWaterName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("WaterName", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *LoadGel) OnWaterVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("WaterVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type LoadGel struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	GelPlate       <-chan execute.ThreadParam
	InPlate        <-chan execute.ThreadParam
	LoadVolume     <-chan execute.ThreadParam
	Protein        <-chan execute.ThreadParam
	SampleName     <-chan execute.ThreadParam
	Water          <-chan execute.ThreadParam
	WaterName      <-chan execute.ThreadParam
	WaterVolume    <-chan execute.ThreadParam
	RunSolution    chan<- execute.ThreadParam
	Status         chan<- execute.ThreadParam
}

type LoadGelParamBlock struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	GelPlate    *wtype.LHPlate
	InPlate     *wtype.LHPlate
	LoadVolume  wunit.Volume
	Protein     *wtype.LHComponent
	SampleName  string
	Water       *wtype.LHComponent
	WaterName   string
	WaterVolume wunit.Volume
}

type LoadGelConfig struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	GelPlate    wtype.FromFactory
	InPlate     wtype.FromFactory
	LoadVolume  wunit.Volume
	Protein     wtype.FromFactory
	SampleName  string
	Water       wtype.FromFactory
	WaterName   string
	WaterVolume wunit.Volume
}

type LoadGelResultBlock struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	RunSolution *wtype.LHSolution
	Status      string
}

type LoadGelJSONBlock struct {
	ID          *execute.ThreadID
	BlockID     *execute.BlockID
	Error       *bool
	GelPlate    **wtype.LHPlate
	InPlate     **wtype.LHPlate
	LoadVolume  *wunit.Volume
	Protein     **wtype.LHComponent
	SampleName  *string
	Water       **wtype.LHComponent
	WaterName   *string
	WaterVolume *wunit.Volume
	RunSolution **wtype.LHSolution
	Status      *string
}

func (c *LoadGel) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("GelPlate", "*wtype.LHPlate", "GelPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("LoadVolume", "wunit.Volume", "LoadVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Protein", "*wtype.LHComponent", "Protein", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SampleName", "string", "SampleName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("WaterName", "string", "WaterName", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("WaterVolume", "wunit.Volume", "WaterVolume", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("RunSolution", "*wtype.LHSolution", "RunSolution", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("LoadGel", "LoadGel", "", false, inp, outp)

	return ci
}
