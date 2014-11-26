// antha/cmd/anthafmt/simplify.go: Part of the Antha language
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

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"reflect"
)

type simplifier struct {
	hasDotImport bool // package file contains: import . "some/import/path"
}

func (s *simplifier) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CompositeLit:
		// array, slice, and map composite literals may be simplified
		outer := n
		var eltType ast.Expr
		switch typ := outer.Type.(type) {
		case *ast.ArrayType:
			eltType = typ.Elt
		case *ast.MapType:
			eltType = typ.Value
		}

		if eltType != nil {
			typ := reflect.ValueOf(eltType)
			for i, x := range outer.Elts {
				px := &outer.Elts[i]
				// look at value of indexed/named elements
				if t, ok := x.(*ast.KeyValueExpr); ok {
					x = t.Value
					px = &t.Value
				}
				ast.Walk(s, x) // simplify x
				// if the element is a composite literal and its literal type
				// matches the outer literal's element type exactly, the inner
				// literal type may be omitted
				if inner, ok := x.(*ast.CompositeLit); ok {
					if match(nil, typ, reflect.ValueOf(inner.Type)) {
						inner.Type = nil
					}
				}
				// if the outer literal's element type is a pointer type *T
				// and the element is & of a composite literal of type T,
				// the inner &T may be omitted.
				if ptr, ok := eltType.(*ast.StarExpr); ok {
					if addr, ok := x.(*ast.UnaryExpr); ok && addr.Op == token.AND {
						if inner, ok := addr.X.(*ast.CompositeLit); ok {
							if match(nil, reflect.ValueOf(ptr.X), reflect.ValueOf(inner.Type)) {
								inner.Type = nil // drop T
								*px = inner      // drop &
							}
						}
					}
				}
			}

			// node was simplified - stop walk (there are no subnodes to simplify)
			return nil
		}

	case *ast.SliceExpr:
		// a slice expression of the form: s[a:len(s)]
		// can be simplified to: s[a:]
		// if s is "simple enough" (for now we only accept identifiers)
		if s.hasDotImport {
			// if dot imports are present, we cannot be certain that an
			// unresolved "len" identifier refers to the predefined len()
			break
		}
		if s, _ := n.X.(*ast.Ident); s != nil && s.Obj != nil {
			// the array/slice object is a single, resolved identifier
			if call, _ := n.High.(*ast.CallExpr); call != nil && len(call.Args) == 1 && !call.Ellipsis.IsValid() {
				// the high expression is a function call with a single argument
				if fun, _ := call.Fun.(*ast.Ident); fun != nil && fun.Name == "len" && fun.Obj == nil {
					// the function called is "len" and it is not locally defined; and
					// because we don't have dot imports, it must be the predefined len()
					if arg, _ := call.Args[0].(*ast.Ident); arg != nil && arg.Obj == s.Obj {
						// the len argument is the array/slice object
						n.High = nil
					}
				}
			}
		}
		// Note: We could also simplify slice expressions of the form s[0:b] to s[:b]
		//       but we leave them as is since sometimes we want to be very explicit
		//       about the lower bound.
		// An example where the 0 helps:
		//       x, y, z := b[0:2], b[2:4], b[4:6]
		// An example where it does not:
		//       x, y := b[:n], b[n:]

	case *ast.RangeStmt:
		// a range of the form: for x, _ = range v {...}
		// can be simplified to: for x = range v {...}
		if ident, _ := n.Value.(*ast.Ident); ident != nil && ident.Name == "_" {
			n.Value = nil
		}
	}

	return s
}

func simplify(f *ast.File) {
	var s simplifier

	// determine if f contains dot imports
	for _, imp := range f.Imports {
		if imp.Name != nil && imp.Name.Name == "." {
			s.hasDotImport = true
			break
		}
	}

	ast.Walk(&s, f)
}
