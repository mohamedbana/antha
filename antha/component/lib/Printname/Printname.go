package Printname

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

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func (e *Printname) requirements() {
	_ = wunit.Make_units

}

// Actions to perform before protocol itself
func (e *Printname) setup(p PrintnameParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// Core process of the protocol: steps to be performed for each input
func (e *Printname) steps(p PrintnameParamBlock, r *PrintnameResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	if p.Name == "Michael Jackson" {
		r.Fullname = fmt.Sprintln(p.Name)
	} else {
		r.Fullname = "there's only one Michael Jackson, we accept no imitators"
	}
	_ = _wrapper.WaitToEnd()

}

// Actions to perform after steps block to analyze data
func (e *Printname) analysis(p PrintnameParamBlock, r *PrintnameResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *Printname) validation(p PrintnameParamBlock, r *PrintnameResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Printname) Complete(params interface{}) {
	p := params.(PrintnameParamBlock)
	if p.Error {
		e.Fullname <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PrintnameResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Fullname <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Fullname <- execute.ThreadParam{Value: r.Fullname, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Printname) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Printname) NewConfig() interface{} {
	return &PrintnameConfig{}
}

func (e *Printname) NewParamBlock() interface{} {
	return &PrintnameParamBlock{}
}

func NewPrintname() interface{} { //*Printname {
	e := new(Printname)
	e.init()
	return e
}

// Mapper function
func (e *Printname) Map(m map[string]interface{}) interface{} {
	var res PrintnameParamBlock
	res.Error = false || m["Name"].(execute.ThreadParam).Error

	vName, is := m["Name"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PrintnameJSONBlock
		json.Unmarshal([]byte(vName.JSONString), &temp)
		res.Name = *temp.Name
	} else {
		res.Name = m["Name"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["Name"].(execute.ThreadParam).ID
	res.BlockID = m["Name"].(execute.ThreadParam).BlockID

	return res
}

func (e *Printname) OnName(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(1, e, e)
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

type Printname struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	Name           <-chan execute.ThreadParam
	Fullname       chan<- execute.ThreadParam
}

type PrintnameParamBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Name    string
}

type PrintnameConfig struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Name    string
}

type PrintnameResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Fullname string
}

type PrintnameJSONBlock struct {
	ID       *execute.ThreadID
	BlockID  *execute.BlockID
	Error    *bool
	Name     *string
	Fullname *string
}

func (c *Printname) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Name", "string", "Name", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Fullname", "string", "Fullname", true, true, nil, nil))

	ci := execute.NewComponentInfo("Printname", "Printname", "", false, inp, outp)

	return ci
}
