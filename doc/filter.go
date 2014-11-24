// antha/doc/filter.go: Part of the Antha language
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


package doc

import "github.com/antha-lang/antha/ast"

type Filter func(string) bool

func matchFields(fields *ast.FieldList, f Filter) bool {
	if fields != nil {
		for _, field := range fields.List {
			for _, name := range field.Names {
				if f(name.Name) {
					return true
				}
			}
		}
	}
	return false
}

func matchDecl(d *ast.GenDecl, f Filter) bool {
	for _, d := range d.Specs {
		switch v := d.(type) {
		case *ast.ValueSpec:
			for _, name := range v.Names {
				if f(name.Name) {
					return true
				}
			}
		case *ast.TypeSpec:
			if f(v.Name.Name) {
				return true
			}
			switch t := v.Type.(type) {
			case *ast.StructType:
				if matchFields(t.Fields, f) {
					return true
				}
			case *ast.InterfaceType:
				if matchFields(t.Methods, f) {
					return true
				}
			}
		}
	}
	return false
}

func filterValues(a []*Value, f Filter) []*Value {
	w := 0
	for _, vd := range a {
		if matchDecl(vd.Decl, f) {
			a[w] = vd
			w++
		}
	}
	return a[0:w]
}

func filterFuncs(a []*Func, f Filter) []*Func {
	w := 0
	for _, fd := range a {
		if f(fd.Name) {
			a[w] = fd
			w++
		}
	}
	return a[0:w]
}

func filterTypes(a []*Type, f Filter) []*Type {
	w := 0
	for _, td := range a {
		n := 0 // number of matches
		if matchDecl(td.Decl, f) {
			n = 1
		} else {
			// type name doesn't match, but we may have matching consts, vars, factories or methods
			td.Consts = filterValues(td.Consts, f)
			td.Vars = filterValues(td.Vars, f)
			td.Funcs = filterFuncs(td.Funcs, f)
			td.Methods = filterFuncs(td.Methods, f)
			n += len(td.Consts) + len(td.Vars) + len(td.Funcs) + len(td.Methods)
		}
		if n > 0 {
			a[w] = td
			w++
		}
	}
	return a[0:w]
}

// Filter eliminates documentation for names that don't pass through the filter f.
// TODO(gri): Recognize "Type.Method" as a name.
//
func (p *Package) Filter(f Filter) {
	p.Consts = filterValues(p.Consts, f)
	p.Vars = filterValues(p.Vars, f)
	p.Types = filterTypes(p.Types, f)
	p.Funcs = filterFuncs(p.Funcs, f)
	p.Doc = "" // don't show top-level package doc
}