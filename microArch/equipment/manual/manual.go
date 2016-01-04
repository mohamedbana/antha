// /equipment/manual/manual.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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

//package manual contains the implementation for a representation of the actions
// of different pieces of equipment being carried out manually, that is, by a human.
// Different levels of humanization may exist since manual drivers may also be used
// to represent the usage of non antha connected equipment.
package manual

import (
	"errors"
	"fmt"

	"encoding/json"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/internal/github.com/twinj/uuid"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipment/manual/cli"
	"github.com/antha-lang/antha/microArch/equipment/manual/grpc"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/frontend/socket"
	"github.com/antha-lang/antha/microArch/logger"
	schedulerLiquidhandling "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

type AnthaManual interface{}

/////// CUI IMPLEMENTATION

//AnthaManualCUI is a piece of equipment that receives orders through a CUI interface
type AnthaManualCUI struct {
	//ID the
	ID         string
	Behaviours []equipment.Behaviour
	Cui        *cli.CUI
}

//NewAnthaManualCUI creates an Antha Manual driver with the given id that implements the following behaviours:
// action.MESSAGE
// action.LH_SETUP
// action.LH_MOVE
// action.LH_MOVE_EXPLICIT
// action.LH_MOVE_RAW
// action.LH_ASPIRATE
// action.LH_DISPENSE
// action.LH_LOAD_TIPS
// action.LH_UNLOAD_TIPS
// action.LH_SET_PIPPETE_SPEED
// action.LH_SET_DRIVE_SPEED
// action.LH_STOP
// action.LH_SET_POSITION_STATE
// action.LH_RESET_PISTONS
// action.LH_WAIT
// action.LH_MIX
// action.IN_INCUBATE
// action.IN_INCUBATE_SHAKE
// action.MLH_CHANGE_TIPS
func NewAnthaManualCUI(id string) *AnthaManualCUI {
	//This handler should be able to do every possible action
	be := make([]equipment.Behaviour, 0)
	be = append(be, *equipment.NewBehaviour(action.MESSAGE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SETUP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_EXPLICIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_RAW, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_ASPIRATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_DISPENSE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_LOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_UNLOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_PIPPETE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_DRIVE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_STOP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_POSITION_STATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_RESET_PISTONS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_WAIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MIX, ""))
	be = append(be, *equipment.NewBehaviour(action.IN_INCUBATE, ""))
	be = append(be, *equipment.NewBehaviour(action.IN_INCUBATE_SHAKE, ""))
	be = append(be, *equipment.NewBehaviour(action.MLH_CHANGE_TIPS, ""))

	eq := new(AnthaManualCUI)
	eq.ID = id
	eq.Behaviours = be

	//Let's init the cli part of the manual driver.
	eq.Cui = cli.NewCUI()
	return eq
}

//GetID returns the string that identifies a piece of equipment. Ideally uuids v4 should be used.
func (e AnthaManualCUI) GetID() string {
	return e.ID
}

//GetEquipmentDefinition returns a description of the equipment device in terms of
// operations it can handle, restrictions, configuration options ...
func (e AnthaManualCUI) GetEquipmentDefinition() {
	return
}

//Perform an action in the equipment. Actions might be transmitted in blocks to the equipment
func (e AnthaManualCUI) Do(actionDescription equipment.ActionDescription) error {
	id := uuid.NewV4()
	levels := make([]cli.MultiLevelMessage, 0)
	levels = append(levels, *cli.NewMultiLevelMessage(fmt.Sprintf("%s", actionDescription.ActionData), nil))
	req := cli.NewCUICommandRequest(id.String(), *cli.NewMultiLevelMessage(
		fmt.Sprintf("%s", actionDescription.Action),
		levels,
	))

	e.Cui.CmdIn <- *req
	res := <-e.Cui.CmdOut
	if res.Error != nil {
		logger.Error(res.Error.Error())
		return errors.New(fmt.Sprintf("Manual Driver fail: id[%s]: %s", res.Id, res.Error))
	}
	logger.Info(fmt.Sprintf("OK: %s.", actionDescription.ActionData))

	return nil
}

//Can queries a piece of equipment about an action execution. The description of the action must meet the constraints
// of the piece of equipment.
func (e AnthaManualCUI) Can(b equipment.ActionDescription) bool {
	for _, eb := range e.Behaviours {
		if eb.Matches(b) {
			return true
		}
	}
	return false
}

//Status should give a description of the current execution status and any future actions queued to the device
func (e AnthaManualCUI) Status() string {
	return "OK"
}

//Shutdown disconnect, turn off, signal whatever is necessary for a graceful shutdown
func (e AnthaManualCUI) Shutdown() error {
	e.Cui.Close()
	return nil
}

//Init driver will be initialized when registered
func (e AnthaManualCUI) Init() error {
	e.Cui.Init()
	e.Cui.RunCLI()
	return nil
}

/////// SOCKET IMPLEMENTATION

//AnthaManualSocket is a piece of equipment that receives orders through a CUI interface
type AnthaManualSocket struct {
	//ID the
	ID         string
	Behaviours []equipment.Behaviour
	Socket     socket.Socket // TODO put reference to monolith socket library in here
}

//NewAnthaManualSocket creates an Antha Manual driver with the given id that implements the following behaviours:
// action.MESSAGE
func NewAnthaManualSocket(id string) *AnthaManualSocket {
	//This handler should be able to do every possible action
	be := make([]equipment.Behaviour, 0)
	be = append(be, *equipment.NewBehaviour(action.MESSAGE, ""))
	/* TODO Remove this block, only for testing purposes */
	be = append(be, *equipment.NewBehaviour(action.LH_SETUP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_EXPLICIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_RAW, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_ASPIRATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_DISPENSE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_LOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_UNLOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_PIPPETE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_DRIVE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_STOP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_POSITION_STATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_RESET_PISTONS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_WAIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MIX, ""))
	be = append(be, *equipment.NewBehaviour(action.IN_INCUBATE, ""))
	be = append(be, *equipment.NewBehaviour(action.IN_INCUBATE_SHAKE, ""))
	be = append(be, *equipment.NewBehaviour(action.MLH_CHANGE_TIPS, ""))
	/* TODO Remove ^^^^^^^^^^^^^^^^^^^^^ */

	eq := new(AnthaManualSocket)
	eq.ID = id
	eq.Behaviours = be

	//Let's init the cli part of the manual driver.
	return eq
}

//GetID returns the string that identifies a piece of equipment. Ideally uuids v4 should be used.
func (e AnthaManualSocket) GetID() string {
	return e.ID
}

//GetEquipmentDefinition returns a description of the equipment device in terms of
// operations it can handle, restrictions, configuration options ...
func (e AnthaManualSocket) GetEquipmentDefinition() {
	return
}

//Perform an action in the equipment. Actions might be transmitted in blocks to the equipment
func (e AnthaManualSocket) Do(actionDescription equipment.ActionDescription) error {
	logger.Debug(
		fmt.Sprintf("Inst in Manual driver %s --> %s", actionDescription.Action, actionDescription.ActionData),
	)

	//switch the different supported actions and call the specific socket implementation for the message.
	switch actionDescription.Action {
	case action.MESSAGE:
		e.Socket.Message(actionDescription.ActionData)
	}

	return nil
}

//Can queries a piece of equipment about an action execution. The description of the action must meet the constraints
// of the piece of equipment.
func (e AnthaManualSocket) Can(b equipment.ActionDescription) bool {
	for _, eb := range e.Behaviours {
		if eb.Matches(b) {
			return true
		}
	}
	return false
}

//Status should give a description of the current execution status and any future actions queued to the device
func (e AnthaManualSocket) Status() string {
	return "OK"
}

//Shutdown disconnect, turn off, signal whatever is necessary for a graceful shutdown
func (e AnthaManualSocket) Shutdown() error {
	return nil
}

//Init driver will be initialized when registered
func (e AnthaManualSocket) Init() error {
	return nil
}

/////// GRPC IMPLEMENTATION
type AnthaManualGrpc struct {
	ID         string
	Behaviours []equipment.Behaviour
	properties *liquidhandling.LHProperties
	driver     *grpc.Driver
	queue      map[execute.ThreadID]*schedulerLiquidhandling.LHRequest
	queueLock  sync.Mutex
	planner    map[execute.ThreadID]*schedulerLiquidhandling.Liquidhandler
	wg         map[execute.ThreadID]*sync.WaitGroup
}

func NewAnthaManualGrpc(id string, uri string) *AnthaManualGrpc {
	driver := grpc.NewDriver(uri)
	//	driver.Go()

	be := make([]equipment.Behaviour, 0)
	//	be = append(be, *equipment.NewBehaviour(action.MESSAGE, ""))
	/* TODO Remove this block, only for testing purposes */
	be = append(be, *equipment.NewBehaviour(action.LH_SETUP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_EXPLICIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MOVE_RAW, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_ASPIRATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_DISPENSE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_LOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_UNLOAD_TIPS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_PIPPETE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_DRIVE_SPEED, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_STOP, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_SET_POSITION_STATE, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_RESET_PISTONS, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_WAIT, ""))
	be = append(be, *equipment.NewBehaviour(action.LH_MIX, ""))
	be = append(be, *equipment.NewBehaviour(action.IN_INCUBATE, ""))
	be = append(be, *equipment.NewBehaviour(action.IN_INCUBATE_SHAKE, ""))
	be = append(be, *equipment.NewBehaviour(action.MLH_CHANGE_TIPS, ""))
	/* TODO Remove ^^^^^^^^^^^^^^^^^^^^^ */

	ret := &AnthaManualGrpc{
		id,
		be,
		nil,
		driver,
		make(map[execute.ThreadID]*schedulerLiquidhandling.LHRequest),
		sync.Mutex{},
		make(map[execute.ThreadID]*schedulerLiquidhandling.Liquidhandler),
		make(map[execute.ThreadID]*sync.WaitGroup),
	}

	return ret
}

//GetID returns the string that identifies a piece of equipment. Ideally uuids v4 should be used.
func (e *AnthaManualGrpc) GetID() string {
	return e.ID
}

//GetEquipmentDefinition returns a description of the equipment device in terms of
// operations it can handle, restrictions, configuration options ...
func (e *AnthaManualGrpc) GetEquipmentDefinition() {
	return
}

//Perform an action in the equipment. Actions might be transmitted in blocks to the equipment
func (e *AnthaManualGrpc) Do(actionDescription equipment.ActionDescription) error {
	switch actionDescription.Action {
	case action.LH_MIX:
		return e.sendMix(actionDescription)
	case action.LH_END:
		return e.end(actionDescription)
	case action.LH_CONFIG:
		return e.configRequest(actionDescription)
	default:
		return fmt.Errorf("Not implemented")
	}
	return nil
}

func (e *AnthaManualGrpc) configRequest(actionDescription equipment.ActionDescription) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(actionDescription.ActionData), &data); err != nil {
		return err
	}
	pBlockId, err := execute.StringToBlockID(data["BLOCKID"].(string))
	blockId := *pBlockId
	if err != nil {
		return err
	}
	var req *schedulerLiquidhandling.LHRequest

	e.queueLock.Lock()
	defer e.queueLock.Unlock()

	if r, ok := e.queue[blockId.ThreadID]; !ok {
		req = schedulerLiquidhandling.NewLHRequest()
		req.BlockID = blockId
		req.Policies = liquidhandling.GetLHPolicyForTest()
		lhplanner := schedulerLiquidhandling.Init(e.properties)

		e.queue[blockId.ThreadID] = req
		e.planner[blockId.ThreadID] = lhplanner
		wg := sync.WaitGroup{}
		e.wg[blockId.ThreadID] = &wg
		wg.Add(blockId.OutputCount)
	} else {
		req = r
	}
	//@jmanart TODO this down needs to come from the configuration step, need to figure out the correct cast
	// XXX XXX XXX this can cause big issues
	req.Input_Setup_Weights["MAX_N_PLATES"] = 4.5
	req.Input_Setup_Weights["MAX_N_WELLS"] = 278.0
	req.Input_Setup_Weights["RESIDUAL_VOLUME_WEIGHT"] = 1.0

	// MIS MAJOR TODO XXX XXX XXX here: this HARD CODE needs to be removed or we can only use one type of plate!
	// this stuff needs to come from the database - have to work out user configurability here also
	// oh dear, this code is wronger than I had realised
	// MIS fix here - only allow a single plate in here
	if len(req.Input_platetypes) == 0 {
		pwc := factory.GetPlateByType("DWST12")
		req.Input_platetypes = append(req.Input_platetypes, pwc)
	}

	return nil
}
func (e *AnthaManualGrpc) sendMix(actionDescription equipment.ActionDescription) error {
	var sol wtype.LHSolution
	err := json.Unmarshal([]byte(actionDescription.ActionData), &sol)
	if err != nil {
		return err
	}

	e.queueLock.Lock()
	defer e.queueLock.Unlock()

	req, ok := e.queue[sol.BlockID.ThreadID]
	if !ok {
		return fmt.Errorf("Request for block id %v not found", sol.BlockID)
	}
	if sol.Platetype != "" {
		req.Output_platetype = factory.GetPlateByType(sol.Platetype)
	} else {
		req.Output_platetype = factory.GetPlateByType("pcrplate_with_cooler")
	}
	req.Output_solutions[sol.ID] = &sol

	wg := e.wg[sol.BlockID.ThreadID]
	wg.Done()
	return nil
}

func (e *AnthaManualGrpc) end(actionDescription equipment.ActionDescription) error {

	blockId, err := execute.StringToBlockID(actionDescription.ActionData)
	if err != nil {
		return err
	}
	logger.Debug("Waiting for ", blockId)
	wg := e.wg[blockId.ThreadID]
	wg.Wait()

	e.queueLock.Lock()
	defer e.queueLock.Unlock()

	req, ok := e.queue[blockId.ThreadID]
	if !ok || req == nil {
		return nil
	}

	planner, ok := e.planner[blockId.ThreadID]
	if !ok {
		return nil
	}

	planner.MakeSolutions(req)

	e.queue[blockId.ThreadID] = nil
	e.planner[blockId.ThreadID] = nil
	logger.Debug("Request Cleanup Done")

	return nil
}

//Can queries a piece of equipment about an action execution. The description of the action must meet the constraints
// of the piece of equipment.
func (e *AnthaManualGrpc) Can(b equipment.ActionDescription) bool {
	for _, eb := range e.Behaviours {
		if eb.Matches(b) {
			return true
		}
	}
	return false
}

//Status should give a description of the current execution status and any future actions queued to the device
func (e *AnthaManualGrpc) Status() string {
	return "OK"
}

//Shutdown disconnect, turn off, signal whatever is necessary for a graceful shutdown
func (e *AnthaManualGrpc) Shutdown() error {
	return nil
}

//Init driver will be initialized when registered
func (e *AnthaManualGrpc) Init() error {
	//e.properties = factory.GetLiquidhandlerByType("GilsonPipetmax")
	//e.properties = factory.GetLiquidhandlerByType("CyBioGeneTheatre")
	p, _ := e.driver.GetCapabilities()
	e.properties = &p
	e.properties.Driver = e.driver
	return nil
}
