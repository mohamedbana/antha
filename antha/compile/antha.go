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
	pkgName     string
	params      []string                        // canonical order of params to ensure deterministic output
	paramTypes  map[string]string               // map from params to type names
	results     []string                        // canonical order of results to ensure deterministic output
	resultTypes map[string]string               // map from results to type names
	configTypes map[string]string               // map from params to config type names
	reuseMap    map[token.Token]map[string]bool // map use of data through execution phases
	intrinsics  map[string]string               // replacements for identifiers in expressions in functions
	types       map[string]string               // replacement for type names in type expressions and type and type lists
	imports     map[string]string               // additional imports with additional expression to suppress unused import
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
		"MixTo":    "execute.MixTo",
		"MixInto":  "execute.MixInto",
		"Mix":      "execute.Mix",
		"Incubate": "execute.Incubate",
		"Call":     "execute.Call",
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
	}
	p.imports = map[string]string{
		"github.com/antha-lang/antha/antha/anthalib/wunit":             "wunit.Make_units",
		"github.com/antha-lang/antha/execute":                          "execute.MixInto",
		"github.com/antha-lang/antha/inject":                           "",
		"github.com/antha-lang/antha/bvendor/golang.org/x/net/context": "",
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
	var imports []string
	for i := range p.imports {
		imports = append(imports, i)
	}
	p.addAnthaImports(src, imports)
	p.sugarAST(src.Decls)
	p.addPreamble(src.Decls)
	p.removeParamDecls(src)
}

// Print out additional go code for each antha file
func (p *compiler) generate() {
	p.genBoilerplate()
	p.genStructs()
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
	var boilerPlateTemplate = `
func _run(_ctx context.Context, value inject.Value) (inject.Value, error) {
	input := &Input_{}
	output := &Output_{}
	if err := inject.Assign(value, input); err != nil {
		return nil, err
	}
	_setup(_ctx, input)
	_steps(_ctx, input, output)
	_analysis(_ctx, input, output)
	_validation(_ctx, input, output)
	return inject.MakeValue(output), nil
}

var (
	{{range .Uses}} _ = {{.}}
	{{end}}
)

func New() interface{} {
	return &Element_{
		inject.CheckedRunner {
			RunFunc: _run,
			In: 	 &Input_{},
			Out: 	 &Output_{},
		},
	}
}
`
	var params struct {
		Uses []string
	}
	for _, use := range p.imports {
		if len(use) > 0 {
			params.Uses = append(params.Uses, use)
		}
	}
	sort.Strings(params.Uses)

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(boilerPlateTemplate))
	t.Execute(&buf, params)
	p.print(buf, newline)
}

// helper function to generate the paramblock and primary element structs
func (p *compiler) genStructs() {
	var structTemplate = `type Element_ struct {
	inject.CheckedRunner
}

type Input_ struct {
	{{range .Params}}{{.Name}} {{.Type}}
	{{end}}
}

type Output_ struct{
	{{range .Results}}{{.Name}} {{.Type}}
	{{end}}
}`

	type VarSpec struct {
		Name, Type string
	}
	type StructSpec struct {
		ElementName string
		Params      []VarSpec
		Results     []VarSpec
	}

	t := template.Must(template.New(p.pkgName + "_struct").Parse(structTemplate))

	var structSpec StructSpec
	structSpec.ElementName = p.pkgName
	for _, name := range p.params {
		ts := p.paramTypes[name]
		structSpec.Params = append(structSpec.Params, VarSpec{name, ts})
	}
	for _, name := range p.results {
		structSpec.Results = append(structSpec.Results, VarSpec{name, p.resultTypes[name]})
	}
	var b bytes.Buffer
	t.Execute(&b, structSpec)
	p.print(b, newline)
}

// Add statements to the beginning of antha decls
func (p *compiler) addPreamble(decls []ast.Decl) {
	//for _, decl := range decls {
	//	switch d := decl.(type) {
	//	case *ast.AnthaDecl:
	//		// `_wrapper := execution.NewWrapper()`
	//		s1Lhs, _ := parser.ParseExpr("_wrapper")
	//		s1Rhs, _ := parser.ParseExpr("execution.NewWrapper(p.ID, p.BlockID, p)")
	//		s1 := &ast.AssignStmt{
	//			Lhs:    []ast.Expr{s1Lhs},
	//			Rhs:    []ast.Expr{s1Rhs},
	//			TokPos: d.TokPos,
	//			Tok:    token.DEFINE,
	//		}
	//		d.Body.List = append([]ast.Stmt{s1, s2}, d.Body.List...)
	//	default:
	//		continue
	//	}
	//}
}

// Update AST for antha semantics
func (p *compiler) sugarAST(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.GenDecl:
		case *ast.AnthaDecl:
			ast.Inspect(decl.Body, p.inspectForParams)
			ast.Inspect(decl.Body, p.inspectForIntrinsics)
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
func (p *compiler) inspectForParams(node ast.Node) bool {
	// sugar if it is a known param
	rewriteIdent := func(node *ast.Ident) {
		if _, ok := p.paramTypes[node.Name]; ok {
			node.Name = "_input." + node.Name
		} else if _, ok := p.resultTypes[node.Name]; ok {
			node.Name = "_output." + node.Name
		}
	}

	rewriteAssignLhs := func(node *ast.AssignStmt) {
		for j := range node.Lhs {
			if ident, ok := node.Lhs[j].(*ast.Ident); ok {
				isUpper := ident.Name[0:1] == strings.ToUpper(ident.Name[0:1])
				if _, fixme := p.resultTypes[ident.Name]; fixme && isUpper {
					ident.Name = "_output." + ident.Name
				}
			}
		}
	}

	switch n := node.(type) {
	case nil:
		return false
	case *ast.AssignStmt:
		rewriteAssignLhs(n)
	case *ast.KeyValueExpr:
		if _, identKey := n.Key.(*ast.Ident); identKey {
			// Skip identifiers that are keys
			ast.Inspect(n.Value, p.inspectForParams)
			return false
		}
	case *ast.Ident:
		rewriteIdent(n)
	case *ast.SelectorExpr:
		// Skip identifiers that are field accesses
		ast.Inspect(n.X, p.inspectForParams)
		return false
	}
	return true
}

// Replace bare antha function names with go qualified names
func (p *compiler) inspectForIntrinsics(node ast.Node) bool {
	switch n := node.(type) {
	case nil:
	case *ast.CallExpr:
		if ident, direct := n.Fun.(*ast.Ident); !direct {
			break
		} else if sugar, ok := p.intrinsics[ident.Name]; !ok {
			break
		} else {
			ident.Name = sugar
			expr, _ := parser.ParseExpr("_ctx")
			n.Args = append([]ast.Expr{expr}, n.Args...)
		}
	default:
	}
	return true
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
