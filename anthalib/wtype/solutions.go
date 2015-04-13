// wtype/solutions.go: Part of the Antha language
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
	"github.com/antha-lang/antha/anthalib/wunit"
)

// interface defining a solution type. Defined as
// being a liquid, having a concentration, a solvent
// and one or more solutes
type Solution interface {
	Liquid
	Concentration() wunit.Concentration
	ConcentrationOf(s string) wunit.Concentration
	Solvent() Liquid
	Solutes() []Physical
}

// interface to define a suspension type
// this also has a solvent and solutes but no concentration
type Suspension interface {
	Liquid
	Solvent() Liquid
	Solutes() []Physical
}

// a solution with water as the solvent
type WaterSolution struct {
	GenericLiquid
	Sltes []GenericPhysical
}

// constructor

func NewWaterSolution(v wunit.Volume, solutenames []string, solutetypes []string, soluteamounts []wunit.Mass) *WaterSolution {
	l := NewGenericLiquid("water", "water", v, nil)
	as := make([]GenericPhysical, len(solutenames))

	for i := 0; i < len(solutenames); i++ {
		s := NewGenericPhysical(solutetypes[i])
		s.SetName(solutenames[i])
		s.SetMass(soluteamounts[i])
		as = append(as, s)
	}
	w := WaterSolution{l, as}

	return &w
}

func (as *WaterSolution) Solvent() Liquid {
	return &(as.GenericLiquid)
}

func (as *WaterSolution) Solutes() []Physical {
	return as.Sltes
}

func (as *WaterSolution) Concentration() wunit.Concentration {
	// only defined if there is exactly one solute

	if len(as.Sltes) != 1 {
		panic("Wtype error: Cannot use Concentration for multiple solutes")
	}

	return as.ConcentrationOf(as.Sltes[0].MatterType())
}

func (as *WaterSolution) ConcentrationOf(s string) wunit.Concentration {
	v := as.GenericLiquid.Volume()
	var m wunit.Mass
	for _, l := range as.Sltes {
		if l.MatterType() == s {
			m = l.Mass()
			break
		}
	}
	cnc := wunit.NewConcentration(m.RawValue()/v.RawValue(), "g/l")
	return cnc
}

// a structure defining a non water solution
type NonWaterSolution struct {
	GenericLiquid
	Sltes []GenericPhysical
}

func (nas *NonWaterSolution) Solvent() Liquid {
	return &nas.GenericLiquid
}

func (nas *NonWaterSolution) Solutes() []Physical {
	return nas.Sltes
}
