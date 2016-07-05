// antha/AnthaStandardLibrary/Packages/enzymes/Find.go: Part of the Antha language
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

package search

import (
	"strings"
)

type Thingfound struct {
	Thing     string
	Positions []int
	Reverse   bool
}

// not perfect yet! issue with byte conversion of certain characters!
func Findall(bigthing string, smallthing string) (positions []int) {

	positions = make([]int, 0)
	count := strings.Count(bigthing, smallthing)
	//// fmt.Println("count", count)
	if count != 0 {

		pos := (strings.Index(bigthing, smallthing))
		restofbigthing := bigthing[(pos + 1):]
		//// fmt.Println("seq", bigthing)
		//// fmt.Println("rest,", restofbigthing)
		//pos = pos
		//// fmt.Println("pos = ", pos)
		for i := 0; i < count; i++ {
			//// fmt.Println("pos = ", pos)
			positions = append(positions, (pos + 1))
			//// fmt.Println("positions", positions)
			pos = pos + (strings.Index(restofbigthing, smallthing) + 1)
			//// fmt.Println("pos2 = ", pos)
			restofbigthing = bigthing[(pos + 1):]
			//// fmt.Println("rest2,", restofbigthing)
		}
	}
	return positions
}

func Findallthings(bigthing string, smallthings []string) (thingsfound []Thingfound) {
	var thingfound Thingfound
	thingsfound = make([]Thingfound, 0)

	for _, thing := range smallthings {
		if strings.Contains(bigthing, thing) {
			thingfound.Thing = thing
			thingfound.Positions = Findall(bigthing, thing)
			thingsfound = append(thingsfound, thingfound)
		}
	}
	return thingsfound
}

func Containsallthings(bigthing string, smallthings []string) (trueornot bool) {
	i := 0
	for _, thing := range smallthings {

		//	if strings.Contains(strings.ToUpper(bigthing), strings.ToUpper(thing)) {
		if strings.Contains(bigthing, thing) {
			i = i + 1
		}
	}
	if i == len(smallthings) {
		trueornot = true
	}

	return trueornot
}
