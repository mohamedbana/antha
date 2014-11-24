// antha/doc/testdata/example.go: Part of the Antha language
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


package testing

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type InternalExample struct {
	Name   string
	F      func()
	Output string
}

func RunExamples(examples []InternalExample) (ok bool) {
	ok = true

	var eg InternalExample

	stdout, stderr := os.Stdout, os.Stderr
	defer func() {
		os.Stdout, os.Stderr = stdout, stderr
		if e := recover(); e != nil {
			fmt.Printf("--- FAIL: %s\npanic: %v\n", eg.Name, e)
			os.Exit(1)
		}
	}()

	for _, eg = range examples {
		if *chatty {
			fmt.Printf("=== RUN: %s\n", eg.Name)
		}

		// capture stdout and stderr
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Stdout, os.Stderr = w, w
		outC := make(chan string)
		go func() {
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, r)
			if err != nil {
				fmt.Fprintf(stderr, "testing: copying pipe: %v\n", err)
				os.Exit(1)
			}
			outC <- buf.String()
		}()

		// run example
		t0 := time.Now()
		eg.F()
		dt := time.Now().Sub(t0)

		// close pipe, restore stdout/stderr, get output
		w.Close()
		os.Stdout, os.Stderr = stdout, stderr
		out := <-outC

		// report any errors
		tstr := fmt.Sprintf("(%.2f seconds)", dt.Seconds())
		if g, e := strings.TrimSpace(out), strings.TrimSpace(eg.Output); g != e {
			fmt.Printf("--- FAIL: %s %s\ngot:\n%s\nwant:\n%s\n",
				eg.Name, tstr, g, e)
			ok = false
		} else if *chatty {
			fmt.Printf("--- PASS: %s %s\n", eg.Name, tstr)
		}
	}

	return
}