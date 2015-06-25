// /anthalib/incubating/driver.go: Part of the Antha language
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

package Incubating

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/driver"
	"github.com/antha-lang/antha/antha/anthalib/manual"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type ManualIncubating struct {
	ManualDriver *manual.Manual
}

func (m *ManualIncubating) Initialize() driver.CommandStatus {
	m.ManualDriver = manual.NewManual()
	m.ManualDriver.Init() //TODO control panics!
	return driver.CommandStatus{
		OK:        true,
		Errorcode: driver.OK,
		Msg:       "",
	}
}

func (m *ManualIncubating) Finalize() driver.CommandStatus {
	return driver.CommandStatus{
		OK:        true,
		Errorcode: driver.OK,
		Msg:       "",
	}
}

func (m *ManualIncubating) Incubate(matter wtype.Matter, time wunit.Time, temp wunit.Temperature) driver.CommandStatus {
	incubatorDescription := fmt.Sprintf("%s", "Incubator") //TODO need to get the incubator name from somewhere!!!!
	matterName := fmt.Sprintf("%s", matter.MatterType())
	message := fmt.Sprintf("Please, Incubate %s in %s for %d @ %d.", matterName, incubatorDescription, time, temp)
	return m.ManualDriver.Message(message)
}
