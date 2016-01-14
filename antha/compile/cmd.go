// antha/compile/cmd.go: Part of the Antha language
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

// Functionality related to supporting the antha command-line interface
package compile

import (
	"bytes"
	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/token"
	"text/template"
)

// GenerateGraphLib builds a go file defining processes (i.e., go struct
// instances) defined by components
func GenerateComponentLib(b *bytes.Buffer, components []string, workingDirectory string, package_ string) {
	var tmplt = `package {{.Package}}

import (
{{$workingDirectory := .WorkingDirectory}}
{{range .Components}}	"{{$workingDirectory}}/{{.}}"
{{end}})

type ComponentDesc struct {
	Name string
	Constructor func() interface{}
}

func GetComponents() []ComponentDesc {
	c := make([]ComponentDesc, 0)
{{range .Components}}	c = append(c, ComponentDesc{Name: "{{.}}", Constructor: {{.}}.New })
{{end}}
	return c
}`

	var spec struct {
		WorkingDirectory string
		Package          string
		Components       []string
	}
	spec.Package = package_
	spec.WorkingDirectory = workingDirectory
	for _, component := range components {
		spec.Components = append(spec.Components, component)
	}

	t := template.Must(template.New("MainGrapher").Parse(tmplt))
	t.Execute(b, spec)
}

func (cfg *Config) GetComponentName(fset *token.FileSet, node interface{}) string {
	nodeSizes := make(map[ast.Node]int)
	var p compiler
	p.init(cfg, fset, nodeSizes)
	p.printNode(node)
	return p.pkgName
}
