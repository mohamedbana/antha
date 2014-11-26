// antha/cmd/anthafmt/long_test.go: Part of the Antha language
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

// This test applies gofmt to all Go/Antha files under -root.
// To test specific files provide a list of comma-separated
// filenames via the -files flag: antha test -files=gofmt.go .

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/printer"
	"github.com/antha-lang/antha/token"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	root    = flag.String("root", runtime.GOROOT(), "test root directory")
	files   = flag.String("files", "", "comma-separated list of files to test")
	ngo     = flag.Int("n", runtime.NumCPU(), "number of goroutines used")
	verbose = flag.Bool("verbose", false, "verbose mode")
	nfiles  int // number of files processed
)

func gofmt(fset *token.FileSet, filename string, src *bytes.Buffer) error {
	f, _, err := parse(fset, filename, src.Bytes(), false)
	if err != nil {
		return err
	}
	ast.SortImports(fset, f)
	src.Reset()
	return (&printer.Config{Mode: printerMode, Tabwidth: tabWidth}).Fprint(src, fset, f)
}

func testFile(t *testing.T, b1, b2 *bytes.Buffer, filename string) {
	// open file
	f, err := os.Open(filename)
	if err != nil {
		t.Error(err)
		return
	}

	// read file
	b1.Reset()
	_, err = io.Copy(b1, f)
	f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// exclude files w/ syntax errors (typically test cases)
	fset := token.NewFileSet()
	if _, _, err = parse(fset, filename, b1.Bytes(), false); err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "ignoring %s\n", err)
		}
		return
	}

	// gofmt file
	if err = gofmt(fset, filename, b1); err != nil {
		t.Errorf("1st gofmt failed: %v", err)
		return
	}

	// make a copy of the result
	b2.Reset()
	b2.Write(b1.Bytes())

	// gofmt result again
	if err = gofmt(fset, filename, b2); err != nil {
		t.Errorf("2nd gofmt failed: %v", err)
		return
	}

	// the first and 2nd result should be identical
	if !bytes.Equal(b1.Bytes(), b2.Bytes()) {
		t.Errorf("gofmt %s not idempotent", filename)
	}
}

func testFiles(t *testing.T, filenames <-chan string, done chan<- int) {
	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	for filename := range filenames {
		testFile(t, b1, b2, filename)
	}
	done <- 0
}

func genFilenames(t *testing.T, filenames chan<- string) {
	defer close(filenames)

	handleFile := func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		if isGoFile(fi) {
			filenames <- filename
			nfiles++
		}
		return nil
	}

	// test Antha files provided via -files, if any
	if *files != "" {
		for _, filename := range strings.Split(*files, ",") {
			fi, err := os.Stat(filename)
			handleFile(filename, fi, err)
		}
		return // ignore files under -root
	}

	// otherwise, test all Antha files under *root
	filepath.Walk(*root, handleFile)
}

func TestAll(t *testing.T) {
	if testing.Short() {
		return
	}

	if *ngo < 1 {
		*ngo = 1 // make sure test is run
	}
	if *verbose {
		fmt.Printf("running test using %d goroutines\n", *ngo)
	}

	// generate filenames
	filenames := make(chan string, 32)
	go genFilenames(t, filenames)

	// launch test goroutines
	done := make(chan int)
	for i := 0; i < *ngo; i++ {
		go testFiles(t, filenames, done)
	}

	// wait for all test goroutines to complete
	for i := 0; i < *ngo; i++ {
		<-done
	}

	if *verbose {
		fmt.Printf("processed %d files\n", nfiles)
	}
}
