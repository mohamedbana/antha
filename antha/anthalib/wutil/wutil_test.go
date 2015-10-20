// wutil/wutil_test.go: Part of the Antha language
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

package wutil

import (
	"fmt"
	"testing"
)

func TestRound(*testing.T) {
	//ExampleRound()
}

func ExampleRound() {
	fmt.Println(RoundInt(-3.3))
	fmt.Println(RoundInt(3.3))
	fmt.Println(RoundInt(-3.9))
	fmt.Println(RoundInt(3.9))
	fmt.Println(RoundInt(-3.5))
	fmt.Println(RoundInt(3.5))

	// Output:
	// -3
	// 3
	// -4
	// 4
	// -4
	// 4
}

func TestNumToAlpha(*testing.T) {
	//ExampleNTA()
}

func ExampleNTA() {
	a := NumToAlpha(10)
	fmt.Println("10: ", a)
	a = NumToAlpha(1)
	fmt.Println("1: ", a)
	a = NumToAlpha(2)
	fmt.Println("2: ", a)
	a = NumToAlpha(27)
	fmt.Println("27: ", a)
	a = NumToAlpha(100)
	fmt.Println("100: ", a)

	// Output:
	// 10:  J
	// 1:  A
	// 2:  B
	// 27:  AA
	// 100:  CV
}

func TestAlphaToNum(*testing.T) {
	//ExampleATN()
}

func ExampleATN() {
	a := AlphaToNum("J")
	fmt.Println("J: ", a)
	a = AlphaToNum("A")
	fmt.Println("A: ", a)
	a = AlphaToNum("B")
	fmt.Println("B: ", a)
	a = AlphaToNum("CV")
	fmt.Println("CV: ", a)
	a = AlphaToNum("AA")
	fmt.Println("AA: ", a)

	// Output:
	// J:  10
	// A:  1
	// B:  2
	// CV:  100
	// AA:  27
}
