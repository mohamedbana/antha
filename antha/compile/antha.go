// antha/compiler/generator.go: Part of the Antha language
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

// This file implements the generation of go code from the
// antha blocks. It is primarily used in nodes.go and compile.go

package compile

import (
	"bytes"
	"fmt"
	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	"log"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

// Augmentation to compiler state to parse Antha files
type antha struct {
	pkgName      string
	params       []*ast.ValueSpec
	results      []*ast.ValueSpec
	paramTypes   map[string]string               // map used to associate params with types
	paramMap     map[string]string               // map used to sugar antha variables
	resultMap    map[string]string               // map used to sugar antha variables
	reuseMap     map[token.Token]map[string]bool // map use of data through execution phases
	intrinsicMap map[string]string               // replacements for identifiers in expressions in functions
	typeMap      map[string]string               // replacement for type names in type expressions and type and type lists
}

func (p *compiler) anthaInit() {
	p.paramMap = make(map[string]string)
	p.resultMap = make(map[string]string)
	p.paramTypes = make(map[string]string)
	p.reuseMap = make(map[token.Token]map[string]bool)
	p.intrinsicMap = map[string]string{
		"MixInto":  "_wrapper.MixInto",
		"Mix":      "_wrapper.Mix",
		"Incubate": "_wrapper.Incubate",
	}
	p.typeMap = map[string]string{
		"Temperature":       "wunit.Temperature",
		"Time":              "wunit.Time",
		"Length":            "wunit.Length",
		"Area":              "wunit.Area",
		"Volume":            "wunit.Volume",
		"Concentration":     "wunit.Concentration",
		"Amount":            "wunit.Amount",
		"Mass":              "wunit.Mass",
		"Angle":             "wunit.Angle",
		"Energy":            "wunit.Energy",
		"SubstanceQuantity": "wunit.SubstantQuantity",
	}
}

// Merges multiple import blocks and then adds paths
func (p *compiler) addAnthaImports(file *ast.File, paths []string) {
	var specs []ast.Spec
	var decls []ast.Decl
	var pos token.Pos // Dummy position to use for generated imports

	for _, d := range file.Decls {
		gd, ok := d.(*ast.GenDecl)
		if pos == token.NoPos {
			pos = d.Pos()
		}
		if !ok || gd.Tok != token.IMPORT {
			decls = append(decls, d)
			continue
		}
		for _, s := range gd.Specs {
			specs = append(specs, s)
		}
	}

	for _, p := range paths {
		specs = append(specs,
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    fmt.Sprintf(`"%s"`, p),
					ValuePos: pos,
				}})
	}

	if len(specs) == 0 {
		if len(decls) != len(file.Decls) {
			// Clean up empty imports
			file.Decls = decls
		}
		return
	}

	//if pos == token.NoPos {
	//	log.Panicf("no declarations in antha file: %s", file.Name)
	//}

	merged := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: pos,
		Rparen: pos,
		Specs:  specs,
	}

	decls = append([]ast.Decl{merged}, decls...)
	file.Decls = decls
	// NB(ddn): tried to sort here, but the following needs proper token.Pos,
	// which are annoying to generate now. A gofmt on the generated file should
	// be just as good.
	//ast.SortImports(token.NewFileSet(), file)
}

// generate all the On... message handlers for the element
func (p *compiler) genHandlers() {
	var b bytes.Buffer

	numParams := len(p.paramMap)

	// Sort
	sortedParams := make([]string, len(p.paramMap))
	i := 0
	for k, _ := range p.paramMap {
		sortedParams[i] = k
		i += 1
	}
	sort.Strings(sortedParams)
	for _, param := range sortedParams {
		p.genOnFunction(param, numParams, b)
	}
}

// Generate message handler for parameter
func (p *compiler) genOnFunction(name string, paramCount int, b bytes.Buffer) {
	type anthaParam struct {
		PkgName    string
		ParamName  string
		ParamCount int
	}

	const tpl = `func (e *{{.PkgName}}) On{{.ParamName}}(param execute.ThreadParam) {
	e.lock.Lock()
	var bag *execute.AsyncBag = e.params[param.ID]
	if bag == nil {
		bag = new(execute.AsyncBag)
		bag.Init({{.ParamCount}}, e, e)
		e.params[param.ID] = bag
	}
	e.lock.Unlock()

	fired := bag.AddValue("{{.ParamName}}", param)
	if fired {
		e.lock.Lock()
		delete(e.params, param.ID)
		e.lock.Unlock()
	}
}`
	var funcVals anthaParam
	funcVals.PkgName = p.pkgName
	funcVals.ParamName = strings.Title(name)
	funcVals.ParamCount = paramCount

	t := template.Must(template.New("On" + name).Parse(tpl))
	t.Execute(&b, funcVals)
	p.print(b, newline)
}

// Return appropriate go type for an antha (type) expr
func (p *compiler) getTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		if v, ok := p.typeMap[t.Name]; ok {
			res = v
		} else {
			res = t.Name
		}
		return
	case *ast.SelectorExpr:
		res = p.getTypeString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType: // note: array types can use a param as the length, so must be allocated and treated as a slice since they can be dynamic
		res = "[]" + p.getTypeString(t.Elt)
		return
	case *ast.StarExpr:
		res = "*" + p.getTypeString(t.X)
		return
	default:
		log.Panicln("Invalid type spec to get type of: ", reflect.TypeOf(e), t)
	}
	return
}

// Return appropriate go configuration type for an antha (type) expr
func (p *compiler) getConfigTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		if v, ok := p.typeMap[t.Name]; ok {
			res = v
		} else {
			res = t.Name
		}
		return
	case *ast.SelectorExpr:
		res = p.getConfigTypeString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType: // note: array types can use a param as the length, so must be allocated and treated as a slice since they can be dynamic
		res = "[]" + p.getConfigTypeString(t.Elt)
		return
	case *ast.StarExpr:
		res = "wtype.FromFactory"
		return
	default:
		log.Panicln("Invalid type spec to get type of: ", reflect.TypeOf(e), t)
	}
	return
}

// Remove antha param declarations since they are fully captured in indexParams()
func (p *compiler) removeParamDecls(src *ast.File) {
	var newd []ast.Decl
	for _, decl := range src.Decls {
		toadd := true
		switch decl.(type) {
		case *ast.GenDecl:
			switch decl.(*ast.GenDecl).Tok {
			case token.PARAMETERS, token.DATA, token.INPUTS, token.OUTPUTS:
				toadd = false
			}
		}
		if toadd {
			newd = append(newd, decl)
		}
	}
	src.Decls = newd
}

// Modify AST
func (p *compiler) transform(src *ast.File) {
	p.addAnthaImports(src,
		[]string{
			"github.com/antha-lang/antha/antha/anthalib/execution",
			"github.com/antha-lang/antha/antha/anthalib/wunit",
			"github.com/antha-lang/antha/antha/execute",
			"github.com/antha-lang/antha/flow",
			"sync",
			"encoding/json",
		},
	)
	p.sugarAST(src.Decls)
	p.addPreamble(src.Decls)
	p.removeParamDecls(src)
}

// Print out additional go code for each antha file
func (p *compiler) generate() {
	p.genBoilerplate()
	p.genHandlers()
	p.genStructs()
	p.genComponentInfo()
}

// Collect information needed in downstream generation passes
func (p *compiler) analyze(src *ast.File) {
	p.indexParams(src.Decls)
	p.pkgName = strings.Title(src.Name.Name)
}

// index all the spec definitions for inputs and outputs to element
func (p *compiler) indexParams(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.PARAMETERS, token.INPUTS:
				for j := range decl.Specs {
					tmp := decl.Specs[j].(*ast.ValueSpec)
					p.params = append(p.params, tmp)
					for k := range tmp.Names {
						//param := strings.Title(tmp.Names[k].String()) //TODO check input data cannot be lowercase
						param := tmp.Names[k].String()
						p.paramTypes[param] = p.getTypeString(tmp.Type)
						p.paramMap[param] = "p." + param
					}
				}
			case token.OUTPUTS, token.DATA:
				for j := range decl.Specs {
					tmp := decl.Specs[j].(*ast.ValueSpec)
					p.results = append(p.results, tmp)
					for k := range tmp.Names {
						if string(tmp.Names[k].String()[0]) == strings.ToUpper(string(tmp.Names[k].String()[0])) {
							result := tmp.Names[k].String()
							p.paramTypes[result] = p.getTypeString(tmp.Type)
							p.resultMap[result] = "r." + result //assign to resultblock, later staged
						} else { //lowercase, assign to resultblock
							result := tmp.Names[k].String()
							p.paramTypes[result] = p.getTypeString(tmp.Type)
							p.resultMap[result] = "r." + result
						}
					}
				}
			}
		case *ast.AnthaDecl:
			var ref ast.AnthaDecl
			ref = *decl
			tok := ref.Tok
			ast.Inspect(ref.Body, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.Ident:
					if _, is := p.resultMap[n.(*ast.Ident).Name]; is {
						varName := n.(*ast.Ident).Name
						if _, ok := p.reuseMap[tok]; !ok {
							p.reuseMap[tok] = make(map[string]bool)
						}
						p.reuseMap[tok][varName] = true
					}
				}
				return true
			})
		default:
			continue
		}
	}
}

// function to generate boilerplate antha element code
func (p *compiler) genBoilerplate() {
	var boilerPlateTemplate = `// AsyncBag functions
func (e *{{.PkgName}}) Complete(params interface{}) {
	p := params.({{.PkgName}}ParamBlock)
	if p.Error {
{{range .Results}}		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}		return
	}
	r := new({{.PkgName}}ResultBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)
	if r.Error {
{{range .Results}}		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}		return
	}
{{range .StepsOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
	e.analysis(p, r)
		if r.Error {
{{range .Results}}		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}		return
	}
{{range .AnalysisOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
	e.validation(p, r)
	if r.Error {
{{range .Results}}		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}		return
	}
{{range .ValidationOutput}}	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
}

// init function, read characterization info from seperate file to validate ranges?
func (e *{{.PkgName}}) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func (e *{{.PkgName}}) NewConfig() interface{} {
	return &{{.PkgName}}Config{}
}

func (e *{{.PkgName}}) NewParamBlock() interface{} {
	return &{{.PkgName}}ParamBlock{}
}

func New{{.PkgName}}() interface{} {//*{{.PkgName}} {
	e := new({{.PkgName}})
	e.init()
	return e
}
{{$p := .PkgName}}
// Mapper function
func (e *{{.PkgName}}) Map(m map[string]interface{}) interface{} {
	var res {{.PkgName}}ParamBlock
	res.Error = false {{range .Params}} || m["{{.Name}}"].(execute.ThreadParam).Error{{end}}
{{range .Params}}
	v{{.Name}}, is := m["{{.Name}}"].(execute.ThreadParam).Value.(execute.JSONValue)
	if is {
		var temp {{$p}}JSONBlock
		json.Unmarshal([]byte(v{{.Name}}.JSONString), &temp)
		res.{{.Name}} = *temp.{{.Name}}
	} else {
		res.{{.Name}} = m["{{.Name}}"].(execute.ThreadParam).Value.({{.Type}})
	}
{{end}}
{{if .HasIdLine}}
	res.ID = m["{{.IdVarName}}"].(execute.ThreadParam).ID
{{end}}
	return res
}

`
	//res.Error = false {{range .Params}} || m["{{.IdVarName}}"].(execute.ThreadParam).Error {{end}}
	type BoilerPlateParamVars struct {
		Name, Type string
	}
	type OutputRedirection struct {
		OutputVariableName, ResultVariableName string
	}
	type BoilerPlateVars struct {
		PkgName          string
		Params           []BoilerPlateParamVars
		HasIdLine        bool
		IdVarName        string
		StepsOutput      []OutputRedirection
		AnalysisOutput   []OutputRedirection
		ValidationOutput []OutputRedirection
		Results          []OutputRedirection
	}

	bpVars := BoilerPlateVars{PkgName: p.pkgName}
	for i := range p.params {
		bpVars.Params = append(bpVars.Params, BoilerPlateParamVars{p.params[i].Names[0].Name, p.getTypeString(p.params[i].Type)})
	}
	if len(p.params) > 0 {
		bpVars.HasIdLine = true
		bpVars.IdVarName = p.params[0].Names[0].Name
	}
	for i := range p.results {
		varName := p.results[i].Names[0].Name
		if string(varName[0]) == strings.ToUpper(string(varName[0])) { // only upper case are external
			bpVars.Results = append(bpVars.Results, OutputRedirection{varName, varName})
			_, an := p.reuseMap[token.ANALYSIS][varName]
			_, va := p.reuseMap[token.VALIDATION][varName]
			if !(an || va) { //steps is the last one it appeared (if it did)
				bpVars.StepsOutput = append(bpVars.StepsOutput, OutputRedirection{varName, varName})
			} else if !va { //analysis is the last one it appeared (if it did)
				bpVars.AnalysisOutput = append(bpVars.AnalysisOutput, OutputRedirection{varName, varName})
			} else { //validation is the last one it appeared (if it did)
				bpVars.ValidationOutput = append(bpVars.ValidationOutput, OutputRedirection{varName, varName})
			}
		}
	}

	var b bytes.Buffer
	t := template.Must(template.New(p.pkgName + "_boilerplate").Parse(boilerPlateTemplate))
	t.Execute(&b, bpVars)
	p.print(b, newline)
}

// helper function to generate the paramblock and primary element structs
func (p *compiler) genStructs() {
	var structTemplate = `type {{.ElementName}} struct {
	flow.Component                    // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once
	params         map[execute.ThreadID]*execute.AsyncBag{{range .Inports}}
	{{.Name}}		<-chan execute.ThreadParam{{end}}{{range .Outports}}
	{{.Name}}		chan<- execute.ThreadParam{{end}}
}

type {{.ElementName}}ParamBlock struct{
	ID		execute.ThreadID
	Error	bool{{range .Params}}
	{{.Name}}		{{.Type}}{{end}}
}

type {{.ElementName}}Config struct{
	ID		execute.ThreadID
	Error	bool{{range .Configs}}
	{{.Name}}		{{.Type}}{{end}}
}

type {{.ElementName}}ResultBlock struct{
	ID		execute.ThreadID
	Error	bool{{range .Results}}
	{{.Name}}		{{.Type}}{{end}}
}

type {{.ElementName}}JSONBlock struct{
	ID			*execute.ThreadID
	Error		*bool{{range .Params}}
	{{.Name}}		*{{.Type}}{{end}}{{range .Results}}
	{{.Name}}		*{{.Type}}{{end}}
}`

	type VarSpec struct {
		Name, Type string
	}
	type StructSpec struct {
		ElementName              string
		Params, Inports, Configs []VarSpec
		Results, Outports        []VarSpec
	}

	t := template.Must(template.New(p.pkgName + "_struct").Parse(structTemplate))

	var structSpec StructSpec
	structSpec.ElementName = p.pkgName
	for _, i := range p.params {
		n := i.Names[0].Name
		ts := p.getTypeString(i.Type)
		cs := p.getConfigTypeString(i.Type)
		structSpec.Params = append(structSpec.Params, VarSpec{n, ts})
		structSpec.Inports = append(structSpec.Inports, VarSpec{n, ts})
		structSpec.Configs = append(structSpec.Configs, VarSpec{n, cs})
	}
	for _, i := range p.results {
		structSpec.Results = append(structSpec.Results, VarSpec{i.Names[0].Name, p.getTypeString(i.Type)})
		if string(i.Names[0].Name[0]) == strings.ToUpper(string(i.Names[0].Name[0])) {
			structSpec.Outports = append(structSpec.Outports, VarSpec{i.Names[0].Name, p.getTypeString(i.Type)})
		}
	}
	var b bytes.Buffer
	t.Execute(&b, structSpec)
	p.print(b, newline)
}

func (p *compiler) genComponentInfo() {
	var tmp = `
func (c *{{.ElementName}}) ComponentInfo() *execute.ComponentInfo {
	inp := make([]execute.PortInfo,0)
	outp := make([]execute.PortInfo,0)
{{range .Inports}}	inp = append(inp, *execute.NewPortInfo( "{{.Name}}", "{{.Type}}", "{{.Name}}", true, true, nil, nil))
{{end}}{{range .Outports}}	outp = append(outp, *execute.NewPortInfo( "{{.Name}}", "{{.Type}}", "{{.Name}}", true, true, nil, nil))
{{end}}
	ci := execute.NewComponentInfo("{{.ElementName}}", "{{.ElementName}}", "", false, inp, outp)

	return ci
}`

	type VarSpec struct {
		Name, Type string
	}
	type StructSpec struct {
		ElementName       string
		Inports, Outports []VarSpec
	}

	t := template.Must(template.New(p.pkgName + "_jsonDesc").Parse(tmp))

	var structSpec StructSpec
	structSpec.ElementName = p.pkgName
	for _, i := range p.params {
		structSpec.Inports = append(structSpec.Inports, VarSpec{i.Names[0].Name, p.getTypeString(i.Type)})
	}
	for _, i := range p.results {
		if string(i.Names[0].Name[0]) == strings.ToUpper(string(i.Names[0].Name[0])) {
			structSpec.Outports = append(structSpec.Outports, VarSpec{i.Names[0].Name, p.getTypeString(i.Type)})
		}
	}
	var b bytes.Buffer
	t.Execute(&b, structSpec)
	p.print(b, newline)

}

// Add statements to the beginning of antha decls
func (p *compiler) addPreamble(decls []ast.Decl) {
	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.AnthaDecl:
			if d.Tok == token.REQUIREMENTS {
				//_ = wunit.Make_units
				s2Lhs, _ := parser.ParseExpr("_")
				s2Rhs, _ := parser.ParseExpr("wunit.Make_units")
				s2 := &ast.AssignStmt{
					Lhs:    []ast.Expr{s2Lhs},
					Rhs:    []ast.Expr{s2Rhs},
					TokPos: d.TokPos,
					Tok:    token.ASSIGN,
				}
				d.Body.List = append([]ast.Stmt{s2}, d.Body.List...)

				continue
			}
			// `_wrapper := execution.NewWrapper()`
			s1Lhs, _ := parser.ParseExpr("_wrapper")
			s1Rhs, _ := parser.ParseExpr("execution.NewWrapper(p.ID)")
			s1 := &ast.AssignStmt{
				Lhs:    []ast.Expr{s1Lhs},
				Rhs:    []ast.Expr{s1Rhs},
				TokPos: d.TokPos,
				Tok:    token.DEFINE,
			}
			// `_ = _wrapper`
			s2Lhs, _ := parser.ParseExpr("_")
			s2Rhs, _ := parser.ParseExpr("_wrapper")
			s2 := &ast.AssignStmt{
				Lhs:    []ast.Expr{s2Lhs},
				Rhs:    []ast.Expr{s2Rhs},
				TokPos: d.TokPos,
				Tok:    token.ASSIGN,
			}
			d.Body.List = append([]ast.Stmt{s1, s2}, d.Body.List...)
		default:
			continue
		}
	}
}

// Update AST for antha semantics
func (p *compiler) sugarAST(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.GenDecl:
		case *ast.AnthaDecl:
			p.sugarForParams(decl.Body)
			p.sugarForIntrinsics(decl.Body)
		default: // ignore all other decls... for now
			continue
		}
	}
}

// Replace bare antha identifiers with go qualified names
func (p *compiler) sugarForParams(body *ast.BlockStmt) {
	list := body.List

	for i := range list {
		ast.Inspect(list[i], func(n ast.Node) bool {
			switch node := n.(type) {
			case nil:
				return false
			case *ast.Ident:
				// sugar the name if it is a known param
				if sugar, ok := p.paramMap[node.Name]; ok {
					node.Name = sugar
				} else if sugar, ok := p.resultMap[node.Name]; ok {
					node.Name = sugar
				}
				// test if the left hand side of an assignment is to a results variable
				// if so, sugar it by replacing the = with a channel operation
			case *ast.AssignStmt:
				// test if a result is being assigned to.
				fix_assign := false
				for j := range node.Lhs {
					if ident, ok := node.Lhs[j].(*ast.Ident); ok {
						isUpper := string(ident.Name[0]) == strings.ToUpper(string(ident.Name[0]))
						if sugar, fixme := p.resultMap[ident.Name]; fixme && isUpper {
							//fix_assign = true
							ident.Name = sugar
						}
					}
				}
				//TODO: generate multiple stmts if needed
				if fix_assign {
					if len(node.Rhs) > 1 {
						// currently undefined so panic
						panic("too many right hand exprs")
					}
					list[i] = &ast.SendStmt{node.Lhs[0], node.TokPos, p.wrapParamExpr(node.Rhs[0])}
				}
			default:

			}
			return true
		})
	}
}

// Replace bare antha function names with go qualified names
func (p *compiler) sugarForIntrinsics(root ast.Node) {
	ast.Inspect(root, func(n ast.Node) bool {
		switch node := n.(type) {
		case nil:
			return false
		case *ast.Ident:
			// sugar the name if it is a known param
			if sugar, ok := p.intrinsicMap[node.Name]; ok {
				node.Name = sugar
			}
			// test if the left hand side of an assignment is to a results variable
			// if so, sugar it by replacing the = with a channel operation
		default:
		}
		return true
	})
}

// Return qualified reference to antha parameter
func (p *compiler) wrapParamExpr(e ast.Expr) (res ast.Expr) {
	type_ := &ast.SelectorExpr{ast.NewIdent("execute"), ast.NewIdent("ThreadParam")}
	exp := &ast.SelectorExpr{ast.NewIdent("p"), ast.NewIdent("ID")}
	var elts []ast.Expr
	elts = append(elts, e)
	elts = append(elts, exp)
	res = &ast.CompositeLit{type_, token.NoPos, elts, token.NoPos}
	return
}
