// /equipmentManager/equipmentManager.go: Part of the Antha language
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

package equipmentManager

import (
	"github.com/antha-lang/antha/microArch/equipment"
)

type EquipmentManager interface {
	// GetActionCandidate returns a list of available equipment based on an
	// action we want to perform or nil if no such action can be found
	GetActionCandidate(actionQuery equipment.ActionDescription) equipment.Equipment
}

//AnthaEquipmentManager implements the EquipmentManager interface using an EquipmentREgistry as the means to keep track
// of the different equipment that is registered in the system at a certain time.
type AnthaEquipmentManager struct {
	//ID the uuid that identifies this piece of equipment
	ID string
	//reg the registry used to persist the state
	reg EquipmentRegistry
}

//NewAnthaScheduler returns a new AnthaEquipmentManager identified by id
func NewAnthaEquipmentManager(id string) *AnthaEquipmentManager {
	ret := new(AnthaEquipmentManager)
	ret.ID = id
	ret.reg = NewMemoryEquipmentRegistry()
	return ret
}

var _ EquipmentManager = (*AnthaEquipmentManager)(nil)

//GetID returns the string id representing a manager, usually a uuid
func (e *AnthaEquipmentManager) GetID() string {
	return e.ID
}

// @implements EquipmentManager
func (e *AnthaEquipmentManager) RegisterEquipment(equipment equipment.Equipment) error {
	return e.reg.RegisterEquipment(equipment)
}

func (e *AnthaEquipmentManager) GetEquipment(equipmentQuery string) []equipment.Equipment {
	ret := make([]equipment.Equipment, 0)
	oneq := e.reg.GetEquipmentByID(equipmentQuery)
	if oneq != nil {
		ret = append(ret, oneq)
	}
	return ret
}
func (e *AnthaEquipmentManager) GetActionCandidate(actionQuery equipment.ActionDescription) equipment.Equipment {
	for _, eq := range e.reg.ListEquipment() {
		if eq.Can(actionQuery) {
			return eq //FIFO!!
		}
	}
	return nil
}
func (e *AnthaEquipmentManager) CheckAvailability(equipmentQuery string) {

}
func (e *AnthaEquipmentManager) Shutdown() error {
	for _, eq := range e.reg.ListEquipment() {
		eq.Shutdown()
	}
	return nil
}

//EquipmentRegistry Stores Equipment information related to a particular execution environment
type EquipmentRegistry interface {
	//RegisterEquipment inserts the given piece of equipment into this registry
	RegisterEquipment(eq equipment.Equipment) error
	//GetEquipmentByID gets us a piece of equipment identified by a certain ID
	GetEquipmentByID(id string) equipment.Equipment
	//ListEquipment wil return a list of all the pieces of equipment this registry knows of
	ListEquipment() map[string]equipment.Equipment
}

//EquipmentRegistry in memory implementation of an equipmentRegistry
type MemoryEquipmentRegistry struct {
	EquipmentList map[string]equipment.Equipment
}

//NewMemoryEquipmentRegistry instantiates a new memory registry
func NewMemoryEquipmentRegistry() *MemoryEquipmentRegistry {
	eq := new(MemoryEquipmentRegistry)
	eq.EquipmentList = make(map[string]equipment.Equipment, 0)
	return eq
}
func (reg *MemoryEquipmentRegistry) RegisterEquipment(eq equipment.Equipment) error {
	if _, exists := reg.EquipmentList[eq.GetID()]; !exists {
		reg.EquipmentList[eq.GetID()] = eq
		return eq.Init()
	}
	return nil
}
func (reg *MemoryEquipmentRegistry) GetEquipmentByID(id string) equipment.Equipment {
	if eq, exists := reg.EquipmentList[id]; exists {
		return eq
	}
	return nil
}
func (reg *MemoryEquipmentRegistry) ListEquipment() map[string]equipment.Equipment {
	return reg.EquipmentList
}
