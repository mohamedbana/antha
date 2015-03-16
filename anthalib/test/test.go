// anthalib//test/test.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package main

import (
	"github.com/antha-lang/antha/anthalib/wunit"
	"github.com/antha-lang/antha/anthalib/wtype"
	"github.com/antha-lang/antha/anthalib/mixer"
	"github.com/antha-lang/antha/anthalib/execution"
)

func main(){
	// this still seems quite heavyweight

	// make two liquids
	w:=wtype.NewGenericLiquid("water", "water", wunit.NewVolume(0.001,"L"))
	t:=wtype.NewGenericLiquid("tartrazine", "water", wunit.NewVolume(0.0001,"L"))
	p:=wtype.GenericSBSFormatPlate()

	// define a volume and a concentration
	v:=wunit.NewVolume(0.0001, "L")
	c:=wunit.NewConcentration(100, "g/l")

	// now try to mix them together
	m:=mixer.Mix(mixer.SampleForTotalVolume(w, v), mixer.SampleForConcentration(t, c))

	m=m

	context:=execution.GetContext()

	context=context
}




