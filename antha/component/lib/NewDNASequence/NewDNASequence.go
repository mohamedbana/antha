package NewDNASequence

import (
	"fmt"
	//"math"
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func (e *NewDNASequence) requirements() {
	_ = wunit.Make_units

}

// Actions to perform before protocol itself
func (e *NewDNASequence) setup(p NewDNASequenceParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// Core process of the protocol: steps to be performed for each input
func (e *NewDNASequence) steps(p NewDNASequenceParamBlock, r *NewDNASequenceResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	if p.Plasmid != p.Linear {
		if p.Plasmid {
			r.DNA = wtype.MakePlasmidDNASequence(p.Gene_name, p.DNA_seq)
		} else if p.Linear {
			r.DNA = wtype.MakeLinearDNASequence(p.Gene_name, p.DNA_seq)
		} else if p.SingleStranded {
			r.DNA = wtype.MakeSingleStrandedDNASequence(p.Gene_name, p.DNA_seq)
		}
	} else {
		fmt.Println("correct conditions not met")
	}
	_ = _wrapper.WaitToEnd()

}

// Actions to perform after steps block to analyze data
func (e *NewDNASequence) analysis(p NewDNASequenceParamBlock, r *NewDNASequenceResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *NewDNASequence) validation(p NewDNASequenceParamBlock, r *NewDNASequenceResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *NewDNASequence) Complete(params interface{}) {
	p := params.(NewDNASequenceParamBlock)
	if p.Error {
		e.DNA <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(NewDNASequenceResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.DNA <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.DNA <- execute.ThreadParam{Value: r.DNA, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *NewDNASequence) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *NewDNASequence) NewConfig() interface{} {
	return &NewDNASequenceConfig{}
}

func (e *NewDNASequence) NewParamBlock() interface{} {
	return &NewDNASequenceParamBlock{}
}

func NewNewDNASequence() interface{} { //*NewDNASequence {
	e := new(NewDNASequence)
	e.init()
	return e
}

// Mapper function
func (e *NewDNASequence) Map(m map[string]interface{}) interface{} {
	var res NewDNASequenceParamBlock
	res.Error = false || m["DNA_seq"].(execute.ThreadParam).Error || m["Gene_name"].(execute.ThreadParam).Error || m["Linear"].(execute.ThreadParam).Error || m["Plasmid"].(execute.ThreadParam).Error || m["SingleStranded"].(execute.ThreadParam).Error

	vDNA_seq, is := m["DNA_seq"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp NewDNASequenceJSONBlock
		json.Unmarshal([]byte(vDNA_seq.JSONString), &temp)
		res.DNA_seq = *temp.DNA_seq
	} else {
		res.DNA_seq = m["DNA_seq"].(execute.ThreadParam).Value.(string)
	}

	vGene_name, is := m["Gene_name"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp NewDNASequenceJSONBlock
		json.Unmarshal([]byte(vGene_name.JSONString), &temp)
		res.Gene_name = *temp.Gene_name
	} else {
		res.Gene_name = m["Gene_name"].(execute.ThreadParam).Value.(string)
	}

	vLinear, is := m["Linear"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp NewDNASequenceJSONBlock
		json.Unmarshal([]byte(vLinear.JSONString), &temp)
		res.Linear = *temp.Linear
	} else {
		res.Linear = m["Linear"].(execute.ThreadParam).Value.(bool)
	}

	vPlasmid, is := m["Plasmid"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp NewDNASequenceJSONBlock
		json.Unmarshal([]byte(vPlasmid.JSONString), &temp)
		res.Plasmid = *temp.Plasmid
	} else {
		res.Plasmid = m["Plasmid"].(execute.ThreadParam).Value.(bool)
	}

	vSingleStranded, is := m["SingleStranded"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp NewDNASequenceJSONBlock
		json.Unmarshal([]byte(vSingleStranded.JSONString), &temp)
		res.SingleStranded = *temp.SingleStranded
	} else {
		res.SingleStranded = m["SingleStranded"].(execute.ThreadParam).Value.(bool)
	}

	res.ID = m["DNA_seq"].(execute.ThreadParam).ID
	res.BlockID = m["DNA_seq"].(execute.ThreadParam).BlockID

	return res
}

func (e *NewDNASequence) OnDNA_seq(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNA_seq", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *NewDNASequence) OnGene_name(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Gene_name", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *NewDNASequence) OnLinear(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Linear", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *NewDNASequence) OnPlasmid(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Plasmid", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *NewDNASequence) OnSingleStranded(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(5, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SingleStranded", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type NewDNASequence struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	DNA_seq        <-chan execute.ThreadParam
	Gene_name      <-chan execute.ThreadParam
	Linear         <-chan execute.ThreadParam
	Plasmid        <-chan execute.ThreadParam
	SingleStranded <-chan execute.ThreadParam
	DNA            chan<- execute.ThreadParam
}

type NewDNASequenceParamBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	DNA_seq        string
	Gene_name      string
	Linear         bool
	Plasmid        bool
	SingleStranded bool
}

type NewDNASequenceConfig struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	DNA_seq        string
	Gene_name      string
	Linear         bool
	Plasmid        bool
	SingleStranded bool
}

type NewDNASequenceResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	DNA     wtype.DNASequence
}

type NewDNASequenceJSONBlock struct {
	ID             *execute.ThreadID
	BlockID        *execute.BlockID
	Error          *bool
	DNA_seq        *string
	Gene_name      *string
	Linear         *bool
	Plasmid        *bool
	SingleStranded *bool
	DNA            *wtype.DNASequence
}

func (c *NewDNASequence) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("DNA_seq", "string", "DNA_seq", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Gene_name", "string", "Gene_name", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Linear", "bool", "Linear", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Plasmid", "bool", "Plasmid", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SingleStranded", "bool", "SingleStranded", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("DNA", "wtype.DNASequence", "DNA", true, true, nil, nil))

	ci := execute.NewComponentInfo("NewDNASequence", "NewDNASequence", "", false, inp, outp)

	return ci
}
