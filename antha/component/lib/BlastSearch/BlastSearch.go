// Example element demonstrating how to perform a BLAST search using the megablast algorithm

package BlastSearch

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	biogo "github.com/antha-lang/antha/internal/github.com/biogo/ncbi/blast"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol

//string //wtype.DNASequence//string
//Name string

// Data which is returned from this protocol; output data

//AnthaSeq wtype.DNASequence

// Physical inputs to this protocol

// Physical outputs from this protocol

func (e *BlastSearch) requirements() {
	_ = wunit.Make_units

}

// Actions to perform before protocol itself
func (e *BlastSearch) setup(p BlastSearchParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// Core process of the protocol: steps to be performed for each input
func (e *BlastSearch) steps(p BlastSearchParamBlock, r *BlastSearchResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	var err error
	var hits []biogo.Hit
	/*
		if Querytype == "PROTEIN" {
		hits, err = blast.MegaBlastP(Query)
		if err != nil {
			fmt.Println(err.Error())
		}

		Hits = fmt.Sprintln(blast.HitSummary(hits))


		} else if Querytype == "DNA" {
		hits, err = blast.MegaBlastN(Query)
		if err != nil {
			fmt.Println(err.Error())
		}

		Hits = fmt.Sprintln(blast.HitSummary(hits))
		}
	*/

	// Convert the sequence to an anthatype
	//AnthaSeq = wtype.MakeLinearDNASequence(Name, DNA)

	// look for orfs
	orf, orftrue := sequences.FindORF(p.AnthaSeq.Seq)

	if orftrue == true && len(orf.DNASeq) == len(p.AnthaSeq.Seq) {
		// if open reading frame is detected, we'll perform a blastP search'
		fmt.Println("ORF detected:", "full sequence length: ", len(p.AnthaSeq.Seq), "ORF length: ", len(orf.DNASeq))
		hits, err = blast.MegaBlastP(orf.ProtSeq)
	} else {
		// otherwise we'll blast the nucleotide sequence
		hits, err = p.AnthaSeq.Blast()
	}
	if err != nil {
		fmt.Println(err.Error())

	} //else {

	r.Hits = fmt.Sprintln(blast.HitSummary(hits))

	// Rename Sequence with ID of top blast hit
	p.AnthaSeq.Nm = hits[0].Id
	_ = _wrapper.WaitToEnd()

	//}

}

// Actions to perform after steps block to analyze data
func (e *BlastSearch) analysis(p BlastSearchParamBlock, r *BlastSearchResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *BlastSearch) validation(p BlastSearchParamBlock, r *BlastSearchResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *BlastSearch) Complete(params interface{}) {
	p := params.(BlastSearchParamBlock)
	if p.Error {
		e.Hits <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(BlastSearchResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Hits <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Hits <- execute.ThreadParam{Value: r.Hits, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *BlastSearch) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *BlastSearch) NewConfig() interface{} {
	return &BlastSearchConfig{}
}

func (e *BlastSearch) NewParamBlock() interface{} {
	return &BlastSearchParamBlock{}
}

func NewBlastSearch() interface{} { //*BlastSearch {
	e := new(BlastSearch)
	e.init()
	return e
}

// Mapper function
func (e *BlastSearch) Map(m map[string]interface{}) interface{} {
	var res BlastSearchParamBlock
	res.Error = false || m["AnthaSeq"].(execute.ThreadParam).Error

	vAnthaSeq, is := m["AnthaSeq"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp BlastSearchJSONBlock
		json.Unmarshal([]byte(vAnthaSeq.JSONString), &temp)
		res.AnthaSeq = *temp.AnthaSeq
	} else {
		res.AnthaSeq = m["AnthaSeq"].(execute.ThreadParam).Value.(wtype.DNASequence)
	}

	res.ID = m["AnthaSeq"].(execute.ThreadParam).ID
	res.BlockID = m["AnthaSeq"].(execute.ThreadParam).BlockID

	return res
}

func (e *BlastSearch) OnAnthaSeq(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(1, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AnthaSeq", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type BlastSearch struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	AnthaSeq       <-chan execute.ThreadParam
	Hits           chan<- execute.ThreadParam
}

type BlastSearchParamBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	AnthaSeq wtype.DNASequence
}

type BlastSearchConfig struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	AnthaSeq wtype.DNASequence
}

type BlastSearchResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Hits    string
}

type BlastSearchJSONBlock struct {
	ID       *execute.ThreadID
	BlockID  *execute.BlockID
	Error    *bool
	AnthaSeq *wtype.DNASequence
	Hits     *string
}

func (c *BlastSearch) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AnthaSeq", "wtype.DNASequence", "AnthaSeq", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Hits", "string", "Hits", true, true, nil, nil))

	ci := execute.NewComponentInfo("BlastSearch", "BlastSearch", "", false, inp, outp)

	return ci
}
