// wutil/numtoalpha.go: Part of the Antha language
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
	"math"
	"strings"
)

func NumToAlphaCountFromZero(n int) string {
	return NumToAlpha(n + 1)
}

// COUNTING FROM 1!
func NumToAlpha(n int) string {
	symbols := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := ""

	if n < 1 {
		return s
	}

	v := 26

	for {
		n -= 1
		r := (n % v)
		s += strings.Split(symbols, "")[r]
		n /= v
		if n <= 0 {
			break
		}
	}

	t := ""

	for i := len(s) - 1; i >= 0; i-- {
		t += string(s[i])
	}

	return t
}

func AlphaToNum(s string) int {
	symbols := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	v := 0

	for i, b := range s {
		c := string(b)
		x := int(math.Pow(float64(len(symbols)), float64(len(s)-(i+1)))) * (strings.Index(symbols, c) + 1)
		v += x
	}

	return v
}

func DecodeCoords(s string) (int, int) {
	var x, y int
	tox := strings.Split(s, ":")
	x = AlphaToNum(tox[0]) - 1
	y = ParseInt(tox[1]) - 1
	return x, y
}
