// /anthalib/factory/make_component_library.go: Part of the Antha language
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

package factory

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func makeComponentLibrary() map[string]*wtype.LHComponent {
	matter := wtype.MakeMatterLib()

	cmap := make(map[string]*wtype.LHComponent)

	A := wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "water"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "tartrazine"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "DNAsolution"
	A.Type = "dna"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["glycerol"]
	A.CName = "restrictionenzyme"
	A.Type = "glycerol"
	A.Smax = 100
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "dna_part"
	A.Type = "dna"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["glycerol"]
	A.CName = "SapI"
	A.Type = "glycerol"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["glycerol"]
	A.CName = "T4Ligase"
	A.Type = "glycerol"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "CutsmartBuffer"
	A.Type = "water"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "ATP"
	A.Type = "water"
	A.Smax = 5.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "standard_cloning_vector_mark_1"
	A.Type = "water"
	A.Smax = 1.0
	cmap[A.CName] = A

	return cmap
}

func GetComponentByType(typ string) *wtype.LHComponent {
	components := makeComponentLibrary()
	c := components[typ]
	return c.Dup()
}

func GetComponentList() []string {
	components := makeComponentLibrary()
	kz := make([]string, len(components))
	x := 0
	for name, _ := range components {
		kz[x] = name
		x += 1
	}
	return kz

}
