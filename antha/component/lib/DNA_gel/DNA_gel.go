// example protocol for loading a DNAgel

package DNA_gel

import (
	//"LiquidHandler"
	//"Labware"
	//"coldplate"
	//"reagents"
	//"Devices"
	//"strconv"
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol (data)

//DNAladder Volume // or should this be a concentration?

//DNAgelruntime time.Duration
//DNAgelwellcapacity Volume
//DNAgelnumberofwells int32
//Organism Taxonomy //= http://www.ncbi.nlm.nih.gov/nuccore/49175990?report=genbank
//Organismgenome Genome
//Target_DNA wtype.DNASequence
//Target_DNAsize float64 //Length
//Runvoltage float64
//AgarosePercentage Percentage
// polyerase kit sets key info such as buffer composition, which effects primer melting temperature for example, along with thermocycle parameters

// Data which is returned from this protocol, and data types

//	NumberofBands[] int
//Bandsizes[] Length
//Bandconc[]Concentration
//Pass bool
//PhotoofDNAgel Image

// Physical Inputs to this protocol with types

//WaterSolution
//WaterSolution //Chemspiderlink // not correct link but similar desirable

//Gel

//DNAladder *wtype.LHComponent//NucleicacidSolution
//Water *wtype.LHComponent//WaterSolution

//DNAgelbuffer *wtype.LHComponent//WaterSolution
//DNAgelNucleicacidintercalator *wtype.LHComponent//ToxicSolution // e.g. ethidium bromide, sybrsafe
//QC_sample *wtype.LHComponent//QC // this is a control
//DNASizeladder *wtype.LHComponent//WaterSolution
//Devices.Gelpowerpack Device
// need to calculate which DNASizeladder is required based on target sequence length and required resolution to distinguish from incorrect assembly possibilities

// Physical outputs from this protocol with types

//Gel
//

// No special requirements on inputs
func (e *DNA_gel) requirements() {
	_ = wunit.Make_units

	// None
	/* QC if negative result should still show band then include QC which will result in band // in reality this may never happen... the primers should be designed within antha too
	   control blank with no template_DNA */
}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func (e *DNA_gel) setup(p DNA_gelParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

	/*control.config.per_DNAgel {
	load DNASizeladder(DNAgelrunvolume) // should run more than one per gel in many cases
	QC := mix (Loadingdye(loadingdyevolume), QC_sample(DNAgelrunvolume-loadingdyevolume))
	load QC(DNAgelrunvolume)
	}*/
}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *DNA_gel) steps(p DNA_gelParamBlock, r *DNA_gelResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	if len(p.Samplenames) != p.Samplenumber {
		panic(fmt.Sprintln("length of sample names:", len(p.Samplenames), "is not equal to sample number:", p.Samplenumber))
	}

	loadedsamples := make([]*wtype.LHSolution, 0)

	var DNAgelloadmix *wtype.LHComponent

	for i := 0; i < p.Samplenumber; i++ {
		// ready to add water to well
		waterSample := mixer.Sample(p.Water, p.Watervol)

		// for troubleshooting
		nothingvol := p.Watervol
		nothingvol.Mvalue = 1.0
		nothingSampletostopitcrashing := mixer.Sample(p.Water, nothingvol)

		// load gel

		if p.Loadingdyeinsample == false {
			DNAgelloadmixsolution := _wrapper.MixInto(
				p.DNAgel,
				mixer.Sample(p.Loadingdye, p.Loadingdyevolume),
				mixer.SampleForTotalVolume(p.Sampletotest, p.DNAgelrunvolume),
			)
			DNAgelloadmix = wtype.SolutionToComponent(DNAgelloadmixsolution)
		} else {
			DNAgelloadmix = p.Sampletotest
		}

		// Ensure  sample will be dispensed appropriately:

		DNAgelloadmix.Type = p.Mixingpolicy
		DNAgelloadmix.CName = p.Samplenames[i] //originalname + strconv.Itoa(i)

		loadedsample := _wrapper.MixInto(
			p.DNAgel,
			waterSample,
			nothingSampletostopitcrashing,
			mixer.Sample(DNAgelloadmix, p.DNAgelrunvolume),
		)

		loadedsamples = append(r.Loadedsamples, loadedsample)
	}
	r.Loadedsamples = loadedsamples
	_ = _wrapper.WaitToEnd()

	// Then run the gel
	/* DNAgel := electrophoresis.Run(Loadedgel,Runvoltage,DNAgelruntime)

		// then analyse
	   	DNAgel.Visualise()
		PCR_product_length = call(assemblydesign_validation).PCR_product_length
		if DNAgel.Numberofbands() == 1
		&& DNAgel.Bandsize(DNAgel[0]) == PCR_product_length {
			Pass = true
			}

		incorrect_assembly_possibilities := assemblydesign_validation.Otherpossibleassemblysizes()

		for _, incorrect := range incorrect_assembly_possibilities {
			if  PCR_product_length == incorrect {
	    pass == false
		S := "matches size of incorrect assembly possibility"
		}

		//cherrypick(positive_colonies,recoverylocation)*/
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *DNA_gel) analysis(p DNA_gelParamBlock, r *DNA_gelResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

	// need the control samples to be completed before doing the analysis

	//

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *DNA_gel) validation(p DNA_gelParamBlock, r *DNA_gelResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

	/* 	if calculatedbandsize == expected {
			stop
		}
		if calculatedbandsize != expected {
		if S == "matches size of incorrect assembly possibility" {
			call(assembly_troubleshoot)
			}
		} // loop at beginning should be designed to split labware resource optimally in the event of any failures e.g. if 96well capacity and 4 failures check 96/4 = 12 colonies of each to maximise chance of getting a hit
	    }
	    if repeat > 2
		stop
	    }
	    if (recoverylocation doesn't grow then use backup or repeat
		}
		if sequencingresults do not match expected then use backup or repeat
	    // TODO: */
}

// AsyncBag functions
func (e *DNA_gel) Complete(params interface{}) {
	p := params.(DNA_gelParamBlock)
	if p.Error {
		e.Loadedsamples <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(DNA_gelResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Loadedsamples <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Loadedsamples <- execute.ThreadParam{Value: r.Loadedsamples, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *DNA_gel) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *DNA_gel) NewConfig() interface{} {
	return &DNA_gelConfig{}
}

func (e *DNA_gel) NewParamBlock() interface{} {
	return &DNA_gelParamBlock{}
}

func NewDNA_gel() interface{} { //*DNA_gel {
	e := new(DNA_gel)
	e.init()
	return e
}

// Mapper function
func (e *DNA_gel) Map(m map[string]interface{}) interface{} {
	var res DNA_gelParamBlock
	res.Error = false || m["DNAgel"].(execute.ThreadParam).Error || m["DNAgelrunvolume"].(execute.ThreadParam).Error || m["InPlate"].(execute.ThreadParam).Error || m["Loadingdye"].(execute.ThreadParam).Error || m["Loadingdyeinsample"].(execute.ThreadParam).Error || m["Loadingdyevolume"].(execute.ThreadParam).Error || m["Mixingpolicy"].(execute.ThreadParam).Error || m["Samplenames"].(execute.ThreadParam).Error || m["Samplenumber"].(execute.ThreadParam).Error || m["Sampletotest"].(execute.ThreadParam).Error || m["Water"].(execute.ThreadParam).Error || m["Watervol"].(execute.ThreadParam).Error

	vDNAgel, is := m["DNAgel"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vDNAgel.JSONString), &temp)
		res.DNAgel = *temp.DNAgel
	} else {
		res.DNAgel = m["DNAgel"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vDNAgelrunvolume, is := m["DNAgelrunvolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vDNAgelrunvolume.JSONString), &temp)
		res.DNAgelrunvolume = *temp.DNAgelrunvolume
	} else {
		res.DNAgelrunvolume = m["DNAgelrunvolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vLoadingdye, is := m["Loadingdye"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vLoadingdye.JSONString), &temp)
		res.Loadingdye = *temp.Loadingdye
	} else {
		res.Loadingdye = m["Loadingdye"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vLoadingdyeinsample, is := m["Loadingdyeinsample"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vLoadingdyeinsample.JSONString), &temp)
		res.Loadingdyeinsample = *temp.Loadingdyeinsample
	} else {
		res.Loadingdyeinsample = m["Loadingdyeinsample"].(execute.ThreadParam).Value.(bool)
	}

	vLoadingdyevolume, is := m["Loadingdyevolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vLoadingdyevolume.JSONString), &temp)
		res.Loadingdyevolume = *temp.Loadingdyevolume
	} else {
		res.Loadingdyevolume = m["Loadingdyevolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vMixingpolicy, is := m["Mixingpolicy"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vMixingpolicy.JSONString), &temp)
		res.Mixingpolicy = *temp.Mixingpolicy
	} else {
		res.Mixingpolicy = m["Mixingpolicy"].(execute.ThreadParam).Value.(string)
	}

	vSamplenames, is := m["Samplenames"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vSamplenames.JSONString), &temp)
		res.Samplenames = *temp.Samplenames
	} else {
		res.Samplenames = m["Samplenames"].(execute.ThreadParam).Value.([]string)
	}

	vSamplenumber, is := m["Samplenumber"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vSamplenumber.JSONString), &temp)
		res.Samplenumber = *temp.Samplenumber
	} else {
		res.Samplenumber = m["Samplenumber"].(execute.ThreadParam).Value.(int)
	}

	vSampletotest, is := m["Sampletotest"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vSampletotest.JSONString), &temp)
		res.Sampletotest = *temp.Sampletotest
	} else {
		res.Sampletotest = m["Sampletotest"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vWater, is := m["Water"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vWater.JSONString), &temp)
		res.Water = *temp.Water
	} else {
		res.Water = m["Water"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vWatervol, is := m["Watervol"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vWatervol.JSONString), &temp)
		res.Watervol = *temp.Watervol
	} else {
		res.Watervol = m["Watervol"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["DNAgel"].(execute.ThreadParam).ID
	res.BlockID = m["DNAgel"].(execute.ThreadParam).BlockID

	return res
}

//func cherrypick ()

func (e *DNA_gel) OnDNAgel(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAgel", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnDNAgelrunvolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAgelrunvolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("InPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnLoadingdye(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Loadingdye", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnLoadingdyeinsample(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Loadingdyeinsample", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnLoadingdyevolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Loadingdyevolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnMixingpolicy(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Mixingpolicy", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnSamplenames(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Samplenames", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnSamplenumber(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Samplenumber", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnSampletotest(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Sampletotest", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnWater(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Water", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *DNA_gel) OnWatervol(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(12, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Watervol", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type DNA_gel struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	DNAgel             <-chan execute.ThreadParam
	DNAgelrunvolume    <-chan execute.ThreadParam
	InPlate            <-chan execute.ThreadParam
	Loadingdye         <-chan execute.ThreadParam
	Loadingdyeinsample <-chan execute.ThreadParam
	Loadingdyevolume   <-chan execute.ThreadParam
	Mixingpolicy       <-chan execute.ThreadParam
	Samplenames        <-chan execute.ThreadParam
	Samplenumber       <-chan execute.ThreadParam
	Sampletotest       <-chan execute.ThreadParam
	Water              <-chan execute.ThreadParam
	Watervol           <-chan execute.ThreadParam
	Loadedsamples      chan<- execute.ThreadParam
}

type DNA_gelParamBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	DNAgel             *wtype.LHPlate
	DNAgelrunvolume    wunit.Volume
	InPlate            *wtype.LHPlate
	Loadingdye         *wtype.LHComponent
	Loadingdyeinsample bool
	Loadingdyevolume   wunit.Volume
	Mixingpolicy       string
	Samplenames        []string
	Samplenumber       int
	Sampletotest       *wtype.LHComponent
	Water              *wtype.LHComponent
	Watervol           wunit.Volume
}

type DNA_gelConfig struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	DNAgel             wtype.FromFactory
	DNAgelrunvolume    wunit.Volume
	InPlate            wtype.FromFactory
	Loadingdye         wtype.FromFactory
	Loadingdyeinsample bool
	Loadingdyevolume   wunit.Volume
	Mixingpolicy       string
	Samplenames        []string
	Samplenumber       int
	Sampletotest       wtype.FromFactory
	Water              wtype.FromFactory
	Watervol           wunit.Volume
}

type DNA_gelResultBlock struct {
	ID            execute.ThreadID
	BlockID       execute.BlockID
	Error         bool
	Loadedsamples []*wtype.LHSolution
}

type DNA_gelJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	DNAgel             **wtype.LHPlate
	DNAgelrunvolume    *wunit.Volume
	InPlate            **wtype.LHPlate
	Loadingdye         **wtype.LHComponent
	Loadingdyeinsample *bool
	Loadingdyevolume   *wunit.Volume
	Mixingpolicy       *string
	Samplenames        *[]string
	Samplenumber       *int
	Sampletotest       **wtype.LHComponent
	Water              **wtype.LHComponent
	Watervol           *wunit.Volume
	Loadedsamples      *[]*wtype.LHSolution
}

func (c *DNA_gel) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("DNAgel", "*wtype.LHPlate", "DNAgel", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAgelrunvolume", "wunit.Volume", "DNAgelrunvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Loadingdye", "*wtype.LHComponent", "Loadingdye", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Loadingdyeinsample", "bool", "Loadingdyeinsample", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Loadingdyevolume", "wunit.Volume", "Loadingdyevolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Mixingpolicy", "string", "Mixingpolicy", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Samplenames", "[]string", "Samplenames", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Samplenumber", "int", "Samplenumber", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sampletotest", "*wtype.LHComponent", "Sampletotest", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Water", "*wtype.LHComponent", "Water", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Watervol", "wunit.Volume", "Watervol", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Loadedsamples", "[]*wtype.LHSolution", "Loadedsamples", true, true, nil, nil))

	ci := execute.NewComponentInfo("DNA_gel", "DNA_gel", "", false, inp, outp)

	return ci
}
