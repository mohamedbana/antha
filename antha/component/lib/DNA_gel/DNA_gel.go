package DNA_gel

import (
	//"LiquidHandler"
	//"Labware"
	//"coldplate"
	//"reagents"
	//"Devices"
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

// Input parameters for this protocol (data)

// or should this be a concentration?

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

	// load gel
	var DNAgelloadmix *wtype.LHComponent

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

	loadedgel := _wrapper.MixInto(
		p.DNAgel,
		mixer.Sample(DNAgelloadmix, p.DNAgelrunvolume),
	)

	r.Loadedgel = loadedgel
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
		e.Loadedgel <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(DNA_gelResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Loadedgel <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Loadedgel <- execute.ThreadParam{Value: r.Loadedgel, ID: p.ID, Error: false}

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
	res.Error = false || m["DNAgel"].(execute.ThreadParam).Error || m["DNAgelrunvolume"].(execute.ThreadParam).Error || m["DNAladder"].(execute.ThreadParam).Error || m["Loadingdye"].(execute.ThreadParam).Error || m["Loadingdyeinsample"].(execute.ThreadParam).Error || m["Loadingdyevolume"].(execute.ThreadParam).Error || m["Sampletotest"].(execute.ThreadParam).Error

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

	vDNAladder, is := m["DNAladder"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vDNAladder.JSONString), &temp)
		res.DNAladder = *temp.DNAladder
	} else {
		res.DNAladder = m["DNAladder"].(execute.ThreadParam).Value.(wunit.Volume)
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

	vSampletotest, is := m["Sampletotest"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp DNA_gelJSONBlock
		json.Unmarshal([]byte(vSampletotest.JSONString), &temp)
		res.Sampletotest = *temp.Sampletotest
	} else {
		res.Sampletotest = m["Sampletotest"].(execute.ThreadParam).Value.(*wtype.LHComponent)
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
		bag.Init(7, e, e)
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
		bag.Init(7, e, e)
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
func (e *DNA_gel) OnDNAladder(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("DNAladder", param)
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
		bag.Init(7, e, e)
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
		bag.Init(7, e, e)
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
		bag.Init(7, e, e)
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
func (e *DNA_gel) OnSampletotest(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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

type DNA_gel struct {
	flow.Component     // component "superclass" embedded
	lock               sync.Mutex
	startup            sync.Once
	params             map[execute.ThreadID]*execute.AsyncBag
	DNAgel             <-chan execute.ThreadParam
	DNAgelrunvolume    <-chan execute.ThreadParam
	DNAladder          <-chan execute.ThreadParam
	Loadingdye         <-chan execute.ThreadParam
	Loadingdyeinsample <-chan execute.ThreadParam
	Loadingdyevolume   <-chan execute.ThreadParam
	Sampletotest       <-chan execute.ThreadParam
	Loadedgel          chan<- execute.ThreadParam
}

type DNA_gelParamBlock struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	DNAgel             *wtype.LHPlate
	DNAgelrunvolume    wunit.Volume
	DNAladder          wunit.Volume
	Loadingdye         *wtype.LHComponent
	Loadingdyeinsample bool
	Loadingdyevolume   wunit.Volume
	Sampletotest       *wtype.LHComponent
}

type DNA_gelConfig struct {
	ID                 execute.ThreadID
	BlockID            execute.BlockID
	Error              bool
	DNAgel             wtype.FromFactory
	DNAgelrunvolume    wunit.Volume
	DNAladder          wunit.Volume
	Loadingdye         wtype.FromFactory
	Loadingdyeinsample bool
	Loadingdyevolume   wunit.Volume
	Sampletotest       wtype.FromFactory
}

type DNA_gelResultBlock struct {
	ID        execute.ThreadID
	BlockID   execute.BlockID
	Error     bool
	Loadedgel *wtype.LHSolution
}

type DNA_gelJSONBlock struct {
	ID                 *execute.ThreadID
	BlockID            *execute.BlockID
	Error              *bool
	DNAgel             **wtype.LHPlate
	DNAgelrunvolume    *wunit.Volume
	DNAladder          *wunit.Volume
	Loadingdye         **wtype.LHComponent
	Loadingdyeinsample *bool
	Loadingdyevolume   *wunit.Volume
	Sampletotest       **wtype.LHComponent
	Loadedgel          **wtype.LHSolution
}

func (c *DNA_gel) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("DNAgel", "*wtype.LHPlate", "DNAgel", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAgelrunvolume", "wunit.Volume", "DNAgelrunvolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("DNAladder", "wunit.Volume", "DNAladder", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Loadingdye", "*wtype.LHComponent", "Loadingdye", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Loadingdyeinsample", "bool", "Loadingdyeinsample", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Loadingdyevolume", "wunit.Volume", "Loadingdyevolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sampletotest", "*wtype.LHComponent", "Sampletotest", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Loadedgel", "*wtype.LHSolution", "Loadedgel", true, true, nil, nil))

	ci := execute.NewComponentInfo("DNA_gel", "DNA_gel", "", false, inp, outp)

	return ci
}
