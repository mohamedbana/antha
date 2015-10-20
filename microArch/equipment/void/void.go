// microArch/equipment/void/void.go: Part of the Antha language
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

package void

import "github.com/antha-lang/antha/microArch/equipment"

//VoidEquipment is able to do everything, but does nothing
type VoidEquipment struct {
	Id string
}

func NewVoidEquipment(id string) *VoidEquipment {
	ret := new(VoidEquipment)
	ret.Id = id
	return ret
}

func (v *VoidEquipment) GetID() string {
	return v.Id
}
func (v *VoidEquipment) Do(actionDescription equipment.ActionDescription) error {
	return nil
}
func (v *VoidEquipment) Can(ac equipment.ActionDescription) bool {
	return true
}
func (v *VoidEquipment) Status() string {
	return ""
}
func (v *VoidEquipment) Init() error {
	return nil
}
func (v *VoidEquipment) Shutdown() error {
	return nil
}
func (v *VoidEquipment) SetInstance(inst interface{}) error {
	return nil
}
func (v *VoidEquipment) SetCommunicationChannel(in chan interface{}, out chan interface{}) error {
	return nil
}
