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
	"fmt"
	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/antha/token"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

// helper function to generate a wrapper main file to convert an element into a standalone binary
// Prerequisite: Antha elements have been indexed in this compiler context (such as by calling this
// after indexParams and sugarAST
func (p *compiler) standAlone(src *ast.File, packageRoute string) []byte {

	var b bytes.Buffer

	p.getMainPackageImports(&b, packageRoute)
	p.getMainApp(&b)
	p.getMainReferenceFunctions(&b)

	return b.Bytes()
}

func (p *compiler) getMainPackageImports(b *bytes.Buffer, packageRoute string) {
	type MainImportsSpec struct {
		LowerPkgName string
		PkgRoute     string
	}

	var element = `package main

import "github.com/antha-lang/antha/antha/execute"
import "github.com/antha-lang/antha/flow"
import "os"
import "io"
import "encoding/json"
import "log"
import {{.LowerPkgName}} "{{.PkgRoute}}"

var (
	exitCode = 0
)

`

	var funcVals MainImportsSpec
	funcVals.LowerPkgName = string(unicode.ToLower(rune(p.pkgName[0]))) + p.pkgName[1:]
	funcVals.PkgRoute = packageRoute

	t := template.Must(template.New("MainImports").Parse(element))
	t.Execute(b, funcVals)

	return
}

//getMainApp puts in b the contents of a go file capable of runing the component described in the
// compiler. This executable file will read its input from stdin and write it to stdout
func (p *compiler) getMainApp(b *bytes.Buffer) {
	type InPortSpec struct {
		InPortName string
		PkgName    string
	}
	type OutPortSpec struct {
		OutPortName string
		PkgName     string
	}
	type AppSpec struct {
		LowerPkgName string
		PkgName      string
		InPorts      []InPortSpec
		OutPorts     []OutPortSpec
	}

	app := `
type App struct {
    flow.Graph
}

func NewApp() *App {
    n := new(App)
    n.InitGraphState()

    n.Add({{.LowerPkgName}}.New{{.PkgName}}(), "{{.PkgName}}")

{{range .InPorts}}	n.MapInPort("{{.InPortName}}", "{{.PkgName}}", "{{.InPortName}}")
{{end}}
{{range .OutPorts}}	n.MapOutPort("{{.OutPortName}}", "{{.PkgName}}", "{{.OutPortName}}")
{{end}}

   return n
}
`

	t := template.Must(template.New("MainApp").Parse(app))
	var appVals AppSpec
	appVals.PkgName = p.pkgName
	appVals.LowerPkgName = string(unicode.ToLower(rune(p.pkgName[0]))) + p.pkgName[1:]
	for _, name := range p.params {
		appVals.InPorts = append(appVals.InPorts, InPortSpec{name, p.pkgName})
	}
	for _, name := range p.results {
		if name[0:1] == strings.ToUpper(name[0:1]) {
			appVals.OutPorts = append(appVals.OutPorts, OutPortSpec{name, p.pkgName})
		}
	}
	t.Execute(b, appVals)

	return
}

//getMainReferenceFunctions puts in b the necessary funcitons for a main function to be able to read/write and execute an element
func (p *compiler) getMainReferenceFunctions(b *bytes.Buffer) {
	type InPortSpec struct {
		InPortName string
	}
	type OutPortSpec struct {
		OutPortName string
	}
	type ReferenceMainSpec struct {
		LowerPkgName string
		PkgName      string
		InPorts      []InPortSpec
		OutPorts     []OutPortSpec
	}

	referenceMainStart := `
func referenceMain() {
    net := NewApp()

{{range .InPorts}}	{{.InPortName}}Chan := make(chan execute.ThreadParam)
    net.SetInPort("{{.InPortName}}", {{.InPortName}}Chan)
{{end}}

{{range .OutPorts}}	{{.OutPortName}}Chan := make(chan execute.ThreadParam)
    net.SetOutPort("{{.OutPortName}}", {{.OutPortName}}Chan)
{{end}}

    flow.RunNet(net)

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	log.SetOutput(os.Stderr)

	go func() {
{{range .InPorts}}		defer close({{.InPortName}}Chan)
{{end}}

		for {
			var p {{.LowerPkgName}}.{{.PkgName}}JSONBlock
			if err := dec.Decode(&p); err != nil {
				if err != io.EOF {
					log.Println("Error decoding", err)
				}
				return
			}
			//log.Print(p)
			if p.ID == nil {
				log.Println("Error, no ID")
				continue
			}
			if p.Error == nil {
				log.Println("Error, no error")
				continue
			}
{{range .InPorts}}			if p.{{.InPortName}} != nil {
				param := execute.ThreadParam{Value: *(p.{{.InPortName}}), ID: *(p.ID), Error: *(p.Error)}
				{{.InPortName}}Chan <- param
			}
{{end}}
		}
	}()

{{range .OutPorts}}	go func() {
		for sequence := range {{.OutPortName}}Chan {
			if err := enc.Encode(&sequence); err != nil {
				log.Println(err)
			}
		}
	}()
{{end}}

	<-net.Wait()
}

func main() {
	referenceMain()
	os.Exit(exitCode)
}
`

	t := template.Must(template.New("MainReferenceMainStart").Parse(referenceMainStart))
	var rMSpec ReferenceMainSpec
	rMSpec.LowerPkgName = string(unicode.ToLower(rune(p.pkgName[0]))) + p.pkgName[1:]
	rMSpec.PkgName = p.pkgName
	for _, name := range p.params {
		rMSpec.InPorts = append(rMSpec.InPorts, InPortSpec{name})
	}
	for _, name := range p.results {
		if name[0:1] == strings.ToUpper(name[0:1]) {
			rMSpec.OutPorts = append(rMSpec.OutPorts, OutPortSpec{name})
		}
	}
	t.Execute(b, rMSpec)

	return
}

type componentSpec struct {
	Name            string
	ConstructorFunc string
	InPorts         []string
	OutPorts        []string
}

type graphRunnerSpec struct {
	Components       []componentSpec
	WorkingDirectory string
	Package          string
}

type componentSpecs []componentSpec

func (a componentSpecs) Len() int {
	return len(a)
}

func (a componentSpecs) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (a componentSpecs) Swap(i, j int) {
	a[j], a[i] = a[i], a[j]
}

func (a *graphRunnerSpec) Sort() {
	sort.Sort(componentSpecs(a.Components))
	for _, c := range a.Components {
		sort.Strings(c.InPorts)
		sort.Strings(c.OutPorts)
	}
}

// GenerateGraphLib builds a go file defining processes (i.e., go struct
// instances) defined by components
func GenerateComponentLib(b *bytes.Buffer, components []execute.ComponentInfo, workingDirectory string, package_ string) {
	var tmplt = `package {{.Package}}

import (
{{$workingDirectory := .WorkingDirectory}}
{{range .Components}}	"{{$workingDirectory}}/{{.Name}}"
{{end}})

type ComponentDesc struct {
	Name string
	Constructor func() interface{}
}

func GetComponents() []ComponentDesc {
	portMap := make(map[string]map[string]bool) //representing component, port name, and true if in
{{range .Components}}	portMap["{{.Name}}"] = make(map[string]bool){{$componentName := .Name}}
{{range .InPorts}}	portMap["{{$componentName}}"]["{{.}}"] = true
{{end}}
{{range .OutPorts}}	portMap["{{$componentName}}"]["{{.}}"] = false
{{end}}
{{end}}
	c := make([]ComponentDesc, 0)
{{range .Components}}	c = append(c, ComponentDesc{Name: "{{.Name}}", Constructor: {{.Name}}.{{.ConstructorFunc}} })
{{end}}
	return c
}`

	var spec graphRunnerSpec
	spec.Package = package_
	spec.WorkingDirectory = workingDirectory
	for _, component := range components {
		var cs componentSpec
		cs.Name = component.Name
		cs.ConstructorFunc = fmt.Sprintf("New%s", component.Name)
		for _, ip := range component.InPorts {
			cs.InPorts = append(cs.InPorts, ip.Id)
		}
		for _, ip := range component.OutPorts {
			cs.OutPorts = append(cs.OutPorts, ip.Id)
		}
		spec.Components = append(spec.Components, cs)
	}

	spec.Sort()

	t := template.Must(template.New("MainGrapher").Parse(tmplt))
	t.Execute(b, spec)
}

// GenerateGraphRunner builds a go file capable of running FBP (Flow-based
// programming) graphs with processes defined by the components included in the
// components argument.
func GenerateGraphRunner(b *bytes.Buffer, components []execute.ComponentInfo, workingDirectory string) {
	var tmplt = `package main

import (
    "flag"
    "fmt"
	"encoding/json"
	"log"
	"os"

    "github.com/antha-lang/antha/antha/execute"
    "github.com/antha-lang/antha/flow"
{{$workingDirectory := .WorkingDirectory}}
{{range .Components}}	"{{$workingDirectory}}/{{.Name}}"
{{end}}
)
func main() {
    flag.Parse()
    if flag.NArg() == 0 {
        fmt.Println("Graph definition file missing")
        return
    }

	portMap := make(map[string]map[string]bool) //representing component, port name, and true if in
{{range .Components}}	portMap["{{.Name}}"] = make(map[string]bool){{$componentName := .Name}}
{{range .InPorts}}	portMap["{{$componentName}}"]["{{$componentName}}"] = true
{{end}}
{{range .OutPorts}}	portMap["{{$componentName}}"]["{{$componentName}}"] = false
{{end}}
{{end}}

{{range .Components}}	flow.Register("{{.Name}}", {{.Name}}.{{.ConstructorFunc}})
{{end}}

    dec := json.NewDecoder(os.Stdin)
    enc := json.NewEncoder(os.Stdout)
    log.SetOutput(os.Stderr)


    graph := flow.LoadJSON(flag.Arg(0))
    if graph == nil {
        fmt.Println("empty graph")
        return
    }
    flow.RunNet(graph)

    //<-graph.Ready()

	for _, port := range graph.GetUnboundOutPorts() {
		ch := make(chan execute.ThreadParam)
		graph.SetOutPort(port.Port, ch)
		go func() {
			for a := range ch {
				if err := enc.Encode(&a); err != nil {
					log.Println(err)
				}
			}
		}()
	}

	inPortMap := make(map[string]chan execute.ThreadParam)
	for _, port := range graph.GetUnboundInPorts() {
		ch := make(chan execute.ThreadParam)
		inPortMap[port.Port] = ch
		graph.SetInPort(port.Port, ch)
	}
	go func() {
		for _, ch := range inPortMap {
			defer close(ch)
		}

		for {
			var p execute.JSONBlock
			if err := dec.Decode(&p); err != nil {
				log.Println("Error decoding", err)
				return
			}
			if p.ID == nil { //TODO add error control in JSONBlock unmarshaling??
				log.Println("Error, no ID")
				continue
			}
			for k, v := range p.Values {
				tmp := make(map[string]interface{})
				tmp[k] = v
				sthg, err := json.Marshal(&tmp)
				if err != nil {
					continue
				}
				if _, exists := inPortMap[k]; exists {
					param := execute.ThreadParam{Value: execute.JSONValue{Name: k, JSONString: string(sthg)}, ID: *p.ID, Error: *p.Error}
					inPortMap[k] <- param
				}
			}
		}
	}()

	<-graph.Wait()
}`

	var spec graphRunnerSpec
	spec.WorkingDirectory = workingDirectory
	for _, component := range components {
		var cs componentSpec
		cs.Name = component.Name
		cs.ConstructorFunc = fmt.Sprintf("New%s", component.Name)
		for _, ip := range component.InPorts {
			cs.InPorts = append(cs.InPorts, ip.Id)
		}
		for _, ip := range component.OutPorts {
			cs.OutPorts = append(cs.OutPorts, ip.Id)
		}
		spec.Components = append(spec.Components, cs)
	}

	spec.Sort()

	t := template.Must(template.New("MainGrapher").Parse(tmplt))
	t.Execute(b, spec)
}

func (cfg *Config) GetFileComponentInfo(fset *token.FileSet, node interface{}) execute.ComponentInfo {
	nodeSizes := make(map[ast.Node]int)
	var p compiler
	p.init(cfg, fset, nodeSizes)
	p.printNode(node)

	var ci execute.ComponentInfo
	ci.InPorts = make([]execute.PortInfo, 0)
	ci.OutPorts = make([]execute.PortInfo, 0)
	ci.Subgraph = false //TODO this might not be the case sometimes
	ci.Name = p.pkgName
	for name, type_ := range p.paramTypes {
		ci.InPorts = append(ci.InPorts, execute.PortInfo{
			Id:   name,
			Type: type_,
		})
	}
	for name, type_ := range p.resultTypes {
		if name[0:1] == strings.ToUpper(name[0:1]) {
			ci.OutPorts = append(ci.OutPorts, execute.PortInfo{
				Id:   name,
				Type: type_,
			})
		}
	}

	return ci
}
