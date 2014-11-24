// antha/doc/doc_test.go: Part of the Antha language
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


package doc

import (
	"github.com/antha-lang/antha/parser"
	"github.com/antha-lang/antha/printer"
	"github.com/antha-lang/antha/token"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

var update = flag.Bool("update", false, "update golden (.out) files")
var files = flag.String("files", "", "consider only Go/Antha test files matching this regular expression")

const dataDir = "testdata"

var templateTxt = readTemplate("template.txt")

func readTemplate(filename string) *template.Template {
	t := template.New(filename)
	t.Funcs(template.FuncMap{
		"node":     nodeFmt,
		"synopsis": synopsisFmt,
		"indent":   indentFmt,
	})
	return template.Must(t.ParseFiles(filepath.Join(dataDir, filename)))
}

func nodeFmt(node interface{}, fset *token.FileSet) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)
	return strings.Replace(strings.TrimSpace(buf.String()), "\n", "\n\t", -1)
}

func synopsisFmt(s string) string {
	const n = 64
	if len(s) > n {
		// cut off excess text and go back to a word boundary
		s = s[0:n]
		if i := strings.LastIndexAny(s, "\t\n "); i >= 0 {
			s = s[0:i]
		}
		s = strings.TrimSpace(s) + " ..."
	}
	return "// " + strings.Replace(s, "\n", " ", -1)
}

func indentFmt(indent, s string) string {
	end := ""
	if strings.HasSuffix(s, "\n") {
		end = "\n"
		s = s[:len(s)-1]
	}
	return indent + strings.Replace(s, "\n", "\n"+indent, -1) + end
}

func isGoFile(fi os.FileInfo) bool {
	name := fi.Name()
	return !fi.IsDir() &&
		len(name) > 0 && name[0] != '.' && // ignore .files
		filepath.Ext(name) == ".go"
}

type bundle struct {
	*Package
	FSet *token.FileSet
}

func test(t *testing.T, mode Mode) {
	// determine file filter
	filter := isGoFile
	if *files != "" {
		rx, err := regexp.Compile(*files)
		if err != nil {
			t.Fatal(err)
		}
		filter = func(fi os.FileInfo) bool {
			return isGoFile(fi) && rx.MatchString(fi.Name())
		}
	}

	// get packages
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dataDir, filter, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	// test packages
	for _, pkg := range pkgs {
		importpath := dataDir + "/" + pkg.Name
		doc := New(pkg, importpath, mode)

		// golden files always use / in filenames - canonicalize them
		for i, filename := range doc.Filenames {
			doc.Filenames[i] = filepath.ToSlash(filename)
		}

		// print documentation
		var buf bytes.Buffer
		if err := templateTxt.Execute(&buf, bundle{doc, fset}); err != nil {
			t.Error(err)
			continue
		}
		got := buf.Bytes()

		// update golden file if necessary
		golden := filepath.Join(dataDir, fmt.Sprintf("%s.%d.golden", pkg.Name, mode))
		if *update {
			err := ioutil.WriteFile(golden, got, 0644)
			if err != nil {
				t.Error(err)
			}
			continue
		}

		// get golden file
		want, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Error(err)
			continue
		}

		// compare
		if !bytes.Equal(got, want) {
			t.Errorf("package %s\n\tgot:\n%s\n\twant:\n%s", pkg.Name, got, want)
		}
	}
}

func Test(t *testing.T) {
	test(t, 0)
	test(t, AllDecls)
	test(t, AllMethods)
}