//Some examples functions
// Calculate rate of reaction, V, of enzyme displaying Micahelis-Menten kinetics with Vmax, Km and [S] declared
// Calculating [S] and V from g/l concentration and looking up molecular weight of named substrate
// Calculating [S] and V from g/l concentration of DNA of known sequence
// Calculating [S] and V from g/l concentration of Protein product of DNA of known sequence

package Datacrunch

import (
	"fmt"
	//"math"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	"encoding/json"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol

//Amount
// i.e. Moles, M

//Amount

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func (e *Datacrunch) requirements() {
	_ = wunit.Make_units

}

// Actions to perform before protocol itself
func (e *Datacrunch) setup(p DatacrunchParamBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// Core process of the protocol: steps to be performed for each input
func (e *Datacrunch) steps(p DatacrunchParamBlock, r *DatacrunchResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper

	// Work out rate of reaction, V of enzyme with Michaelis-Menten kinetics and [S], Km and Vmax declared
	//Using declared values for S and unit of S
	km := wunit.NewAmount(p.Km, p.Kmunit) //.SIValue()
	s := wunit.NewAmount(p.S, p.Sunit)    //.SIValue()

	r.V = ((s.SIValue() * p.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate based on conc and looking up molecular weight in pubchem

	// Look up properties
	substrate_mw := pubchem.MakeMolecule(p.Substrate_name)

	// calculate moles
	submoles := sequences.Moles(p.SubstrateConc, substrate_mw.MolecularWeight, p.SubstrateVol)
	// calculate molar concentration
	submolarconc := sequences.GtoMolarConc(p.SubstrateConc, substrate_mw.MolecularWeight)

	// make a new amount
	s = wunit.NewAmount(submolarconc, "M")

	// use michaelis menton equation
	v_substrate_name := ((s.SIValue() * p.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate from DNA Sequence
	// calculate molar concentration
	dna_mw := sequences.MassDNA(p.DNA_seq, false, false)
	dnamolarconc := sequences.GtoMolarConc(p.DNAConc, dna_mw)

	// make a new amount
	s = wunit.NewAmount(dnamolarconc, "M")

	// use michaelis menton equation
	v_dna := ((s.SIValue() * p.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate from Protein product of dna Sequence

	// translate
	orf, orftrue := sequences.FindORF(p.DNA_seq)
	var protein_mw float64
	if orftrue == true {
		protein_mw_kDA := sequences.Molecularweight(orf)
		protein_mw = protein_mw_kDA * 1000
		r.Orftrue = orftrue
	}

	// calculate molar concentration
	proteinmolarconc := sequences.GtoMolarConc(p.ProteinConc, protein_mw)

	// make a new amount
	s = wunit.NewAmount(submolarconc, "M")

	// use michaelis menton equation
	v_protein := ((s.SIValue() * p.Vmax) / (s.SIValue() + km.SIValue()))

	// print report
	r.Status = fmt.Sprintln(
		"Rate, V of enzyme at substrate conc", p.S, p.Sunit,
		"of enzyme with Km", km.ToString(),
		"and Vmax", p.Vmax, p.Vmaxunit,
		"=", r.V, p.Vunit, ".",
		"Substrate =", p.Substrate_name, ". We have", p.SubstrateVol.ToString(), "of", p.Substrate_name, "at concentration of", p.SubstrateConc.ToString(),
		"Therefore... Moles of", p.Substrate_name, "=", submoles, "Moles.",
		"Molar Concentration of", p.Substrate_name, "=", submolarconc, "Mol/L.",
		"Rate, V = ", v_substrate_name, p.Vmaxunit,
		"Substrate =", "DNA Sequence of", p.Gene_name, "We have", "concentration of", p.DNAConc.ToString(),
		"Therefore... Molar conc", "=", dnamolarconc, "Mol/L",
		"Rate, V = ", v_dna, p.Vmaxunit,
		"Substrate =", "protein from DNA sequence", p.Gene_name, ".",
		"We have", "concentration of", p.ProteinConc.ToString(),
		"Therefore... Molar conc", "=", proteinmolarconc, "Mol/L",
		"Rate, V = ", v_protein, p.Vmaxunit)
	_ = _wrapper.WaitToEnd()

}

// Actions to perform after steps block to analyze data
func (e *Datacrunch) analysis(p DatacrunchParamBlock, r *DatacrunchResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

func (e *Datacrunch) validation(p DatacrunchParamBlock, r *DatacrunchResultBlock) {
	_wrapper := execution.NewWrapper(p.ID,
		p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *Datacrunch) Complete(params interface{}) {
	p := params.(DatacrunchParamBlock)
	if p.Error {
		e.Orftrue <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Status <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.V <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(DatacrunchResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Orftrue <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Status <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.V <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Orftrue <- execute.ThreadParam{Value: r.Orftrue, ID: p.ID, Error: false}

	e.Status <- execute.ThreadParam{Value: r.Status, ID: p.ID, Error: false}

	e.V <- execute.ThreadParam{Value: r.V, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *Datacrunch) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *Datacrunch) NewConfig() interface{} {
	return &DatacrunchConfig{}
}

func (e *Datacrunch) NewParamBlock() interface{} {
	return &DatacrunchParamBlock{}
}

func NewDatacrunch() interface{} { //*Datacrunch {
	e := new(Datacrunch)
	e.init()
	return e
}

// Mapper function
func (e *Datacrunch) Map(m map[string]interface{}) interface{} {
	var res DatacrunchParamBlock
	res.Error = false || m["DNAConc"].(execute.ThreadParam).Error || m["DNA_seq"].(execute.ThreadParam).Error || m["Gene_name"].(execute.ThreadParam).Error || m["Km"].(execute.ThreadParam).Error || m["Kmunit"].(execute.ThreadParam).Error || m["ProteinConc"].(execute.ThreadParam).Error || m["S"].(execute.ThreadParam).Error || m["SubstrateConc"].(execute.ThreadParam).Error || m["SubstrateVol"].(execute.ThreadParam).Error || m["Substrate_name"].(execute.ThreadParam).Error || m["Sunit"].(execute.ThreadParam).Error || m["Vmax"].(execute.ThreadParam).Error || m["Vmaxunit"].(execute.ThreadParam).Error || m["Vunit"].(execute.ThreadParam).Error

	vDNAConc, is := m["DNAConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vDNAConc.JSONString), &temp)
		res.DNAConc = *temp.DNAConc
	} else {
		res.DNAConc = m["DNAConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vDNA_seq, is := m["DNA_seq"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vDNA_seq.JSONString), &temp)
		res.DNA_seq = *temp.DNA_seq
	} else {
		res.DNA_seq = m["DNA_seq"].(execute.ThreadParam).Value.(string)
	}

	vGene_name, is := m["Gene_name"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vGene_name.JSONString), &temp)
		res.Gene_name = *temp.Gene_name
	} else {
		res.Gene_name = m["Gene_name"].(execute.ThreadParam).Value.(string)
	}

	vKm, is := m["Km"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vKm.JSONString), &temp)
		res.Km = *temp.Km
	} else {
		res.Km = m["Km"].(execute.ThreadParam).Value.(float64)
	}

	vKmunit, is := m["Kmunit"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vKmunit.JSONString), &temp)
		res.Kmunit = *temp.Kmunit
	} else {
		res.Kmunit = m["Kmunit"].(execute.ThreadParam).Value.(string)
	}

	vProteinConc, is := m["ProteinConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vProteinConc.JSONString), &temp)
		res.ProteinConc = *temp.ProteinConc
	} else {
		res.ProteinConc = m["ProteinConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vS, is := m["S"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vS.JSONString), &temp)
		res.S = *temp.S
	} else {
		res.S = m["S"].(execute.ThreadParam).Value.(float64)
	}

	vSubstrateConc, is := m["SubstrateConc"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vSubstrateConc.JSONString), &temp)
		res.SubstrateConc = *temp.SubstrateConc
	} else {
		res.SubstrateConc = m["SubstrateConc"].(execute.ThreadParam).Value.(wunit.Concentration)
	}

	vSubstrateVol, is := m["SubstrateVol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vSubstrateVol.JSONString), &temp)
		res.SubstrateVol = *temp.SubstrateVol
	} else {
		res.SubstrateVol = m["SubstrateVol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vSubstrate_name, is := m["Substrate_name"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vSubstrate_name.JSONString), &temp)
		res.Substrate_name = *temp.Substrate_name
	} else {
		res.Substrate_name = m["Substrate_name"].(execute.ThreadParam).Value.(string)
	}

	vSunit, is := m["Sunit"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vSunit.JSONString), &temp)
		res.Sunit = *temp.Sunit
	} else {
		res.Sunit = m["Sunit"].(execute.ThreadParam).Value.(string)
	}

	vVmax, is := m["Vmax"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vVmax.JSONString), &temp)
		res.Vmax = *temp.Vmax
	} else {
		res.Vmax = m["Vmax"].(execute.ThreadParam).Value.(float64)
	}

	vVmaxunit, is := m["Vmaxunit"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vVmaxunit.JSONString), &temp)
		res.Vmaxunit = *temp.Vmaxunit
	} else {
		res.Vmaxunit = m["Vmaxunit"].(execute.ThreadParam).Value.(string)
	}

	vVunit, is := m["Vunit"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DatacrunchJSONBlock
		json.Unmarshal([]byte(vVunit.JSONString), &temp)
		res.Vunit = *temp.Vunit
	} else {
		res.Vunit = m["Vunit"].(execute.ThreadParam).Value.(string)
	}

	res.ID = m["DNAConc"].(execute.ThreadParam).ID
	res.BlockID = m["DNAConc"].(execute.ThreadParam).BlockID

	return res
}

func (e *Datacrunch) OnDNAConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnDNA_seq(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
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
func (e *Datacrunch) OnGene_name(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
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
func (e *Datacrunch) OnKm(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Km", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnKmunit(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Kmunit", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnProteinConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ProteinConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnS(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("S", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnSubstrateConc(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SubstrateConc", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnSubstrateVol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SubstrateVol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnSubstrate_name(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Substrate_name", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnSunit(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Sunit", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnVmax(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vmax", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnVmaxunit(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vmaxunit", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *Datacrunch) OnVunit(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(14, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Vunit", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type Datacrunch struct {
	flow.Component // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag
	DNAConc        <-chan execute.ThreadParam
	DNA_seq        <-chan execute.ThreadParam
	Gene_name      <-chan execute.ThreadParam
	Km             <-chan execute.ThreadParam
	Kmunit         <-chan execute.ThreadParam
	ProteinConc    <-chan execute.ThreadParam
	S              <-chan execute.ThreadParam
	SubstrateConc  <-chan execute.ThreadParam
	SubstrateVol   <-chan execute.ThreadParam
	Substrate_name <-chan execute.ThreadParam
	Sunit          <-chan execute.ThreadParam
	Vmax           <-chan execute.ThreadParam
	Vmaxunit       <-chan execute.ThreadParam
	Vunit          <-chan execute.ThreadParam
	Orftrue        chan<- execute.ThreadParam
	Status         chan<- execute.ThreadParam
	V              chan<- execute.ThreadParam
}

type DatacrunchParamBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	DNAConc        wunit.Concentration
	DNA_seq        string
	Gene_name      string
	Km             float64
	Kmunit         string
	ProteinConc    wunit.Concentration
	S              float64
	SubstrateConc  wunit.Concentration
	SubstrateVol   wunit.Volume
	Substrate_name string
	Sunit          string
	Vmax           float64
	Vmaxunit       string
	Vunit          string
}

type DatacrunchConfig struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	DNAConc        wunit.Concentration
	DNA_seq        string
	Gene_name      string
	Km             float64
	Kmunit         string
	ProteinConc    wunit.Concentration
	S              float64
	SubstrateConc  wunit.Concentration
	SubstrateVol   wunit.Volume
	Substrate_name string
	Sunit          string
	Vmax           float64
	Vmaxunit       string
	Vunit          string
}

type DatacrunchResultBlock struct {
	ID      execute.ThreadID
	BlockID execute.BlockID
	Error   bool
	Orftrue bool
	Status  string
	V       float64
}

type DatacrunchJSONBlock struct {
	ID             *execute.ThreadID
	BlockID        *execute.BlockID
	Error          *bool
	DNAConc        *wunit.Concentration
	DNA_seq        *string
	Gene_name      *string
	Km             *float64
	Kmunit         *string
	ProteinConc    *wunit.Concentration
	S              *float64
	SubstrateConc  *wunit.Concentration
	SubstrateVol   *wunit.Volume
	Substrate_name *string
	Sunit          *string
	Vmax           *float64
	Vmaxunit       *string
	Vunit          *string
	Orftrue        *bool
	Status         *string
	V              *float64
}

func (c *Datacrunch) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("DNAConc", "wunit.Concentration", "DNAConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNA_seq", "string", "DNA_seq", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Gene_name", "string", "Gene_name", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Km", "float64", "Km", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Kmunit", "string", "Kmunit", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ProteinConc", "wunit.Concentration", "ProteinConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("S", "float64", "S", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SubstrateConc", "wunit.Concentration", "SubstrateConc", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SubstrateVol", "wunit.Volume", "SubstrateVol", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Substrate_name", "string", "Substrate_name", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sunit", "string", "Sunit", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vmax", "float64", "Vmax", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vmaxunit", "string", "Vmaxunit", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Vunit", "string", "Vunit", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Orftrue", "bool", "Orftrue", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Status", "string", "Status", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("V", "float64", "V", true, true, nil, nil))

	ci := execute.NewComponentInfo("Datacrunch", "Datacrunch", "", false, inp, outp)

	return ci
}
