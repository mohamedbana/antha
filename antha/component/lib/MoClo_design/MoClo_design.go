// This protocol is intended to design assembly parts using the MoClo assembly standard.
// Overhangs for a part are chosen according to the designated class of each part (e.g. promoter).
// The MoClo standard is hierarchical so the enzyme is chosen based on the level of assembly.
// i.e. first level 0 parts are made which may comprise of a promoter, 5prime upstream part, coding sequene, and terminator.
// Level 0 parts can then be assembled together by using level 1 enzymes and overhangs.
// currently this protocol only supports level 0 steps.
// see http://journals.plos.org/plosone/article?id=10.1371/journal.pone.0016765

package MoClo_design

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"strings"
	"sync"
)

// Input parameters for this protocol (data)

//MoClo
// of assembly standard

// labels e.g. pro = promoter

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order
// desired sequence to end up with after assembly

// Input Requirement specification
func (e *MoClo_design) requirements() {
	_ = wunit.Make_units

	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func (e *MoClo_design) setup(p MoClo_designParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *MoClo_design) steps(p MoClo_designParamBlock, r *MoClo_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	//var msg string
	// set warnings reported back to user to none initially
	warnings := make([]string, 1)
	warnings[0] = "none"

	/* find sequence data from keyword; looking it up by a given name in an inventory
	   or by biobrick ID from iGem parts registry */
	partsinorder := make([]wtype.DNASequence, 0)
	var partDNA = wtype.DNASequence{"", "", false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, ""}

	r.Status = "all parts available"
	for _, part := range p.Partsinorder {

		if strings.Contains(part, "BBa_") == true {

			partDNA.Nm = part
			partproperties := igem.LookUp([]string{part})
			partDNA.Seq = partproperties.Sequence(part)
			//partDNA.Seq = igem.GetSequence(part)

			/* We can add logic to check the status of parts too and return a warning if the part
			   is not characterised */

			if strings.Contains(partproperties.Results(part), "Works") != true {

				warnings = make([]string, 0)
				//		warning := fmt.Sprintln("iGem part", part, "results =",  igem.GetResults(part), "rating",igem.GetRating(part), "part type",igem.GetType(part), "part decription =", igem.GetDescription(part), "Categories",igem.GetCategories(part))
				warning := fmt.Sprintln("iGem part", part, "results =", partproperties.Results(part), "rating", partproperties.Rating(part), "part type", partproperties.Type(part), "part decription =", partproperties.Description(part), "Categories", partproperties.Categories(part))

				warnings = append(warnings, warning)

			}
		} else {
			partDNA = Inventory.Partslist[part]

		}

		if partDNA.Seq == "" || partDNA.Nm == "" {
			r.Status = text.Print("part: "+partDNA.Nm, partDNA.Seq+": not found in Inventory so element aborted!")
		}
		partsinorder = append(partsinorder, partDNA)
	}
	// lookup vector sequence
	vectordata := Inventory.Partslist[p.Vector]

	//lookup restriction enzyme
	restrictionenzyme := enzymes.Enzymelookup[p.AssemblyStandard][p.Level]

	// (1) Add standard overhangs using chosen assembly standard
	r.PartswithOverhangs = enzymes.MakeStandardTypeIIsassemblyParts(partsinorder, p.AssemblyStandard, p.Level, p.PartMoClotypesinorder)

	// OR (2) Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	//PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{p.Constructname, restrictionenzyme.Name, vectordata, r.PartswithOverhangs}
	status, numberofassemblies, _, newDNASequence, err := enzymes.Assemblysimulator(assembly)

	endreport := "Only run in the event of assembly failure"
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

	r.Warnings = strings.Join(warnings, ";")

	partsummary := make([]string, 0)
	for _, part := range r.PartswithOverhangs {
		partsummary = append(partsummary, text.Print(part.Nm, part.Seq))
	}

	partstoorder := text.Print("PartswithOverhangs: ", partsummary)

	// Print status
	if r.Status != "all parts available" {
		r.Status = fmt.Sprintln(r.Status)
	} else {
		r.Status = fmt.Sprintln(
			text.Print("simulator status: ", status),
			text.Print("Endreport after digestion: ", endreport),
			text.Print("Warnings:", r.Warnings),
			text.Print("Simulationpass=", r.Simulationpass),
			text.Print("NewDNASequence: ", r.NewDNASequence),
			partstoorder,
		)
	}
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *MoClo_design) analysis(p MoClo_designParamBlock, r *MoClo_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *MoClo_design) validation(p MoClo_designParamBlock, r *MoClo_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *MoClo_design) Complete(params interface{}) {
	p := params.(MoClo_designParamBlock)
	if p.Error {
		e.NewDNASequence <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.PartswithOverhangs <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Simulationpass <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Warnings <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(MoClo_designResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.NewDNASequence <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.PartswithOverhangs <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Simulationpass <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Warnings <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.NewDNASequence <- execute.ThreadParam{Value: r.NewDNASequence, ID: p.ID, Error: false}

	e.PartswithOverhangs <- execute.ThreadParam{Value: r.PartswithOverhangs, ID: p.ID, Error: false}

	e.Simulationpass <- execute.ThreadParam{Value: r.Simulationpass, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Warnings <- execute.ThreadParam{Value: r.Warnings, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *MoClo_design) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *MoClo_design) NewConfig() interface{} {
	return &MoClo_designConfig{}
}

func (e *MoClo_design) NewParamBlock() interface{} {
	return &MoClo_designParamBlock{}
}

func NewMoClo_design() interface{} { //*MoClo_design {
	e := new(MoClo_design)
	e.init()
	return e
}

// Mapper function
func (e *MoClo_design) Map(m map[string]interface{}) interface{} {
	var res MoClo_designParamBlock
	res.Error = false || m["AssemblyStandard"].(execute.ThreadParam).Error || m["Constructname"].(execute.ThreadParam).Error || m["Level"].(execute.ThreadParam).Error || m["PartMoClotypesinorder"].(execute.ThreadParam).Error || m["Partsinorder"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error

	vAssemblyStandard, is := m["AssemblyStandard"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MoClo_designJSONBlock
		json.Unmarshal([]byte(vAssemblyStandard.JSONString), &temp)
		res.AssemblyStandard = *temp.AssemblyStandard
	} else {
		res.AssemblyStandard = m["AssemblyStandard"].(execute.ThreadParam).Value.(string)
	}

	vConstructname, is := m["Constructname"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MoClo_designJSONBlock
		json.Unmarshal([]byte(vConstructname.JSONString), &temp)
		res.Constructname = *temp.Constructname
	} else {
		res.Constructname = m["Constructname"].(execute.ThreadParam).Value.(string)
	}

	vLevel, is := m["Level"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MoClo_designJSONBlock
		json.Unmarshal([]byte(vLevel.JSONString), &temp)
		res.Level = *temp.Level
	} else {
		res.Level = m["Level"].(execute.ThreadParam).Value.(string)
	}

	vPartMoClotypesinorder, is := m["PartMoClotypesinorder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MoClo_designJSONBlock
		json.Unmarshal([]byte(vPartMoClotypesinorder.JSONString), &temp)
		res.PartMoClotypesinorder = *temp.PartMoClotypesinorder
	} else {
		res.PartMoClotypesinorder = m["PartMoClotypesinorder"].(execute.ThreadParam).Value.([]string)
	}

	vPartsinorder, is := m["Partsinorder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MoClo_designJSONBlock
		json.Unmarshal([]byte(vPartsinorder.JSONString), &temp)
		res.Partsinorder = *temp.Partsinorder
	} else {
		res.Partsinorder = m["Partsinorder"].(execute.ThreadParam).Value.([]string)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp MoClo_designJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["AssemblyStandard"].(execute.ThreadParam).ID
	res.BlockID = m["AssemblyStandard"].(execute.ThreadParam).BlockID

	return res
}

func (e *MoClo_design) OnAssemblyStandard(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AssemblyStandard", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MoClo_design) OnConstructname(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
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
func (e *MoClo_design) OnLevel(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Level", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MoClo_design) OnPartMoClotypesinorder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("PartMoClotypesinorder", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MoClo_design) OnPartsinorder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Partsinorder", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *MoClo_design) OnVector(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
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

type MoClo_design struct {
	flow.Component        // component "superclass" embedded
	lock                  sync.Mutex
	startup               sync.Once
	params                map[execute.ThreadID]*execute.AsyncBag
	AssemblyStandard      <-chan execute.ThreadParam
	Constructname         <-chan execute.ThreadParam
	Level                 <-chan execute.ThreadParam
	PartMoClotypesinorder <-chan execute.ThreadParam
	Partsinorder          <-chan execute.ThreadParam
	Vector                <-chan execute.ThreadParam
	NewDNASequence        chan<- execute.ThreadParam
	PartswithOverhangs    chan<- execute.ThreadParam
	Simulationpass        chan<- execute.ThreadParam
	Status                chan<- execute.ThreadParam
	Warnings              chan<- execute.ThreadParam
}

type MoClo_designParamBlock struct {
	ID                    execute.ThreadID
	BlockID               execute.BlockID
	Error                 bool
	AssemblyStandard      string
	Constructname         string
	Level                 string
	PartMoClotypesinorder []string
	Partsinorder          []string
	Vector                string
}

type MoClo_designConfig struct {
	ID                    execute.ThreadID
	BlockID               execute.BlockID
	Error                 bool
	AssemblyStandard      string
	Constructname         string
	Level                 string
	PartMoClotypesinorder []string
	Partsinorder          []string
	Vector                string
}

type MoClo_designResultBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Status             string
	Warnings           string
}

type MoClo_designJSONBlock struct {
	ID                    *execute.ThreadID
	BlockID               *execute.BlockID
	Error                 *bool
	AssemblyStandard      *string
	Constructname         *string
	Level                 *string
	PartMoClotypesinorder *[]string
	Partsinorder          *[]string
	Vector                *string
	NewDNASequence        *wtype.DNASequence
	PartswithOverhangs    *[]wtype.DNASequence
	Simulationpass        *bool
	Status                *string
	Warnings              *string
}

func (c *MoClo_design) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AssemblyStandard", "string", "AssemblyStandard", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Constructname", "string", "Constructname", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Level", "string", "Level", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartMoClotypesinorder", "[]string", "PartMoClotypesinorder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Partsinorder", "[]string", "Partsinorder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "string", "Vector", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("NewDNASequence", "wtype.DNASequence", "NewDNASequence", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("PartswithOverhangs", "[]wtype.DNASequence", "PartswithOverhangs", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Simulationpass", "bool", "Simulationpass", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Warnings", "string", "Warnings", true, true, nil, nil))

	ci := execute.NewComponentInfo("MoClo_design", "MoClo_design", "", false, inp, outp)

	return ci
}
