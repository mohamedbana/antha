package PCR

import (
	"encoding/json"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

/*type Polymerase struct {
	wtype.LHComponent
	Rate_BPpers float64
	Fidelity_errorrate float64 // could dictate how many colonies are checked in validation!
	Extensiontemp Temperature
	Hotstart bool
	StockConcentration Concentration // this is normally in U?
	TargetConcentration Concentration
	// this is also a glycerol solution rather than a watersolution!
}
*/

// Input parameters for this protocol (data)

// PCRprep parameters:

// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume

//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...
//FullTemplatesequence string // better to use Sid's type system here after proof of concept
//FullTemplatelength int	// clearly could be calculated from the sequence... Sid will have a method to do this already so check!
//TargetTemplatesequence string // better to use Sid's type system here after proof of concept
//TargetTemplatelengthinBP int

// Reaction parameters: (could be a entered as thermocycle parameters type possibly?)

//Denaturationtemp Temperature

// Should be calculated from primer and template binding
// should be calculated from template length and polymerase rate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// e.g. DMSO

// Physical outputs from this protocol with types

func (e *PCR) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *PCR) setup(p PCRParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *PCR) steps(p PCRParamBlock, r *PCRResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.SampleForTotalVolume(p.Buffer, p.ReactionVolume)
	samples = append(samples, bufferSample)
	templateSample := mixer.Sample(p.Template, p.Templatevolume)
	samples = append(samples, templateSample)
	dntpSample := mixer.SampleForConcentration(p.DNTPS, p.DNTPconc)
	samples = append(samples, dntpSample)
	FwdPrimerSample := mixer.SampleForConcentration(p.FwdPrimer, p.FwdPrimerConc)
	samples = append(samples, FwdPrimerSample)
	RevPrimerSample := mixer.SampleForConcentration(p.RevPrimer, p.RevPrimerConc)
	samples = append(samples, RevPrimerSample)

	for _, additive := range p.Additives {
		additiveSample := mixer.SampleForConcentration(additive, p.Additiveconc)
		samples = append(samples, additiveSample)
	}

	polySample := mixer.SampleForConcentration(p.PCRPolymerase, p.TargetpolymeraseConcentration)
	samples = append(samples, polySample)
	reaction := _wrapper.MixInto(p.OutPlate, samples...)

	// thermocycle parameters called from enzyme lookup:

	polymerase := p.PCRPolymerase.CName

	extensionTemp := enzymes.DNApolymerasetemps[polymerase]["extensiontemp"]
	meltingTemp := enzymes.DNApolymerasetemps[polymerase]["meltingtemp"]

	// initial Denaturation

	_wrapper.Incubate(reaction, meltingTemp, p.InitDenaturationtime, false)

	for i := 0; i < p.Numberofcycles; i++ {

		// Denature

		_wrapper.Incubate(reaction, meltingTemp, p.Denaturationtime, false)

		// Anneal
		_wrapper.Incubate(reaction, p.AnnealingTemp, p.Annealingtime, false)

		//extensiontime := TargetTemplatelengthinBP/PCRPolymerase.RateBPpers // we'll get type issues here so leave it out for now

		// Extend
		_wrapper.Incubate(reaction, extensionTemp, p.Extensiontime, false)

	}
	// Final Extension
	_wrapper.Incubate(reaction, extensionTemp, p.Finalextensiontime, false)

	// all done
	r.Reaction = reaction
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *PCR) analysis(p PCRParamBlock, r *PCRResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *PCR) validation(p PCRParamBlock, r *PCRResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *PCR) Complete(params interface{}) {
	p := params.(PCRParamBlock)
	if p.Error {
		e.Reaction <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PCRResultBlock)
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
func (e *PCR) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *PCR) NewConfig() interface{} {
	return &PCRConfig{}
}

func (e *PCR) NewParamBlock() interface{} {
	return &PCRParamBlock{}
}

func NewPCR() interface{} { //*PCR {
	e := new(PCR)
	e.init()
	return e
}

// Mapper function
func (e *PCR) Map(m map[string]interface{}) interface{} {
	var res PCRParamBlock
	res.Error = false || m["Additiveconc"].(execute.ThreadParam).Error || m["Additives"].(execute.ThreadParam).Error || m["AnnealingTemp"].(execute.ThreadParam).Error || m["Annealingtime"].(execute.ThreadParam).Error || m["Buffer"].(execute.ThreadParam).Error || m["DNTPS"].(execute.ThreadParam).Error || m["DNTPconc"].(execute.ThreadParam).Error || m["Denaturationtime"].(execute.ThreadParam).Error || m["Extensiontemp"].(execute.ThreadParam).Error || m["Extensiontime"].(execute.ThreadParam).Error || m["Finalextensiontime"].(execute.ThreadParam).Error || m["FwdPrimer"].(execute.ThreadParam).Error || m["FwdPrimerConc"].(execute.ThreadParam).Error || m["InitDenaturationtime"].(execute.ThreadParam).Error || m["Numberofcycles"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["PCRPolymerase"].(execute.ThreadParam).Error || m["ReactionVolume"].(execute.ThreadParam).Error || m["RevPrimer"].(execute.ThreadParam).Error || m["RevPrimerConc"].(execute.ThreadParam).Error || m["TargetpolymeraseConcentration"].(execute.ThreadParam).Error || m["Template"].(execute.ThreadParam).Error || m["Templatevolume"].(execute.ThreadParam).Error

	vAdditiveconc, is := m["Additiveconc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vAdditiveconc.JSONString), &temp)
		res.Additiveconc = *temp.Additiveconc
	} else {
		res.Additiveconc = m["Additiveconc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vAdditives, is := m["Additives"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vAdditives.JSONString), &temp)
		res.Additives = *temp.Additives
	} else {
		res.Additives = m["Additives"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vAnnealingTemp, is := m["AnnealingTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vAnnealingTemp.JSONString), &temp)
		res.AnnealingTemp = *temp.AnnealingTemp
	} else {
		res.AnnealingTemp = m["AnnealingTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vAnnealingtime, is := m["Annealingtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vAnnealingtime.JSONString), &temp)
		res.Annealingtime = *temp.Annealingtime
	} else {
		res.Annealingtime = m["Annealingtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vBuffer, is := m["Buffer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vBuffer.JSONString), &temp)
		res.Buffer = *temp.Buffer
	} else {
		res.Buffer = m["Buffer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDNTPS, is := m["DNTPS"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vDNTPS.JSONString), &temp)
		res.DNTPS = *temp.DNTPS
	} else {
		res.DNTPS = m["DNTPS"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDNTPconc, is := m["DNTPconc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vDNTPconc.JSONString), &temp)
		res.DNTPconc = *temp.DNTPconc
	} else {
		res.DNTPconc = m["DNTPconc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vDenaturationtime, is := m["Denaturationtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vDenaturationtime.JSONString), &temp)
		res.Denaturationtime = *temp.Denaturationtime
	} else {
		res.Denaturationtime = m["Denaturationtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vExtensiontemp, is := m["Extensiontemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vExtensiontemp.JSONString), &temp)
		res.Extensiontemp = *temp.Extensiontemp
	} else {
		res.Extensiontemp = m["Extensiontemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vExtensiontime, is := m["Extensiontime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vExtensiontime.JSONString), &temp)
		res.Extensiontime = *temp.Extensiontime
	} else {
		res.Extensiontime = m["Extensiontime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vFinalextensiontime, is := m["Finalextensiontime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vFinalextensiontime.JSONString), &temp)
		res.Finalextensiontime = *temp.Finalextensiontime
	} else {
		res.Finalextensiontime = m["Finalextensiontime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vFwdPrimer, is := m["FwdPrimer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vFwdPrimer.JSONString), &temp)
		res.FwdPrimer = *temp.FwdPrimer
	} else {
		res.FwdPrimer = m["FwdPrimer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vFwdPrimerConc, is := m["FwdPrimerConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vFwdPrimerConc.JSONString), &temp)
		res.FwdPrimerConc = *temp.FwdPrimerConc
	} else {
		res.FwdPrimerConc = m["FwdPrimerConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vInitDenaturationtime, is := m["InitDenaturationtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vInitDenaturationtime.JSONString), &temp)
		res.InitDenaturationtime = *temp.InitDenaturationtime
	} else {
		res.InitDenaturationtime = m["InitDenaturationtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vNumberofcycles, is := m["Numberofcycles"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vNumberofcycles.JSONString), &temp)
		res.Numberofcycles = *temp.Numberofcycles
	} else {
		res.Numberofcycles = m["Numberofcycles"].(execute.ThreadParam).Value.(int)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPCRPolymerase, is := m["PCRPolymerase"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vPCRPolymerase.JSONString), &temp)
		res.PCRPolymerase = *temp.PCRPolymerase
	} else {
		res.PCRPolymerase = m["PCRPolymerase"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vReactionVolume, is := m["ReactionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vReactionVolume.JSONString), &temp)
		res.ReactionVolume = *temp.ReactionVolume
	} else {
		res.ReactionVolume = m["ReactionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vRevPrimer, is := m["RevPrimer"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vRevPrimer.JSONString), &temp)
		res.RevPrimer = *temp.RevPrimer
	} else {
		res.RevPrimer = m["RevPrimer"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vRevPrimerConc, is := m["RevPrimerConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vRevPrimerConc.JSONString), &temp)
		res.RevPrimerConc = *temp.RevPrimerConc
	} else {
		res.RevPrimerConc = m["RevPrimerConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vTargetpolymeraseConcentration, is := m["TargetpolymeraseConcentration"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vTargetpolymeraseConcentration.JSONString), &temp)
		res.TargetpolymeraseConcentration = *temp.TargetpolymeraseConcentration
	} else {
		res.TargetpolymeraseConcentration = m["TargetpolymeraseConcentration"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vTemplate, is := m["Template"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vTemplate.JSONString), &temp)
		res.Template = *temp.Template
	} else {
		res.Template = m["Template"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vTemplatevolume, is := m["Templatevolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PCRJSONBlock
		json.Unmarshal([]byte(vTemplatevolume.JSONString), &temp)
		res.Templatevolume = *temp.Templatevolume
	} else {
		res.Templatevolume = m["Templatevolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["Additiveconc"].(execute.ThreadParam).ID
	res.BlockID = m["Additiveconc"].(execute.ThreadParam).BlockID

	return res
}

func (e *PCR) OnAdditiveconc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Additiveconc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnAdditives(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Additives", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnAnnealingTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AnnealingTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnAnnealingtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Annealingtime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnBuffer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
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
func (e *PCR) OnDNTPS(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNTPS", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnDNTPconc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNTPconc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnDenaturationtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Denaturationtime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnExtensiontemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Extensiontemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnExtensiontime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Extensiontime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnFinalextensiontime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Finalextensiontime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnFwdPrimer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("FwdPrimer", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnFwdPrimerConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("FwdPrimerConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnInitDenaturationtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("InitDenaturationtime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnNumberofcycles(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Numberofcycles", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
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
func (e *PCR) OnPCRPolymerase(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PCRPolymerase", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnReactionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
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
func (e *PCR) OnRevPrimer(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RevPrimer", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnRevPrimerConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RevPrimerConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnTargetpolymeraseConcentration(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TargetpolymeraseConcentration", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnTemplate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Template", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PCR) OnTemplatevolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Templatevolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type PCR struct {
	flow.Component                // component "superclass" embedded
	lock                          sync.Mutex
	startup                       sync.Once
	params                        map[execute.ThreadID]*execute.AsyncBag
	Additiveconc                  <-chan execute.ThreadParam
	Additives                     <-chan execute.ThreadParam
	AnnealingTemp                 <-chan execute.ThreadParam
	Annealingtime                 <-chan execute.ThreadParam
	Buffer                        <-chan execute.ThreadParam
	DNTPS                         <-chan execute.ThreadParam
	DNTPconc                      <-chan execute.ThreadParam
	Denaturationtime              <-chan execute.ThreadParam
	Extensiontemp                 <-chan execute.ThreadParam
	Extensiontime                 <-chan execute.ThreadParam
	Finalextensiontime            <-chan execute.ThreadParam
	FwdPrimer                     <-chan execute.ThreadParam
	FwdPrimerConc                 <-chan execute.ThreadParam
	InitDenaturationtime          <-chan execute.ThreadParam
	Numberofcycles                <-chan execute.ThreadParam
	OutPlate                      <-chan execute.ThreadParam
	PCRPolymerase                 <-chan execute.ThreadParam
	ReactionVolume                <-chan execute.ThreadParam
	RevPrimer                     <-chan execute.ThreadParam
	RevPrimerConc                 <-chan execute.ThreadParam
	TargetpolymeraseConcentration <-chan execute.ThreadParam
	Template                      <-chan execute.ThreadParam
	Templatevolume                <-chan execute.ThreadParam
	Reaction                      chan<- execute.ThreadParam
}

type PCRParamBlock struct {
	ID                            execute.ThreadID
	BlockID                       execute.BlockID
	Error                         bool
	Additiveconc                  wunit.Concentration
	Additives                     []*wtype.LHComponent
	AnnealingTemp                 wunit.Temperature
	Annealingtime                 wunit.Time
	Buffer                        *wtype.LHComponent
	DNTPS                         *wtype.LHComponent
	DNTPconc                      wunit.Concentration
	Denaturationtime              wunit.Time
	Extensiontemp                 wunit.Temperature
	Extensiontime                 wunit.Time
	Finalextensiontime            wunit.Time
	FwdPrimer                     *wtype.LHComponent
	FwdPrimerConc                 wunit.Concentration
	InitDenaturationtime          wunit.Time
	Numberofcycles                int
	OutPlate                      *wtype.LHPlate
	PCRPolymerase                 *wtype.LHComponent
	ReactionVolume                wunit.Volume
	RevPrimer                     *wtype.LHComponent
	RevPrimerConc                 wunit.Concentration
	TargetpolymeraseConcentration wunit.Concentration
	Template                      *wtype.LHComponent
	Templatevolume                wunit.Volume
}

type PCRConfig struct {
	ID                            execute.ThreadID
	BlockID                       execute.BlockID
	Error                         bool
	Additiveconc                  wunit.Concentration
	Additives                     []wtype.FromFactory
	AnnealingTemp                 wunit.Temperature
	Annealingtime                 wunit.Time
	Buffer                        wtype.FromFactory
	DNTPS                         wtype.FromFactory
	DNTPconc                      wunit.Concentration
	Denaturationtime              wunit.Time
	Extensiontemp                 wunit.Temperature
	Extensiontime                 wunit.Time
	Finalextensiontime            wunit.Time
	FwdPrimer                     wtype.FromFactory
	FwdPrimerConc                 wunit.Concentration
	InitDenaturationtime          wunit.Time
	Numberofcycles                int
	OutPlate                      wtype.FromFactory
	PCRPolymerase                 wtype.FromFactory
	ReactionVolume                wunit.Volume
	RevPrimer                     wtype.FromFactory
	RevPrimerConc                 wunit.Concentration
	TargetpolymeraseConcentration wunit.Concentration
	Template                      wtype.FromFactory
	Templatevolume                wunit.Volume
}

type PCRResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Reaction *wtype.LHSolution
}

type PCRJSONBlock struct {
	ID                            *execute.ThreadID
	BlockID                       *execute.BlockID
	Error                         *bool
	Additiveconc                  *wunit.Concentration
	Additives                     *[]*wtype.LHComponent
	AnnealingTemp                 *wunit.Temperature
	Annealingtime                 *wunit.Time
	Buffer                        **wtype.LHComponent
	DNTPS                         **wtype.LHComponent
	DNTPconc                      *wunit.Concentration
	Denaturationtime              *wunit.Time
	Extensiontemp                 *wunit.Temperature
	Extensiontime                 *wunit.Time
	Finalextensiontime            *wunit.Time
	FwdPrimer                     **wtype.LHComponent
	FwdPrimerConc                 *wunit.Concentration
	InitDenaturationtime          *wunit.Time
	Numberofcycles                *int
	OutPlate                      **wtype.LHPlate
	PCRPolymerase                 **wtype.LHComponent
	ReactionVolume                *wunit.Volume
	RevPrimer                     **wtype.LHComponent
	RevPrimerConc                 *wunit.Concentration
	TargetpolymeraseConcentration *wunit.Concentration
	Template                      **wtype.LHComponent
	Templatevolume                *wunit.Volume
	Reaction                      **wtype.LHSolution
}

func (c *PCR) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Additiveconc", "wunit.Concentration", "Additiveconc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Additives", "[]*wtype.LHComponent", "Additives", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("AnnealingTemp", "wunit.Temperature", "AnnealingTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Annealingtime", "wunit.Time", "Annealingtime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Buffer", "*wtype.LHComponent", "Buffer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNTPS", "*wtype.LHComponent", "DNTPS", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNTPconc", "wunit.Concentration", "DNTPconc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Denaturationtime", "wunit.Time", "Denaturationtime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Extensiontemp", "wunit.Temperature", "Extensiontemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Extensiontime", "wunit.Time", "Extensiontime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Finalextensiontime", "wunit.Time", "Finalextensiontime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("FwdPrimer", "*wtype.LHComponent", "FwdPrimer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("FwdPrimerConc", "wunit.Concentration", "FwdPrimerConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InitDenaturationtime", "wunit.Time", "InitDenaturationtime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Numberofcycles", "int", "Numberofcycles", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PCRPolymerase", "*wtype.LHComponent", "PCRPolymerase", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ReactionVolume", "wunit.Volume", "ReactionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RevPrimer", "*wtype.LHComponent", "RevPrimer", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RevPrimerConc", "wunit.Concentration", "RevPrimerConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TargetpolymeraseConcentration", "wunit.Concentration", "TargetpolymeraseConcentration", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Template", "*wtype.LHComponent", "Template", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Templatevolume", "wunit.Volume", "Templatevolume", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Reaction", "*wtype.LHSolution", "Reaction", true, true, nil, nil))

	ci := execute.NewComponentInfo("PCR", "PCR", "", false, inp, outp)

	return ci
}
