package MakeMedia

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"strconv"
	"sync"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

//SolidComponentMasses []Volume //Mass // Should be Mass

//  +/- x  e.g. 7.0 +/- 0.2

//LiqComponentkeys	[]string
//Solidcomponentkeys	[]string // name or barcode id
//Acidkey string
//Basekey string

// Physical Inputs to this protocol with types

/*SolidComponents		[]*wtype.LHComponent // should be new type or field indicating solid and mass
Acid				*wtype.LHComponent
Base 				*wtype.LHComponent
*/

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func (e *MakeMedia) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *MakeMedia) setup(p MakeMediaParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *MakeMedia) steps(p MakeMediaParamBlock, r *MakeMediaResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	recipestring := make([]string, 0)
	var step string
	stepcounter := 1 // counting from 1 is what makes us human
	liquids := make([]*wtype.LHComponent, 0)
	step = text.Print("Recipe for: ", p.Name)
	recipestring = append(recipestring, step)

	for i, liq := range p.LiqComponents {
		liqsamp := mixer.Sample(liq, p.LiqComponentVolumes[i])
		liquids = append(liquids, liqsamp)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add "+p.LiqComponentVolumes[i].ToString()+" of "+liq.CName)
		recipestring = append(recipestring, step)
		stepcounter++
	}

	//solids := make([]*wtype.LHComponent,0)

	/*for k, sol := range SolidComponents {
		solsamp := mixer.Sample(sol,SolidComponentMasses[k])
		liquids = append(liquids,solsamp)
		step = text.Print("Step" + strconv.Itoa(stepcounter) + ": ", "add " + SolidComponentMasses[k].ToString() + " of " + sol.CName)
		recipestring = append(recipestring,step)
		stepcounter = stepcounter + k
	}*/

	watersample := mixer.SampleForTotalVolume(p.Water, p.TotalVolume)
	liquids = append(liquids, watersample)
	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add up to "+p.TotalVolume.ToString()+" of "+p.Water.CName)
	recipestring = append(recipestring, step)
	stepcounter++

	// Add pH handling functions and driver calls etc...

	description := fmt.Sprint("adjust pH to ", p.PH_setPoint, " +/-", p.PH_tolerance, " for temp ", p.PH_setPointTemp.ToString(), "C")
	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", description)
	recipestring = append(recipestring, step)
	stepcounter++

	/*
		prepH := MixInto(Vessel,liquids...)

		pHactual := prepH.Measure("pH")

		step = text.Print("pH measured = ", pHactual)
		recipestring = append(recipestring,step)

		//pHactual = wutil.Roundto(pHactual,PH_tolerance)

		pHmax := PH_setpoint + PH_tolerance
		pHmin := PH_setpoint - PH_tolerance

		if pHactual < pHmax || pHactual < pHmin {
			// basically just a series of sample, stir, wait and recheck pH
		Media, newph, componentadded = prepH.AdjustpH(PH_setPoint, pHactual, PH_setPointTemp,Acid,Base)

		step = text.Print("Adjusted pH = ", newpH)
		recipestring = append(recipestring,step)

		step = text.Print("Component added = ", componentadded.Vol + componentadded.Vunit + " of " + componentadded.Conc + componentadded.Cunit + " " + componentadded.CName + )
		recipestring = append(recipestring,step)
		}
	*/
	r.Media = _wrapper.MixInto(p.Vessel, liquids...)

	r.Status = fmt.Sprintln(recipestring)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *MakeMedia) analysis(p MakeMediaParamBlock, r *MakeMediaResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *MakeMedia) validation(p MakeMediaParamBlock, r *MakeMediaResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *MakeMedia) Complete(params interface{}) {
	p := params.(MakeMediaParamBlock)
	if p.Error {
		e.Media <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(MakeMediaResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Media <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Media <- execute.ThreadParam{Value: r.Media, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *MakeMedia) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *MakeMedia) NewConfig() interface{} {
	return &MakeMediaConfig{}
}

func (e *MakeMedia) NewParamBlock() interface{} {
	return &MakeMediaParamBlock{}
}

func NewMakeMedia() interface{} { //*MakeMedia {
	e := new(MakeMedia)
	e.init()
	return e
}

// Mapper function
func (e *MakeMedia) Map(m map[string]interface{}) interface{} {
	var res MakeMediaParamBlock
	res.Error = false || m["LiqComponentVolumes"].(execute.ThreadParam).Error || m["LiqComponents"].(execute.ThreadParam).Error || m["Name"].(execute.ThreadParam).Error || m["PH_setPoint"].(execute.ThreadParam).Error || m["PH_setPointTemp"].(execute.ThreadParam).Error || m["PH_tolerance"].(execute.ThreadParam).Error || m["TotalVolume"].(execute.ThreadParam).Error || m["Vessel"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error

	vLiqComponentVolumes, is := m["LiqComponentVolumes"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vLiqComponentVolumes.JSONString), &temp)
		res.LiqComponentVolumes = *temp.LiqComponentVolumes
	} else {
		res.LiqComponentVolumes = m["LiqComponentVolumes"].(execute.ThreadParam).Value.([]wunit.Volume)
	}

	vLiqComponents, is := m["LiqComponents"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vLiqComponents.JSONString), &temp)
		res.LiqComponents = *temp.LiqComponents
	} else {
		res.LiqComponents = m["LiqComponents"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vName, is := m["Name"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vName.JSONString), &temp)
		res.Name = *temp.Name
	} else {
		res.Name = m["Name"].(execute.ThreadParam).Value.(string)
	}

	vPH_setPoint, is := m["PH_setPoint"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vPH_setPoint.JSONString), &temp)
		res.PH_setPoint = *temp.PH_setPoint
	} else {
		res.PH_setPoint = m["PH_setPoint"].(execute.ThreadParam).Value.(float64)
	}

	vPH_setPointTemp, is := m["PH_setPointTemp"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vPH_setPointTemp.JSONString), &temp)
		res.PH_setPointTemp = *temp.PH_setPointTemp
	} else {
		res.PH_setPointTemp = m["PH_setPointTemp"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	vPH_tolerance, is := m["PH_tolerance"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vPH_tolerance.JSONString), &temp)
		res.PH_tolerance = *temp.PH_tolerance
	} else {
		res.PH_tolerance = m["PH_tolerance"].(execute.ThreadParam).Value.(float64)
	}

	vTotalVolume, is := m["TotalVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vTotalVolume.JSONString), &temp)
		res.TotalVolume = *temp.TotalVolume
	} else {
		res.TotalVolume = m["TotalVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vVessel, is := m["Vessel"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vVessel.JSONString), &temp)
		res.Vessel = *temp.Vessel
	} else {
		res.Vessel = m["Vessel"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MakeMediaJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["LiqComponentVolumes"].(execute.ThreadParam).ID
	res.BlockID = m["LiqComponentVolumes"].(execute.ThreadParam).BlockID

	return res
}

/*
type Mole struct {
	number float64
}*/

func (e *MakeMedia) OnLiqComponentVolumes(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("LiqComponentVolumes", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnLiqComponents(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("LiqComponents", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Name", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnPH_setPoint(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PH_setPoint", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnPH_setPointTemp(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PH_setPointTemp", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnPH_tolerance(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PH_tolerance", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnTotalVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("TotalVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnVessel(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vessel", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MakeMedia) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
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

type MakeMedia struct {
	flow.Component      // component "superclass" embedded
	lock                sync.Mutex
	startup             sync.Once
	params              map[execute.ThreadID]*execute.AsyncBag
	LiqComponentVolumes <-chan execute.ThreadParam
	LiqComponents       <-chan execute.ThreadParam
	Name                <-chan execute.ThreadParam
	PH_setPoint         <-chan execute.ThreadParam
	PH_setPointTemp     <-chan execute.ThreadParam
	PH_tolerance        <-chan execute.ThreadParam
	TotalVolume         <-chan execute.ThreadParam
	Vessel              <-chan execute.ThreadParam
	Water               <-chan execute.ThreadParam
	Media               chan<- execute.ThreadParam
	Status              chan<- execute.ThreadParam
}

type MakeMediaParamBlock struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	LiqComponentVolumes []wunit.Volume
	LiqComponents       []*wtype.LHComponent
	Name                string
	PH_setPoint         float64
	PH_setPointTemp     wunit.Temperature
	PH_tolerance        float64
	TotalVolume         wunit.Volume
	Vessel              *wtype.LHPlate
	Water               *wtype.LHComponent
}

type MakeMediaConfig struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	LiqComponentVolumes []wunit.Volume
	LiqComponents       []wtype.FromFactory
	Name                string
	PH_setPoint         float64
	PH_setPointTemp     wunit.Temperature
	PH_tolerance        float64
	TotalVolume         wunit.Volume
	Vessel              wtype.FromFactory
	Water               wtype.FromFactory
}

type MakeMediaResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Media   *wtype.LHSolution
	Status  string
}

type MakeMediaJSONBlock struct {
	ID                  *execute.ThreadID
	BlockID             *execute.BlockID
	Error               *bool
	LiqComponentVolumes *[]wunit.Volume
	LiqComponents       *[]*wtype.LHComponent
	Name                *string
	PH_setPoint         *float64
	PH_setPointTemp     *wunit.Temperature
	PH_tolerance        *float64
	TotalVolume         *wunit.Volume
	Vessel              **wtype.LHPlate
	Water               **wtype.LHComponent
	Media               **wtype.LHSolution
	Status              *string
}

func (c *MakeMedia) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("LiqComponentVolumes", "[]wunit.Volume", "LiqComponentVolumes", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("LiqComponents", "[]*wtype.LHComponent", "LiqComponents", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Name", "string", "Name", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PH_setPoint", "float64", "PH_setPoint", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PH_setPointTemp", "wunit.Temperature", "PH_setPointTemp", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PH_tolerance", "float64", "PH_tolerance", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("TotalVolume", "wunit.Volume", "TotalVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vessel", "*wtype.LHPlate", "Vessel", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Media", "*wtype.LHSolution", "Media", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("MakeMedia", "MakeMedia", "", false, inp, outp)

	return ci
}
