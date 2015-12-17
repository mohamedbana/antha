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
// 2 Royal College St, London NW1 0NH UK

package factory

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/logger"
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
	A.CName = "Some component in factory"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "10x_M9Salts"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "100x_MEMVitamins"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "Yeast extract"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "Tryptone"
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
	A.CName = "Yellow"
	A.Type = "viscous"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "Blue"
	A.Type = "viscous"
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "Green"
	A.Type = "viscous"
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
	A.CName = "bsa"
	A.Type = "water"
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
	A.GenericMatter = matter["glycerol"]
	A.CName = "EcoRI"
	A.Type = "glycerol"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["glycerol"]
	A.CName = "EnzMastermix: 1/2 SapI; 1/2 T4 Ligase"
	A.Type = "glycerol"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "TypeIIsbuffer: 2/11 10xCutsmart; 1/11 1mM ATP; 8/11 Water"
	A.Type = "water"
	A.Smax = 9999
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
	A.CName = "SapI_Mastermix: 1/5 SapI; 1/5 T4 Ligase; 2/5 Cutsmart; 1/5 1mM ATP"
	A.Type = "water"
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "standard_cloning_vector_mark_1"
	A.Type = "dna"
	A.Smax = 1.0
	cmap[A.CName] = A
	// solutions needed for PCR example:

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["glycerol"]
	A.CName = "Q5Polymerase"
	A.Type = "glycerol"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["glycerol"]
	A.CName = "GoTaq_ green 2x mastermix"
	A.Type = "glycerol"
	A.Smax = 9999.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "DMSO"
	A.Type = "water"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "pET_GFP"
	A.Type = "water"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "HC"
	A.Type = "water"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "PrimerFw"
	A.Type = "dna"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "PrimerRev"
	A.Type = "dna"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "template_part"
	A.Type = "dna"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.GenericMatter = matter["water"]
	A.CName = "DNTPs"
	A.Type = "water"
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A
	return cmap
}

func GetComponentByType(typ string) *wtype.LHComponent {
	components := makeComponentLibrary()
	c := components[typ]
	if c == nil {
		logger.Fatal(fmt.Sprintf("Component %s not found", typ))
		panic(fmt.Errorf("Component %s not found", typ)) //TODO refactor to errors
	}
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
