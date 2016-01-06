package Plotdata

import (
	"encoding/json"
	graph "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol (data)

/*datarange*/
/*datarange*/

//	HeaderRange []string

// Data which is returned from this protocol, and data types

//	OutputData       []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Plotdata) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Plotdata) setup(p PlotdataParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Plotdata) steps(p PlotdataParamBlock, r *PlotdataResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	// now plot the graph

	// the data points

	plot := graph.Plot(p.Xvalues, p.Yvaluearray)

	graph.Export(plot, p.Exportedfilename)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Plotdata) analysis(p PlotdataParamBlock, r *PlotdataResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Plotdata) validation(p PlotdataParamBlock, r *PlotdataResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Plotdata) Complete(params interface{}) {
	p := params.(PlotdataParamBlock)
	if p.Error {
		return
	}
	r := new(PlotdataResultBlock)
	defer func() {
		if res := recover(); res != nil {
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Plotdata) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Plotdata) NewConfig() interface{} {
	return &PlotdataConfig{}
}

func (e *Plotdata) NewParamBlock() interface{} {
	return &PlotdataParamBlock{}
}

func NewPlotdata() interface{} { //*Plotdata {
	e := new(Plotdata)
	e.init()
	return e
}

// Mapper function
func (e *Plotdata) Map(m map[string]interface{}) interface{} {
	var res PlotdataParamBlock
	res.Error = false || m["Exportedfilename"].(execute.ThreadParam).Error || m["Xvalues"].(execute.ThreadParam).Error || m["Yvaluearray"].(execute.ThreadParam).Error

	vExportedfilename, is := m["Exportedfilename"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlotdataJSONBlock
		json.Unmarshal([]byte(vExportedfilename.JSONString), &temp)
		res.Exportedfilename = *temp.Exportedfilename
	} else {
		res.Exportedfilename = m["Exportedfilename"].(execute.ThreadParam).Value.(string)
	}

	vXvalues, is := m["Xvalues"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlotdataJSONBlock
		json.Unmarshal([]byte(vXvalues.JSONString), &temp)
		res.Xvalues = *temp.Xvalues
	} else {
		res.Xvalues = m["Xvalues"].(execute.ThreadParam).Value.([]float64)
	}

	vYvaluearray, is := m["Yvaluearray"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PlotdataJSONBlock
		json.Unmarshal([]byte(vYvaluearray.JSONString), &temp)
		res.Yvaluearray = *temp.Yvaluearray
	} else {
		res.Yvaluearray = m["Yvaluearray"].(execute.ThreadParam).Value.([][]float64)
	}

	res.ID = m["Exportedfilename"].(execute.ThreadParam).ID
	res.BlockID = m["Exportedfilename"].(execute.ThreadParam).BlockID

	return res
}

func (e *Plotdata) OnExportedfilename(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(3, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Exportedfilename", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Plotdata) OnXvalues(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(3, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Xvalues", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Plotdata) OnYvaluearray(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(3, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Yvaluearray", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Plotdata struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	Exportedfilename <-chan execute.ThreadParam
	Xvalues          <-chan execute.ThreadParam
	Yvaluearray      <-chan execute.ThreadParam
}

type PlotdataParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Exportedfilename string
	Xvalues          []float64
	Yvaluearray      [][]float64
}

type PlotdataConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Exportedfilename string
	Xvalues          []float64
	Yvaluearray      [][]float64
}

type PlotdataResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
}

type PlotdataJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	Exportedfilename *string
	Xvalues          *[]float64
	Yvaluearray      *[][]float64
}

func (c *Plotdata) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Exportedfilename", "string", "Exportedfilename", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Xvalues", "[]float64", "Xvalues", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Yvaluearray", "[][]float64", "Yvaluearray", true, true, nil, nil))

	ci := execute.NewComponentInfo("Plotdata", "Plotdata", "", false, inp, outp)

	return ci
}
