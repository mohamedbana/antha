// antha/doc/testdata/d2.go: Part of the Antha language
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

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test cases for sort order of declarations.

package d

// C1 should be second.
const C1 = 1

// C0 should be first.
const C0 = 0

// V1 should be second.
var V1 uint

// V0 should be first.
var V0 uintptr

// CAx constants should appear after CBx constants.
const (
	CA2 = iota // before CA1
	CA1        // before CA0
	CA0        // at end
)

// VAx variables should appear after VBx variables.
var (
	VA2 int // before VA1
	VA1 int // before VA0
	VA0 int // at end
)

// T1 should be second.
type T1 struct{}

// T0 should be first.
type T0 struct{}

// F1 should be second.
func F1() {}

// F0 should be first.
func F0() {}
