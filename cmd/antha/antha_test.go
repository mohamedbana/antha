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
// 1 Royal College St, London NW1 0NH UK

package main

import (
	"testing"
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

func TestGetPackageRoute(t *testing.T) {
	//TODO assert errors too
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %v.", err)
	}
	goPath := os.Getenv("GOPATH")
	//check that we are inside goPath
	if strings.HasPrefix(wd, goPath) { //TODO we should have another way of testing this, probably changing dir?
		test1, err := getPackageRoute(nil)
		if err != nil {
			t.Errorf("On getPackageRoute(nil) %v.", err)
		}
		if strings.HasPrefix(*test1, string(filepath.Separator)) {
			t.Errorf("On getPackageRoute(nil) got path '%s' starting with %c.", test1, filepath.Separator)
		}
		if strings.HasSuffix(*test1, string(filepath.Separator)) {
			t.Errorf("On getPackageRoute(nil) got path '%s' ending with %c.", test1, filepath.Separator)
		}
	}
	//Build a map for tests, must be dynamical because os.Getenv("GOPATH")...
	if strings.HasSuffix(goPath,string(filepath.Separator)) {
		goPath = strings.TrimRight(goPath, string(filepath.Separator))
	}
	fs := filepath.Separator
	testPaths := make(map[string]string)
	testPaths[fmt.Sprintf("%s%csrc%cgithub.com%cSynthace%cexamples%cadd%cadd.an", goPath,fs, fs, fs, fs, fs, fs)] =
		fmt.Sprintf("github.com%cSynthace%cexamples%cadd", fs, fs, fs)
	testPaths[fmt.Sprintf("%s%csrc%cgithub.com%cSynthace%cexamples%cadd%csum.an", goPath,fs, fs, fs, fs, fs, fs)] =
		fmt.Sprintf("github.com%cSynthace%cexamples%cadd", fs, fs, fs)
	testPaths[fmt.Sprintf("%s%csrc%cgithub.com%cSynthace%cexamples%csum%cadd.an", goPath,fs, fs, fs, fs, fs, fs)] =
		fmt.Sprintf("github.com%cSynthace%cexamples%csum", fs, fs, fs)
	for in, ex := range testPaths {
		res, err := getPackageRoute(&in)
		if err != nil {
			t.Errorf("On using %s. Got error %v.",in, err)
		}
		if ex != *res {
			t.Errorf("On using [%s]. Expecting [%s] got [%s].",in, ex, *res)
		}
	}
}
