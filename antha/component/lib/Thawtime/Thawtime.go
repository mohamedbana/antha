//status = compiles and calculates; need to fill in correct parameters and check units
//currently using dummy values only so won't be accurate yet!
// Once working move from floats to antha types and units
package Thawtime

import (
	"encoding/json"
	"fmt" // we need this go library to print
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

// all of our functions used here are in the Thaw.go file in the eng package which this points to
//"github.com/montanaflynn/stats" // a rounding function is used from this third party library

// Many of the real parameters required will be looked up via the specific labware (platetype) and liquid type which are being used.

/* e.g. the sample volume as frozen by a previous storage protocol;
could be known or measured via liquid height detection on some liquid handlers */

// These should be captured via sensors just prior to execution

// This will be monitored via the thermometer in the freezer in which the sample was stored

/* This will offer another knob to tweak (in addition to the other parameters) as a means to improve
the correlation over time as we see how accurate the calculator is in practice */

func (e *Thawtime) requirements() {
	_ = wunit.Make_units

}
func (e *Thawtime) setup(p ThawtimeParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *Thawtime) steps(p ThawtimeParamBlock, r *ThawtimeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	/*  Step 1. we need a mass for the following equations so we calculate this by looking up
	the liquid density and multiplying by the fill volume using this function from the engineering library */

	//fillvolume:= Fillvolume.SIValue()

	mass := eng.Massfromvolume(p.Fillvolume, p.Liquid)

	/*  Step 2. Required heat energy to melt the solid is calculated using the calculated mass along with the latent heat of melting
	which we find via a liquid class look up package which is not required for import here since it's imported from the engineering library */

	q := eng.Q(p.Liquid, mass)

	/*  Step 3. Heat will be transferred via both convection through the air and conduction through the plate walls.
	Let's first work out the heat energy transferred via convection, this uses an empirical parameter,
	the convective heat transfer coefficient of air (HC_air), this is calculated via another function in the eng library.
	In future we will make this process slightly more sophisticated by adding conditions, since this empirical equation is
	only validated between air velocities 2 - 20 m/s. It could also be adjusted to calculate heat transfer if the sample
	is agitated on a shaker to speed up thawing. */

	hc_air := eng.Hc_air(p.Airvelocity)

	/*  Step 4. The rate of heat transfer by convection is then calculated using this value combined with the temperature differential
	(measured by the temp sensor) and surface area dictated by the plate type (another look up called from the eng library!)*/

	convection := eng.ConvectionPowertransferred(hc_air, p.Platetype, p.SurfaceTemp, p.BulkTemp)

	/*  Step 5. We now estimate the heat transfer rate via conduction. For this we need to know the thermal conductivity of the plate material
	along with the wall thickness. As before, both of these are looked up via the labware library called by this function in the eng library */

	conduction := eng.ConductionPowertransferred(p.Platetype, p.SurfaceTemp, p.BulkTemp)

	/*  Step 6. We're now ready to estimate the thawtime needed by simply dividing the estimated heat required to melt/thaw (i.e. q from step 2)
	by the combined rate of heat transfer estimated to occur via both convection and conduction */
	r.Estimatedthawtime = eng.Thawtime(convection, conduction, q)

	/* Step 7. Since there're a lot of assumptions here (liquid behaves as water, no change in temperature gradient, no heat transferred via radiation,
	imprecision in the estimates and 	empirical formaulas) we'll multiply by a fudgefactor to be safer that we've definitely thawed,
	this (and all parameters!) can be adjusted over time as we see emprically how reliable this function is as more datapoints are collected */
	r.Thawtimeused = wunit.NewTime(r.Estimatedthawtime.SIValue()*p.Fudgefactor, "s")

	r.Status = fmt.Sprintln("For", mass.ToString(), "of", p.Liquid, "in", p.Platetype,
		"Thawtime required =", r.Estimatedthawtime.ToString(),
		"Thawtime used =", r.Thawtimeused.ToString(),
		"power required =", q, "J",
		"HC_air (convective heat transfer coefficient=", hc_air,
		"Convective power=", convection, "J/s",
		"conductive power=", conduction, "J/s")
	_ = _wrapper.WaitToEnd()

}
func (e *Thawtime) analysis(p ThawtimeParamBlock, r *ThawtimeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *Thawtime) validation(p ThawtimeParamBlock, r *ThawtimeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Thawtime) Complete(params interface{}) {
	p := params.(ThawtimeParamBlock)
	if p.Error {
		e.Estimatedthawtime <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Thawtimeused <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(ThawtimeResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Estimatedthawtime <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Thawtimeused <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Estimatedthawtime <- execute.ThreadParam{Value: r.Estimatedthawtime, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Thawtimeused <- execute.ThreadParam{Value: r.Thawtimeused, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Thawtime) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Thawtime) NewConfig() interface{} {
	return &ThawtimeConfig{}
}

func (e *Thawtime) NewParamBlock() interface{} {
	return &ThawtimeParamBlock{}
}

func NewThawtime() interface{} { //*Thawtime {
	e := new(Thawtime)
	e.init()
	return e
}

// Mapper function
func (e *Thawtime) Map(m map[string]interface{}) interface{} {
	var res ThawtimeParamBlock
	res.Error = false || m["Airvelocity"].(execute.ThreadParam).Error || m["BulkTemp"].(execute.ThreadParam).Error || m["Fillvolume"].(execute.ThreadParam).Error || m["Fudgefactor"].(execute.ThreadParam).Error || m["Liquid"].(execute.ThreadParam).Error || m["Platetype"].(execute.ThreadParam).Error || m["SurfaceTemp"].(execute.ThreadParam).Error

	vAirvelocity, is := m["Airvelocity"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vAirvelocity.JSONString), &temp)
		res.Airvelocity = *temp.Airvelocity
	} else {
		res.Airvelocity = m["Airvelocity"].(execute.ThreadParam).Value.(float64)
	}

	vBulkTemp, is := m["BulkTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vBulkTemp.JSONString), &temp)
		res.BulkTemp = *temp.BulkTemp
	} else {
		res.BulkTemp = m["BulkTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vFillvolume, is := m["Fillvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vFillvolume.JSONString), &temp)
		res.Fillvolume = *temp.Fillvolume
	} else {
		res.Fillvolume = m["Fillvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vFudgefactor, is := m["Fudgefactor"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vFudgefactor.JSONString), &temp)
		res.Fudgefactor = *temp.Fudgefactor
	} else {
		res.Fudgefactor = m["Fudgefactor"].(execute.ThreadParam).Value.(float64)
	}

	vLiquid, is := m["Liquid"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vLiquid.JSONString), &temp)
		res.Liquid = *temp.Liquid
	} else {
		res.Liquid = m["Liquid"].(execute.ThreadParam).Value.(string)
	}

	vPlatetype, is := m["Platetype"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vPlatetype.JSONString), &temp)
		res.Platetype = *temp.Platetype
	} else {
		res.Platetype = m["Platetype"].(execute.ThreadParam).Value.(string)
	}

	vSurfaceTemp, is := m["SurfaceTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ThawtimeJSONBlock
		json.Unmarshal([]byte(vSurfaceTemp.JSONString), &temp)
		res.SurfaceTemp = *temp.SurfaceTemp
	} else {
		res.SurfaceTemp = m["SurfaceTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	res.ID = m["Airvelocity"].(execute.ThreadParam).ID
	res.BlockID = m["Airvelocity"].(execute.ThreadParam).BlockID

	return res
}

func (e *Thawtime) OnAirvelocity(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *Thawtime) OnBulkTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("BulkTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Thawtime) OnFillvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Fillvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Thawtime) OnFudgefactor(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Fudgefactor", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Thawtime) OnLiquid(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *Thawtime) OnPlatetype(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *Thawtime) OnSurfaceTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SurfaceTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Thawtime struct {
	flow.Component    // component "superclass" embedded
	lock              sync.Mutex
	startup           sync.Once
	params            map[execute.ThreadID]*execute.AsyncBag
	Airvelocity       <-chan execute.ThreadParam
	BulkTemp          <-chan execute.ThreadParam
	Fillvolume        <-chan execute.ThreadParam
	Fudgefactor       <-chan execute.ThreadParam
	Liquid            <-chan execute.ThreadParam
	Platetype         <-chan execute.ThreadParam
	SurfaceTemp       <-chan execute.ThreadParam
	Estimatedthawtime chan<- execute.ThreadParam
	Status            chan<- execute.ThreadParam
	Thawtimeused      chan<- execute.ThreadParam
}

type ThawtimeParamBlock struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	Airvelocity float64
	BulkTemp    wunit.Temperature
	Fillvolume  wunit.Volume
	Fudgefactor float64
	Liquid      string
	Platetype   string
	SurfaceTemp wunit.Temperature
}

type ThawtimeConfig struct {
	ID          execute.ThreadID
	BlockID     execute.BlockID
	Error       bool
	Airvelocity float64
	BulkTemp    wunit.Temperature
	Fillvolume  wunit.Volume
	Fudgefactor float64
	Liquid      string
	Platetype   string
	SurfaceTemp wunit.Temperature
}

type ThawtimeResultBlock struct {
	ID                execute.ThreadID
	BlockID           execute.BlockID
	Error             bool
	Estimatedthawtime wunit.Time
	Status            string
	Thawtimeused      wunit.Time
}

type ThawtimeJSONBlock struct {
	ID                *execute.ThreadID
	BlockID           *execute.BlockID
	Error             *bool
	Airvelocity       *float64
	BulkTemp          *wunit.Temperature
	Fillvolume        *wunit.Volume
	Fudgefactor       *float64
	Liquid            *string
	Platetype         *string
	SurfaceTemp       *wunit.Temperature
	Estimatedthawtime *wunit.Time
	Status            *string
	Thawtimeused      *wunit.Time
}

func (c *Thawtime) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Airvelocity", "float64", "Airvelocity", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("BulkTemp", "wunit.Temperature", "BulkTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Fillvolume", "wunit.Volume", "Fillvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Fudgefactor", "float64", "Fudgefactor", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Liquid", "string", "Liquid", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Platetype", "string", "Platetype", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SurfaceTemp", "wunit.Temperature", "SurfaceTemp", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Estimatedthawtime", "wunit.Time", "Estimatedthawtime", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Thawtimeused", "wunit.Time", "Thawtimeused", true, true, nil, nil))

	ci := execute.NewComponentInfo("Thawtime", "Thawtime", "", false, inp, outp)

	return ci
}
