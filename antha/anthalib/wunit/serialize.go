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
// 2 Royal College St, London NW1 0NH UK

package wunit

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type stringer interface {
	ToString() string
}

func marshal(x stringer) ([]byte, error) {
	var s *string
	if x != nil {
		r := x.ToString()
		s = &r
	}
	return json.Marshal(s)
}

func unmarshal(b []byte) (value float64, unit string, err error) {
	var s *string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	} else if s == nil {
		return
	}
	if _, err = fmt.Fscanf(strings.NewReader(*s), "%e%s", &value, &unit); err != nil {
		if err == io.EOF {
			err = nil
			unit = ""
			if _, err = fmt.Fscanf(strings.NewReader(*s), "%e", &value); err != nil {
				return
			}
		}
		return
	}
	return
}

func (m Volume) MarshalJSON() ([]byte, error) {
	return marshal(m)
}

func (m *Volume) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewVolume(value, unit)
	}
	return nil
}

func (m *Temperature) MarshalJSON() ([]byte, error) {
	return marshal(m)
}

func (m *Temperature) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewTemperature(value, unit)
	}
	return nil
}

func (m *Concentration) MarshalJSON() ([]byte, error) {
	return marshal(m)
}

func (m *Concentration) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewConcentration(value, unit)
	}
	return nil
}

func (m *Time) MarshalJSON() ([]byte, error) {
	return marshal(m)

}

func (m *Time) UnmarshalJSON(b []byte) error {
	if value, unit, err := unmarshal(b); err != nil {
		return err
	} else if unit != "" {
		*m = NewTime(value, unit)
	}
	return nil
}
