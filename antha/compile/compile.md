---
layout: default
type: api
navgroup: docs
shortname: antha/compile
title: antha/compile
apidocs:
  published: 2015-06-25
  antha_version: 0.0.2
  package: antha/compile
---
# compile
--
    import "."


## Usage

#### func  Fprint

```go
func Fprint(output io.Writer, fset *token.FileSet, node interface{}) error
```
Fprint "pretty-prints" an AST node to output. It calls Config.Fprint with
default settings.

#### func  GenerateComponentLib

```go
func GenerateComponentLib(b *bytes.Buffer, components []execute.ComponentInfo, workingDirectory string, package_ string)
```
GenerateGraphLib builds a go file defining processes (i.e., go class instances)
defined by components

#### func  GenerateGraphRunner

```go
func GenerateGraphRunner(b *bytes.Buffer, components []execute.ComponentInfo, workingDirectory string)
```
GenerateGraphRunner builds a go file capable of running fbp graphs with
processes defined by the components included in the components argument.

#### type CommentedNode

```go
type CommentedNode struct {
	Node     interface{} // *ast.File, or ast.Expr, ast.Decl, ast.Spec, or ast.Stmt
	Comments []*ast.CommentGroup
}
```

A CommentedNode bundles an AST node and corresponding comments. It may be
provided as argument to any of the Fprint functions.

#### type Config

```go
type Config struct {
	Mode     Mode // default: 0
	Tabwidth int  // default: 8
	Indent   int  // default: 0 (all code is indented at least by this much)
}
```

A Config node controls the output of Fprint.

#### func (*Config) Fprint

```go
func (cfg *Config) Fprint(output io.Writer, fset *token.FileSet, node interface{}) error
```
Fprint "pretty-prints" an AST node to output for a given configuration cfg.
Position information is interpreted relative to the file set fset. The node type
must be *ast.File, *CommentedNode, []ast.Decl, []ast.Stmt, or
assignment-compatible to ast.Expr, ast.Decl, ast.Spec, or ast.Stmt.

#### func (*Config) GetFileComponentInfo

```go
func (cfg *Config) GetFileComponentInfo(fset *token.FileSet, node interface{}) execute.ComponentInfo
```

#### func (*Config) MainFprint

```go
func (cfg *Config) MainFprint(output io.Writer, fset *token.FileSet, node interface{}, packageRoute string) error
```
MainFprintf prints the contents of a main function for the file given in the
buffer given as parameter

#### type Mode

```go
type Mode uint
```

A Mode value is a set of flags (or 0). They control printing.

```go
const (
	RawFormat Mode = 1 << iota // do not use a tabwriter; if set, UseSpaces is ignored
	TabIndent                  // use tabs for indentation independent of UseSpaces
	UseSpaces                  // use spaces instead of tabs for alignment
	SourcePos                  // emit //line comments to preserve original source positions
)
```
