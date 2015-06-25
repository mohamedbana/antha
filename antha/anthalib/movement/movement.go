// /anthalib/movement/movement.go: Part of the Antha language
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

package movement

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/driver"
	"github.com/antha-lang/antha/antha/anthalib/manual"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type ManualMovement struct {
	ManualDriver *manual.Manual
}

func (m *ManualMovement) Initialize() driver.CommandStatus {
	m.ManualDriver = manual.NewManual()
	m.ManualDriver.Init() //TODO control panics!
	return driver.CommandStatus{
		OK:        true,
		Errorcode: driver.OK,
		Msg:       "",
	}
}

func (m *ManualMovement) Finalize() driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: driver.OK,
		Msg:       "",
	}
}

func (m *ManualMovement) Move(entity wtype.Entity, final wtype.Location) driver.CommandStatus {
	initialLocationDescription := fmt.Sprintf("%s", entity.Location().Location_Name())
	finalLocationDescription := fmt.Sprintf("%s", final.Location_Name())
	entityName := fmt.Sprintf("%s", entity.Name())
	message := fmt.Sprintf("Please, move %s from %s to %s", entityName, initialLocationDescription, finalLocationDescription)
	return m.ManualDriver.Message(message)
}

func (m *ManualMovement) Stop() driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: driver.OK,
		Msg:       "",
	}
}

func (m *ManualMovement) Wait(time float64) driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: driver.OK,
		Msg:       "",
	}
}
