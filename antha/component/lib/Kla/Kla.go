/* Islam, R. S., Tisi, D., Levy, M. S. & Lye, G. J. Scale-up of Escherichia coli growth and recombinant protein expression conditions from microwell to laboratory and pilot scale based on matched kLa. Biotechnol. Bioeng. 99, 1128–1139 (2008).

equation (6)

func kLa_squaremicrowell = (3.94 x 10E-4) * (D/dv)* ai * RE^1.91 * exp ^ (a * Fr^b) // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case

kla = dimensionless
	var D = diffusion coefficient, m2 􏰀 s􏰁1
	var dv = microwell vessel diameter, m
	var ai = initial specific surface area, m􏰁1
	var RE = Reynolds number, (ro * n * dv * 2/mu), dimensionless
		var	ro	= density, kg 􏰀/ m􏰁3
		var	n 	= shaking frequency, s􏰁1
		var	mu	= viscosity, kg 􏰀/ m􏰁 /􏰀 s
	const exp = Eulers number, 2.718281828

	var Fr = Froude number = dt(2 * math.Pi * n)^2 /(2 * g), (dimensionless)
		var dt = shaking amplitude, m
		const g = acceleration due to gravity, m 􏰀/ s􏰁2
	const	a = constant
	const	b = constant
*/
// make type /time and units of /hour and per second
// check accuracy against literature and experimental values
package Kla

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Setpoints"
//"github.com/montanaflynn/stats"

//diffusion coefficient, m2 􏰀 s􏰁1 // from wikipedia: Oxygen (dis) - Water (l) 	@25 degrees C 	2.10x10−5 cm2/s // should call from elsewhere really
// add temp etc?

func (e *Kla) requirements() {
	_ = wunit.Make_units

}
func (e *Kla) setup(p KlaParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *Kla) steps(p KlaParamBlock, r *KlaResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	dv := labware.Labwaregeometry[p.Platetype]["dv"] // microwell vessel diameter, m 0.017 //
	ai := labware.Labwaregeometry[p.Platetype]["ai"] // initial specific surface area, /m 96.0
	//var RE = Reynolds number, (ro * n * dv * 2/mu), dimensionless

	ro := liquidclasses.Liquidclass[p.Liquid]["ro"] //density, kg 􏰀/ m􏰁3 999.7 // environment dependent
	mu := liquidclasses.Liquidclass[p.Liquid]["mu"] //0.001           environment dependent                        //liquidclasses.Liquidclass[liquid]["mu"] viscosity, kg 􏰀/ m􏰁 /􏰀 s

	n := p.Rpm / 60 //shaking frequency, s􏰁1
	//const exp = Eulers number, 2.718281828

	//Fr = Froude number = dt(2 * math.Pi * n)^2 /(2 * g), (dimensionless)
	dt := devices.Shaker[p.Shakertype]["dt"] //0.008                                  //shaking amplitude, m // move to shaker package

	a := labware.Labwaregeometry[p.Platetype]["a"] //0.88   //
	b := labware.Labwaregeometry[p.Platetype]["b"] //1.24

	Fr := eng.Froude(dt, n, eng.G)
	Re := eng.RE(ro, n, mu, dv)
	r.Necessaryshakerspeed = eng.Shakerspeed(p.TargetRE, ro, mu, dv)

	Vl := p.Fillvolume.SIValue()
	Sigma := liquidclasses.Liquidclass[p.Liquid]["sigma"]

	// Check Ncrit! original paper used this to calculate speed in shallow round well plates... double check paper

	// add loop to use correct formula dependent on Platetype etc...
	/*Criticalshakerspeed := "error"
	if labware.Labwaregeometry[Platetype]["numberofwellsides"] == 4.0 {*/
	r.Ncrit = eng.Ncrit_srw(Sigma, dv, Vl, ro, dt)
	//}
	/*if i == 4.0 {
		Criticalshakerspeed := "error"
	}
	*/
	//Criticalshakerspeed := stats.Round(eng.Ncrit_srw(Sigma, dv, Vl , ro , dt ),3)

	if Re > 5E3 {
		r.Flowstate = fmt.Sprintln("Flowstate = Turbulent flow")
	}

	//klainputs :=fmt.Sprintln("D",D,"dv", dv,"ai", ai,"Re", Re,"a", a,"Fr", Fr,"b", b)

	r.CalculatedKla = eng.KLa_squaremicrowell(p.D, dv, ai, Re, a, Fr, b)

	r.Status = fmt.Sprintln("TargetRE = ", p.TargetRE, "Calculated Reynolds number = ", Re, "shakerspeedrequired for targetRE= ", r.Necessaryshakerspeed*60, "Froude number = ", Fr, "kla =", r.CalculatedKla, "/h", "Ncrit	=", r.Ncrit, "/S")
	_ = _wrapper.WaitToEnd()

	//CalculatedKla = setpoints.CalculateKlasquaremicrowell(Platetype, Liquid, Rpm, Shakertype, TargetRE, D)

	/*
		fmt.Println("shakerspeedrequired= ", stats.Round(Necessaryshakerspeed*60, 3))
		fmt.Println("Froude number = ", stats.Round(Fr, 3))
		fmt.Println("kla =", stats.Round(CalculatedKla, 3))
		fmt.Println("Shaker speed required for turbulence	=", Criticalshakerspeed,"/S")*/
	//fmt.Println("=", (Criticalshakerspeed*60), 3),"rpm")
}
func (e *Kla) analysis(p KlaParamBlock, r *KlaResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

} // works in either analysis or steps sections

func (e *Kla) validation(p KlaParamBlock, r *KlaResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

	/*if Evaporatedliquid > Volumeperwell {
	panic("not enough liquid") */
}

// AsyncBag functions
func (e *Kla) Complete(params interface{}) {
	p := params.(KlaParamBlock)
	if p.Error {
		e.CalculatedKla <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Flowstate <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Ncrit <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Necessaryshakerspeed <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(KlaResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.CalculatedKla <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Flowstate <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Ncrit <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Necessaryshakerspeed <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.CalculatedKla <- execute.ThreadParam{Value: r.CalculatedKla, ID: p.ID, Error: false}

	e.Flowstate <- execute.ThreadParam{Value: r.Flowstate, ID: p.ID, Error: false}

	e.Ncrit <- execute.ThreadParam{Value: r.Ncrit, ID: p.ID, Error: false}

	e.Necessaryshakerspeed <- execute.ThreadParam{Value: r.Necessaryshakerspeed, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Kla) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Kla) NewConfig() interface{} {
	return &KlaConfig{}
}

func (e *Kla) NewParamBlock() interface{} {
	return &KlaParamBlock{}
}

func NewKla() interface{} { //*Kla {
	e := new(Kla)
	e.init()
	return e
}

// Mapper function
func (e *Kla) Map(m map[string]interface{}) interface{} {
	var res KlaParamBlock
	res.Error = false || m["D"].(execute.ThreadParam).Error || m["Fillvolume"].(execute.ThreadParam).Error || m["Liquid"].(execute.ThreadParam).Error || m["Platetype"].(execute.ThreadParam).Error || m["Rpm"].(execute.ThreadParam).Error || m["Shakertype"].(execute.ThreadParam).Error || m["TargetRE"].(execute.ThreadParam).Error

	vD, is := m["D"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vD.JSONString), &temp)
		res.D = *temp.D
	} else {
		res.D = m["D"].(execute.ThreadParam).Value.(float64)
	}

	vFillvolume, is := m["Fillvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vFillvolume.JSONString), &temp)
		res.Fillvolume = *temp.Fillvolume
	} else {
		res.Fillvolume = m["Fillvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vLiquid, is := m["Liquid"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vLiquid.JSONString), &temp)
		res.Liquid = *temp.Liquid
	} else {
		res.Liquid = m["Liquid"].(execute.ThreadParam).Value.(string)
	}

	vPlatetype, is := m["Platetype"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vPlatetype.JSONString), &temp)
		res.Platetype = *temp.Platetype
	} else {
		res.Platetype = m["Platetype"].(execute.ThreadParam).Value.(string)
	}

	vRpm, is := m["Rpm"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vRpm.JSONString), &temp)
		res.Rpm = *temp.Rpm
	} else {
		res.Rpm = m["Rpm"].(execute.ThreadParam).Value.(float64)
	}

	vShakertype, is := m["Shakertype"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vShakertype.JSONString), &temp)
		res.Shakertype = *temp.Shakertype
	} else {
		res.Shakertype = m["Shakertype"].(execute.ThreadParam).Value.(string)
	}

	vTargetRE, is := m["TargetRE"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp KlaJSONBlock
		json.Unmarshal([]byte(vTargetRE.JSONString), &temp)
		res.TargetRE = *temp.TargetRE
	} else {
		res.TargetRE = m["TargetRE"].(execute.ThreadParam).Value.(float64)
	}

	res.ID = m["D"].(execute.ThreadParam).ID
	res.BlockID = m["D"].(execute.ThreadParam).BlockID

	return res
}

func (e *Kla) OnD(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("D", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Kla) OnFillvolume(param execute.ThreadParam) {
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
func (e *Kla) OnLiquid(param execute.ThreadParam) {
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
func (e *Kla) OnPlatetype(param execute.ThreadParam) {
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
func (e *Kla) OnRpm(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Rpm", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Kla) OnShakertype(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Shakertype", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Kla) OnTargetRE(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TargetRE", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Kla struct {
	flow.Component       // component "superclass" embedded
	lock                 sync.Mutex
	startup              sync.Once
	params               map[execute.ThreadID]*execute.AsyncBag
	D                    <-chan execute.ThreadParam
	Fillvolume           <-chan execute.ThreadParam
	Liquid               <-chan execute.ThreadParam
	Platetype            <-chan execute.ThreadParam
	Rpm                  <-chan execute.ThreadParam
	Shakertype           <-chan execute.ThreadParam
	TargetRE             <-chan execute.ThreadParam
	CalculatedKla        chan<- execute.ThreadParam
	Flowstate            chan<- execute.ThreadParam
	Ncrit                chan<- execute.ThreadParam
	Necessaryshakerspeed chan<- execute.ThreadParam
	Status               chan<- execute.ThreadParam
}

type KlaParamBlock struct {
	ID         execute.ThreadID
	BlockID    execute.BlockID
	Error      bool
	D          float64
	Fillvolume wunit.Volume
	Liquid     string
	Platetype  string
	Rpm        float64
	Shakertype string
	TargetRE   float64
}

type KlaConfig struct {
	ID         execute.ThreadID
	BlockID    execute.BlockID
	Error      bool
	D          float64
	Fillvolume wunit.Volume
	Liquid     string
	Platetype  string
	Rpm        float64
	Shakertype string
	TargetRE   float64
}

type KlaResultBlock struct {
	ID                   execute.ThreadID
	BlockID              execute.BlockID
	Error                bool
	CalculatedKla        float64
	Flowstate            string
	Ncrit                float64
	Necessaryshakerspeed float64
	Status               string
}

type KlaJSONBlock struct {
	ID                   *execute.ThreadID
	BlockID              *execute.BlockID
	Error                *bool
	D                    *float64
	Fillvolume           *wunit.Volume
	Liquid               *string
	Platetype            *string
	Rpm                  *float64
	Shakertype           *string
	TargetRE             *float64
	CalculatedKla        *float64
	Flowstate            *string
	Ncrit                *float64
	Necessaryshakerspeed *float64
	Status               *string
}

func (c *Kla) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("D", "float64", "D", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Fillvolume", "wunit.Volume", "Fillvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Liquid", "string", "Liquid", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Platetype", "string", "Platetype", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Rpm", "float64", "Rpm", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Shakertype", "string", "Shakertype", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TargetRE", "float64", "TargetRE", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("CalculatedKla", "float64", "CalculatedKla", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Flowstate", "string", "Flowstate", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Ncrit", "float64", "Ncrit", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Necessaryshakerspeed", "float64", "Necessaryshakerspeed", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Kla", "Kla", "", false, inp, outp)

	return ci
}
