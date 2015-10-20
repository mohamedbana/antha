// antha/doc/testdata/d1.go: Part of the Antha language
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
// 2 Royal College St, London NW1 0NH UK

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test cases for sort order of declarations.

package d

// C2 should be third.
const C2 = 2

// V2 should be third.
var V2 int

// CBx constants should appear before CAx constants.
const (
	CB2 = iota // before CB1
	CB1        // before CB0
	CB0        // at end
)

// VBx variables should appear before VAx variables.
var (
	VB2 int // before VB1
	VB1 int // before VB0
	VB0 int // at end
)

const (
	// Single const declarations inside ()'s are considered ungrouped
	// and show up in sorted order.
	Cungrouped = 0
)

var (
	// Single var declarations inside ()'s are considered ungrouped
	// and show up in sorted order.
	Vungrouped = 0
)

// T2 should be third.
type T2 struct{}

// Grouped types are sorted nevertheless.
type (
	// TG2 should be third.
	TG2 struct{}

	// TG1 should be second.
	TG1 struct{}

	// TG0 should be first.
	TG0 struct{}
)

// F2 should be third.
func F2() {}
