// wtype/wtype.go: Part of the Antha language
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

package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// base type for defining materials
type Matter interface {
	MatterType() string
	MeltingPoint() wunit.Temperature
	BoilingPoint() wunit.Temperature
	SpecificHeatCapacity() wunit.SpecificHeatCapacity
}

// a sample of matter
type Physical interface {
	// embedded class for dealing with type of material
	Matter
	// identifier of sample
	Name() string
	SetName(string) string
	// mass of sample
	Mass() wunit.Mass
	SetMass(wunit.Mass) wunit.Mass
	// volume occupied by sample
	Volume() wunit.Volume
	SetVolume(wunit.Volume) wunit.Volume
	// temperature of object
	Temperature() wunit.Temperature
	SetTemperature(t wunit.Temperature)
	// ratio of mass to volume
	Density() wunit.Density
}

// The Entity interface declares that this object is an independently movable thing
type Entity interface {
	// Entities must be solid objects
	Solid
	// since it can be moved independently, an Entity must have a location
	Location() Location
	//SetLocation updates the position of this entity
	SetLocation(newLocation Location) error
}

// solid state
type Solid interface {
	Physical
	Shape() Shape
}

// liquid state
type Liquid interface {
	Physical
	Viscosity() float64
	// take some of this liquid
	Sample(v wunit.Volume) Liquid
	Container() LiquidContainer
	Add(v wunit.Volume)
	GetSmax() float64
	GetVisc() float64
	GetExtra() map[string]interface{}
	GetConc() float64
	GetCunit() string
	GetVunit() string
	GetStockConcentration() float64
}

// so far the best definition of this is not-solid-or-liquid...
type Gas interface {
	Physical
	Gas()
}

// to be composed with an X to make a SealedX
type Sealed interface {
	IsSealed()
}

type AnthaObject struct {
	ID   string
	Inst string
	Name string
}

func NewAnthaObject(name string) AnthaObject {
	id := GetUUID()
	ao := AnthaObject{id, "", name}
	return ao
}

/*************************/
