// wtype/wtype_test.go: Part of the Antha language
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
	"fmt"
	"testing"
)

func testBS(bs BioSequence) {
	fmt.Println(bs.Sequence())
}

func TestOne(*testing.T) {
	dna := DNASequence{"test", "ACCACACATAGCTAGCTAGCTAG", false, false, Overhang{}, Overhang{}, ""}
	testBS(&dna)
}

func ExampleOne() {
	dna := DNASequence{"test", "ACCACACATAGCTAGCTAGCTAG", false, false, Overhang{}, Overhang{}, ""}
	testBS(&dna)
	// Output:
	// ACCACACATAGCTAGCTAGCTAG
}

func TestLocations(*testing.T) {
	nl := NewLocation("liquidhandler", 9, NewShape("box", "", 0, 0, 0))
	nl2 := NewLocation("anotherliquidhandler", 9, NewShape("box", "", 0, 0, 0))
	fmt.Println("Location ", nl.Location_Name(), " ", nl.Location_ID(), " and location ", nl2.Location_Name(), " ", nl2.Location_ID(), " are the same? ", SameLocation(nl, nl2, 0))

	fmt.Println("Location ", nl.Positions()[0].Location_Name(), " and location ", nl.Positions()[1].Location_Name(), " are the same? ", SameLocation(nl.Positions()[0], nl.Positions()[1], 0), " share a parent? ", SameLocation(nl.Positions()[0], nl.Positions()[1], 1))

	fmt.Println("Locations ", nl.Location_Name(), " and ", nl.Positions()[0].Location_Name(), " share a parent? ", SameLocation(nl, nl.Positions()[0], 1))
}

func TestWellCoords(*testing.T) {
	fmt.Println("Testing Well Coords")
	wc := MakeWellCoordsA1("A1")
	fmt.Println(wc.FormatA1())
	fmt.Println(wc.Format1A())
	fmt.Println(wc.FormatXY())
	fmt.Println(wc.X, " ", wc.Y)
	wc = MakeWellCoordsXYsep("X1", "Y1")
	fmt.Println(wc.FormatA1())
	fmt.Println(wc.Format1A())
	fmt.Println(wc.FormatXY())
	fmt.Println(wc.X, " ", wc.Y)
	wc = MakeWellCoordsXY("X1Y1")
	fmt.Println(wc.FormatA1())
	fmt.Println(wc.Format1A())
	fmt.Println(wc.FormatXY())
	fmt.Println(wc.X, " ", wc.Y)
	fmt.Println("Finished Testing Well Coords")
}

func TestLHComponentSampleStuff(t *testing.T) {
	var c LHComponent

	faux := c.IsSample()

	if faux {
		t.Fatal("IsSample() must return false on new components")
	}

	c.SetIsSample(true)

	vrai := c.IsSample()

	if !vrai {
		t.Fatal("IsSample() must return true following SetIsSample(true)")
	}

	c.SetIsSample(false)

	faux = c.IsSample()

	if faux {
		t.Fatal("IsSample() must return false following SetIsSample(false)")
	}

	// now the same from NewLHComponent

	c2 := NewLHComponent()

	faux = c2.IsSample()

	if faux {
		t.Fatal("IsSample() must return false on new components")
	}

	c2.SetIsSample(true)

	vrai = c2.IsSample()

	if !vrai {
		t.Fatal("IsSample() must return true following SetIsSample(true)")
	}

	c2.SetIsSample(false)

	faux = c2.IsSample()

	if faux {
		t.Fatal("IsSample() must return false following SetIsSample(false)")
	}

	// finally need to make sure sample works
	// grrr import cycle not allowed: honestly I think Sample should just be an
	// instance method of LHComponent now anyway
	/*

		c.CName = "YOMAMMA"
		s := mixer.Sample(c, wunit.NewVolume(10.0, "ul"))

		vrai = s.IsSample()

		if !vrai {
			t.Fatal("IsSample() must return true for results of mixer.Sample()")
		}
		s = mixer.SampleForConcentration(c, wunit.NewConcentration(10.0, "mol/l"))

		vrai = s.IsSample()

		if !vrai {
			t.Fatal("IsSample() must return true for results of mixer.SampleForConcentration()")
		}
		s = mixer.SampleForTotalVolume(c, wunit.NewVolume(10.0, "ul"))

		vrai = s.IsSample()

		if !vrai {
			t.Fatal("IsSample() must return true for results of mixer.SampleForTotalVolume()")
		}
	*/
}
