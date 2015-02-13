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
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/anthalib/wutil"
	"errors"
)


// mix needs to define the interface with liquid handling
// in order to do this it has to make the appropriate liquid handling
// request structure
// the functions in this package mostly convert the wtype representations
// into the simpler map representations which are the interface with the
// network layer
type SampleComponent map[string]interface{}


// take all of this liquid
func SampleAll(l wtype.Liquid)SampleComponent{
	return Sample(l, l.Volume())
}

// take a sample of volume v from this liquid
func Sample(l wtype.Liquid, v wunit.Volume) SampleComponent{
	ret:=make(map[string]interface{},5)

	ret["name"]=l.Name()
	ret["vol"]=v.RawValue()
	ret["vunit"]=v.Unit().PrefixedSymbol()

	return ret
}

// take a sample of this liquid and aim for a particular concentration
func SampleForConcentration(l wtype.Liquid, c wunit.Concentration)SampleComponent{
	ret:=make(map[string]interface{},5)
	ret["name"]=l.Name()
	ret["conc"]=c.RawValue()
	ret["cunit"]=c.Unit().PrefixedSymbol()
	return ret
}

// take a sample of this liquid to be used to make the solution up to 
// a particular total volume
func SampleForTotalVolume(l wtype.Liquid, v wunit.Volume)SampleComponent{
	ret:=make(map[string]interface{}, 5)
	ret["name"]=l.Name()
	ret["tvol"]=v.RawValue()
	ret["vunit"]=v.Unit().PrefixedSymbol()

	return ret
}

// mix the specified SampleComponents together
// the destination will be the location of the first sample
func Mix(components ...SampleComponent)SampleComponent{
		// this needs to be a container of some kind
		var d wtype.LiquidContainer

		if(components[0]["container"]!=nil){
			d=components[0]["container"].(wtype.LiquidContainer)
		} else{
			wutil.Error(errors.New("No container specified for mix"))
		}

		// d is essentially another SampleComponent

		return MixInto(d, components...)
}

// mix the specified SampleComponents together into the destination
// specified as the first argument
func MixInto(destination wtype.LiquidContainer, components ...SampleComponent)SampleComponent{
	// we must respect the order in which things are mixed.
	// the convention is that mix(X,Y) corresponds to "Add Y to X"

	ret:=make(map[string] interface{})

	// we use the first argument to specify the destination

	ret["containertype"]=destination.ContainerType()

	ret["components"]=components

	// this translates to the component ordering in the resulting solution
	for i,cmp:=range(components){
		cmp["order"]=i
	}

	return ret
}
