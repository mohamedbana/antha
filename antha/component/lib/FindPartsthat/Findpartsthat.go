// example protocol which allows a primitive method for searching the igem registry
// for parts with specified functions or a specified status (e.g. A = available or "Works", or results != none)
// see the igem package ("github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem")
// and igem website for more details about how to make the most of this http://parts.igem.org/Registry_API

package FindPartsthat

import (
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	"fmt"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	//	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"encoding/json"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"strings"
	"sync"
)

// Input parameters for this protocol (data)

//Constructname 				string
// e.g. promoter
// e.g. arsenic, reporter, alkane, logic gate

//RestrictionsitetoAvoid		[]string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

//Partsfound	[]wtype.DNASequence // map[string]wtype.DNASequence
//map[string][]string

// Input Requirement specification
func (e *FindPartsthat) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *FindPartsthat) setup(p FindPartsthatParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *FindPartsthat) steps(p FindPartsthatParamBlock, r *FindPartsthatResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	//var msg string
	// set warnings reported back to user to none initially
	//	warnings := make([]string,0)
	BackupParts := make([]string, 0)
	status := ""
	joinedstatus := make([]string, 0)
	// Look up parts from registry according to properties (this will take a couple of minutes the first time)

	parts := make([][]string, 0)
	subparts := make([]string, 0)

	partstatus := ""

	if p.OnlyreturnAvailableParts {
		partstatus = "A"
	}

	// first we'll parse the igem registry based on the short description contained in the fasta header for each part sequence
	for _, desc := range p.Parttypes {

		subparts = igem.FilterRegistry([]string{desc, partstatus})
		status = text.Print(desc+" :", subparts)
		joinedstatus = append(joinedstatus, status)
		parts = append(parts, subparts)
	}

	othercriteria := ""
	if p.OnlyreturnWorkingparts {
		othercriteria = "WORKS"
	}

	for i, subparts := range parts {

		partdetails := igem.LookUp(subparts)
		// now we can get detailed information of all of those records to interrogate further
		// this can be slow if there are many parts to check (~2 seconds per block of 14 parts)

		for _, subpart := range subparts {

			if strings.Contains(partdetails.Description(subpart), p.Partdescriptions[i]) &&
				strings.Contains(partdetails.Results(subpart), othercriteria) {
				BackupParts = append(BackupParts, subpart)

			}
			r.FulllistBackupParts = append(r.FulllistBackupParts, BackupParts)
		}
	}
	/*
		if len(warnings) != 0 {
		Warnings = fmt.Errorf(strings.Join(warnings,";"))
		}else{Warnings = nil}
	*/

	r.FulllistBackupParts = parts
	r.Status = strings.Join(joinedstatus, " ; ")

	// Print status
	if r.Status != "all parts available" {
		r.Status = fmt.Sprintln(r.Status)
	} else {
		r.Status = fmt.Sprintln(
			"Warnings:", r.Warnings.Error(),
			"Back up parts found (Reported to work!)", p.Parts,
			"Back up parts found (Reported to work!)", r.FulllistBackupParts,
		)
	}
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *FindPartsthat) analysis(p FindPartsthatParamBlock, r *FindPartsthatResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *FindPartsthat) validation(p FindPartsthatParamBlock, r *FindPartsthatResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *FindPartsthat) Complete(params interface{}) {
	p := params.(FindPartsthatParamBlock)
	if p.Error {
		e.FulllistBackupParts <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Warnings <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(FindPartsthatResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.FulllistBackupParts <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Warnings <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.FulllistBackupParts <- execute.ThreadParam{Value: r.FulllistBackupParts, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Warnings <- execute.ThreadParam{Value: r.Warnings, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *FindPartsthat) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *FindPartsthat) NewConfig() interface{} {
	return &FindPartsthatConfig{}
}

func (e *FindPartsthat) NewParamBlock() interface{} {
	return &FindPartsthatParamBlock{}
}

func NewFindPartsthat() interface{} { //*FindPartsthat {
	e := new(FindPartsthat)
	e.init()
	return e
}

// Mapper function
func (e *FindPartsthat) Map(m map[string]interface{}) interface{} {
	var res FindPartsthatParamBlock
	res.Error = false || m["OnlyreturnAvailableParts"].(execute.ThreadParam).Error || m["OnlyreturnWorkingparts"].(execute.ThreadParam).Error || m["Partdescriptions"].(execute.ThreadParam).Error || m["Parts"].(execute.ThreadParam).Error || m["Parttypes"].(execute.ThreadParam).Error

	vOnlyreturnAvailableParts, is := m["OnlyreturnAvailableParts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp FindPartsthatJSONBlock
		json.Unmarshal([]byte(vOnlyreturnAvailableParts.JSONString), &temp)
		res.OnlyreturnAvailableParts = *temp.OnlyreturnAvailableParts
	} else {
		res.OnlyreturnAvailableParts = m["OnlyreturnAvailableParts"].(execute.ThreadParam).Value.(bool)
	}

	vOnlyreturnWorkingparts, is := m["OnlyreturnWorkingparts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp FindPartsthatJSONBlock
		json.Unmarshal([]byte(vOnlyreturnWorkingparts.JSONString), &temp)
		res.OnlyreturnWorkingparts = *temp.OnlyreturnWorkingparts
	} else {
		res.OnlyreturnWorkingparts = m["OnlyreturnWorkingparts"].(execute.ThreadParam).Value.(bool)
	}

	vPartdescriptions, is := m["Partdescriptions"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp FindPartsthatJSONBlock
		json.Unmarshal([]byte(vPartdescriptions.JSONString), &temp)
		res.Partdescriptions = *temp.Partdescriptions
	} else {
		res.Partdescriptions = m["Partdescriptions"].(execute.ThreadParam).Value.([]string)
	}

	vParts, is := m["Parts"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp FindPartsthatJSONBlock
		json.Unmarshal([]byte(vParts.JSONString), &temp)
		res.Parts = *temp.Parts
	} else {
		res.Parts = m["Parts"].(execute.ThreadParam).Value.([][]string)
	}

	vParttypes, is := m["Parttypes"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp FindPartsthatJSONBlock
		json.Unmarshal([]byte(vParttypes.JSONString), &temp)
		res.Parttypes = *temp.Parttypes
	} else {
		res.Parttypes = m["Parttypes"].(execute.ThreadParam).Value.([]string)
	}

	res.ID = m["OnlyreturnAvailableParts"].(execute.ThreadParam).ID
	res.BlockID = m["OnlyreturnAvailableParts"].(execute.ThreadParam).BlockID

	return res
}

func (e *FindPartsthat) OnOnlyreturnAvailableParts(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OnlyreturnAvailableParts", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *FindPartsthat) OnOnlyreturnWorkingparts(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OnlyreturnWorkingparts", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *FindPartsthat) OnPartdescriptions(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Partdescriptions", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *FindPartsthat) OnParts(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Parts", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *FindPartsthat) OnParttypes(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Parttypes", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type FindPartsthat struct {
	flow.Component           // component "superclass" embedded
	lock                     sync.Mutex
	startup                  sync.Once
	params                   map[execute.ThreadID]*execute.AsyncBag
	OnlyreturnAvailableParts <-chan execute.ThreadParam
	OnlyreturnWorkingparts   <-chan execute.ThreadParam
	Partdescriptions         <-chan execute.ThreadParam
	Parts                    <-chan execute.ThreadParam
	Parttypes                <-chan execute.ThreadParam
	FulllistBackupParts      chan<- execute.ThreadParam
	Status                   chan<- execute.ThreadParam
	Warnings                 chan<- execute.ThreadParam
}

type FindPartsthatParamBlock struct {
	ID                       execute.ThreadID
	BlockID                  execute.BlockID
	Error                    bool
	OnlyreturnAvailableParts bool
	OnlyreturnWorkingparts   bool
	Partdescriptions         []string
	Parts                    [][]string
	Parttypes                []string
}

type FindPartsthatConfig struct {
	ID                       execute.ThreadID
	BlockID                  execute.BlockID
	Error                    bool
	OnlyreturnAvailableParts bool
	OnlyreturnWorkingparts   bool
	Partdescriptions         []string
	Parts                    [][]string
	Parttypes                []string
}

type FindPartsthatResultBlock struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	FulllistBackupParts [][]string
	Status              string
	Warnings            error
}

type FindPartsthatJSONBlock struct {
	ID                       *execute.ThreadID
	BlockID                  *execute.BlockID
	Error                    *bool
	OnlyreturnAvailableParts *bool
	OnlyreturnWorkingparts   *bool
	Partdescriptions         *[]string
	Parts                    *[][]string
	Parttypes                *[]string
	FulllistBackupParts      *[][]string
	Status                   *string
	Warnings                 *error
}

func (c *FindPartsthat) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("OnlyreturnAvailableParts", "bool", "OnlyreturnAvailableParts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OnlyreturnWorkingparts", "bool", "OnlyreturnWorkingparts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Partdescriptions", "[]string", "Partdescriptions", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Parts", "[][]string", "Parts", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Parttypes", "[]string", "Parttypes", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("FulllistBackupParts", "[][]string", "FulllistBackupParts", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Warnings", "error", "Warnings", true, true, nil, nil))

	ci := execute.NewComponentInfo("FindPartsthat", "FindPartsthat", "", false, inp, outp)

	return ci
}
