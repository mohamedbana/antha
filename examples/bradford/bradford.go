// Example bradford protocol.
// Computes the standard curve from a linear regression
// TODO: implement replicates from parameters
package bradford
	
import "github.com/antha-lang/antha/execute"
import "github.com/antha-lang/goflow"
import "sync"
import "log"
import "bytes"
import "encoding/json"
import "io"


//import "github.com/antha-lang/antha/examples/bradford"

// import the antha PlateReader device, and a third party go library from github
import (
	"PlateReader"
	"github.com/sajari/regression"
)

// Input parameters for this protocol (data)

// Note: 1 replicate means experiment is in duplicate, etc.

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// None

// No special requirements on inputs
func (e *Bradford) requirements() {
	// None
}

// Condititions run on startup
// Including configuring an controls required, and the blocking level needed
// for them (in this case, per plate of samples processed)
func (e *Bradford) setup(p ParamBlock) {
	control.Config(config.per_plate)

	var control_curve [p.ControlCurvePoints + 1]WaterSolution

	for i := 0; i < p.ControlCurvePoints; i++ {
		go func(i) {
			if i == p.ControlCurvePoints {
				control_curve[i] = mix(distilled_water(p.SampleVolume), bradford_reagent(p.BradfordVolume))
			} else {
				control_curve[i] = serial_dilute(control_protein(p.SampleVolume), p.ControlCurvePoints, p.ControlCurveDilutionFactor, i)
			}
			control_absorbance[i] = plate_reader.read(control_curve[i], p.ReadFrequency)
		}()
	}
}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *Bradford) steps(p ParamBlock) {
	var product = mix(p.Sample(p.SampleVolume) + p.BradfordReagent(p.BradfordVolume))
	e.SampleAbsorbance <- execute.ThreadParam{PlateReader.ReadAbsorbance(product, Wavelength), p.ID}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *Bradford) analysis(p ParamBlock, r ResultBlock) {
	// need the control samples to be completed before doing the analysis
	control.WaitForCompletion()
	// Need to compute the linear curve y = m * x + c
	var r regression.Regression
	r.SetObservedName("Absorbance")
	r.SetVarName(0, "Concentration")
	r.AddDataPoint(regression.DataPoint{Observed: p.ControlCurvePoints + 1, Variables: ControlAbsorbance})
	r.AddDataPoint(regression.DataPoint{Observed: p.ControlCurvePoints + 1, Variables: ControlConcentrations})
	r.RunLinearRegression()
	m := r.GetRegCoeff(0)
	c := r.GetRegCoeff(1)
	e.RSquared <- execute.ThreadParam{r.Rsquared, p.ID}

	e.ProteinConc <- execute.ThreadParam{(e.SampleAbsorbance - c) / m, p.ID}
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *Bradford) validation(p ParamBlock, r ResultBlock) {
	if e.SampleAbsorbance > 1 {
		panic("Sample likely needs further dilution")
	}
	if e.RSquared < 0.9 {
		warn("Low r_squared on standard curve")
	}
	if e.RSquared < 0.7 {
		panic("Bad r_squared on standard curve")
	}
	// TODO: add test of replicate variance
}
// AsyncBag functions
func (e *Bradford) Complete(params interface{}) {
	p := params.(ParamBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p)
	
}

// empty function for interface support
func (e *Bradford) anthaElement() {}

// init function, read characterization info from seperate file to validate ranges?
func (e *Bradford) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func New() *Bradford {
	e := new(Bradford)
	e.init()
	return e
}

// Mapper function
func (e *Bradford) Map(m map[string]interface{}) interface{} {
	var res ParamBlock

	res.SampleVolume = m["SampleVolume"].(execute.ThreadParam).Value.(Volume)	

	res.BradfordVolume = m["BradfordVolume"].(execute.ThreadParam).Value.(Volume)	

	res.ReadFrequency = m["ReadFrequency"].(execute.ThreadParam).Value.(Wavelength)	

	res.ControlCurvePoints = m["ControlCurvePoints"].(execute.ThreadParam).Value.(uint32)	

	res.ControlCurveDilutionFactor = m["ControlCurveDilutionFactor"].(execute.ThreadParam).Value.(uint32)	

	res.ReplicateCount = m["ReplicateCount"].(execute.ThreadParam).Value.(uint32)	

	res.Sample = m["Sample"].(execute.ThreadParam).Value.(WaterSolution)	

	res.BradfordReagent = m["BradfordReagent"].(execute.ThreadParam).Value.(WaterSolution)	

	res.ControlProtein = m["ControlProtein"].(execute.ThreadParam).Value.(WaterSolution)	

	res.DistilledWater = m["DistilledWater"].(execute.ThreadParam).Value.(WaterSolution)	

	return res
}


type Bradford struct {
	flow.Component                    // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once	
	params         map[execute.ThreadID]*execute.AsyncBag
	SampleVolume          <-chan execute.ThreadParam
	BradfordVolume          <-chan execute.ThreadParam
	ReadFrequency          <-chan execute.ThreadParam
	ControlCurvePoints          <-chan execute.ThreadParam
	ControlCurveDilutionFactor          <-chan execute.ThreadParam
	ReplicateCount          <-chan execute.ThreadParam
	Sample          <-chan execute.ThreadParam
	BradfordReagent          <-chan execute.ThreadParam
	ControlProtein          <-chan execute.ThreadParam
	DistilledWater          <-chan execute.ThreadParam
	SampleAbsorbance      chan<- execute.ThreadParam
	ProteinConc      chan<- execute.ThreadParam
	RSquared      chan<- execute.ThreadParam
	control_absorbance      chan<- execute.ThreadParam
	control_concentrations      chan<- execute.ThreadParam
}

type ParamBlock struct {
	ID        execute.ThreadID
	SampleVolume Volume
	BradfordVolume Volume
	ReadFrequency Wavelength
	ControlCurvePoints uint32
	ControlCurveDilutionFactor uint32
	ReplicateCount uint32
	Sample WaterSolution
	BradfordReagent WaterSolution
	ControlProtein WaterSolution
	DistilledWater WaterSolution
}

type ResultBlock struct {
	ID        execute.ThreadID
	SampleAbsorbance Absorbance
	ProteinConc Concentration
	RSquared float32
	control_absorbance []Absorbance
	control_concentrations []float64
}

type JSONBlock struct {
	ID        *execute.ThreadID
	SampleAbsorbance *Absorbance
	ProteinConc *Concentration
	RSquared *float32
	control_absorbance *[]Absorbance
	control_concentrations *[]float64
}

