// /execute/types_test.go: Part of the Antha language
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

package execute

import (
	"encoding/json"
	"testing"
)

func TestSimpleValuesJSONBlockUnmarshalJSON(t *testing.T) {
	tx := `{"ID":"1", "A":5, "B":2, "Error": false }`

	idOne := ThreadID("1")
	falseBool := false
	ref := JSONBlock{
		ID:     &idOne,
		Error:  &falseBool,
		Values: map[string]interface{}{"A": 5, "B": 2},
	}

	var res JSONBlock
	err := json.Unmarshal([]byte(tx), &res)
	if err != nil {
		t.Errorf("Unable to unmarshal %s.", tx)
	}

	if string(*ref.ID) != string(*res.ID) {
		t.Errorf("IDs don't match. Expected %s, got %s.", string(*ref.ID), string(*res.ID))
	}
	if *ref.Error != *res.Error {
		t.Errorf("Errors don't match. Expected %s, got %s.", *ref.Error, *res.Error)
	}

	//let's compare the values
	t2, err := json.Marshal(ref.Values["A"])
	if err != nil {
		t.Errorf("Unable to marshal %v.", err)
	}
	var pint int
	err = json.Unmarshal(t2, &pint)

	if ref.Values["A"] != pint {
		t.Errorf("Values don't match. Expected %v, got %v.", ref.Values["A"], res.Values["A"])
	}
}
