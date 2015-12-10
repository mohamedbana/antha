package Transfer

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

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Transfer) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Transfer) setup(p TransferParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Transfer) steps(p TransferParamBlock, r *TransferResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	sample := mixer.Sample(p.Startingsolution, p.LiquidVolume)
	r.FinalSolution = _wrapper.MixInto(p.OutPlate, sample)

	r.Status = p.LiquidVolume.ToString() + " of " + p.Liquidname + " was mixed into " + p.OutPlate.Type
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Transfer) analysis(p TransferParamBlock, r *TransferResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Transfer) validation(p TransferParamBlock, r *TransferResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Transfer) Complete(params interface{}) {
	p := params.(TransferParamBlock)
	if p.Error {
		e.FinalSolution <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TransferResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.FinalSolution <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.FinalSolution <- execute.ThreadParam{Value: r.FinalSolution, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Transfer) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Transfer) NewConfig() interface{} {
	return &TransferConfig{}
}

func (e *Transfer) NewParamBlock() interface{} {
	return &TransferParamBlock{}
}

func NewTransfer() interface{} { //*Transfer {
	e := new(Transfer)
	e.init()
	return e
}

// Mapper function
func (e *Transfer) Map(m map[string]interface{}) interface{} {
	var res TransferParamBlock
	res.Error = false || m["LiquidVolume"].(execute.ThreadParam).Error || m["Liquidname"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Startingsolution"].(execute.ThreadParam).Error

	vLiquidVolume, is := m["LiquidVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransferJSONBlock
		json.Unmarshal([]byte(vLiquidVolume.JSONString), &temp)
		res.LiquidVolume = *temp.LiquidVolume
	} else {
		res.LiquidVolume = m["LiquidVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vLiquidname, is := m["Liquidname"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransferJSONBlock
		json.Unmarshal([]byte(vLiquidname.JSONString), &temp)
		res.Liquidname = *temp.Liquidname
	} else {
		res.Liquidname = m["Liquidname"].(execute.ThreadParam).Value.(string)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransferJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vStartingsolution, is := m["Startingsolution"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TransferJSONBlock
		json.Unmarshal([]byte(vStartingsolution.JSONString), &temp)
		res.Startingsolution = *temp.Startingsolution
	} else {
		res.Startingsolution = m["Startingsolution"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["LiquidVolume"].(execute.ThreadParam).ID
	res.BlockID = m["LiquidVolume"].(execute.ThreadParam).BlockID

	return res
}

func (e *Transfer) OnLiquidVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("LiquidVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transfer) OnLiquidname(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Liquidname", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Transfer) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
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
func (e *Transfer) OnStartingsolution(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Startingsolution", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Transfer struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	LiquidVolume     <-chan execute.ThreadParam
	Liquidname       <-chan execute.ThreadParam
	OutPlate         <-chan execute.ThreadParam
	Startingsolution <-chan execute.ThreadParam
	FinalSolution    chan<- execute.ThreadParam
	Status           chan<- execute.ThreadParam
}

type TransferParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	LiquidVolume     wunit.Volume
	Liquidname       string
	OutPlate         *wtype.LHPlate
	Startingsolution *wtype.LHComponent
}

type TransferConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	LiquidVolume     wunit.Volume
	Liquidname       string
	OutPlate         wtype.FromFactory
	Startingsolution wtype.FromFactory
}

type TransferResultBlock struct {
	ID            execute.ThreadID
	BlockID       execute.BlockID
	Error         bool
	FinalSolution *wtype.LHSolution
	Status        string
}

type TransferJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	LiquidVolume     *wunit.Volume
	Liquidname       *string
	OutPlate         **wtype.LHPlate
	Startingsolution **wtype.LHComponent
	FinalSolution    **wtype.LHSolution
	Status           *string
}

func (c *Transfer) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("LiquidVolume", "wunit.Volume", "LiquidVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Liquidname", "string", "Liquidname", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Startingsolution", "*wtype.LHComponent", "Startingsolution", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("FinalSolution", "*wtype.LHSolution", "FinalSolution", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("Transfer", "Transfer", "", false, inp, outp)

	return ci
}
