// antha/doc/headscan.go: Part of the Antha language
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

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

/*
	The headscan command extracts comment headings from package files;
	it is used to detect false positives which may require an adjustment
	to the comment formatting heuristics in comment.go.

	Usage: headscan [-root root_directory]

	By default, the $GOROOT/src directory is scanned.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/antha-lang/antha/antha/doc"
	"github.com/antha-lang/antha/antha/parser"
	"github.com/antha-lang/antha/antha/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	root    = flag.String("root", filepath.Join(runtime.GOROOT(), "src"), "root of filesystem tree to scan")
	verbose = flag.Bool("v", false, "verbose mode")
)

const (
	html_h    = "<h3>"
	html_endh = "</h3>\n"
)

func isGoFile(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), ".go") &&
		!strings.HasSuffix(fi.Name(), "_test.go")
}

func appendHeadings(list []string, comment string) []string {
	var buf bytes.Buffer
	doc.ToHTML(&buf, comment, nil)
	for s := buf.String(); ; {
		i := strings.Index(s, html_h)
		if i < 0 {
			break
		}
		i += len(html_h)
		j := strings.Index(s, html_endh)
		if j < 0 {
			list = append(list, s[i:]) // incorrect HTML
			break
		}
		list = append(list, s[i:j])
		s = s[j+len(html_endh):]
	}
	return list
}

func main() {
	flag.Parse()
	fset := token.NewFileSet()
	nheadings := 0
	err := filepath.Walk(*root, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			return nil
		}
		pkgs, err := parser.ParseDir(fset, path, isGoFile, parser.ParseComments)
		if err != nil {
			if *verbose {
				fmt.Fprintln(os.Stderr, err)
			}
			return nil
		}
		for _, pkg := range pkgs {
			d := doc.New(pkg, path, doc.Mode(0))
			list := appendHeadings(nil, d.Doc)
			for _, d := range d.Consts {
				list = appendHeadings(list, d.Doc)
			}
			for _, d := range d.Types {
				list = appendHeadings(list, d.Doc)
			}
			for _, d := range d.Vars {
				list = appendHeadings(list, d.Doc)
			}
			for _, d := range d.Funcs {
				list = appendHeadings(list, d.Doc)
			}
			if len(list) > 0 {
				// directories may contain multiple packages;
				// print path and package name
				fmt.Printf("%s (package %s)\n", path, pkg.Name)
				for _, h := range list {
					fmt.Printf("\t%s\n", h)
				}
				nheadings += len(list)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(nheadings, "headings found")
}
