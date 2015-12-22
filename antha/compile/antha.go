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
// 2 Royal College St, London NW1 0NH UK

// This file implements the generation of go code from the
// antha blocks. It is primarily used in nodes.go and compile.go

package compile

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
)

// Augmentation to compiler state to parse Antha files
type antha struct {
	pkgName     string
	params      []string                        // canonical order of params to ensure deterministic output
	paramTypes  map[string]string               // map from params to type names
	results     []string                        // canonical order of results to ensure deterministic output
	resultTypes map[string]string               // map from results to type names
	configTypes map[string]string               // map from params to config type names
	reuseMap    map[token.Token]map[string]bool // map use of data through execution phases
	intrinsics  map[string]string               // replacements for identifiers in expressions in functions
	types       map[string]string               // replacement for type names in type expressions and type and type lists
}

func (p *compiler) anthaInit() {
	p.paramTypes = make(map[string]string)
	p.resultTypes = make(map[string]string)
	p.configTypes = make(map[string]string)
	p.reuseMap = make(map[token.Token]map[string]bool)
	for _, tok := range []token.Token{token.ANALYSIS, token.VALIDATION, token.STEPS, token.SETUP, token.REQUIREMENTS} {
		p.reuseMap[tok] = make(map[string]bool)
	}
	p.intrinsics = map[string]string{
		"MixTo":        "_wrapper.MixTo",
		"MixInto":      "_wrapper.MixInto",
		"Mix":          "_wrapper.Mix",
		"Incubate":     "_wrapper.Incubate",
		"Centrifuge":   "_wrapper.Centrifuge",
		"Electroshock": "_wrapper.Electroshock",
		"ReadEM":       "_wrapper.ReadEM",
	}
	p.types = map[string]string{
		"Temperature":          "wunit.Temperature",
		"Time":                 "wunit.Time",
		"Length":               "wunit.Length",
		"Area":                 "wunit.Area",
		"Volume":               "wunit.Volume",
		"Concentration":        "wunit.Concentration",
		"Amount":               "wunit.Amount",
		"Mass":                 "wunit.Mass",
		"Angle":                "wunit.Angle",
		"Energy":               "wunit.Energy",
		"SubstanceQuantity":    "wunit.SubstanceQuantity",
		"Force":                "wunit.Force",
		"Pressure":             "wunit.Pressure",
		"SpecificHeatCapacity": "wunit.SpecificHeatCapacity",
		"Density":              "wunit.Density",
		"FlowRate":             "wunit.FlowRate",
		"Velocity":             "wunit.Velocity",
		"Rate":                 "wunit.Rate",
		"Resistance":           "wunit.Resistance",
		"Capacitance":          "wunit.Capacitance",
		"Voltage":              "wunit.Voltage",
		"Component":            "wtype.LHComponent",
		"Solution":             "wtype.LHComponent",
		"Plate":                "wtype.LHPlate",
		"Tipbox":               "wtype.LHTipbox",
		"Tip":                  "wtype.LHTip",
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

// Generate message handler for parameter
func (p *compiler) genOnFunction(name string, paramCount int) {
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

	var b bytes.Buffer
	t.Execute(&b, funcVals)
	p.print(b, newline)
}

// Return appropriate go type string for an antha (type) expr
func (p *compiler) getTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		if v, ok := p.types[t.Name]; ok {
			res = v
		} else {
			res = t.Name
		}
		return
	case *ast.SelectorExpr:
		res = p.getTypeString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType:
		// note: array types can use a param as the length, so must be allocated and treated as a slice since they can be dynamic
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
		if v, ok := p.types[t.Name]; ok {
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
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
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
			"github.com/antha-lang/antha/microArch/execution",
			"github.com/antha-lang/antha/antha/anthalib/wunit",
			"github.com/antha-lang/antha/antha/execute",
			"github.com/antha-lang/antha/flow",
			"sync",
			"runtime/debug",
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

	for _, param := range p.params {
		p.genOnFunction(param, len(p.params))
	}

	p.genStructs()
	p.genComponentInfo()
}

func sortKeys(m map[string]string) []string {
	sorted := make([]string, len(m))
	i := 0
	for k, _ := range m {
		sorted[i] = k
		i += 1
	}
	sort.Strings(sorted)
	return sorted
}

// Collect information needed in downstream generation passes
func (p *compiler) analyze(src *ast.File) {
	p.recordParams(src.Decls)
	p.recordOutputUse(src.Decls)
	p.pkgName = strings.Title(src.Name.Name)

	p.params = sortKeys(p.paramTypes)
	p.results = sortKeys(p.resultTypes)
}

// index all the spec definitions for inputs and outputs to element
func (p *compiler) recordParams(decls []ast.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.PARAMETERS, token.INPUTS:
				for _, spec := range decl.Specs {
					spec := spec.(*ast.ValueSpec)
					for _, name := range spec.Names {
						//TODO check input data cannot be lowercase
						param := name.String()
						p.paramTypes[param] = p.getTypeString(spec.Type)
						p.configTypes[param] = p.getConfigTypeString(spec.Type)
					}
				}
			case token.OUTPUTS, token.DATA:
				for _, spec := range decl.Specs {
					spec := spec.(*ast.ValueSpec)
					for _, name := range spec.Names {
						param := name.String()
						p.resultTypes[param] = p.getTypeString(spec.Type)
					}
				}
			}
		}
	}
}

// Track use of output parameters in antha blocks
func (p *compiler) recordOutputUse(decls []ast.Decl) {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case *ast.AnthaDecl:
			tok := decl.Tok
			ast.Inspect(decl.Body, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.Ident:
					if _, is := p.resultTypes[n.Name]; is {
						varName := n.Name
						p.reuseMap[tok][varName] = true
					}
				}
				return true
			})
		}
	}
}

// function to generate boilerplate antha element code //TODO insert blockID
func (p *compiler) genBoilerplate() {
	var boilerPlateTemplate = `// AsyncBag functions
func (e *{{.PkgName}}) Complete(params interface{}) {
	p := params.({{.PkgName}}ParamBlock)
	if p.Error {
{{range .Results}}		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: nil, ID: p.ID, Error: true}
{{end}}		return
	}
	r := new({{.PkgName}}ResultBlock)
	defer func() {
		if res := recover(); res != nil {
{{range .Results}}		e.{{.OutputVariableName}} <- execute.ThreadParam{Value: res, ID: p.ID, Error: true}
{{end}}		execute.AddError(&execute.RuntimeError{BaseError: res, Stack: debug.Stack()})
			return
		}
	}()
	e.startup.Do(func() { e.setup(p) })
	e.steps(p, r)
{{range .StepsOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
	e.analysis(p, r)
{{range .AnalysisOutput}}
	e.{{.OutputVariableName}} <- execute.ThreadParam{Value: r.{{.ResultVariableName}}, ID: p.ID, Error: false}
{{end}}
	e.validation(p, r)
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
	res.BlockID = m["{{.IdVarName}}"].(execute.ThreadParam).BlockID
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
	for _, name := range p.params {
		bpVars.Params = append(bpVars.Params, BoilerPlateParamVars{name, p.paramTypes[name]})
	}
	// TODO(ddn): Check if first or alphabetically first param can serve as ID
	if len(p.params) > 0 {
		bpVars.HasIdLine = true
		bpVars.IdVarName = p.params[0]
	}
	for _, name := range p.results {
		if name[0:1] == strings.ToUpper(name[0:1]) { // only upper case are external
			bpVars.Results = append(bpVars.Results, OutputRedirection{name, name})
			_, an := p.reuseMap[token.ANALYSIS][name]
			_, va := p.reuseMap[token.VALIDATION][name]
			if !(an || va) { //steps is the last one it appeared (if it did)
				bpVars.StepsOutput = append(bpVars.StepsOutput, OutputRedirection{name, name})
			} else if !va { //analysis is the last one it appeared (if it did)
				bpVars.AnalysisOutput = append(bpVars.AnalysisOutput, OutputRedirection{name, name})
			} else { //validation is the last one it appeared (if it did)
				bpVars.ValidationOutput = append(bpVars.ValidationOutput, OutputRedirection{name, name})
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
	BlockID	execute.BlockID
	Error	bool{{range .Params}}
	{{.Name}}		{{.Type}}{{end}}
}

type {{.ElementName}}Config struct{
	ID		execute.ThreadID
	BlockID	execute.BlockID
	Error	bool{{range .Configs}}
	{{.Name}}		{{.Type}}{{end}}
}

type {{.ElementName}}ResultBlock struct{
	ID		execute.ThreadID
	BlockID	execute.BlockID
	Error	bool{{range .Results}}
	{{.Name}}		{{.Type}}{{end}}
}

type {{.ElementName}}JSONBlock struct{
	ID			*execute.ThreadID
	BlockID		*execute.BlockID
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
	for _, name := range p.params {
		ts := p.paramTypes[name]
		cs := p.configTypes[name]
		structSpec.Params = append(structSpec.Params, VarSpec{name, ts})
		structSpec.Inports = append(structSpec.Inports, VarSpec{name, ts})
		structSpec.Configs = append(structSpec.Configs, VarSpec{name, cs})
	}
	for _, name := range p.results {
		structSpec.Results = append(structSpec.Results, VarSpec{name, p.resultTypes[name]})
		if name[0:1] == strings.ToUpper(name[0:1]) {
			structSpec.Outports = append(structSpec.Outports, VarSpec{name, p.resultTypes[name]})
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
	for _, name := range p.params {
		structSpec.Inports = append(structSpec.Inports, VarSpec{name, p.paramTypes[name]})
	}
	for _, name := range p.results {
		if name[0:1] == strings.ToUpper(name[0:1]) {
			structSpec.Outports = append(structSpec.Outports, VarSpec{name, p.resultTypes[name]})
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
			s1Rhs, _ := parser.ParseExpr("execution.NewWrapper(p.ID, p.BlockID, p)")
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

			s3Lhs, _ := parser.ParseExpr("_")
			s3Rhs, _ := parser.ParseExpr("_wrapper.WaitToEnd()") //TODO won't handle intermediate waiting
			s3 := &ast.AssignStmt{
				Lhs:    []ast.Expr{s3Lhs},
				Rhs:    []ast.Expr{s3Rhs},
				TokPos: d.TokPos,
				Tok:    token.ASSIGN,
			}

			d.Body.List = append([]ast.Stmt{s1, s2}, d.Body.List...)
			d.Body.List = append(d.Body.List, s3)
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
			p.sugarForTypes(decl.Body)
		default:
			continue
		}
	}
}

// Return appropriate nested SelectorExpr for the replacement for Identifier
func (p *compiler) sugarForIdent(t *ast.Ident) ast.Expr {
	if v, ok := p.types[t.Name]; ok {
		cs := strings.Split(v, ".")
		var base ast.Expr = &ast.Ident{NamePos: t.NamePos, Name: cs[0]}
		for _, c := range cs[1:] {
			base = &ast.SelectorExpr{X: base, Sel: &ast.Ident{NamePos: t.NamePos, Name: c}}
		}
		return base
	} else {
		return t
	}
}

// Return appropriate go type for an antha (type) expr
func (p *compiler) sugarForType(t ast.Expr) ast.Expr {
	switch t := t.(type) {
	case nil:
	case *ast.Ident:
		return p.sugarForIdent(t)
	case *ast.ParenExpr:
		return &ast.ParenExpr{Lparen: t.Lparen, X: p.sugarForType(t.X), Rparen: t.Rparen}
	case *ast.SelectorExpr:
		return t
	case *ast.StarExpr:
		return &ast.StarExpr{Star: t.Star, X: p.sugarForType(t.X)}
	case *ast.ArrayType:
		return &ast.ArrayType{Lbrack: t.Lbrack, Len: t.Len, Elt: p.sugarForType(t.Elt)}
	case *ast.StructType:
		return &ast.StructType{Struct: t.Struct, Fields: p.sugarForFieldList(t.Fields), Incomplete: t.Incomplete}
	case *ast.FuncType:
		return &ast.FuncType{Func: t.Func, Params: p.sugarForFieldList(t.Params), Results: p.sugarForFieldList(t.Results)}
	case *ast.InterfaceType:
		return &ast.InterfaceType{Interface: t.Interface, Methods: p.sugarForFieldList(t.Methods), Incomplete: t.Incomplete}
	case *ast.MapType:
		return &ast.MapType{Map: t.Map, Key: p.sugarForType(t.Key), Value: p.sugarForType(t.Value)}
	case *ast.ChanType:
		return &ast.ChanType{Begin: t.Begin, Arrow: t.Arrow, Dir: t.Dir, Value: p.sugarForType(t.Value)}
	default:
		log.Panicf("unexpected expression %s of type %s", t, reflect.TypeOf(t))
	}

	return t
}

func (p *compiler) sugarForFieldList(t *ast.FieldList) *ast.FieldList {
	if t == nil {
		return nil
	}
	var fields []*ast.Field
	for _, f := range t.List {
		fields = append(fields, &ast.Field{Doc: f.Doc, Names: f.Names, Type: p.sugarForType(f.Type), Tag: f.Tag, Comment: f.Comment})
	}
	return &ast.FieldList{Opening: t.Opening, List: fields, Closing: t.Closing}
}

// Replace bare antha types with go qualified names
func (p *compiler) sugarForTypes(root ast.Node) ast.Node {
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		case nil:
			return false
		case *ast.FuncLit:
			n.Type = p.sugarForType(n.Type).(*ast.FuncType)
			return false
		case *ast.CompositeLit:
			n.Type = p.sugarForType(n.Type)
			return false
		case *ast.TypeAssertExpr:
			n.Type = p.sugarForType(n.Type)
			return false
		case *ast.ValueSpec:
			n.Type = p.sugarForType(n.Type)
			return false
		default:
			return true
		}
	})
	return root
}

// Replace bare antha identifiers with go qualified names
func (p *compiler) sugarForParams(body *ast.BlockStmt) {
	// TODO(ddn): merge with sugarForTypes or clone and adopt a similar strategy to
	// restrict replacements to unqualified variable use. Right now, the code below
	// can rewrite field accesses, etc.
	list := body.List

	for i := range list {
		ast.Inspect(list[i], func(n ast.Node) bool {
			switch node := n.(type) {
			case nil:
				return false
			case *ast.Ident:
				// sugar the name if it is a known param
				if _, ok := p.paramTypes[node.Name]; ok {
					node.Name = "p." + node.Name
				} else if _, ok := p.resultTypes[node.Name]; ok {
					node.Name = "r." + node.Name
				}
				// test if the left hand side of an assignment is to a results variable
				// if so, sugar it by replacing the = with a channel operation
			case *ast.AssignStmt:
				// test if a result is being assigned to.
				fix_assign := false
				for j := range node.Lhs {
					if ident, ok := node.Lhs[j].(*ast.Ident); ok {
						isUpper := ident.Name[0:1] == strings.ToUpper(ident.Name[0:1])
						if _, fixme := p.resultTypes[ident.Name]; fixme && isUpper {
							//fix_assign = true
							ident.Name = "r." + ident.Name
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
	// TODO(ddn): merge with sugarForTypes or clone and adopt a similar strategy to
	// restrict replacements to unqualified variable use. Right now, the code below
	// can rewrite field accesses, etc.
	ast.Inspect(root, func(n ast.Node) bool {
		switch n := n.(type) {
		case nil:
			return false
		case *ast.Ident:
			// sugar the name if it is a known param
			if sugar, ok := p.intrinsics[n.Name]; ok {
				n.Name = sugar
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
