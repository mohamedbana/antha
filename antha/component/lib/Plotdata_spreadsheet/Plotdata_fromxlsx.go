package Plotdata_spreadsheet

import (
	"fmt"
	//"math/rand"
	//"github.com/antha-lang/antha/internal/github.com/montanaflynn/stats"
	"encoding/json"
	graph "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol (data)

//                                                                         = "plotinumdata.xlsx"
//                                                                        = 0
/*datarange*/ //  = []string{"a4", "a16"}                                                           // row in A1 format i.e string{A,E} would indicate all data between those points
/*datarange*/ //= [][]string{[]string{"b4", "b16"}, []string{"c4", "c16"}, []string{"d4", "d16"}} // column in A1 format i.e string{1,12} would indicate all data between those points
//= "Excelfile.jpg"
//	HeaderRange []string // if Bycolumn == true, format would be e.g. string{A1,E1} else e.g. string{A1,A12}

// Data which is returned from this protocol, and data types

//	OutputData       []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *Plotdata_spreadsheet) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *Plotdata_spreadsheet) setup(p Plotdata_spreadsheetParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Plotdata_spreadsheet) steps(p Plotdata_spreadsheetParamBlock, r *Plotdata_spreadsheetResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	// Get some data.

	file, err := spreadsheet.OpenFile(p.Filename)

	sheet := file.Sheets[p.Sheet]

	Xdatarange, err := spreadsheet.ConvertMinMaxtoArray(p.Xminmax)
	if err != nil {
		fmt.Println(p.Xminmax, Xdatarange)
		panic(err)
	}
	fmt.Println(Xdatarange)

	Ydatarangearray := make([][]string, 0)
	for i, Yminmax := range p.Yminmaxarray {
		Ydatarange, err := spreadsheet.ConvertMinMaxtoArray(Yminmax)
		if err != nil {
			panic(err)
		}
		if len(Xdatarange) != len(Ydatarange) {
			panicmessage := fmt.Errorf("for index", i, "of array", "len(Xdatarange) != len(Ydatarange)")
			panic(panicmessage.Error())
		}
		Ydatarangearray = append(Ydatarangearray, Ydatarange)
		fmt.Println(Ydatarange)
	}

	// now plot the graph

	// the data points

	graph.PlotfromMinMaxpairs(sheet, p.Xminmax, p.Yminmaxarray, p.Exportedfilename)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Plotdata_spreadsheet) analysis(p Plotdata_spreadsheetParamBlock, r *Plotdata_spreadsheetResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Plotdata_spreadsheet) validation(p Plotdata_spreadsheetParamBlock, r *Plotdata_spreadsheetResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Plotdata_spreadsheet) Complete(params interface{}) {
	p := params.(Plotdata_spreadsheetParamBlock)
	if p.Error {
		return
	}
	r := new(Plotdata_spreadsheetResultBlock)
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
func (e *Plotdata_spreadsheet) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Plotdata_spreadsheet) NewConfig() interface{} {
	return &Plotdata_spreadsheetConfig{}
}

func (e *Plotdata_spreadsheet) NewParamBlock() interface{} {
	return &Plotdata_spreadsheetParamBlock{}
}

func NewPlotdata_spreadsheet() interface{} { //*Plotdata_spreadsheet {
	e := new(Plotdata_spreadsheet)
	e.init()
	return e
}

// Mapper function
func (e *Plotdata_spreadsheet) Map(m map[string]interface{}) interface{} {
	var res Plotdata_spreadsheetParamBlock
	res.Error = false || m["Exportedfilename"].(execute.ThreadParam).Error || m["Filename"].(execute.ThreadParam).Error || m["Sheet"].(execute.ThreadParam).Error || m["Xminmax"].(execute.ThreadParam).Error || m["Yminmaxarray"].(execute.ThreadParam).Error

	vExportedfilename, is := m["Exportedfilename"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Plotdata_spreadsheetJSONBlock
		json.Unmarshal([]byte(vExportedfilename.JSONString), &temp)
		res.Exportedfilename = *temp.Exportedfilename
	} else {
		res.Exportedfilename = m["Exportedfilename"].(execute.ThreadParam).Value.(string)
	}

	vFilename, is := m["Filename"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Plotdata_spreadsheetJSONBlock
		json.Unmarshal([]byte(vFilename.JSONString), &temp)
		res.Filename = *temp.Filename
	} else {
		res.Filename = m["Filename"].(execute.ThreadParam).Value.(string)
	}

	vSheet, is := m["Sheet"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Plotdata_spreadsheetJSONBlock
		json.Unmarshal([]byte(vSheet.JSONString), &temp)
		res.Sheet = *temp.Sheet
	} else {
		res.Sheet = m["Sheet"].(execute.ThreadParam).Value.(int)
	}

	vXminmax, is := m["Xminmax"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Plotdata_spreadsheetJSONBlock
		json.Unmarshal([]byte(vXminmax.JSONString), &temp)
		res.Xminmax = *temp.Xminmax
	} else {
		res.Xminmax = m["Xminmax"].(execute.ThreadParam).Value.([]string)
	}

	vYminmaxarray, is := m["Yminmaxarray"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Plotdata_spreadsheetJSONBlock
		json.Unmarshal([]byte(vYminmaxarray.JSONString), &temp)
		res.Yminmaxarray = *temp.Yminmaxarray
	} else {
		res.Yminmaxarray = m["Yminmaxarray"].(execute.ThreadParam).Value.([][]string)
	}

	res.ID = m["Exportedfilename"].(execute.ThreadParam).ID
	res.BlockID = m["Exportedfilename"].(execute.ThreadParam).BlockID

	return res
}

func (e *Plotdata_spreadsheet) OnExportedfilename(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
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
func (e *Plotdata_spreadsheet) OnFilename(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Filename", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Plotdata_spreadsheet) OnSheet(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Sheet", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Plotdata_spreadsheet) OnXminmax(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Xminmax", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Plotdata_spreadsheet) OnYminmaxarray(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Yminmaxarray", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Plotdata_spreadsheet struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	Exportedfilename <-chan execute.ThreadParam
	Filename         <-chan execute.ThreadParam
	Sheet            <-chan execute.ThreadParam
	Xminmax          <-chan execute.ThreadParam
	Yminmaxarray     <-chan execute.ThreadParam
}

type Plotdata_spreadsheetParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Exportedfilename string
	Filename         string
	Sheet            int
	Xminmax          []string
	Yminmaxarray     [][]string
}

type Plotdata_spreadsheetConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	Exportedfilename string
	Filename         string
	Sheet            int
	Xminmax          []string
	Yminmaxarray     [][]string
}

type Plotdata_spreadsheetResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
}

type Plotdata_spreadsheetJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	Exportedfilename *string
	Filename         *string
	Sheet            *int
	Xminmax          *[]string
	Yminmaxarray     *[][]string
}

func (c *Plotdata_spreadsheet) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Exportedfilename", "string", "Exportedfilename", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Filename", "string", "Filename", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sheet", "int", "Sheet", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Xminmax", "[]string", "Xminmax", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Yminmaxarray", "[][]string", "Yminmaxarray", true, true, nil, nil))

	ci := execute.NewComponentInfo("Plotdata_spreadsheet", "Plotdata_spreadsheet", "", false, inp, outp)

	return ci
}
