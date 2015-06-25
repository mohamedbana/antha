// execution/liquidhandlerservice.go: Part of the Antha language
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

package execution

import (
	"errors"
	"fmt"
	lhdriver "github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/factory"
	"github.com/antha-lang/antha/antha/anthalib/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/execute"
	"sync"
)

// the liquid handler holds channels for communicating
// with the liquid handling service provider
type LiquidHandlingService struct {
	Properties   *lhdriver.LHProperties
	RequestQueue map[execute.ThreadID]*liquidhandling.LHRequest
	lock         *sync.Mutex
}

func NewLiquidHandlingService(p *lhdriver.LHProperties) *LiquidHandlingService {
	var lhs LiquidHandlingService
	lhs.Init(p)
	return &lhs
}

// Initialize the liquid handling service
func (lhs *LiquidHandlingService) Init(p *lhdriver.LHProperties) {
	lhs.Properties = p
	lhs.RequestQueue = make(map[execute.ThreadID]*liquidhandling.LHRequest)
	lhs.lock = new(sync.Mutex)
}

func initRequest(rq *liquidhandling.LHRequest, id execute.ThreadID) *liquidhandling.LHRequest {
	rq.BlockID = string(id)
	rq.Input_Setup_Weights = plateInitWeights(id)
	rq.Policies = initLHPolicies()
	initDriverConfig(rq, id)

	return rq
}

func initDriverConfig(rq *liquidhandling.LHRequest, id execute.ThreadID) {
	ctx := GetContext()
	cfg := ctx.ConfigService.GetConfig(id)
	rq.Robotfn = cfg["SQLITE_FILE_IN"].(string)
	rq.Outputfn = cfg["SQLITE_FILE_OUT"].(string)
}

func initLHPolicies() *lhdriver.LHPolicyRuleSet {
	pol := lhdriver.GetLHPolicyForTest()
	return pol
}

func plateInitWeights(id execute.ThreadID) map[string]float64 {
	ret := make(map[string]float64, 3)
	ctx := GetContext()
	cfg := ctx.ConfigService.GetConfig(id)

	ret["MAX_N_WELLS"] = cfg["MAX_N_WELLS"].(float64)
	ret["MAX_N_PLATES"] = cfg["MAX_N_PLATES"].(float64)
	ret["RESIDUAL_VOLUME_WEIGHT"] = cfg["RESIDUAL_VOLUME_WEIGHT"].(float64)
	return ret
}

func (lhs *LiquidHandlingService) MakeMixRequest(solution *wtype.LHSolution) *liquidhandling.LHRequest {
	lhs.lock.Lock()
	defer lhs.lock.Unlock()
	rq, ok := lhs.RequestQueue[execute.ThreadID(solution.BlockID)]

	if !ok {
		// if we don't have a request with this ID, make a new one
		rq = liquidhandling.NewLHRequest()
		rq = initRequest(rq, execute.ThreadID(solution.BlockID))
	}

	if solution.Platetype != "" {
		rq.Output_platetype = factory.GetPlateByType(solution.Platetype)
	}
	rq.Output_solutions[solution.ID] = solution
	lhs.RequestQueue[execute.ThreadID(rq.BlockID)] = rq

	return rq
}

func (lhs *LiquidHandlingService) ConfigureRequest(id execute.ThreadID, name string, value interface{}) error {

	// get the request out

	rq, ok := lhs.RequestQueue[id]

	if !ok {
		panic(fmt.Sprintf("LiquidHandlingService: No request with id %s", id))
	}

	switch name {
	case "input_platetype":
		plate := value.(*wtype.LHPlate)
		rq.Input_platetypes = append(rq.Input_platetypes, plate)
	case "output_platetype":
		plate := value.(*wtype.LHPlate)
		rq.Output_platetype = plate
	case "tip_type":
		tips := value.(*wtype.LHTipbox)
		rq.Tip_Type = tips
	case "input_setup_weights":
		input_setup_weights := value.(map[string]float64)
		rq.Input_Setup_Weights = input_setup_weights
	default:
		return errors.New(fmt.Sprintf("No such parameter %s", name))

	}

	lhs.RequestQueue[id] = rq

	return nil
}

func (lhs *LiquidHandlingService) SetDriverInputFilename(id execute.ThreadID, fn string) error {
	lhs.lock.Lock()
	defer lhs.lock.Unlock()
	lhr, ok := lhs.RequestQueue[id]

	if !ok {
		return errors.New(fmt.Sprintf("No such ID: %s", id))
	}

	lhr.Robotfn = fn

	lhs.RequestQueue[id] = lhr

	return nil
}

func (lhs *LiquidHandlingService) SetDriverOutputFilename(id execute.ThreadID, fn string) error {
	lhs.lock.Lock()
	defer lhs.lock.Unlock()
	lhr, ok := lhs.RequestQueue[id]

	if !ok {
		return errors.New(fmt.Sprintf("No such ID: %s", id))
	}

	lhr.Outputfn = fn

	lhs.RequestQueue[id] = lhr

	return nil
}

func (lhs *LiquidHandlingService) Run() error {
	lhs.lock.Lock()
	defer lhs.lock.Unlock()

	for _, rq := range lhs.RequestQueue {
		// each block gets executed separately
		liquidhandler := liquidhandling.Init(lhs.Properties)
		liquidhandler.MakeSolutions(rq)
	}

	// clear the queue

	lhs.RequestQueue = make(map[execute.ThreadID]*liquidhandling.LHRequest)
	return nil
}
