// antha/doc/testdata/e.go: Part of the Antha language
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

// The package e is a go/doc test for embedded methods.
package e

// ----------------------------------------------------------------------------
// Conflicting methods M must not show up.

type t1 struct{}

// t1.M should not appear as method in a Tx type.
func (t1) M() {}

type t2 struct{}

// t2.M should not appear as method in a Tx type.
func (t2) M() {}

// T1 has no embedded (level 1) M method due to conflict.
type T1 struct {
	t1
	t2
}

// ----------------------------------------------------------------------------
// Higher-level method M wins over lower-level method M.

// T2 has only M as top-level method.
type T2 struct {
	t1
}

// T2.M should appear as method of T2.
func (T2) M() {}

// ----------------------------------------------------------------------------
// Higher-level method M wins over lower-level conflicting methods M.

type t1e struct {
	t1
}

type t2e struct {
	t2
}

// T3 has only M as top-level method.
type T3 struct {
	t1e
	t2e
}

// T3.M should appear as method of T3.
func (T3) M() {}

// ----------------------------------------------------------------------------
// Don't show conflicting methods M embedded via an exported and non-exported
// type.

// T1 has no embedded (level 1) M method due to conflict.
type T4 struct {
	t2
	T2
}

// ----------------------------------------------------------------------------
// Don't show embedded methods of exported anonymous fields unless AllMethods
// is set.

type T4 struct{}

// T4.M should appear as method of T5 only if AllMethods is set.
func (*T4) M() {}

type T5 struct {
	T4
}

// ----------------------------------------------------------------------------
// Recursive type declarations must not lead to endless recursion.

type U1 struct {
	*U1
}

// U1.M should appear as method of U1.
func (*U1) M() {}

type U2 struct {
	*U3
}

// U2.M should appear as method of U2 and as method of U3 only if AllMethods is set.
func (*U2) M() {}

type U3 struct {
	*U2
}

// U3.N should appear as method of U3 and as method of U2 only if AllMethods is set.
func (*U3) N() {}

type U4 struct {
	*u5
}

// U4.M should appear as method of U4.
func (*U4) M() {}

type u5 struct {
	*U4
}

// ----------------------------------------------------------------------------
// A higher-level embedded type (and its methods) wins over the same type (and
// its methods) embedded at a lower level.

type V1 struct {
	*V2
	*V5
}

type V2 struct {
	*V3
}

type V3 struct {
	*V4
}

type V4 struct {
	*V5
}

type V5 struct {
	*V6
}

type V6 struct{}

// V4.M should appear as method of V2 and V3 if AllMethods is set.
func (*V4) M() {}

// V6.M should appear as method of V1 and V5 if AllMethods is set.
func (*V6) M() {}
