package Mastermix_reactions

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

// if buffer is being added

// add as many as possible option

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// optional if nil this is ignored

// Physical outputs from this protocol with types

func (e *Mastermix_reactions) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Mastermix_reactions) setup(p Mastermix_reactionsParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Mastermix_reactions) steps(p Mastermix_reactionsParamBlock, r *Mastermix_reactionsResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	// work out volume to top up to in each case (per reaction):
	topupVolumeperreacttion := p.TotalVolumeperreaction.SIValue() - p.VolumetoLeaveforDNAperreaction.SIValue()

	// multiply by number of reactions per mastermix
	topupVolume := wunit.NewVolume(float64(p.Reactionspermastermix)*topupVolumeperreacttion, "l")

	if len(p.Components) != len(p.ComponentVolumesperReaction) {
		panic("len(Components) != len(OtherComponentVolumes)")
	}

	mastermixes := make([]*wtype.LHSolution, 0)

	if p.AliquotbyRow {
		panic("MixTo based method coming soon!")
	} else {
		for i := 0; i < p.NumberofMastermixes; i++ {

			eachmastermix := make([]*wtype.LHComponent, 0)

			if p.TopUpBuffer != nil {
				bufferSample := mixer.SampleForTotalVolume(p.TopUpBuffer, topupVolume)
				eachmastermix = append(eachmastermix, bufferSample)
			}

			for k, component := range p.Components {
				if k == len(p.Components) {
					component.Type = "NeedToMix"
				}

				// multiply volume of each component by number of reactions per mastermix
				adjustedvol := wunit.NewVolume(float64(p.Reactionspermastermix)*p.ComponentVolumesperReaction[k].SIValue(), "l")

				componentSample := mixer.Sample(component, adjustedvol)
				eachmastermix = append(eachmastermix, componentSample)
			}

			mastermix := _wrapper.MixInto(p.OutPlate, eachmastermix...)
			mastermixes = append(mastermixes, mastermix)

		}

	}
	r.Mastermixes = mastermixes
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Mastermix_reactions) analysis(p Mastermix_reactionsParamBlock, r *Mastermix_reactionsResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Mastermix_reactions) validation(p Mastermix_reactionsParamBlock, r *Mastermix_reactionsResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Mastermix_reactions) Complete(params interface{}) {
	p := params.(Mastermix_reactionsParamBlock)
	if p.Error {
		e.Mastermixes <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(Mastermix_reactionsResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Mastermixes <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Mastermixes <- execute.ThreadParam{Value: r.Mastermixes, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Mastermix_reactions) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Mastermix_reactions) NewConfig() interface{} {
	return &Mastermix_reactionsConfig{}
}

func (e *Mastermix_reactions) NewParamBlock() interface{} {
	return &Mastermix_reactionsParamBlock{}
}

func NewMastermix_reactions() interface{} { //*Mastermix_reactions {
	e := new(Mastermix_reactions)
	e.init()
	return e
}

// Mapper function
func (e *Mastermix_reactions) Map(m map[string]interface{}) interface{} {
	var res Mastermix_reactionsParamBlock
	res.Error = false || m["AliquotbyRow"].(execute.ThreadParam).Error || m["ComponentVolumesperReaction"].(execute.ThreadParam).Error || m["Components"].(execute.ThreadParam).Error || m["Inplate"].(execute.ThreadParam).Error || m["NumberofMastermixes"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Reactionspermastermix"].(execute.ThreadParam).Error || m["TopUpBuffer"].(execute.ThreadParam).Error || m["TotalVolumeperreaction"].(execute.ThreadParam).Error || m["VolumetoLeaveforDNAperreaction"].(execute.ThreadParam).Error

	vAliquotbyRow, is := m["AliquotbyRow"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vAliquotbyRow.JSONString), &temp)
		res.AliquotbyRow = *temp.AliquotbyRow
	} else {
		res.AliquotbyRow = m["AliquotbyRow"].(execute.ThreadParam).Value.(bool)
	}

	vComponentVolumesperReaction, is := m["ComponentVolumesperReaction"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vComponentVolumesperReaction.JSONString), &temp)
		res.ComponentVolumesperReaction = *temp.ComponentVolumesperReaction
	} else {
		res.ComponentVolumesperReaction = m["ComponentVolumesperReaction"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vComponents, is := m["Components"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vComponents.JSONString), &temp)
		res.Components = *temp.Components
	} else {
		res.Components = m["Components"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vInplate, is := m["Inplate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vInplate.JSONString), &temp)
		res.Inplate = *temp.Inplate
	} else {
		res.Inplate = m["Inplate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vNumberofMastermixes, is := m["NumberofMastermixes"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vNumberofMastermixes.JSONString), &temp)
		res.NumberofMastermixes = *temp.NumberofMastermixes
	} else {
		res.NumberofMastermixes = m["NumberofMastermixes"].(execute.ThreadParam).Value.(int)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vReactionspermastermix, is := m["Reactionspermastermix"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vReactionspermastermix.JSONString), &temp)
		res.Reactionspermastermix = *temp.Reactionspermastermix
	} else {
		res.Reactionspermastermix = m["Reactionspermastermix"].(execute.ThreadParam).Value.(int)
	}

	vTopUpBuffer, is := m["TopUpBuffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vTopUpBuffer.JSONString), &temp)
		res.TopUpBuffer = *temp.TopUpBuffer
	} else {
		res.TopUpBuffer = m["TopUpBuffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vTotalVolumeperreaction, is := m["TotalVolumeperreaction"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vTotalVolumeperreaction.JSONString), &temp)
		res.TotalVolumeperreaction = *temp.TotalVolumeperreaction
	} else {
		res.TotalVolumeperreaction = m["TotalVolumeperreaction"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vVolumetoLeaveforDNAperreaction, is := m["VolumetoLeaveforDNAperreaction"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Mastermix_reactionsJSONBlock
		json.Unmarshal([]byte(vVolumetoLeaveforDNAperreaction.JSONString), &temp)
		res.VolumetoLeaveforDNAperreaction = *temp.VolumetoLeaveforDNAperreaction
	} else {
		res.VolumetoLeaveforDNAperreaction = m["VolumetoLeaveforDNAperreaction"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["AliquotbyRow"].(execute.ThreadParam).ID
	res.BlockID = m["AliquotbyRow"].(execute.ThreadParam).BlockID

	return res
}

func (e *Mastermix_reactions) OnAliquotbyRow(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AliquotbyRow", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnComponentVolumesperReaction(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ComponentVolumesperReaction", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnComponents(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Components", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnInplate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Inplate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnNumberofMastermixes(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("NumberofMastermixes", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
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
func (e *Mastermix_reactions) OnReactionspermastermix(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Reactionspermastermix", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnTopUpBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TopUpBuffer", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnTotalVolumeperreaction(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TotalVolumeperreaction", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Mastermix_reactions) OnVolumetoLeaveforDNAperreaction(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(10, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VolumetoLeaveforDNAperreaction", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Mastermix_reactions struct {
	flow.Component                 // component "superclass" embedded
	lock                           sync.Mutex
	startup                        sync.Once
	params                         map[execute.ThreadID]*execute.AsyncBag
	AliquotbyRow                   <-chan execute.ThreadParam
	ComponentVolumesperReaction    <-chan execute.ThreadParam
	Components                     <-chan execute.ThreadParam
	Inplate                        <-chan execute.ThreadParam
	NumberofMastermixes            <-chan execute.ThreadParam
	OutPlate                       <-chan execute.ThreadParam
	Reactionspermastermix          <-chan execute.ThreadParam
	TopUpBuffer                    <-chan execute.ThreadParam
	TotalVolumeperreaction         <-chan execute.ThreadParam
	VolumetoLeaveforDNAperreaction <-chan execute.ThreadParam
	Mastermixes                    chan<- execute.ThreadParam
	Status                         chan<- execute.ThreadParam
}

type Mastermix_reactionsParamBlock struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	AliquotbyRow                   bool
	ComponentVolumesperReaction    []wunit.Volume
	Components                     []*wtype.LHComponent
	Inplate                        *wtype.LHPlate
	NumberofMastermixes            int
	OutPlate                       *wtype.LHPlate
	Reactionspermastermix          int
	TopUpBuffer                    *wtype.LHComponent
	TotalVolumeperreaction         wunit.Volume
	VolumetoLeaveforDNAperreaction wunit.Volume
}

type Mastermix_reactionsConfig struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	AliquotbyRow                   bool
	ComponentVolumesperReaction    []wunit.Volume
	Components                     []wtype.FromFactory
	Inplate                        wtype.FromFactory
	NumberofMastermixes            int
	OutPlate                       wtype.FromFactory
	Reactionspermastermix          int
	TopUpBuffer                    wtype.FromFactory
	TotalVolumeperreaction         wunit.Volume
	VolumetoLeaveforDNAperreaction wunit.Volume
}

type Mastermix_reactionsResultBlock struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	Mastermixes []*wtype.LHSolution
	Status      string
}

type Mastermix_reactionsJSONBlock struct {
	ID                             *execute.ThreadID
	BlockID                        *execute.BlockID
	Error                          *bool
	AliquotbyRow                   *bool
	ComponentVolumesperReaction    *[]wunit.Volume
	Components                     *[]*wtype.LHComponent
	Inplate                        **wtype.LHPlate
	NumberofMastermixes            *int
	OutPlate                       **wtype.LHPlate
	Reactionspermastermix          *int
	TopUpBuffer                    **wtype.LHComponent
	TotalVolumeperreaction         *wunit.Volume
	VolumetoLeaveforDNAperreaction *wunit.Volume
	Mastermixes                    *[]*wtype.LHSolution
	Status                         *string
}

func (c *Mastermix_reactions) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AliquotbyRow", "bool", "AliquotbyRow", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ComponentVolumesperReaction", "[]wunit.Volume", "ComponentVolumesperReaction", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Components", "[]*wtype.LHComponent", "Components", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Inplate", "*wtype.LHPlate", "Inplate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("NumberofMastermixes", "int", "NumberofMastermixes", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Reactionspermastermix", "int", "Reactionspermastermix", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TopUpBuffer", "*wtype.LHComponent", "TopUpBuffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TotalVolumeperreaction", "wunit.Volume", "TotalVolumeperreaction", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VolumetoLeaveforDNAperreaction", "wunit.Volume", "VolumetoLeaveforDNAperreaction", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Mastermixes", "[]*wtype.LHSolution", "Mastermixes", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Mastermix_reactions", "Mastermix_reactions", "", false, inp, outp)

	return ci
}
