package Transformation

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

//= 50.(uL)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Transformation) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Transformation) setup(p TransformationParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Transformation) steps(p TransformationParamBlock, r *TransformationResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	competetentcellmix := mixer.Sample(p.ReadyCompCells, p.CompetentCellvolumeperassembly)
	transformationmix := make([]*wtype.LHComponent, 0)
	transformationmix = append(transformationmix, competetentcellmix)
	DNAsample := mixer.Sample(p.Reaction, p.Reactionvolume)
	transformationmix = append(transformationmix, DNAsample)

	transformedcells := _wrapper.MixInto(p.OutPlate, transformationmix...)

	_wrapper.Incubate(transformedcells, p.Postplasmidtemp, p.Postplasmidtime, false)

	r.Transformedcells = transformedcells
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Transformation) analysis(p TransformationParamBlock, r *TransformationResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Transformation) validation(p TransformationParamBlock, r *TransformationResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Transformation) Complete(params interface{}) {
	p := params.(TransformationParamBlock)
	if p.Error {
		e.Transformedcells <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TransformationResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Transformedcells <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Transformedcells <- execute.ThreadParam{Value: r.Transformedcells, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Transformation) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Transformation) NewConfig() interface{} {
	return &TransformationConfig{}
}

func (e *Transformation) NewParamBlock() interface{} {
	return &TransformationParamBlock{}
}

func NewTransformation() interface{} { //*Transformation {
	e := new(Transformation)
	e.init()
	return e
}

// Mapper function
func (e *Transformation) Map(m map[string]interface{}) interface{} {
	var res TransformationParamBlock
	res.Error = false || m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Postplasmidtemp"].(execute.ThreadParam).Error || m["Postplasmidtime"].(execute.ThreadParam).Error || m["Reaction"].(execute.ThreadParam).Error || m["Reactionvolume"].(execute.ThreadParam).Error || m["ReadyCompCells"].(execute.ThreadParam).Error

	vCompetentCellvolumeperassembly, is := m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vCompetentCellvolumeperassembly.JSONString), &temp)
		res.CompetentCellvolumeperassembly = *temp.CompetentCellvolumeperassembly
	} else {
		res.CompetentCellvolumeperassembly = m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPostplasmidtemp, is := m["Postplasmidtemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vPostplasmidtemp.JSONString), &temp)
		res.Postplasmidtemp = *temp.Postplasmidtemp
	} else {
		res.Postplasmidtemp = m["Postplasmidtemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vPostplasmidtime, is := m["Postplasmidtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vPostplasmidtime.JSONString), &temp)
		res.Postplasmidtime = *temp.Postplasmidtime
	} else {
		res.Postplasmidtime = m["Postplasmidtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReaction, is := m["Reaction"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vReaction.JSONString), &temp)
		res.Reaction = *temp.Reaction
	} else {
		res.Reaction = m["Reaction"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vReactionvolume, is := m["Reactionvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vReactionvolume.JSONString), &temp)
		res.Reactionvolume = *temp.Reactionvolume
	} else {
		res.Reactionvolume = m["Reactionvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vReadyCompCells, is := m["ReadyCompCells"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransformationJSONBlock
		json.Unmarshal([]byte(vReadyCompCells.JSONString), &temp)
		res.ReadyCompCells = *temp.ReadyCompCells
	} else {
		res.ReadyCompCells = m["ReadyCompCells"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["CompetentCellvolumeperassembly"].(execute.ThreadParam).ID
	res.BlockID = m["CompetentCellvolumeperassembly"].(execute.ThreadParam).BlockID

	return res
}

func (e *Transformation) OnCompetentCellvolumeperassembly(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("CompetentCellvolumeperassembly", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *Transformation) OnPostplasmidtemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Postplasmidtemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation) OnPostplasmidtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Postplasmidtime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation) OnReaction(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Reaction", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation) OnReactionvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Reactionvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation) OnReadyCompCells(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ReadyCompCells", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Transformation struct {
	flow.Component                 // component "superclass" embedded
	lock                           sync.Mutex
	startup                        sync.Once
	params                         map[execute.ThreadID]*execute.AsyncBag
	CompetentCellvolumeperassembly <-chan execute.ThreadParam
	OutPlate                       <-chan execute.ThreadParam
	Postplasmidtemp                <-chan execute.ThreadParam
	Postplasmidtime                <-chan execute.ThreadParam
	Reaction                       <-chan execute.ThreadParam
	Reactionvolume                 <-chan execute.ThreadParam
	ReadyCompCells                 <-chan execute.ThreadParam
	Transformedcells               chan<- execute.ThreadParam
}

type TransformationParamBlock struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       *wtype.LHPlate
	Postplasmidtemp                wunit.Temperature
	Postplasmidtime                wunit.Time
	Reaction                       *wtype.LHComponent
	Reactionvolume                 wunit.Volume
	ReadyCompCells                 *wtype.LHComponent
}

type TransformationConfig struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       wtype.FromFactory
	Postplasmidtemp                wunit.Temperature
	Postplasmidtime                wunit.Time
	Reaction                       wtype.FromFactory
	Reactionvolume                 wunit.Volume
	ReadyCompCells                 wtype.FromFactory
}

type TransformationResultBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Transformedcells *wtype.LHSolution
}

type TransformationJSONBlock struct {
	ID                             *execute.ThreadID
	BlockID                        *execute.BlockID
	Error                          *bool
	CompetentCellvolumeperassembly *wunit.Volume
	OutPlate                       **wtype.LHPlate
	Postplasmidtemp                *wunit.Temperature
	Postplasmidtime                *wunit.Time
	Reaction                       **wtype.LHComponent
	Reactionvolume                 *wunit.Volume
	ReadyCompCells                 **wtype.LHComponent
	Transformedcells               **wtype.LHSolution
}

func (c *Transformation) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("CompetentCellvolumeperassembly", "wunit.Volume", "CompetentCellvolumeperassembly", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Postplasmidtemp", "wunit.Temperature", "Postplasmidtemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Postplasmidtime", "wunit.Time", "Postplasmidtime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Reaction", "*wtype.LHComponent", "Reaction", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Reactionvolume", "wunit.Volume", "Reactionvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReadyCompCells", "*wtype.LHComponent", "ReadyCompCells", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Transformedcells", "*wtype.LHSolution", "Transformedcells", true, true, nil, nil))

	ci := execute.NewComponentInfo("Transformation", "Transformation", "", false, inp, outp)

	return ci
}
