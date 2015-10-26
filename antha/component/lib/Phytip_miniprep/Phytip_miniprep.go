//
package Phytip_miniprep

import (
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/unitoperations"
	//"github.com/antha-lang/antha/antha/anthalib/mixer"
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
	"time"
)

//Torr

// cubesensor streams:
// in pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa
// in Kelvin
// Percentage // density water vapor (kg/m3)

//float64// time

/* RBvolume Volume // 150ul
RBflowrate Rate
RBpause Time // seconds
RBcycles int

LBvolume Volume
LBflowrate Rate
LBpause Time
LBcycles int

PBvolume Volume
PBflowrate Rate
PBpause Time
PBcycles int

Equilibrationvolume Volume
Equilibrationflowrate Rate
Equilibrationpause Time
Equilibrationcycles int

Airdispensevolume Volume
Airdispenseflowrate Rate
Airdispensepause Time
Airdispensecycles int



Airaspiratevolume Volume
Airaspirateflowrate Rate
Airaspiratepause Time
Airaspiratecylces int

Capturevoume Volume
Captureflowrate Rate
Capturepause Time
Capturecycles int

Washbuffervolume [] Volume
Washbufferflowrate [] Rate
Washbufferpause [] Time
Washbuffercycles [] int



Elutionbuffervolume Volume
Elutionflowrate Rate
Elutionpause Time
Elutioncycles int

*/
//or

/* PlasmidConc Concentration
Storagelocation Location
Storageconditions StorageHistory
Plasmidbuffer Composition */ // is this all inferred from a PLasmid solution  type anyway?

//
// wtype.LHTip
//Devices.Chiller
//Devices.Vacuum_manifold
//Devices.Blotting_position

//unitoperations.Pellet // wrong type?

//RB wtype.LHComponent //Watersolution
//LB wtype.LHComponent //Watersolution
//PB wtype.LHComponent //Watersolution
//Watersolution // equilibration buffer
//Gas
//Washbuffer []wtype.LHComponent //Watersolution
//Elutionbuffer *wtype.LHComponent //Watersolution

//Solution //PlasmidSolution

func (e *Phytip_miniprep) requirements() {
	_ = wunit.Make_units

}
func (e *Phytip_miniprep) setup(p Phytip_miniprepParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *Phytip_miniprep) steps(p Phytip_miniprepParamBlock, r *Phytip_miniprepResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper

	resuspension, _ := unitoperations.Resuspend(p.Cellpellet, p.Resuspensionstep, p.Tips)
	lysate, _ := unitoperations.Chromatography(resuspension, p.Lysisstep, p.Tips)
	precipitate, _ := unitoperations.Chromatography(lysate, p.Precipitationstep, p.Tips)

	_, columnready := unitoperations.Chromatography(p.Water, p.Equilibrationstep, p.Phytips)

	_, readyforcapture := unitoperations.Chromatography(p.Air, p.Airstep, columnready)
	capture, readyforcapture := unitoperations.Chromatography(precipitate, p.Capturestep, readyforcapture)

	for _, washstep := range p.Washsteps {
		_, readyforcapture = unitoperations.Chromatography(capture, washstep, readyforcapture)
	}
	readyfordrying := unitoperations.Blot(readyforcapture, p.Blotcycles, p.Blottime)

	/*if Vacuum == true {
		drytips := unitoperations.Dry(Tips,Drytime,Vacuumstrength)


		//parameters required for evaporation calculator
		Liquid := Washsteps[0].Pipetstep.Name //ethanol?
		// lookup properties via liquidclasses package to workout evaporation time using Evaporationrate element?


		//Platetype := Phytips.tip //.surfacearea? labware.phytip.surfacearea?
		Volumeperwell := (Washsteps[0].Pipetstep.Volume.SIValue() / 10) // assume max 10% residual volume for now??

		drytimerequired := Evaporation.Estimatedevaporationtime(Airvelocity, Liquid, Platetype, Volumeperwell)


	} else {*/drytips := unitoperations.Dry(readyfordrying, p.Drytime, p.Vacuumstrength) //}

	r.PlasmidDNAsolution, _ = unitoperations.Chromatography(p.Elutionstep.Buffer, p.Elutionstep, drytips)
	_ = _wrapper.WaitToEnd()

}
func (e *Phytip_miniprep) analysis(p Phytip_miniprepParamBlock, r *Phytip_miniprepResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *Phytip_miniprep) validation(p Phytip_miniprepParamBlock, r *Phytip_miniprepResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Phytip_miniprep) Complete(params interface{}) {
	p := params.(Phytip_miniprepParamBlock)
	if p.Error {
		e.PlasmidDNAsolution <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(Phytip_miniprepResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.PlasmidDNAsolution <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.PlasmidDNAsolution <- execute.ThreadParam{Value: r.PlasmidDNAsolution, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Phytip_miniprep) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Phytip_miniprep) NewConfig() interface{} {
	return &Phytip_miniprepConfig{}
}

func (e *Phytip_miniprep) NewParamBlock() interface{} {
	return &Phytip_miniprepParamBlock{}
}

func NewPhytip_miniprep() interface{} { //*Phytip_miniprep {
	e := new(Phytip_miniprep)
	e.init()
	return e
}

// Mapper function
func (e *Phytip_miniprep) Map(m map[string]interface{}) interface{} {
	var res Phytip_miniprepParamBlock
	res.Error = false || m["Air"].(execute.ThreadParam).Error || m["Airstep"].(execute.ThreadParam).Error || m["Blotcycles"].(execute.ThreadParam).Error || m["Blottime"].(execute.ThreadParam).Error || m["Capturestep"].(execute.ThreadParam).Error || m["Cellpellet"].(execute.ThreadParam).Error || m["Cellpelletmass"].(execute.ThreadParam).Error || m["Drytime"].(execute.ThreadParam).Error || m["Elutionstep"].(execute.ThreadParam).Error || m["Equilibrationstep"].(execute.ThreadParam).Error || m["Lysisstep"].(execute.ThreadParam).Error || m["Pa"].(execute.ThreadParam).Error || m["Phytips"].(execute.ThreadParam).Error || m["Precipitationstep"].(execute.ThreadParam).Error || m["Relativehumidity"].(execute.ThreadParam).Error || m["Resuspensionstep"].(execute.ThreadParam).Error || m["Temp"].(execute.ThreadParam).Error || m["Time"].(execute.ThreadParam).Error || m["Tips"].(execute.ThreadParam).Error || m["Vacuum"].(execute.ThreadParam).Error || m["Vacuumstrength"].(execute.ThreadParam).Error || m["Washsteps"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vAir, is := m["Air"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vAir.JSONString), &temp)
		res.Air = *temp.Air
	} else {
		res.Air = m["Air"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vAirstep, is := m["Airstep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vAirstep.JSONString), &temp)
		res.Airstep = *temp.Airstep
	} else {
		res.Airstep = m["Airstep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vBlotcycles, is := m["Blotcycles"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vBlotcycles.JSONString), &temp)
		res.Blotcycles = *temp.Blotcycles
	} else {
		res.Blotcycles = m["Blotcycles"].(execute.ThreadParam).Value.(int)
	}

	vBlottime, is := m["Blottime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vBlottime.JSONString), &temp)
		res.Blottime = *temp.Blottime
	} else {
		res.Blottime = m["Blottime"].(execute.ThreadParam).Value.(time.Duration)
	}

	vCapturestep, is := m["Capturestep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vCapturestep.JSONString), &temp)
		res.Capturestep = *temp.Capturestep
	} else {
		res.Capturestep = m["Capturestep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vCellpellet, is := m["Cellpellet"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vCellpellet.JSONString), &temp)
		res.Cellpellet = *temp.Cellpellet
	} else {
		res.Cellpellet = m["Cellpellet"].(execute.ThreadParam).Value.(*wtype.Physical)
	}

	vCellpelletmass, is := m["Cellpelletmass"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vCellpelletmass.JSONString), &temp)
		res.Cellpelletmass = *temp.Cellpelletmass
	} else {
		res.Cellpelletmass = m["Cellpelletmass"].(execute.ThreadParam).Value.(wunit.Mass)
	}

	vDrytime, is := m["Drytime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vDrytime.JSONString), &temp)
		res.Drytime = *temp.Drytime
	} else {
		res.Drytime = m["Drytime"].(execute.ThreadParam).Value.(time.Duration)
	}

	vElutionstep, is := m["Elutionstep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vElutionstep.JSONString), &temp)
		res.Elutionstep = *temp.Elutionstep
	} else {
		res.Elutionstep = m["Elutionstep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vEquilibrationstep, is := m["Equilibrationstep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vEquilibrationstep.JSONString), &temp)
		res.Equilibrationstep = *temp.Equilibrationstep
	} else {
		res.Equilibrationstep = m["Equilibrationstep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vLysisstep, is := m["Lysisstep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vLysisstep.JSONString), &temp)
		res.Lysisstep = *temp.Lysisstep
	} else {
		res.Lysisstep = m["Lysisstep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vPa, is := m["Pa"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vPa.JSONString), &temp)
		res.Pa = *temp.Pa
	} else {
		res.Pa = m["Pa"].(execute.ThreadParam).Value.(float64)
	}

	vPhytips, is := m["Phytips"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vPhytips.JSONString), &temp)
		res.Phytips = *temp.Phytips
	} else {
		res.Phytips = m["Phytips"].(execute.ThreadParam).Value.(unitoperations.Column)
	}

	vPrecipitationstep, is := m["Precipitationstep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vPrecipitationstep.JSONString), &temp)
		res.Precipitationstep = *temp.Precipitationstep
	} else {
		res.Precipitationstep = m["Precipitationstep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vRelativehumidity, is := m["Relativehumidity"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vRelativehumidity.JSONString), &temp)
		res.Relativehumidity = *temp.Relativehumidity
	} else {
		res.Relativehumidity = m["Relativehumidity"].(execute.ThreadParam).Value.(float64)
	}

	vResuspensionstep, is := m["Resuspensionstep"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vResuspensionstep.JSONString), &temp)
		res.Resuspensionstep = *temp.Resuspensionstep
	} else {
		res.Resuspensionstep = m["Resuspensionstep"].(execute.ThreadParam).Value.(unitoperations.Chromstep)
	}

	vTemp, is := m["Temp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vTemp.JSONString), &temp)
		res.Temp = *temp.Temp
	} else {
		res.Temp = m["Temp"].(execute.ThreadParam).Value.(float64)
	}

	vTime, is := m["Time"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vTime.JSONString), &temp)
		res.Time = *temp.Time
	} else {
		res.Time = m["Time"].(execute.ThreadParam).Value.(time.Duration)
	}

	vTips, is := m["Tips"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vTips.JSONString), &temp)
		res.Tips = *temp.Tips
	} else {
		res.Tips = m["Tips"].(execute.ThreadParam).Value.(unitoperations.Column)
	}

	vVacuum, is := m["Vacuum"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vVacuum.JSONString), &temp)
		res.Vacuum = *temp.Vacuum
	} else {
		res.Vacuum = m["Vacuum"].(execute.ThreadParam).Value.(bool)
	}

	vVacuumstrength, is := m["Vacuumstrength"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vVacuumstrength.JSONString), &temp)
		res.Vacuumstrength = *temp.Vacuumstrength
	} else {
		res.Vacuumstrength = m["Vacuumstrength"].(execute.ThreadParam).Value.(float64)
	}

	vWashsteps, is := m["Washsteps"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vWashsteps.JSONString), &temp)
		res.Washsteps = *temp.Washsteps
	} else {
		res.Washsteps = m["Washsteps"].(execute.ThreadParam).Value.([]unitoperations.Chromstep)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Phytip_miniprepJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["Air"].(execute.ThreadParam).ID
	res.BlockID = m["Air"].(execute.ThreadParam).BlockID

	return res
}

func (e *Phytip_miniprep) OnAir(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Air", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnAirstep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Airstep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnBlotcycles(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Blotcycles", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnBlottime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Blottime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnCapturestep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Capturestep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnCellpellet(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Cellpellet", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnCellpelletmass(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Cellpelletmass", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnDrytime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Drytime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnElutionstep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Elutionstep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnEquilibrationstep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Equilibrationstep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnLysisstep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Lysisstep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnPa(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Pa", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnPhytips(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Phytips", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnPrecipitationstep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Precipitationstep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnRelativehumidity(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Relativehumidity", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnResuspensionstep(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Resuspensionstep", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Temp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnTime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Time", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnTips(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Tips", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnVacuum(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vacuum", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnVacuumstrength(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vacuumstrength", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnWashsteps(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Washsteps", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Phytip_miniprep) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(23, e, e)
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

type Phytip_miniprep struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	Air                <-chan execute.ThreadParam
	Airstep            <-chan execute.ThreadParam
	Blotcycles         <-chan execute.ThreadParam
	Blottime           <-chan execute.ThreadParam
	Capturestep        <-chan execute.ThreadParam
	Cellpellet         <-chan execute.ThreadParam
	Cellpelletmass     <-chan execute.ThreadParam
	Drytime            <-chan execute.ThreadParam
	Elutionstep        <-chan execute.ThreadParam
	Equilibrationstep  <-chan execute.ThreadParam
	Lysisstep          <-chan execute.ThreadParam
	Pa                 <-chan execute.ThreadParam
	Phytips            <-chan execute.ThreadParam
	Precipitationstep  <-chan execute.ThreadParam
	Relativehumidity   <-chan execute.ThreadParam
	Resuspensionstep   <-chan execute.ThreadParam
	Temp               <-chan execute.ThreadParam
	Time               <-chan execute.ThreadParam
	Tips               <-chan execute.ThreadParam
	Vacuum             <-chan execute.ThreadParam
	Vacuumstrength     <-chan execute.ThreadParam
	Washsteps          <-chan execute.ThreadParam
	Water              <-chan execute.ThreadParam
	PlasmidDNAsolution chan<- execute.ThreadParam
}

type Phytip_miniprepParamBlock struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Air               *wtype.LHComponent
	Airstep           unitoperations.Chromstep
	Blotcycles        int
	Blottime          time.Duration
	Capturestep       unitoperations.Chromstep
	Cellpellet        *wtype.Physical
	Cellpelletmass    wunit.Mass
	Drytime           time.Duration
	Elutionstep       unitoperations.Chromstep
	Equilibrationstep unitoperations.Chromstep
	Lysisstep         unitoperations.Chromstep
	Pa                float64
	Phytips           unitoperations.Column
	Precipitationstep unitoperations.Chromstep
	Relativehumidity  float64
	Resuspensionstep  unitoperations.Chromstep
	Temp              float64
	Time              time.Duration
	Tips              unitoperations.Column
	Vacuum            bool
	Vacuumstrength    float64
	Washsteps         []unitoperations.Chromstep
	Water             *wtype.LHComponent
}

type Phytip_miniprepConfig struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Air               wtype.FromFactory
	Airstep           unitoperations.Chromstep
	Blotcycles        int
	Blottime          time.Duration
	Capturestep       unitoperations.Chromstep
	Cellpellet        wtype.FromFactory
	Cellpelletmass    wunit.Mass
	Drytime           time.Duration
	Elutionstep       unitoperations.Chromstep
	Equilibrationstep unitoperations.Chromstep
	Lysisstep         unitoperations.Chromstep
	Pa                float64
	Phytips           unitoperations.Column
	Precipitationstep unitoperations.Chromstep
	Relativehumidity  float64
	Resuspensionstep  unitoperations.Chromstep
	Temp              float64
	Time              time.Duration
	Tips              unitoperations.Column
	Vacuum            bool
	Vacuumstrength    float64
	Washsteps         []unitoperations.Chromstep
	Water             wtype.FromFactory
}

type Phytip_miniprepResultBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	PlasmidDNAsolution *wtype.LHComponent
}

type Phytip_miniprepJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	Air                **wtype.LHComponent
	Airstep            *unitoperations.Chromstep
	Blotcycles         *int
	Blottime           *time.Duration
	Capturestep        *unitoperations.Chromstep
	Cellpellet         **wtype.Physical
	Cellpelletmass     *wunit.Mass
	Drytime            *time.Duration
	Elutionstep        *unitoperations.Chromstep
	Equilibrationstep  *unitoperations.Chromstep
	Lysisstep          *unitoperations.Chromstep
	Pa                 *float64
	Phytips            *unitoperations.Column
	Precipitationstep  *unitoperations.Chromstep
	Relativehumidity   *float64
	Resuspensionstep   *unitoperations.Chromstep
	Temp               *float64
	Time               *time.Duration
	Tips               *unitoperations.Column
	Vacuum             *bool
	Vacuumstrength     *float64
	Washsteps          *[]unitoperations.Chromstep
	Water              **wtype.LHComponent
	PlasmidDNAsolution **wtype.LHComponent
}

func (c *Phytip_miniprep) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Air", "*wtype.LHComponent", "Air", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Airstep", "unitoperations.Chromstep", "Airstep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Blotcycles", "int", "Blotcycles", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Blottime", "time.Duration", "Blottime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Capturestep", "unitoperations.Chromstep", "Capturestep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Cellpellet", "*wtype.Physical", "Cellpellet", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Cellpelletmass", "wunit.Mass", "Cellpelletmass", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Drytime", "time.Duration", "Drytime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Elutionstep", "unitoperations.Chromstep", "Elutionstep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Equilibrationstep", "unitoperations.Chromstep", "Equilibrationstep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Lysisstep", "unitoperations.Chromstep", "Lysisstep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Pa", "float64", "Pa", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Phytips", "unitoperations.Column", "Phytips", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Precipitationstep", "unitoperations.Chromstep", "Precipitationstep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Relativehumidity", "float64", "Relativehumidity", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Resuspensionstep", "unitoperations.Chromstep", "Resuspensionstep", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Temp", "float64", "Temp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Time", "time.Duration", "Time", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Tips", "unitoperations.Column", "Tips", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vacuum", "bool", "Vacuum", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vacuumstrength", "float64", "Vacuumstrength", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Washsteps", "[]unitoperations.Chromstep", "Washsteps", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("PlasmidDNAsolution", "*wtype.LHComponent", "PlasmidDNAsolution", true, true, nil, nil))

	ci := execute.NewComponentInfo("Phytip_miniprep", "Phytip_miniprep", "", false, inp, outp)

	return ci
}
