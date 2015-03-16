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
// 1 Royal College St, London NW1 0NH UK

package mixer

import (
	"errors"
	"github.com/antha-lang/antha/anthalib/liquidhandling"
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/anthalib/wutil"
)

// mix needs to define the interface with liquid handling
// in order to do this it has to make the appropriate liquid handling
// request structure

// take all of this liquid
func SampleAll(l wtype.Liquid) *liquidhandling.LHComponent {
	return Sample(l, l.Volume())
}

// take a sample of volume v from this liquid
func Sample(l wtype.Liquid, v wunit.Volume) *liquidhandling.LHComponent {
	ret := liquidhandling.NewLHComponent()

	ret.CName = l.Name()
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()

	return ret
}

// take a sample of this liquid and aim for a particular concentration
func SampleForConcentration(l wtype.Liquid, c wunit.Concentration) *liquidhandling.LHComponent {
	ret := liquidhandling.NewLHComponent()
	ret.CName = l.Name()
	ret.Conc = c.RawValue()
	ret.Cunit = c.Unit().PrefixedSymbol()
	return ret
}

// take a sample of this liquid to be used to make the solution up to
// a particular total volume
func SampleForTotalVolume(l wtype.Liquid, v wunit.Volume) *liquidhandling.LHComponent {
	ret := liquidhandling.NewLHComponent()
	ret.CName = l.Name()
	ret.Tvol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()

	return ret
}

// mix the specified liquidhandling.LHComponents together
// the destination will be the location of the first sample
func Mix(components ...*liquidhandling.LHComponent) *liquidhandling.LHSolution {
	// this needs to be a container of some kind
	var d wtype.LiquidContainer

	if components[0].Container != nil {
		d = components[0].Container
	} else {
		wutil.Error(errors.New("No container specified for mix"))
	}

	// d is essentially another liquidhandling.LHComponent

	return MixInto(d, components...)
}

// mix the specified liquidhandling.LHComponents together into the destination
// specified as the first argument
func MixInto(destination wtype.LiquidContainer, components ...*liquidhandling.LHComponent) *liquidhandling.LHSolution {
	// we must respect the order in which things are mixed.
	// the convention is that mix(X,Y) corresponds to "Add Y to X"

	ret := liquidhandling.NewLHSolution()

	// we use the first argument to specify the destination

	ret.ContainerType = destination.ContainerType()

	ret.Components = components

	// this translates to the component ordering in the resulting solution
	for i, cmp := range components {
		cmp.Order = i
	}

	return ret
}
