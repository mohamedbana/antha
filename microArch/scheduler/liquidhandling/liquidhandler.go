// liquidhandling/Liquidhandler.go: Part of the Antha language
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
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
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

// NB for flexibility we should not make the properties object part of this but rather
// send it in as an argument

type Liquidhandler struct {
	Properties       *liquidhandling.LHProperties
	SetupAgent       func(*LHRequest, *liquidhandling.LHProperties) *LHRequest
	LayoutAgent      func(*LHRequest, *liquidhandling.LHProperties) *LHRequest
	ExecutionPlanner func(*LHRequest, *liquidhandling.LHProperties) *LHRequest
	PolicyManager    *LHPolicyManager
	Counter          int
	Once             sync.Once
}

// initialize the liquid handling structure
func Init(properties *liquidhandling.LHProperties) *Liquidhandler {
	lh := Liquidhandler{}
	lh.SetupAgent = BasicSetupAgent
	lh.LayoutAgent = BasicLayoutAgent
	lh.ExecutionPlanner = AdvancedExecutionPlanner2
	lh.Properties = properties
	return &lh
}

// high-level function which requests planning and execution for an incoming set of
// solutions
func (this *Liquidhandler) MakeSolutions(request *LHRequest) *LHRequest {
	// the minimal request which is possible defines what solutions are to be made
	if len(request.Output_solutions) == 0 {
		return request
	}

	f := func() {
		this.Plan(request)
		this.Execute(request)

		// output some info on the final setup

		OutputSetup(this.Properties)
	}

	this.Once.Do(f)

	return request
}

// run the request via the driver
func (this *Liquidhandler) Execute(request *LHRequest) error {
	// set up the robot

	this.do_setup(request)

	instructions := (*request).Instructions

	if instructions == nil {
		RaiseError("Cannot execute request: no instructions")
	}

	// some timing info for the log (only) for now

	timer := this.Properties.GetTimer()
	var d time.Duration

	for _, ins := range instructions {
		//		logger.Debug(fmt.Sprintln(liquidhandling.InsToString(ins)))
		ins.(liquidhandling.TerminalRobotInstruction).OutputTo(this.Properties.Driver)

		if timer != nil {
			d += timer.TimeFor(ins)
		}
	}

	logger.Debug(fmt.Sprintf("Total time estimate: %s", d.String()))

	return nil
}

func (this *Liquidhandler) do_setup(rq *LHRequest) {
	this.Properties.Driver.RemoveAllPlates()
	for position, plateid := range this.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		plate := this.Properties.PlateLookup[plateid]
		name := plate.(wtype.Named).GetName()
		this.Properties.Driver.AddPlateTo(position, plate, name)
	}

	// XXX -- this needs to check this won't be an error
	this.Properties.Driver.(liquidhandling.ExtendedLiquidhandlingDriver).UpdateMetaData(this.Properties)
}

// This runs the following steps in order:
// - determine required inputs
// - request inputs	--- should be moved out
// - define robot setup
// - define output layout
// - generate the robot instructions
// - request consumables and other device setups e.g. heater setting
//
// as described above, steps only have an effect if the required inputs are
// not defined beforehand
//
// so essentially the idea is to parameterise all requests to liquid handlers
// using a Command structure called LHRequest
//
// Depending on its state of completeness, the request structure may be executable
// immediately or may need some additional definition. The purpose of the liquid
// handling service is to provide methods to invoke when parts of the request need
// further definition.
//
// when running a request we should be able to provide mechanisms for pushing requests
// back into the queue to allow them to be cached
//
// this should be OK since the LHRequest parameterises all state including instructions
// for asynchronous drivers we have to determine how far the program got before it was
// paused, which should be tricky but possible.
//
// need to find a good way to codify the rules of the system:
// essentially the question is what happens to inputs pre-defined.
// I will define this asap
//

func (this *Liquidhandler) Plan(request *LHRequest) {
	// convert requests to volumes and determine required stock concentrations
	solutions, stockconcs := solution_setup(request, this.Properties)
	request.Output_solutions = solutions
	request.Stockconcs = stockconcs

	// looks at components, determines what inputs are required and
	// requests them
	request = this.GetInputs(request)

	// define the input plates

	request = input_plate_setup(request)

	// set up the mapping of the outputs
	// this assumes the input plates are set
	request = this.Layout(request)

	// define the output plates
	request = output_plate_setup(request)

	// next we need to determine the liquid handler setup
	request = this.Setup(request)

	// now make instructions
	request = this.ExecutionPlan(request)

	// define the tip boxes - this will depend on the execution plan
	request = this.Tip_box_setup(request)
}

// request the inputs which are needed to run the plan, unless they have already
// been requested
func (this *Liquidhandler) GetInputs(request *LHRequest) *LHRequest {

	if this.Counter > 0 {
		return request
	}
	this.Counter += 1

	solutions := (*request).Output_solutions
	inputs := make(map[string][]*wtype.LHComponent, 3)

	order := make(map[string]map[string]int, 3)

	for _, solution := range solutions {
		// components are either other solutions or come in as inputs
		// this needs solving too
		components := solution.Components

		for _, component := range components {
			component.Destination = solution.ID
			cmps, ok := inputs[component.CName]
			if !ok {
				cmps = make([]*wtype.LHComponent, 0, 3)
			}

			cmps = append(cmps, component)
			inputs[component.CName] = cmps

			for j := 0; j < len(components); j++ {
				if component.Order < components[j].Order {
					m, ok := order[component.CName]
					if !ok {
						m = make(map[string]int, len(components))
						order[component.CName] = m
					}

					m[components[j].CName] += 1
				} else {
					m, ok := order[components[j].CName]
					if !ok {
						m = make(map[string]int, len(components))
						order[components[j].CName] = m
					}
					m[component.CName] += 1
				}
			}

		}
	}

	// define component ordering

	component_order := DefineOrderOrFail(order)
	(*request).Input_order = component_order

	var requestinputs map[string][]*wtype.LHComponent
	requestinputs = request.Input_solutions

	if len(requestinputs) == 0 {
		requestinputs = make(map[string][]*wtype.LHComponent, 5)
	}

	// add any new inputs

	for k, v := range inputs {
		if requestinputs[k] == nil {
			requestinputs[k] = v
		}
	}

	(*request).Input_solutions = requestinputs

	// fix some tips in place
	// TODO this has to be sorted out
	// SERIOUSLY
	max_n_tipboxes := len(this.Properties.Tip_preferences)

	fmt.Println("MAX N TIPBOXES: ", max_n_tipboxes)

	for i := 0; i < max_n_tipboxes; i++ {
		//		this.Properties.AddTipBox(request.Tip_Type.Dup())//TODO get this from where it comes, quick hack now!
		//		this.Properties.AddTipBox(factory.GetTipboxByType("Gilson20"))

		// XXX this needs attention: we shouldn't allow this HARD CODE
		// in future we need to use the validation mechanism to trap this way earlier
		// MARKED FOR DELETION --- THIS NOW IS HANDLED ELSEWHERE
		if request.Tip_Type == nil || request.Tip_Type.GenericSolid == nil {
			logger.Debug(fmt.Sprintf("LiquidHandling model is %q", this.Properties.Model))
			if this.Properties.Model == "Pipetmax" {
				// original
				this.Properties.AddTipBox(factory.GetTipboxByType("Gilson20"))
				this.Properties.Tips = make([]*wtype.LHTip, 1)
				this.Properties.Tips[0] = factory.GetTipboxByType("Gilson20").Tiptype

				// larger vol
				/*this.Properties.AddTipBox(factory.GetTipboxByType("Gilson200"))

				this.Properties.Tips = make([]*wtype.LHTip, 1)
				this.Properties.Tips[0] = factory.GetTipboxByType("Gilson200").Tiptype*/
			} else { //if this.Properties.Model == "GeneTheatre" { //TODO handle general case differently
				this.Properties.AddTipBox(factory.GetTipboxByType("CyBio50Tipbox"))
				this.Properties.Tips = make([]*wtype.LHTip, 1)
				this.Properties.Tips[0] = factory.GetTipboxByType("CyBio50Tipbox").Tiptype
			}
		} else {
			//this.Properties.AddTipBox(factory.GetTipboxByType("Gilson200"))
			this.Properties.AddTipBox(factory.GetTipboxByType(request.Tip_Type.Name()))
		}
	}

	// finally we have to add a waste

	var waste *wtype.LHTipwaste
	// again we don't want this to happen
	// MARKED FOR DELETION... SHOULD BE HANDLED ELSEWHERE
	if this.Properties.Model == "Pipetmax" {
		waste = factory.GetTipwasteByType("Gilsontipwaste")
	} else { //if this.Properties.Model == "GeneTheatre" { //TODO handle general case differently
		waste = factory.GetTipwasteByType("CyBiotipwaste")
	}

	this.Properties.AddTipWaste(waste)

	return request
}

func DefineOrderOrFail(mapin map[string]map[string]int) []string {
	cmps := make([]string, 0, 1)

	for name, _ := range mapin {
		cmps = append(cmps, name)
	}

	ord := make([][]string, len(cmps))

	mx := 0
	for i := 0; i < len(cmps); i++ {
		cnt := 0
		for j := 0; j < len(cmps); j++ {
			if i == j {
				continue
			}

			// only one side can be > 0

			c1 := mapin[cmps[i]][cmps[j]]
			c2 := mapin[cmps[j]][cmps[i]]

			if c1 > 0 && c2 > 0 {
				log.Fatal("CANNOT DEAL WITH INCONSISTENT COMPONENT ORDERING")
			}

			// if c1 > 0 we add to the count

			if c1 > 0 {
				cnt += 1
			}
		}

		a := ord[cnt]

		if a == nil {
			a = make([]string, 0, 3)
		}

		a = append(a, cmps[i])
		if cnt > mx {
			mx = cnt
		}
		ord[cnt] = a
	}

	ret := make([]string, 0, len(cmps))

	// take in reverse order
	if len(cmps) > 0 {
		for j := mx; j >= 0; j-- {
			a := ord[j]
			if a == nil {
				continue
			}

			for _, name := range a {
				ret = append(ret, name)
			}
		}
	}

	return ret
}

// define which labware to use
func (this *Liquidhandler) GetPlates(plates map[string]*wtype.LHPlate, major_layouts map[int][]string, ptype *wtype.LHPlate) map[string]*wtype.LHPlate {
	if plates == nil {
		plates = make(map[string]*wtype.LHPlate, len(major_layouts))

		// assign new plates
		for i := 0; i < len(major_layouts); i++ {
			//newplate := wtype.New_Plate(ptype)
			newplate := factory.GetPlateByType(ptype.Type)
			plates[newplate.ID] = newplate
		}
	}

	// we should know how many plates we need
	for k, plate := range plates {
		if plate.Inst == "" {
			//stockrequest := execution.GetContext().StockMgr.RequestStock(makePlateStockRequest(plate))
			//plate.Inst = stockrequest["inst"].(string)
		}

		plates[k] = plate
	}

	return plates
}

// generate setup for the robot
func (this *Liquidhandler) Setup(request *LHRequest) *LHRequest {
	// assign the plates to positions
	// this needs to be parameterizable
	return this.SetupAgent(request, this.Properties)
}

// generate the output layout
func (this *Liquidhandler) Layout(request *LHRequest) *LHRequest {
	// assign the results to destinations
	// again needs to be parameterized

	return this.LayoutAgent(request, this.Properties)
}

// make the instructions for executing this request
func (this *Liquidhandler) ExecutionPlan(request *LHRequest) *LHRequest {
	// finally define the instructions which will enact the transfers
	// this is quite involved, we need a strategy to do this

	return this.ExecutionPlanner(request, this.Properties)
}

func OutputSetup(robot *liquidhandling.LHProperties) {
	logger.Debug("DECK SETUP INFO")
	logger.Debug("Tipboxes: ")

	for k, v := range robot.Tipboxes {
		logger.Debug(fmt.Sprintf("%s %s: %s", k, robot.PlateIDLookup[k], v.Type))
	}

	logger.Debug("Plates:")

	for k, v := range robot.Plates {
		logger.Debug(fmt.Sprintf("%s %s: %s", k, robot.PlateIDLookup[k], v.PlateName))
	}

	logger.Debug("Tipwastes: ")

	for k, v := range robot.Tipwastes {
		logger.Debug(fmt.Sprintf("%s %s: %s capacity %d", k, robot.PlateIDLookup[k], v.Type, v.Capacity))
	}

}
