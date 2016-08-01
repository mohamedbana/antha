// Part of the Antha language
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

// prints pre-formatted colour text using ansi codes
package text

import (
	"fmt"

	"github.com/mgutz/ansi"
)

// Prints a description highlighted in red followed by the value in unformatted text
func Print(description string, value interface{}) (fmtd string) {

	fmtd = fmt.Sprintln(ansi.Color(description, "red"), value)
	return
}

func Annotate(stringtoannotate string, colour string) (fmtd string) {
	// currently leads to undesired behaviour with string manipulation
	//fmtd = fmt.Sprint(ansi.Color(stringtoannotate, colour))
	_ = colour
	fmtd = stringtoannotate
	return
}

/*
func Printfield( value interface{}) (fmtd string) {

	switch myValue := value.(type){
		case string:
		fmt.Println(myValue)
		case Hit:
		fmt.Printf("%+v\n", myValue)
		default:
		// fmt.Println("Type not handled: ", reflect.TypeOf(value))
	}

	//a := &Hsp{Len: "afoo"}
	val := reflect.Indirect(reflect.ValueOf(value))
	desc := fmt.Sprint(val.Type().Field(0).Name)

	fmtd = fmt.Sprint(ansi.Color(desc, "red"), value)
	return
}
*/
