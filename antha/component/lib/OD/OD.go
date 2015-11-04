// Example OD measurement protocol.
// Computes the OD and dry cell weight estimate from absorbance reading
// TODO: implement replicates from parameters
package OD

import (
	//"liquid handler"
	"encoding/json"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"sync"
)

//"standard_labware"
// Input parameters for this protocol (data)

//= uL(100)
//= uL(0)
//Total_volume Volume//= ul (sample_volume+diluent_volume)
//Wavelength //= nm(600)
//Diluent_type //= (PBS)
//= (0.25)
//Replicate_count uint32 //= 1 // Note: 1 replicate means experiment is in duplicate, etc.
// calculate path length? - takes place under plate reader since this will only be necessary for plate reader protocols? labware?
// Data which is returned from this protocol, and data types
//= 0.0533
//WellCrosssectionalArea float64// should be calculated from plate and well type automatically

//Absorbance
//Absorbance
//(pathlength corrected)

//R_squared float32
//Control_absorbance [control_curve_points+1]float64//Absorbance
//Control_concentrations [control_curve_points+1]float64

// Physical Inputs to this protocol with types

//Culture

// Physical outputs from this protocol with types

// None

func (e *OD) requirements() {
	_ = wunit.Make_units

	// sufficient sample volume available to sacrifice
}
func (e *OD) setup(p ODParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

	/*control.Config(config.per_plate)
	var control_blank[total_volume]WaterSolution

	blank_absorbance = platereader.Read(ODplate,control_blank, wavelength)*/
}
func (e *OD) steps(p ODParamBlock, r *ODResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	var product *wtype.LHSolution //WaterSolution

	for {
		product = _wrapper.MixInto(p.ODplate, mixer.Sample(p.Sampletotest, p.Sample_volume), mixer.Sample(p.Diluent, p.Diluent_volume))
		/*Is it necessary to include platetype in Read function?
		or is the info on volume, opacity, pathlength etc implied in LHSolution?*/
		r.Sample_absorbance = platereader.ReadAbsorbance(*p.ODplate, *product, p.Wlength)

		if r.Sample_absorbance.Reading < 1 {
			break
		}
		p.Diluent_volume.Mvalue += 1 //diluent_volume = diluent_volume + 1

	}
	_ = _wrapper.WaitToEnd()

} // serial dilution or could write element for finding optimum dilution or search historical data
func (e *OD) analysis(p ODParamBlock, r *ODResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	// Need to substract blank from measurement; normalise to path length of 1cm for OD value; apply conversion factor to estimate dry cell weight

	r.Blankcorrected_absorbance = platereader.Blankcorrect(r.Sample_absorbance, p.Blank_absorbance)
	volumetopathlengthconversionfactor := wunit.NewLength(p.Heightof100ulinm, "m")                        //WellCrosssectionalArea
	r.OD = platereader.PathlengthCorrect(volumetopathlengthconversionfactor, r.Blankcorrected_absorbance) // 0.0533 could be written as function of labware and liquid volume (or measureed height)
	r.Estimateddrycellweight_conc = wunit.NewConcentration(r.OD.Reading*p.ODtoDCWconversionfactor, "g/L")
	_ = _wrapper.WaitToEnd()

}
func (e *OD) validation(p ODParamBlock, r *ODResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

	/*
		if Sample_absorbance > 1 {
		panic("Sample likely needs further dilution")
		}
		if Sample_absorbance < 0.1 {
		warn("Low OD, sample likely needs increased volume")
		}
		}*/
	// TODO: add test of replicate variance
}

// AsyncBag functions
func (e *OD) Complete(params interface{}) {
	p := params.(ODParamBlock)
	if p.Error {
		e.Blankcorrected_absorbance <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Estimateddrycellweight_conc <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.OD <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Sample_absorbance <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(ODResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Blankcorrected_absorbance <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Estimateddrycellweight_conc <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.OD <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Sample_absorbance <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(res)
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.analysis(p, r)

	e.Blankcorrected_absorbance <- execute.ThreadParam{Value: r.Blankcorrected_absorbance, ID: p.ID, Error: false}

	e.Estimateddrycellweight_conc <- execute.ThreadParam{Value: r.Estimateddrycellweight_conc, ID: p.ID, Error: false}

	e.OD <- execute.ThreadParam{Value: r.OD, ID: p.ID, Error: false}

	e.Sample_absorbance <- execute.ThreadParam{Value: r.Sample_absorbance, ID: p.ID, Error: false}

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *OD) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *OD) NewConfig() interface{} {
	return &ODConfig{}
}

func (e *OD) NewParamBlock() interface{} {
	return &ODParamBlock{}
}

func NewOD() interface{} { //*OD {
	e := new(OD)
	e.init()
	return e
}

// Mapper function
func (e *OD) Map(m map[string]interface{}) interface{} {
	var res ODParamBlock
	res.Error = false || m["Blank_absorbance"].(execute.ThreadParam).Error || m["Diluent"].(execute.ThreadParam).Error || m["Diluent_volume"].(execute.ThreadParam).Error || m["Heightof100ulinm"].(execute.ThreadParam).Error || m["ODplate"].(execute.ThreadParam).Error || m["ODtoDCWconversionfactor"].(execute.ThreadParam).Error || m["Sample_volume"].(execute.ThreadParam).Error || m["Sampletotest"].(execute.ThreadParam).Error || m["Wlength"].(execute.ThreadParam).Error

	vBlank_absorbance, is := m["Blank_absorbance"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vBlank_absorbance.JSONString), &temp)
		res.Blank_absorbance = *temp.Blank_absorbance
	} else {
		res.Blank_absorbance = m["Blank_absorbance"].(execute.ThreadParam).Value.(wtype.Absorbance)
	}

	vDiluent, is := m["Diluent"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vDiluent.JSONString), &temp)
		res.Diluent = *temp.Diluent
	} else {
		res.Diluent = m["Diluent"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vDiluent_volume, is := m["Diluent_volume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vDiluent_volume.JSONString), &temp)
		res.Diluent_volume = *temp.Diluent_volume
	} else {
		res.Diluent_volume = m["Diluent_volume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vHeightof100ulinm, is := m["Heightof100ulinm"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vHeightof100ulinm.JSONString), &temp)
		res.Heightof100ulinm = *temp.Heightof100ulinm
	} else {
		res.Heightof100ulinm = m["Heightof100ulinm"].(execute.ThreadParam).Value.(float64)
	}

	vODplate, is := m["ODplate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vODplate.JSONString), &temp)
		res.ODplate = *temp.ODplate
	} else {
		res.ODplate = m["ODplate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vODtoDCWconversionfactor, is := m["ODtoDCWconversionfactor"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vODtoDCWconversionfactor.JSONString), &temp)
		res.ODtoDCWconversionfactor = *temp.ODtoDCWconversionfactor
	} else {
		res.ODtoDCWconversionfactor = m["ODtoDCWconversionfactor"].(execute.ThreadParam).Value.(float64)
	}

	vSample_volume, is := m["Sample_volume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vSample_volume.JSONString), &temp)
		res.Sample_volume = *temp.Sample_volume
	} else {
		res.Sample_volume = m["Sample_volume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vSampletotest, is := m["Sampletotest"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vSampletotest.JSONString), &temp)
		res.Sampletotest = *temp.Sampletotest
	} else {
		res.Sampletotest = m["Sampletotest"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vWlength, is := m["Wlength"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp ODJSONBlock
		json.Unmarshal([]byte(vWlength.JSONString), &temp)
		res.Wlength = *temp.Wlength
	} else {
		res.Wlength = m["Wlength"].(execute.ThreadParam).Value.(float64)
	}

	res.ID = m["Blank_absorbance"].(execute.ThreadParam).ID
	res.BlockID = m["Blank_absorbance"].(execute.ThreadParam).BlockID

	return res
}

func (e *OD) OnBlank_absorbance(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Blank_absorbance", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnDiluent(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Diluent", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnDiluent_volume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Diluent_volume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnHeightof100ulinm(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Heightof100ulinm", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnODplate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ODplate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnODtoDCWconversionfactor(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("ODtoDCWconversionfactor", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnSample_volume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Sample_volume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *OD) OnSampletotest(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
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
func (e *OD) OnWlength(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(9, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Wlength", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type OD struct {
	flow.Component              // component "superclass" embedded
	lock                        sync.Mutex
	startup                     sync.Once
	params                      map[execute.ThreadID]*execute.AsyncBag
	Blank_absorbance            <-chan execute.ThreadParam
	Diluent                     <-chan execute.ThreadParam
	Diluent_volume              <-chan execute.ThreadParam
	Heightof100ulinm            <-chan execute.ThreadParam
	ODplate                     <-chan execute.ThreadParam
	ODtoDCWconversionfactor     <-chan execute.ThreadParam
	Sample_volume               <-chan execute.ThreadParam
	Sampletotest                <-chan execute.ThreadParam
	Wlength                     <-chan execute.ThreadParam
	Blankcorrected_absorbance   chan<- execute.ThreadParam
	Estimateddrycellweight_conc chan<- execute.ThreadParam
	OD                          chan<- execute.ThreadParam
	Sample_absorbance           chan<- execute.ThreadParam
}

type ODParamBlock struct {
	ID                      execute.ThreadID
	BlockID                 execute.BlockID
	Error                   bool
	Blank_absorbance        wtype.Absorbance
	Diluent                 *wtype.LHComponent
	Diluent_volume          wunit.Volume
	Heightof100ulinm        float64
	ODplate                 *wtype.LHPlate
	ODtoDCWconversionfactor float64
	Sample_volume           wunit.Volume
	Sampletotest            *wtype.LHComponent
	Wlength                 float64
}

type ODConfig struct {
	ID                      execute.ThreadID
	BlockID                 execute.BlockID
	Error                   bool
	Blank_absorbance        wtype.Absorbance
	Diluent                 wtype.FromFactory
	Diluent_volume          wunit.Volume
	Heightof100ulinm        float64
	ODplate                 wtype.FromFactory
	ODtoDCWconversionfactor float64
	Sample_volume           wunit.Volume
	Sampletotest            wtype.FromFactory
	Wlength                 float64
}

type ODResultBlock struct {
	ID                          execute.ThreadID
	BlockID                     execute.BlockID
	Error                       bool
	Blankcorrected_absorbance   wtype.Absorbance
	Estimateddrycellweight_conc wunit.Concentration
	OD                          wtype.Absorbance
	Sample_absorbance           wtype.Absorbance
}

type ODJSONBlock struct {
	ID                          *execute.ThreadID
	BlockID                     *execute.BlockID
	Error                       *bool
	Blank_absorbance            *wtype.Absorbance
	Diluent                     **wtype.LHComponent
	Diluent_volume              *wunit.Volume
	Heightof100ulinm            *float64
	ODplate                     **wtype.LHPlate
	ODtoDCWconversionfactor     *float64
	Sample_volume               *wunit.Volume
	Sampletotest                **wtype.LHComponent
	Wlength                     *float64
	Blankcorrected_absorbance   *wtype.Absorbance
	Estimateddrycellweight_conc *wunit.Concentration
	OD                          *wtype.Absorbance
	Sample_absorbance           *wtype.Absorbance
}

func (c *OD) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Blank_absorbance", "wtype.Absorbance", "Blank_absorbance", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Diluent", "*wtype.LHComponent", "Diluent", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Diluent_volume", "wunit.Volume", "Diluent_volume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Heightof100ulinm", "float64", "Heightof100ulinm", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ODplate", "*wtype.LHPlate", "ODplate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("ODtoDCWconversionfactor", "float64", "ODtoDCWconversionfactor", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sample_volume", "wunit.Volume", "Sample_volume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Sampletotest", "*wtype.LHComponent", "Sampletotest", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Wlength", "float64", "Wlength", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Blankcorrected_absorbance", "wtype.Absorbance", "Blankcorrected_absorbance", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Estimateddrycellweight_conc", "wunit.Concentration", "Estimateddrycellweight_conc", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("OD", "wtype.Absorbance", "OD", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Sample_absorbance", "wtype.Absorbance", "Sample_absorbance", true, true, nil, nil))

	ci := execute.NewComponentInfo("OD", "OD", "", false, inp, outp)

	return ci
}
