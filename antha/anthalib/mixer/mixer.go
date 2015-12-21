// mixer/mixer.go: Part of the Antha language
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

package mixer

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// mix needs to define the interface with liquid handling
// in order to do this it has to make the appropriate liquid handling
// request structure

// take all of this liquid
func SampleAll(l wtype.Liquid) *wtype.LHComponent {
	return Sample(l, l.Volume())
}

// below need to account for having locations for liquids specified...

// take a sample of volume v from this liquid
func Sample(l wtype.Liquid, v wunit.Volume) *wtype.LHComponent {
	ret := wtype.NewLHComponent()

	ret.CName = l.Name()
	ret.Type = l.GetType()
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	if l.Container() != nil {
		ret.LContainer = l.Container().(*wtype.LHWell)
	}

	return ret
}

// take a sample of this liquid and aim for a particular concentration
func SampleForConcentration(l wtype.Liquid, c wunit.Concentration) *wtype.LHComponent {
	ret := wtype.NewLHComponent()
	ret.CName = l.Name()
	ret.Type = l.GetType()
	ret.Conc = c.RawValue()
	ret.Cunit = c.Unit().PrefixedSymbol()
	ret.CName = l.Name()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	ret.LContainer = l.Container().(*wtype.LHWell)
	return ret
}

func SampleMass(s wtype.Liquid, m wunit.Mass, d wunit.Density) *wtype.LHComponent {

	// calculate volume to add from density
	v := wunit.MasstoVolume(m, d)

	ret := wtype.NewLHComponent()
	ret.CName = s.Name()
	ret.Type = s.GetType()
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = s.GetExtra()
	ret.Smax = s.GetSmax()
	ret.Visc = s.GetVisc()
	if s.Container() != nil {
		ret.LContainer = s.Container().(*wtype.LHWell)
	}

	return ret
}

// take a sample of this liquid to be used to make the solution up to
// a particular total volume
func SampleForTotalVolume(l wtype.Liquid, v wunit.Volume) *wtype.LHComponent {
	ret := wtype.NewLHComponent()
	ret.CName = l.Name()
	ret.Type = l.GetType()
	ret.Tvol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.LContainer = l.Container().(*wtype.LHWell)
	ret.CName = l.Name()
	ret.Extra = l.GetExtra()
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()

	return ret
}

func SampleSolidtoLiquid(s wtype.Powder, m wunit.Mass, d wunit.Density) *wtype.LHComponent {

	// calculate volume to add from density
	v := wunit.MasstoVolume(m, d)

	ret := wtype.NewLHComponent()
	ret.CName = s.Name()
	ret.Type = s.GetType()
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = s.GetExtra()
	ret.Smax = s.GetSmax()
	ret.Visc = s.GetVisc()
	if s.Container() != nil {
		ret.LContainer = s.Container().(*wtype.LHWell)
	}

	return ret
}

// Temp hack to mix solutions
func MixLiquidstemp(liquids ...*wtype.LHSolution) *wtype.LHSolution {
	// we must respect the order in which things are mixed.
	// the convention is that mix(X,Y) corresponds to "Add Y to X"

	ret := wtype.NewLHSolution()

	ret.Components = make([]*wtype.LHComponent, 0)
	for _, liquid := range liquids {
		for _, component := range liquid.Components {
			ret.Components = append(ret.Components, component)
		}
	}
	// this translates to the component ordering in the resulting solution
	for i, cmp := range liquids {
		cmp.Order = i
	}

	return ret
}

// mix the specified wtype.LHComponents together
// and leave the destination TBD
func Mix(components ...*wtype.LHComponent) *wtype.LHSolution {
	// we must respect the order in which things are mixed.
	// the convention is that mix(X,Y) corresponds to "Add Y to X"

	ret := wtype.NewLHSolution()

	ret.Components = components

	// this translates to the component ordering in the resulting solution
	for i, cmp := range components {
		cmp.Order = i
	}

	return ret
}

// mix the specified wtype.LHComponents together into the destination
// specified as the first argument
func MixInto(destination *wtype.LHPlate, components ...*wtype.LHComponent) *wtype.LHSolution {
	// we must respect the order in which things are mixed.
	// the convention is that mix(X,Y) corresponds to "Add Y to X"

	ret := wtype.NewLHSolution()

	// we use the first argument to specify the destination

	ret.ContainerType = destination.Type
	ret.Platetype = destination.Type

	ret.Components = components

	// this translates to the component ordering in the resulting solution
	for i, cmp := range components {
		cmp.Order = i
	}

	return ret
}

// mix the specified wtype.LHComponents together into the destination
// specified as the first argument
func MixTo(destination *wtype.LHPlate, address string, components ...*wtype.LHComponent) *wtype.LHSolution {
	// we must respect the order in which things are mixed.
	// the convention is that mix(X,Y) corresponds to "Add Y to X"

	ret := wtype.NewLHSolution()

	// we use the first argument to specify the destination

	ret.ContainerType = destination.Type
	ret.Platetype = destination.Type
	ret.Welladdress = address
	ret.Components = components

	// this translates to the component ordering in the resulting solution
	for i, cmp := range components {
		cmp.Order = i
	}

	return ret
}
