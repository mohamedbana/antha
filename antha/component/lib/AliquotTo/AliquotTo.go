// variant of aliquot.an whereby the low level MixTo command is used to pipette by row

package AliquotTo

import (
	"encoding/json"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"strconv"
	"sync"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func (e *AliquotTo) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *AliquotTo) setup(p AliquotToParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *AliquotTo) steps(p AliquotToParamBlock, r *AliquotToResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	number := p.SolutionVolume.SIValue() / p.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < p.NumberofAliquots {
		panic("Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHSolution, 0)

	// work out well coordinates for any plate
	wellpositionarray := make([]string, 0)

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//k := 0
	for j := 0; j < p.OutPlate.WlsY; j++ {
		for i := 0; i < p.OutPlate.WlsX; i++ { //countingfrom1iswhatmakesushuman := j + 1
			//k = k + 1
			wellposition := string(alphabet[j]) + strconv.Itoa(i+1)
			//fmt.Println(wellposition, k)
			wellpositionarray = append(wellpositionarray, wellposition)
		}

	}

	for k := 0; k < p.NumberofAliquots; k++ {
		if p.Solution.Type == "dna" {
			p.Solution.Type = "DoNotMix"
		}
		aliquotSample := mixer.Sample(p.Solution, p.VolumePerAliquot)
		aliquot := _wrapper.MixTo(p.OutPlate, wellpositionarray[k], aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	r.Aliquots = aliquots
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *AliquotTo) analysis(p AliquotToParamBlock, r *AliquotToResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *AliquotTo) validation(p AliquotToParamBlock, r *AliquotToResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *AliquotTo) Complete(params interface{}) {
	p := params.(AliquotToParamBlock)
	if p.Error {
		e.Aliquots <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(AliquotToResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Aliquots <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Aliquots <- execute.ThreadParam{Value: r.Aliquots, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *AliquotTo) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *AliquotTo) NewConfig() interface{} {
	return &AliquotToConfig{}
}

func (e *AliquotTo) NewParamBlock() interface{} {
	return &AliquotToParamBlock{}
}

func NewAliquotTo() interface{} { //*AliquotTo {
	e := new(AliquotTo)
	e.init()
	return e
}

// Mapper function
func (e *AliquotTo) Map(m map[string]interface{}) interface{} {
	var res AliquotToParamBlock
	res.Error = false || m["InPlate"].(execute.ThreadParam).Error || m["NumberofAliquots"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Solution"].(execute.ThreadParam).Error || m["SolutionVolume"].(execute.ThreadParam).Error || m["VolumePerAliquot"].(execute.ThreadParam).Error

	vInPlate, is := m["InPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AliquotToJSONBlock
		json.Unmarshal([]byte(vInPlate.JSONString), &temp)
		res.InPlate = *temp.InPlate
	} else {
		res.InPlate = m["InPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vNumberofAliquots, is := m["NumberofAliquots"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AliquotToJSONBlock
		json.Unmarshal([]byte(vNumberofAliquots.JSONString), &temp)
		res.NumberofAliquots = *temp.NumberofAliquots
	} else {
		res.NumberofAliquots = m["NumberofAliquots"].(execute.ThreadParam).Value.(int)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AliquotToJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vSolution, is := m["Solution"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AliquotToJSONBlock
		json.Unmarshal([]byte(vSolution.JSONString), &temp)
		res.Solution = *temp.Solution
	} else {
		res.Solution = m["Solution"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vSolutionVolume, is := m["SolutionVolume"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AliquotToJSONBlock
		json.Unmarshal([]byte(vSolutionVolume.JSONString), &temp)
		res.SolutionVolume = *temp.SolutionVolume
	} else {
		res.SolutionVolume = m["SolutionVolume"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vVolumePerAliquot, is := m["VolumePerAliquot"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp AliquotToJSONBlock
		json.Unmarshal([]byte(vVolumePerAliquot.JSONString), &temp)
		res.VolumePerAliquot = *temp.VolumePerAliquot
	} else {
		res.VolumePerAliquot = m["VolumePerAliquot"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["InPlate"].(execute.ThreadParam).ID
	res.BlockID = m["InPlate"].(execute.ThreadParam).BlockID

	return res
}

func (e *AliquotTo) OnInPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
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
func (e *AliquotTo) OnNumberofAliquots(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("NumberofAliquots", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *AliquotTo) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OutPlate", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *AliquotTo) OnSolution(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Solution", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *AliquotTo) OnSolutionVolume(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("SolutionVolume", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *AliquotTo) OnVolumePerAliquot(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(6, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VolumePerAliquot", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type AliquotTo struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	InPlate          <-chan execute.ThreadParam
	NumberofAliquots <-chan execute.ThreadParam
	OutPlate         <-chan execute.ThreadParam
	Solution         <-chan execute.ThreadParam
	SolutionVolume   <-chan execute.ThreadParam
	VolumePerAliquot <-chan execute.ThreadParam
	Aliquots         chan<- execute.ThreadParam
}

type AliquotToParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	InPlate          *wtype.LHPlate
	NumberofAliquots int
	OutPlate         *wtype.LHPlate
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type AliquotToConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	InPlate          wtype.FromFactory
	NumberofAliquots int
	OutPlate         wtype.FromFactory
	Solution         wtype.FromFactory
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type AliquotToResultBlock struct {
	ID       execute.ThreadID
	BlockID  execute.BlockID
	Error    bool
	Aliquots []*wtype.LHSolution
}

type AliquotToJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	InPlate          **wtype.LHPlate
	NumberofAliquots *int
	OutPlate         **wtype.LHPlate
	Solution         **wtype.LHComponent
	SolutionVolume   *wunit.Volume
	VolumePerAliquot *wunit.Volume
	Aliquots         *[]*wtype.LHSolution
}

func (c *AliquotTo) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("InPlate", "*wtype.LHPlate", "InPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("NumberofAliquots", "int", "NumberofAliquots", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Solution", "*wtype.LHComponent", "Solution", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("SolutionVolume", "wunit.Volume", "SolutionVolume", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VolumePerAliquot", "wunit.Volume", "VolumePerAliquot", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Aliquots", "[]*wtype.LHSolution", "Aliquots", true, true, nil, nil))

	ci := execute.NewComponentInfo("AliquotTo", "AliquotTo", "", false, inp, outp)

	return ci
}
