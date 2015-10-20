// antha/doc/testdata/a0.go: Part of the Antha language
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

// comment 0
package a

//BUG(uid): bug0

//TODO(uid): todo0

// A note with some spaces after it, should be ignored (watch out for
// emacs modes that remove trailing whitespace).
//NOTE(uid):

// SECBUG(uid): sec hole 0
// need to fix asap

// Multiple notes may be in the same comment group and should be
// recognized individually. Notes may start in the middle of a
// comment group as long as they start at the beginning of an
// individual comment.
//
// NOTE(foo): 1 of 4 - this is the first line of note 1
// - note 1 continues on this 2nd line
// - note 1 continues on this 3rd line
// NOTE(foo): 2 of 4
// NOTE(bar): 3 of 4
/* NOTE(bar): 4 of 4 */
// - this is the last line of note 4
//
//

// NOTE(bam): This note which contains a (parenthesized) subphrase
//            must appear in its entirety.

// NOTE(xxx) The ':' after the marker and uid is optional.

// NOTE(): NO uid - should not show up.
// NOTE()  NO uid - should not show up.
