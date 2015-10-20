// antha/doc/synopsis_test.go: Part of the Antha language
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

package doc

import "testing"

var tests = []struct {
	txt string
	fsl int
	syn string
}{
	{"", 0, ""},
	{"foo", 3, "foo"},
	{"foo.", 4, "foo."},
	{"foo.bar", 7, "foo.bar"},
	{"  foo.  ", 6, "foo."},
	{"  foo\t  bar.\n", 12, "foo bar."},
	{"  foo\t  bar.\n", 12, "foo bar."},
	{"a  b\n\nc\r\rd\t\t", 12, "a b c d"},
	{"a  b\n\nc\r\rd\t\t  . BLA", 15, "a b c d ."},
	{"Package poems by T.S.Eliot. To rhyme...", 27, "Package poems by T.S.Eliot."},
	{"Package poems by T. S. Eliot. To rhyme...", 29, "Package poems by T. S. Eliot."},
	{"foo implements the foo ABI. The foo ABI is...", 27, "foo implements the foo ABI."},
	{"Package\nfoo. ..", 12, "Package foo."},
	{"P . Q.", 3, "P ."},
	{"P. Q.   ", 8, "P. Q."},
	{"Package Καλημέρα κόσμε.", 36, "Package Καλημέρα κόσμε."},
	{"Package こんにちは 世界\n", 31, "Package こんにちは 世界"},
	{"Package こんにちは。世界", 26, "Package こんにちは。"},
	{"Package 안녕．世界", 17, "Package 안녕．"},
	{"Package foo does bar.", 21, "Package foo does bar."},
	{"Copyright 2012 Google, Inc. Package foo does bar.", 27, ""},
	{"All Rights reserved. Package foo does bar.", 20, ""},
	{"All rights reserved. Package foo does bar.", 20, ""},
	{"Authors: foo@bar.com. Package foo does bar.", 21, ""},
}

func TestSynopsis(t *testing.T) {
	for _, e := range tests {
		fsl := firstSentenceLen(e.txt)
		if fsl != e.fsl {
			t.Errorf("got fsl = %d; want %d for %q\n", fsl, e.fsl, e.txt)
		}
		syn := Synopsis(e.txt)
		if syn != e.syn {
			t.Errorf("got syn = %q; want %q for %q\n", syn, e.syn, e.txt)
		}
	}
}
