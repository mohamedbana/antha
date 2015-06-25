// /anthalib/wunit/serialize.go: Part of the Antha language
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

package wunit

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (m *Volume) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", m.ToString())), nil
}

func (m *Volume) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var value float64
	var unit string
	if _, err := fmt.Fscanf(strings.NewReader(s), "%e%s", &value, &unit); err != nil {
		return err
	}
	*m = NewVolume(value, unit)
	return nil
}

func (m *Temperature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", m.ToString())), nil

}

func (m *Temperature) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var value float64
	var unit string
	if _, err := fmt.Fscanf(strings.NewReader(s), "%e%s", &value, &unit); err != nil {
		return err
	}
	*m = NewTemperature(value, unit)
	return nil
}

func (m *Concentration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", m.ToString())), nil

}

func (m *Concentration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var value float64
	var unit string
	if _, err := fmt.Fscanf(strings.NewReader(s), "%e%s", &value, &unit); err != nil {
		return err
	}
	*m = NewConcentration(value, unit)
	return nil
}

func (m *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", m.ToString())), nil

}

func (m *Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var value float64
	var unit string
	if _, err := fmt.Fscanf(strings.NewReader(s), "%e%s", &value, &unit); err != nil {
		return err
	}
	*m = NewTime(value, unit)
	return nil
}
