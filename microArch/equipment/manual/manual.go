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
// 1 Royal College St, London NW1 0NH UK

package manual

import (
	"errors"
	"fmt"

	"github.com/antha-lang/antha/internal/github.com/nu7hatch/gouuid"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/action"
	"github.com/antha-lang/antha/microArch/equipment/manual/cli"
	"github.com/antha-lang/antha/microArch/logger"
)

type AnthaManual struct {
	ID         string
	Behaviours []equipment.Behaviour
	Cui        cli.CUI
}

func NewAnthaManual(id string) *AnthaManual {
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

	eq := new(AnthaManual)
	eq.ID = id
	eq.Behaviours = be

	//Let's init the cli part of the manual driver.
	eq.Cui = *cli.NewCUI()
	return eq
}

func (e AnthaManual) GetID() string {
	return e.ID
}

func (e AnthaManual) GetEquipmentDefinition() {
	return
}
func (e AnthaManual) Do(actionDescription equipment.ActionDescription) error {
	//	fmt.Println("BEEN ASKED TO DO ", actionDescription.Action, " --> ", actionDescription.ActionData)

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	levels := make([]cli.MultiLevelMessage, 0)
	levels = append(levels, *cli.NewMultiLevelMessage(fmt.Sprintf("%s", actionDescription.ActionData), nil))
	req := cli.NewCUICommandRequest(id.String(), *cli.NewMultiLevelMessage(
		fmt.Sprintf("%s", actionDescription.Action),
		levels,
	))

	e.Cui.CmdIn <- *req
	res := <-e.Cui.CmdOut
	if res.Error != nil {
		(*logger.GetLogger()).Log(*logger.NewLogEntry(e.GetID(), logger.ERROR, res.Error.Error(), ""))
		return errors.New(fmt.Sprintf("Manual Driver fail: id[%s]: %s", res.Id, res.Error))
	}
	(*logger.GetLogger()).Log(*logger.NewLogEntry(e.GetID(), logger.INFO, fmt.Sprintf("OK: %s.", actionDescription.ActionData), ""))

	return nil
}

func (e AnthaManual) Can(b equipment.ActionDescription) bool {
	for _, eb := range e.Behaviours {
		if eb.Matches(b) {
			return true
		}
	}
	return false
}

func (e AnthaManual) Status() string {
	return "OK"
}

func (e AnthaManual) Shutdown() error {
	e.Cui.Close()
	return nil
}

func (e AnthaManual) Init() error {
	e.Cui.Init()
	e.Cui.RunCLI()
	return nil
}
