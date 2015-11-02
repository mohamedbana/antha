package TypeIISAssembly_design

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/REBASE"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"strconv"
	"strings"
	"sync"
)

//"github.com/mgutz/ansi"

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func (e *TypeIISAssembly_design) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *TypeIISAssembly_design) setup(p TypeIISAssembly_designParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *TypeIISAssembly_design) steps(p TypeIISAssembly_designParamBlock, r *TypeIISAssembly_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
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
			partDNA.Seq = igem.GetSequence(part)

			/* We can add logic to check the status of parts too and return a warning if the part
			   is not characeterised */

			if strings.Contains(igem.GetResults(part), "Works") != true {

				warnings = make([]string, 0)
				warning := fmt.Sprintln("iGem part", part, "results =", igem.GetResults(part), "rating", igem.GetRating(part), "part type", igem.GetType(part), "part decription =", igem.GetDescription(part), "Categories", igem.GetCategories(part))
				warnings = append(warnings, warning)

			}
		} else {
			partDNA = Inventory.Partslist[part]

		}

		if partDNA.Seq == "" || partDNA.Nm == "" {
			r.Status = fmt.Sprintln("part not found in Inventory so element aborted!")
		}
		partsinorder = append(partsinorder, partDNA)
	}

	// or Look up parts from registry according to properties (this will take a couple of minutes the first time)
	subparts := igem.FilterRegistry([]string{"Fluorescent", "A "})
	partdetails := igem.PartProperties(subparts)
	//fmt.Println(partdetails)

	// this can be slow if there are many parts to check (~2 seconds per block of 14 parts)
	for _, subpart := range subparts {
		if strings.Contains(igem.GetDescriptionfromSubset(subpart, partdetails), "RED") &&
			strings.Contains(igem.GetResultsfromSubset(subpart, partdetails), "WORKS") {
			r.BackupParts = append(r.BackupParts, subpart)

		}
	}

	// lookup vector sequence
	vectordata := Inventory.Partslist[p.Vector]

	//lookup restriction enzyme
	restrictionenzyme := enzymes.Enzymelookup[p.AssemblyStandard][p.Level]

	// (1) Add standard overhangs using chosen assembly standard
	//PartswithOverhangs = enzymes.MakeStandardTypeIIsassemblyParts(partsinorder, AssemblyStandard, Level, PartMoClotypesinorder)

	// OR (2) Add overhangs for scarfree assembly based on part seqeunces only, i.e. no Assembly standard
	r.PartswithOverhangs = enzymes.MakeScarfreeCustomTypeIIsassemblyParts(partsinorder, vectordata, restrictionenzyme)

	// perfrom mock digest to test fragement overhangs (fragments are hidden by using _, )
	_, stickyends5, stickyends3 := enzymes.TypeIIsdigest(vectordata, restrictionenzyme)

	// Check that assembly is feasible with designed parts by simulating assembly of the sequences with the chosen enzyme
	assembly := enzymes.Assemblyparameters{p.Constructname, restrictionenzyme.Name, vectordata, r.PartswithOverhangs}
	status, numberofassemblies, sitesfound, newDNASequence, _ := enzymes.Assemblysimulator(assembly)

	// The default sitesfound produced from the assembly simulator only checks to SapI and BsaI so we'll repeat with the enzymes declared in parameters
	// first lookup enzyme properties
	enzlist := make([]wtype.LogicalRestrictionEnzyme, 0)
	for _, site := range p.RestrictionsitetoAvoid {
		enzsite := rebase.EnzymeLookup(site)
		enzlist = append(enzlist, enzsite)
	}
	othersitesfound := enzymes.Restrictionsitefinder(newDNASequence, enzlist)

	for _, site := range sitesfound {
		othersitesfound = append(othersitesfound, site)
	}

	// Now let's find out the size of fragments we would get if digested with a common site cutter
	tspEI := rebase.EnzymeLookup("TspEI")

	Testdigestionsizes := enzymes.RestrictionMapper(newDNASequence, tspEI)

	// allow the data to be exported by capitalising the first letter of the variable
	r.Sitesfound = othersitesfound

	r.NewDNASequence = newDNASequence
	if status == "Yay! this should work" && numberofassemblies == 1 {

		r.Simulationpass = true
	}

	r.Warnings = strings.Join(warnings, ";")

	// Export sequences to order into a fasta file

	partswithOverhangs := make([]*wtype.DNASequence, 0)
	for i, part := range r.PartswithOverhangs {
		_ = enzymes.ExportFastaDir(p.Constructname, strconv.Itoa(i+1), &part)
		partswithOverhangs = append(partswithOverhangs, &part)

	}
	_ = enzymes.Makefastaserial(p.Constructname, partswithOverhangs)

	//partstoorder := ansi.Color(fmt.Sprintln("PartswithOverhangs", PartswithOverhangs),"red")
	partstoorder := fmt.Sprintln("PartswithOverhangs", r.PartswithOverhangs)

	// Print status
	if r.Status != "all parts available" {
		r.Status = fmt.Sprintln(r.Status)
	} else {
		r.Status = fmt.Sprintln(
			"Warnings:", r.Warnings,
			"Simulationpass=", r.Simulationpass,
			"Back up parts found (which work)", r.BackupParts,
			"NewDNASequence", r.NewDNASequence,
			//"partonewithoverhangs", partonewithoverhangs,
			//"Vector",vectordata,
			"Vector digest:", stickyends5, stickyends3,
			partstoorder,
			"Sitesfound", r.Sitesfound,
			"Partsinorder=", p.Partsinorder, partsinorder,
			"Test digestion sizes with TspEI", Testdigestionsizes,
		//"Restriction Enzyme=",restrictionenzyme,
		)
	}
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *TypeIISAssembly_design) analysis(p TypeIISAssembly_designParamBlock, r *TypeIISAssembly_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *TypeIISAssembly_design) validation(p TypeIISAssembly_designParamBlock, r *TypeIISAssembly_designResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *TypeIISAssembly_design) Complete(params interface{}) {
	p := params.(TypeIISAssembly_designParamBlock)
	if p.Error {
		e.BackupParts <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.NewDNASequence <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.PartswithOverhangs <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Simulationpass <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Sitesfound <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Warnings <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(TypeIISAssembly_designResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.BackupParts <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.NewDNASequence <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.PartswithOverhangs <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Simulationpass <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Sitesfound <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Warnings <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.BackupParts <- execute.ThreadParam{Value: r.BackupParts, ID: p.ID, Error: false}

	e.NewDNASequence <- execute.ThreadParam{Value: r.NewDNASequence, ID: p.ID, Error: false}

	e.PartswithOverhangs <- execute.ThreadParam{Value: r.PartswithOverhangs, ID: p.ID, Error: false}

	e.Simulationpass <- execute.ThreadParam{Value: r.Simulationpass, ID: p.ID, Error: false}

	e.Sitesfound <- execute.ThreadParam{Value: r.Sitesfound, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.Warnings <- execute.ThreadParam{Value: r.Warnings, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *TypeIISAssembly_design) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *TypeIISAssembly_design) NewConfig() interface{} {
	return &TypeIISAssembly_designConfig{}
}

func (e *TypeIISAssembly_design) NewParamBlock() interface{} {
	return &TypeIISAssembly_designParamBlock{}
}

func NewTypeIISAssembly_design() interface{} { //*TypeIISAssembly_design {
	e := new(TypeIISAssembly_design)
	e.init()
	return e
}

// Mapper function
func (e *TypeIISAssembly_design) Map(m map[string]interface{}) interface{} {
	var res TypeIISAssembly_designParamBlock
	res.Error = false || m["AssemblyStandard"].(execute.ThreadParam).Error || m["Constructname"].(execute.ThreadParam).Error || m["Level"].(execute.ThreadParam).Error || m["PartMoClotypesinorder"].(execute.ThreadParam).Error || m["Partsinorder"].(execute.ThreadParam).Error || m["RestrictionsitetoAvoid"].(execute.ThreadParam).Error || m["Vector"].(execute.ThreadParam).Error

	vAssemblyStandard, is := m["AssemblyStandard"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vAssemblyStandard.JSONString), &temp)
		res.AssemblyStandard = *temp.AssemblyStandard
	} else {
		res.AssemblyStandard = m["AssemblyStandard"].(execute.ThreadParam).Value.(string)
	}

	vConstructname, is := m["Constructname"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vConstructname.JSONString), &temp)
		res.Constructname = *temp.Constructname
	} else {
		res.Constructname = m["Constructname"].(execute.ThreadParam).Value.(string)
	}

	vLevel, is := m["Level"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vLevel.JSONString), &temp)
		res.Level = *temp.Level
	} else {
		res.Level = m["Level"].(execute.ThreadParam).Value.(string)
	}

	vPartMoClotypesinorder, is := m["PartMoClotypesinorder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vPartMoClotypesinorder.JSONString), &temp)
		res.PartMoClotypesinorder = *temp.PartMoClotypesinorder
	} else {
		res.PartMoClotypesinorder = m["PartMoClotypesinorder"].(execute.ThreadParam).Value.([]string)
	}

	vPartsinorder, is := m["Partsinorder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vPartsinorder.JSONString), &temp)
		res.Partsinorder = *temp.Partsinorder
	} else {
		res.Partsinorder = m["Partsinorder"].(execute.ThreadParam).Value.([]string)
	}

	vRestrictionsitetoAvoid, is := m["RestrictionsitetoAvoid"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vRestrictionsitetoAvoid.JSONString), &temp)
		res.RestrictionsitetoAvoid = *temp.RestrictionsitetoAvoid
	} else {
		res.RestrictionsitetoAvoid = m["RestrictionsitetoAvoid"].(execute.ThreadParam).Value.([]string)
	}

	vVector, is := m["Vector"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp TypeIISAssembly_designJSONBlock
		json.Unmarshal([]byte(vVector.JSONString), &temp)
		res.Vector = *temp.Vector
	} else {
		res.Vector = m["Vector"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["AssemblyStandard"].(execute.ThreadParam).ID
	res.BlockID = m["AssemblyStandard"].(execute.ThreadParam).BlockID

	return res
}

func (e *TypeIISAssembly_design) OnAssemblyStandard(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *TypeIISAssembly_design) OnConstructname(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *TypeIISAssembly_design) OnLevel(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *TypeIISAssembly_design) OnPartMoClotypesinorder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *TypeIISAssembly_design) OnPartsinorder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *TypeIISAssembly_design) OnRestrictionsitetoAvoid(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("RestrictionsitetoAvoid", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *TypeIISAssembly_design) OnVector(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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

type TypeIISAssembly_design struct {
	flow.Component         // component "superclass" embedded
	lock                   sync.Mutex
	startup                sync.Once
	params                 map[execute.ThreadID]*execute.AsyncBag
	AssemblyStandard       <-chan execute.ThreadParam
	Constructname          <-chan execute.ThreadParam
	Level                  <-chan execute.ThreadParam
	PartMoClotypesinorder  <-chan execute.ThreadParam
	Partsinorder           <-chan execute.ThreadParam
	RestrictionsitetoAvoid <-chan execute.ThreadParam
	Vector                 <-chan execute.ThreadParam
	BackupParts            chan<- execute.ThreadParam
	NewDNASequence         chan<- execute.ThreadParam
	PartswithOverhangs     chan<- execute.ThreadParam
	Simulationpass         chan<- execute.ThreadParam
	Sitesfound             chan<- execute.ThreadParam
	Status                 chan<- execute.ThreadParam
	Warnings               chan<- execute.ThreadParam
}

type TypeIISAssembly_designParamBlock struct {
	ID                     execute.ThreadID
	BlockID                execute.BlockID
	Error                  bool
	AssemblyStandard       string
	Constructname          string
	Level                  string
	PartMoClotypesinorder  []string
	Partsinorder           []string
	RestrictionsitetoAvoid []string
	Vector                 string
}

type TypeIISAssembly_designConfig struct {
	ID                     execute.ThreadID
	BlockID                execute.BlockID
	Error                  bool
	AssemblyStandard       string
	Constructname          string
	Level                  string
	PartMoClotypesinorder  []string
	Partsinorder           []string
	RestrictionsitetoAvoid []string
	Vector                 string
}

type TypeIISAssembly_designResultBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	BackupParts        []string
	NewDNASequence     wtype.DNASequence
	PartswithOverhangs []wtype.DNASequence
	Simulationpass     bool
	Sitesfound         []enzymes.Restrictionsites
	Status             string
	Warnings           string
}

type TypeIISAssembly_designJSONBlock struct {
	ID                     *execute.ThreadID
	BlockID                *execute.BlockID
	Error                  *bool
	AssemblyStandard       *string
	Constructname          *string
	Level                  *string
	PartMoClotypesinorder  *[]string
	Partsinorder           *[]string
	RestrictionsitetoAvoid *[]string
	Vector                 *string
	BackupParts            *[]string
	NewDNASequence         *wtype.DNASequence
	PartswithOverhangs     *[]wtype.DNASequence
	Simulationpass         *bool
	Sitesfound             *[]enzymes.Restrictionsites
	Status                 *string
	Warnings               *string
}

func (c *TypeIISAssembly_design) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AssemblyStandard", "string", "AssemblyStandard", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Constructname", "string", "Constructname", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Level", "string", "Level", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("PartMoClotypesinorder", "[]string", "PartMoClotypesinorder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Partsinorder", "[]string", "Partsinorder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("RestrictionsitetoAvoid", "[]string", "RestrictionsitetoAvoid", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vector", "string", "Vector", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("BackupParts", "[]string", "BackupParts", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("NewDNASequence", "wtype.DNASequence", "NewDNASequence", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("PartswithOverhangs", "[]wtype.DNASequence", "PartswithOverhangs", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Simulationpass", "bool", "Simulationpass", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Sitesfound", "[]enzymes.Restrictionsites", "Sitesfound", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Warnings", "string", "Warnings", true, true, nil, nil))

	ci := execute.NewComponentInfo("TypeIISAssembly_design", "TypeIISAssembly_design", "", false, inp, outp)

	return ci
}
