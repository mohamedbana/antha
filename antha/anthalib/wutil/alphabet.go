// anthalib//wutil/displaymap.go: Part of the Antha language
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

package wutil

var Alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
	"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF"}

var (
	alphabet string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func MakeAlphabetArray() (alphabetarray []string) {

	alphabetarray = make([]string, 0)
	startercharacter := ""

	for j := 0; j < (len(alphabet)); j++ {

		for i := 0; i < (len(alphabet)); i++ {
			character := startercharacter + string(alphabet[i])

			alphabetarray = append(alphabetarray, character)
		}
		startercharacter = string(alphabet[j])

	}
	return
}
