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
// 2 Royal College St, London NW1 0NH UK

package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
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

type WaterSolution *LHComponent

type NonwaterSolution *LHComponent
