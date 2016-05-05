// antha/AnthaStandardLibrary/Packages/enzymes/Annotatedseq.go: Part of the Antha language
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

// Package for interacting with and manipulating dna sequences in extension to methods available in wtype
package sequences

import (
	"fmt"
	//. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func ORFs2Features(orfs []ORF) (features []wtype.Feature) {

	features = make([]wtype.Feature, 0)

	for i, orf := range orfs {
		// currently just names each orf + number of orf. Add Rename orf function and sort by struct field function to run first to put orfs in order
		reverse := false
		if strings.ToUpper(orf.Direction) == strings.ToUpper("REVERSE") {
			reverse = true
		}
		feature := wtype.Feature{"orf" + strconv.Itoa(i), "orf", reverse, orf.StartPosition, orf.EndPosition, orf.DNASeq, orf.ProtSeq}
		features = append(features, feature)
	}
	return
}

func MakeFeature(name string, seq string, sequencetype string, class string, reverse string) (feature wtype.Feature) {

	feature.Name = name
	feature.DNASeq = strings.ToUpper(seq)

	//features := make([]Feature,0)
	//feature := Feature
	//fmt.Println("len seq =", len(seq))
	feature.Class = class
	if sequencetype == "aa" {
		feature.DNASeq = RevTranslatetoNstring(seq)
		feature.Protseq = seq
		feature.StartPosition = 1
		feature.EndPosition = len(feature.DNASeq)
		fmt.Println("len seq =", len(feature.DNASeq))
	} else {
		feature.DNASeq = seq
		feature.StartPosition = 1
		feature.EndPosition = len(seq)
		fmt.Println("len seq =", len(seq))
	}

	if reverse == "Reverse" {
		feature.Reverse = true
	}
	if feature.Class == "orf" {
		orf, orftrue := FindORF(seq)
		if orftrue == true {
			fmt.Println("orftrue!)")
			feature.Protseq = orf.ProtSeq
			feature.StartPosition = orf.StartPosition
			feature.EndPosition = orf.EndPosition
		}
	}
	return feature
}
