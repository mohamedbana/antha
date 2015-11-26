package SumVolume

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

//"github.com/antha-lang/antha/antha/anthalib/wunit"
// Input parameters for this protocol

//D Concentration
//E float64

// Data which is returned from this protocol

//DmolarConc wunit.MolarConcentration

// Physical inputs to this protocol

// Physical outputs from this protocol

func (e *SumVolume) requirements() {
	_ = wunit.Make_units

}

// Actions to perform before protocol itself
func (e *SumVolume) setup(p SumVolumeParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// Core process of the protocol: steps to be performed for each input
func (e *SumVolume) steps(p SumVolumeParamBlock, r *SumVolumeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	//var Dmassconc wunit.MassConcentration = D

	/*	molarmass := wunit.NewAmount(E,"M")

		var Dnew = wunit.MoleculeConcentration{D,E}

		mass := wunit.NewMass(1,"g")

		DmolarConc = Dnew.AsMolar(mass)
	*/
	r.Sum = *(wunit.CopyVolume(&p.A))
	(&r.Sum).Add(&p.B)
	r.Status = fmt.Sprintln(
		"Sum of", p.A.ToString(), "and", p.B.ToString(), "=", r.Sum.ToString(), "Temp=", p.C.ToString(),
	)
	_ = _wrapper.WaitToEnd()

	//"D Concentration in g/l", D, "D concentration in M/l", DmolarConc)
}

// Actions to perform after steps block to analyze data
func (e *SumVolume) analysis(p SumVolumeParamBlock, r *SumVolumeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *SumVolume) validation(p SumVolumeParamBlock, r *SumVolumeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *SumVolume) Complete(params interface{}) {
	p := params.(SumVolumeParamBlock)
	if p.Error {
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Sum <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(SumVolumeResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Sum <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Sum <- execute.ThreadParam{Value: r.Sum, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *SumVolume) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *SumVolume) NewConfig() interface{} {
	return &SumVolumeConfig{}
}

func (e *SumVolume) NewParamBlock() interface{} {
	return &SumVolumeParamBlock{}
}

func NewSumVolume() interface{} { //*SumVolume {
	e := new(SumVolume)
	e.init()
	return e
}

// Mapper function
func (e *SumVolume) Map(m map[string]interface{}) interface{} {
	var res SumVolumeParamBlock
	res.Error = false || m["A"].(execute.ThreadParam).Error || m["B"].(execute.ThreadParam).Error || m["C"].(execute.ThreadParam).Error

	vA, is := m["A"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SumVolumeJSONBlock
		json.Unmarshal([]byte(vA.JSONString), &temp)
		res.A = *temp.A
	} else {
		res.A = m["A"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vB, is := m["B"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SumVolumeJSONBlock
		json.Unmarshal([]byte(vB.JSONString), &temp)
		res.B = *temp.B
	} else {
		res.B = m["B"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vC, is := m["C"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp SumVolumeJSONBlock
		json.Unmarshal([]byte(vC.JSONString), &temp)
		res.C = *temp.C
	} else {
		res.C = m["C"].(execute.ThreadParam).Value.(wunit.Temperature)
	}

	res.ID = m["A"].(execute.ThreadParam).ID
	res.BlockID = m["A"].(execute.ThreadParam).BlockID

	return res
}

func (e *SumVolume) OnA(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(3, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("A", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SumVolume) OnB(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(3, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("B", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *SumVolume) OnC(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(3, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("C", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type SumVolume struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	A              <-chan execute.ThreadParam
	B              <-chan execute.ThreadParam
	C              <-chan execute.ThreadParam
	Status         chan<- execute.ThreadParam
	Sum            chan<- execute.ThreadParam
}

type SumVolumeParamBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	A       wunit.Volume
	B       wunit.Volume
	C       wunit.Temperature
}

type SumVolumeConfig struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	A       wunit.Volume
	B       wunit.Volume
	C       wunit.Temperature
}

type SumVolumeResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Status  string
	Sum     wunit.Volume
}

type SumVolumeJSONBlock struct {
	ID      *execute.ThreadID
	BlockID *execute.BlockID
	Error   *bool
	A       *wunit.Volume
	B       *wunit.Volume
	C       *wunit.Temperature
	Status  *string
	Sum     *wunit.Volume
}

func (c *SumVolume) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("A", "wunit.Volume", "A", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("B", "wunit.Volume", "B", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("C", "wunit.Temperature", "C", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Sum", "wunit.Volume", "Sum", true, true, nil, nil))

	ci := execute.NewComponentInfo("SumVolume", "SumVolume", "", false, inp, outp)

	return ci
}
