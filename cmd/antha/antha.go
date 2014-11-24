// antha/cmd/antha/antha.go: Part of the Antha language
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

package main

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/printer"
	"github.com/antha-lang/antha/scanner"
	"github.com/antha-lang/antha/token"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// execution variables
var (
	exitCode   = 0
	fileSet    = token.NewFileSet() // per process FileSet
	parserMode parser.Mode
)

// parameters to control code formatting
const (
	tabWidth    = 8
	printerMode = printer.UseSpaces | printer.TabIndent
)

// command line parameters
var (
	// main operation modes
	trace     = flag.Bool("trace", false, "show AST trace")
	allErrors = flag.Bool("errors", false, "report all errors (not just the first 10 on different lines)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: antha [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// utility function
func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func initParserMode() {
	parserMode = parser.ParseComments
	if *allErrors {
		parserMode |= parser.AllErrors
	}
	if *trace {
		parserMode |= parser.Trace
	}
}

// Utility function to check file extesion
func isAnthaFile(f os.FileInfo) bool {
	// ignore non-Antha or Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".an")
}

func main() {
	// call gofmtMain in a separate function
	// so that it can use defer and have them
	// run before the exit.
	anthaMain()
	os.Exit(exitCode)
}

func anthaMain() {
	flag.Usage = usage
	flag.Parse()

	initParserMode()

	// try to parse standard input if no files or directories were passed in
	if flag.NArg() == 0 {
		if err := processFile("<standard input>", os.Stdin, os.Stdout, true); err != nil {
			report(err)
		}
		return
	}

	// parse every filename or directory passed in as input
	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				report(err)
			}
		}
	}
}

// recursive support function for walkpath
func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isAnthaFile(f) {
		err = processFile(path, nil, os.Stdout, false)
	}
	if err != nil {
		report(err)
	}
	return nil
}

// If in == nil, the source is the contents of the file with the given filename.
func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	file, adjust, err := parse(fileSet, filename, src, stdin)
	if err != nil {
		return err
	}

	// if this isn't an antha file, bail
	if file.Tok != token.PROTOCOL {
		return errors.New(filename + " is not a valid Antha file")
	}

	ast.SortImports(fileSet, file)

	err = translateAST(fileSet, file)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = (&printer.Config{Mode: printerMode, Tabwidth: tabWidth}).Fprint(&buf, fileSet, file)
	if err != nil {
		return err
	}
	res := buf.Bytes()
	if adjust != nil {
		res = adjust(src, res)
	}

	// save the output as a translated .go file in same location as .an
	err = ioutil.WriteFile(strings.TrimSuffix(filename, ".an")+".go", res, 0777)
	if err != nil {
		return err
	}

	_, err = out.Write(res)

	return err
}

// parse parses src, which was read from filename,
// as a Antha source file or statement list.
func parse(fset *token.FileSet, filename string, src []byte, stdin bool) (*ast.File, func(orig, src []byte) []byte, error) {
	// Try as whole source file.
	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	// If the error is that the source file didn't begin with a
	// package line and this is standard input, fall through to
	// try as a source fragment.  Stop and return on any other error.
	if !stdin || !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	// If this is a declaration list, make it a source file
	// by inserting a package clause.
	// Insert using a ;, not a newline, so that the line numbers
	// in psrc match the ones in src.
	psrc := append([]byte("package p;"), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the package clause.
			// Anthafmt has turned the ; into a \n.
			src = src[len("package p\n"):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	// If the error is that the source file didn't begin with a
	// declaration, fall through to try as a statement list.
	// Stop and return on any other error.
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	// If this is a statement list, make it a source file
	// by inserting a package clause and turning the list
	// into a function body.  This handles expressions too.
	// Insert using a ;, not a newline, so that the line numbers
	// in fsrc match the ones in src.
	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			// Remove the wrapping.
			// Anthafmt has turned the ; into a \n\n.
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			// Anthafmt has also indented the function body one level.
			// Remove that indent.
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	// Failed, and out of options.
	return nil, nil, err
}

// Utility function for matchSpace
func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}

// matchSpace reformats src to use the same space context as orig.
// 1) If orig begins with blank lines, matchSpace inserts them at the beginning of src.
// 2) matchSpace copies the indentation of the first non-blank line in orig
//    to every non-blank line in src.
// 3) matchSpace copies the trailing space from orig and uses it in place
//   of src's trailing space.
func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], line[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' { // not blank
			b.Write(indent)
		}
		b.Write(line)
	}
	b.Write(after)
	return b.Bytes()
}

// convert antha AST to go AST
func translateAST(fset *token.FileSet, f *ast.File) error {
	// Create an ast.CommentMap from the ast.File's comments.
	// This helps keeping the association between comments
	// and AST nodes.
	//cmap := ast.NewCommentMap(fset, f, f.Comments)

	//fmt.Println(len(cmap))

	return nil
}