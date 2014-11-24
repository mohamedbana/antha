// antha/doc/testdata/c.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
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


package c

import "a"

// ----------------------------------------------------------------------------
// Test that empty declarations don't cause problems

const ()

type ()

var ()

// ----------------------------------------------------------------------------
// Test that types with documentation on both, the Decl and the Spec node
// are handled correctly.

// A (should see this)
type A struct{}

// B (should see this)
type (
	B struct{}
)

type (
	// C (should see this)
	C struct{}
)

// D (should not see this)
type (
	// D (should see this)
	D struct{}
)

// E (should see this for E2 and E3)
type (
	// E1 (should see this)
	E1 struct{}
	E2 struct{}
	E3 struct{}
	// E4 (should see this)
	E4 struct{}
)

// ----------------------------------------------------------------------------
// Test that local and imported types are different when
// handling anonymous fields.

type T1 struct{}

func (t1 *T1) M() {}

// T2 must not show methods of local T1
type T2 struct {
	a.T1 // not the same as locally declared T1
}