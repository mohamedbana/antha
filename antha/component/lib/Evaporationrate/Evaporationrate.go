// antha/component/lib/Evaporationrate/Evaporationrate.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

/* Evaporation calculator based on
http://www.engineeringtoolbox.com/evaporation-water-surface-d_690.html

This engineering function may need to be improved to account for vapour pressure and surface tension

gs = Θ A (xs - x) / 3600         (1)

or

gh = Θ A (xs - x)

where

gs = amount of evaporated water per second (kg/s)

gh = amount of evaporated water per hour (kg/h)

Θ = (25 + 19 v) = evaporation coefficient (kg/m2h)

v = velocity of air above the water surface (m/s)

A = water surface area (m2)

xs = humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)

x = humidity ratio in the air (kg/kg) (kg H2O in kg Dry Air) */

package Evaporationrate

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

// ul

// cubesensor streams:
//wunit.Pressure // in pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa
// input in deg C will be converted to Kelvin
// Percentage // density water vapor (kg/m3)

// // velocity of air above water in m/s ; could be calculated or measured

// time

// ul/h
// ul

func (e *Evaporationrate) requirements() {
	_ = wunit.Make_units

}
func (e *Evaporationrate) setup(p EvaporationrateParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *Evaporationrate) steps(p EvaporationrateParamBlock, r *EvaporationrateResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *Evaporationrate) analysis(p EvaporationrateParamBlock, r *EvaporationrateResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper

	tempinKelvin := (p.Temp.SIValue() + 273.15)
	var PWS float64 = eng.Pws(tempinKelvin)
	var pw float64 = eng.Pw(p.Relativehumidity, PWS) // vapour partial pressure in Pascals
	var Gh = (eng.Θ(p.Liquid, p.Airvelocity) *
		(labware.Labwaregeometry[p.Platetype]["Surfacearea"] *
			((eng.Xs(PWS, p.Pa)) - (eng.X(pw, p.Pa))))) // Gh is rate of evaporation in kg/h
	evaporatedliquid := (Gh * (p.Executiontime.SIValue() / 3600))                            // in kg
	evaporatedliquid = (evaporatedliquid * liquidclasses.Liquidclass[p.Liquid]["ro"]) / 1000 // converted to litres
	r.Evaporatedliquid = wunit.NewVolume((evaporatedliquid * 1000000), "ul")                 // convert to ul

	r.Evaporationrateestimate = Gh * 1000000 // ul/h if declared in parameters or data it doesn't need declaring again

	estimatedevaporationtime := p.Volumeperwell.ConvertTo(wunit.ParsePrefixedUnit("ul")) / r.Evaporationrateestimate
	r.Estimatedevaporationtime = wunit.NewTime((estimatedevaporationtime * 3600), "s")

	r.Status = fmt.Sprintln("Well Surface Area=",
		(labware.Labwaregeometry[p.Platetype]["Surfacearea"])*1000000, "mm2",
		"evaporation rate =", Gh*1000000, "ul/h",
		"total evaporated liquid =", r.Evaporatedliquid.ToString(), "after", p.Executiontime.ToString(),
		"estimated evaporation time = ", r.Estimatedevaporationtime.ToString())
	_ = _wrapper.WaitToEnd()

} // works in either analysis or steps sections

func (e *Evaporationrate) validation(p EvaporationrateParamBlock, r *EvaporationrateResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper

	if r.Evaporatedliquid.SIValue() > p.Volumeperwell.SIValue() {
		panic("not enough liquid")
	}
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Evaporationrate) Complete(params interface{}) {
	p := params.(EvaporationrateParamBlock)
	if p.Error {
		e.Estimatedevaporationtime <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Evaporatedliquid <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Evaporationrateestimate <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(EvaporationrateResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Estimatedevaporationtime <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Evaporatedliquid <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Evaporationrateestimate <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.analysis(p, r)

	e.Estimatedevaporationtime <- execute.ThreadParam{Value: r.Estimatedevaporationtime, ID: p.ID, Error: false}

	e.Evaporationrateestimate <- execute.ThreadParam{Value: r.Evaporationrateestimate, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.validation(p, r)
	e.Evaporatedliquid <- execute.ThreadParam{Value: r.Evaporatedliquid, ID: p.ID, Error: false}

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Evaporationrate) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Evaporationrate) NewConfig() interface{} {
	return &EvaporationrateConfig{}
}

func (e *Evaporationrate) NewParamBlock() interface{} {
	return &EvaporationrateParamBlock{}
}

func NewEvaporationrate() interface{} { //*Evaporationrate {
	e := new(Evaporationrate)
	e.init()
	return e
}

// Mapper function
func (e *Evaporationrate) Map(m map[string]interface{}) interface{} {
	var res EvaporationrateParamBlock
	res.Error = false || m["Airvelocity"].(execute.ThreadParam).Error || m["Executiontime"].(execute.ThreadParam).Error || m["Liquid"].(execute.ThreadParam).Error || m["Pa"].(execute.ThreadParam).Error || m["Platetype"].(execute.ThreadParam).Error || m["Relativehumidity"].(execute.ThreadParam).Error || m["Temp"].(execute.ThreadParam).Error || m["Volumeperwell"].(execute.ThreadParam).Error

	vAirvelocity, is := m["Airvelocity"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vAirvelocity.JSONString), &temp)
		res.Airvelocity = *temp.Airvelocity
	} else {
		res.Airvelocity = m["Airvelocity"].(execute.ThreadParam).Value.(float64)
	}

	vExecutiontime, is := m["Executiontime"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vExecutiontime.JSONString), &temp)
		res.Executiontime = *temp.Executiontime
	} else {
		res.Executiontime = m["Executiontime"].(execute.ThreadParam).Value.(wunit.Time)
	}

	vLiquid, is := m["Liquid"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vLiquid.JSONString), &temp)
		res.Liquid = *temp.Liquid
	} else {
		res.Liquid = m["Liquid"].(execute.ThreadParam).Value.(string)
	}

	vPa, is := m["Pa"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vPa.JSONString), &temp)
		res.Pa = *temp.Pa
	} else {
		res.Pa = m["Pa"].(execute.ThreadParam).Value.(float64)
	}

	vPlatetype, is := m["Platetype"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vPlatetype.JSONString), &temp)
		res.Platetype = *temp.Platetype
	} else {
		res.Platetype = m["Platetype"].(execute.ThreadParam).Value.(string)
	}

	vRelativehumidity, is := m["Relativehumidity"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vRelativehumidity.JSONString), &temp)
		res.Relativehumidity = *temp.Relativehumidity
	} else {
		res.Relativehumidity = m["Relativehumidity"].(execute.ThreadParam).Value.(float64)
	}

	vTemp, is := m["Temp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vTemp.JSONString), &temp)
		res.Temp = *temp.Temp
	} else {
		res.Temp = m["Temp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vVolumeperwell, is := m["Volumeperwell"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp EvaporationrateJSONBlock
		json.Unmarshal([]byte(vVolumeperwell.JSONString), &temp)
		res.Volumeperwell = *temp.Volumeperwell
	} else {
		res.Volumeperwell = m["Volumeperwell"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["Airvelocity"].(execute.ThreadParam).ID
	res.BlockID = m["Airvelocity"].(execute.ThreadParam).BlockID

	return res
}

// Go helper functions:

//Functions for rounding numbers to a specified number of decimal places (places):
/*func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}
*/
/* This function calculates Θ required for the evaporation calculator based on air velocity above the sample;
this will be important in a laminar flow cabinet, fume cabinet and when the plates are mixing:
*/

/*: 0.62198 * pws / (pa - pws), // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
"x":  0.62198 * pw / (pa - pw), */

func (e *Evaporationrate) OnAirvelocity(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Airvelocity", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Evaporationrate) OnExecutiontime(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Executiontime", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Evaporationrate) OnLiquid(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Liquid", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Evaporationrate) OnPa(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *Evaporationrate) OnPlatetype(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Platetype", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Evaporationrate) OnRelativehumidity(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *Evaporationrate) OnTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
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
func (e *Evaporationrate) OnVolumeperwell(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(8, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Volumeperwell", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Evaporationrate struct {
	flow.Component           // component "superclass" embedded
	lock                     sync.Mutex
	startup                  sync.Once
	params                   map[execute.ThreadID]*execute.AsyncBag
	Airvelocity              <-chan execute.ThreadParam
	Executiontime            <-chan execute.ThreadParam
	Liquid                   <-chan execute.ThreadParam
	Pa                       <-chan execute.ThreadParam
	Platetype                <-chan execute.ThreadParam
	Relativehumidity         <-chan execute.ThreadParam
	Temp                     <-chan execute.ThreadParam
	Volumeperwell            <-chan execute.ThreadParam
	Estimatedevaporationtime chan<- execute.ThreadParam
	Evaporatedliquid         chan<- execute.ThreadParam
	Evaporationrateestimate  chan<- execute.ThreadParam
	Status                   chan<- execute.ThreadParam
}

type EvaporationrateParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Airvelocity      float64
	Executiontime    wunit.Time
	Liquid           string
	Pa               float64
	Platetype        string
	Relativehumidity float64
	Temp             wunit.Temperature
	Volumeperwell    wunit.Volume
}

type EvaporationrateConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Airvelocity      float64
	Executiontime    wunit.Time
	Liquid           string
	Pa               float64
	Platetype        string
	Relativehumidity float64
	Temp             wunit.Temperature
	Volumeperwell    wunit.Volume
}

type EvaporationrateResultBlock struct {
	ID                       execute.ThreadID
	BlockID                  execute.BlockID
	Error                    bool
	Estimatedevaporationtime wunit.Time
	Evaporatedliquid         wunit.Volume
	Evaporationrateestimate  float64
	Status                   string
}

type EvaporationrateJSONBlock struct {
	ID                       *execute.ThreadID
	BlockID                  *execute.BlockID
	Error                    *bool
	Airvelocity              *float64
	Executiontime            *wunit.Time
	Liquid                   *string
	Pa                       *float64
	Platetype                *string
	Relativehumidity         *float64
	Temp                     *wunit.Temperature
	Volumeperwell            *wunit.Volume
	Estimatedevaporationtime *wunit.Time
	Evaporatedliquid         *wunit.Volume
	Evaporationrateestimate  *float64
	Status                   *string
}

func (c *Evaporationrate) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Airvelocity", "float64", "Airvelocity", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Executiontime", "wunit.Time", "Executiontime", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Liquid", "string", "Liquid", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Pa", "float64", "Pa", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Platetype", "string", "Platetype", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Relativehumidity", "float64", "Relativehumidity", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Temp", "wunit.Temperature", "Temp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Volumeperwell", "wunit.Volume", "Volumeperwell", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Estimatedevaporationtime", "wunit.Time", "Estimatedevaporationtime", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Evaporatedliquid", "wunit.Volume", "Evaporatedliquid", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Evaporationrateestimate", "float64", "Evaporationrateestimate", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Evaporationrate", "Evaporationrate", "", false, inp, outp)

	return ci
}
