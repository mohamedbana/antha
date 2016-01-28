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

/*
const (
	orf = iota
	Promoter
	Ribosomebindingsite
	TranslationInitSite
	Origin
	Marker
	Misc
)
*/
type Feature struct {
	Name          string
	Class         string //int // defined by constants above
	Reverse       bool
	StartPosition int
	EndPosition   int
	DNASeq        string
	Protseq       string
}

func (feat *Feature) Coordinates() (pair []int) {
	pair[0] = feat.StartPosition
	pair[1] = feat.EndPosition
	return
}

type AnnotatedSeq struct {
	//Nm string
	wtype.DNASequence
	Features []Feature
}

func Annotate(dnaseq wtype.DNASequence, features []Feature) (annotated AnnotatedSeq) {
	annotated.DNASequence = dnaseq
	annotated.Features = features
	return
}

func AddFeatures(annotated AnnotatedSeq, features []Feature) (updated AnnotatedSeq) {

	for _, feature := range features {
		annotated.Features = append(annotated.Features, feature)
	}
	return
}

func ORFs2Features(orfs []ORF) (features []Feature) {

	features = make([]Feature, 0)

	for i, orf := range orfs {
		// currently just names each orf + number of orf. Add Rename orf function and sort by struct field function to run first to put orfs in order
		reverse := false
		if strings.ToUpper(orf.Direction) == strings.ToUpper("REVERSE") {
			reverse = true
		}
		feature := Feature{"orf" + strconv.Itoa(i), "orf", reverse, orf.StartPosition, orf.EndPosition, orf.DNASeq, orf.ProtSeq}
		features = append(features, feature)
	}
	return
}

func MakeFeature(name string, seq string, sequencetype string, class string, reverse string) (feature Feature) {

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

func MakeAnnotatedSeq(name string, seq string, circular bool, features []Feature) (annotated AnnotatedSeq, err error) {
	annotated.Nm = name
	//annotated.Seq.Nm = name
	annotated.Seq = seq //.Seq = seq
	annotated.Plasmid = circular

	for _, feature := range features {
		if strings.Contains(seq, feature.DNASeq) {
			feature.StartPosition = strings.Index(seq, feature.DNASeq)
			feature.EndPosition = feature.EndPosition + feature.StartPosition
		} else if strings.Contains(seq, RevComp(feature.DNASeq)) {
			feature.StartPosition = strings.Index(seq, feature.DNASeq)
			feature.EndPosition = feature.EndPosition + feature.StartPosition
			err = fmt.Errorf(feature.Name, " Feature only found in reverse direction")
		} else {
			err = fmt.Errorf(feature.Name, " not found in sequence")
		}
	}
	annotated.Features = features
	return
}

func ConcatenateFeatures(name string, featuresinorder []Feature) (annotated AnnotatedSeq) {

	annotated.Nm = name
	//annotated.Seq.Nm = name
	annotated.Seq = featuresinorder[0].DNASeq
	annotated.Features = make([]Feature, 0)
	annotated.Features = append(annotated.Features, featuresinorder[0])
	for i := 1; i < len(featuresinorder); i++ {
		nextfeature := featuresinorder[i]
		nextfeature.StartPosition = nextfeature.StartPosition + annotated.Features[i-1].EndPosition
		nextfeature.EndPosition = nextfeature.EndPosition + annotated.Features[i-1].EndPosition
		annotated.Seq = annotated.Seq + featuresinorder[i].DNASeq
		annotated.Features = append(annotated.Features, nextfeature)
	}
	return annotated
}
