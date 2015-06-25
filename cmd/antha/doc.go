// antha/cmd/andtha/doc.go: Part of the Antha language
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

/*
Antha is the primary cross-compiler for antha.
It parses Antha Element definitions to create the underyling Go source
which is in turn compiled using the regular Go tool chain.

Without an explicit input parameter it parses from stdin, and will output to stdout.
Given a directory or antha source file, it will generate .go files for each antha input.

Usage:
	anthamt [flags] [path ...]

The flags are:
	-trace
		Shows the entire AST while parsing, to debug parse errors
	-errors
		Print all (including spurious) errors.
*/
package main
