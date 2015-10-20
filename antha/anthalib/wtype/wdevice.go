// wtype/wdevice.go: Part of the Antha language
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

package wtype

import "github.com/antha-lang/antha/antha/anthalib/wunit"

// device interface type
type Device interface {
	//	Solid
	Manufacturer() string
	Type() string
	Ready() bool
}

// Pipetter unit
type Pipetter interface {
	Aspirate(l Liquid)
	Dispense(l Liquid)
	MovePipetteTo(c Coordinates)
	MoveSpeedLimits() wunit.MeasurementLimits
	PipetteSpeedLimits() wunit.MeasurementLimits
	PipetteVolumeLimits() wunit.MeasurementLimits
}

// something which can move Entities about
// should be defined as generally as possible
type Mover interface {
	Grab(e Entity) bool
	Drop(s *Slot) Entity
	MoveTo(c Coordinates)
	MaxWeight() wunit.Mass
	Gripper() VariableSlot
}

// device capable of increasing the temperature
type Heater interface {
	Heat(p Physical, t wunit.Temperature)
	HeatingRate() wunit.Measurement
}

// device capable of decreasing the temperature
type Chiller interface {
	Cool(p Physical, t wunit.Temperature)
	CoolingRate() wunit.Measurement
}

// device capable of sealing labware
type Sealer interface {
	Seal(s Solid) Sealed
}

// device capable of desealing labware
type DeSealer interface {
	Peel(s Sealed) Solid
}

// a holder on a device which can contain labware
type Slot interface {
	SolidContainer
	Dimensions() Geometry
	CanHold(e Entity) bool
}

// a slot which can change size
type VariableSlot interface {
	Slot
	Capabilities() wunit.MeasurementLimits
}
