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
// 1 Royal College St, London NW1 0NH UK

package compile

import (
	"bytes"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/antha/parser"
	"testing"
)

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
