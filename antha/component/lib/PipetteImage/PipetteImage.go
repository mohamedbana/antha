// variant of aliquot.an whereby the low level MixTo command is used to pipette by row

package PipetteImage

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	//"strconv"
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/flow"
	"github.com/antha-lang/antha/microArch/execution"
	"runtime/debug"
	"sync"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//InPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func (e *PipetteImage) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *PipetteImage) setup(p PipetteImageParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *PipetteImage) steps(p PipetteImageParamBlock, r *PipetteImageResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	chosencolourpalette := image.AvailablePalettes[p.Palettename]
	positiontocolourmap, _ := image.ImagetoPlatelayout(p.Imagefilename, p.OutPlate, chosencolourpalette)

	//Pixels = image.PipetteImagetoPlate(OutPlate, positiontocolourmap, AvailableColours, Colourcomponents, VolumePerWell)

	componentmap, err := image.MakestringtoComponentMap(p.AvailableColours, p.Colourcomponents)
	if err != nil {
		panic(err)
	}

	solutions := make([]*wtype.LHSolution, 0)

	counter := 0
	// currently set up to only pipette if yellow (to make visualisation easier in trilution simulator
	for locationkey, colour := range positiontocolourmap {

		component := componentmap[image.Colourcomponentmap[colour]]

		if component.Type == "dna" {
			component.Type = "DoNotMix"
		}
		fmt.Println(image.Colourcomponentmap[colour])

		if p.OnlythisColour != "" {

			if image.Colourcomponentmap[colour] == p.OnlythisColour {
				counter = counter + 1
				fmt.Println("wells", counter)
				pixelSample := mixer.Sample(component, p.VolumePerWell)
				solution := _wrapper.MixTo(p.OutPlate, locationkey, pixelSample)
				solutions = append(solutions, solution)
			}

		} else {
			fmt.Println("component.Type=", component.CName)
			if component.CName != "white" {
				counter = counter + 1
				fmt.Println("wells", counter)
				pixelSample := mixer.Sample(component, p.VolumePerWell)
				solution := _wrapper.MixTo(p.OutPlate, locationkey, pixelSample)
				solutions = append(solutions, solution)
			}
		}
	}

	r.Numberofpixels = len(r.Pixels)
	fmt.Println("Pixels =", r.Numberofpixels)
	r.Pixels = solutions
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *PipetteImage) analysis(p PipetteImageParamBlock, r *PipetteImageResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *PipetteImage) validation(p PipetteImageParamBlock, r *PipetteImageResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *PipetteImage) Complete(params interface{}) {
	p := params.(PipetteImageParamBlock)
	if p.Error {
		e.Numberofpixels <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Pixels <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PipetteImageResultBlock)
	defer func() {
		if res := recover(); res != nil {
			e.Numberofpixels <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			e.Pixels <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
			execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)

	e.Numberofpixels <- execute.ThreadParam{Value: r.Numberofpixels, ID: p.ID, Error: false}

	e.Pixels <- execute.ThreadParam{Value: r.Pixels, ID: p.ID, Error: false}

	e.analysis(p, r)

	e.validation(p, r)

}

// init function, read characterization info from seperate file to validate ranges?
func (e *PipetteImage) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *PipetteImage) NewConfig() interface{} {
	return &PipetteImageConfig{}
}

func (e *PipetteImage) NewParamBlock() interface{} {
	return &PipetteImageParamBlock{}
}

func NewPipetteImage() interface{} { //*PipetteImage {
	e := new(PipetteImage)
	e.init()
	return e
}

// Mapper function
func (e *PipetteImage) Map(m map[string]interface{}) interface{} {
	var res PipetteImageParamBlock
	res.Error = false || m["AvailableColours"].(execute.ThreadParam).Error || m["Colourcomponents"].(execute.ThreadParam).Error || m["Imagefilename"].(execute.ThreadParam).Error || m["OnlythisColour"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["Palettename"].(execute.ThreadParam).Error || m["VolumePerWell"].(execute.ThreadParam).Error

	vAvailableColours, is := m["AvailableColours"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vAvailableColours.JSONString), &temp)
		res.AvailableColours = *temp.AvailableColours
	} else {
		res.AvailableColours = m["AvailableColours"].(execute.ThreadParam).Value.([]string)
	}

	vColourcomponents, is := m["Colourcomponents"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vColourcomponents.JSONString), &temp)
		res.Colourcomponents = *temp.Colourcomponents
	} else {
		res.Colourcomponents = m["Colourcomponents"].(execute.ThreadParam).Value.([]*wtype.LHComponent)
	}

	vImagefilename, is := m["Imagefilename"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vImagefilename.JSONString), &temp)
		res.Imagefilename = *temp.Imagefilename
	} else {
		res.Imagefilename = m["Imagefilename"].(execute.ThreadParam).Value.(string)
	}

	vOnlythisColour, is := m["OnlythisColour"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vOnlythisColour.JSONString), &temp)
		res.OnlythisColour = *temp.OnlythisColour
	} else {
		res.OnlythisColour = m["OnlythisColour"].(execute.ThreadParam).Value.(string)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vPalettename, is := m["Palettename"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vPalettename.JSONString), &temp)
		res.Palettename = *temp.Palettename
	} else {
		res.Palettename = m["Palettename"].(execute.ThreadParam).Value.(string)
	}

	vVolumePerWell, is := m["VolumePerWell"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImageJSONBlock
		json.Unmarshal([]byte(vVolumePerWell.JSONString), &temp)
		res.VolumePerWell = *temp.VolumePerWell
	} else {
		res.VolumePerWell = m["VolumePerWell"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	res.ID = m["AvailableColours"].(execute.ThreadParam).ID
	res.BlockID = m["AvailableColours"].(execute.ThreadParam).BlockID

	return res
}

func (e *PipetteImage) OnAvailableColours(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("AvailableColours", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage) OnColourcomponents(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Colourcomponents", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage) OnImagefilename(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Imagefilename", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage) OnOnlythisColour(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("OnlythisColour", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage) OnOutPlate(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
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
func (e *PipetteImage) OnPalettename(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Palettename", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage) OnVolumePerWell(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VolumePerWell", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type PipetteImage struct {
	flow.Component   // component "superclass" embedded
	lock             sync.Mutex
	startup          sync.Once
	params           map[execute.ThreadID]*execute.AsyncBag
	AvailableColours <-chan execute.ThreadParam
	Colourcomponents <-chan execute.ThreadParam
	Imagefilename    <-chan execute.ThreadParam
	OnlythisColour   <-chan execute.ThreadParam
	OutPlate         <-chan execute.ThreadParam
	Palettename      <-chan execute.ThreadParam
	VolumePerWell    <-chan execute.ThreadParam
	Numberofpixels   chan<- execute.ThreadParam
	Pixels           chan<- execute.ThreadParam
}

type PipetteImageParamBlock struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	AvailableColours []string
	Colourcomponents []*wtype.LHComponent
	Imagefilename    string
	OnlythisColour   string
	OutPlate         *wtype.LHPlate
	Palettename      string
	VolumePerWell    wunit.Volume
}

type PipetteImageConfig struct {
	ID               execute.ThreadID
	BlockID          execute.BlockID
	Error            bool
	AvailableColours []string
	Colourcomponents []wtype.FromFactory
	Imagefilename    string
	OnlythisColour   string
	OutPlate         wtype.FromFactory
	Palettename      string
	VolumePerWell    wunit.Volume
}

type PipetteImageResultBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	Numberofpixels int
	Pixels         []*wtype.LHSolution
}

type PipetteImageJSONBlock struct {
	ID               *execute.ThreadID
	BlockID          *execute.BlockID
	Error            *bool
	AvailableColours *[]string
	Colourcomponents *[]*wtype.LHComponent
	Imagefilename    *string
	OnlythisColour   *string
	OutPlate         **wtype.LHPlate
	Palettename      *string
	VolumePerWell    *wunit.Volume
	Numberofpixels   *int
	Pixels           *[]*wtype.LHSolution
}

func (c *PipetteImage) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("AvailableColours", "[]string", "AvailableColours", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Colourcomponents", "[]*wtype.LHComponent", "Colourcomponents", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Imagefilename", "string", "Imagefilename", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OnlythisColour", "string", "OnlythisColour", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Palettename", "string", "Palettename", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VolumePerWell", "wunit.Volume", "VolumePerWell", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Numberofpixels", "int", "Numberofpixels", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Pixels", "[]*wtype.LHSolution", "Pixels", true, true, nil, nil))

	ci := execute.NewComponentInfo("PipetteImage", "PipetteImage", "", false, inp, outp)

	return ci
}
