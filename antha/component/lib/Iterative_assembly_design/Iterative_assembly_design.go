// This protocol is based on scarfree design so please look at that first.
// The protocol is intended to design assembly parts using the first enzyme
// which is found to be feasible to use from a list of ApprovedEnzymes enzymes . If no enzyme
// from the list is feasible to use (i.e. due to the presence of existing restriction sites in a part)
// all typeIIs enzymes will be screened to find feasible backup options

package Iterative_assembly_design

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func (e *Iterative_assembly_design) requirements() {
	_ = wunit.Make_units

	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func (e *Iterative_assembly_design) setup(p Iterative_assembly_designParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Iterative_assembly_design) steps(p Iterative_assembly_designParamBlock, r *Iterative_assembly_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	//var msg string
	// set warnings reported back to user to none initially

	warnings := make([]string, 0)
	sitefound := false
	Enzyme := "No enzymes which passed with these sequences"
	// make an empty array of DNA Sequences ready to fill
	partsinorder := make([]wtype.DNASequence, 0)

	r.Status = "all parts available"
	for i, part := range p.Seqsinorder {
		if strings.Contains(part, "BBa_") {
			part = igem.GetSequence(part)
		}
		partDNA := wtype.MakeLinearDNASequence("Part "+strconv.Itoa(i), part)

		partsinorder = append(partsinorder, partDNA)
	}
	// Find all possible typeIIs enzymes we could use for these sequences (i.e. non cutters of all parts)
	possibilities := lookup.FindEnzymeNamesofClass("TypeIIs")
	var backupoption string
	for _, possibility := range possibilities {
		// check number of sites per part !
		enz := lookup.EnzymeLookup(possibility)

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})
			if len(info) != 0 {
				if info[0].Sitefound == true {
					sitefound = true
					break
				}
			}
		}
		if sitefound == false {
			backupoption = possibility
			r.BackupEnzymes = append(r.BackupEnzymes, backupoption)
		}
	}

	sitefound = false
	for _, Enzyme := range p.ApprovedEnzymes {

		// check number of sites per part !
		enz := lookup.EnzymeLookup(Enzyme)

		for _, part := range partsinorder {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})
			if len(info) != 0 {
				if info[0].Sitefound == true {
					sitefound = true
					break
				}
			}
		}
		if sitefound == false {
			r.EnzymeUsed = enz
		}
	}

	if sitefound != true {
		fmt.Println("enzyme used", r.EnzymeUsed)
		Enzyme = r.EnzymeUsed.Name

		// make vector into an antha type DNASequence
		vectordata := wtype.MakePlasmidDNASequence("Vector", p.Vector)

		//lookup restriction enzyme
		restrictionenzyme, err := lookup.TypeIIsLookup(r.EnzymeUsed.Name)
		if err != nil {
			text.Print("Error", err.Error())
		}

		//  Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
		r.PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

		// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
		assembly := enzymes.Assemblyparameters{p.Constructname, restrictionenzyme.Name, vectordata, r.PartswithOverhangs}
		status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

		endreport := "Endreport only run in the event of assembly simulation failure"
		//sites := "Restriction mapper only run in the event of assembly simulation failure"
		r.NewDNASequence = newDNASequence
		if err == nil && numberofassemblies == 1 {

			r.Simulationpass = true
		} else {

			warnings = append(warnings, status)
			// perform mock digest to test fragement overhangs (fragments are hidden by using _, )
			_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

			allends := make([]string, 0)
			ends := ""

			ends = text.Print(vectordata.Nm+" 5 Prime end: ", stickyends5)
			allends = append(allends, ends)
			ends = text.Print(vectordata.Nm+" 3 Prime end: ", stickyends3)
			allends = append(allends, ends)

			for _, part := range r.PartswithOverhangs {
				_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(part, restrictionenzyme)
				ends = text.Print(part.Nm+" 5 Prime end: ", stickyends5)
				allends = append(allends, ends)
				ends = text.Print(part.Nm+" 3 Prime end: ", stickyends3)
				allends = append(allends, ends)
			}
			endreport = strings.Join(allends, " ")
		}

		// check number of sites per part !
		enz := lookup.EnzymeLookup(Enzyme)
		sites := make([]int, 0)
		multiple := make([]string, 0)
		for _, part := range r.PartswithOverhangs {

			info := enzymes.Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})

			sitepositions := enzymes.SitepositionString(info[0])

			sites = append(sites, info[0].Numberofsites)
			sitepositions = text.Print(part.Nm+" "+Enzyme+" positions:", sitepositions)
			multiple = append(multiple, sitepositions)
		}

		if len(warnings) == 0 {
			warnings = append(warnings, "none")
		}
		r.Warnings = fmt.Errorf(strings.Join(warnings, ";"))

		partsummary := make([]string, 0)
		for _, part := range r.PartswithOverhangs {
			partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
		}

		partstoorder := text.Print("PartswithOverhangs: ", partsummary)

		r.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Sites per part for "+Enzyme, sites),
			text.Print("Positions: ", multiple),
			text.Print("Warnings:", r.Warnings.Error()),
			text.Print("Simulationpass=", r.Simulationpass),
			text.Print("NewDNASequence: ", r.NewDNASequence),
			partstoorder)

	}
	// Print status
	if r.Status != "all parts available" {
		r.Status = fmt.Sprintln(r.Status,
			text.Print("Backup Enzymes: ", r.BackupEnzymes))
	} else if sitefound == true {
		r.Status = fmt.Sprintln(text.Print("No Enzyme found to be compatible from approved list", p.ApprovedEnzymes),
			text.Print("Backup Enzymes: ", r.BackupEnzymes))

	} else {
		r.Status = fmt.Sprintln(r.Status,
			text.Print("Backup Enzymes: ", r.BackupEnzymes))

	}
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Iterative_assembly_design) analysis(p Iterative_assembly_designParamBlock, r *Iterative_assembly_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Iterative_assembly_design) validation(p Iterative_assembly_designParamBlock, r *Iterative_assembly_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Iterative_assembly_design) Complete(params interface{}) {
	p := params.(Iterative_assembly_designParamBlock)
	if p.Error {
		e.BackupEnzymes <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.EnzymeUsed <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.NewDNASequence <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.PartswithOverhangs <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Simulationpass <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Warnings <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(Iterative_assembly_designResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.BackupEnzymes <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.EnzymeUsed <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.NewDNASequence <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.PartswithOverhangs <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Simulationpass <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Warnings <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.BackupEnzymes <- execute.ThreadParam{Value: r.BackupEnzymes, ID: p.ID, Error: false}

	e.EnzymeUsed <- execute.ThreadParam{Value: r.EnzymeUsed, ID: p.ID, Error: false}

	e.NewDNASequence <- execute.ThreadParam{Value: r.NewDNASequence, ID: p.ID, Error: false}

	e.PartswithOverhangs <- execute.ThreadParam{Value: r.PartswithOverhangs, ID: p.ID, Error: false}

	e.Simulationpass <- execute.ThreadParam{Value: r.Simulationpass, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Warnings <- execute.ThreadParam{Value: r.Warnings, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Iterative_assembly_design) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Iterative_assembly_design) NewConfig() interface{} {
	return &Iterative_assembly_designConfig{}
}

func (e *Iterative_assembly_design) NewParamBlock() interface{} {
	return &Iterative_assembly_designParamBlock{}
}

func NewIterative_assembly_design() interface{} { //*Iterative_assembly_design {
	e := new(Iterative_assembly_design)
	e.init()
	return e
}

// Mapper function
func (e *Iterative_assembly_design) Map(m map[string]interface{}) interface{} {
	var res Iterative_assembly_designParamBlock
	res.Error = false || m["ApprovedEnzymes"].(execute.ThreadParam).Error || m["Constructname"].(execute.ThreadParam).Error || m["Seqsinorder"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error

	vApprovedEnzymes, is := m["ApprovedEnzymes"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Iterative_assembly_designJSONBlock
		json.Unmarshal([]byte(vApprovedEnzymes.JSONString), &temp)
		res.ApprovedEnzymes = *temp.ApprovedEnzymes
	} else {
		res.ApprovedEnzymes = m["ApprovedEnzymes"].(execute.ThreadParam).Value.([]string)
	}

	vConstructname, is := m["Constructname"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Iterative_assembly_designJSONBlock
		json.Unmarshal([]byte(vConstructname.JSONString), &temp)
		res.Constructname = *temp.Constructname
	} else {
		res.Constructname = m["Constructname"].(execute.ThreadParam).Value.(string)
	}

	vSeqsinorder, is := m["Seqsinorder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Iterative_assembly_designJSONBlock
		json.Unmarshal([]byte(vSeqsinorder.JSONString), &temp)
		res.Seqsinorder = *temp.Seqsinorder
	} else {
		res.Seqsinorder = m["Seqsinorder"].(execute.ThreadParam).Value.([]string)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp Iterative_assembly_designJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["ApprovedEnzymes"].(execute.ThreadParam).ID
	res.BlockID = m["ApprovedEnzymes"].(execute.ThreadParam).BlockID

	return res
}

func (e *Iterative_assembly_design) OnApprovedEnzymes(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ApprovedEnzymes", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Iterative_assembly_design) OnConstructname(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Constructname", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Iterative_assembly_design) OnSeqsinorder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Seqsinorder", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Iterative_assembly_design) OnVector(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(4, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vector", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Iterative_assembly_design struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	ApprovedEnzymes    <-chan execute.ThreadParam
	Constructname      <-chan execute.ThreadParam
	Seqsinorder        <-chan execute.ThreadParam
	Vector             <-chan execute.ThreadParam
	BackupEnzymes      chan<- execute.ThreadParam
	EnzymeUsed         chan<- execute.ThreadParam
	NewDNASequence     chan<- execute.ThreadParam
	PartswithOverhangs chan<- execute.ThreadParam
	Simulationpass     chan<- execute.ThreadParam
	Status             chan<- execute.ThreadParam
	Warnings           chan<- execute.ThreadParam
}

type Iterative_assembly_designParamBlock struct {
	ID              execute.ThreadID
	BlockID         execute.BlockID
	Error           bool
	ApprovedEnzymes []string
	Constructname   string
	Seqsinorder     []string
	Vector          string
}

type Iterative_assembly_designConfig struct {
	ID              execute.ThreadID
	BlockID         execute.BlockID
	Error           bool
	ApprovedEnzymes []string
	Constructname   string
	Seqsinorder     []string
	Vector          string
}

type Iterative_assembly_designResultBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	BackupEnzymes      []string
	EnzymeUsed         wtype.RestrictionEnzyme
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           error
}

type Iterative_assembly_designJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	ApprovedEnzymes    *[]string
	Constructname      *string
	Seqsinorder        *[]string
	Vector             *string
	BackupEnzymes      *[]string
	EnzymeUsed         *wtype.RestrictionEnzyme
	NewDNASequence     *wtype.DNASequence
	PartswithOverhangs *[]wtype.DNASequence
	Simulationpass     *bool
	Status             *string
	Warnings           *error
}

func (c *Iterative_assembly_design) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("ApprovedEnzymes", "[]string", "ApprovedEnzymes", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Constructname", "string", "Constructname", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Seqsinorder", "[]string", "Seqsinorder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "string", "Vector", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("BackupEnzymes", "[]string", "BackupEnzymes", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("EnzymeUsed", "wtype.RestrictionEnzyme", "EnzymeUsed", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("NewDNASequence", "wtype.DNASequence", "NewDNASequence", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("PartswithOverhangs", "[]wtype.DNASequence", "PartswithOverhangs", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Simulationpass", "bool", "Simulationpass", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Warnings", "error", "Warnings", true, true, nil, nil))

	ci := execute.NewComponentInfo("Iterative_assembly_design", "Iterative_assembly_design", "", false, inp, outp)

	return ci
}
