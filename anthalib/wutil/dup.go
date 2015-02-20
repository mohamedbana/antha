// anthalib//wutil/dup.go: Part of the Antha language
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

package wutil

import "github.com/antha-lang/antha/anthalib/wtype"
import "encoding/json"

// functions to copy objects about
// be VERY CAREFUL with dup - it duplicates inst and id so it breaks
// many of the assumptions tied to these if not used with care

// an EXACT duplicate, including UUIDs which ordinarily should change
func Dup(s1 map[string]interface{}) map[string]interface{} {

	s2 := make(map[string]interface{}, len(s1))

	b, err := json.Marshal(s1)

	// proper error handling needed
	if err != nil {
		Error(err)
	}

	err = json.Unmarshal(b, &s2)

	// proper error handling needed
	if err != nil {
		Error(err)
	}

	return s2
}

// this changes the UUID appropriately; more likely to be useful in practice
// also ensures "inst" is not copied since we cannot physically duplicate anything
func Copy(s1 map[string]interface{}) map[string]interface{} {
	s2 := Dup(s1)
	s2["id"] = wtype.GetUUID()
	s2["inst"] = nil
	return s2
}
