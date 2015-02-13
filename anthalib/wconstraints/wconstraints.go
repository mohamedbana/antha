// wconstraints/wconstraints.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package wconstraints

import (
	"github.com/antha-lang/antha/anthalib/wunit"
	"time"
)


// general constraint handling interface
type Constraints interface{
	Add(c Constraint) 
	Remove(c Constraint) 
}

// Interface to define a constraint 
// -- the main requirement is that the constraint can be checked
type Constraint interface{
	// constraints return true when violated
	Check(i interface {}) bool
}

// time constraint - simply keeps track of a duration
// by recording when the constraint was created and how long it is
type TimeConstraint struct{
	Start time.Time
	Length time.Duration
}

// test whether time is up
func (tc *TimeConstraint)Check(i interface{}) bool{
	t:=i.(time.Time)
	d:=t.Sub(tc.Start)
	if d > tc.Length{
		return true
	} else{
		return false
	}
}

// Temperature constraint -- sample must be kept within
// the specified limits
// defines a max and a min temperature
type TempConstraint struct{
	Min wunit.Temperature
	Max wunit.Temperature
}

// Test the constraint against the current temperature
func (tc TempConstraint)Check(i interface{})bool{
	t:=i.(wunit.Temperature)

	v:=t.Value()

	if v<tc.Min.Value() || v > tc.Max.Value(){
		return true
	}

	return false
}

