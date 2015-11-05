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

package enzymes

import (
	"fmt"
	"strings"
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

type AnnotatedSeq struct {
	Nm       string
	Seq      string
	features []Feature
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

func ConcatenateFeatures(name string, featuresinorder []Feature) (annotated AnnotatedSeq) {

	annotated.Nm = name
	annotated.Seq = featuresinorder[0].DNASeq
	annotated.features = make([]Feature, 0)
	annotated.features = append(annotated.features, featuresinorder[0])
	for i := 1; i < len(featuresinorder); i++ {
		nextfeature := featuresinorder[i]
		nextfeature.StartPosition = nextfeature.StartPosition + annotated.features[i-1].EndPosition
		nextfeature.EndPosition = nextfeature.EndPosition + annotated.features[i-1].EndPosition
		annotated.Seq = annotated.Seq + featuresinorder[i].DNASeq
		annotated.features = append(annotated.features, nextfeature)
	}
	return annotated
}
