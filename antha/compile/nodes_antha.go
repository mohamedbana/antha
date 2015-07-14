// Implementations called from nodes.go
package compile

import (
	"bytes"
	"github.com/antha-lang/antha/antha/ast"
	"github.com/antha-lang/antha/antha/token"
	"text/template"
)

// Called from node.go walker to print antha code blocks
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
