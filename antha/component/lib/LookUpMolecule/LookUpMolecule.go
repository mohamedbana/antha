// example of how to look up molecule properties from pubchem
package LookUpMolecule

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

// Name of compound or array of multiple compounds

// molecule type is returned consisting of name, formula, molecular weight and chemical ID (CID)

// or JSON structure if preferred

// status to be printed out in manual driver console

func (e *LookUpMolecule) requirements() {
	_ = wunit.Make_units

}
func (e *LookUpMolecule) setup(p LookUpMoleculeParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}
func (e *LookUpMolecule) steps(p LookUpMoleculeParamBlock, r *LookUpMoleculeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper

	// method of making molecule from name
	r.Compoundprops = pubchem.MakeMolecule(p.Compound)

	// or returning properties in JSON structure
	r.Jsonstring = pubchem.Compoundproperties(p.Compound)

	// method of making a list of compounds from names
	r.List = pubchem.MakeMolecules(p.Compoundlist)

	// Print out status
	r.Status = fmt.Sprintln("Returned data from",
		p.Compound, "=",
		r.Compoundprops.Moleculename,
		r.Compoundprops.MolecularWeight,
		r.Compoundprops.MolecularFormula,
		r.Compoundprops.CID,
		"Data in JSON format =", r.Jsonstring,
		"List=", r.List)
	_ = _wrapper.WaitToEnd()

}
func (e *LookUpMolecule) analysis(p LookUpMoleculeParamBlock, r *LookUpMoleculeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *LookUpMolecule) validation(p LookUpMoleculeParamBlock, r *LookUpMoleculeResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *LookUpMolecule) Complete(params interface{}) {
	p := params.(LookUpMoleculeParamBlock)
	if p.Error {
		e.Compoundprops <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Jsonstring <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.List <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(LookUpMoleculeResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Compoundprops <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Jsonstring <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.List <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Compoundprops <- execute.ThreadParam{Value: r.Compoundprops, ID: p.ID, Error: false}

	e.Jsonstring <- execute.ThreadParam{Value: r.Jsonstring, ID: p.ID, Error: false}

	e.List <- execute.ThreadParam{Value: r.List, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *LookUpMolecule) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *LookUpMolecule) NewConfig() interface{} {
	return &LookUpMoleculeConfig{}
}

func (e *LookUpMolecule) NewParamBlock() interface{} {
	return &LookUpMoleculeParamBlock{}
}

func NewLookUpMolecule() interface{} { //*LookUpMolecule {
	e := new(LookUpMolecule)
	e.init()
	return e
}

// Mapper function
func (e *LookUpMolecule) Map(m map[string]interface{}) interface{} {
	var res LookUpMoleculeParamBlock
	res.Error = false || m["Compound"].(execute.ThreadParam).Error || m["Compoundlist"].(execute.ThreadParam).Error

	vCompound, is := m["Compound"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LookUpMoleculeJSONBlock
		json.Unmarshal([]byte(vCompound.JSONString), &temp)
		res.Compound = *temp.Compound
	} else {
		res.Compound = m["Compound"].(execute.ThreadParam).Value.(string)
	}

	vCompoundlist, is := m["Compoundlist"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp LookUpMoleculeJSONBlock
		json.Unmarshal([]byte(vCompoundlist.JSONString), &temp)
		res.Compoundlist = *temp.Compoundlist
	} else {
		res.Compoundlist = m["Compoundlist"].(execute.ThreadParam).Value.([]string)
	}

	res.ID = m["Compound"].(execute.ThreadParam).ID
	res.BlockID = m["Compound"].(execute.ThreadParam).BlockID

	return res
}

func (e *LookUpMolecule) OnCompound(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Compound", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *LookUpMolecule) OnCompoundlist(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(2, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Compoundlist", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type LookUpMolecule struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	Compound       <-chan execute.ThreadParam
	Compoundlist   <-chan execute.ThreadParam
	Compoundprops  chan<- execute.ThreadParam
	Jsonstring     chan<- execute.ThreadParam
	List           chan<- execute.ThreadParam
	Status         chan<- execute.ThreadParam
}

type LookUpMoleculeParamBlock struct {
	ID           execute.ThreadID
	BlockID      execute.BlockID
	Error        bool
	Compound     string
	Compoundlist []string
}

type LookUpMoleculeConfig struct {
	ID           execute.ThreadID
	BlockID      execute.BlockID
	Error        bool
	Compound     string
	Compoundlist []string
}

type LookUpMoleculeResultBlock struct {
	ID            execute.ThreadID
	BlockID       execute.BlockID
	Error         bool
	Compoundprops pubchem.Molecule
	Jsonstring    string
	List          []pubchem.Molecule
	Status        string
}

type LookUpMoleculeJSONBlock struct {
	ID            *execute.ThreadID
	BlockID       *execute.BlockID
	Error         *bool
	Compound      *string
	Compoundlist  *[]string
	Compoundprops *pubchem.Molecule
	Jsonstring    *string
	List          *[]pubchem.Molecule
	Status        *string
}

func (c *LookUpMolecule) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Compound", "string", "Compound", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Compoundlist", "[]string", "Compoundlist", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Compoundprops", "pubchem.Molecule", "Compoundprops", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Jsonstring", "string", "Jsonstring", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("List", "[]pubchem.Molecule", "List", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))

	ci := execute.NewComponentInfo("LookUpMolecule", "LookUpMolecule", "", false, inp, outp)

	return ci
}
