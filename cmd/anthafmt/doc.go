// antha/cmd/anthafmt/doc.go: Part of the Antha language
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

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Anthafmt formats Antha programs.
It uses tabs (width = 8) for indentation and blanks for alignment.

Without an explicit path, it processes the standard input.  Given a file,
it operates on that file; given a directory, it operates on all .go files in
that directory, recursively.  (Files starting with a period are ignored.)
By default, anthamt prints the reformatted sources to standard output.

Usage:
	anthamt [flags] [path ...]

The flags are:
	-d
		Do not print reformatted sources to standard output.
		If a file's formatting is different than anthamt's, print diffs
		to standard output.
	-e
		Print all (including spurious) errors.
	-l
		Do not print reformatted sources to standard output.
		If a file's formatting is different from anthamt's, print its name
		to standard output.
	-r rule
		Apply the rewrite rule to the source before reformatting.
	-s
		Try to simplify code (after applying the rewrite rule, if any).
	-w
		Do not print reformatted sources to standard output.
		If a file's formatting is different from anthamt's, overwrite it
		with anthamt's version.

Debugging support:
	-cpuprofile filename
		Write cpu profile to the specified file.


The rewrite rule specified with the -r flag must be a string of the form:

	pattern -> replacement

Both pattern and replacement must be valid Go/Antha expressions.
In the pattern, single-character lowercase identifiers serve as
wildcards matching arbitrary sub-expressions; those expressions
will be substituted for the same identifiers in the replacement.

When anthamt reads from standard input, it accepts either a full Antha program
or a program fragment.  A program fragment must be a syntactically
valid declaration list, statement list, or expression.  When formatting
such a fragment, anthamt preserves leading indentation as well as leading
and trailing spaces, so that individual sections of a Go/Antha program can be
formatted by piping them through anthamt.

Examples

To check files for unnecessary parentheses:

	anthamt -r '(a) -> a' -l *.go

To remove the parentheses:

	anthamt -r '(a) -> a' -w *.go

To convert the package tree from explicit slice upper bounds to implicit ones:

	anthamt -r 'α[β:len(α)] -> α[β:]' -w $GOROOT/src/pkg

The simplify command

When invoked with -s anthamt will make the following source transformations where possible.

	An array, slice, or map composite literal of the form:
		[]T{T{}, T{}}
	will be simplified to:
		[]T{ { }, { } }

	A slice expression of the form:
		s[a:len(s)]
	will be simplified to:
		s[a:]

	A range of the form:
		for x, _ = range v {...}
	will be simplified to:
		for x = range v {...}
*/
package main

// BUG(rsc): The implementation of -r is a bit slow.
