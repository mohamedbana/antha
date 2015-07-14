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
	"log"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/execute"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"

	"fmt"
)

// general antha node print function
func (p *compiler) anthaDecl(d *ast.AnthaDecl) {
	switch d.Tok {
	case token.STEPS, token.SETUP, token.REQUIREMENTS, token.ANALYSIS, token.VALIDATION:
		ctx := new(anthaContext)
		ctx.init(p.pkgName, d.Tok)
		p.anthaSig(d, ctx)
	default:
		panic("Bad anthaDecl")

	}
}

// placeholder function
func (p *compiler) anthaParamDecl(d *ast.GenDecl) {
	p.setComment(d.Doc)
	p.print(d.Pos(), d.Tok, blank)

	if d.Lparen.IsValid() {
		// group of parenthesized declarations
		p.print(d.Lparen, token.LPAREN)
		if n := len(d.Specs); n > 0 {
			p.print(indent, formfeed)
			if n > 1 && (d.Tok == token.CONST || d.Tok == token.VAR) {
				// two or more grouped const/var declarations:
				// determine if the type column must be kept
				keepType := keepTypeColumn(d.Specs)
				var line int
				for i, s := range d.Specs {
					if i > 0 {
						p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) > 0)
					}
					p.recordLine(&line)
					p.valueSpec(s.(*ast.ValueSpec), keepType[i])
				}
			} else {
				var line int
				for i, s := range d.Specs {
					if i > 0 {
						p.linebreak(p.lineFor(s.Pos()), 1, ignore, p.linesFrom(line) > 0)
					}
					p.recordLine(&line)
					p.spec(s, n, false)
				}
			}
			p.print(unindent, formfeed)
		}
		p.print(d.Rparen, token.RPAREN)

	} else {
		// single declaration
		// TODO: should this be illegal syntax?
		p.spec(d.Specs[0], 1, true)
	}
}

// Outputs the Antha block signatures
func (p *compiler) anthaSig(d *ast.AnthaDecl, c *anthaContext) {
	p.setComment(d.Doc)

	p.print(d.Pos(), c.getSignature())

	// adjust the padding to the body
	p.adjBlock(p.distanceFrom(d.Pos()), vtab, d.Body)
}

// antha generator context structure
// Note: this code is not thread safe
type anthaContext struct {
	PkgName   string
	emitStart bool
	emitEnd   bool
	tok       token.Token
}

// Special case blocks
//	PARAMETERS
//	DATA
//	INPUTS
//	OUTPUTS

var anthaSigs = map[token.Token]string{
	token.STEPS:        `func (e *{{.PkgName}}) steps(p {{.PkgName}}ParamBlock, r *{{.PkgName}}ResultBlock)`,
	token.REQUIREMENTS: `func (e *{{.PkgName}}) requirements()`,
	token.SETUP:        `func (e *{{.PkgName}}) setup(p {{.PkgName}}ParamBlock)`,
	token.ANALYSIS:     `func (e *{{.PkgName}}) analysis(p {{.PkgName}}ParamBlock, r *{{.PkgName}}ResultBlock)`,
	token.VALIDATION:   `func (e *{{.PkgName}}) validation(p {{.PkgName}}ParamBlock, r *{{.PkgName}}ResultBlock)`,
}

// templates for any start lines in block
var anthaStarts = map[token.Token]string{
	token.STEPS:        ``,
	token.REQUIREMENTS: ``,
	token.SETUP:        ``,
	token.ANALYSIS:     ``,
	token.VALIDATION:   ``,
}

// templates for any end lines in block
var anthaEnds = map[token.Token]string{
	token.STEPS:        ``,
	token.REQUIREMENTS: ``,
	token.SETUP:        ``,
	token.ANALYSIS:     ``,
	token.VALIDATION:   ``,
}

// init the context to the appropriate type of antha block
func (c *anthaContext) init(name string, tok token.Token) {
	c.emitStart = true
	c.emitEnd = true
	c.tok = tok
	c.PkgName = name
}

// simple helper function to generate the appropriate block
// function signatures
func (c *anthaContext) getSignature() (b bytes.Buffer) {
	t := template.Must(template.New(c.tok.String() + "_sig").Parse(anthaSigs[c.tok]))
	t.Execute(&b, c)
	return
}

// simple helper function to generate the beginning of an antha block
// depending on the block type
func (c *anthaContext) beginBlock() (b bytes.Buffer, empty bool) {
	if !c.emitStart {
		empty = true
		return
	}
	t := template.Must(template.New(c.tok.String() + "_start").Parse(anthaStarts[c.tok]))
	t.Execute(&b, c)
	empty = b.Len() > 0
	c.emitStart = false
	return
}

// simple helper function to generate the ending of an antha block
// depending on the block type
func (c *anthaContext) endBlock() (b bytes.Buffer, empty bool) {
	if !c.emitEnd {
		empty = true
		return
	}
	t := template.Must(template.New(c.tok.String() + "_end").Parse(anthaEnds[c.tok]))
	t.Execute(&b, c)
	empty = b.Len() > 0
	c.emitEnd = false
	return
}

func (p *compiler) addAnthaImports() {
	const tpl = `

import "github.com/antha-lang/antha/antha/anthalib/wunit"
import "github.com/antha-lang/antha/antha/anthalib/execution"
import "github.com/antha-lang/antha/antha/execute"
import "github.com/antha-lang/antha/flow"
import "sync"
import "encoding/json"
`
	var b bytes.Buffer
	t := template.Must(template.New("imports").Parse(tpl))
	t.Execute(&b, nil)
	p.print(b, newline)
}

// finds all the importspec decls, and adds the antha decls to the last one
/*func (p *compiler) addAnthaImports(list []ast.Decl) {

	for _, d := range list {
		d, ok := d.(*GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			break
		}

		if !d.Lparen.IsValid() {
			// Not a block: sorted by default.
			continue
		}
	}
}*/

// utility function to generate all the On... message handlers for the element
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
		p.getOnFunction(param, numParams, b)
	}
}

// helper function to generate a parameter or inputs signatures
func (p *compiler) getOnFunction(name string, paramCount int, b bytes.Buffer) {
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
	return
}

// utility function to get type strings from an expr
func getTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		res = getAnthaLibReference(t.Name)
		return
	case *ast.SelectorExpr:
		res = getTypeString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType: // note: array types can use a param as the length, so must be allocated and treated as a slice since they can be dynamic
		res = "[]" + getTypeString(t.Elt)
		return
	case *ast.StarExpr:
		res = "*" + getTypeString(t.X)
		return
	default:
		log.Panicln("Invalid type spec to get type of: ", reflect.TypeOf(e), t)
	}
	return
}

func getAnthaLibReference(name string) string {
	switch name{
	case "Temperature":
		return "wunit."+name
	case "Time":
		return "wunit."+name
	case "Length":
		return "wunit."+name
	case "Area":
		return "wunit."+name
	case "Volume":
		return "wunit."+name
	case "Amount":
		return "wunit."+name
	case "Mass":
		return "wunit."+name
	case "Angle":
		return "wunit."+name
	case "Energy":
		return "wunit."+name
	case "SubstanceQuantity":
		return "wunit."+name
	}
	return name
}

func getConfigString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		res = getAnthaLibReference(t.Name)
		return
	case *ast.SelectorExpr:
		res = getConfigString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType: // note: array types can use a param as the length, so must be allocated and treated as a slice since they can be dynamic
		res = "[]" + getConfigString(t.Elt)
		return
	case *ast.StarExpr:
		res = "wtype.FromFactory"
		return
	default:
		log.Panicln("Invalid type spec to get type of: ", reflect.TypeOf(e), t)
	}
	return
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
						p.paramTypes[param] = getTypeString(tmp.Type)
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
							p.paramTypes[result] = getTypeString(tmp.Type)
							p.resultMap[result] = "r." + result //assign to resultblock, later staged
						} else { //lowercase, assign to resultblock
							result := tmp.Names[k].String()
							p.paramTypes[result] = getTypeString(tmp.Type)
							p.resultMap[result] = "r." + result
						}
					}
				}
			}

		case *ast.AnthaDecl:
			var ref ast.AnthaDecl
			ref = *decl
			for _, item := range ref.Body.List {
				p.analyseStmtVariableUse(ref.Tok, item)
			}
		default:
			continue
		}
	}
}

func (p *compiler) analyseStmtVariableUse(tok token.Token, n interface{}) {
	switch n.(type) {
	case *ast.ExprStmt:
		p.analyseExprVariableUse(tok, n.(*ast.ExprStmt).X)
	case *ast.AssignStmt:
		for _, i := range n.(*ast.AssignStmt).Lhs {
			p.analyseExprVariableUse(tok, i)
		}
		for _, i := range n.(*ast.AssignStmt).Rhs {
			p.analyseExprVariableUse(tok, i)
		}
	case *ast.SendStmt:
		p.analyseExprVariableUse(tok, n.(*ast.SendStmt).Chan)
		p.analyseExprVariableUse(tok, n.(*ast.SendStmt).Value)
	case *ast.DeclStmt:
	case *ast.IfStmt:
		p.analyseStmtVariableUse(tok, n.(*ast.IfStmt).Init)
		p.analyseExprVariableUse(tok, n.(*ast.IfStmt).Cond)
		p.analyseStmtVariableUse(tok, n.(*ast.IfStmt).Else)
		for i := range n.(*ast.IfStmt).Body.List {
			p.analyseStmtVariableUse(tok, i)
		}
	default:
	}
}

func (p *compiler) analyseExprVariableUse(tok token.Token, n interface{}) {
	switch n.(type) {
	case *ast.BadExpr:
		//fmt.Println("BadExpr") //Nothing to do
	case *ast.Ident:
		//fmt.Println("Ident")
		if _, is := p.resultMap[n.(*ast.Ident).Name]; is {
			p.addVariableUse(tok, n.(*ast.Ident).Name)
		}

	case *ast.Ellipsis:
		//fmt.Println("Ellipsis") //Not likely to happen if no funcs are defined
		p.analyseExprVariableUse(tok, n.(*ast.Ellipsis).Elt)
	case *ast.BasicLit:
		//fmt.Println("BasicLit") //Nothing to do
	case *ast.FuncLit:
		//fmt.Println("FuncLit")
		p.analyseStmtVariableUse(tok, n.(*ast.FuncLit).Body)
	case *ast.CompositeLit:
		//fmt.Println("CompositeLit")
		p.analyseExprVariableUse(tok, n.(*ast.CompositeLit).Type)
		for _, i := range n.(*ast.CompositeLit).Elts {
			p.analyseExprVariableUse(tok, i)
		}
	case *ast.ParenExpr:
		//fmt.Println("ParenExpr")
		p.analyseExprVariableUse(tok, n.(*ast.ParenExpr).X)
	case *ast.SelectorExpr:
		//fmt.Println("SelectorExpr")
		p.analyseExprVariableUse(tok, n.(*ast.SelectorExpr).X)
	case *ast.IndexExpr:
		//fmt.Println("IndexExpr")
		p.analyseExprVariableUse(tok, n.(*ast.IndexExpr).X)
		p.analyseExprVariableUse(tok, n.(*ast.IndexExpr).Index)
	case *ast.SliceExpr:
		//fmt.Println("SliceExpr")
		p.analyseExprVariableUse(tok, n.(*ast.SliceExpr).X)
		p.analyseExprVariableUse(tok, n.(*ast.SliceExpr).Low)
		p.analyseExprVariableUse(tok, n.(*ast.SliceExpr).High)
		p.analyseExprVariableUse(tok, n.(*ast.SliceExpr).Max)
	case *ast.TypeAssertExpr:
		//fmt.Println("TypeAssertExpr")
		p.analyseExprVariableUse(tok, n.(*ast.TypeAssertExpr).X)
		p.analyseExprVariableUse(tok, n.(*ast.TypeAssertExpr).Type)
	case *ast.CallExpr:
		//fmt.Println("CallExpr")
		p.analyseExprVariableUse(tok, n.(*ast.CallExpr).Fun)
		for _, i := range n.(*ast.CallExpr).Args {
			p.analyseExprVariableUse(tok, i)
		}
	case *ast.StarExpr:
		//fmt.Println("StarExpr")
		p.analyseExprVariableUse(tok, n.(*ast.StarExpr).X)
	case *ast.UnaryExpr:
		//fmt.Println("UnaryExpr")
		p.analyseExprVariableUse(tok, n.(*ast.UnaryExpr).X)
	case *ast.BinaryExpr:
		//fmt.Println("BinaryExpr")
		p.analyseExprVariableUse(tok, n.(*ast.BinaryExpr).X)
		p.analyseExprVariableUse(tok, n.(*ast.BinaryExpr).Y)
	case *ast.KeyValueExpr:
		//fmt.Println("KeyValueExpr")
		p.analyseExprVariableUse(tok, n.(*ast.KeyValueExpr).Key)
		p.analyseExprVariableUse(tok, n.(*ast.KeyValueExpr).Value)
	case *ast.ArrayType:
		//fmt.Println("ArrayType")
		p.analyseExprVariableUse(tok, n.(*ast.ArrayType).Len)
		p.analyseExprVariableUse(tok, n.(*ast.ArrayType).Elt)
	case *ast.StructType:
		//fmt.Println("StructType") //Nothing to do
	case *ast.FuncType:
		//fmt.Println("FuncType") //Nothing to do ... TODO check closures?
	case *ast.InterfaceType:
		//fmt.Println("InterfaceType") //Nothing to do
	case *ast.MapType:
		//fmt.Println("MapType")
		p.analyseExprVariableUse(tok, n.(*ast.MapType).Key)
		p.analyseExprVariableUse(tok, n.(*ast.MapType).Value)
	case *ast.ChanType:
		//fmt.Println("ChanType")
		p.analyseExprVariableUse(tok, n.(*ast.ChanType).Value)
	case nil:
		//ignore
	default:
		//		fmt.Println("default -> ", reflect.TypeOf(n))
	}
}

func (p *compiler) addVariableUse(tok token.Token, varName string) { // TODO change tok to Tok, not string
	if _, ok := p.reuseMap[tok]; !ok {
		p.reuseMap[tok] = make(map[string]bool)
	}
	p.reuseMap[tok][varName] = true
}

// function to generate boilerplate antha element code
func (p *compiler) genBoilerplate() {
	var boilerPlateTemplate = `// AsyncBag functions
func (e *{{.PkgName}}) Complete(params interface{}) {
	p := params.({{.PkgName}}ParamBlock)
	if p.Error {
{{range .Results}}
		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}
		return
	}
	r := new({{.PkgName}}ResultBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)
	if r.Error {
{{range .Results}}
		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}
		return
	}
{{range .StepsOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
	e.analysis(p, r)
		if r.Error {

{{range .Results}}
		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}
		return
	}
{{range .AnalysisOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
	e.validation(p, r)
		if r.Error {
{{range .Results}}
		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}
		return
	}
{{range .ValidationOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
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
		bpVars.Params = append(bpVars.Params, BoilerPlateParamVars{p.params[i].Names[0].Name, getTypeString(p.params[i].Type)})
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
		ts := getTypeString(i.Type)
		cs := getConfigString(i.Type)
		structSpec.Params = append(structSpec.Params, VarSpec{n, ts})
		structSpec.Inports = append(structSpec.Inports, VarSpec{n, ts})
		structSpec.Configs = append(structSpec.Configs, VarSpec{n, cs})
	}
	for _, i := range p.results {
		structSpec.Results = append(structSpec.Results, VarSpec{i.Names[0].Name, getTypeString(i.Type)})
		if string(i.Names[0].Name[0]) == strings.ToUpper(string(i.Names[0].Name[0])) {
			structSpec.Outports = append(structSpec.Outports, VarSpec{i.Names[0].Name, getTypeString(i.Type)})
		}
	}
	var b bytes.Buffer
	t.Execute(&b, structSpec)
	p.print(b, newline)
}

func (p *compiler) genExtraMethods() {
	p.genComponentInfo()
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
		structSpec.Inports = append(structSpec.Inports, VarSpec{i.Names[0].Name, getTypeString(i.Type)})
	}
	for _, i := range p.results {
		if string(i.Names[0].Name[0]) == strings.ToUpper(string(i.Names[0].Name[0])) {
			structSpec.Outports = append(structSpec.Outports, VarSpec{i.Names[0].Name, getTypeString(i.Type)})
		}
	}
	var b bytes.Buffer
	t.Execute(&b, structSpec)
	p.print(b, newline)

}

// helper function to make an ImportSpec node. Pass a nil name string to
// use a default namespace
func (p *compiler) getImportSpec(path string, name *string) (spec *ast.ImportSpec) {
	spec = new(ast.ImportSpec)
	if name != nil {
		spec.Name = ast.NewIdent(*name)
	}
	spec.Path = &ast.BasicLit{ValuePos: token.NoPos, Kind: token.STRING, Value: path}
	return
}

func (p *compiler) addPreamble(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.AnthaDecl:
			if decl.Tok == token.REQUIREMENTS {
				//_ = wunit.Make_units
				s2Lhs, _ := parser.ParseExpr("_")
				s2Rhs, _ := parser.ParseExpr("wunit.Make_units")
				s2 := &ast.AssignStmt{
					Lhs:    []ast.Expr{s2Lhs},
					Rhs:    []ast.Expr{s2Rhs},
					TokPos: decl.TokPos,
					Tok:    token.ASSIGN,
				}
				decl.Body.List = append([]ast.Stmt{s2}, decl.Body.List...)

				continue
			}
			// `_wrapper := execution.NewWrapper()`
			s1Lhs, _ := parser.ParseExpr("_wrapper")
			s1Rhs, _ := parser.ParseExpr("execution.NewWrapper(p.ID)")
			s1 := &ast.AssignStmt{
				Lhs:    []ast.Expr{s1Lhs},
				Rhs:    []ast.Expr{s1Rhs},
				TokPos: decl.TokPos,
				Tok:    token.DEFINE,
			}
			// `_ = _wrapper`
			s2Lhs, _ := parser.ParseExpr("_")
			s2Rhs, _ := parser.ParseExpr("_wrapper")
			s2 := &ast.AssignStmt{
				Lhs:    []ast.Expr{s2Lhs},
				Rhs:    []ast.Expr{s2Rhs},
				TokPos: decl.TokPos,
				Tok:    token.ASSIGN,
			}
			decl.Body.List = append([]ast.Stmt{s1, s2}, decl.Body.List...)
		default:
			continue
		}
	}
}

// updates in place the AST by altering any idents to the appropriate
// antha structures, and altering any assignments to
// channel operations
func (p *compiler) sugarAST(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.GenDecl:
			//fmt.Println("->", decl.Tok)
		case *ast.AnthaDecl:
			// fix the AST
			p.sugarStmts(decl.Body.List)
		default: // ignore all other decls... for now
			continue
		}
	}
}

func (p *compiler) sugarForParams(list []ast.Stmt) {
	for i := range list {
		if stmt, isAssign := list[i].(*ast.AssignStmt); isAssign {
			// test if a result is being assigned to.
			fix_assign := false
			for j := range stmt.Lhs {
				if ident, ok := stmt.Lhs[j].(*ast.Ident); ok {
					isUpper := string(ident.Name[0]) == strings.ToUpper(string(ident.Name[0]))
					if sugar, fixme := p.resultMap[ident.Name]; fixme && isUpper {
						//fix_assign = true
						ident.Name = sugar
					}
				}
			}
			//TODO: generate multiple stmts if needed
			if fix_assign {
				if len(stmt.Rhs) > 1 {
					// currently undefined so panic
					panic("too many right hand exprs")
				}
				list[i] = &ast.SendStmt{stmt.Lhs[0], stmt.TokPos, p.wrapExpr(stmt.Rhs[0])}
			}
		}

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
				fix_assign := false
				for j := range node.Lhs {
					if ident, ok := node.Lhs[j].(*ast.Ident); ok {
						isUpper := string(ident.Name[0]) == strings.ToUpper(string(ident.Name[0]))
						if sugar, fixme := p.resultMap[ident.Name]; fixme && isUpper { //fix_assign = true
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
					list[i] = &ast.SendStmt{node.Lhs[0], node.TokPos, p.wrapExpr(node.Rhs[0])}
				}
			default:

			}
			return true
		})
	}
}

func (p *compiler) sugarForIntrinsics(list []ast.Stmt) {
	for i := range list {
		ast.Inspect(list[i], func(n ast.Node) bool {
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
}

// Master function to handle any AST modifications needed for sugaring Antha code
// TODO: refactor into more granular functions as more gets added
func (p *compiler) sugarStmts(list []ast.Stmt) {
	p.sugarForParams(list)
	p.sugarForIntrinsics(list)
}

// helper function to wrap the expr in assignments
func (p *compiler) wrapExpr(e ast.Expr) (res ast.Expr) {
	Type := &ast.SelectorExpr{ast.NewIdent("execute"), ast.NewIdent("ThreadParam")}
	exp := &ast.SelectorExpr{ast.NewIdent("p"), ast.NewIdent("ID")}
	var elts []ast.Expr
	elts = append(elts, e)
	elts = append(elts, exp)
	res = &ast.CompositeLit{Type, token.NoPos, elts, token.NoPos}
	return
}

// helper function to generate a factory function for Antha Elements
func (p *compiler) getFactory() (b bytes.Buffer) {
	// internal struct rather than just string for future flexibility
	type factorySpec struct {
		Name string
	}
	var genTpl = `
	case "{{.Name}}":
		e = {{.Name}}.New()`
	var specs []factorySpec
	var startTpl = `
	func NewElement(name string) (e *AnthaElement) {
		select name {
			
	`
	var endTpl = `
		default:
			e = nil
	}
	return
}	
	`
	for protocol := range p.anthaImports {
		specs = append(specs, factorySpec{protocol})
	}

	t1 := template.Must(template.New("factory_start").Parse(startTpl))
	t1.Execute(&b, nil)

	t2 := template.Must(template.New("factory_elements").Parse(genTpl))
	for spec := range specs {
		t2.Execute(&b, spec)
	}

	t3 := template.Must(template.New("factory_end").Parse(endTpl))
	t3.Execute(&b, nil)

	return
}

// helper function to generate a wrapper main file to convert an element into a standalone binary
// Prerequisite: Antha elements have been indexed in this compiler context (such as by calling this
// after indexParams and sugarAST
func (p *compiler) standAlone(src *ast.File, packageRoute string) {

	var b bytes.Buffer

	p.GetMainPackageImports(&b, packageRoute)
	p.GetMainApp(&b)
	p.GetMainReferenceFunctions(&b)

	p.mainOutput = b.Bytes()

	return
}

func (p *compiler) GetMainPackageImports(b *bytes.Buffer, packageRoute string) {
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

//GetMainApp puts in b the contents of a go file capable of runing the component described in the
// compiler. This executable file will read its input from stdin and write it to stdout
func (p *compiler) GetMainApp(b *bytes.Buffer) {
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
	for i := range p.params {
		appVals.InPorts = append(appVals.InPorts, InPortSpec{p.params[i].Names[0].Name, p.pkgName})
	}
	for i := range p.results {
		if string(p.results[i].Names[0].String()[0]) == strings.ToUpper(string(p.results[i].Names[0].String()[0])) {
			appVals.OutPorts = append(appVals.OutPorts, OutPortSpec{p.results[i].Names[0].Name, p.pkgName})
		}
	}
	t.Execute(b, appVals)

	return
}

//GetMainReferenceFunctions puts in b the necessary funcitons for a main function to be able to read/write and execute an element
func (p *compiler) GetMainReferenceFunctions(b *bytes.Buffer) {
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
	for i := range p.params {
		rMSpec.InPorts = append(rMSpec.InPorts, InPortSpec{p.params[i].Names[0].Name})
	}
	for i := range p.results {
		if string(p.results[i].Names[0].String()[0]) == strings.ToUpper(string(p.results[i].Names[0].String()[0])) {
			rMSpec.OutPorts = append(rMSpec.OutPorts, OutPortSpec{p.results[i].Names[0].Name})
		}
	}
	t.Execute(b, rMSpec)

	return
}

// GenerateGraphLib builds a go file defining processes (i.e., go struct
// instances) defined by components
func GenerateComponentLib(b *bytes.Buffer, components []execute.ComponentInfo, workingDirectory string, package_ string) {
	type ComponentSpec struct {
		Name            string
		ConstructorFunc string
		InPorts         []string
		OutPorts        []string
	}
	type GraphRunnerSpec struct {
		Components       []ComponentSpec
		WorkingDirectory string
		Package          string
	}
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
{{range .InPorts}}	portMap["{{$componentName}}"]["{{$componentName}}"] = true
{{end}}
{{range .OutPorts}}	portMap["{{$componentName}}"]["{{$componentName}}"] = false
{{end}}
{{end}}
	c := make([]ComponentDesc, 0)
{{range .Components}}	c = append(c, ComponentDesc{Name: "{{.Name}}", Constructor: {{.Name}}.{{.ConstructorFunc}} })
{{end}}
	return c
}`

	var graphRunnerSpec GraphRunnerSpec
	graphRunnerSpec.Package = package_
	graphRunnerSpec.WorkingDirectory = workingDirectory
	for _, component := range components {
		var cs ComponentSpec
		cs.Name = component.Name
		cs.ConstructorFunc = fmt.Sprintf("New%s", component.Name)
		for _, ip := range component.InPorts {
			cs.InPorts = append(cs.InPorts, ip.Id)
		}
		for _, ip := range component.OutPorts {
			cs.OutPorts = append(cs.OutPorts, ip.Id)
		}
		graphRunnerSpec.Components = append(graphRunnerSpec.Components, cs)
	}

	t := template.Must(template.New("MainGrapher").Parse(tmplt))
	t.Execute(b, graphRunnerSpec)
}

// GenerateGraphRunner builds a go file capable of running FBP (Flow-based
// programming) graphs with processes defined by the components included in the
// components argument.
func GenerateGraphRunner(b *bytes.Buffer, components []execute.ComponentInfo, workingDirectory string) {
	type ComponentSpec struct {
		Name            string
		ConstructorFunc string
		InPorts         []string
		OutPorts        []string
	}
	type GraphRunnerSpec struct {
		Components       []ComponentSpec
		WorkingDirectory string
	}
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

	var graphRunnerSpec GraphRunnerSpec
	graphRunnerSpec.WorkingDirectory = workingDirectory
	for _, component := range components {
		var cs ComponentSpec
		cs.Name = component.Name
		cs.ConstructorFunc = fmt.Sprintf("New%s", component.Name)
		for _, ip := range component.InPorts {
			cs.InPorts = append(cs.InPorts, ip.Id)
		}
		for _, ip := range component.OutPorts {
			cs.OutPorts = append(cs.OutPorts, ip.Id)
		}
		graphRunnerSpec.Components = append(graphRunnerSpec.Components, cs)
	}

	t := template.Must(template.New("MainGrapher").Parse(tmplt))
	t.Execute(b, graphRunnerSpec)
}
