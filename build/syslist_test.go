// antha/build/syslist_test.go: Part of the Antha language
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


package build

import (
	"runtime"
	"testing"
)

var (
	thisOS    = runtime.GOOS
	thisArch  = runtime.GOARCH
	otherOS   = anotherOS()
	otherArch = anotherArch()
)

func anotherOS() string {
	if thisOS != "darwin" {
		return "darwin"
	}
	return "linux"
}

func anotherArch() string {
	if thisArch != "amd64" {
		return "amd64"
	}
	return "386"
}

type GoodFileTest struct {
	name   string
	result bool
}

var tests = []GoodFileTest{
	{"file.go", true},
	{"file.c", true},
	{"file_foo.go", true},
	{"file_" + thisArch + ".go", true},
	{"file_" + otherArch + ".go", false},
	{"file_" + thisOS + ".go", true},
	{"file_" + otherOS + ".go", false},
	{"file_" + thisOS + "_" + thisArch + ".go", true},
	{"file_" + otherOS + "_" + thisArch + ".go", false},
	{"file_" + thisOS + "_" + otherArch + ".go", false},
	{"file_" + otherOS + "_" + otherArch + ".go", false},
	{"file_foo_" + thisArch + ".go", true},
	{"file_foo_" + otherArch + ".go", false},
	{"file_" + thisOS + ".c", true},
	{"file_" + otherOS + ".c", false},
}

func TestGoodOSArch(t *testing.T) {
	for _, test := range tests {
		if Default.goodOSArchFile(test.name, make(map[string]bool)) != test.result {
			t.Fatalf("goodOSArchFile(%q) != %v", test.name, test.result)
		}
	}
}