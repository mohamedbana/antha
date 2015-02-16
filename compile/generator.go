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
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/token"
	"log"
	"reflect"
	"strings"
	"text/template"
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
	token.STEPS:        `func (e *{{.PkgName}}) steps(p ParamBlock)`,
	token.REQUIREMENTS: `func (e *{{.PkgName}}) requirements()`,
	token.SETUP:        `func (e *{{.PkgName}}) setup(p ParamBlock)`,
	token.ANALYSIS:     `func (e *{{.PkgName}}) analysis(p ParamBlock, r ResultBlock)`,
	token.VALIDATION:   `func (e *{{.PkgName}}) validation(p ParamBlock, r ResultBlock)`,
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
	
import "github.com/antha-lang/antha/execute"
import "github.com/antha-lang/goflow"
import "sync"
import "log"
import "bytes"
import "encoding/json"
import "io"

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
	for param := range p.paramMap {
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
		bag = new(AsyncBag)
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
	return
}

// utility function to get type strings from an expr
func getTypeString(e ast.Expr) (res string) {
	switch t := e.(type) {
	case *ast.Ident:
		res = t.Name
		return
	case *ast.SelectorExpr:
		res = getTypeString(t.X) + "." + t.Sel.Name
		return
	case *ast.ArrayType: // note: array types can use a param as the length, so must be allocated and treated as a slice since they can be dynamic
		res = "[]" + getTypeString(t.Elt)
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
						param := strings.Title(tmp.Names[k].String())
						p.paramTypes[param] = getTypeString(tmp.Type)
						p.paramMap[param] = "p." + param
					}
				}
			case token.OUTPUTS, token.DATA:
				for j := range decl.Specs {
					tmp := decl.Specs[j].(*ast.ValueSpec)
					p.results = append(p.results, tmp)
					for k := range tmp.Names {
						result := strings.Title(tmp.Names[k].String())
						p.paramTypes[result] = getTypeString(tmp.Type)
						p.resultMap[result] = "e." + result
					}

				}
			}
		default:
			continue
		}
	}
}

// function to generate boilerplate antha element code
func (p *compiler) genBoilerplate() {
	var tpl = `// AsyncBag functions
func (e *{{.}}) Complete(params interface{}) {
	p := params.(ParamBlock)
	e.startup.Do(func() { e.setup(p) })
	e.steps(p)
	
}

// empty function for interface support
func (e *{{.}}) anthaElement() {}

// init function, read characterization info from seperate file to validate ranges?
func (e *{{.}}) init() {
	e.params = make(map[execute.ThreadID]*execute.AsyncBag)
}

func New() *{{.}} {
	e := new({{.}})
	e.init()
	return e
}

// Mapper function
func (e *{{.}}) Map(m map[string]interface{}) interface{} {
	var res ParamBlock
`
	type mapLine struct {
		Name, Type string
	}

	var lineTpl = `
	res.{{.Name}} = m["{{.Name}}"].(execute.ThreadParam).Value.({{.Type}})	
`
	var endTpl = `
	return res
}

`
	var b bytes.Buffer
	t := template.Must(template.New(p.pkgName + "_boilerplate").Parse(tpl))
	t.Execute(&b, p.pkgName)
	t2 := template.Must(template.New(p.pkgName + "_lineTpl").Parse(lineTpl))

	for i := range p.params {
		var spec mapLine
		spec.Name = p.params[i].Names[0].Name
		spec.Type = getTypeString(p.params[i].Type)
		t2.Execute(&b, spec)
	}
	t3 := template.Must(template.New(p.pkgName + "_boiler_end").Parse(endTpl))
	t3.Execute(&b, nil)

	p.print(b, newline)
}

// helper function to generate the paramblock and primary element structs
// TODO: clean this up using template functions instead
func (p *compiler) genStructs() {
	type typeSpec struct {
		Name, Type string
	}
	var element = `type {{.}} struct {
	flow.Component                    // component "superclass" embedded
	lock           sync.Mutex
	startup        sync.Once	
	params         map[execute.ThreadID]*execute.AsyncBag`
	var param = `
	{{.}}          <-chan execute.ThreadParam`
	var result = `
	{{.}}      chan<- execute.ThreadParam`
	var paramBlock = `
type ParamBlock struct {
	ID        execute.ThreadID`
	var resultBlock = `
type ResultBlock struct {
	ID        execute.ThreadID`
	var paramEntry = `
	{{.Name}} {{.Type}}`

	// template for JSON decoder struct
	var jsonBlock = `
type JSONBlock struct {
	ID        *execute.ThreadID`
	var jsonEntry = `
	{{.Name}} *{{.Type}}`

	var b bytes.Buffer

	t1 := template.Must(template.New(p.pkgName + "_struct").Parse(element))
	t1.Execute(&b, p.pkgName)

	t2 := template.Must(template.New(p.pkgName + "_struct_param").Parse(param))
	for i := range p.params {
		t2.Execute(&b, p.params[i].Names[0].Name)
	}
	t3 := template.Must(template.New(p.pkgName + "_struct_result").Parse(result))
	for i := range p.results {
		t3.Execute(&b, p.results[i].Names[0].Name)
	}
	b.WriteString("\n}\n")

	t4 := template.Must(template.New(p.pkgName + "_paramBlock").Parse(paramBlock))
	t4.Execute(&b, p.pkgName)

	t5 := template.Must(template.New(p.pkgName + "_param_entry").Parse(paramEntry))
	for i := range p.params {
		var spec typeSpec
		spec.Name = p.params[i].Names[0].Name
		spec.Type = getTypeString(p.params[i].Type)
		t5.Execute(&b, spec)
	}
	b.WriteString("\n}\n")

	t6 := template.Must(template.New(p.pkgName + "_resultBlock").Parse(resultBlock))
	t6.Execute(&b, p.pkgName)

	for i := range p.results {
		var spec typeSpec
		spec.Name = p.results[i].Names[0].Name
		spec.Type = getTypeString(p.results[i].Type)
		t5.Execute(&b, spec)
	}
	b.WriteString("\n}\n")

	t7 := template.Must(template.New(p.pkgName + "_jsonBlock").Parse(jsonBlock))
	t7.Execute(&b, p.pkgName)

	t8 := template.Must(template.New(p.pkgName + "_json_entry").Parse(jsonEntry))

	for i := range p.results {
		var spec typeSpec
		spec.Name = p.results[i].Names[0].Name
		spec.Type = getTypeString(p.results[i].Type)
		t8.Execute(&b, spec)
	}
	b.WriteString("\n}\n")

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

// updates in place the AST by altering any idents to the appropriate
// antha structures, and altering any assignments to
// channel operations
func (p *compiler) sugarAST(d []ast.Decl) {
	for i := range d {
		switch decl := d[i].(type) {
		case *ast.AnthaDecl:
			// fix the AST
			p.sugarStmts(decl.Body.List)
		default: // ignore all other decls... for now
			continue
		}
	}
}

// Master function to handle any AST modifications needed for sugaring Antha code
// TODO: refactor into more granular functions as more gets added
func (p *compiler) sugarStmts(list []ast.Stmt) {
	//fmt.Println("sugarStmts")
	for i := range list {
		//fmt.Println(list[i], reflect.TypeOf(list[i]))

		if stmt, isAssign := list[i].(*ast.AssignStmt); isAssign {
			// test if a result is being assigned to.
			fix_assign := false
			for j := range stmt.Lhs {
				if ident, ok := stmt.Lhs[j].(*ast.Ident); ok {
					if sugar, fixme := p.resultMap[ident.Name]; fixme {
						fix_assign = true
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
				//fmt.Println(n, reflect.TypeOf(n))

				// sugar the name if it is a known param
				if sugar, ok := p.paramMap[node.Name]; ok {
					node.Name = sugar
				} else if sugar, ok := p.resultMap[node.Name]; ok {
					node.Name = sugar
				}
				// test if the left hand side of an assignment is to a results variable
				// if so, sugar it by replacing the = with a channel operation
			case *ast.AssignStmt:
			default:
				//fmt.Println(n, reflect.TypeOf(n))

			}
			return true
		})
	}
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
func (p *compiler) standAlone(src *ast.File) {

}
