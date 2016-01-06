package Transformation_complete

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

//= 2 (hours)

//Shakerspeed float64 // correct type?

//Plateoutdilution float64

/*ReactionVolume wunit.Volume
PartConc wunit.Concentration
VectorConc wunit.Concentration
AtpVol wunit.Volume
ReVol wunit.Volume
LigVol wunit.Volume
ReactionTemp wunit.Temperature
ReactionTime wunit.Time
InactivationTemp wunit.Temperature
InactivationTime wunit.Time
*/

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Transformation_complete) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Transformation_complete) setup(p Transformation_completeParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Transformation_complete) steps(p Transformation_completeParamBlock, r *Transformation_completeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	competentcells := make([]*wtype.LHComponent, 0)
	competentcells = append(competentcells, p.CompetentCells)
	readycompetentcells := _wrapper.MixInto(p.OutPlate, competentcells...)            // readycompetentcells IS now a LHSolution
	_wrapper.Incubate(readycompetentcells, p.Preplasmidtemp, p.Preplasmidtime, false) // we can incubate an LHSolution so this is fine

	readycompetentcellsComp := wtype.SolutionToComponent(readycompetentcells)

	competetentcellmix := mixer.Sample(readycompetentcellsComp, p.CompetentCellvolumeperassembly) // ERROR! mixer.Sample needs a liquid, not an LHSolution! however, the typeIIs method worked with a *wtype.LHComponent from inputs!
	transformationmix := make([]*wtype.LHComponent, 0)
	transformationmix = append(transformationmix, competetentcellmix)
	DNAsample := mixer.Sample(p.Reaction, p.Reactionvolume)
	transformationmix = append(transformationmix, DNAsample)

	transformedcells := _wrapper.MixInto(p.OutPlate, transformationmix...)

	_wrapper.Incubate(transformedcells, p.Postplasmidtemp, p.Postplasmidtime, false)

	transformedcellsComp := wtype.SolutionToComponent(transformedcells)

	recoverymix := make([]*wtype.LHComponent, 0)
	recoverymixture := mixer.Sample(p.Recoverymedium, p.Recoveryvolume)

	recoverymix = append(recoverymix, transformedcellsComp) // ERROR! transformedcells is now an LHSolution, not a liquid, so can't be used here
	recoverymix = append(recoverymix, recoverymixture)
	recoverymix2 := _wrapper.MixInto(p.OutPlate, recoverymix...)

	_wrapper.Incubate(recoverymix2, p.Recoverytemp, p.Recoverytime, true)

	recoverymix2Comp := wtype.SolutionToComponent(recoverymix2)

	plateout := mixer.Sample(recoverymix2Comp, p.Plateoutvolume) // ERROR! recoverymix2 is now an LHSolution, not a liquid, so can't be used here
	platedculture := _wrapper.MixInto(p.AgarPlate, plateout)

	r.Platedculture = platedculture
	_ = _wrapper.WaitToEnd()

	/*atpSample := mixer.Sample(Atp, AtpVol)
	samples = append(samples, atpSample)
	vectorSample := mixer.SampleForConcentration(Vector, VectorConc)
	samples = append(samples, vectorSample)

	for _, part := range Parts {
		partSample := mixer.SampleForConcentration(part, PartConc)
		samples = append(samples, partSample)
	}

	reSample := mixer.Sample(RestrictionEnzyme, ReVol)
	samples = append(samples, reSample)
	ligSample := mixer.Sample(Ligase, LigVol)
	samples = append(samples, ligSample)


	// incubate the reaction mixture

	Incubate(reaction, ReactionTemp, ReactionTime, false)

	// inactivate

	Incubate(reaction, InactivationTemp, InactivationTime, false)

	// all done
	Reaction = reaction

	readycompetentcells := Incubate (CompetentCells,Preplasmidtemp, Preplasmidtime, false)


	product := Mix (Reaction(ReactionVolume), readycompetentcells(CompetentCellvolumeperassembly))
	transformedcells := Incubate (product, Postplasmidtime,Postplasmidtemp,false)
	recoverymixture := Mix (transformedcells, Recoverymedium (Recoveryvolume)) // or alternative recovery medium
	Incubate (recoverymixture, Recoverytime, Recoverytemp, Shakerspeed)
	platedculture := MixInto(AgarPlate, Plateoutvolume)

	Platedculture = platedculture

	*/
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Transformation_complete) analysis(p Transformation_completeParamBlock, r *Transformation_completeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Transformation_complete) validation(p Transformation_completeParamBlock, r *Transformation_completeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Transformation_complete) Complete(params interface{}) {
	p := params.(Transformation_completeParamBlock)
	if p.Error {
		e.Platedculture <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(Transformation_completeResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Platedculture <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Platedculture <- execute.ThreadParam{Value: r.Platedculture, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Transformation_complete) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Transformation_complete) NewConfig() interface{} {
	return &Transformation_completeConfig{}
}

func (e *Transformation_complete) NewParamBlock() interface{} {
	return &Transformation_completeParamBlock{}
}

func NewTransformation_complete() interface{} { //*Transformation_complete {
	e := new(Transformation_complete)
	e.init()
	return e
}

// Mapper function
func (e *Transformation_complete) Map(m map[string]interface{}) interface{} {
	var res Transformation_completeParamBlock
	res.Error = false || m["AgarPlate"].(execute.ThreadParam).Error || m["CompetentCells"].(execute.ThreadParam).Error || m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Plateoutvolume"].(execute.ThreadParam).Error || m["Postplasmidtemp"].(execute.ThreadParam).Error || m["Postplasmidtime"].(execute.ThreadParam).Error || m["Preplasmidtemp"].(execute.ThreadParam).Error || m["Preplasmidtime"].(execute.ThreadParam).Error || m["Reaction"].(execute.ThreadParam).Error || m["Reactionvolume"].(execute.ThreadParam).Error || m["Recoverymedium"].(execute.ThreadParam).Error || m["Recoverytemp"].(execute.ThreadParam).Error || m["Recoverytime"].(execute.ThreadParam).Error || m["Recoveryvolume"].(execute.ThreadParam).Error

	vAgarPlate, is := m["AgarPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vAgarPlate.JSONString), &temp)
		res.AgarPlate = *temp.AgarPlate
	} else {
		res.AgarPlate = m["AgarPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vCompetentCells, is := m["CompetentCells"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vCompetentCells.JSONString), &temp)
		res.CompetentCells = *temp.CompetentCells
	} else {
		res.CompetentCells = m["CompetentCells"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vCompetentCellvolumeperassembly, is := m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vCompetentCellvolumeperassembly.JSONString), &temp)
		res.CompetentCellvolumeperassembly = *temp.CompetentCellvolumeperassembly
	} else {
		res.CompetentCellvolumeperassembly = m["CompetentCellvolumeperassembly"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPlateoutvolume, is := m["Plateoutvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vPlateoutvolume.JSONString), &temp)
		res.Plateoutvolume = *temp.Plateoutvolume
	} else {
		res.Plateoutvolume = m["Plateoutvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vPostplasmidtemp, is := m["Postplasmidtemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vPostplasmidtemp.JSONString), &temp)
		res.Postplasmidtemp = *temp.Postplasmidtemp
	} else {
		res.Postplasmidtemp = m["Postplasmidtemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vPostplasmidtime, is := m["Postplasmidtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vPostplasmidtime.JSONString), &temp)
		res.Postplasmidtime = *temp.Postplasmidtime
	} else {
		res.Postplasmidtime = m["Postplasmidtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vPreplasmidtemp, is := m["Preplasmidtemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vPreplasmidtemp.JSONString), &temp)
		res.Preplasmidtemp = *temp.Preplasmidtemp
	} else {
		res.Preplasmidtemp = m["Preplasmidtemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vPreplasmidtime, is := m["Preplasmidtime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vPreplasmidtime.JSONString), &temp)
		res.Preplasmidtime = *temp.Preplasmidtime
	} else {
		res.Preplasmidtime = m["Preplasmidtime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vReaction, is := m["Reaction"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vReaction.JSONString), &temp)
		res.Reaction = *temp.Reaction
	} else {
		res.Reaction = m["Reaction"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vReactionvolume, is := m["Reactionvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vReactionvolume.JSONString), &temp)
		res.Reactionvolume = *temp.Reactionvolume
	} else {
		res.Reactionvolume = m["Reactionvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vRecoverymedium, is := m["Recoverymedium"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vRecoverymedium.JSONString), &temp)
		res.Recoverymedium = *temp.Recoverymedium
	} else {
		res.Recoverymedium = m["Recoverymedium"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vRecoverytemp, is := m["Recoverytemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vRecoverytemp.JSONString), &temp)
		res.Recoverytemp = *temp.Recoverytemp
	} else {
		res.Recoverytemp = m["Recoverytemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vRecoverytime, is := m["Recoverytime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vRecoverytime.JSONString), &temp)
		res.Recoverytime = *temp.Recoverytime
	} else {
		res.Recoverytime = m["Recoverytime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vRecoveryvolume, is := m["Recoveryvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Transformation_completeJSONBlock
		json.Unmarshal([]byte(vRecoveryvolume.JSONString), &temp)
		res.Recoveryvolume = *temp.Recoveryvolume
	} else {
		res.Recoveryvolume = m["Recoveryvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["AgarPlate"].(execute.ThreadParam).ID
	res.BlockID = m["AgarPlate"].(execute.ThreadParam).BlockID

	return res
}

func (e *Transformation_complete) OnAgarPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AgarPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnCompetentCells(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("CompetentCells", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnCompetentCellvolumeperassembly(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
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
func (e *Transformation_complete) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
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
func (e *Transformation_complete) OnPlateoutvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Plateoutvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnPostplasmidtemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
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
func (e *Transformation_complete) OnPostplasmidtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
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
func (e *Transformation_complete) OnPreplasmidtemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Preplasmidtemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnPreplasmidtime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Preplasmidtime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnReaction(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
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
func (e *Transformation_complete) OnReactionvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
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
func (e *Transformation_complete) OnRecoverymedium(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoverymedium", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnRecoverytemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoverytemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnRecoverytime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoverytime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transformation_complete) OnRecoveryvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(15, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Recoveryvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Transformation_complete struct {
	flow.Component                 // component "superclass" embedded
	lock                           sync.Mutex
	startup                        sync.Once
	params                         map[execute.ThreadID]*execute.AsyncBag
	AgarPlate                      <-chan execute.ThreadParam
	CompetentCells                 <-chan execute.ThreadParam
	CompetentCellvolumeperassembly <-chan execute.ThreadParam
	OutPlate                       <-chan execute.ThreadParam
	Plateoutvolume                 <-chan execute.ThreadParam
	Postplasmidtemp                <-chan execute.ThreadParam
	Postplasmidtime                <-chan execute.ThreadParam
	Preplasmidtemp                 <-chan execute.ThreadParam
	Preplasmidtime                 <-chan execute.ThreadParam
	Reaction                       <-chan execute.ThreadParam
	Reactionvolume                 <-chan execute.ThreadParam
	Recoverymedium                 <-chan execute.ThreadParam
	Recoverytemp                   <-chan execute.ThreadParam
	Recoverytime                   <-chan execute.ThreadParam
	Recoveryvolume                 <-chan execute.ThreadParam
	Platedculture                  chan<- execute.ThreadParam
}

type Transformation_completeParamBlock struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	AgarPlate                      *wtype.LHPlate
	CompetentCells                 *wtype.LHComponent
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       *wtype.LHPlate
	Plateoutvolume                 wunit.Volume
	Postplasmidtemp                wunit.Temperature
	Postplasmidtime                wunit.Time
	Preplasmidtemp                 wunit.Temperature
	Preplasmidtime                 wunit.Time
	Reaction                       *wtype.LHComponent
	Reactionvolume                 wunit.Volume
	Recoverymedium                 *wtype.LHComponent
	Recoverytemp                   wunit.Temperature
	Recoverytime                   wunit.Time
	Recoveryvolume                 wunit.Volume
}

type Transformation_completeConfig struct {
	ID                             execute.ThreadID
	BlockID                        execute.BlockID
	Error                          bool
	AgarPlate                      wtype.FromFactory
	CompetentCells                 wtype.FromFactory
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       wtype.FromFactory
	Plateoutvolume                 wunit.Volume
	Postplasmidtemp                wunit.Temperature
	Postplasmidtime                wunit.Time
	Preplasmidtemp                 wunit.Temperature
	Preplasmidtime                 wunit.Time
	Reaction                       wtype.FromFactory
	Reactionvolume                 wunit.Volume
	Recoverymedium                 wtype.FromFactory
	Recoverytemp                   wunit.Temperature
	Recoverytime                   wunit.Time
	Recoveryvolume                 wunit.Volume
}

type Transformation_completeResultBlock struct {
	ID            execute.ThreadID
	BlockID       execute.BlockID
	Error         bool
	Platedculture *wtype.LHSolution
}

type Transformation_completeJSONBlock struct {
	ID                             *execute.ThreadID
	BlockID                        *execute.BlockID
	Error                          *bool
	AgarPlate                      **wtype.LHPlate
	CompetentCells                 **wtype.LHComponent
	CompetentCellvolumeperassembly *wunit.Volume
	OutPlate                       **wtype.LHPlate
	Plateoutvolume                 *wunit.Volume
	Postplasmidtemp                *wunit.Temperature
	Postplasmidtime                *wunit.Time
	Preplasmidtemp                 *wunit.Temperature
	Preplasmidtime                 *wunit.Time
	Reaction                       **wtype.LHComponent
	Reactionvolume                 *wunit.Volume
	Recoverymedium                 **wtype.LHComponent
	Recoverytemp                   *wunit.Temperature
	Recoverytime                   *wunit.Time
	Recoveryvolume                 *wunit.Volume
	Platedculture                  **wtype.LHSolution
}

func (c *Transformation_complete) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AgarPlate", "*wtype.LHPlate", "AgarPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("CompetentCells", "*wtype.LHComponent", "CompetentCells", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("CompetentCellvolumeperassembly", "wunit.Volume", "CompetentCellvolumeperassembly", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Plateoutvolume", "wunit.Volume", "Plateoutvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Postplasmidtemp", "wunit.Temperature", "Postplasmidtemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Postplasmidtime", "wunit.Time", "Postplasmidtime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Preplasmidtemp", "wunit.Temperature", "Preplasmidtemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Preplasmidtime", "wunit.Time", "Preplasmidtime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Reaction", "*wtype.LHComponent", "Reaction", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Reactionvolume", "wunit.Volume", "Reactionvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoverymedium", "*wtype.LHComponent", "Recoverymedium", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoverytemp", "wunit.Temperature", "Recoverytemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoverytime", "wunit.Time", "Recoverytime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Recoveryvolume", "wunit.Volume", "Recoveryvolume", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Platedculture", "*wtype.LHSolution", "Platedculture", true, true, nil, nil))

	ci := execute.NewComponentInfo("Transformation_complete", "Transformation_complete", "", false, inp, outp)

	return ci
}
