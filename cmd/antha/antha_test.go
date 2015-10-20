// /antha/antha_test.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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

package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportPathFail(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		os.Chdir(wd)
	}()

	if err := os.Chdir(filepath.FromSlash("/")); err != nil {
		t.Error(err)
	}

	if _, err := getImportPath(""); err == nil {
		t.Error("expecting error for invalid working directory but succeeded instead")
	}
}

func TestImportPath(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	test1, err := getImportPath("")
	if err != nil {
		t.Error(err)
	}
	if strings.HasPrefix(test1, "/") {
		t.Errorf("got path '%s' starting with '/'.", test1)
	}
	if strings.HasSuffix(test1, "/") {
		t.Errorf("got path '%s' ending with '/'.", test1)
	}
	testPaths := make(map[string]string)
	testPaths[filepath.Join(goPath, "src", "github.com", "Synthace", "examples", "add")] =
		path.Join("github.com", "Synthace", "examples", "add")
	for in, ex := range testPaths {
		res, err := getImportPath(in)
		if err != nil {
			t.Errorf("On using %s. Got error %v.", in, err)
		}
		if ex != res {
			t.Errorf("On using [%s]. Expecting [%s] got [%s].", in, ex, res)
		}
	}
}
