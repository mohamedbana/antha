// /antharun/frontend.go: Part of the Antha language
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

package main

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/internal/github.com/twinj/uuid"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/manual"
	"github.com/antha-lang/antha/microArch/equipment/manual/cli"
	"github.com/antha-lang/antha/microArch/equipmentManager"
	"github.com/antha-lang/antha/microArch/logger"
)

type Frontend struct {
	equipmentManager *equipmentManager.AnthaEquipmentManager
	cui              *cli.CUI
}

func NewCUIFrontend() (*Frontend, error) {
	fe := &Frontend{}

	eid := uuid.NewV4()
	// TODO need to shudown equipmentmanager here on error
	fe.equipmentManager = equipmentManager.NewAnthaEquipmentManager(eid.String())
	fee := equipmentManager.EquipmentManager(fe.equipmentManager)
	equipmentManager.SetEquipmentManager(fee)

	mid := uuid.NewV4()
	md := manual.NewAnthaManualCUI(mid.String())
	mdd := equipment.Equipment(md)
	fe.equipmentManager.RegisterEquipment(mdd)

	//cui logger middleware
	logger.RegisterMiddleware(md.Cui)

	return fe, nil
}

func NewRemoteFrontend(driverAddress string) (*Frontend, error) {
	fe := &Frontend{}

	eid := uuid.NewV4()
	// TODO need to shudown equipmentmanager here on error
	fe.equipmentManager = equipmentManager.NewAnthaEquipmentManager(eid.String())
	fee := equipmentManager.EquipmentManager(fe.equipmentManager)
	equipmentManager.SetEquipmentManager(fee)

	remoteManual := manual.NewAnthaManualGrpc(uuid.NewV4().String(), driverAddress)
	if err := remoteManual.Init(); err != nil {
		return nil, err
	}

	err := fe.equipmentManager.RegisterEquipment(remoteManual)
	return fe, err
}

func (fe *Frontend) Shutdown() {
	fe.equipmentManager.Shutdown()
}

func (fe *Frontend) SendAlert(msg interface{}) error {
	if fe.cui != nil {
		var mml cli.MultiLevelMessage
		switch typedMessage := msg.(type) {
		case *wtype.LHSolution:
			mesc := make([]cli.MultiLevelMessage, 0)
			for _, c := range typedMessage.Components {
				mesc = append(mesc, *cli.NewMultiLevelMessage(fmt.Sprintf("%s, %g", c.CName, c.Conc), nil))
			}
			mesC := *cli.NewMultiLevelMessage("Reagents", mesc)
			mesc1 := make([]cli.MultiLevelMessage, 0)
			mesc1 = append(mesc1, mesC)
			mml = *cli.NewMultiLevelMessage(fmt.Sprintf("%s @ %s", typedMessage.SName, typedMessage.Welladdress), mesc1)
		default:
			mml = *cli.NewMultiLevelMessage(fmt.Sprintf("%v", typedMessage), nil)
		}
		mesC := make([]cli.MultiLevelMessage, 0)
		mesC = append(mesC, mml)
		req := cli.NewCUICommandRequest("Alert", *cli.NewMultiLevelMessage(
			"Output",
			mesC,
		))

		fe.cui.CmdIn <- *req
		res := <-fe.cui.CmdOut
		return res.Error
	} else {
		fmt.Println(msg)
	}
	return nil
}