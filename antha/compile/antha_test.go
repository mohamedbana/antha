// /compile/generator_test.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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

package compile

import (
	"bytes"
	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	"testing"
)

func TestTypeSugaring(t *testing.T) {
	nodeSizes := make(map[ast.Node]int)
	cfg := &Config{}
	compiler := &compiler{}
	fset := token.NewFileSet()
	compiler.init(cfg, fset, nodeSizes)

	expr, err := parser.ParseExpr("func(x Volume) Concentration { x := Volume }")
	if err != nil {
		t.Fatal(err)
	}
	desired, err := parser.ParseExpr("func(x wunit.Volume) wunit.Concentration { x := Volume }")
	if err != nil {
		t.Fatal(err)
	}

	compiler.sugarForTypes(expr)
	var buf1, buf2 bytes.Buffer
	if err := compiler.Fprint(&buf1, fset, expr); err != nil {
		t.Fatal(err)
	}
	if err := compiler.Fprint(&buf2, fset, desired); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Errorf("wanted\n'''%s'''\ngot\n'''%s'''\n", buf2.String(), buf1.String())
	}
}

// Verify that generation of an empty component produces a valid go file
func TestGenerateGraphRunnerOfEmptyComponent(t *testing.T) {
	comp := execute.ComponentInfo{Name: "Test", Description: "", Icon: "", Subgraph: false, InPorts: nil, OutPorts: nil}
	componentLibrary := []execute.ComponentInfo{comp}
	var buf bytes.Buffer
	GenerateGraphRunner(&buf, componentLibrary, "")
	_, err := parser.ParseFile(fset, "", buf.String(), parser.AllErrors)
	if err != nil {
		t.Error(err)
	}
}

// Verify that parser does fail sometimes
func TestBadGenerateGraphRunnerOfEmptyComponent(t *testing.T) {
	comp := execute.ComponentInfo{Name: "Test", Description: "", Icon: "", Subgraph: false, InPorts: nil, OutPorts: nil}
	componentLibrary := []execute.ComponentInfo{comp}
	var buf bytes.Buffer
	GenerateGraphRunner(&buf, componentLibrary, "")
	buf.WriteString("invalid tokens at end of program")
	_, err := parser.ParseFile(fset, "", buf.String(), parser.AllErrors)
	if err == nil {
		t.Error("expected illegal program")
	}
}

// Verify that generation of an empty component produces a valid go file
func TestGenerateLibOfEmptyComponent(t *testing.T) {
	comp := execute.ComponentInfo{Name: "Test", Description: "", Icon: "", Subgraph: false, InPorts: nil, OutPorts: nil}
	componentLibrary := []execute.ComponentInfo{comp}
	var buf bytes.Buffer
	GenerateComponentLib(&buf, componentLibrary, "", "main")
	_, err := parser.ParseFile(fset, "", buf.String(), parser.AllErrors)
	if err != nil {
		t.Error(err)
	}
}

// Verify that parser does fail sometimes
func TestBadGenerateLibOfEmptyComponent(t *testing.T) {
	comp := execute.ComponentInfo{Name: "Test", Description: "", Icon: "", Subgraph: false, InPorts: nil, OutPorts: nil}
	componentLibrary := []execute.ComponentInfo{comp}
	var buf bytes.Buffer
	GenerateComponentLib(&buf, componentLibrary, "", "main")
	buf.WriteString("invalid tokens at end of program")
	_, err := parser.ParseFile(fset, "", buf.String(), parser.AllErrors)
	if err == nil {
		t.Error("expected illegal program")
	}
}
