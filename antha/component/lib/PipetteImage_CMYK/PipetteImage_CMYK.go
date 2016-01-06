// Generates instructions to pipette out a defined image onto a defined plate using a defined palette of colours
package PipetteImage_CMYK

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

func (e *PipetteImage_CMYK) requirements() {
	_ = wunit.Make_units

}

// Conditions to run on startup
func (e *PipetteImage_CMYK) setup(p PipetteImage_CMYKParamBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// The core process for this protocol, with the steps to be performed
// for every input
func (e *PipetteImage_CMYK) steps(p PipetteImage_CMYKParamBlock, r *PipetteImage_CMYKResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper

	chosencolourpalette := image.AvailablePalettes["WebSafe"]
	positiontocolourmap, _ := image.ImagetoPlatelayout(p.Imagefilename, p.OutPlate, chosencolourpalette)

	//Pixels = image.PipetteImagetoPlate(OutPlate, positiontocolourmap, AvailableColours, Colourcomponents, VolumePerWell)

	/*componentmap, err := image.MakestringtoComponentMap(AvailableColours, Colourcomponents)
	if err != nil {
		panic(err)
	}*/

	solutions := make([]*wtype.LHSolution, 0)

	counter := 0

	//solutions := image.PipetteImagebyBlending(OutPlate, positiontocolourmap,Cyan, Magenta, Yellow,Black, VolumeForFullcolour)

	for locationkey, colour := range positiontocolourmap {
		counter = counter + 1
		components := make([]*wtype.LHComponent, 0)

		cmyk := image.ColourtoCMYK(colour)

		cyanvol := wunit.NewVolume((float64(cmyk.C) * p.VolumeForFullcolour.SIValue()), "l")
		yellowvol := wunit.NewVolume((float64(cmyk.Y) * p.VolumeForFullcolour.SIValue()), "l")
		magentavol := wunit.NewVolume((float64(cmyk.M) * p.VolumeForFullcolour.SIValue()), "l")
		blackvol := wunit.NewVolume((float64(cmyk.K) * p.VolumeForFullcolour.SIValue()), "l")

		cyanSample := mixer.Sample(p.Cyan, cyanvol)
		components = append(components, cyanSample)
		yellowSample := mixer.Sample(p.Yellow, yellowvol)
		components = append(components, yellowSample)
		magentaSample := mixer.Sample(p.Magenta, magentavol)
		components = append(components, magentaSample)
		blackSample := mixer.Sample(p.Black, blackvol)
		components = append(components, blackSample)
		solution := _wrapper.MixTo(p.OutPlate, locationkey, components...)
		solutions = append(solutions, solution)
	}

	/*
		for locationkey, colour := range positiontocolourmap {

			component := componentmap[image.Colourcomponentmap[colour]]

			if component.Type == "dna" {
				component.Type = "DoNotMix"
			}
			fmt.Println(image.Colourcomponentmap[colour])

			if OnlythisColour !="" {

			if image.Colourcomponentmap[colour] == OnlythisColour{
				counter = counter + 1
				fmt.Println("wells",counter)
			pixelSample := mixer.Sample(component, VolumePerWell)
			solution := MixTo(OutPlate, locationkey, pixelSample)
			solutions = append(solutions, solution)
				}

			}else{
				if component.CName !="white"{
				counter = counter + 1
				fmt.Println("wells",counter)
			pixelSample := mixer.Sample(component, VolumePerWell)
			solution := MixTo(OutPlate, locationkey, pixelSample)
			solutions = append(solutions, solution)
			}
			}
		}
	*/
	r.Pixels = solutions
	r.Numberofpixels = len(r.Pixels)
	fmt.Println("Pixels =", r.Numberofpixels)
	_ = _wrapper.WaitToEnd()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func (e *PipetteImage_CMYK) analysis(p PipetteImage_CMYKParamBlock, r *PipetteImage_CMYKResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func (e *PipetteImage_CMYK) validation(p PipetteImage_CMYKParamBlock, r *PipetteImage_CMYKResultBlock) {
	_wrapper := execution.NewWrapper(p.ID, p.BlockID, p)
	_ = _wrapper
	_ = _wrapper.WaitToEnd()

}

// AsyncBag functions
func (e *PipetteImage_CMYK) Complete(params interface{}) {
	p := params.(PipetteImage_CMYKParamBlock)
	if p.Error {
		e.Numberofpixels <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		e.Pixels <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
		return
	}
	r := new(PipetteImage_CMYKResultBlock)
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
func (e *PipetteImage_CMYK) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *PipetteImage_CMYK) NewConfig() interface{} {
	return &PipetteImage_CMYKConfig{}
}

func (e *PipetteImage_CMYK) NewParamBlock() interface{} {
	return &PipetteImage_CMYKParamBlock{}
}

func NewPipetteImage_CMYK() interface{} { //*PipetteImage_CMYK {
	e := new(PipetteImage_CMYK)
	e.init()
	return e
}

// Mapper function
func (e *PipetteImage_CMYK) Map(m map[string]interface{}) interface{} {
	var res PipetteImage_CMYKParamBlock
	res.Error = false || m["Black"].(execute.ThreadParam).Error || m["Cyan"].(execute.ThreadParam).Error || m["Imagefilename"].(execute.ThreadParam).Error || m["Magenta"].(execute.ThreadParam).Error || m["OutPlate"].(execute.ThreadParam).Error || m["VolumeForFullcolour"].(execute.ThreadParam).Error || m["Yellow"].(execute.ThreadParam).Error

	vBlack, is := m["Black"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vBlack.JSONString), &temp)
		res.Black = *temp.Black
	} else {
		res.Black = m["Black"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vCyan, is := m["Cyan"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vCyan.JSONString), &temp)
		res.Cyan = *temp.Cyan
	} else {
		res.Cyan = m["Cyan"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vImagefilename, is := m["Imagefilename"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vImagefilename.JSONString), &temp)
		res.Imagefilename = *temp.Imagefilename
	} else {
		res.Imagefilename = m["Imagefilename"].(execute.ThreadParam).Value.(string)
	}

	vMagenta, is := m["Magenta"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vMagenta.JSONString), &temp)
		res.Magenta = *temp.Magenta
	} else {
		res.Magenta = m["Magenta"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	vOutPlate, is := m["OutPlate"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vOutPlate.JSONString), &temp)
		res.OutPlate = *temp.OutPlate
	} else {
		res.OutPlate = m["OutPlate"].(execute.ThreadParam).Value.(*wtype.LHPlate)
	}

	vVolumeForFullcolour, is := m["VolumeForFullcolour"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vVolumeForFullcolour.JSONString), &temp)
		res.VolumeForFullcolour = *temp.VolumeForFullcolour
	} else {
		res.VolumeForFullcolour = m["VolumeForFullcolour"].(execute.ThreadParam).Value.(wunit.Volume)
	}

	vYellow, is := m["Yellow"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp PipetteImage_CMYKJSONBlock
		json.Unmarshal([]byte(vYellow.JSONString), &temp)
		res.Yellow = *temp.Yellow
	} else {
		res.Yellow = m["Yellow"].(execute.ThreadParam).Value.(*wtype.LHComponent)
	}

	res.ID = m["Black"].(execute.ThreadParam).ID
	res.BlockID = m["Black"].(execute.ThreadParam).BlockID

	return res
}

func (e *PipetteImage_CMYK) OnBlack(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Black", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage_CMYK) OnCyan(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Cyan", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage_CMYK) OnImagefilename(param execute.ThreadParam) {
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
func (e *PipetteImage_CMYK) OnMagenta(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Magenta", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage_CMYK) OnOutPlate(param execute.ThreadParam) {
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
func (e *PipetteImage_CMYK) OnVolumeForFullcolour(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("VolumeForFullcolour", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}
func (e *PipetteImage_CMYK) OnYellow(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init(7, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("Yellow", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}

type PipetteImage_CMYK struct {
	flow.Component      // component "superclass" embedded
	lock                sync.Mutex
	startup             sync.Once
	params              map[execute.ThreadID]*execute.AsyncBag
	Black               <-chan execute.ThreadParam
	Cyan                <-chan execute.ThreadParam
	Imagefilename       <-chan execute.ThreadParam
	Magenta             <-chan execute.ThreadParam
	OutPlate            <-chan execute.ThreadParam
	VolumeForFullcolour <-chan execute.ThreadParam
	Yellow              <-chan execute.ThreadParam
	Numberofpixels      chan<- execute.ThreadParam
	Pixels              chan<- execute.ThreadParam
}

type PipetteImage_CMYKParamBlock struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	Black               *wtype.LHComponent
	Cyan                *wtype.LHComponent
	Imagefilename       string
	Magenta             *wtype.LHComponent
	OutPlate            *wtype.LHPlate
	VolumeForFullcolour wunit.Volume
	Yellow              *wtype.LHComponent
}

type PipetteImage_CMYKConfig struct {
	ID                  execute.ThreadID
	BlockID             execute.BlockID
	Error               bool
	Black               wtype.FromFactory
	Cyan                wtype.FromFactory
	Imagefilename       string
	Magenta             wtype.FromFactory
	OutPlate            wtype.FromFactory
	VolumeForFullcolour wunit.Volume
	Yellow              wtype.FromFactory
}

type PipetteImage_CMYKResultBlock struct {
	ID             execute.ThreadID
	BlockID        execute.BlockID
	Error          bool
	Numberofpixels int
	Pixels         []*wtype.LHSolution
}

type PipetteImage_CMYKJSONBlock struct {
	ID                  *execute.ThreadID
	BlockID             *execute.BlockID
	Error               *bool
	Black               **wtype.LHComponent
	Cyan                **wtype.LHComponent
	Imagefilename       *string
	Magenta             **wtype.LHComponent
	OutPlate            **wtype.LHPlate
	VolumeForFullcolour *wunit.Volume
	Yellow              **wtype.LHComponent
	Numberofpixels      *int
	Pixels              *[]*wtype.LHSolution
}

func (c *PipetteImage_CMYK) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo, 0)
	outp := make([]execute.PortInfo, 0)
	inp = append(inp, *execute.NewPortInfo("Black", "*wtype.LHComponent", "Black", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Cyan", "*wtype.LHComponent", "Cyan", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Imagefilename", "string", "Imagefilename", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Magenta", "*wtype.LHComponent", "Magenta", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("OutPlate", "*wtype.LHPlate", "OutPlate", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("VolumeForFullcolour", "wunit.Volume", "VolumeForFullcolour", true, true, nil, nil))
	inp = append(inp, *execute.NewPortInfo("Yellow", "*wtype.LHComponent", "Yellow", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Numberofpixels", "int", "Numberofpixels", true, true, nil, nil))
	outp = append(outp, *execute.NewPortInfo("Pixels", "[]*wtype.LHSolution", "Pixels", true, true, nil, nil))

	ci := execute.NewComponentInfo("PipetteImage_CMYK", "PipetteImage_CMYK", "", false, inp, outp)

	return ci
}
