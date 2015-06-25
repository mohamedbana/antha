// equipmentmanagerservice/equipmentmanagerservice.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/anthalib/driver/liquidhandling"
	"github.com/antha-lang/antha/antha/anthalib/factory"
)

// holds channels for communicating with the equipment manager
type EquipmentManagerService struct {
	RequestsIn       chan EquipmentManagerRequest
	RequestsOut      map[string]chan EquipmentManagerRequest
	devicelist       map[string][]string
	deviceproperties map[string]*liquidhandling.LHProperties
	devicequeues     map[string]map[string]interface{}
}

// get properties for a device
func (ems *EquipmentManagerService) GetEquipmentProperties(deviceclass string) interface{} {
	return ems.deviceproperties[ems.devicelist[deviceclass][0]]
}

// bit of a short-term fix

//returns a list of devices known to the system
func (ems *EquipmentManagerService) GetDeviceListByClass(class string) []string {
	return ems.devicelist[class]
}

// get properties files describing the handlers themselves
func (ems *EquipmentManagerService) GetLiquidHandlerProperties(devname string) *liquidhandling.LHProperties {
	return ems.deviceproperties[devname]
}

// ask for some equipment
func (ems *EquipmentManagerService) RequestEquipment(rin EquipmentManagerRequest) chan EquipmentManagerRequest {
	ems.RequestsIn <- rin
	ch := make(chan EquipmentManagerRequest)
	ems.RequestsOut[rin["ID"].(string)] = ch
	return ch
}

// initialize the equipment manager
// needs to read config from somewhere
func (ems *EquipmentManagerService) Init() {
	ems.RequestsIn = make(chan EquipmentManagerRequest, 5)
	ems.RequestsOut = make(map[string]chan EquipmentManagerRequest, 100)

	ems.devicelist = make(map[string][]string)

	// TODO CONFIG -- this needs to take place in an external configuration call

	ems.devicelist["liquidhandler"] = make([]string, 2)
	ems.devicelist["liquidhandler"][0] = "Manual"
	ems.devicelist["liquidhandler"][1] = "GilsonPipetmax"

	ems.deviceproperties = make(map[string]*liquidhandling.LHProperties)
	ems.deviceproperties["GilsonPipetmax"] = factory.GetLiquidhandlerByType("GilsonPipetmax")
	ems.deviceproperties["Manual"] = factory.GetLiquidhandlerByType("Manual")

	ems.devicequeues = make(map[string]map[string]interface{})
	ems.devicequeues["liquidhandler"] = make(map[string]interface{})
	ems.devicequeues["liquidhandler"]["GilsonPipetmax"] = NewLiquidHandlingService(ems.deviceproperties["GilsonPipetmax"])
	ems.devicequeues["liquidhandler"]["Manual"] = NewLiquidHandlingService(ems.deviceproperties["Manual"])

	go func() {
		equipmentmanagerDaemon(ems)
	}()
}

// Daemon for passing requests through to the service
// the new pattern is:
// request comes in, channel comes out. manager stores channel
// when request is serviced the output channel is retrieved and
// fed the output
func equipmentmanagerDaemon(ems *EquipmentManagerService) {
	for {
		rin := <-ems.RequestsIn
		rsp := ems.handleRequest(rin)

		id := rin["ID"].(string)

		rout, ok := ems.RequestsOut[id]

		if !ok {
			panic("No corresponding output channel for request")
		}

		delete(ems.RequestsOut, id)

		rout <- rsp
	}
}

func (ems *EquipmentManagerService) MakeDeviceRequest(devicetype, devicename string) chan EquipmentManagerRequest {
	emr := NewEquipmentManagerRequest()

	emr["requestType"] = "DeviceRequest"
	emr["deviceType"] = devicetype

	if devicename != "" {
		emr["deviceName"] = devicename
	}

	return ems.RequestEquipment(emr)
}

func MakeDeviceResponse() EquipmentManagerRequest {
	emr := NewEquipmentManagerRequest()
	emr["status"] = "FAIL"
	return emr
}

func (ems *EquipmentManagerService) handleRequest(emr EquipmentManagerRequest) EquipmentManagerRequest {
	emrsp := MakeDeviceResponse()
	switch emr["requestType"] {
	case "DeviceRequest":
		// in future this will potentially need to check the queues for the various devices

		switch emr["deviceType"] {
		case "liquidhandler":
			devname, ok := emr["deviceName"].(string)
			if !ok {
				devnamelist := ems.GetDeviceListByClass("liquidhandler")
				devname = devnamelist[0]
			}

			emrsp["devicequeue"] = ems.devicequeues["liquidhandler"][devname]
			emrsp["status"] = "SUCCESS"
		}

	}
	return emrsp
}
