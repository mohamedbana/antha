// liquidhandling/liquidhandler.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	//"github.com/antha-lang/antha/anthalib/wutil"
	"github.com/antha-lang/antha/anthalib/execution"
)

// the liquid handler structure defines the interface to a particular liquid handling
// platform. The structure holds the following items:
// - an LHRequest structure defining the characteristics of the platform
// - a channel for communicating with the liquid handler
// additionally three functions are defined to implement platform-specific
// implementation requirements
// in each case the LHRequest structure passed in has some additional information
// added and is then passed out. Features which are already defined (e.g. by the
// scheduler or the user) are respected as constraints and will be left unchanged.
// The three functions define
// - setup (SetupAgent): How sources are assigned to plates and plates to positions
// - layout (LayoutAgent): how experiments are assigned to outputs
// - execution (ExecutionPlanner): generates instructions to implement the required plan
//
// The general mechanism by which requests which refer to specific items as opposed to
// those which only state that an item of a particular kind is required is by the definition
// of an 'inst' tag in the request structure with a guid. If this is defined and valid
// it indicates that this item in the request (e.g. a plate, stock etc.) is a specific
// instance. If this is absent then the GUID will either be created or requested
//
type liquidhandler struct {
	Properties       *LHProperties
	Handler          chan *LHRequest
	SetupAgent       func(*LHRequest, *LHProperties) *LHRequest
	LayoutAgent      func(*LHRequest, *LHProperties) *LHRequest
	ExecutionPlanner func(*LHRequest, *LHProperties) *LHRequest
}

// initialize the liquid handling structure
func Init(properties *LHProperties) *liquidhandler {
	// explicitly enforce a few things

	//minvol:=properties.Cmnvol

	minvol := properties.CurrConf.Minvol

	if minvol == 0.0 {
		panic("liquidhandler Initialization error: Must provide description of capabilities")
	}

	lh := liquidhandler{}
	ch := make(chan *LHRequest, 10)
	go RunLiquidHandler(&ch)
	lh.SetupAgent = BasicSetupAgent
	lh.LayoutAgent = BasicLayoutAgent
	lh.ExecutionPlanner = AdvancedExecutionPlanner
	lh.Properties = properties
	lh.Handler = ch
	return &lh
}

// tell the liquid handler to run
func RunLiquidHandler(*chan *LHRequest) {
	// what's on the other end of the channel is an interesting question
}

// temporary function for testing
func raiseError(err string) {
	// TODO remove this
	fmt.Println("liquidhandling raiseError called: remove this to win")
	fmt.Println(err)
	fmt.Println("error done")
	panic("NO")
}

// high-level function which requests planning and execution for an incoming set of
// solutions
func (this *liquidhandler) MakeSolutions(request *LHRequest) *LHRequest {
	// the minimal request which is possible defines what solutions are to be made
	if request.Output_solutions == nil {
		raiseError("No solutions defined")
	}
	this.Plan(request)
	this.Execute(request)
	return request
}

// run the request - this blocks for inputs then runs
func (this *liquidhandler) Execute(request *LHRequest) {
	// wait for inputs to be available
	// well now... this can actually do things now.
	// need to think about the

	output := NewOutputInterface("not used yet")

	instructions := (*request).Instructions

	if len(instructions) == 0 {
		raiseError("Cannot execute request: no instructions")
	}

	for _, ins := range instructions {
		// these should all be transfers at the top level
		tfr := ins.(TransferInstruction)
		arr := SimpleOutput(tfr, *request)
		for _, sins := range arr {
			s := output.Output(sins)
			fmt.Println(s)
		}
	}
}

// This runs the following steps in order:
// - determine required inputs
// - request inputs
// - define robot setup
// - define output layout
// - generate the robot instructions
// - request consumables and other device setups e.g. heater setting
//
// as described above, steps only have an effect if the required inputs are
// not defined beforehand
//
func (this *liquidhandler) Plan(request *LHRequest) {
	// convert requests to volumes and determine required stock concentrations
	solutions, stockconcs := solution_setup(request, this.Properties)
	request.Output_solutions = solutions
	request.Stockconcs = stockconcs

	// looks at components, determines what inputs are required and
	// requests them
	(*request).Input_solutions = this.GetInputs(request)

	// map components to input plates -- this should just be another call to layout
	// this breaks the pattern established above, needs fixing
	inplates, inass := input_plate_setup(request)
	(*request).Input_plates = inplates
	(*request).Input_assignments = inass

	// set up the mapping of the outputs
	// more pattern-breaking, this needs to be tidied up
	request = this.Layout(request)

	// next define input and output plates... again the function will fill in what is missing
	input_plates := this.GetPlates(request.Input_plates, request.Input_major_group_layouts, request.Input_platetype)
	request.Input_plates = input_plates

	output_plates := this.GetPlates(request.Output_plates, request.Output_major_group_layouts, request.Output_platetype)
	request.Output_plates = output_plates

	// now make instructions
	// yet another violation of the pattern, this needs fixing
	request = this.ExecutionPlan(request)

	// next we need to determine the liquid handler setup
	request = this.Setup(request)
}

// request the inputs which are needed to run the plan, unless they have already
// been requested
func (this *liquidhandler) GetInputs(request *LHRequest) map[string][]*LHComponent {
	solutions := (*request).Output_solutions
	inputs := make(map[string][]*LHComponent, 3)

	for _, solution := range solutions {
		// components are either other solutions or come in as inputs

		components := solution.Components
		for _, component := range components {
			component.Destination = solution.ID

			cmps, ok := inputs[component.Name]
			if !ok {
				cmps = make([]*LHComponent, 0, 3)
			}

			cmps = append(cmps, component)
			inputs[component.Name] = cmps
		}
	}

	var requestinputs map[string][]*LHComponent
	requestinputs = request.Input_solutions

	if len(requestinputs) == 0 {
		requestinputs = make(map[string][]*LHComponent, 5)
	} else {
		requestinputs = request.Input_solutions
	}

	// add any new inputs

	for k, v := range inputs {
		if requestinputs[k] == nil {
			requestinputs[k] = v
		}
	}

	requestinputs = this.MakeStockRequest(requestinputs)
	return requestinputs
}

func makeStockRequest(sample *LHComponent) execution.StockRequest {
	stockrequest := make(execution.StockRequest, 3)
	stockrequest["Name"] = sample.Name
	stockrequest["Volume"] = sample.Vol
	stockrequest["Concentration"] = sample.Conc
	stockrequest["SampleID"] = sample.ID
	return stockrequest
}

// make sure the stocks are coming
func (this *liquidhandler) MakeStockRequest(inputs map[string][]*LHComponent) map[string][]*LHComponent {
	for k, inputarr := range inputs {
		fmt.Println("Input", k)
		for i, sample := range inputarr {
			// we make requests only if samples have not already been reserved
			if sample.Inst == "" {
				rslt := execution.GetContext().StockMgr.RequestStock(makeStockRequest(sample))
				sample.Inst = rslt["inst"].(string)
			}
			inputarr[i] = sample
		}
	}

	return inputs
}

func makePlateStockRequest(plate *LHPlate) execution.StockRequest {
	ret := make(execution.StockRequest, 3)
	ret["name"] = plate.PlateName
	ret["type"] = plate.Type
	ret["plateid"] = plate.ID

	return ret
}

// define which labware to use
// and request specific instances
func (this *liquidhandler) GetPlates(plates map[string]*LHPlate, major_layouts map[int][]string, ptype *LHPlate) map[string]*LHPlate {
	if plates == nil {
		plates = make(map[string]*LHPlate, len(major_layouts))

		// assign new plates
		for i := 0; i < len(major_layouts); i++ {
			newplate := new_plate(ptype)
			plates[newplate.ID] = newplate
		}
	}

	// we should know how many plates we need
	for k, plate := range plates {
		if plate.Inst == "" {
			stockrequest := execution.GetContext().StockMgr.RequestStock(makePlateStockRequest(plate))
			plate.Inst = stockrequest["inst"].(string)
		}

		plates[k] = plate
	}

	return plates
}

// generate setup for the robot
func (this *liquidhandler) Setup(request *LHRequest) *LHRequest {
	// assign the plates to positions
	// this needs to be parameterizable
	return this.SetupAgent(request, this.Properties)
}

// generate the output layout
func (this *liquidhandler) Layout(request *LHRequest) *LHRequest {
	// assign the results to destinations
	// again needs to be parameterized

	return this.LayoutAgent(request, this.Properties)
}

// make the instructions for executing this request
func (this *liquidhandler) ExecutionPlan(request *LHRequest) *LHRequest {
	// finally define the instructions which will enact the transfers
	// this is quite involved, we need a strategy to do this

	return this.ExecutionPlanner(request, this.Properties)
}

// find a new tip

func (this *liquidhandler) GetTips(ntips int, volume float64) {

}
